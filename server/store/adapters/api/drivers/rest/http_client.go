package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/dal"
)

// httpClient wraps http.httpClient with context support and convenience methods
type httpClient struct {
	httpClient *http.Client
	baseURL    string
	headers    map[string][]string
}

// ClientConfig holds configuration for creating a new Client
type ClientConfig struct {
	BaseURL             url.URL
	Timeout             time.Duration
	MaxIdleConns        int
	MaxIdleConnsPerHost int
	IdleConnTimeout     time.Duration
	Headers             map[string][]string

	DSN dal.DSN
}

// newClient creates a new HTTP client wrapper with the given configuration
func newClient(config ClientConfig) *httpClient {
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}
	if config.MaxIdleConns == 0 {
		config.MaxIdleConns = 100
	}
	if config.MaxIdleConnsPerHost == 0 {
		config.MaxIdleConnsPerHost = 10
	}
	if config.IdleConnTimeout == 0 {
		config.IdleConnTimeout = 90 * time.Second
	}

	transport := &http.Transport{
		MaxIdleConns:          config.MaxIdleConns,
		MaxIdleConnsPerHost:   config.MaxIdleConnsPerHost,
		IdleConnTimeout:       config.IdleConnTimeout,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	// @todo should this be moved down tho where auth headers are handled?
	if strings.ToLower(config.DSN.AuthType) == "basic" {
		config.BaseURL.User = url.UserPassword(config.DSN.Username, config.DSN.Password)
	}

	return &httpClient{
		httpClient: &http.Client{
			Transport: transport,
		},
		baseURL: config.BaseURL.String(),
		headers: config.Headers,
	}
}

// buildURL constructs the full URL from base and path
func (c *httpClient) buildURL(path string) (string, error) {
	if c.baseURL == "" {
		return path, nil
	}

	base, err := url.Parse(c.baseURL)
	if err != nil {
		return "", fmt.Errorf("invalid base URL: %w", err)
	}

	rel, err := url.Parse(path)
	if err != nil {
		return "", fmt.Errorf("invalid path: %w", err)
	}

	return base.ResolveReference(rel).String(), nil
}

// Do executes an HTTP request with context support
func (c *httpClient) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	if c.headers == nil {
		c.headers = map[string][]string{}
	}

	// Apply default headers
	for k, vv := range c.headers {
		if req.Header.Get(k) == "" {
			for _, v := range vv {
				req.Header.Add(k, v)
			}
		}
	}

	if _, ok := c.headers["Content-Type"]; !ok {
		req.Header.Set("Content-Type", "application/json")
	}

	// Ensure context is attached
	req = req.WithContext(ctx)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		if ctx.Err() == context.Canceled {
			return nil, fmt.Errorf("request cancelled: %w", err)
		}
		if ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("request timeout: %w", err)
		}
		return nil, err
	}

	return resp, nil
}

// Get performs a GET request
func (c *httpClient) Get(ctx context.Context, path string, headers map[string][]string) (*http.Response, error) {
	fullURL, err := c.buildURL(path)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, err
	}

	for k, vv := range headers {
		for _, v := range vv {
			req.Header.Add(k, v)
		}
	}

	return c.Do(ctx, req)
}

// Post performs a POST request with a body
func (c *httpClient) Post(ctx context.Context, path string, body interface{}, headers map[string][]string) (*http.Response, error) {
	return c.doWithBody(ctx, http.MethodPost, path, body, headers)
}

// Put performs a PUT request with a body
func (c *httpClient) Put(ctx context.Context, path string, body interface{}, headers map[string][]string) (*http.Response, error) {
	return c.doWithBody(ctx, http.MethodPut, path, body, headers)
}

// Patch performs a PATCH request with a body
func (c *httpClient) Patch(ctx context.Context, path string, body interface{}, headers map[string][]string) (*http.Response, error) {
	return c.doWithBody(ctx, http.MethodPatch, path, body, headers)
}

// Delete performs a DELETE request
func (c *httpClient) Delete(ctx context.Context, path string, headers map[string][]string) (*http.Response, error) {
	fullURL, err := c.buildURL(path)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, fullURL, nil)
	if err != nil {
		return nil, err
	}

	for k, vv := range headers {
		for _, v := range vv {
			req.Header.Set(k, v)
		}
	}

	return c.Do(ctx, req)
}

// Head performs a HEAD request
func (c *httpClient) Head(ctx context.Context, path string, headers map[string][]string) (*http.Response, error) {
	fullURL, err := c.buildURL(path)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodHead, fullURL, nil)
	if err != nil {
		return nil, err
	}

	for k, vv := range headers {
		for _, v := range vv {
			req.Header.Set(k, v)
		}
	}

	return c.Do(ctx, req)
}

// Options performs an OPTIONS request
func (c *httpClient) Options(ctx context.Context, path string, headers map[string][]string) (*http.Response, error) {
	fullURL, err := c.buildURL(path)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodOptions, fullURL, nil)
	if err != nil {
		return nil, err
	}

	for k, vv := range headers {
		for _, v := range vv {
			req.Header.Set(k, v)
		}
	}

	return c.Do(ctx, req)
}

// doWithBody is a helper for methods that send a body
func (c *httpClient) doWithBody(ctx context.Context, method, path string, body interface{}, headers map[string][]string) (*http.Response, error) {
	fullURL, err := c.buildURL(path)
	if err != nil {
		return nil, err
	}

	var bodyReader io.Reader
	if body != nil {
		switch v := body.(type) {
		case []byte:
			bodyReader = bytes.NewReader(v)
		case string:
			bodyReader = bytes.NewReader([]byte(v))
		case io.Reader:
			bodyReader = v
		default:
			// Assume JSON serialization
			jsonData, err := json.Marshal(body)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal body: %w", err)
			}
			bodyReader = bytes.NewReader(jsonData)
			if headers == nil {
				headers = make(map[string][]string)
			}
			if len(headers["Content-Type"]) == 0 {
				headers["Content-Type"] = []string{"application/json"}
			}
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, fullURL, bodyReader)
	if err != nil {
		return nil, err
	}

	for k, vv := range headers {
		for _, v := range vv {
			req.Header.Set(k, v)
		}
	}

	return c.Do(ctx, req)
}

// SetHeader sets a default header for all requests
func (c *httpClient) SetHeader(key, value string) {
	if c.headers == nil {
		c.headers = make(map[string][]string)
	}

	if c.headers[key] == nil {
		c.headers[key] = []string{}
	}

	c.headers[key] = append(c.headers[key], value)
}

// SetAuth sets the Authorization header
func (c *httpClient) SetAuth(authType, token string) {
	c.SetHeader("Authorization", fmt.Sprintf("%s %s", authType, token))
}

// SetBearerToken sets Bearer token authentication
func (c *httpClient) SetBearerToken(token string) {
	c.SetAuth("Bearer", token)
}

// Close closes idle connections
func (c *httpClient) Close() {
	c.httpClient.CloseIdleConnections()
}
