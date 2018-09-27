package rest

import (
	"github.com/crusttech/crust/crm/rest/handlers"
	"github.com/crusttech/crust/crm/service"
	"github.com/crusttech/crust/internal/auth"
	"github.com/go-chi/chi"
)

func MountRoutes(jwtAuth auth.TokenEncoder) func(chi.Router) {
	var (
		fieldSvc   = service.Field()
		moduleSvc  = service.Module()
		contentSvc = service.Content()
	)

	var (
		field  = Field{}.New(fieldSvc)
		module = Module{}.New(moduleSvc, contentSvc)
	)

	// Initialize handers & controllers.
	return func(r chi.Router) {
		// Protect all _private_ routes
		r.Group(func(r chi.Router) {
			r.Use(auth.MiddlewareValidOnly)

			handlers.NewField(field).MountRoutes(r)
			handlers.NewModule(module).MountRoutes(r)
		})
	}
}
