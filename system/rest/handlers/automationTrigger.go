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
	AutomationTriggerAPI interface {
		List(context.Context, *request.AutomationTriggerList) (interface{}, error)
		Create(context.Context, *request.AutomationTriggerCreate) (interface{}, error)
		Update(context.Context, *request.AutomationTriggerUpdate) (interface{}, error)
		Read(context.Context, *request.AutomationTriggerRead) (interface{}, error)
		Delete(context.Context, *request.AutomationTriggerDelete) (interface{}, error)
		Undelete(context.Context, *request.AutomationTriggerUndelete) (interface{}, error)
	}

	// HTTP API interface
	AutomationTrigger struct {
		List     func(http.ResponseWriter, *http.Request)
		Create   func(http.ResponseWriter, *http.Request)
		Update   func(http.ResponseWriter, *http.Request)
		Read     func(http.ResponseWriter, *http.Request)
		Delete   func(http.ResponseWriter, *http.Request)
		Undelete func(http.ResponseWriter, *http.Request)
	}
)

func NewAutomationTrigger(h AutomationTriggerAPI) *AutomationTrigger {
	return &AutomationTrigger{
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAutomationTriggerList()
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
			params := request.NewAutomationTriggerCreate()
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
			params := request.NewAutomationTriggerUpdate()
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
			params := request.NewAutomationTriggerRead()
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
			params := request.NewAutomationTriggerDelete()
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
			params := request.NewAutomationTriggerUndelete()
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
	}
}

func (h AutomationTrigger) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/automation/triggers/", h.List)
		r.Post("/automation/triggers/", h.Create)
		r.Put("/automation/triggers/{triggerID}", h.Update)
		r.Get("/automation/triggers/{triggerID}", h.Read)
		r.Delete("/automation/triggers/{triggerID}", h.Delete)
		r.Post("/automation/triggers/{triggerID}/undelete", h.Undelete)
	})
}
