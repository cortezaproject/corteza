package rest

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/titpetric/factory/resputil"

	"github.com/crusttech/crust/auth/rest/handlers"
	"github.com/crusttech/crust/auth/service"
	"github.com/crusttech/crust/config"
)

func MountRoutes(oidcConfig *config.OIDC, jwtAuth jwtEncodeCookieSetter) func(chi.Router) {
	var userSvc = service.User()

	oidc, err := OpenIdConnect(oidcConfig, userSvc, jwtAuth)
	if err != nil {
		log.Errorf("Could not initialize OIDC:", err.Error())
	}

	// Initialize handers & controllers.
	return func(r chi.Router) {
		handlers.NewAuth(Auth{}.New(userSvc, jwtAuth)).MountRoutes(r)

		if oidc != nil {
			r.Route("/oidc", func(r chi.Router) {
				r.Get("/", oidc.HandleRedirect)
				r.Get("/callback", oidc.HandleOAuth2Callback)
			})
		}

		r.Get("/jwt", func(w http.ResponseWriter, r *http.Request) {
			if c, err := r.Cookie("jwt"); err != nil {
				resputil.JSON(w, "")
			} else {
				resputil.JSON(w, c.Value)
			}
		})
	}
}
