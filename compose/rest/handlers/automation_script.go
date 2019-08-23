package handlers

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `automation_script.go`, `automation_script.util.go` or `automation_script_test.go` to
	implement your API calls, helper functions and tests. The file `automation_script.go`
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
type AutomationScriptAPI interface {
	List(context.Context, *request.AutomationScriptList) (interface{}, error)
	Create(context.Context, *request.AutomationScriptCreate) (interface{}, error)
	Read(context.Context, *request.AutomationScriptRead) (interface{}, error)
	Update(context.Context, *request.AutomationScriptUpdate) (interface{}, error)
	Delete(context.Context, *request.AutomationScriptDelete) (interface{}, error)
	Runnable(context.Context, *request.AutomationScriptRunnable) (interface{}, error)
	Run(context.Context, *request.AutomationScriptRun) (interface{}, error)
	Test(context.Context, *request.AutomationScriptTest) (interface{}, error)
}

// HTTP API interface
type AutomationScript struct {
	List     func(http.ResponseWriter, *http.Request)
	Create   func(http.ResponseWriter, *http.Request)
	Read     func(http.ResponseWriter, *http.Request)
	Update   func(http.ResponseWriter, *http.Request)
	Delete   func(http.ResponseWriter, *http.Request)
	Runnable func(http.ResponseWriter, *http.Request)
	Run      func(http.ResponseWriter, *http.Request)
	Test     func(http.ResponseWriter, *http.Request)
}

func NewAutomationScript(h AutomationScriptAPI) *AutomationScript {
	return &AutomationScript{
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAutomationScriptList()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("AutomationScript.List", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.List(r.Context(), params)
			if err != nil {
				logger.LogControllerError("AutomationScript.List", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("AutomationScript.List", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Create: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAutomationScriptCreate()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("AutomationScript.Create", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Create(r.Context(), params)
			if err != nil {
				logger.LogControllerError("AutomationScript.Create", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("AutomationScript.Create", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Read: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAutomationScriptRead()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("AutomationScript.Read", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Read(r.Context(), params)
			if err != nil {
				logger.LogControllerError("AutomationScript.Read", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("AutomationScript.Read", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Update: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAutomationScriptUpdate()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("AutomationScript.Update", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Update(r.Context(), params)
			if err != nil {
				logger.LogControllerError("AutomationScript.Update", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("AutomationScript.Update", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Delete: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAutomationScriptDelete()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("AutomationScript.Delete", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Delete(r.Context(), params)
			if err != nil {
				logger.LogControllerError("AutomationScript.Delete", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("AutomationScript.Delete", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Runnable: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAutomationScriptRunnable()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("AutomationScript.Runnable", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Runnable(r.Context(), params)
			if err != nil {
				logger.LogControllerError("AutomationScript.Runnable", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("AutomationScript.Runnable", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Run: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAutomationScriptRun()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("AutomationScript.Run", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Run(r.Context(), params)
			if err != nil {
				logger.LogControllerError("AutomationScript.Run", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("AutomationScript.Run", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Test: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAutomationScriptTest()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("AutomationScript.Test", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Test(r.Context(), params)
			if err != nil {
				logger.LogControllerError("AutomationScript.Test", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("AutomationScript.Test", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
	}
}

func (h AutomationScript) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/namespace/{namespaceID}/automation/script/", h.List)
		r.Post("/namespace/{namespaceID}/automation/script/", h.Create)
		r.Get("/namespace/{namespaceID}/automation/script/{scriptID}", h.Read)
		r.Post("/namespace/{namespaceID}/automation/script/{scriptID}", h.Update)
		r.Delete("/namespace/{namespaceID}/automation/script/{scriptID}", h.Delete)
		r.Get("/namespace/{namespaceID}/automation/script/runnable", h.Runnable)
		r.Post("/namespace/{namespaceID}/automation/script/{scriptID}/run", h.Run)
		r.Post("/namespace/{namespaceID}/automation/script/test", h.Test)
	})
}
