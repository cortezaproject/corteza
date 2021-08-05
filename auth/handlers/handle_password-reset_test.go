package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/cortezaproject/corteza-server/auth/request"
	"github.com/cortezaproject/corteza-server/auth/settings"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/stretchr/testify/require"
)

func Test_requestPasswordResetForm(t *testing.T) {
	var (
		ctx  = context.Background()
		user = makeMockUser(ctx)

		req = &http.Request{
			URL: &url.URL{},
		}

		authService  authService
		authHandlers *AuthHandlers
		authReq      *request.AuthReq

		authSettings = &settings.Settings{}

		rq = require.New(t)
	)

	authReq = prepareClientAuthReq(ctx, req, user)
	authHandlers = prepareClientAuthHandlers(ctx, authService, authSettings)

	payload := map[string]string{"foo": "bar"}
	authReq.SetKV(payload)

	err := authHandlers.requestPasswordResetForm(authReq)

	rq.NoError(err)
	rq.Equal(TmplRequestPasswordReset, authReq.Template)
	rq.Equal(authReq.Data["form"], payload)
}

func Test_resetPasswordForm(t *testing.T) {
	var (
		ctx = context.Background()
		req = &http.Request{
			URL: &url.URL{},
		}

		authService  authService
		authHandlers *AuthHandlers
		authReq      *request.AuthReq
	)

	tcc := []testingExpect{
		{
			name:     "request reset success",
			payload:  map[string]string(nil),
			alerts:   []request.Alert(nil),
			link:     GetLinks().ResetPassword,
			template: TmplResetPassword,
			fn: func(_ *settings.Settings) {
				req.URL = &url.URL{RawQuery: "token=NOT_EMPTY"}

				authService = &authServiceMocked{
					validatePasswordResetToken: func(ctx context.Context, token string) (user *types.User, err error) {
						u := makeMockUser(ctx)
						u.SetRoles([]uint64{})

						return u, nil
					},
				}
			},
		},
		{
			name:     "invalid password reset token",
			payload:  map[string]string(nil),
			alerts:   []request.Alert{{Type: "warning", Text: "Invalid or expired password reset token, please repeat password reset request."}},
			link:     GetLinks().RequestPasswordReset,
			template: TmplResetPassword,
			fn: func(_ *settings.Settings) {
				req.URL = &url.URL{RawQuery: "token=NOT_EMPTY"}

				authService = &authServiceMocked{
					validatePasswordResetToken: func(ctx context.Context, token string) (user *types.User, err error) {
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

			authReq = prepareClientAuthReq(ctx, req, nil)
			authHandlers = prepareClientAuthHandlers(ctx, authService, authSettings)

			// unset so we get to the main functionality
			authReq.AuthUser = nil

			err := authHandlers.resetPasswordForm(authReq)

			rq.NoError(err)
			rq.Equal(tc.template, authReq.Template)
			rq.Equal(tc.payload, authReq.GetKV())
			rq.Equal(tc.alerts, authReq.NewAlerts)
			rq.Equal(tc.link, authReq.RedirectTo)
		})
	}
}

func Test_requestPasswordReset(t *testing.T) {
	var (
		ctx  = context.Background()
		user = makeMockUser(ctx)

		req = &http.Request{}

		authService  authService
		authHandlers *AuthHandlers
		authReq      *request.AuthReq
	)

	tcc := []testingExpect{
		{
			name:    "request reset success",
			payload: map[string]string(nil),
			alerts:  []request.Alert(nil),
			link:    GetLinks().PasswordResetRequested,
			fn: func(_ *settings.Settings) {
				authService = &authServiceMocked{
					sendPasswordResetToken: func(ctx context.Context, email string) (err error) {
						return nil
					},
				}
			},
		},
		{
			name:    "request reset disabled",
			payload: map[string]string(nil),
			alerts:  []request.Alert{{Type: "danger", Text: "Password reset disabled"}},
			link:    GetLinks().Login,
			fn: func(_ *settings.Settings) {
				authService = &authServiceMocked{
					sendPasswordResetToken: func(ctx context.Context, email string) (err error) {
						return service.AuthErrPasswordResetDisabledByConfig()
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
			user.Meta = &types.UserMeta{}

			authSettings := &settings.Settings{}

			tc.fn(authSettings)

			authReq = prepareClientAuthReq(ctx, req, user)
			authHandlers = prepareClientAuthHandlers(ctx, authService, authSettings)

			err := authHandlers.requestPasswordResetProc(authReq)

			rq.NoError(err)
			rq.Equal(tc.payload, authReq.GetKV())
			rq.Equal(tc.alerts, authReq.NewAlerts)
			rq.Equal(tc.link, authReq.RedirectTo)
		})
	}
}

func Test_requestPasswordProc(t *testing.T) {
	var (
		ctx  = context.Background()
		user = makeMockUser(ctx)

		req = &http.Request{}

		authService  authService
		authHandlers *AuthHandlers
		authReq      *request.AuthReq
	)

	tcc := []testingExpect{
		{
			name:    "reset password success",
			payload: map[string]string(nil),
			alerts:  []request.Alert{{Type: "primary", Text: "Password successfully reset.", Html: ""}},
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
			name:    "reset password disabled",
			payload: map[string]string(nil),
			alerts:  []request.Alert{{Type: "danger", Text: "Password reset disabled", Html: ""}},
			link:    GetLinks().Login,
			fn: func(_ *settings.Settings) {
				authService = &authServiceMocked{
					setPassword: func(ctx context.Context, userID uint64, password string) (err error) {
						return service.AuthErrPasswordResetDisabledByConfig()
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

			authReq = prepareClientAuthReq(ctx, req, user)
			authHandlers = prepareClientAuthHandlers(ctx, authService, authSettings)

			err := authHandlers.resetPasswordProc(authReq)

			rq.NoError(err)
			rq.Equal(tc.payload, authReq.GetKV())
			rq.Equal(tc.alerts, authReq.NewAlerts)
			rq.Equal(tc.link, authReq.RedirectTo)
		})
	}
}
