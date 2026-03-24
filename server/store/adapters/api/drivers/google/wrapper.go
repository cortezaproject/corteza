package google

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/cortezaproject/corteza/server/store/adapters/api/cred_registry"
)

type (
	Wrapper struct {
		client       *http.Client
		baseURL      string
		connectionID uint64
	}
)

func NewWrapper(baseURL string, connectionID uint64) *Wrapper {
	transport := &http.Transport{
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   10,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	return &Wrapper{
		client: &http.Client{
			Transport: transport,
			Timeout:   30 * time.Second,
		},
		baseURL:      baseURL,
		connectionID: connectionID,
	}
}

// Run executes an HTTP request against the Google API.
func (w *Wrapper) Run(ctx context.Context, method string, path string, payload []byte, extraHeaders map[string][]string) (statusCode int, outHeaders map[string][]string, rsp []byte, err error) {
	token, err := cred_registry.Default().GetAccessToken(ctx, w.connectionID)
	if err != nil {
		return 0, nil, nil, fmt.Errorf("failed to get access token: %w", err)
	}

	baseURL := strings.TrimSuffix(w.baseURL, "/")
	if path != "" && !strings.HasPrefix(path, "/") && !strings.HasPrefix(path, "?") {
		path = "/" + path
	}
	fullURL := baseURL + path

	var bodyReader io.Reader
	if payload != nil {
		bodyReader = bytes.NewReader(payload)
	}

	req, err := http.NewRequestWithContext(ctx, method, fullURL, bodyReader)
	if err != nil {
		return 0, nil, nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")

	for k, vv := range extraHeaders {
		for _, v := range vv {
			req.Header.Add(k, v)
		}
	}

	resp, err := w.client.Do(req)
	if err != nil {
		return 0, nil, nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, resp.Header, nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return resp.StatusCode, resp.Header, nil, fmt.Errorf("google API returned %d: %s", resp.StatusCode, string(body))
	}

	return resp.StatusCode, resp.Header, body, nil
}
