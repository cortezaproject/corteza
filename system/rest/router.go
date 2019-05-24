package rest

import (
	"github.com/go-chi/chi"

	"github.com/cortezaproject/corteza-server/internal/auth"
	"github.com/cortezaproject/corteza-server/system/rest/handlers"
)

func MountRoutes(r chi.Router) {
	NewExternalAuth().ApiServerRoutes(r)

	// Provide raw `/auth` handlers
	handlers.NewAuth((Auth{}).New()).MountRoutes(r)

	handlers.NewAuthInternal((AuthInternal{}).New()).MountRoutes(r)

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
