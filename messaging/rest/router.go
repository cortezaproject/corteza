package rest

import (
	"github.com/go-chi/chi"

	"github.com/cortezaproject/corteza-server/messaging/rest/handlers"
	"github.com/cortezaproject/corteza-server/pkg/auth"
)

func MountRoutes(r chi.Router) {
	// Initialize handlers & controllers.
	r.Group(func(r chi.Router) {
		handlers.NewAttachment(Attachment{}.New()).MountRoutes(r)
	})

	// Protect all _private_ routes
	r.Group(func(r chi.Router) {
		r.Use(auth.MiddlewareValidOnly)
		handlers.NewActivity(Activity{}.New()).MountRoutes(r)
		handlers.NewChannel(Channel{}.New()).MountRoutes(r)
		handlers.NewMessage(Message{}.New()).MountRoutes(r)
		handlers.NewSearch(Search{}.New()).MountRoutes(r)
		handlers.NewStatus(Status{}.New()).MountRoutes(r)
		handlers.NewCommands(Commands{}.New()).MountRoutes(r)
		handlers.NewPermissions(Permissions{}.New()).MountRoutes(r)
	})
}
