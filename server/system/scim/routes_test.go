package scim

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cortezaproject/corteza/server/pkg/options"
	"github.com/stretchr/testify/require"
)

func TestGuard(t *testing.T) {
	var (
		next = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })

		call = func(secret, header string) int {
			req := httptest.NewRequest(http.MethodGet, "/scim/Users", nil)
			if header != "" {
				req.Header.Set("Authorization", header)
			}

			rec := httptest.NewRecorder()
			Guard(options.SCIMOpt{Secret: secret})(next).ServeHTTP(rec, req)
			return rec.Code
		}
	)

	require.Equal(t, http.StatusOK, call("s3cret", "Bearer s3cret"))

	require.Equal(t, http.StatusForbidden, call("s3cret", ""))
	require.Equal(t, http.StatusForbidden, call("s3cret", "Bearer wrong!"))
	require.Equal(t, http.StatusForbidden, call("s3cret", "Bearer s3cret "))
	require.Equal(t, http.StatusForbidden, call("s3cret", "Basic  s3cret"), "prefix must be checked")

	// secret that is not set must not let anyone in
	require.Equal(t, http.StatusForbidden, call("", ""))
	require.Equal(t, http.StatusForbidden, call("", "Bearer "))
	require.Equal(t, http.StatusForbidden, call("", "1234567"), "any 7 characters used to pass")
}
