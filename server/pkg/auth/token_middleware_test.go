package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ok is the handler the middleware chain protects. Reaching it means the
// request was allowed through.
func ok(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusTeapot) }

func serve(t *testing.T, mw func(http.Handler) http.Handler, authorization string) int {
	t.Helper()

	verifier, err := TokenVerifierMiddlewareWithSecretSigner("this-is-a-test-secret-long-enough")
	require.NoError(t, err)

	r := httptest.NewRequest(http.MethodGet, "/api/whatever", nil)
	if authorization != "" {
		r.Header.Set("Authorization", authorization)
	}

	w := httptest.NewRecorder()
	verifier(mw(http.HandlerFunc(ok))).ServeHTTP(w, r)
	return w.Code
}

// TestHttpTokenValidatorRejectsMalformedToken pins the fix for a token that
// fails to parse being reported as a server fault.
//
// jwtauth returns a raw parse error, which carries no Kind, so
// ProperlyServeHTTP rendered it as 500. A credential that does not verify is
// rejected, not broken.
func TestHttpTokenValidatorRejectsMalformedToken(t *testing.T) {
	code := serve(t, HttpTokenValidator("api"), "Bearer not-a-real-token")
	assert.Equal(t, http.StatusUnauthorized, code, "a malformed token is 401, not 500")
}

// TestHttpTokenValidatorPassesAnonymous documents behaviour that looks like a
// hole and is not: the validator checks the scope of a token that is present,
// it does not require one. REST relies on that — its handlers reach RBAC, which
// denies the anonymous role. Anywhere that discloses something before RBAC runs
// needs HttpAuthenticatedOnly as well; see /api/mcp.
func TestHttpTokenValidatorPassesAnonymous(t *testing.T) {
	code := serve(t, HttpTokenValidator("api"), "")
	assert.Equal(t, http.StatusTeapot, code, "no token reaches the handler by design")
}

func TestHttpAuthenticatedOnly(t *testing.T) {
	t.Run("rejects a request with no identity", func(t *testing.T) {
		code := serve(t, HttpAuthenticatedOnly(), "")
		assert.Equal(t, http.StatusUnauthorized, code)
	})

	t.Run("rejects a malformed token", func(t *testing.T) {
		code := serve(t, HttpAuthenticatedOnly(), "Bearer not-a-real-token")
		assert.Equal(t, http.StatusUnauthorized, code)
	})
}
