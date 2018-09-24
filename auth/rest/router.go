package rest

import (
	"context"
	"log"
	"net/http"

	"github.com/crusttech/crust/auth/repository"
	"github.com/crusttech/crust/auth/types"
	"github.com/crusttech/crust/internal/auth"
	"github.com/go-chi/chi"
	"github.com/titpetric/factory/resputil"

	"github.com/crusttech/crust/auth/service"
	"github.com/crusttech/crust/internal/config"
)

type (
	checkResponse struct {
		JWT  string      `json:"jwt"`
		User *types.User `json:"user"`
	}
)

func MountRoutes(oidcConfig *config.OIDC, jwtAuth jwtEncodeCookieSetter) func(chi.Router) {
	var userSvc = service.User()
	var ctx = context.Background()

	oidc, err := OpenIdConnect(ctx, oidcConfig, userSvc, jwtAuth, repository.NewSettings(ctx))
	if err != nil {
		log.Print("Could not initialize OIDC:", err.Error())
	}

	// Initialize handers & controllers.
	return func(r chi.Router) {
		if oidc != nil {
			r.Route("/oidc", func(r chi.Router) {
				r.Get("/", oidc.HandleRedirect)
				r.Get("/callback", oidc.HandleOAuth2Callback)
			})
		}

		r.Get("/check", func(w http.ResponseWriter, r *http.Request) {
			if c, err := r.Cookie("jwt"); err == nil {
				ctx := r.Context()

				if identity := auth.GetIdentityFromContext(ctx); identity != nil && identity.Valid() {
					if user, err := service.DefaultUser.With(ctx).FindByID(identity.Identity()); err == nil {
						resputil.JSON(w, checkResponse{
							JWT:  c.Value,
							User: user,
						})

						return
					}
				}

				// Did not send response, assuming invalid cookie
				jwtAuth.SetCookie(w, r, nil)
			}

			resputil.JSON(w, "")
		})

		r.Delete("/check", func(w http.ResponseWriter, r *http.Request) {
			jwtAuth.SetCookie(w, r, nil)
		})
	}
}
