package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cortezaproject/corteza/server/pkg/options"
	"github.com/stretchr/testify/require"
)

func TestHandleCORS(t *testing.T) {
	const (
		allowOrigin      = "Access-Control-Allow-Origin"
		allowCredentials = "Access-Control-Allow-Credentials"
		allowMethods     = "Access-Control-Allow-Methods"
	)

	tests := []struct {
		name      string
		origins   []string
		method    string
		origin    string
		preflight bool

		wantOrigin      string
		wantCredentials string
		wantMethods     string
		wantNext        bool
	}{
		{
			name:            "any origin reflected by default",
			origins:         options.CorsAnyOrigin(),
			method:          http.MethodGet,
			origin:          "https://evil.example",
			wantOrigin:      "https://evil.example",
			wantCredentials: "true",
			wantNext:        true,
		},
		{
			name:            "any origin preflight by default",
			origins:         options.CorsAnyOrigin(),
			method:          http.MethodOptions,
			origin:          "https://evil.example",
			preflight:       true,
			wantOrigin:      "https://evil.example",
			wantCredentials: "true",
			wantMethods:     http.MethodGet,
		},
		{
			name:            "listed origin reflected",
			origins:         []string{"https://crm.example.com"},
			method:          http.MethodPost,
			origin:          "https://crm.example.com",
			wantOrigin:      "https://crm.example.com",
			wantCredentials: "true",
			wantNext:        true,
		},
		{
			name:     "unlisted origin gets no headers but is still served",
			origins:  []string{"https://crm.example.com"},
			method:   http.MethodPost,
			origin:   "https://evil.example",
			wantNext: true,
		},
		{
			name:            "listed origin preflight",
			origins:         []string{"https://crm.example.com"},
			method:          http.MethodOptions,
			origin:          "https://crm.example.com",
			preflight:       true,
			wantOrigin:      "https://crm.example.com",
			wantCredentials: "true",
			wantMethods:     http.MethodGet,
		},
		{
			name:      "unlisted origin preflight",
			origins:   []string{"https://crm.example.com"},
			method:    http.MethodOptions,
			origin:    "https://evil.example",
			preflight: true,
		},
		{
			name:            "subdomain wildcard",
			origins:         []string{"https://*.example.com"},
			method:          http.MethodGet,
			origin:          "https://crm.example.com",
			wantOrigin:      "https://crm.example.com",
			wantCredentials: "true",
			wantNext:        true,
		},
		{
			name:     "subdomain wildcard does not match other domains",
			origins:  []string{"https://*.example.com"},
			method:   http.MethodGet,
			origin:   "https://example.org",
			wantNext: true,
		},
		{
			name:     "same-origin request without origin header",
			origins:  []string{"https://crm.example.com"},
			method:   http.MethodGet,
			wantNext: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				req    = require.New(t)
				called bool
				next   = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { called = true })
				rec    = httptest.NewRecorder()
				r      = httptest.NewRequest(tt.method, "/api/system/users", nil)
			)

			if tt.origin != "" {
				r.Header.Set("Origin", tt.origin)
			}

			if tt.preflight {
				r.Header.Set("Access-Control-Request-Method", http.MethodGet)
				r.Header.Set("Access-Control-Request-Headers", "Authorization")
			}

			handleCORS(tt.origins)(next).ServeHTTP(rec, r)

			req.Equal(tt.wantNext, called)
			req.Equal(tt.wantOrigin, rec.Header().Get(allowOrigin))
			req.Equal(tt.wantCredentials, rec.Header().Get(allowCredentials))
			req.Equal(tt.wantMethods, rec.Header().Get(allowMethods))
		})
	}
}
