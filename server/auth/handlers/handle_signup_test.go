package handlers

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/cortezaproject/corteza-server/auth/request"
	"github.com/cortezaproject/corteza-server/auth/settings"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/stretchr/testify/require"
)

func Test_signupForm(t *testing.T) {
	var (
		user = makeMockUser()

		req = &http.Request{
			URL: &url.URL{},
		}

		authService  authService
		authHandlers *AuthHandlers
		authReq      *request.AuthReq

		authSettings = &settings.Settings{}

		rq = require.New(t)
	)

	authHandlers = prepareClientAuthHandlers(authService, authSettings)
	authReq = prepareClientAuthReq(authHandlers, req, user)

	userForm := map[string]string{
		"email":  user.Email,
		"handle": user.Handle,
		"name":   user.Name,
	}

	authReq.SetKV(userForm)

	err := authHandlers.signupForm(authReq)

	rq.NoError(err)
	rq.Equal(TmplSignup, authReq.Template)
	rq.Equal(userForm, authReq.Data["form"])
}

func Test_signupProc(t *testing.T) {
	var (
		user = makeMockUser()

		req = &http.Request{
			PostForm: url.Values{},
		}

		authService  authService
		authHandlers *AuthHandlers
		authReq      *request.AuthReq
	)

	tcc := []testingExpect{
		{
			name:    "success email confirmed",
			err:     "",
			alerts:  []request.Alert{{Type: "primary", Text: "signup.alerts.signup-successful", Html: ""}},
			link:    GetLinks().Profile,
			payload: map[string]string(nil),
			fn: func(_ *settings.Settings) {
				authService = &authServiceMocked{
					internalSignUp: func(c context.Context, user *types.User, s string) (u *types.User, err error) {
						u = &types.User{
							EmailConfirmed: true,
							Meta:           &types.UserMeta{},
						}

						u.Meta.SecurityPolicy.MFA.EnforcedEmailOTP = true
						u.Meta.SecurityPolicy.MFA.EnforcedTOTP = false

						err = nil

						return
					},
				}
			},
		},
		{
			name:    "success email unconfirmed",
			err:     "",
			alerts:  []request.Alert(nil),
			link:    GetLinks().PendingEmailConfirmation,
			payload: map[string]string(nil),
			fn: func(_ *settings.Settings) {
				authService = &authServiceMocked{
					internalSignUp: func(c context.Context, user *types.User, s string) (u *types.User, err error) {
						return &types.User{EmailConfirmed: false}, nil
					},
				}
			},
		},
		{
			name:    "internal signup disabled",
			err:     "",
			alerts:  []request.Alert{{Type: "danger", Text: "signup.alerts.signup-disabled", Html: ""}},
			link:    GetLinks().Login,
			payload: map[string]string(nil),
			fn: func(_ *settings.Settings) {
				authService = &authServiceMocked{
					internalSignUp: func(c context.Context, user *types.User, s string) (u *types.User, err error) {
						return nil, service.AuthErrInternalSignupDisabledByConfig()
					},
				}
			},
		},
		{
			name:    "invalid email format",
			err:     "",
			alerts:  []request.Alert(nil),
			link:    GetLinks().Signup,
			payload: map[string]string{"email": "", "error": "invalid email", "handle": "", "name": ""},
			fn: func(_ *settings.Settings) {
				authService = &authServiceMocked{
					internalSignUp: func(c context.Context, user *types.User, s string) (u *types.User, err error) {
						return nil, service.AuthErrInvalidEmailFormat()
					},
				}
			},
		},
		{
			name:    "invalid handle",
			err:     "",
			alerts:  []request.Alert(nil),
			link:    GetLinks().Signup,
			payload: map[string]string{"email": "", "error": "invalid handle", "handle": "", "name": ""},
			fn: func(_ *settings.Settings) {
				authService = &authServiceMocked{
					internalSignUp: func(c context.Context, user *types.User, s string) (u *types.User, err error) {
						return nil, service.AuthErrInvalidHandle()
					},
				}
			},
		},
		{
			name:    "password not secure",
			err:     "",
			alerts:  []request.Alert(nil),
			link:    GetLinks().Signup,
			payload: map[string]string{"email": "", "error": "provided password is not secure; use longer password to follow the password policy", "handle": "", "name": ""},
			fn: func(_ *settings.Settings) {
				authService = &authServiceMocked{
					internalSignUp: func(c context.Context, user *types.User, s string) (u *types.User, err error) {
						return nil, service.AuthErrPasswordNotSecure()
					},
				}
			},
		},
		{
			name:    "invalid credentials",
			err:     "",
			alerts:  []request.Alert(nil),
			link:    GetLinks().Signup,
			payload: map[string]string{"email": "", "error": "invalid username and password combination", "handle": "", "name": ""},
			fn: func(_ *settings.Settings) {
				authService = &authServiceMocked{
					internalSignUp: func(c context.Context, user *types.User, s string) (u *types.User, err error) {
						return nil, service.AuthErrInvalidCredentials()
					},
				}
			},
		},
		{
			name:    "handle not unique",
			err:     "handle not unique",
			alerts:  []request.Alert(nil),
			link:    GetLinks().Signup,
			payload: map[string]string{"email": "test@mail.com", "error": "handle not unique", "handle": "test-user", "name": "Test User"},
			postForm: url.Values{
				"email":  {"test@mail.com"},
				"handle": {"test-user"},
				"name":   {"Test User"},
			},
			fn: func(_ *settings.Settings) {
				authService = &authServiceMocked{
					internalSignUp: func(c context.Context, user *types.User, s string) (u *types.User, err error) {
						return nil, service.UserErrHandleNotUnique()
					},
				}
			},
		},
		{
			name:    "email not unique",
			err:     "email not unique",
			alerts:  []request.Alert(nil),
			link:    GetLinks().Signup,
			payload: map[string]string{"email": "test@mail.com", "error": "email not unique", "handle": "test-user", "name": "Test User"},
			postForm: url.Values{
				"email":  {"test@mail.com"},
				"handle": {"test-user"},
				"name":   {"Test User"},
			},
			fn: func(_ *settings.Settings) {
				authService = &authServiceMocked{
					internalSignUp: func(c context.Context, user *types.User, s string) (u *types.User, err error) {
						return nil, service.UserErrEmailNotUnique()
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
			req.PostForm = tc.postForm

			authSettings := &settings.Settings{}

			tc.fn(authSettings)

			authHandlers = prepareClientAuthHandlers(authService, authSettings)
			authReq = prepareClientAuthReq(authHandlers, req, user)

			err := authHandlers.signupProc(authReq)

			if len(tc.err) > 0 {
				rq.Equal(tc.err, err.Error())
			} else {
				rq.NoError(err)
			}
			rq.Equal(tc.template, authReq.Template)
			rq.Equal(tc.payload, authReq.GetKV())
			rq.Equal(tc.alerts, authReq.NewAlerts)
			rq.Equal(tc.link, authReq.RedirectTo)
		})
	}
}
