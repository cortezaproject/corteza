package handlers

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_requiresFormPostValidation(t *testing.T) {
	rq := require.New(t)

	rq.True(requiresFormPostValidation(http.MethodPost))
	rq.True(requiresFormPostValidation(http.MethodPut))
	rq.True(requiresFormPostValidation(http.MethodPatch))

	rq.False(requiresFormPostValidation(http.MethodGet))
	rq.False(requiresFormPostValidation(http.MethodHead))
	rq.False(requiresFormPostValidation(http.MethodDelete))
}

// Test_validFormPost_oidcCallbackQuery reproduces the Azure AD
// (openid-connect) callback that broke login: a GET to
// /auth/external/openid-connect.azure/callback whose `code` query
// param alone is 1472 bytes, well over maxPostValueLength (1024).
//
// validFormPost on its own still flags it - that's expected, and
// exactly why requiresFormPostValidation must gate it out for GET.
func Test_validFormPost_oidcCallbackQuery(t *testing.T) {
	rq := require.New(t)

	code := strings.Repeat("a", 1472)
	state := strings.Repeat("b", 88)
	sessionState := strings.Repeat("c", 36)

	req, err := http.NewRequest(
		http.MethodGet,
		"/auth/external/openid-connect.azure/callback?code="+code+"&state="+state+"&session_state="+sessionState,
		nil,
	)
	rq.NoError(err)
	rq.NoError(req.ParseForm())

	rq.False(requiresFormPostValidation(req.Method))
	rq.False(validFormPost(req))
}

func Test_validFormPost_postWithinLimits(t *testing.T) {
	rq := require.New(t)

	form := url.Values{}
	form.Set("username", "mock.user")
	form.Set("password", "secret")

	req, err := http.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(form.Encode()))
	rq.NoError(err)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rq.NoError(req.ParseForm())

	rq.True(requiresFormPostValidation(req.Method))
	rq.True(validFormPost(req))
}

func Test_validFormPost_postExceedingLimits(t *testing.T) {
	rq := require.New(t)

	form := url.Values{}
	form.Set("username", strings.Repeat("a", maxPostValueLength+1))

	req, err := http.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(form.Encode()))
	rq.NoError(err)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rq.NoError(req.ParseForm())

	rq.True(requiresFormPostValidation(req.Method))
	rq.False(validFormPost(req))
}
