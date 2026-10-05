package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testSecret = "this-is-a-test-secret-long-enough"

// ok is the handler the middleware chain protects. Reaching it means the
// request was allowed through.
func ok(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusTeapot) }

func serve(t *testing.T, mw func(http.Handler) http.Handler, authorization string) int {
	t.Helper()

	verifier, err := TokenVerifierMiddlewareWithSecretSigner(testSecret)
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

// signedAt returns a bearer token signed with the test secret whose iat lies
// `offset` away from now.
func signedAt(t *testing.T, offset time.Duration) string {
	t.Helper()

	token := jwt.New()
	require.NoError(t, token.Set(jwt.JwtIDKey, "test-access"))
	require.NoError(t, token.Set(jwt.SubjectKey, "1"))
	require.NoError(t, token.Set(jwt.IssuedAtKey, time.Now().Add(offset).Unix()))
	require.NoError(t, token.Set(jwt.ExpirationKey, time.Now().Add(time.Hour).Unix()))
	require.NoError(t, token.Set("scope", "api"))

	signed, err := jwt.Sign(token, jwa.HS512, []byte(testSecret))
	require.NoError(t, err)
	return "Bearer " + string(signed)
}

// TestVerifierAllowsClockSkew covers a token whose iat is slightly ahead of the
// verifying clock, as when time sync steps the clock back right after the token
// was issued: it is accepted, while one issued well in the future is not.
func TestVerifierAllowsClockSkew(t *testing.T) {
	prev := TokenIssuer
	t.Cleanup(func() { TokenIssuer = prev })

	var err error
	TokenIssuer, err = NewTokenIssuer(WithLookup(func(context.Context, string) error { return nil }))
	require.NoError(t, err)

	t.Run("iat a few seconds ahead", func(t *testing.T) {
		code := serve(t, HttpTokenValidator("api"), signedAt(t, 3*time.Second))
		assert.Equal(t, http.StatusTeapot, code)
	})

	t.Run("iat far ahead", func(t *testing.T) {
		code := serve(t, HttpTokenValidator("api"), signedAt(t, 10*time.Minute))
		assert.Equal(t, http.StatusUnauthorized, code)
	})
}
