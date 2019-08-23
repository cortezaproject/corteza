package handlers

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `automation_trigger.go`, `automation_trigger.util.go` or `automation_trigger_test.go` to
	implement your API calls, helper functions and tests. The file `automation_trigger.go`
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
type AutomationTriggerAPI interface {
	List(context.Context, *request.AutomationTriggerList) (interface{}, error)
	Create(context.Context, *request.AutomationTriggerCreate) (interface{}, error)
	Read(context.Context, *request.AutomationTriggerRead) (interface{}, error)
	Update(context.Context, *request.AutomationTriggerUpdate) (interface{}, error)
	Delete(context.Context, *request.AutomationTriggerDelete) (interface{}, error)
}

// HTTP API interface
type AutomationTrigger struct {
	List   func(http.ResponseWriter, *http.Request)
	Create func(http.ResponseWriter, *http.Request)
	Read   func(http.ResponseWriter, *http.Request)
	Update func(http.ResponseWriter, *http.Request)
	Delete func(http.ResponseWriter, *http.Request)
}

func NewAutomationTrigger(h AutomationTriggerAPI) *AutomationTrigger {
	return &AutomationTrigger{
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAutomationTriggerList()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("AutomationTrigger.List", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.List(r.Context(), params)
			if err != nil {
				logger.LogControllerError("AutomationTrigger.List", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("AutomationTrigger.List", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Create: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAutomationTriggerCreate()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("AutomationTrigger.Create", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Create(r.Context(), params)
			if err != nil {
				logger.LogControllerError("AutomationTrigger.Create", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("AutomationTrigger.Create", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Read: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAutomationTriggerRead()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("AutomationTrigger.Read", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Read(r.Context(), params)
			if err != nil {
				logger.LogControllerError("AutomationTrigger.Read", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("AutomationTrigger.Read", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Update: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAutomationTriggerUpdate()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("AutomationTrigger.Update", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Update(r.Context(), params)
			if err != nil {
				logger.LogControllerError("AutomationTrigger.Update", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("AutomationTrigger.Update", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Delete: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAutomationTriggerDelete()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("AutomationTrigger.Delete", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Delete(r.Context(), params)
			if err != nil {
				logger.LogControllerError("AutomationTrigger.Delete", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("AutomationTrigger.Delete", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
	}
}

func (h AutomationTrigger) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/namespace/{namespaceID}/automation/script/{scriptID}/trigger/", h.List)
		r.Post("/namespace/{namespaceID}/automation/script/{scriptID}/trigger/", h.Create)
		r.Get("/namespace/{namespaceID}/automation/script/{scriptID}/trigger/{triggerID}", h.Read)
		r.Post("/namespace/{namespaceID}/automation/script/{scriptID}/trigger/{triggerID}", h.Update)
		r.Delete("/namespace/{namespaceID}/automation/script/{scriptID}/trigger/{triggerID}", h.Delete)
	})
}
