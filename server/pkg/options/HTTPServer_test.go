package options

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHttpServerOpt_GetCorsAllowedOrigins(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  []string
	}{
		{"default", "http://*,https://*", CorsAnyOrigin()},
		{"empty", "", CorsAnyOrigin()},
		{"whitespace and commas", " , ,", CorsAnyOrigin()},
		{"wildcard", "*", CorsAnyOrigin()},
		{"wildcard among origins", "https://a.example.com,*", CorsAnyOrigin()},
		{"single", "https://example.com", []string{"https://example.com"}},
		{"multiple", "https://a.example.com,https://*.b.example.com", []string{"https://a.example.com", "https://*.b.example.com"}},
		{"trimmed", " https://a.example.com/ , https://b.example.com ,", []string{"https://a.example.com", "https://b.example.com"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := HttpServerOpt{CorsAllowedOrigins: tt.value}
			require.Equal(t, tt.want, o.GetCorsAllowedOrigins())
		})
	}
}

func TestHttpServerOpt_CorsAllowedOriginsDefault(t *testing.T) {
	// make sure env is unset while restoring it after the test
	t.Setenv("HTTP_CORS_ALLOWED_ORIGINS", "")
	require.NoError(t, os.Unsetenv("HTTP_CORS_ALLOWED_ORIGINS"))

	require.Equal(t, CorsAnyOrigin(), HttpServer().GetCorsAllowedOrigins())
}
