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

func MountRoutes(oidcConfig *config.OIDC, socialConfig *config.Social, jwtEncoder auth.TokenEncoder) func(chi.Router) {
	var err error
	var userSvc = service.User()
	var ctx = context.Background()
	var oidc *openIdConnect

	if oidcConfig.Enabled {
		oidc, err = OpenIdConnect(ctx, oidcConfig, userSvc, jwtEncoder, repository.NewSettings(ctx, repository.DB(ctx)))
		if err != nil {
			log.Println("Could not initialize OIDC:", err.Error())
		}
	} else {
		log.Println("OIDC is disabled")
	}

	// Initialize handers & controllers.
	return func(r chi.Router) {
		if oidcConfig.Enabled && oidc != nil {
			r.Route("/oidc", func(r chi.Router) {
				r.Get("/", oidc.HandleRedirect)
				r.Get("/callback", oidc.HandleOAuth2Callback)
			})
		}

		NewSocial(socialConfig, jwtEncoder).MountRoutes(r)

		// Provide raw `/auth` handlers
		Auth{}.New().Handlers(jwtEncoder).MountRoutes(r)

		// Protect all _private_ routes
		r.Group(func(r chi.Router) {
			r.Use(auth.MiddlewareValidOnly)

			handlers.NewUser(User{}.New()).MountRoutes(r)
			handlers.NewRole(Role{}.New()).MountRoutes(r)
			handlers.NewOrganisation(Organisation{}.New()).MountRoutes(r)
			handlers.NewPermissions(Permissions{}.New()).MountRoutes(r)
			handlers.NewApplication(Application{}.New()).MountRoutes(r)
		})
	}
}
