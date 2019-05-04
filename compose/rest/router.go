package rest

import (
	"github.com/go-chi/chi"

	"github.com/crusttech/crust/compose/rest/handlers"
	"github.com/crusttech/crust/internal/auth"
)

func MountRoutes() func(chi.Router) {
	var (
		namespace    = Namespace{}.New()
		module       = Module{}.New()
		record       = Record{}.New()
		page         = Page{}.New()
		chart        = Chart{}.New()
		trigger      = Trigger{}.New()
		notification = Notification{}.New()
		attachment   = Attachment{}.New()
	)

	// Initialize handlers & controllers.
	return func(r chi.Router) {

		r.Group(func(r chi.Router) {
			handlers.NewPermissions(Permissions{}.New()).MountRoutes(r)
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
			handlers.NewTrigger(trigger).MountRoutes(r)
			handlers.NewNotification(notification).MountRoutes(r)
		})

		// Use alternative handlers that support file serving
		handlers.NewAttachment(attachment).MountRoutes(r)
	}
}
