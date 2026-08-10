package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testResource() OAuthProtectedResource {
	return OAuthProtectedResource{
		BasePath: "/api/mcp",
		Issuer:   "https://human.example.tld/auth",
		Scopes:   []string{"api", "profile"},
	}
}

// metadata serves the document for path and decodes it.
func metadata(t *testing.T, o OAuthProtectedResource, path string, header map[string]string) (int, map[string]any) {
	t.Helper()

	r := httptest.NewRequest(http.MethodGet, path, nil)
	r.Host = "human.example.tld"
	for k, v := range header {
		r.Header.Set(k, v)
	}

	w := httptest.NewRecorder()
	o.MetadataHandler().ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		return w.Code, nil
	}

	out := make(map[string]any)
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &out))
	return w.Code, out
}

// TestProtectedResourceEcho covers the requirement the rest of the handshake
// hangs off: the resource identifier has to come back as the client asked for
// it, path and all, because the client compares the two.
func TestProtectedResourceEcho(t *testing.T) {
	var (
		o = testResource()

		cases = []struct {
			name     string
			path     string
			resource string
		}{
			{"bare probe answers for the base path",
				WellKnownProtectedResource,
				"http://human.example.tld/api/mcp"},
			{"exact surface",
				WellKnownProtectedResource + "/api/mcp",
				"http://human.example.tld/api/mcp"},
			{"group-narrowed surface",
				WellKnownProtectedResource + "/api/mcp/configuring",
				"http://human.example.tld/api/mcp/configuring"},
			{"trailing slash is not part of the identifier",
				WellKnownProtectedResource + "/api/mcp/usage/",
				"http://human.example.tld/api/mcp/usage"},
		}
	)

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			code, doc := metadata(t, o, c.path, nil)
			require.Equal(t, http.StatusOK, code)

			assert.Equal(t, c.resource, doc["resource"])
			assert.Equal(t, []any{"https://human.example.tld/auth"}, doc["authorization_servers"])
			assert.Equal(t, []any{"api", "profile"}, doc["scopes_supported"])
			assert.Equal(t, []any{"header"}, doc["bearer_methods_supported"])
		})
	}
}

// TestProtectedResourceRefusesForeignPaths keeps the handler from declaring
// paths it does not protect. Everything unrouted on this server answers 200, so
// a client probing a foreign path would otherwise be told it is an OAuth
// resource and sent off to authenticate for nothing.
func TestProtectedResourceRefusesForeignPaths(t *testing.T) {
	o := testResource()

	for _, path := range []string{
		WellKnownProtectedResource + "/api/system",
		WellKnownProtectedResource + "/api/mcpsomething",
		WellKnownProtectedResource + "/",
	} {
		t.Run(path, func(t *testing.T) {
			code, _ := metadata(t, o, path, nil)

			// "/" trims to empty, which is the bare probe and is answered.
			if path == WellKnownProtectedResource+"/" {
				assert.Equal(t, http.StatusOK, code)
				return
			}

			assert.Equal(t, http.StatusNotFound, code)
		})
	}
}

// TestProtectedResourceOrigin covers the deployment shapes. Publishing the
// wrong origin produces a document that parses and cannot be connected to.
func TestProtectedResourceOrigin(t *testing.T) {
	var cases = []struct {
		name          string
		sslTerminated bool
		header        map[string]string
		resource      string
	}{
		{"plain http", false, nil,
			"http://human.example.tld/api/mcp"},
		{"tls terminated at a proxy that says nothing", true, nil,
			"https://human.example.tld/api/mcp"},
		{"proxy header wins over the option", false,
			map[string]string{"X-Forwarded-Proto": "https"},
			"https://human.example.tld/api/mcp"},
		{"proxy header wins the other way too", true,
			map[string]string{"X-Forwarded-Proto": "http"},
			"http://human.example.tld/api/mcp"},
		{"forwarded host", false,
			map[string]string{"X-Forwarded-Host": "mcp.public.tld"},
			"http://mcp.public.tld/api/mcp"},
		{"first hop of a chained header", false,
			map[string]string{"X-Forwarded-Proto": "https, http", "X-Forwarded-Host": "mcp.public.tld, internal:8080"},
			"https://mcp.public.tld/api/mcp"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			o := testResource()
			o.SslTerminated = c.sslTerminated

			code, doc := metadata(t, o, WellKnownProtectedResource+"/api/mcp", c.header)
			require.Equal(t, http.StatusOK, code)
			assert.Equal(t, c.resource, doc["resource"])
		})
	}
}

// challenged runs the middleware over a handler that answers with status.
func challenged(t *testing.T, path string, status int) *httptest.ResponseRecorder {
	t.Helper()

	r := httptest.NewRequest(http.MethodPost, path, nil)
	r.Host = "human.example.tld"

	w := httptest.NewRecorder()
	testResource().Challenge()(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(status)
	})).ServeHTTP(w, r)

	return w
}

func TestChallengeOnUnauthorized(t *testing.T) {
	t.Run("401 carries the pointer to its own resource path", func(t *testing.T) {
		w := challenged(t, "/api/mcp/configuring", http.StatusUnauthorized)

		assert.Equal(t,
			`Bearer resource_metadata="http://human.example.tld/.well-known/oauth-protected-resource/api/mcp/configuring", scope="api profile"`,
			w.Header().Get("WWW-Authenticate"),
		)
	})

	t.Run("anything else is left alone", func(t *testing.T) {
		// A pointer on a 200 is ignored by clients and misleading to read.
		for _, status := range []int{http.StatusOK, http.StatusForbidden, http.StatusInternalServerError} {
			w := challenged(t, "/api/mcp", status)
			assert.Empty(t, w.Header().Get("WWW-Authenticate"), "status %d", status)
		}
	})
}

// TestChallengeWriterStreams guards the streamable HTTP transport. Wrapping a
// ResponseWriter hides whatever interfaces the original implemented, and an MCP
// server that cannot flush cannot stream.
func TestChallengeWriterStreams(t *testing.T) {
	var (
		w       = httptest.NewRecorder()
		wrapped = &challengeWriter{ResponseWriter: w, challenge: "irrelevant"}
	)

	f, ok := any(wrapped).(http.Flusher)
	require.True(t, ok, "challengeWriter must implement http.Flusher")

	_, _ = wrapped.Write([]byte("event: hello\n\n"))
	f.Flush()

	assert.True(t, w.Flushed, "flush must reach the underlying writer")
	assert.Equal(t, w, wrapped.Unwrap(), "ResponseController must be able to unwrap")
}
