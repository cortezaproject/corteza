package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/cortezaproject/corteza/server/pkg/options"
	"github.com/cortezaproject/corteza/server/system/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type (
	redirectTestClientService struct {
		client *types.AuthClient
	}
)

func (s redirectTestClientService) Lookup(context.Context, interface{}) (*types.AuthClient, error) {
	return s.client, nil
}

func (redirectTestClientService) Confirmed(context.Context, uint64) (types.AuthConfirmedClientSet, error) {
	return nil, nil
}

func (redirectTestClientService) Revoke(context.Context, uint64, uint64) error {
	return nil
}

func Test_oauth2AuthorizeRedirectURI(t *testing.T) {
	var (
		opt = options.AuthOpt{BaseURL: "http://localhost:8090/auth"}

		tcc = []struct {
			name      string
			clientURI string
			redirect  string
			handled   bool
		}{
			{name: "own origin, client without redirect URI", redirect: "http://localhost:8090/compose/auth", handled: true},
			{name: "arbitrary domain, client without redirect URI", redirect: "http://evil.com/steal"},
			{name: "userinfo bypass, client without redirect URI", redirect: "http://localhost:8090@evil.com/steal"},
			{name: "host prefix bypass, client without redirect URI", redirect: "http://localhost:8090.evil.com/steal"},
			{name: "registered", clientURI: "https://app.example.com/callback", redirect: "https://app.example.com/callback", handled: true},
			{name: "registered host prefix bypass", clientURI: "https://app.example.com/callback", redirect: "https://app.example.com.evil.com/callback"},
			{name: "registered userinfo bypass", clientURI: "https://app.example.com", redirect: "https://app.example.com@evil.com/callback"},
			{name: "own origin is not allowed when client has redirect URI", clientURI: "https://app.example.com", redirect: "http://localhost:8090/compose/auth"},
		}
	)

	for _, tc := range tcc {
		t.Run(tc.name, func(t *testing.T) {
			var (
				rq      = require.New(t)
				handled bool

				h = &AuthHandlers{
					Log: zap.NewNop(),
					Opt: opt,
					ClientService: redirectTestClientService{client: &types.AuthClient{
						ID:          42,
						Handle:      "corteza-webapp",
						RedirectURI: tc.clientURI,
						Trusted:     true,
						Enabled:     true,
					}},
					OAuth2: &oauth2ServiceMocked{
						handleAuthorizeRequest: func(w http.ResponseWriter, r *http.Request) error {
							handled = true
							return nil
						},
					},
				}

				form = url.Values{
					"response_type": []string{"code"},
					"client_id":     []string{"corteza-webapp"},
					"redirect_uri":  []string{tc.redirect},
					"scope":         []string{"api"},
				}

				authReq = prepareClientAuthReq(h, &http.Request{Form: form, PostForm: url.Values{}}, makeMockUser())
			)

			rq.NoError(h.oauth2Authorize(authReq))
			rq.Equal(tc.handled, handled)

			if tc.handled {
				rq.Equal(-1, authReq.Status)
				return
			}

			// must end with an error page and never with a redirect
			rq.Equal(http.StatusBadRequest, authReq.Status)
			rq.Equal(TmplInternalError, authReq.Template)
			rq.Empty(authReq.RedirectTo)
			rq.Empty(authReq.Response.(*httptest.ResponseRecorder).Header().Get("Location"))
		})
	}
}
