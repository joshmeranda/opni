package migrate

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/rancher/opni/pkg/logger"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/snappy"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/pkg/errors"
	promConfig "github.com/prometheus/common/config"
	"github.com/prometheus/common/version"
	"github.com/prometheus/prometheus/prompb"
)

var (
	userAgent = fmt.Sprintf("Prometheus/%s", version.Version)
	lg        = logger.New().Named("migrator")
)

const (
	maxErrMsgLen = 512
)

type (
	RecoveryAction uint8
	eventType      uint8
)

const (
	Retry RecoveryAction = iota
	Skip
	Abort
)

const (
	onTimeout eventType = iota
	onError
)

type ClientConfig struct {
	URL           string
	Timeout       time.Duration
	OnTimeout     RecoveryAction
	OnTimeoutStr  string
	OnErr         RecoveryAction
	OnErrStr      string
	MaxRetry      int
	RetryDelay    time.Duration
	CustomHeaders map[string][]string
}

func (c ClientConfig) shouldRetry(on eventType, retries int, message string) (retry bool) {
	makeDecision := func(r RecoveryAction) bool {
		switch r {
		case Retry:
			if c.MaxRetry != 0 {
				if retries > c.MaxRetry {
					lg.Errorf("exceeded retrying limit in %s. Erring out.", message)
					return false
				}
			}

			lg.Infof("%s. Retrying again after delay", message)
			time.Sleep(c.RetryDelay)

			return true
		case Abort:
		}

		return false
	}

	switch on {
	case onTimeout:
		return makeDecision(c.OnTimeout)
	case onError:
		return makeDecision(c.OnErr)
	default:
		panic("invalid 'on'. Should be 'error' or 'timeout'")
	}
}

func RecoveryActionEnum(typ string) (RecoveryAction, error) {
	switch typ {
	case "retry":
		return Retry, nil
	case "skip":
		return Skip, nil
	case "abort":
		return Abort, nil
	default:
		return 0, fmt.Errorf("invalid slab enums: expected from ['retry', 'skip', 'abort']")
	}
}

func ParseClientInfo(conf *ClientConfig) error {
	wrapErr := func(typ string, err error) error {
		return fmt.Errorf("enum conversion '%s': %w", typ, err)
	}

	enum, err := RecoveryActionEnum(conf.OnTimeoutStr)
	if err != nil {
		return wrapErr("on-timeout", err)
	}
	conf.OnTimeout = enum

	enum, err = RecoveryActionEnum(conf.OnErrStr)
	if err != nil {
		return wrapErr("on-error", err)
	}
	conf.OnErr = enum

	return nil
}

type Client struct {
	remoteName    string
	url           *promConfig.URL
	cHTTP         *http.Client
	clientConfig  ClientConfig
	customHeaders map[string][]string
}

// NewClient creates a new read or write client. The `clientType` should be either `read` or `write`. The client type
// is used to get the auth from the auth store. If the `clientType` is other than the ones specified, then auth may not work.
func NewClient(remoteName string, clientConfig ClientConfig, httpConfig promConfig.HTTPClientConfig) (*Client, error) {
	parsedUrl, err := url.Parse(clientConfig.URL)
	if err != nil {
		return nil, fmt.Errorf("parsing-%s-url: %w", remoteName, err)
	}
	httpClient, err := promConfig.NewClientFromConfig(httpConfig, fmt.Sprintf("%s_client", remoteName), promConfig.WithHTTP2Disabled())
	if err != nil {
		return nil, err
	}

	t := httpClient.Transport
	httpClient.Transport = &nethttp.Transport{
		RoundTripper: t,
	}

	return &Client{
		remoteName:    remoteName,
		url:           &promConfig.URL{URL: parsedUrl},
		cHTTP:         httpClient,
		clientConfig:  clientConfig,
		customHeaders: clientConfig.CustomHeaders,
	}, nil
}

func (c *Client) addCustomHeaders(h http.Header) {
	for k, vs := range c.customHeaders {
		for _, v := range vs {
			h.Add(k, v)
		}
	}
}

// Config returns the client runtime.
func (c *Client) Config() ClientConfig {
	return c.clientConfig
}

