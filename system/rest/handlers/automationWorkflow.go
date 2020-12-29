package handlers

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
//

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/api"
	"github.com/cortezaproject/corteza-server/system/rest/request"
	"github.com/go-chi/chi"
	"net/http"
)

type (
	// Internal API interface
	AutomationWorkflowAPI interface {
		List(context.Context, *request.AutomationWorkflowList) (interface{}, error)
		Create(context.Context, *request.AutomationWorkflowCreate) (interface{}, error)
		Update(context.Context, *request.AutomationWorkflowUpdate) (interface{}, error)
		Read(context.Context, *request.AutomationWorkflowRead) (interface{}, error)
		Delete(context.Context, *request.AutomationWorkflowDelete) (interface{}, error)
		Undelete(context.Context, *request.AutomationWorkflowUndelete) (interface{}, error)
		Test(context.Context, *request.AutomationWorkflowTest) (interface{}, error)
	}

	// HTTP API interface
	AutomationWorkflow struct {
		List     func(http.ResponseWriter, *http.Request)
		Create   func(http.ResponseWriter, *http.Request)
		Update   func(http.ResponseWriter, *http.Request)
		Read     func(http.ResponseWriter, *http.Request)
		Delete   func(http.ResponseWriter, *http.Request)
		Undelete func(http.ResponseWriter, *http.Request)
		Test     func(http.ResponseWriter, *http.Request)
	}
)

func NewAutomationWorkflow(h AutomationWorkflowAPI) *AutomationWorkflow {
	return &AutomationWorkflow{
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAutomationWorkflowList()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.List(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Create: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAutomationWorkflowCreate()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Create(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Update: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAutomationWorkflowUpdate()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Update(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Read: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAutomationWorkflowRead()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Read(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Delete: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAutomationWorkflowDelete()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Delete(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Undelete: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAutomationWorkflowUndelete()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Undelete(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Test: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAutomationWorkflowTest()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Test(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
	}
}

func (h AutomationWorkflow) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/automation/workflows/", h.List)
		r.Post("/automation/workflows/", h.Create)
		r.Put("/automation/workflows/{workflowID}", h.Update)
		r.Get("/automation/workflows/{workflowID}", h.Read)
		r.Delete("/automation/workflows/{workflowID}", h.Delete)
		r.Post("/automation/workflows/{workflowID}/undelete", h.Undelete)
		r.Post("/automation/workflows/{workflowID}/test", h.Test)
	})
}
