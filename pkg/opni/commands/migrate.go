//go:build !noagentv1

package commands

import (
	"context"
	"fmt"
	"github.com/cortexproject/cortex/pkg/storage/bucket"
	"github.com/cortexproject/cortex/pkg/util/log"
	"github.com/cortexproject/cortex/tools/thanosconvert"
	"github.com/inhies/go-bytesize"
	"github.com/prometheus/common/config"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql/parser"
	"github.com/rancher/opni/pkg/migrate"
	plan "github.com/rancher/opni/pkg/migrate/planner"
	"github.com/rancher/opni/pkg/migrate/reader"
	"github.com/rancher/opni/pkg/migrate/writer"
	"github.com/spf13/cobra"
	"github.com/weaveworks/common/logging"
	"os"
	"sigs.k8s.io/yaml"
	"strconv"
	"strings"
	"time"
)

const (
	migrationJobName       = "prom-migrator"
	defaultTimeout         = time.Minute * 5
	defaultRetryDelay      = time.Second
	defaultStartTime       = "1970-01-01T00:00:00+00:00" // RFC3339 based time.Unix from 0 seconds.
	defaultMaxReadDuration = time.Hour * 2
	defaultLaIncrement     = time.Minute
)

// timeNowUnix returns the surrent Unix timestamp.
// time.Now().Unix() call is wrapped so it can be mocked during testing.
var timeNowUnix = func() int64 {
	return time.Now().Unix()
}

type migrateConfig struct {
	name  string
	start string
	end   string

	// generated from start and end
	minTimestamp    int64
	minTimestampSec int64

	maxTimestamp             int64
	maxTimestampSec          int64
	maxSlabSizeBytes         int64
	maxSlabSizeHumanReadable string
	concurrentPull           int
	concurrentPush           int
	readerClientConfig       migrate.ClientConfig
	maxReadDuration          time.Duration // Max range of slab in terms of time.
	writerClientConfig       migrate.ClientConfig
	garbageCollectOnPush     bool
	laIncrement              time.Duration

	readerTls config.TLSConfig
	writerTls config.TLSConfig

	readerMetricsMatcher string
	readerLabelsMatcher  []*labels.Matcher
}

func migrateData(cfg migrateConfig) error {
	if err := validateConf(&cfg); err != nil {
		return fmt.Errorf("config is not valid: %w", err)

	}

	maxSlabSizeBytes, err := bytesize.Parse(cfg.maxSlabSizeHumanReadable)
	if err != nil {
		return fmt.Errorf("error parsing byte-size: %w", err)
	}
	cfg.maxSlabSizeBytes = int64(maxSlabSizeBytes)

	planConfig := &plan.Config{
		MinTimestamp:       cfg.minTimestamp,
		MaxTimestamp:       cfg.maxTimestamp,
		JobName:            cfg.name,
		SlabSizeLimitBytes: cfg.maxSlabSizeBytes,
		MaxReadDuration:    cfg.maxReadDuration,
		LaIncrement:        cfg.laIncrement,
		NumStores:          cfg.concurrentPull,
	}
	planner, err := plan.NewPlan(planConfig)
	if err != nil {
		return fmt.Errorf("could not create plan: %w", err)
	}

	readErrChan := make(chan error)
	writeErrChan := make(chan error)
	slabChan := make(chan *plan.Slab)

	cont, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	readerConfig := reader.Config{
		Context:         cont,
		ClientConfig:    cfg.readerClientConfig,
		Plan:            planner,
		HTTPConfig:      config.HTTPClientConfig{TLSConfig: cfg.readerTls},
		ConcurrentPulls: cfg.concurrentPull,
		SlabChan:        slabChan,
		MetricsMatchers: cfg.readerLabelsMatcher,
	}
	read, err := reader.NewReader(readerConfig)
	if err != nil {
		return fmt.Errorf("could not create reader: %w", err)
	}

	writerConfig := writer.Config{
		Context:              cont,
		ClientConfig:         cfg.writerClientConfig,
		HTTPConfig:           config.HTTPClientConfig{TLSConfig: cfg.writerTls},
		GarbageCollectOnPush: cfg.garbageCollectOnPush,
		MigrationJobName:     cfg.name,
		ConcurrentPush:       cfg.concurrentPush,
		SlabChan:             slabChan,
	}
	write, err := writer.NewWriter(writerConfig)
	if err != nil {
		return fmt.Errorf("could not create writer: %w", err)
	}

	read.Run(readErrChan)
	write.Run(writeErrChan)

loop:
	for {
		select {
		case err = <-readErrChan:
			if err != nil {
				cancelFunc()
				return fmt.Errorf("error running reader: %w", err)
			}
		case err, ok := <-writeErrChan:
			_ = err
			cancelFunc() // As in any ideal case, the reader will always exit normally first.
			if ok {
				return fmt.Errorf("error running writer: %w", err)
			}
			break loop
		}
	}

	lg.Infof("finished migration")

	return nil
}

