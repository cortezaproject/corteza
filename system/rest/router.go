package rest

import (
	"context"
	"log"

	"github.com/go-chi/chi"

	"github.com/crusttech/crust/internal/auth"
	"github.com/crusttech/crust/internal/config"
	"github.com/crusttech/crust/system/repository"
	"github.com/crusttech/crust/system/rest/handlers"
	"github.com/crusttech/crust/system/service"
)

func MountRoutes(oidcConfig *config.OIDC, jwtAuth jwtEncodeCookieSetter) func(chi.Router) {
	var userSvc = service.User()
	var ctx = context.Background()

	oidc, err := OpenIdConnect(ctx, oidcConfig, userSvc, jwtAuth, repository.NewSettings(ctx, repository.DB(ctx)))
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

		// Provide raw `/auth` handlers
		Auth{}.New().Handlers(jwtAuth).MountRoutes(r)

		// Protect all _private_ routes
		r.Group(func(r chi.Router) {
			r.Use(auth.MiddlewareValidOnly)

			handlers.NewUser(User{}.New()).MountRoutes(r)
			handlers.NewTeam(Team{}.New()).MountRoutes(r)
			handlers.NewOrganisation(Organisation{}.New()).MountRoutes(r)
		})
	}
}
