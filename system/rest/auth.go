package rest

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
	"github.com/titpetric/factory/resputil"

	"github.com/crusttech/crust/internal/auth"
	"github.com/crusttech/crust/internal/payload"
	"github.com/crusttech/crust/internal/payload/outgoing"
	"github.com/crusttech/crust/system/internal/service"
	"github.com/crusttech/crust/system/rest/handlers"
	"github.com/crusttech/crust/system/rest/request"
)

var _ = errors.Wrap

type (
	Auth struct {
		jwt auth.TokenEncoder
	}

	checkResponse struct {
		JWT  string         `json:"jwt"`
		User *outgoing.User `json:"user"`
	}
)

func (Auth) New() *Auth {
	return &Auth{}
}

func (ctrl *Auth) Check(ctx context.Context, r *request.AuthCheck) (interface{}, error) {
	return nil, errors.New("Not implemented: Auth.check")
}

func (ctrl *Auth) Logout(ctx context.Context, r *request.AuthLogout) (interface{}, error) {
	return nil, errors.New("Not implemented: Auth.logout")
}

// Handlers() func ignores "std" crust controllers
//
// Crush handlers are too abstract for our auth needs so we need (direct access to htt.ResponseWriter)
func (ctrl *Auth) Handlers(jwtEncoder auth.TokenEncoder) *handlers.Auth {
	h := handlers.NewAuth(ctrl)
	// Check JWT if valid
	h.Check = func(w http.ResponseWriter, r *http.Request) {
		if c, err := r.Cookie("jwt"); err == nil {
			ctx := r.Context()

			if identity := auth.GetIdentityFromContext(ctx); identity != nil && identity.Valid() {
				if user, err := service.DefaultUser.With(ctx).FindByID(identity.Identity()); err == nil {
					jwtEncoder.SetCookie(w, r, user)

					resputil.JSON(w, checkResponse{
						JWT:  c.Value,
						User: payload.User(user),
					})

					return
				}
			}

			// Did not send response, assuming invalid cookie
			jwtEncoder.SetCookie(w, r, nil)
		} else {
			resputil.JSON(w, err)
		}
	}

	h.Logout = func(w http.ResponseWriter, r *http.Request) {
		jwtEncoder.SetCookie(w, r, nil)
	}

	return h
}
