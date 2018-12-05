package handlers

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `workflow.go`, `workflow.util.go` or `workflow_test.go` to
	implement your API calls, helper functions and tests. The file `workflow.go`
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
type WorkflowAPI interface {
	List(context.Context, *request.WorkflowList) (interface{}, error)
	Create(context.Context, *request.WorkflowCreate) (interface{}, error)
	Get(context.Context, *request.WorkflowGet) (interface{}, error)
	Update(context.Context, *request.WorkflowUpdate) (interface{}, error)
	Delete(context.Context, *request.WorkflowDelete) (interface{}, error)
}

// HTTP API interface
type Workflow struct {
	List   func(http.ResponseWriter, *http.Request)
	Create func(http.ResponseWriter, *http.Request)
	Get    func(http.ResponseWriter, *http.Request)
	Update func(http.ResponseWriter, *http.Request)
	Delete func(http.ResponseWriter, *http.Request)
}

func NewWorkflow(wh WorkflowAPI) *Workflow {
	return &Workflow{
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewWorkflowList()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return wh.List(r.Context(), params)
			})
		},
		Create: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewWorkflowCreate()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return wh.Create(r.Context(), params)
			})
		},
		Get: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewWorkflowGet()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return wh.Get(r.Context(), params)
			})
		},
		Update: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewWorkflowUpdate()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return wh.Update(r.Context(), params)
			})
		},
		Delete: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewWorkflowDelete()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return wh.Delete(r.Context(), params)
			})
		},
	}
}

func (wh *Workflow) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Route("/workflow", func(r chi.Router) {
			r.Get("/", wh.List)
			r.Post("/", wh.Create)
			r.Get("/{workflowID}", wh.Get)
			r.Post("/{workflowID}", wh.Update)
			r.Delete("/{workflowID}", wh.Delete)
		})
	})
}
