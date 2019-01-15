package handlers

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `trigger.go`, `trigger.util.go` or `trigger_test.go` to
	implement your API calls, helper functions and tests. The file `trigger.go`
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
type TriggerAPI interface {
	List(context.Context, *request.TriggerList) (interface{}, error)
	Create(context.Context, *request.TriggerCreate) (interface{}, error)
	Read(context.Context, *request.TriggerRead) (interface{}, error)
	Update(context.Context, *request.TriggerUpdate) (interface{}, error)
	Delete(context.Context, *request.TriggerDelete) (interface{}, error)
}

// HTTP API interface
type Trigger struct {
	List   func(http.ResponseWriter, *http.Request)
	Create func(http.ResponseWriter, *http.Request)
	Read   func(http.ResponseWriter, *http.Request)
	Update func(http.ResponseWriter, *http.Request)
	Delete func(http.ResponseWriter, *http.Request)
}

func NewTrigger(th TriggerAPI) *Trigger {
	return &Trigger{
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewTriggerList()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return th.List(r.Context(), params)
			})
		},
		Create: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewTriggerCreate()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return th.Create(r.Context(), params)
			})
		},
		Read: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewTriggerRead()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return th.Read(r.Context(), params)
			})
		},
		Update: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewTriggerUpdate()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return th.Update(r.Context(), params)
			})
		},
		Delete: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewTriggerDelete()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return th.Delete(r.Context(), params)
			})
		},
	}
}

func (th *Trigger) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Route("/trigger", func(r chi.Router) {
			r.Get("/", th.List)
			r.Post("/", th.Create)
			r.Get("/{triggerID}", th.Read)
			r.Post("/{triggerID}", th.Update)
			r.Delete("/{triggerID}", th.Delete)
		})
	})
}
