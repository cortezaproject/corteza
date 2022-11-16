package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/cortezaproject/corteza/server/system/service"

	"github.com/cortezaproject/corteza/server/auth/request"
	"github.com/cortezaproject/corteza/server/auth/settings"
	"github.com/cortezaproject/corteza/server/system/types"
	"github.com/stretchr/testify/require"
)

func Test_createPasswordForm(t *testing.T) {
	var (
		req = &http.Request{
			URL: &url.URL{},
		}

		user         *types.User
		authService  authService
		authHandlers *AuthHandlers
		authReq      *request.AuthReq
	)

	tcc := []testingExpect{
		{
			name:     "password create success",
			payload:  map[string]string(nil),
			alerts:   []request.Alert(nil),
			link:     GetLinks().CreatePassword,
			template: TmplCreatePassword,
			fn: func(_ *settings.Settings) {
				req.URL = &url.URL{RawQuery: "token=NOT_EMPTY"}

				authService = &authServiceMocked{
					passwordSet: func(ctx context.Context, s string) bool {
						return false
					},
					validatePasswordCreateToken: func(ctx context.Context, token string) (user *types.User, err error) {
						u := makeMockUser()
						u.SetRoles()

						return u, nil
					},
				}
			},
		},
		{
			name:     "invalid password create token",
			payload:  map[string]string(nil),
			alerts:   []request.Alert{{Type: "warning", Text: "create-password.alerts.invalid-expired-password-token"}},
			link:     GetLinks().Login,
			template: TmplCreatePassword,
			fn: func(_ *settings.Settings) {
				req.URL = &url.URL{RawQuery: "token=NOT_EMPTY"}

				authService = &authServiceMocked{
					passwordSet: func(ctx context.Context, s string) bool {
						return false
					},
					validatePasswordCreateToken: func(ctx context.Context, token string) (user *types.User, err error) {
						return nil, errors.New("invalid token")
					},
				}
			},
		},
		{
			name:     "invalid password create request",
			payload:  map[string]string(nil),
			alerts:   []request.Alert{{Type: "warning", Text: "create-password.alerts.invalid-expired-password-token"}},
			link:     GetLinks().Login,
			template: TmplCreatePassword,
			fn: func(_ *settings.Settings) {
				req.URL = &url.URL{RawQuery: "token=NOT_EMPTY"}
				user = makeMockUser()
				authService = &authServiceMocked{
					passwordSet: func(ctx context.Context, s string) bool {
						return false
					},
					validatePasswordCreateToken: func(ctx context.Context, token string) (user *types.User, err error) {
						return nil, errors.New("invalid token")
					},
				}
			},
		},
	}

	for _, tc := range tcc {
		t.Run(tc.name, func(t *testing.T) {
			rq := require.New(t)

			// reset from previous
			req.Form = url.Values{}
			req.PostForm = url.Values{}

			authSettings := &settings.Settings{}

			tc.fn(authSettings)

			authHandlers = prepareClientAuthHandlers(authService, authSettings)
			authReq = prepareClientAuthReq(authHandlers, req, user)

			// unset so we get to the main functionality
			authReq.AuthUser = nil

			err := authHandlers.createPasswordForm(authReq)

			rq.NoError(err)
			rq.Equal(tc.template, authReq.Template)
			rq.Equal(tc.payload, authReq.GetKV())
			rq.Equal(tc.alerts, authReq.NewAlerts)
			rq.Equal(tc.link, authReq.RedirectTo)
		})
	}
}

func Test_createPasswordProc(t *testing.T) {
	var (
		user = makeMockUser()

		req = &http.Request{}

		authService  authService
		authHandlers *AuthHandlers
		authReq      *request.AuthReq
	)

	tcc := []testingExpect{
		{
			name:    "create password success",
			payload: map[string]string(nil),
			alerts:  []request.Alert{{Type: "primary", Text: "create-password.alerts.password-create-success", Html: ""}},
			link:    GetLinks().Profile,
			fn: func(_ *settings.Settings) {
				authService = &authServiceMocked{
					setPassword: func(ctx context.Context, userID uint64, password string) (err error) {
						return nil
					},
				}
			},
		},
		{
			name:    "create password disabled",
			payload: map[string]string(nil),
			alerts:  []request.Alert{{Type: "danger", Text: "create-password.alert.password-create-disabled", Html: ""}},
			link:    GetLinks().Login,
			fn: func(_ *settings.Settings) {
				authService = &authServiceMocked{
					setPassword: func(ctx context.Context, userID uint64, password string) (err error) {
						return service.AuthErrPasswordCreateDisabledByConfig()
					},
				}
			},
		},
	}

	for _, tc := range tcc {
		t.Run(tc.name, func(t *testing.T) {
			rq := require.New(t)

			authSettings := &settings.Settings{}

			tc.fn(authSettings)

			authHandlers = prepareClientAuthHandlers(authService, authSettings)
			authReq = prepareClientAuthReq(authHandlers, req, user)

			err := authHandlers.createPasswordProc(authReq)

			rq.NoError(err)
			rq.Equal(tc.payload, authReq.GetKV())
			rq.Equal(tc.alerts, authReq.NewAlerts)
			rq.Equal(tc.link, authReq.RedirectTo)
		})
	}
}
