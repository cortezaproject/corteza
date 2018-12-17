package rest

import (
	"github.com/go-chi/chi"

	"github.com/crusttech/crust/crm/rest/handlers"
	"github.com/crusttech/crust/internal/auth"
)

func MountRoutes(jwtAuth auth.TokenEncoder) func(chi.Router) {
	var (
		field    = Field{}.New()
		module   = Module{}.New()
		page     = Page{}.New()
		workflow = Workflow{}.New()
		job      = Job{}.New()
	)

	// Initialize handers & controllers.
	return func(r chi.Router) {
		// Protect all _private_ routes
		r.Group(func(r chi.Router) {
			r.Use(auth.MiddlewareValidOnly)

			handlers.NewField(field).MountRoutes(r)
			handlers.NewPage(page).MountRoutes(r)
			handlers.NewModule(module).MountRoutes(r)
			handlers.NewWorkflow(workflow).MountRoutes(r)
			handlers.NewJob(job).MountRoutes(r)
		})
	}
}
