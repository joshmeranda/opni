package writer

import (
	"context"
	"fmt"
	"github.com/rancher/opni/pkg/migrate"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/snappy"
	"github.com/prometheus/common/config"
	"github.com/prometheus/prometheus/prompb"
	"github.com/prometheus/prometheus/storage/remote"
)

// shard is an atomic unit of writing. Each shard accepts a sharded time-series set and
// pushes the respective time-series via a POST request. More shards will lead to bigger
// throughput if the target system can handle it.
type shard struct {
	ctx          context.Context
	client       *migrate.Client
	queue        chan *[]prompb.TimeSeries
	copyShardSet *shardsSet
}

func (shard *shard) run(shardIndex int) {
	defer func() {
		lg.Infof("shard-%d is inactive", shardIndex)
	}()
	lg.Infof("shard-%d is active", shardIndex)

	for {
		select {
		case <-shard.ctx.Done():
			return
		case timeSeries, ok := <-shard.queue:
			if !ok {
				return
			}

			if len(*timeSeries) != 0 {
				if err := shard.sendSamplesWithBackoff(timeSeries); err != nil {
					shard.copyShardSet.errChan <- err
					return
				}
			} else {
				lg.Warn("found empty time-series shard shard %d", shardIndex)
			}

			shard.copyShardSet.errChan <- nil
		}
	}
}

// sendSamplesWithBackoff to the remote storage with backoff for recoverable errors.
func (shard *shard) sendSamplesWithBackoff(samples *[]prompb.TimeSeries) error {
	buf := bytePool.Get().(*[]byte)
	defer func() {
		*buf = (*buf)[:0]
		bytePool.Put(buf)
	}()

	req, err := buildWriteRequest(*samples, *buf)
	if err != nil {
		return err
	}

	backoff := backOffRetryDuration
	nonRecvRetries := 0
	*buf = req

	for {
		select {
		case <-shard.ctx.Done():
			return shard.ctx.Err()
		default:
		}

		if err := shard.client.Write(shard.ctx, *buf); err != nil {
			if _, ok := err.(remote.RecoverableError); !ok {
				switch r := shard.client.Config(); r.OnErr {
				case migrate.Retry:
					if r.MaxRetry != 0 {
						if nonRecvRetries >= r.MaxRetry {
							return fmt.Errorf("exceeded retrying limit in non-recoverable error. Aborting")
						}

						nonRecvRetries++
					}

					lg.Warn("received non-recoverable error, retrying again after delay", r.RetryDelay)
					time.Sleep(r.RetryDelay)

					continue
				case migrate.Skip:
					lg.Warn("received non-recoverable error, skipping current slab")
					return nil
				case migrate.Abort:
				}
				return err
			}

			time.Sleep(backoff)

			continue
		}

		return nil
	}
}

func buildWriteRequest(samples []prompb.TimeSeries, buf []byte) ([]byte, error) {
	req := &prompb.WriteRequest{
		Timeseries: samples,
	}
	data, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}

	// snappy uses len() to see if it needs to allocate a new slice. Make the
	// buffer as long as possible.
	//if buf != nil {
	//	buf = buf[0:cap(buf)]
	//}

	compressed := snappy.Encode(buf, data)
	return compressed, nil
}

// shardsSet represents a set of shards. It consists of configuration related to handling and management of shards.
type shardsSet struct {
	num         int
	set         []*shard
	cancelFuncs []context.CancelFunc

	errChan chan error
}

func newShardsSet(writerCtx context.Context, httpConfig config.HTTPClientConfig, clientConfig migrate.ClientConfig, numShards int) (*shardsSet, error) {
	set := make([]*shard, numShards)
	cancelFuncs := make([]context.CancelFunc, numShards)

	ss := &shardsSet{
		num:     numShards,
		errChan: make(chan error, numShards),
	}

	for i := 0; i < numShards; i++ {
		client, err := migrate.NewClient(fmt.Sprintf("writer-shard-%d", i), clientConfig, httpConfig)
		if err != nil {
			return nil, fmt.Errorf("could not write-shard-client-%d: %w", i, err)
		}

		ctx, cancelFunc := context.WithCancel(writerCtx)
		shard := &shard{
			ctx:          ctx,
			client:       client,
			queue:        make(chan *[]prompb.TimeSeries),
			copyShardSet: ss,
		}

		cancelFuncs[i] = cancelFunc
		set[i] = shard
	}

	ss.set = set
	ss.cancelFuncs = cancelFuncs

	// Run the shards.
	for i := 0; i < ss.num; i++ {
		go ss.set[i].run(i)
	}

	return ss, nil
}

// scheduleTS batches the time-series based on the available shards.
func (s *shardsSet) scheduleTS(ts []prompb.TimeSeries) (numSigExpected int) {
	var batches = make([][]prompb.TimeSeries, s.num)
	for _, series := range ts {
		hash := migrate.HashLabels(labelsSliceToLabels(series.GetLabels()))
		batchIndex := hash % uint64(s.num)
		batches[batchIndex] = append(batches[batchIndex], series)
	}

	// Feed to the shards.
	for shardIndex := 0; shardIndex < s.num; shardIndex++ {
		batchIndex := shardIndex
		if len(batches[batchIndex]) == 0 {
			// We do not want "wait state" to happen on the shards.
			continue
		}

		numSigExpected++
		s.set[shardIndex].queue <- &batches[batchIndex]
	}

	return
}

const backOffRetryDuration = time.Second * 1

var bytePool = sync.Pool{New: func() interface{} { return new([]byte) }}

func labelsSliceToLabels(ls []prompb.Label) prompb.Labels {
	var labels prompb.Labels

	labels.Labels = ls

	return labels
}
