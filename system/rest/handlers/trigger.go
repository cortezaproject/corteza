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

	"net/http"

	"github.com/go-chi/chi"
	"github.com/titpetric/factory/resputil"

	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/system/rest/request"
)

// Internal API interface
type TriggerAPI interface {
	List(context.Context, *request.TriggerList) (interface{}, error)
	Fire(context.Context, *request.TriggerFire) (interface{}, error)
}

// HTTP API interface
type Trigger struct {
	List func(http.ResponseWriter, *http.Request)
	Fire func(http.ResponseWriter, *http.Request)
}

func NewTrigger(h TriggerAPI) *Trigger {
	return &Trigger{
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewTriggerList()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Trigger.List", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.List(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Trigger.List", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Trigger.List", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Fire: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewTriggerFire()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Trigger.Fire", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Fire(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Trigger.Fire", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Trigger.Fire", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
	}
}

func (h Trigger) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/trigger/", h.List)
		r.Post("/trigger/", h.Fire)
	})
}
