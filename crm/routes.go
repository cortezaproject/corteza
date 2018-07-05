package crm

import (
	"github.com/go-chi/chi"
)

func MountRoutes(r chi.Router) {
	module := ModuleHandlers{}.new()
	types := TypesHandlers{}.new()
	r.Route("/module", func(r chi.Router) {
		r.Get("/list", module.List)
		r.Post("/edit", module.Edit)
		r.Get("/content/list", module.ContentList)
		r.Post("/content/edit", module.ContentEdit)
		r.Delete("/content/delete", module.ContentDelete)
	})
	r.Route("/types", func(r chi.Router) {
		r.Get("/list", types.List)
		r.Get("/type/{id}", types.Type)
	})
}