package crm

import (
	"github.com/go-chi/chi"
)

func MountRoutes(r chi.Router) {
	types := TypesHandlers{}.new()
	r.Route("/types", func(r chi.Router) {
		r.Get("/list", types.List)
		r.Get("/type/{id}", types.Type)
	})
}
