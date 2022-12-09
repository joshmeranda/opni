package planner

import (
	"context"
	"fmt"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/prompb"
	"github.com/rancher/opni/pkg/migrate"
	"math"
	"sync"
)

// todo: should be merged into panner
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
	minTimestamp int64 // inclusive.
	maxTimestamp int64 // exclusive.
	done         bool
	timeseries   []*prompb.TimeSeries
	stores       []store
	numStores    int
	plan         *Plan
	percentDone  float64
}

func NewSlab(plan *Plan, minTimestamp int64, maxTimestamp int64) (*Slab, error) {
	if plan.config.MinTimestamp > minTimestamp || plan.config.MaxTimestamp < minTimestamp {
		return nil, fmt.Errorf("minTimestamp is not between global min and max timestamps")
	} else if plan.config.MinTimestamp > maxTimestamp || plan.config.MaxTimestamp < maxTimestamp {
		return nil, fmt.Errorf("maxTimestamp is not between global min and max timestamps")
	} else if minTimestamp > maxTimestamp {
		return nil, fmt.Errorf("minTimestamp cannot be > then maxTimestamp")
	}

	slab := slabPool.Get().(*Slab)

	slab.percentDone = math.Min(100, float64(maxTimestamp-plan.config.MinTimestamp)*100/float64(plan.config.MaxTimestamp-plan.config.MinTimestamp))
	slab.numStores = plan.config.NumStores
	slab.minTimestamp = minTimestamp
	slab.maxTimestamp = maxTimestamp
	slab.plan = plan

	if cap(slab.stores) < slab.numStores {
		slab.stores = make([]store, slab.numStores)
	} else {
		slab.stores = slab.stores[:slab.numStores]
	}

	slab.initStores()

	return slab, nil
}

// initStores initializes the stores.
func (slab *Slab) initStores() {
	proceed := slab.minTimestamp
	increment := (slab.maxTimestamp - slab.minTimestamp) / int64(slab.numStores)

	for storeIndex := 0; storeIndex < slab.numStores; storeIndex++ {
		slab.stores[storeIndex] = store{
			id:           storeIndex,
			minTimestamp: proceed,
			maxTimestamp: proceed + increment,
		}

		proceed += increment
	}

	// To ensure that we cover the entire range, which may get missed due to integer division in increment, we update
	// the maxTimestamp of the last store to be the maxTimestamp of the slab
	slab.stores[slab.numStores-1].maxTimestamp = slab.maxTimestamp
}

// Fetch starts fetching the samples from remote read storage based on the matchers. It takes care of concurrent pulls as well.
func (slab *Slab) Fetch(ctx context.Context, client *migrate.Client, matchers []*labels.Matcher) (err error) {
	totalRequests := slab.numStores
	cancelFuncs := make([]context.CancelFunc, totalRequests)
	responseChan := make(chan interface{}, totalRequests)

	for i := 0; i < totalRequests; i++ {
		readRequest, err := migrate.CreatePrombQuery(slab.stores[i].minTimestamp, slab.stores[i].maxTimestamp, matchers)
		if err != nil {
			return fmt.Errorf("could not create promb query: %w", err)
		}

		cctx, cancelFunc := context.WithCancel(ctx)
		cancelFuncs[i] = cancelFunc

		go client.ReadConcurrent(cctx, readRequest, i, responseChan)
	}

	var bytesCompressed int
	var bytesUncompressed int
	pendingResponses := make([]*migrate.PrompbResponse, totalRequests)

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
			return fmt.Errorf("recevied error during read: %w", response)
		default:
			panic("invalid response type")
		}
	}

	close(responseChan)

	if totalRequests > 1 {
		ts, err := slab.mergeResponses(pendingResponses)
		if err != nil {
			return fmt.Errorf("could not merge responses into single slab: %w", err)
		}
		slab.timeseries = ts
	} else {
		if l := len(pendingResponses); l > 1 || l == 0 {
			return fmt.Errorf("fatal: expected a single pending request when totalRequest is 1. Received pending %d requests", l)
		}

		slab.timeseries = pendingResponses[0].Result.Timeseries
	}

	slab.plan.update(bytesUncompressed, slab.percentDone)

	return nil
}

func (slab *Slab) mergeResponses(subSlabs []*migrate.PrompbResponse) ([]*prompb.TimeSeries, error) {
	timeseries := make(map[string]*prompb.TimeSeries) // todo: could probably be a slice

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

	series := make([]*prompb.TimeSeries, len(timeseries))
	i := 0

	for _, ts := range timeseries {
		series[i] = ts
		i++
	}

	return series, nil
}

// Series returns the time-series in the slab.
func (slab *Slab) Series() []*prompb.TimeSeries {
	return slab.timeseries
}

// Done updates the text and sets the spinner to done.
func (slab *Slab) Done() error {
	slab.done = true
	return nil
}

// IsEmpty returns true if the slab does not contain any time-series.
func (slab *Slab) IsEmpty() bool {
	return len(slab.timeseries) == 0
}
