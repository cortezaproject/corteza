package handlers

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `settings.go`, `settings.util.go` or `settings_test.go` to
	implement your API calls, helper functions and tests. The file `settings.go`
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
type SettingsAPI interface {
	List(context.Context, *request.SettingsList) (interface{}, error)
	Update(context.Context, *request.SettingsUpdate) (interface{}, error)
	Get(context.Context, *request.SettingsGet) (interface{}, error)
	Set(context.Context, *request.SettingsSet) (interface{}, error)
}

// HTTP API interface
type Settings struct {
	List   func(http.ResponseWriter, *http.Request)
	Update func(http.ResponseWriter, *http.Request)
	Get    func(http.ResponseWriter, *http.Request)
	Set    func(http.ResponseWriter, *http.Request)
}

func NewSettings(h SettingsAPI) *Settings {
	return &Settings{
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewSettingsList()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Settings.List", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.List(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Settings.List", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Settings.List", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Update: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewSettingsUpdate()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Settings.Update", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Update(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Settings.Update", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Settings.Update", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Get: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewSettingsGet()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Settings.Get", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Get(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Settings.Get", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Settings.Get", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Set: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewSettingsSet()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Settings.Set", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Set(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Settings.Set", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Settings.Set", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
	}
}

func (h Settings) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/settings/", h.List)
		r.Patch("/settings/", h.Update)
		r.Get("/settings/{key}", h.Get)
		r.Put("/settings/{key}", h.Set)
	})
}
