package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/cortezaproject/corteza-server/auth/request"
	"github.com/cortezaproject/corteza-server/auth/settings"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func Test_oauth2AuthorizeSuccess(t *testing.T) {
	var (
		ctx  = context.Background()
		user = makeMockUser(ctx)

		req = &http.Request{
			Form:     url.Values{},
			PostForm: url.Values{},
		}

		oauthService oauth2Service
		authService  authService
		authHandlers *AuthHandlers
		authReq      *request.AuthReq

		authSettings = &settings.Settings{}
	)

	tcc := []testingExpect{
		{
			name:     "authorization success",
			payload:  -1,
			err:      "",
			template: "",
			fn: func() {
				oauthService = &oauth2ServiceMocked{
					handleAuthorizeRequest: func(w http.ResponseWriter, r *http.Request) error {
						return nil
					},
				}
			},
		},
		{
			name:     "authorization failure",
			payload:  http.StatusInternalServerError,
			err:      "not authorized",
			template: TmplInternalError,
			fn: func() {
				oauthService = &oauth2ServiceMocked{
					handleAuthorizeRequest: func(w http.ResponseWriter, r *http.Request) error {
						return errors.New("not authorized")
					},
				}
			},
		},
	}

	for _, tc := range tcc {
		t.Run(tc.name, func(t *testing.T) {
			rq := require.New(t)

			tc.fn()

			authReq = prepareClientAuthReq(ctx, req, user)
			authHandlers = &AuthHandlers{
				Log:         zap.NewNop(),
				AuthService: authService,
				Settings:    authSettings,
				OAuth2:      oauthService,
			}

			err := authHandlers.oauth2Authorize(authReq)

			rq.NoError(err)
			rq.Equal(tc.template, authReq.Template)
			rq.Equal(tc.payload, authReq.Status)

			if tc.err != "" {
				rq.EqualError(errors.New(tc.err), authReq.Data["error"].(error).Error())
			}
		})
	}
}

func Test_oauth2AuthorizeSuccessSetParams(t *testing.T) {
	var (
		ctx  = context.Background()
		user = makeMockUser(ctx)

		req = &http.Request{
			Form:     url.Values{},
			PostForm: url.Values{},
		}

		authService  authService
		authHandlers *AuthHandlers
		authReq      *request.AuthReq

		authSettings = &settings.Settings{}

		rq = require.New(t)
	)

	oauthService := &oauth2ServiceMocked{
		handleAuthorizeRequest: func(w http.ResponseWriter, r *http.Request) error {
			return nil
		},
	}

	authReq = prepareClientAuthReq(ctx, req, user)
	authReq.Session.Values["oauth2AuthParams"] = url.Values{"foo": []string{"bar"}}

	authHandlers = &AuthHandlers{
		Log:         zap.NewNop(),
		AuthService: authService,
		Settings:    authSettings,
		OAuth2:      oauthService,
	}

	err := authHandlers.oauth2Authorize(authReq)

	rq.NoError(err)
	rq.Equal("", authReq.Template)
	rq.Equal(-1, authReq.Status)
	rq.Equal(nil, authReq.Data["error"])
}
