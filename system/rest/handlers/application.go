package handlers

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `application.go`, `application.util.go` or `application_test.go` to
	implement your API calls, helper functions and tests. The file `application.go`
	is only generated the first time, and will not be overwritten if it exists.
*/

import (
	"context"
	"github.com/go-chi/chi"
	"net/http"

	"github.com/titpetric/factory/resputil"

	"github.com/crusttech/crust/system/rest/request"
)

// Internal API interface
type ApplicationAPI interface {
	List(context.Context, *request.ApplicationList) (interface{}, error)
	Create(context.Context, *request.ApplicationCreate) (interface{}, error)
	Update(context.Context, *request.ApplicationUpdate) (interface{}, error)
	Read(context.Context, *request.ApplicationRead) (interface{}, error)
	Delete(context.Context, *request.ApplicationDelete) (interface{}, error)
}

// HTTP API interface
type Application struct {
	List   func(http.ResponseWriter, *http.Request)
	Create func(http.ResponseWriter, *http.Request)
	Update func(http.ResponseWriter, *http.Request)
	Read   func(http.ResponseWriter, *http.Request)
	Delete func(http.ResponseWriter, *http.Request)
}

func NewApplication(ah ApplicationAPI) *Application {
	return &Application{
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewApplicationList()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return ah.List(r.Context(), params)
			})
		},
		Create: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewApplicationCreate()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return ah.Create(r.Context(), params)
			})
		},
		Update: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewApplicationUpdate()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return ah.Update(r.Context(), params)
			})
		},
		Read: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewApplicationRead()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return ah.Read(r.Context(), params)
			})
		},
		Delete: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewApplicationDelete()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return ah.Delete(r.Context(), params)
			})
		},
	}
}

func (ah *Application) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Route("/application", func(r chi.Router) {
			r.Get("/", ah.List)
			r.Post("/", ah.Create)
			r.Put("/{applicationID}", ah.Update)
			r.Get("/{applicationID}", ah.Read)
			r.Delete("/{applicationID}", ah.Delete)
		})
	})
}
