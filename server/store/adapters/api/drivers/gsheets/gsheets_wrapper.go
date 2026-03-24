package gsheets

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cortezaproject/corteza/server/store/adapters/api/drivers/google"
)

type sheetsWrapper struct {
	*google.Wrapper
}

func newWrapper(baseURL string, connectionID uint64) *sheetsWrapper {
	return &sheetsWrapper{
		Wrapper: google.NewWrapper(baseURL, connectionID),
	}
}

func (w *sheetsWrapper) Run(ctx context.Context, method string, path string, payload []byte, extraHeaders map[string][]string) (int, map[string][]string, []byte, error) {
	statusCode, headers, rsp, err := w.Wrapper.Run(ctx, method, path, payload, extraHeaders)
	if err != nil {
		return statusCode, headers, rsp, err
	}

	// For GET operations, transform array-of-arrays into JSON objects
	if method == "GET" && statusCode >= 200 && statusCode < 300 {
		rsp, err = transformSheetsResponse(rsp)
		if err != nil {
			return statusCode, headers, nil, fmt.Errorf("failed to transform sheets response: %w", err)
		}
	}

	return statusCode, headers, rsp, nil
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
