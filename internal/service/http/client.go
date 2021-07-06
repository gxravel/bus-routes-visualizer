package http

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	api "github.com/gxravel/bus-routes-visualizer/internal/api/http"
	"github.com/gxravel/bus-routes-visualizer/internal/config"
	ierr "github.com/gxravel/bus-routes-visualizer/internal/errors"
	log "github.com/gxravel/bus-routes-visualizer/internal/logger"
	"github.com/gxravel/bus-routes-visualizer/internal/visualizercontext"
	"github.com/pkg/errors"
)

type httpClient struct {
	*http.Client

	statusHandlers map[int]func(*http.Response) error
	headers        map[string]string
}

const (
	defaultTimeout = time.Second * 30
)

type httpOpt func(*httpClient) error

func newCustomClient(opts ...httpOpt) *httpClient {
	c := &httpClient{
		Client: &http.Client{
			Timeout: defaultTimeout,
		},
		statusHandlers: make(map[int]func(*http.Response) error),
	}

	for _, opt := range opts {
		if err := opt(c); err != nil {
			panic(err)
		}
	}

	return c
}

func withTimeout(cfg *config.Config) httpOpt {
	return func(doer *httpClient) error {
		doer.Client = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   cfg.RemoteServices.DefaultTimeout,
					KeepAlive: cfg.RemoteServices.DefaultTimeout,
				}).DialContext,
				IdleConnTimeout:       cfg.RemoteServices.DefaultTimeout,
				TLSHandshakeTimeout:   cfg.RemoteServices.DefaultTimeout,
				ExpectContinueTimeout: 1 * time.Second,

				MaxConnsPerHost:     cfg.RemoteServices.DefaultMaxConns,
				MaxIdleConns:        cfg.RemoteServices.DefaultMaxConns,
				MaxIdleConnsPerHost: cfg.RemoteServices.DefaultMaxConns,
			},
			Timeout: cfg.RemoteServices.DefaultTimeout,
		}
		return nil
	}
}

// nolint
func withHeaders(h map[string]string) httpOpt {
	return func(doer *httpClient) error {
		doer.headers = h
		return nil
	}
}

func withUseInsecureTLS(conf *config.Config) httpOpt {
	logger := log.Default()
	return func(doer *httpClient) error {
		tr, ok := doer.Client.Transport.(*http.Transport)
		if !ok {
			logger.Warn("transport is not initialized")
			return nil
		}

		tlsConf := tr.TLSClientConfig
		if tlsConf == nil {
			tlsConf = &tls.Config{} //nolint
			tr.TLSClientConfig = tlsConf
		}

		if conf.RemoteServices.SkipTLSVerify {
			logger.Warn("skip tls verification for http client")
			tlsConf.InsecureSkipVerify = true
		}

		return nil
	}
}

func (c *httpClient) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	for k, v := range c.headers {
		req.Header.Set(k, v)
	}

	logger := log.FromContext(ctx).WithFields(
		"url", req.URL,
		"method", req.Method,
	)

	defer func(start time.Time) {
		logger.WithField("duration", time.Since(start)).Debug("call to http client response")
	}(time.Now())

	if req.Body != nil {
		if req.Header.Get(api.HeaderContentType) == "" {
			req.Header.Set(api.HeaderContentType, api.MIMEApplicationJSON.String())
		}
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		logger = logger.WithErr(err)
		return nil, errors.Wrapf(err, "can not DO call to %s: ", req.URL.String())
	}

	// for defer
	logger = logger.WithFields(
		"code", resp.Status,
	)

	return resp, nil
}

func (c *httpClient) processResponseStatuses(req *http.Request, resp *http.Response) error {
	switch {
	case resp.StatusCode >= 200 && resp.StatusCode < 300:
		return nil
	default:
		if handler, ok := c.statusHandlers[resp.StatusCode]; ok && handler != nil {
			return handler(resp)
		}

		return ierr.NewProviderAPIError(
			fmt.Sprintf("wrong status: %s when calling %s", resp.Status, req.URL),
			resp.StatusCode,
		)
	}
}

func (c *httpClient) DoJSON(ctx context.Context, req *http.Request, result interface{}) error {
	resp, err := c.Do(ctx, req)
	if err != nil {
		return err
	}

	if err = c.processResponseStatuses(req, resp); err != nil {
		return err
	}

	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(result)
	return errors.Wrap(err, "failed to decode response: ")
}

func (c *httpClient) processRequest(ctx context.Context, method, url string, body, result interface{}) error {
	if ctx == nil {
		ctx = context.Background()
	}
	if method == "" {
		return errors.New("empty method in processRequest")
	}
	if url == "" {
		return errors.New("empty url in processRequest")
	}

	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return errors.Wrap(err, "preparing request")
	}
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return ierr.NewReason(ierr.ErrBadRequest)
		}
		req.Body = io.NopCloser(bytes.NewReader(b))
	}

	setAuthToken(ctx, req)

	if result != nil {
		return c.DoJSON(ctx, req.WithContext(ctx), result)
	}

	resp, err := c.Do(ctx, req.WithContext(ctx))
	if err != nil {
		return err
	}
	if err = resp.Body.Close(); err != nil {
		return err
	}

	return c.processResponseStatuses(req, resp)
}

func setAuthToken(ctx context.Context, r *http.Request) {
	if token := visualizercontext.GetToken(ctx); token != "" {
		r.Header.Set("Authorization", "Bearer "+token)
	}
}
