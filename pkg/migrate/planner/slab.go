package planner

import (
	"context"
	"fmt"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/prompb"
	"github.com/rancher/opni/pkg/migrate"
	"github.com/rancher/opni/pkg/migrate/utils"
	"sync"
)

var slabPool = sync.Pool{New: func() interface{} { return new(Slab) }}

func PutSlab(s *Slab) {
	s.stores = s.stores[:0]
	slabPool.Put(s)
}

type store struct {
	id           int
	minTimestamp int64
	maxTimestamp int64
}

// Slab represents an in-memory storage for data that is fetched by the reader.
type Slab struct {
	minTimestamp         int64 // inclusive.
	maxTimestamp         int64 // exclusive.
	done                 bool
	timeseries           []*prompb.TimeSeries
	stores               []store
	numStores            int
	numBytesCompressed   int
	numBytesUncompressed int
	plan                 *Plan // We keep a copy of plan so that each slab has the authority to update the stats of the planner.
	percentDone          float64
}

// initStores initializes the stores.
func (s *Slab) initStores() {
	var (
		proceed   = s.minTimestamp
		increment = (s.maxTimestamp - s.minTimestamp) / int64(s.numStores)
	)

	for storeIndex := 0; storeIndex < s.numStores; storeIndex++ {
		s.stores[storeIndex] = store{
			id:           storeIndex,
			minTimestamp: proceed,
			maxTimestamp: proceed + increment,
		}

		proceed += increment
	}

	// To ensure that we cover the entire range, which may get missed due to integer division in increment, we update
	// the maxTimestamp of the last store  to be the maxTimestamp of the slab.
	s.stores[s.numStores-1].maxTimestamp = s.maxTimestamp
}

// Fetch starts fetching the samples from remote read storage based on the matchers. It takes care of concurrent pulls as well.
func (s *Slab) Fetch(ctx context.Context, client *utils.Client, mint, maxt int64, matchers []*labels.Matcher) (err error) {
	var (
		totalRequests = s.numStores
		cancelFuncs   = make([]context.CancelFunc, totalRequests)
		responseChan  = make(chan interface{}, totalRequests)
	)

	for i := 0; i < totalRequests; i++ {
		readRequest, err := migrate.CreatePrombQuery(s.stores[i].minTimestamp, s.stores[i].maxTimestamp, matchers)
		if err != nil {
			return fmt.Errorf("could not create promb query: %w", err)
		}

		cctx, cancelFunc := context.WithCancel(ctx)
		cancelFuncs[i] = cancelFunc

		go client.ReadConcurrent(cctx, readRequest, i, responseChan)
	}

	var (
		bytesCompressed   int
		bytesUncompressed int
		// pendingResponses is used to ensure that concurrent fetch results are appending in sorted order of time.
		pendingResponses = make([]*migrate.PrompbResponse, totalRequests)
	)

	for i := 0; i < totalRequests; i++ {
		resp := <-responseChan

		switch response := resp.(type) {
		case *migrate.PrompbResponse:
			bytesCompressed += response.NumBytesCompressed
			bytesUncompressed += response.NumBytesUncompressed

			pendingResponses[response.ID] = response
		case error:
			for _, cancelFnc := range cancelFuncs {
				cancelFnc()
			}
			return fmt.Errorf("executing client-read: %w", response)
		default:
			panic("invalid response type")
		}
	}

	close(responseChan)

	if totalRequests > 1 {
		ts, err := s.mergeSubSlabsToSlab(pendingResponses)
		if err != nil {
			return fmt.Errorf("could not merge subSlabs into a slab: %w", err)
		}
		s.timeseries = ts
	} else {
		// Short path optimization. When not fetching concurrently, we can skip from memory allocations.
		if l := len(pendingResponses); l > 1 || l == 0 {
			return fmt.Errorf("fatal: expected a single pending request when totalRequest is 1. Received pending %d requests", l)
		}

		s.timeseries = pendingResponses[0].Result.Timeseries
	}

	// We set compressed bytes in slab since those are the bytes that will be pushed over the network to the write
	// storage after snappy compression. The pushed bytes are not exactly the bytesCompressed since while pushing, we
	// add the progress metric. But, the size of progress metric along with the sample is negligible. So, it is safe to
	// consider bytesCompressed in such a scenario.
	s.numBytesCompressed = bytesCompressed
	s.numBytesUncompressed = bytesUncompressed
	s.plan.update(bytesUncompressed)

	// todo: these clog stdout we might want to limit output
	lg.Infof("read %%%f done", s.percentDone)

	return nil
}

func (s *Slab) mergeSubSlabsToSlab(subSlabs []*migrate.PrompbResponse) ([]*prompb.TimeSeries, error) {
	timeseries := make(map[string]*prompb.TimeSeries)

	for i := 0; i < len(subSlabs); i++ {
		ts := subSlabs[i].Result.Timeseries

		for _, series := range ts {
			labelsStr := migrate.GetLabelsStr(series.Labels)

			if s, ok := timeseries[labelsStr]; len(series.Samples) > 0 && ok {
				l := len(s.Samples)

				if len(s.Samples) > 1 && s.Samples[l-1].Timestamp == series.Samples[0].Timestamp {
					return nil, fmt.Errorf("overlapping samples detected when fetching data from read storage")
				}

				s.Samples = append(s.Samples, series.Samples...)
			} else {
				timeseries[labelsStr] = series
			}
		}
	}

	// Form series slice after combining the concurrent responses.
	series := make([]*prompb.TimeSeries, len(timeseries))
	i := 0

	for _, ts := range timeseries {
		series[i] = ts
		i++
	}

	return series, nil
}

// Series returns the time-series in the slab.
func (s *Slab) Series() []*prompb.TimeSeries {
	return s.timeseries
}

// UpdateProgressSeries returns a time-series after appending a sample to the progress-metric.
func (s *Slab) UpdateProgressSeries(ts *prompb.TimeSeries) *prompb.TimeSeries {
	ts.Samples = []prompb.Sample{{Timestamp: s.MaxTimestamp(), Value: 1}} // One sample per slab.
	return ts
}

// Done updates the text and sets the spinner to done.
func (s *Slab) Done() error {
	s.done = true
	return nil
}

// IsEmpty returns true if the slab does not contain any time-series.
func (s *Slab) IsEmpty() bool {
	return len(s.timeseries) == 0
}

// MinTimestamp returns the minTimestamp of the slab (inclusive).
func (s *Slab) MinTimestamp() int64 {
	return s.minTimestamp
}

// MaxTimestamp returns the maxTimestamp of the slab (exclusive).
func (s *Slab) MaxTimestamp() int64 {
	return s.maxTimestamp
}
