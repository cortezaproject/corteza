package rest

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
	"github.com/titpetric/factory/resputil"

	"github.com/cortezaproject/corteza-server/internal/auth"
	"github.com/cortezaproject/corteza-server/internal/payload"
	"github.com/cortezaproject/corteza-server/internal/payload/outgoing"
	"github.com/cortezaproject/corteza-server/system/internal/service"
	"github.com/cortezaproject/corteza-server/system/rest/request"
)

var _ = errors.Wrap

type (
	Auth struct {
		tokenEncoder auth.TokenEncoder
		authSettings authServiceSettingsProvider
		authSvc      service.AuthService
	}

	authServiceSettingsProvider interface {
		Format() map[string]interface{}
	}

	exchangeResponse struct {
		JWT  string         `json:"jwt"`
		User *outgoing.User `json:"user"`
	}

	checkResponse struct {
		JWT  string         `json:"jwt"`
		User *outgoing.User `json:"user"`
	}
)

func (Auth) New() *Auth {
	return &Auth{
		tokenEncoder: auth.DefaultJwtHandler,
		authSettings: service.DefaultAuthSettings,
		authSvc:      service.DefaultAuth,
	}
}

func (ctrl *Auth) Check(ctx context.Context, r *request.AuthCheck) (interface{}, error) {
	return func(w http.ResponseWriter, r *http.Request) {
		if identity := auth.GetIdentityFromContext(ctx); identity != nil && identity.Valid() {
			if user, err := service.DefaultUser.With(ctx).FindByID(identity.Identity()); err == nil {
				svc := ctrl.authSvc.With(ctx)

				if err = svc.LoadRoleMemberships(user); err != nil {
					resputil.JSON(w, err)
					return
				} else {
					resputil.JSON(w, checkResponse{
						JWT:  ctrl.tokenEncoder.Encode(user),
						User: payload.User(user),
					})
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
	return ctrl.authSettings.Format(), nil
}

func (ctrl *Auth) ExchangeAuthToken(ctx context.Context, r *request.AuthExchangeAuthToken) (interface{}, error) {
	var svc = ctrl.authSvc.With(ctx)

	user, err := svc.ValidateAuthRequestToken(r.Token)
	if err != nil {
		return nil, err
	}

	if err = svc.LoadRoleMemberships(user); err != nil {
		return nil, err
	}

	return &exchangeResponse{
		JWT:  ctrl.tokenEncoder.Encode(user),
		User: payload.User(user),
	}, nil
}
