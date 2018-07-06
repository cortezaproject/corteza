package crm

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `.go`, `.util.go` or `_test.go` to
	implement your API calls, helper functions and tests. The file `.go`
	is only generated the first time, and will not be overwritten if it exists.
*/

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
