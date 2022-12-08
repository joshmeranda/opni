package planner

import (
	"context"
	"fmt"
	"github.com/rancher/opni/pkg/logger"
	"github.com/rancher/opni/pkg/migrate"
	"github.com/rancher/opni/pkg/migrate/utils"
	"go.uber.org/zap/zapcore"
	"sync/atomic"
	"time"

	"github.com/prometheus/common/config"
	"github.com/prometheus/prometheus/model/labels"
)

var (
	//second = time.Second.Milliseconds()
	minute = time.Minute.Milliseconds()
	lg     = logger.New(logger.WithLogLevel(zapcore.InfoLevel)).Named("planner")
)

// Config represents configuration for the planner.
type Config struct {
	MinTimestamp         int64
	MaxTimestamp         int64
	SlabSizeLimitBytes   int64
	NumStores            int
	ProgressEnabled      bool
	JobName              string
	ProgressMetricName   string // Name for progress metric.
	ProgressClientConfig utils.ClientConfig
	HTTPConfig           config.HTTPClientConfig
	LaIncrement          time.Duration
	MaxReadDuration      time.Duration
}

// fetchLastPushedMaxTimestamp fetches the max timestamp of the last slab pushed to remote-write storage.
func (c *Config) fetchLastPushedMaxTimestamp() (lastPushedMaxt int64, found bool, err error) {
	query, err := migrate.CreatePrombQuery(c.MinTimestamp, c.MaxTimestamp, []*labels.Matcher{
		labels.MustNewMatcher(labels.MatchEqual, labels.MetricName, c.ProgressMetricName),
		labels.MustNewMatcher(labels.MatchEqual, migrate.LabelJob, c.JobName),
	})
	if err != nil {
		return -1, false, fmt.Errorf("could not create fetch-last-pushed-maxTimestamp promb query: %w", err)
	}

	readClient, err := utils.NewClient("reader-last-maxTimestamp-pushed", c.ProgressClientConfig, c.HTTPConfig)
	if err != nil {
		return -1, false, fmt.Errorf("could not create fetch-last-pushed-maxTimestamp reader: %w", err)
	}

	result, _, _, err := readClient.Read(context.Background(), query)
	if err != nil {
		return -1, false, fmt.Errorf("fetch-last-pushed-maxTimestamp query result: %w", err)
	}

	ts := result.Timeseries
	if len(ts) == 0 {
		return -1, false, nil
	}

	for _, series := range ts {
		for i := len(series.Samples) - 1; i >= 0; i-- {
			if series.Samples[i].Timestamp > lastPushedMaxt {
				lastPushedMaxt = series.Samples[i].Timestamp
			}
		}
	}

	if lastPushedMaxt == 0 {
		return -1, false, nil
	}

	return lastPushedMaxt, true, nil
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
	deltaIncRegion     int64 // Time region for which the time-range delta can continue to increase by laIncrement.
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
func (p *Plan) ShouldProceed() bool {
	return p.nextMint < p.config.MaxTimestamp
}

// todo: needed?
// update updates the details of the planner that are dependent on previous fetch stats.
func (p *Plan) update(numBytes int) {
	atomic.StoreInt64(&p.lastNumBytes, int64(numBytes))
}

// NextSlab returns a new slab after allocating the time-range for fetch.
func (p *Plan) NextSlab() (reference *Slab, err error) {
	timeDelta, err := p.config.determineTimeDelta(atomic.LoadInt64(&p.lastNumBytes), p.config.SlabSizeLimitBytes, p.lastTimeRangeDelta)
	if err != nil {
		return nil, fmt.Errorf("could not determine time delta: %w", err)
	}

	mint := p.nextMint
	maxt := mint + timeDelta

	if maxt > p.config.MaxTimestamp {
		maxt = p.config.MaxTimestamp
	}

	p.nextMint = maxt
	p.lastTimeRangeDelta = timeDelta
	bRef, err := p.createSlab(mint, maxt)

	if err != nil {
		return nil, fmt.Errorf("could not create-slab: %w", err)
	}

	return bRef, nil
}

// createSlab creates a new slab and returns reference to the slab for faster write and read operations.
func (p *Plan) createSlab(minTimestamp, maxTimestamp int64) (ref *Slab, err error) {
	if err = p.validateT(minTimestamp, maxTimestamp); err != nil {
		return nil, fmt.Errorf("could not create-slab: %w", err)
	}

	percent := float64(maxTimestamp-p.config.MinTimestamp) * 100 / float64(p.config.MaxTimestamp-p.config.MinTimestamp)
	if percent > 100 {
		percent = 100
	}

	ref = slabPool.Get().(*Slab)
	ref.numStores = p.config.NumStores
	ref.percentDone = percent

	if cap(ref.stores) < p.config.NumStores {
		ref.stores = make([]store, p.config.NumStores)
	} else {
		ref.stores = ref.stores[:p.config.NumStores]
	}

	ref.minTimestamp = minTimestamp
	ref.maxTimestamp = maxTimestamp
	ref.plan = p

	ref.initStores()

	return
}

func (p *Plan) validateT(minTimestamp, maxTimestamp int64) error {
	switch {
	case p.config.MinTimestamp > minTimestamp || p.config.MaxTimestamp < minTimestamp:
		return fmt.Errorf("invalid minTimestamp: %d: global-minTimestamp: %d and global-maxTimestamp: %d", minTimestamp, p.config.MinTimestamp, p.config.MaxTimestamp)
	case p.config.MinTimestamp > maxTimestamp || p.config.MaxTimestamp < maxTimestamp:
		return fmt.Errorf("invalid maxTimestamp: %d: global-minTimestamp: %d and global-maxTimestamp: %d", minTimestamp, p.config.MinTimestamp, p.config.MaxTimestamp)
	case minTimestamp > maxTimestamp:
		return fmt.Errorf("minTimestamp cannot be greater than maxTimestamp: minTimestamp: %d and maxTimestamp: %d", minTimestamp, maxTimestamp)
	}

	return nil
}
