package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cortezaproject/corteza/server/pkg/options"
	"github.com/stretchr/testify/require"
)

func Test_WellKnownOpenIDConfiguration(t *testing.T) {
	var (
		req = require.New(t)

		svc = service{opt: options.AuthOpt{BaseURL: "https://example.tld/auth"}}
		rec = httptest.NewRecorder()
	)

	svc.WellKnownOpenIDConfiguration()(rec, httptest.NewRequest(http.MethodGet, "/.well-known/openid-configuration", nil))

	var doc map[string]interface{}
	req.NoError(json.Unmarshal(rec.Body.Bytes(), &doc))

	// issuer must match the "iss" claim of the issued ID tokens
	req.Equal("https://example.tld/auth", doc["issuer"])
	req.Equal("https://example.tld/auth/oauth2/authorize", doc["authorization_endpoint"])
	req.Equal("https://example.tld/auth/oauth2/token", doc["token_endpoint"])
	req.Equal("https://example.tld/auth/oauth2/userinfo", doc["userinfo_endpoint"])
}
