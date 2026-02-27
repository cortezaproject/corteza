package gsheets

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/cortezaproject/corteza/server/store/adapters/api/cred_registry"
)

type (
	gsheetsWrapper struct {
		client       *http.Client
		baseURL      string
		connectionID uint64
	}
)

func newWrapper(baseURL string, connectionID uint64) *gsheetsWrapper {
	return &gsheetsWrapper{
		client:       &http.Client{},
		baseURL:      baseURL,
		connectionID: connectionID,
	}
}

// Run executes an HTTP request against the Sheets v4 API.
//
// For GET operations the Sheets response is
// {"values":[["h1","h2"],["v1","v2"],...]} (array of arrays, row 0 = headers).
// This method transforms it into [{"h1":"v1","h2":"v2"},...] so the
// existing apidal iterator/table codec can process it unchanged.
func (w *gsheetsWrapper) Run(ctx context.Context, method string, path string, payload []byte, extraHeaders map[string][]string) (statusCode int, outHeaders map[string][]string, rsp []byte, err error) {
	token, err := cred_registry.Default().GetAccessToken(ctx, w.connectionID)
	if err != nil {
		return 0, nil, nil, fmt.Errorf("failed to get access token: %w", err)
	}

	fullURL := w.baseURL + path

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
		return resp.StatusCode, resp.Header, nil, fmt.Errorf("sheets API returned %d: %s", resp.StatusCode, string(body))
	}

	// For GET operations, transform array-of-arrays into JSON objects
	if method == "GET" {
		body, err = transformSheetsResponse(body)
		if err != nil {
			return resp.StatusCode, resp.Header, nil, fmt.Errorf("failed to transform sheets response: %w", err)
		}
	}

	return resp.StatusCode, resp.Header, body, nil
}

// transformSheetsResponse converts:
//
//	{"values": [["h1","h2"], ["v1","v2"], ["v3","v4"]]}
//
// into:
//
//	[{"h1":"v1","h2":"v2"}, {"h1":"v3","h2":"v4"}]
func transformSheetsResponse(data []byte) ([]byte, error) {
	var resp struct {
		Values [][]any `json:"values"`
	}

	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse sheets response: %w", err)
	}

	if len(resp.Values) == 0 {
		return []byte("[]"), nil
	}

	// Row 0 = headers
	headers := make([]string, len(resp.Values[0]))
	for i, h := range resp.Values[0] {
		headers[i] = fmt.Sprintf("%v", h)
	}

	// Rows 1+ = data
	rows := make([]map[string]any, 0, len(resp.Values)-1)
	for _, row := range resp.Values[1:] {
		obj := make(map[string]any, len(headers))
		for i, h := range headers {
			if i < len(row) {
				obj[h] = row[i]
			} else {
				obj[h] = ""
			}
		}
		rows = append(rows, obj)
	}

	out := make([]map[string]any, len(rows))
	for i, row := range rows {
		out[i] = func(i int, row map[string]any) map[string]any {
			aux := map[string]any{
				"id": fmt.Sprintf("%d", i+2),
			}
			for k, v := range row {
				aux[k] = v
			}
			return aux
		}(i, row)
	}

	return json.Marshal(out)
}
