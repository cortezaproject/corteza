package options

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAuthOpt_GetDefaultRedirectURIs(t *testing.T) {
	t.Setenv("DOMAIN_WEBAPP", "app.example.com")
	t.Setenv("HTTP_SSL_TERMINATED", "true")

	t.Run("origins of base URL and webapp", func(t *testing.T) {
		o := AuthOpt{BaseURL: "https://api.example.com/auth"}
		require.Equal(t, []string{"https://api.example.com", "https://app.example.com"}, o.GetDefaultRedirectURIs())
	})

	t.Run("extra from env, no duplicates", func(t *testing.T) {
		o := AuthOpt{
			BaseURL:             "https://api.example.com/auth",
			DefaultRedirectURIs: "http://localhost:8080  https://api.example.com https://other.example.com/cb",
		}

		require.Equal(t,
			[]string{"https://api.example.com", "https://app.example.com", "http://localhost:8080", "https://other.example.com/cb"},
			o.GetDefaultRedirectURIs(),
		)
	})

	t.Run("invalid base URL is skipped", func(t *testing.T) {
		o := AuthOpt{BaseURL: "/auth"}
		require.Equal(t, []string{"https://app.example.com"}, o.GetDefaultRedirectURIs())
	})
}
