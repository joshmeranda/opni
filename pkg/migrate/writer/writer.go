// This file and its contents are licensed under the Apache License 2.0.
// Please see the included NOTICE for copyright information and
// LICENSE for a copy of the license.

package writer

import (
	"context"
	"fmt"
	"github.com/rancher/opni/pkg/logger"
	"github.com/rancher/opni/pkg/migrate"
	"github.com/rancher/opni/pkg/migrate/utils"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/prometheus/common/config"
	"github.com/prometheus/prometheus/prompb"

	"github.com/rancher/opni/pkg/migrate/planner"
)

var (
	gcRunner sync.Once
	lg       = logger.New().Named("writer")
)

func gc(collect <-chan struct{}) {
	for range collect {
		runtime.GC() // This is blocking, hence we run asynchronously.
	}
}

// Config is config for writer.
type Config struct {
	Context          context.Context
	MigrationJobName string // Label value to the progress metric.
	ConcurrentPush   int
	ClientConfig     utils.ClientConfig
	HTTPConfig       config.HTTPClientConfig

	ProgressEnabled    bool
	ProgressMetricName string

	GarbageCollectOnPush bool

	//nolint
	sigGC chan struct{}

	SigSlabRead chan *planner.Slab
	SigSlabStop chan struct{}
}

type Writer struct {
	Config
	shardsSet          *shardsSet
	slabsPushed        int64
	progressTimeSeries *prompb.TimeSeries
}

// NewWriter returns a new remote write. It is responsible for writing to the remote write storage.
func NewWriter(config Config) (*Writer, error) {
	ss, err := newShardsSet(config.Context, config.HTTPConfig, config.ClientConfig, config.ConcurrentPush)
	if err != nil {
		return nil, fmt.Errorf("creating shards: %w", err)
	}
	write := &Writer{
		Config:    config,
		shardsSet: ss,
	}
	if config.ProgressEnabled {
		write.progressTimeSeries = &prompb.TimeSeries{
			Labels: migrate.LabelSet(config.ProgressMetricName, config.MigrationJobName),
		}
	}
	if config.GarbageCollectOnPush {
		write.sigGC = make(chan struct{}, 1)
	}
	return write, nil
}

// Run runs the remote-writer. It waits for the remote-reader to give access to the in-memory
// data-block that is written after the most recent fetch. After reading the block and storing
// the data locally, it gives back the writing access to the remote-reader for further fetches.
func (w *Writer) Run(errChan chan<- error) {
	var (
		err    error
		shards = w.shardsSet
	)

	go func() {
		defer func() {
			for _, cancelFunc := range shards.cancelFuncs {
				cancelFunc()
			}

			lg.Info("writer is down")
			close(errChan)
		}()

		lg.Info("writer is up")
		isErrSig := func(ref *planner.Slab, expectedSigs int) bool {
			for i := 0; i < expectedSigs; i++ {
				if err := <-shards.errChan; err != nil {
					errChan <- fmt.Errorf("remote-write run: %w", err)
					return true
				}
			}

			return false
		}

		for {
			select {
			case <-w.Context.Done():
				return
			case slabRef, ok := <-w.SigSlabRead:
				if !ok {
					return
				}

				numSigExpected := shards.scheduleTS(timeseriesRefToTimeseries(slabRef.Series()))
				if isErrSig(slabRef, numSigExpected) {
					return
				}

				if w.progressTimeSeries != nil {
					numSigExpected := w.pushProgressMetric(slabRef.UpdateProgressSeries(w.progressTimeSeries))
					if isErrSig(slabRef, numSigExpected) {
						return
					}
				}

				atomic.AddInt64(&w.slabsPushed, 1)
				if err = slabRef.Done(); err != nil {
					errChan <- fmt.Errorf("remote-write run: %w", err)
					return
				}

				planner.PutSlab(slabRef)
				w.collectGarbage()
			}
		}
	}()
}

func (w *Writer) collectGarbage() {
	if w.sigGC == nil {
		// sigGC is nil if w.GarbageCollectOnPush is set to false.
		return
	}

	gcRunner.Do(func() {
		go gc(w.sigGC)
	})

	select {
	case w.sigGC <- struct{}{}:
	default:
		// Skip GC as its already running.
	}
}

func (w *Writer) pushProgressMetric(series *prompb.TimeSeries) int {
	shards := w.shardsSet
	return shards.scheduleTS([]prompb.TimeSeries{*series})
}
