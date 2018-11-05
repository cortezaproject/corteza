package rest

import (
	"context"

	"github.com/pkg/errors"

	"github.com/crusttech/crust/system/rest/handlers"
	"github.com/crusttech/crust/system/rest/request"
)

var _ = errors.Wrap

type Auth struct{}

func (Auth) New() *Auth {
	return &Auth{}
}

func (ctrl *Auth) Check(ctx context.Context, r *request.AuthCheck) (interface{}, error) {
	return nil, errors.New("Not implemented: Auth.check")
}

func (ctrl *Auth) Logout(ctx context.Context, r *request.AuthLogout) (interface{}, error) {
	return nil, errors.New("Not implemented: Auth.logout")
}

func (ctrl *Auth) Handlers(jwtAuth jwtEncodeCookieSetter) *handlers.Auth {
	h := handlers.NewAuth(ctrl)
	// Check JWT if valid
	h.Check = func(w http.ResponseWriter, r *http.Request) {
		if c, err := r.Cookie("jwt"); err == nil {
			ctx := r.Context()

			if identity := auth.GetIdentityFromContext(ctx); identity != nil && identity.Valid() {
				if user, err := service.DefaultUser.With(ctx).FindByID(identity.Identity()); err == nil {
					resputil.JSON(w, checkResponse{
						JWT:  c.Value,
						User: payload.User(user),
					})

					return
				}
			}

			// Did not send response, assuming invalid cookie
			jwtAuth.SetCookie(w, r, nil)
		} else {
			resputil.JSON(w, err)
		}
	}
	h.Logout = func(w http.ResponseWriter, r *http.Request) {
		jwtAuth.SetCookie(w, r, nil)
	}
	return h
}
