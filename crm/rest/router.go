package rest

import (
	"github.com/go-chi/chi"

	"github.com/crusttech/crust/crm/rest/handlers"
	"github.com/crusttech/crust/crm/service"
	"github.com/crusttech/crust/internal/auth"
)

func MountRoutes(jwtAuth auth.TokenEncoder) func(chi.Router) {
	var (
		fieldSvc    = service.Field()
		moduleSvc   = service.Module()
		contentSvc  = service.Content()
		pageSvc     = service.Page()
		workflowSvc = service.Workflow()
		jobSvc      = service.Job()
	)

	var (
		field    = Field{}.New(fieldSvc)
		module   = Module{}.New(moduleSvc, contentSvc)
		page     = Page{}.New(pageSvc)
		workflow = Workflow{}.New(workflowSvc)
		job      = Job{}.New(jobSvc)
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
