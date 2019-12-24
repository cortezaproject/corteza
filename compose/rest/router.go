package rest

import (
	"github.com/go-chi/chi"

	"github.com/cortezaproject/corteza-server/compose/rest/handlers"
	"github.com/cortezaproject/corteza-server/pkg/auth"
)

func MountRoutes(r chi.Router) {
	var (
		namespace    = Namespace{}.New()
		module       = Module{}.New()
		record       = Record{}.New()
		page         = Page{}.New()
		chart        = Chart{}.New()
		notification = Notification{}.New()
		attachment   = Attachment{}.New()
	)

	// Initialize handlers & controllers.
	r.Group(func(r chi.Router) {
		// Use alternative handlers that support file serving
		handlers.NewAttachment(attachment).MountRoutes(r)
	})

	// Protect all _private_ routes
	r.Group(func(r chi.Router) {
		r.Use(auth.MiddlewareValidOnly)
		r.Use(middlewareAllowedAccess)

		handlers.NewNamespace(namespace).MountRoutes(r)
		handlers.NewPage(page).MountRoutes(r)
		handlers.NewModule(module).MountRoutes(r)
		handlers.NewRecord(record).MountRoutes(r)
		handlers.NewChart(chart).MountRoutes(r)
		handlers.NewNotification(notification).MountRoutes(r)
		handlers.NewPermissions(Permissions{}.New()).MountRoutes(r)
		handlers.NewSettings(Settings{}.New()).MountRoutes(r)
	})
}
