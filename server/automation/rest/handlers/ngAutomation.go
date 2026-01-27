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
	"github.com/cortezaproject/corteza/server/automation/rest/request"
	"github.com/cortezaproject/corteza/server/pkg/api"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type (
	// Internal API interface
	NgAutomationAPI interface {
		List(context.Context, *request.NgAutomationList) (interface{}, error)
		Create(context.Context, *request.NgAutomationCreate) (interface{}, error)
		Update(context.Context, *request.NgAutomationUpdate) (interface{}, error)
		Read(context.Context, *request.NgAutomationRead) (interface{}, error)
		Delete(context.Context, *request.NgAutomationDelete) (interface{}, error)
		Undelete(context.Context, *request.NgAutomationUndelete) (interface{}, error)
		Test(context.Context, *request.NgAutomationTest) (interface{}, error)
		Exec(context.Context, *request.NgAutomationExec) (interface{}, error)
	}

	// HTTP API interface
	NgAutomation struct {
		List     func(http.ResponseWriter, *http.Request)
		Create   func(http.ResponseWriter, *http.Request)
		Update   func(http.ResponseWriter, *http.Request)
		Read     func(http.ResponseWriter, *http.Request)
		Delete   func(http.ResponseWriter, *http.Request)
		Undelete func(http.ResponseWriter, *http.Request)
		Test     func(http.ResponseWriter, *http.Request)
		Exec     func(http.ResponseWriter, *http.Request)
	}
)

func NewNgAutomation(h NgAutomationAPI) *NgAutomation {
	return &NgAutomation{
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewNgAutomationList()
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
			params := request.NewNgAutomationCreate()
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
			params := request.NewNgAutomationUpdate()
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
			params := request.NewNgAutomationRead()
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
			params := request.NewNgAutomationDelete()
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
			params := request.NewNgAutomationUndelete()
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
			params := request.NewNgAutomationTest()
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
		Exec: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewNgAutomationExec()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Exec(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
	}
}

func (h NgAutomation) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/ng-automation/", h.List)
		r.Post("/ng-automation/", h.Create)
		r.Put("/ng-automation/{automationID}", h.Update)
		r.Get("/ng-automation/{automationID}", h.Read)
		r.Delete("/ng-automation/{automationID}", h.Delete)
		r.Post("/ng-automation/{automationID}/undelete", h.Undelete)
		r.Post("/ng-automation/{automationID}/test", h.Test)
		r.Post("/ng-automation/{automationID}/exec", h.Exec)
	})
}
