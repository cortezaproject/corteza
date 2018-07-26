package server

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `team.go`, `team.util.go` or `team_test.go` to
	implement your API calls, helper functions and tests. The file `team.go`
	is only generated the first time, and will not be overwritten if it exists.
*/

import (
	"github.com/go-chi/chi"
	"net/http"
)

func (th *TeamHandlers) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Route("/teams", func(r chi.Router) {
			r.Get("/", th.List)
			r.Put("/", th.Create)
			r.Post("/{teamID}", th.Edit)
			r.Get("/{teamID}", th.Read)
			r.Delete("/{teamID}", th.Remove)
			r.Post("/{teamID}/archive", th.Archive)
			r.Post("/{teamID}/move", th.Move)
			r.Post("/{teamID}/merge", th.Merge)
		})
	})
}
