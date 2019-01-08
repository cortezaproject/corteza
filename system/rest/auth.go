package rest

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
	"github.com/titpetric/factory/resputil"

	"github.com/crusttech/crust/internal/auth"
	"github.com/crusttech/crust/internal/payload"
	"github.com/crusttech/crust/internal/payload/outgoing"
	"github.com/crusttech/crust/system/rest/handlers"
	"github.com/crusttech/crust/system/rest/request"
	"github.com/crusttech/crust/system/service"
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

func (ctrl *Auth) Login(ctx context.Context, r *request.AuthLogin) (interface{}, error) {
	return nil, errors.New("Not implemented: Auth.login")
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
	h.Login = func(w http.ResponseWriter, r *http.Request) {
		params := request.NewAuthLogin()
		ctx := r.Context()

		// parse request to fill parameters
		err := params.Fill(r)
		if err != nil {
			resputil.JSON(w, err)
			return
		}

		userSvc := service.User().With(ctx)

		// check email and username for login
		user, err := userSvc.FindByEmail(params.Username)
		if err != nil {
			user, err = userSvc.FindByUsername(params.Username)
		}

		// can't find user
		if err != nil {
			resputil.JSON(w, err)
			return
		}

		// validate password
		if user.ValidatePassword(params.Password) {
			jwtEncoder.SetCookie(w, r, user)
			resputil.JSON(w, user, err)
			return
		}

		resputil.JSON(w, errors.New("Password doesn't match"))
	}
	return h
}
