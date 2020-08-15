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
	"github.com/go-chi/chi"
	"github.com/titpetric/factory/resputil"
	"net/http"

	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/system/rest/request"
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
				logger.LogParamError("Automation.List", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.List(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Automation.List", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Automation.List", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Bundle: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAutomationBundle()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Automation.Bundle", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Bundle(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Automation.Bundle", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Automation.Bundle", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		TriggerScript: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAutomationTriggerScript()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Automation.TriggerScript", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.TriggerScript(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Automation.TriggerScript", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Automation.TriggerScript", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
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
