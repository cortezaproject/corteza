package appstore

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Client abstracts communication with the appstore middleware.
type Client interface {
	// ListPage calls GET /v1/connections?page=N&limit=M — returns one page of lightweight summaries, no config blobs.
	ListPage(ctx context.Context, page, limit int) ([]ConnectionSummary, error)

	// GetConnection calls GET /v1/connections/{id}/config — returns the full connection with all config blobs.
	GetConnection(ctx context.Context, id string) (*Connection, error)
}

type HTTPClient struct {
	baseURL string
	apiKey  string
	http    *http.Client
}

// New returns an HTTPClient. Returns nil when url is empty (disables catalog integration).
func New(url, apiKey string) Client {
	if url == "" {
		return nil
	}
	return &HTTPClient{
		baseURL: url,
		apiKey:  apiKey,
		http:    &http.Client{Timeout: 15 * time.Second},
	}
}

func (c *HTTPClient) ListPage(ctx context.Context, page, limit int) ([]ConnectionSummary, error) {
	url := fmt.Sprintf("%s/v1/connections?page=%d&limit=%d", c.baseURL, page, limit)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("appstore returned %d", resp.StatusCode)
	}

	var out []ConnectionSummary
	if err = json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *HTTPClient) GetConnection(ctx context.Context, id string) (*Connection, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/v1/connections/"+id+"/config", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("connection %q not found in appstore", id)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("appstore returned %d", resp.StatusCode)
	}

	var out Connection
	if err = json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}