func (c *Client) Read(ctx context.Context, query *prompb.Query) (result *prompb.QueryResult, numBytesCompressed int, numBytesUncompressed int, err error) {
	req := &prompb.ReadRequest{
		Queries: []*prompb.Query{
			query,
		},
	}
	data, err := proto.Marshal(req)
	if err != nil {
		return nil, -1, -1, errors.Wrapf(err, "unable to marshal read request")
	}

	numAttempts := -1
retry:
	numAttempts++
	if numAttempts > 1 {
		// Retrying case. Reset the context so that timeout context works.
		// If we do not reset the context, retrying with the timeout will not work
		// since the context parent has already been cancelled.
		ctx.Done()
		//
		// Note: We use httpReq.WithContext(ctx) which means if a timeout occurs and we are asked
		// to retry, the context in the httpReq has already been cancelled. Hence, we should create
		// the new context for retry with a new httpReq otherwise we are reusing the already cancelled
		// httpReq.
		ctx = context.Background()
	}

	// Recompress the data on retry, otherwise the bytes.NewReader will make it invalid,
	// leading to corrupt snappy data output on remote-read server.
	compressed := snappy.Encode(nil, data)

	httpReq, err := http.NewRequest("POST", c.url.String(), bytes.NewReader(compressed))
	if err != nil {
		return nil, -1, -1, errors.Wrap(err, "unable to create request")
	}

	c.addCustomHeaders(httpReq.Header)
	httpReq.Header.Add("Content-Encoding", "snappy")
	httpReq.Header.Add("Accept-Encoding", "snappy")
	httpReq.Header.Set("Content-Type", "application/x-protobuf")
	httpReq.Header.Set("User-Agent", userAgent)
	httpReq.Header.Set("X-Prometheus-Remote-Read-Version", "0.1.0")

	ctx, cancel := context.WithTimeout(ctx, c.clientConfig.Timeout)
	defer cancel()

	httpReq = httpReq.WithContext(ctx)

	httpResp, err := c.cHTTP.Do(httpReq)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			if c.clientConfig.shouldRetry(onTimeout, numAttempts, "read request timeout") {
				goto retry
			}
			return nil, 0, 0, fmt.Errorf("error sending request: %w", err)
		}
		err = fmt.Errorf("error sending request: %w", err)
		if c.clientConfig.shouldRetry(onError, numAttempts, fmt.Sprintf("error:\n%s\n", err.Error())) {
			goto retry
		}
		return nil, 0, 0, err
	}
	defer func() {
		_, _ = io.Copy(io.Discard, httpResp.Body)
		_ = httpResp.Body.Close()
	}()

	var reader bytes.Buffer
	_, _ = io.Copy(&reader, httpResp.Body)

	compressed, err = io.ReadAll(bytes.NewReader(reader.Bytes()))
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("error reading response. HTTP status code: %s", httpResp.Status))
		if c.clientConfig.shouldRetry(onError, numAttempts, fmt.Sprintf("error:\n%s\n", err.Error())) {
			goto retry
		}
		return nil, -1, -1, err
	}

	if httpResp.StatusCode/100 != 2 {
		if c.clientConfig.shouldRetry(onError, numAttempts, fmt.Sprintf("received status code: %d", httpResp.StatusCode)) {
			goto retry
		}
		return nil, -1, -1, errors.Errorf("remote server %s returned HTTP status %s: %s", c.url.String(), httpResp.Status, strings.TrimSpace(string(compressed)))
	}

	uncompressed, err := snappy.Decode(nil, compressed)
	if err != nil {
		if c.clientConfig.shouldRetry(onError, numAttempts, fmt.Sprintf("error:\n%s\n", err.Error())) {
			goto retry
		}
		return nil, -1, -1, errors.Wrap(err, "error reading response")
	}

	var resp prompb.ReadResponse
	err = proto.Unmarshal(uncompressed, &resp)
	if err != nil {
		err = errors.Wrap(err, "unable to unmarshal response body")
		if c.clientConfig.shouldRetry(onError, numAttempts, fmt.Sprintf("error:\n%s\n", err.Error())) {
			goto retry
		}
		return nil, -1, -1, err
	}

	if len(resp.Results) != len(req.Queries) {
		err = errors.Errorf("responses: want %d, got %d", len(req.Queries), len(resp.Results))
		if c.clientConfig.shouldRetry(onError, numAttempts, fmt.Sprintf("error:\n%s\n", err.Error())) {
			goto retry
		}
		return nil, -1, -1, err
	}

	return resp.Results[0], len(compressed), len(uncompressed), nil
}

// ReadConcurrent calls the Read and responds on the channels.
func (c *Client) ReadConcurrent(ctx context.Context, query *prompb.Query, shardID int, responseChan chan<- interface{}) {
	result, numBytesCompressed, numBytesUncompressed, err := c.Read(ctx, query)
	if err != nil {
		responseChan <- fmt.Errorf("read-channels: %w", err)
		return
	}

	responseChan <- &PrompbResponse{
		ID:                   shardID,
		Result:               result,
		NumBytesCompressed:   numBytesCompressed,
		NumBytesUncompressed: numBytesUncompressed,
	}
}

// RecoverableError is an error for which the send samples can retry with a backoff.
type RecoverableError struct {
	error
}

func (c *Client) Write(ctx context.Context, req []byte) error {
	httpReq, err := http.NewRequest("POST", c.url.String(), bytes.NewReader(req))
	if err != nil {
		return err
	}

	c.addCustomHeaders(httpReq.Header)
	httpReq.Header.Add("Content-Encoding", "snappy")
	httpReq.Header.Set("Content-Type", "application/x-protobuf")
	httpReq.Header.Set("User-Agent", userAgent)
	httpReq.Header.Set("X-Prometheus-Remote-Write-Version", "0.1.0")
	ctx, cancel := context.WithTimeout(ctx, c.clientConfig.Timeout)
	defer cancel()

	httpReq = httpReq.WithContext(ctx)

	httpResp, err := c.cHTTP.Do(httpReq)
	if err != nil {
		// Errors from Client.Do are from (for example) network errors, so are
		// recoverable.
		return RecoverableError{err}
	}

	defer func() {
		_, _ = io.Copy(io.Discard, httpResp.Body)
		_ = httpResp.Body.Close()
	}()

	if httpResp.StatusCode/100 != 2 {
		scanner := bufio.NewScanner(io.LimitReader(httpResp.Body, maxErrMsgLen))
		line := ""
		if scanner.Scan() {
			line = scanner.Text()
		}
		err = errors.Errorf("server returned HTTP status %s: %s", httpResp.Status, line)
	}
	if httpResp.StatusCode/100 == 5 {
		return RecoverableError{err}
	}
	return err
}
