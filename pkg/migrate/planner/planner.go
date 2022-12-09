package planner

import (
	"fmt"
	"github.com/rancher/opni/pkg/logger"
	"go.uber.org/zap/zapcore"
	"sync/atomic"
	"time"
)

var (
	minute = time.Minute.Milliseconds()
	lg     = logger.New(logger.WithLogLevel(zapcore.InfoLevel)).Named("planner")
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

func (c *Config) determineTimeDelta(numBytes, limit int64, prevTimeDelta int64) (int64, error) {
	errOnBelowMinute := func(delta int64) (int64, error) {
		if delta/minute < 1 {
			// This results in an infinite loop, if not errored out. Occurs when the max number of bytes is so low
			// that even the size of 1 minute slab of migration exceeds the bytes limit, resulting in further (exponential)
			// decrease of time-range.
			//
			// Hence, stop migration and report in logs.
			return delta, fmt.Errorf("slab time-range less than a minute. Please increase slab/read limitations")
		}
		return delta, nil
	}

	switch {
	case numBytes <= limit/2:
		// We increase the time-range linearly for the next fetch if the current time-range fetch resulted in size that is
		// less than half the max read size. This continues till we reach the maximum time-range delta.
		return errOnBelowMinute(c.clampTimeDelta(prevTimeDelta + c.LaIncrement.Milliseconds()))
	case numBytes > limit:
		// Decrease size the time exponentially so that bytes size can be controlled (preventing OOM).
		lg.Infof("decreasing time-range delta to %d minute(s) since size beyond permittable limits", prevTimeDelta/(2*minute))
		return errOnBelowMinute(prevTimeDelta / 2)
	}

	// Here, the numBytes is between the max increment-time size limit (i.e., limit/2) and the max read limit
	// (i.e., increment-time size limit < numBytes <= max read limit). This region is an ideal case of
	// balance between migration speed and memory utilization.
	//
	// Example: If the limit is 500MB, then the max increment-time size limit will be 250MB. This means that till the numBytes is below
	// 250MB, the time-range for next fetch will continue to increase by 1 minute (on the previous fetch time-range). However,
	// the moment any slab comes between 250MB and 500MB, we stop to increment the time-range delta further. This helps
	// keeping the migration tool in safe memory limits.
	return prevTimeDelta, nil
}

func (c *Config) clampTimeDelta(t int64) int64 {
	if t > c.MaxReadDuration.Milliseconds() {
		lg.Debugf("Exceeded '%s', setting back to 'max-read-duration' value", c.MaxReadDuration.String())
		return c.MaxReadDuration.Milliseconds()
	}
	return t
}

type Plan struct {
	config *Config

	nextMint           int64
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
		config:         config,
		nextMint:       config.MinTimestamp,
		deltaIncRegion: config.SlabSizeLimitBytes / 2, // 50% of the total slab size limit.
	}, nil
}

// ShouldProceed reports if there is more data to be fetched.
func (plan *Plan) ShouldProceed() bool {
	return plan.nextMint < plan.config.MaxTimestamp
}

// update updates the details of the planner that are dependent on previous fetch stats.
func (plan *Plan) update(numBytes int) {
	atomic.StoreInt64(&plan.lastNumBytes, int64(numBytes))
}

// NextSlab returns a new slab after allocating the time-range for fetch.
func (plan *Plan) NextSlab() (reference *Slab, err error) {
	timeDelta, err := plan.config.determineTimeDelta(atomic.LoadInt64(&plan.lastNumBytes), plan.config.SlabSizeLimitBytes, plan.lastTimeRangeDelta)
	if err != nil {
		return nil, fmt.Errorf("could not determine time delta: %w", err)
	}

	mint := plan.nextMint
	maxt := mint + timeDelta

	if maxt > plan.config.MaxTimestamp {
		maxt = plan.config.MaxTimestamp
	}

	plan.nextMint = maxt
	plan.lastTimeRangeDelta = timeDelta
	bRef, err := NewSlab(plan, mint, maxt)

	if err != nil {
		return nil, fmt.Errorf("could not create-slab: %w", err)
	}

	return bRef, nil
}

func (plan *Plan) reportPercentDone(percentDone float64) {
	plan.percentDone = percentDone

	if plan.percentDone >= plan.nextPercentThreshold || plan.percentDone >= 100 {
		plan.nextPercentThreshold += 10
		lg.Infof("read %%%f done", plan.percentDone)
	}
}
