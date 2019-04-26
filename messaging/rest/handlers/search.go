package handlers

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `search.go`, `search.util.go` or `search_test.go` to
	implement your API calls, helper functions and tests. The file `search.go`
	is only generated the first time, and will not be overwritten if it exists.
*/

import (
	"context"

	"net/http"

	"github.com/go-chi/chi"
	"github.com/titpetric/factory/resputil"

	"github.com/crusttech/crust/messaging/rest/request"
)

// Internal API interface
type SearchAPI interface {
	Messages(context.Context, *request.SearchMessages) (interface{}, error)
}

// HTTP API interface
type Search struct {
	Messages func(http.ResponseWriter, *http.Request)
}

func NewSearch(sh SearchAPI) *Search {
	return &Search{
		Messages: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewSearchMessages()
			if err := params.Fill(r); err != nil {
				resputil.JSON(w, err)
				return
			}
			if value, err := sh.Messages(r.Context(), params); err != nil {
				resputil.JSON(w, err)
				return
			} else {
				switch fn := value.(type) {
				case func(http.ResponseWriter, *http.Request):
					fn(w, r)
					return
				}
				resputil.JSON(w, value)
				return
			}
		},
	}
}

func (sh *Search) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/search/messages", sh.Messages)
	})
}
