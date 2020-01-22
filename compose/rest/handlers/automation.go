package handlers

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `automation.go`, `automation.util.go` or `automation_test.go` to
	implement your API calls, helper functions and tests. The file `automation.go`
	is only generated the first time, and will not be overwritten if it exists.
*/

import (
	"context"

	"net/http"

	"github.com/go-chi/chi"
	"github.com/titpetric/factory/resputil"

	"github.com/cortezaproject/corteza-server/compose/rest/request"
	"github.com/cortezaproject/corteza-server/pkg/logger"
)

// Internal API interface
type AutomationAPI interface {
	List(context.Context, *request.AutomationList) (interface{}, error)
	Bundle(context.Context, *request.AutomationBundle) (interface{}, error)
	TriggerScript(context.Context, *request.AutomationTriggerScript) (interface{}, error)
}

// HTTP API interface
type Automation struct {
	List          func(http.ResponseWriter, *http.Request)
	Bundle        func(http.ResponseWriter, *http.Request)
	TriggerScript func(http.ResponseWriter, *http.Request)
}

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
