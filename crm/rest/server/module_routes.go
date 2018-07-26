package server

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `module.go`, `module.util.go` or `module_test.go` to
	implement your API calls, helper functions and tests. The file `module.go`
	is only generated the first time, and will not be overwritten if it exists.
*/

import (
	"github.com/go-chi/chi"
	"net/http"
)

func (mh *ModuleHandlers) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Route("/module", func(r chi.Router) {
			r.Get("/", mh.List)
			r.Post("/", mh.Create)
			r.Get("/{id}", mh.Read)
			r.Post("/{id}", mh.Edit)
			r.Delete("/{id}", mh.Delete)
			r.Get("/{module}/content", mh.ContentList)
			r.Post("/{module}/content", mh.ContentCreate)
			r.Get("/{module}/content/{id}", mh.ContentRead)
			r.Post("/{module}/content/{id}", mh.ContentEdit)
			r.Delete("/{module}/content/{id}", mh.ContentDelete)
		})
	})
}
