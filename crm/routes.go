package crm

import (
	"github.com/go-chi/chi"
)

func MountRoutes(r chi.Router) {
	modules := ModulesHandlers{}.new()
	types := TypesHandlers{}.new()
	r.Route("/modules", func(r chi.Router) {
		r.Get("/list", modules.List)
		r.Post("/edit", modules.Edit)
		r.Get("/content/list", modules.ContentList)
		r.Post("/content/edit", modules.ContentEdit)
		r.Delete("/content/delete", modules.ContentDelete)
	})
	r.Route("/types", func(r chi.Router) {
		r.Get("/list", types.List)
		r.Get("/type/{id}", types.Type)
	})
}
