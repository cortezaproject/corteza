package rest

import (
	"github.com/go-chi/chi"

	"github.com/crusttech/crust/crm/rest/handlers"
	"github.com/crusttech/crust/internal/auth"
)

func MountRoutes(jwtAuth auth.TokenEncoder) func(chi.Router) {
	var (
		module       = Module{}.New()
		page         = Page{}.New()
		chart        = Chart{}.New()
		trigger      = Trigger{}.New()
		notification = Notification{}.New()
	)

	// Initialize handers & controllers.
	return func(r chi.Router) {
		// Protect all _private_ routes
		r.Group(func(r chi.Router) {
			r.Use(auth.MiddlewareValidOnly)

			handlers.NewPage(page).MountRoutes(r)
			handlers.NewModule(module).MountRoutes(r)
			handlers.NewChart(chart).MountRoutes(r)
			handlers.NewTrigger(trigger).MountRoutes(r)
			handlers.NewNotification(notification).MountRoutes(r)
		})
	}
}