func parseRFC3339(s string) (int64, error) {
	if s == "" {
		return timeNowUnix(), nil
	}

	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return -1, fmt.Errorf("cannot convert '%s' to timestamp: %w", s, err)
	}

	return t.Unix(), nil
}

func toTimestamp(s string) (timestamp int, err error) {
	timestamp, err = strconv.Atoi(s)

	if err != nil {
		if t, err := parseRFC3339(s); err != nil {
			return 0, fmt.Errorf("timestamp must be either unix timestamp or in RFC3339 format: %w", err)
		} else {
			timestamp = int(t)
		}
	}

	return
}

func convertTimeStrFlagsToTs(conf *migrateConfig) error {
	startTimestamp, err := toTimestamp(conf.start)
	if err != nil {
		return fmt.Errorf("start timestamp invalid: %w", err)
	}

	conf.minTimestampSec = int64(startTimestamp)

	if conf.end == "" {
		conf.end = fmt.Sprintf("%d", timeNowUnix())
	}

	endTimestamp, err := toTimestamp(conf.end)

	if err != nil {
		return fmt.Errorf("end timestamp invalid: %w", err)
	}

	conf.maxTimestamp = int64(endTimestamp)

	// remote-storages tend to respond to time in milliseconds. So, we convert the received values in seconds to milliseconds.
	conf.minTimestamp = conf.minTimestampSec * 1000
	conf.maxTimestamp = conf.maxTimestamp * 1000
	return nil
}

func convertMetricsMatcherStrToLabelMatchers(conf *migrateConfig) error {
	ms, err := parser.ParseExpr(conf.readerMetricsMatcher)
	if err != nil {
		return fmt.Errorf("parsing '-reader-metrics-matcher': %w", err)
	}

	vs, ok := ms.(*parser.VectorSelector)
	if !ok {
		return fmt.Errorf("'-reader-metrics-matcher' should always be a valid vector_selector as per PromQL")
	}

	conf.readerLabelsMatcher = vs.LabelMatchers

	return nil
}

func validateConf(conf *migrateConfig) error {
	if err := convertTimeStrFlagsToTs(conf); err != nil {
		return fmt.Errorf("%w", err)
	}

	if err := migrate.ParseClientInfo(&conf.readerClientConfig); err != nil {
		return fmt.Errorf("parsing reader-client info: %w", err)
	}

	if err := migrate.ParseClientInfo(&conf.writerClientConfig); err != nil {
		return fmt.Errorf("parsing writer-client info: %w", err)
	}
	if err := convertMetricsMatcherStrToLabelMatchers(conf); err != nil {
		return fmt.Errorf("validate '-reader-metrics-matcher': %w", err)
	}

	switch {
	case conf.minTimestampSec > conf.maxTimestamp:
		return fmt.Errorf("end time cannot be before start time")
	case conf.laIncrement < time.Minute:
		return fmt.Errorf("'--slab-range-increment' cannot be less than 1 minute")
	case conf.maxReadDuration < time.Minute:
		return fmt.Errorf("'--max-read-duration' cannot be less than 1 minute")
	}

	return nil
}

func migrateMetaData(configFile string, dryRun bool) error {
	// effectively a copy of https://github.com/cortexproject/cortex/blob/master/cmd/thanosconvert/main.go
	var cfg bucket.Config

	buffer, err := os.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("could not read migrateConfig file at '%s': %w", configFile, err)
	}

	if err := yaml.Unmarshal(buffer, &cfg); err != nil {
		return fmt.Errorf("could not parse migrateConfig file '%s': %w", configFile, err)
	}

	ctx := context.Background()

	logger, err := log.NewPrometheusLogger(logging.Level{}, logging.Format{})
	converter, err := thanosconvert.NewThanosBlockConverter(ctx, cfg, dryRun, logger)
	if err != nil {
		return fmt.Errorf("could not creat converter")
	}

	iterCtx := context.Background()
	if iterCtx == nil {
		return fmt.Errorf("context was nil")
	}
	results, err := converter.Run(iterCtx)

	for user, res := range results {
		fmt.Printf("User %s:\n", user)
		fmt.Printf("  Converted %d:\n  %s", len(res.ConvertedBlocks), strings.Join(res.ConvertedBlocks, ","))
		fmt.Printf("  Unchanged %d:\n  %s", len(res.UnchangedBlocks), strings.Join(res.UnchangedBlocks, ","))
		fmt.Printf("  Failed %d:\n  %s", len(res.FailedBlocks), strings.Join(res.FailedBlocks, ","))
	}

	return nil
}

