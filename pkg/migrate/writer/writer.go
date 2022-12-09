package writer

import (
	"context"
	"fmt"
	"github.com/rancher/opni/pkg/logger"
	"github.com/rancher/opni/pkg/migrate"
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

func timeseriesRefToTimeseries(tsr []*prompb.TimeSeries) (ts []prompb.TimeSeries) {
	ts = make([]prompb.TimeSeries, len(tsr))

	for i := range tsr {
		ts[i] = *tsr[i]
	}

	return
}

// Config is config for writer.
type Config struct {
	Context          context.Context
	MigrationJobName string // Label value to the progress metric.
	ConcurrentPush   int
	ClientConfig     migrate.ClientConfig
	HTTPConfig       config.HTTPClientConfig

	GarbageCollectOnPush bool

	//nolint
	signalGarbageCollector chan struct{}

	SlabChan chan *planner.Slab
}

type Writer struct {
	Config
	shardsSet   *shardsSet
	slabsPushed int64
}

func NewWriter(config Config) (*Writer, error) {
	ss, err := newShardsSet(config.Context, config.HTTPConfig, config.ClientConfig, config.ConcurrentPush)
	if err != nil {
		return nil, fmt.Errorf("error creating shards: %w", err)
	}

	write := &Writer{
		Config:    config,
		shardsSet: ss,
	}

	if config.GarbageCollectOnPush {
		write.signalGarbageCollector = make(chan struct{}, 1)
	}

	return write, nil
}

func (w *Writer) Run(errChan chan<- error) {
	var err error
	shards := w.shardsSet

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
			case slab, ok := <-w.SlabChan:
				if !ok {
					return
				}

				numSigExpected := shards.scheduleTS(timeseriesRefToTimeseries(slab.Series()))
				if isErrSig(slab, numSigExpected) {
					return
				}

				atomic.AddInt64(&w.slabsPushed, 1)
				if err = slab.Done(); err != nil {
					errChan <- fmt.Errorf("remote-write run: %w", err)
					return
				}

				planner.PutSlab(slab)
				w.collectGarbage()
			}
		}
	}()
}

func (w *Writer) collectGarbage() {
	if w.signalGarbageCollector == nil {
		// signalGarbageCollector is nil if w.GarbageCollectOnPush is set to false.
		return
	}

	gcRunner.Do(func() {
		go gc(w.signalGarbageCollector)
	})

	select {
	case w.signalGarbageCollector <- struct{}{}:
	default:
		// Skip GC as its already running.
	}
}

func (w *Writer) pushProgressMetric(series *prompb.TimeSeries) int {
	shards := w.shardsSet
	return shards.scheduleTS([]prompb.TimeSeries{*series})
}
