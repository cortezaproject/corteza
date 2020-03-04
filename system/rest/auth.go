package rest

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
	"github.com/titpetric/factory/resputil"

	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/payload"
	"github.com/cortezaproject/corteza-server/pkg/payload/outgoing"
	"github.com/cortezaproject/corteza-server/system/rest/request"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
)

var _ = errors.Wrap

type (
	Auth struct {
		tokenEncoder auth.TokenEncoder
		settings     *types.Settings
		authSvc      service.AuthService
	}

	authUserResponse struct {
		JWT  string           `json:"jwt"`
		User *authUserPayload `json:"user"`
	}

	authUserPayload struct {
		*outgoing.User
		Roles []string `json:"roles"`
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
					resputil.JSON(w, err)
				} else {
					resputil.JSON(w, p)

				}

				return
			}
		}

		resputil.JSON(w, errors.New("not authenticated"))
	}, nil
}

func (ctrl *Auth) Logout(ctx context.Context, r *request.AuthLogout) (interface{}, error) {
	return true, nil
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
			"externalProviders": ext.Providers.Valid(ctrl.settings),
		}
	)

	if err := ctrl.authSvc.With(ctx).CanRegister(); err != nil {
		// f["internalSignUpEnabled"] = false
		out["signUpDisabled"] = err.Error()
	}

	return out, nil
}

func (ctrl *Auth) ExchangeAuthToken(ctx context.Context, r *request.AuthExchangeAuthToken) (interface{}, error) {
	var svc = ctrl.authSvc.With(ctx)

	user, err := svc.ValidateAuthRequestToken(r.Token)
	if err != nil {
		return nil, err
	}

	return ctrl.makePayload(ctx, user)
}

func (ctrl *Auth) makePayload(ctx context.Context, user *types.User) (*authUserResponse, error) {
	var svc = ctrl.authSvc.With(ctx)
	if err := svc.LoadRoleMemberships(user); err != nil {
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
