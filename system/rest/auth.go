package rest

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/api"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/payload"
	"github.com/cortezaproject/corteza-server/pkg/payload/outgoing"
	"github.com/cortezaproject/corteza-server/system/rest/request"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/pkg/errors"
	"net/http"
)

var _ = errors.Wrap

type (
	Auth struct {
		tokenEncoder auth.TokenEncoder
		settings     *types.AppSettings
		authSvc      authUserService
	}

	authUserResponse struct {
		JWT  string           `json:"jwt"`
		User *authUserPayload `json:"user"`
	}

	authUserPayload struct {
		*outgoing.User
		Roles []string `json:"roles"`
	}

	authUserService interface {
		Impersonate(ctx context.Context, userID uint64) (*types.User, error)
		ValidateAuthRequestToken(ctx context.Context, token string) (user *types.User, err error)
		CanRegister(ctx context.Context) error
		LoadRoleMemberships(ctx context.Context, user *types.User) error
	}
)

func (Auth) New() *Auth {
	return &Auth{
		tokenEncoder: auth.DefaultJwtHandler,
		settings:     service.CurrentSettings,
		authSvc:      service.DefaultAuth,
	}
}

func (ctrl *Auth) Check(ctx context.Context, r *request.AuthCheck) (interface{}, error) {
	return func(w http.ResponseWriter, r *http.Request) {
		if identity := auth.GetIdentityFromContext(ctx); identity != nil && identity.Valid() {
			if user, err := service.DefaultUser.With(ctx).FindByID(identity.Identity()); err == nil {
				var p *authUserResponse

				if p, err = ctrl.makePayload(ctx, user); err != nil {
					api.Send(w, r, err)
				} else {
					api.Send(w, r, p)

				}

				return
			}
		}

		api.Send(w, r, errors.New("not authenticated"))
	}, nil
}

func (ctrl *Auth) Logout(ctx context.Context, r *request.AuthLogout) (interface{}, error) {
	return true, nil
}

// Impersonate implements impersonation functionality
//
// This is experimental and internals will most likely change in the future:
func (ctrl *Auth) Impersonate(ctx context.Context, r *request.AuthImpersonate) (interface{}, error) {
	u, err := ctrl.authSvc.Impersonate(ctx, r.UserID)
	if err != nil {
		return nil, err
	}

	return ctrl.makePayload(ctx, u)
}

func (ctrl *Auth) Settings(ctx context.Context, r *request.AuthSettings) (interface{}, error) {
	var (
		int = ctrl.settings.Auth.Internal
		ext = ctrl.settings.Auth.External

		out = map[string]interface{}{
			"internalEnabled":                         int.Enabled,
			"internalPasswordResetEnabled":            int.PasswordReset.Enabled,
			"internalSignUpEmailConfirmationRequired": int.Signup.EmailConfirmationRequired,
			"internalSignUpEnabled":                   int.Signup.Enabled,

			"externalEnabled":   ext.Enabled,
			"externalProviders": ext.Providers.Valid(),
		}
	)

	if err := ctrl.authSvc.CanRegister(ctx); err != nil {
		// f["internalSignUpEnabled"] = false
		out["signUpDisabled"] = err.Error()
	}

	return out, nil
}

func (ctrl *Auth) ExchangeAuthToken(ctx context.Context, r *request.AuthExchangeAuthToken) (interface{}, error) {
	user, err := ctrl.authSvc.ValidateAuthRequestToken(ctx, r.Token)
	if err != nil {
		return nil, err
	}

	return ctrl.makePayload(ctx, user)
}

func (ctrl *Auth) makePayload(ctx context.Context, user *types.User) (*authUserResponse, error) {
	if err := ctrl.authSvc.LoadRoleMemberships(ctx, user); err != nil {
		return nil, err
	}

	return &authUserResponse{
		JWT: ctrl.tokenEncoder.Encode(user),
		User: &authUserPayload{
			User:  payload.User(user),
			Roles: payload.Uint64stoa(user.Roles()),
		},
	}, nil
}
