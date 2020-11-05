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
	"github.com/cortezaproject/corteza-server/compose/rest/request"
	"github.com/cortezaproject/corteza-server/pkg/api"
	"github.com/go-chi/chi"
	"net/http"
)

type (
	// Internal API interface
	AutomationAPI interface {
		List(context.Context, *request.AutomationList) (interface{}, error)
		Bundle(context.Context, *request.AutomationBundle) (interface{}, error)
		TriggerScript(context.Context, *request.AutomationTriggerScript) (interface{}, error)
	}

	// HTTP API interface
	Automation struct {
		List          func(http.ResponseWriter, *http.Request)
		Bundle        func(http.ResponseWriter, *http.Request)
		TriggerScript func(http.ResponseWriter, *http.Request)
	}
)

func NewAutomation(h AutomationAPI) *Automation {
	return &Automation{
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAutomationList()
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
		Bundle: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAutomationBundle()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Bundle(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		TriggerScript: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAutomationTriggerScript()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.TriggerScript(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
	}
}

func (h Automation) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/automation/", h.List)
		r.Get("/automation/{bundle}-{type}.{ext}", h.Bundle)
		r.Post("/automation/trigger", h.TriggerScript)
	})
}
