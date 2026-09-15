package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSkipCORS(t *testing.T) {
	const widgetPath = "/api/widget/v1/"

	tests := []struct {
		name       string
		path       string
		method     string
		wantNext   bool
		wantOrigin string
	}{
		{
			name:       "preflight outside widget routes answered by global CORS",
			path:       "/api/system/users",
			method:     http.MethodOptions,
			wantOrigin: "https://crm.example.com",
		},
		{
			name:       "request outside widget routes gets global CORS headers",
			path:       "/api/system/users",
			method:     http.MethodGet,
			wantNext:   true,
			wantOrigin: "https://crm.example.com",
		},
		{
			name:     "preflight on widget routes passed to widget handler",
			path:     "/api/widget/v1/session",
			method:   http.MethodOptions,
			wantNext: true,
		},
		{
			name:     "request on widget routes skips global CORS headers",
			path:     "/api/widget/v1/config",
			method:   http.MethodGet,
			wantNext: true,
		},
		{
			name:       "similar prefix is not skipped",
			path:       "/api/widget/v10/config",
			method:     http.MethodGet,
			wantNext:   true,
			wantOrigin: "https://crm.example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				req    = require.New(t)
				called bool
				next   = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { called = true })
				rec    = httptest.NewRecorder()
				r      = httptest.NewRequest(tt.method, tt.path, nil)
			)

			r.Header.Set("Origin", "https://crm.example.com")
			if tt.method == http.MethodOptions {
				r.Header.Set("Access-Control-Request-Method", http.MethodPost)
			}

			skipCORS(widgetPath, handleCORS([]string{"https://crm.example.com"}))(next).ServeHTTP(rec, r)

			req.Equal(tt.wantNext, called)
			req.Equal(tt.wantOrigin, rec.Header().Get("Access-Control-Allow-Origin"))
		})
	}
}