func BuildMigrateCmd() *cobra.Command {
	var cfg migrateConfig

	command := &cobra.Command{
		Use:   "migrate",
		Short: "migrate blocks to opni cortex",
		Long:  "See subcommands for more information.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg.laIncrement = time.Minute * 2

			//if err := migrateData(options); err != nil {
			if err := migrateData(cfg); err != nil {
				lg.Errorf("could not migrate data from '%s' -> '%s': %s", cfg.readerClientConfig.URL, cfg.writerClientConfig.URL, err)
				return nil
			}

			//if err := migrateMetaData(configFile, dryRun); err != nil {
			//	return fmt.Errorf("error migrating data meta: %w", err)
			//}

			return nil
		},
	}

	//command.Flags().StringVar(&cfg.name, "migration-name", migrationJobName, "Name for the current migration that is to be carried out.")
	cfg.name = migrationJobName

	command.Flags().StringVar(&cfg.start, "start", defaultStartTime, fmt.Sprintf("Start time (in RFC3339 format, like '%s', or in number of seconds since the unix epoch, UTC) from which the data migration is to be carried out. (inclusive)", defaultStartTime))

	command.Flags().StringVar(&cfg.end, "end", "", fmt.Sprintf("End time (in RFC3339 format, like '%s', or in number of seconds since the unix epoch, UTC) for carrying out data migration (exclusive). ", defaultStartTime)+
		"By default if this value is unset, then 'end' will be set to the time at which migration is starting. ")

	command.Flags().StringVar(&cfg.maxSlabSizeHumanReadable, "max-read-size", "500MB", "(units: B, KB, MB, GB, TB, PB) the maximum size of data that should be read at a single time. "+
		"More the read size, faster will be the migration but higher will be the memory usage. Example: 250MB.")

	command.Flags().BoolVar(&cfg.garbageCollectOnPush, "gc-on-push", false, "Run garbage collector after every slab is pushed. "+
		"This may lead to better memory management since GC is kick right after each slab to clean unused memory blocks.")

	command.Flags().IntVar(&cfg.concurrentPush, "concurrent-push", 1, "Concurrent push enables pushing of slabs concurrently. "+
		"Each slab is divided into 'concurrent-push' (value) parts and then pushed to the remote-write storage concurrently. This may lead to higher throughput on "+
		"the remote-write storage provided it is capable of handling the load. Note: Larger shards count will lead to significant memory usage.")

	command.Flags().IntVar(&cfg.concurrentPull, "concurrent-pull", 1, "Concurrent pull enables fetching of data concurrently. Note: setting this value too high may cause TLS handshake error on the read storage side or may lead to starvation of fetch requests")

	//command.Flags().DurationVar(&cfg.maxReadDuration, "max-read-duration", defaultMaxReadDuration, "Maximum range of slab that can be achieved through consecutive 'la-increment'. "+
	//	"This defaults to '2 hours' since assuming the migration to be from Prometheus TSDB. Increase this duration if you are not fetching from Prometheus "+
	//	"or if you are fine with slow read on Prometheus side post 2 hours slab time-range.")
	cfg.maxReadDuration = defaultMaxReadDuration

	//command.Flags().DurationVar(&cfg.laIncrement, "slab-range-increment", defaultLaIncrement, "Amount of time-range to be incremented in successive slab.")
	cfg.laIncrement = defaultLaIncrement

	command.Flags().StringVar(&cfg.readerClientConfig.URL, "reader-url", "", "URL address for the storage where the data is to be read from.")
	_ = command.MarkFlagRequired("reader-url")

	command.Flags().DurationVar(&cfg.readerClientConfig.Timeout, "reader-timeout", defaultTimeout, "Timeout for fetching data from read storage. "+
		"This timeout is also used to fetch the progress metric.")

	command.Flags().DurationVar(&cfg.readerClientConfig.RetryDelay, "reader-retry-delay", defaultRetryDelay, "Duration to wait after a 'read-timeout' "+
		"before prom-migrator retries to fetch the slab. RetryDelay is used only if 'retry' option is set in OnTimeout or OnErr.")

	command.Flags().IntVar(&cfg.readerClientConfig.MaxRetry, "reader-max-retries", 0, "Maximum number of retries before erring out. "+
		"Setting this to 0 will make the retry process forever until the process is completed. Note: If you want not to retry, "+
		"change the value of on-timeout or on-error to non-retry options.")

	command.Flags().StringVar(&cfg.readerClientConfig.OnTimeoutStr, "reader-on-timeout", "retry", "When a timeout happens during the read process, how should the reader behave. Valid options: ['retry', 'skip', 'abort']")

	command.Flags().StringVar(&cfg.readerClientConfig.OnErrStr, "reader-on-error", "abort", "When an error occurs during read process, how should the reader behave. Valid options: ['retry', 'skip', 'abort'].")

	command.Flags().StringVar(&cfg.readerMetricsMatcher, "reader-metrics-matcher", `{__name__=~".+"}`, "Metrics vector selector to read data for migration.")

	if cfg.readerClientConfig.CustomHeaders == nil {
		cfg.readerClientConfig.CustomHeaders = make(map[string][]string)
	}

	//command.Flags().Var(&utils.HeadersFlag{Headers: cfg.readerClientConfig.CustomHeaders}, "reader-http-header", "HTTP header to send with all the reader requests. It uses the format `key:value`, for example `-reader-http-header=\"X-Scope-OrgID:42\"`. Can be set multiple times to define several headers or multiple values for the same header.")

	command.Flags().StringVar(&cfg.writerClientConfig.URL, "writer-url", "", "URL address for the storage where the data migration is to be written.")
	_ = command.MarkFlagRequired("writer-url")

	command.Flags().DurationVar(&cfg.writerClientConfig.Timeout, "writer-timeout", defaultTimeout, "Timeout for pushing data to write storage.")

	command.Flags().DurationVar(&cfg.writerClientConfig.RetryDelay, "writer-retry-delay", defaultRetryDelay, "Duration to wait after a 'write-timeout' "+
		"before prom-migrator retries to push the slab. RetryDelay is used only if 'retry' option is set in OnTimeout or OnErr.")

	command.Flags().IntVar(&cfg.writerClientConfig.MaxRetry, "writer-max-retries", 0, "Maximum number of retries before erring out. "+
		"Setting this to 0 will make the retry process forever until the process is completed. Note: If you want not to retry, "+
		"change the value of on-timeout or on-error to non-retry options.")

	command.Flags().StringVar(&cfg.writerClientConfig.OnTimeoutStr, "writer-on-timeout", "retry", "When a timeout happens during the write process, how should the writer behave. Valid options: ['retry', 'skip', 'abort'].")

	command.Flags().StringVar(&cfg.writerClientConfig.OnErrStr, "writer-on-error", "abort", "When an error occurs during write process, how should the writer behave. Valid options: ['retry', 'skip', 'abort'].")

	if cfg.writerClientConfig.CustomHeaders == nil {
		cfg.writerClientConfig.CustomHeaders = make(map[string][]string)
	}
	//command.Flags().Var(&utils.HeadersFlag{Headers: cfg.writerClientConfig.CustomHeaders}, "writer-http-header", "HTTP header to send with all the writer requests. It uses the format `key:value`, for example `-writer-http-header=\"X-Scope-OrgID:42\"`. Can be set multiple times to define several headers or multiple values for the same header.")

	// TLS configurations.
	// Reader.
	command.Flags().StringVar(&cfg.readerTls.CAFile, "reader-tls-ca-file", "", "TLS CA file for remote-read component.")

	command.Flags().StringVar(&cfg.readerTls.CertFile, "reader-tls-cert-file", "", "TLS certificate file for remote-read component.")

	command.Flags().StringVar(&cfg.readerTls.KeyFile, "reader-tls-key-file", "", "TLS key file for remote-read component.")

	command.Flags().StringVar(&cfg.readerTls.ServerName, "reader-tls-server-name", "", "TLS server name for remote-read component.")

	command.Flags().BoolVar(&cfg.readerTls.InsecureSkipVerify, "reader-tls-insecure-skip-verify", false, "TLS insecure skip verify for remote-read component.")

	// Writer.
	command.Flags().StringVar(&cfg.writerTls.CAFile, "writer-tls-ca-file", "", "TLS CA file for remote-writer component.")

	command.Flags().StringVar(&cfg.writerTls.CertFile, "writer-tls-cert-file", "", "TLS certificate file for remote-writer component.")

	command.Flags().StringVar(&cfg.writerTls.KeyFile, "writer-tls-key-file", "", "TLS key file for remote-writer component.")

	command.Flags().StringVar(&cfg.writerTls.ServerName, "writer-tls-server-name", "", "TLS server name for remote-writer component.")

	command.Flags().BoolVar(&cfg.writerTls.InsecureSkipVerify, "writer-tls-insecure-skip-verify", false, "TLS insecure skip verify for remote-writer component.")

	return command
}

func init() {
	AddCommandsToGroup(Utilities, BuildMigrateCmd())
}
