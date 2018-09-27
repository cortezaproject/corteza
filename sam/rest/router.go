package rest

import (
	"github.com/crusttech/crust/internal/auth"
	"github.com/crusttech/crust/sam/rest/handlers"
	"github.com/go-chi/chi"
)

func MountRoutes() func(chi.Router) {
	// Initialize handers & controllers.
	return func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(auth.MiddlewareValidOnly404)
			handlers.NewAttachmentDownloadable(Attachment{}.New()).MountRoutes(r)
		})

		// Protect all _private_ routes
		r.Group(func(r chi.Router) {
			r.Use(auth.MiddlewareValidOnly)

			handlers.NewChannel(Channel{}.New()).MountRoutes(r)
			handlers.NewMessage(Message{}.New()).MountRoutes(r)
			handlers.NewOrganisation(Organisation{}.New()).MountRoutes(r)
			handlers.NewTeam(Team{}.New()).MountRoutes(r)
			handlers.NewUser(User{}.New()).MountRoutes(r)
		})
	}
}
