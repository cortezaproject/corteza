package oauth2

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateRedirectURI(t *testing.T) {
	var (
		cases = []struct {
			name     string
			allowed  []string
			redirect string
			valid    bool
		}{
			{
				name:     "exact match",
				allowed:  []string{"https://app.example.com/callback"},
				redirect: "https://app.example.com/callback",
				valid:    true,
			},
			{
				name:     "origin allows any path",
				allowed:  []string{"https://app.example.com"},
				redirect: "https://app.example.com/compose/auth",
				valid:    true,
			},
			{
				name:     "origin with trailing slash allows any path",
				allowed:  []string{"https://app.example.com/"},
				redirect: "https://app.example.com/compose/auth",
				valid:    true,
			},
			{
				name:     "sub path",
				allowed:  []string{"https://app.example.com/callback"},
				redirect: "https://app.example.com/callback/deep",
				valid:    true,
			},
			{
				name:     "query string is ignored",
				allowed:  []string{"https://app.example.com/callback"},
				redirect: "https://app.example.com/callback?foo=bar",
				valid:    true,
			},
			{
				name:     "host is case insensitive",
				allowed:  []string{"https://app.example.com"},
				redirect: "https://APP.Example.com/callback",
				valid:    true,
			},
			{
				name:     "second of many",
				allowed:  []string{"https://one.example.com", "https://two.example.com"},
				redirect: "https://two.example.com/callback",
				valid:    true,
			},
			{
				name:     "wildcard allows anything",
				allowed:  []string{"*"},
				redirect: "https://anywhere.example.org/callback",
				valid:    true,
			},

			{
				name:     "nothing allowed",
				allowed:  nil,
				redirect: "https://app.example.com/callback",
			},
			{
				name:     "arbitrary domain",
				allowed:  []string{"http://localhost:8090"},
				redirect: "http://evil.com/steal",
			},
			{
				name:     "userinfo bypass",
				allowed:  []string{"http://localhost:8090"},
				redirect: "http://localhost:8090@evil.com/steal",
			},
			{
				name:     "host prefix bypass",
				allowed:  []string{"http://localhost:8090"},
				redirect: "http://localhost:8090.evil.com/steal",
			},
			{
				name:     "host suffix bypass",
				allowed:  []string{"https://app.example.com/callback"},
				redirect: "https://app.example.com.evil.com/callback",
			},
			{
				name:     "path prefix bypass",
				allowed:  []string{"https://app.example.com/callback"},
				redirect: "https://app.example.com/callbackevil",
			},
			{
				name:     "path traversal",
				allowed:  []string{"https://app.example.com/callback"},
				redirect: "https://app.example.com/callback/../evil",
			},
			{
				name:     "different scheme",
				allowed:  []string{"https://app.example.com"},
				redirect: "http://app.example.com/callback",
			},
			{
				name:     "different port",
				allowed:  []string{"http://localhost:8090"},
				redirect: "http://localhost:9999/callback",
			},
			{
				name:     "scheme relative",
				allowed:  []string{"https://app.example.com"},
				redirect: "//evil.com/callback",
			},
			{
				name:     "relative",
				allowed:  []string{"https://app.example.com"},
				redirect: "/callback",
			},
			{
				name:     "javascript scheme",
				allowed:  []string{"https://app.example.com"},
				redirect: "javascript:alert(1)",
			},
			{
				name:     "backslash host confusion",
				allowed:  []string{"https://app.example.com"},
				redirect: "https://evil.com\\@app.example.com/callback",
			},
			{
				name:     "empty",
				allowed:  []string{"https://app.example.com"},
				redirect: "",
			},
			{
				name:     "allowed without scheme never matches",
				allowed:  []string{"app.example.com"},
				redirect: "https://app.example.com/callback",
			},
		}
	)

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := ValidateRedirectURI(c.allowed, c.redirect)
			if c.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestAllowedRedirectURIs(t *testing.T) {
	var (
		def = []string{"https://corteza.example.com"}
	)

	// client without redirect URIs falls back to defaults instead of allowing everything
	require.Equal(t, def, AllowedRedirectURIs("", def))
	require.Equal(t, def, AllowedRedirectURIs("   ", def))

	// client with redirect URIs uses its own list
	require.Equal(t,
		[]string{"https://one.example.com", "https://two.example.com/cb"},
		AllowedRedirectURIs("https://one.example.com  https://two.example.com/cb", def),
	)
}
