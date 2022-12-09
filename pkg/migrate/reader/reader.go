package reader

import (
	"context"
	"fmt"
	"github.com/prometheus/common/config"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/rancher/opni/pkg/logger"
	"github.com/rancher/opni/pkg/migrate"

	plan "github.com/rancher/opni/pkg/migrate/planner"
)

var (
	lg = logger.New().Named("reader")
)

// Config is config for reader.
type Config struct {
	Context      context.Context
	ClientConfig migrate.ClientConfig
	Plan         *plan.Plan
	HTTPConfig   config.HTTPClientConfig

	ConcurrentPulls int

	SlabChan chan *plan.Slab

	MetricsMatchers []*labels.Matcher
}

type Reader struct {
	Config
	client *migrate.Client
}

func NewReader(config Config) (*Reader, error) {
	client, err := migrate.NewClient(fmt.Sprintf("reader-%d", 1), config.ClientConfig, config.HTTPConfig)
	if err != nil {
		return nil, fmt.Errorf("could not create read client: %w", err)
	}

	read := &Reader{
		Config: config,
		client: client,
	}

	return read, nil
}

// Run runs the remote read and starts fetching the samples from the read storage.
func (r *Reader) Run(errChan chan<- error) {
	var err error
	var slab *plan.Slab

	go func() {
		defer func() {
			close(r.SlabChan)
			lg.Info("reader is down")
			close(errChan)
		}()

		lg.Info("reader is up")
		select {
		case <-r.Context.Done():
			return
		default:
		}

		for r.Plan.ShouldProceed() {
			select {
			case <-r.Context.Done():
				return
			default:
			}

			slab, err = r.Plan.NextSlab()
			if err != nil {
				errChan <- fmt.Errorf("read error: %w", err)
				return
			}

			matchers := r.Config.MetricsMatchers
			if len(matchers) == 0 {
				lg.Debugf("empty matchers received, matching everything")
				matchers = []*labels.Matcher{labels.MustNewMatcher(labels.MatchRegexp, labels.MetricName, ".+")}
			}

			err = slab.Fetch(r.Context, r.client, matchers)
			if err != nil {
				errChan <- fmt.Errorf("read error: %w", err)
				return
			}

			if slab.IsEmpty() {
				continue
			}

			r.SlabChan <- slab
			slab = nil
		}
	}()
}
