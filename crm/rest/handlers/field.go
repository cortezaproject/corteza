package handlers

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `field.go`, `field.util.go` or `field_test.go` to
	implement your API calls, helper functions and tests. The file `field.go`
	is only generated the first time, and will not be overwritten if it exists.
*/

import (
	"context"
	"github.com/go-chi/chi"
	"net/http"

	"github.com/titpetric/factory/resputil"

	"github.com/crusttech/crust/crm/rest/request"
)

// Internal API interface
type FieldAPI interface {
	List(context.Context, *request.FieldList) (interface{}, error)
	Type(context.Context, *request.FieldType) (interface{}, error)
}

// HTTP API interface
type Field struct {
	List func(http.ResponseWriter, *http.Request)
	Type func(http.ResponseWriter, *http.Request)
}

func NewField(fh FieldAPI) *Field {
	return &Field{
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewFieldList()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return fh.List(r.Context(), params)
			})
		},
		Type: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewFieldType()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return fh.Type(r.Context(), params)
			})
		},
	}
}

func (fh *Field) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Route("/field", func(r chi.Router) {
			r.Get("/", fh.List)
			r.Get("/{id}", fh.Type)
		})
	})
}
