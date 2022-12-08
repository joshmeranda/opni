// This file and its contents are licensed under the Apache License 2.0.
// Please see the included NOTICE for copyright information and
// LICENSE for a copy of the license.

package reader

import (
	"context"
	"fmt"
	"github.com/rancher/opni/pkg/logger"
	"github.com/rancher/opni/pkg/migrate/utils"

	"github.com/prometheus/common/config"
	"github.com/prometheus/prometheus/model/labels"

	plan "github.com/rancher/opni/pkg/migrate/planner"
)

var (
	lg = logger.New().Named("reader")
)

// Config is config for reader.
type Config struct {
	Context      context.Context
	ClientConfig utils.ClientConfig
	Plan         *plan.Plan
	HTTPConfig   config.HTTPClientConfig

	ConcurrentPulls int

	SigSlabRead chan *plan.Slab // To the writer.
	SigSlabStop chan struct{}

	MetricsMatchers []*labels.Matcher
}

type Reader struct {
	Config
	client *utils.Client
}

// NewReader creates a new Reader. It creates a ReadClient that is imported from Prometheus remote storage.
// Reader takes help of plan to understand how to create fetchers.
func NewReader(config Config) (*Reader, error) {
	rc, err := utils.NewClient(fmt.Sprintf("reader-%d", 1), config.ClientConfig, config.HTTPConfig)
	if err != nil {
		return nil, fmt.Errorf("could not creat read-client: %w", err)
	}

	read := &Reader{
		Config: config,
		client: rc,
	}

	return read, nil
}

// Run runs the remote read and starts fetching the samples from the read storage.
func (r *Reader) Run(errChan chan<- error) {
	var (
		err     error
		slabRef *plan.Slab
	)

	go func() {
		defer func() {
			close(r.SigSlabRead)
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
			case <-r.SigSlabStop:
				return
			default:
			}

			slabRef, err = r.Plan.NextSlab()
			if err != nil {
				errChan <- fmt.Errorf("remote-read run: %w", err)
				return
			}

			ms := r.Config.MetricsMatchers
			if len(ms) == 0 {
				lg.Info("empty matchers received, matching everything")
				ms = []*labels.Matcher{labels.MustNewMatcher(labels.MatchRegexp, labels.MetricName, ".+")}
			}

			err = slabRef.Fetch(r.Context, r.client, slabRef.MinTimestamp(), slabRef.MaxTimestamp(), ms)
			if err != nil {
				errChan <- fmt.Errorf("remote-read run: %w", err)
				return
			}

			if slabRef.IsEmpty() {
				continue
			}

			r.SigSlabRead <- slabRef
			slabRef = nil
		}
	}()
}
