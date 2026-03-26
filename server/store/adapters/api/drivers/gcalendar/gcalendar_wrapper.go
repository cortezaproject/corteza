package gcalendar

import (
	"context"
	"encoding/json"

	"github.com/cortezaproject/corteza/server/store/adapters/api/drivers/google"
)

type calendarWrapper struct {
	*google.Wrapper
}

func newWrapper(baseURL string, connectionID uint64) *calendarWrapper {
	return &calendarWrapper{
		Wrapper: google.NewWrapper(baseURL, connectionID),
	}
}

func (w *calendarWrapper) Run(ctx context.Context, method string, path string, payload []byte, extraHeaders map[string][]string) (int, map[string][]string, []byte, error) {
	statusCode, headers, rsp, err := w.Wrapper.Run(ctx, method, path, payload, extraHeaders)
	if err != nil {
		return statusCode, headers, rsp, err
	}

	// For GET operations, pull `items` array to root out if the JSON has one
	if method == "GET" && statusCode >= 200 && statusCode < 300 {
		var peek struct {
			Items json.RawMessage `json:"items"`
		}
		if err := json.Unmarshal(rsp, &peek); err == nil && len(peek.Items) > 0 {
			rsp = peek.Items
		}
	}

	return statusCode, headers, rsp, nil
}
