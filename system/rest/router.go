package rest

import (
	"github.com/go-chi/chi"

	"github.com/crusttech/crust/internal/auth"
	"github.com/crusttech/crust/system/rest/handlers"
)

func MountRoutes(jwtEncoder auth.TokenEncoder) func(chi.Router) {
	// Initialize handers & controllers.
	return func(r chi.Router) {
		NewSocial(jwtEncoder).MountRoutes(r)

		// Provide raw `/auth` handlers
		handlers.NewAuth((Auth{}).New(jwtEncoder)).MountRoutes(r)

		handlers.NewAuthInternal((AuthInternal{}).New(jwtEncoder)).MountRoutes(r)

		// Protect all _private_ routes
		r.Group(func(r chi.Router) {
			r.Use(auth.MiddlewareValidOnly)

			handlers.NewUser(User{}.New()).MountRoutes(r)
			handlers.NewRole(Role{}.New()).MountRoutes(r)
			handlers.NewOrganisation(Organisation{}.New()).MountRoutes(r)
			handlers.NewPermissions(Permissions{}.New()).MountRoutes(r)
			handlers.NewApplication(Application{}.New()).MountRoutes(r)
			handlers.NewSettings(Settings{}.New()).MountRoutes(r)
		})
	}
}
