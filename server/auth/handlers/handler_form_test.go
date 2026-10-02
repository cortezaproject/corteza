package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_parseRequestForm(t *testing.T) {
	long := strings.Repeat("x", maxPostValueLength*2)

	t.Run("long query value from an identity provider callback is accepted", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/auth/external/openid-connect.m365/callback?code="+long+"&state=abc", nil)

		tooLarge, err := parseRequestForm(r)
		require.NoError(t, err)
		require.False(t, tooLarge)
		require.Equal(t, long, r.Form.Get("code"))
	})

	t.Run("non multipart body is not an error", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader("email=a%40b.c"))
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		tooLarge, err := parseRequestForm(r)
		require.NoError(t, err)
		require.False(t, tooLarge)
		require.Equal(t, "a@b.c", r.PostForm.Get("email"))
	})

	t.Run("content type on a bodiless GET is not an error", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/auth/?code=abc", nil)
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		tooLarge, err := parseRequestForm(r)
		require.NoError(t, err)
		require.False(t, tooLarge)
	})

	t.Run("oversized posted field is reported without a parse error", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(url.Values{"email": {long}}.Encode()))
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		tooLarge, err := parseRequestForm(r)
		require.NoError(t, err)
		require.True(t, tooLarge)
	})

	t.Run("too many posted fields are reported", func(t *testing.T) {
		vv := url.Values{}
		for i := 0; i <= maxPostFields; i++ {
			vv.Set(strings.Repeat("f", i+1), "v")
		}
		r := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(vv.Encode()))
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		tooLarge, err := parseRequestForm(r)
		require.NoError(t, err)
		require.True(t, tooLarge)
	})
}
