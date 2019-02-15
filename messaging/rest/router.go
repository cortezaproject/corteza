package rest

import (
	"github.com/go-chi/chi"

	"github.com/crusttech/crust/internal/auth"
	"github.com/crusttech/crust/messaging/rest/handlers"
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
			handlers.NewSearch(Search{}.New()).MountRoutes(r)
			handlers.NewPermissions(Permissions{}.New()).MountRoutes(r)
		})
	}
}
