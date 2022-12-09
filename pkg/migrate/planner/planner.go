package planner

import (
	"fmt"
	"github.com/rancher/opni/pkg/logger"
	"go.uber.org/zap/zapcore"
	"sync/atomic"
	"time"
)

var (
	lg = logger.New(logger.WithLogLevel(zapcore.InfoLevel)).Named("planner")
)

// Config represents configuration for the planner.
type Config struct {
	MinTimestamp       int64
	MaxTimestamp       int64
	SlabSizeLimitBytes int64
	NumStores          int
	JobName            string
	LaIncrement        time.Duration
	MaxReadDuration    time.Duration
}

type Plan struct {
	config *Config

	nextMinTimestamp   int64
	lastNumBytes       int64
	lastTimeRangeDelta int64
	deltaIncRegion     int64

	percentDone          float64
	nextPercentThreshold float64
}

func NewPlan(config *Config) (*Plan, error) {
	if config.MinTimestamp >= config.MaxTimestamp {
		return nil, fmt.Errorf("min timestamp cannot be >= max timestamp")
	}

	return &Plan{
		config:           config,
		nextMinTimestamp: config.MinTimestamp,
		deltaIncRegion:   config.SlabSizeLimitBytes / 2, // 50% of the total slab size limit.
	}, nil
}

// ShouldProceed reports if there is more data to be fetched.
func (plan *Plan) ShouldProceed() bool {
	return plan.nextMinTimestamp < plan.config.MaxTimestamp
}

// nextTimeDelta determines if we can safely pull a larger time range from the prometheus server given read and slab
// size limits.
func (plan *Plan) nextTimeDelta(numBytes, limit int64, prevTimeDelta int64) (int64, error) {
	var delta int64

	if numBytes <= limit/2 {
		delta = prevTimeDelta + plan.config.LaIncrement.Milliseconds()

		if delta > plan.config.MaxReadDuration.Milliseconds() {
			delta = plan.config.MaxReadDuration.Milliseconds()
		}
	} else if numBytes > limit {
		delta = prevTimeDelta / 2
	} else {
		return prevTimeDelta, nil
	}

	if delta/time.Minute.Milliseconds() < 1 {
		return delta, fmt.Errorf("slab time range less than a minute, increate slab / read limitations")
	}

	return delta, nil
}

// NextSlab returns a new slab after allocating the time-range for fetch.
func (plan *Plan) NextSlab() (reference *Slab, err error) {
	timeDelta, err := plan.nextTimeDelta(atomic.LoadInt64(&plan.lastNumBytes), plan.config.SlabSizeLimitBytes, plan.lastTimeRangeDelta)
	if err != nil {
		return nil, fmt.Errorf("could not determine time delta: %w", err)
	}

	mint := plan.nextMinTimestamp
	maxt := mint + timeDelta

	if maxt > plan.config.MaxTimestamp {
		maxt = plan.config.MaxTimestamp
	}

	plan.nextMinTimestamp = maxt
	plan.lastTimeRangeDelta = timeDelta
	bRef, err := NewSlab(plan, mint, maxt)

	if err != nil {
		return nil, fmt.Errorf("could not create-slab: %w", err)
	}

	return bRef, nil
}

// update the details of the planner that are dependent on previous fetch stats.
func (plan *Plan) update(numBytes int, percentDone float64) {
	atomic.StoreInt64(&plan.lastNumBytes, int64(numBytes))

	plan.percentDone = percentDone
	if plan.percentDone >= plan.nextPercentThreshold || plan.percentDone >= 100 {
		plan.nextPercentThreshold += 10
		lg.Infof("read %%%f done", plan.percentDone)
	}
}
