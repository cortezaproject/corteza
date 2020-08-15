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
	SettingsAPI interface {
		List(context.Context, *request.SettingsList) (interface{}, error)
		Update(context.Context, *request.SettingsUpdate) (interface{}, error)
		Get(context.Context, *request.SettingsGet) (interface{}, error)
		Set(context.Context, *request.SettingsSet) (interface{}, error)
		Current(context.Context, *request.SettingsCurrent) (interface{}, error)
	}

	// HTTP API interface
	Settings struct {
		List    func(http.ResponseWriter, *http.Request)
		Update  func(http.ResponseWriter, *http.Request)
		Get     func(http.ResponseWriter, *http.Request)
		Set     func(http.ResponseWriter, *http.Request)
		Current func(http.ResponseWriter, *http.Request)
	}
)

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
		Current: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewSettingsCurrent()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Settings.Current", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Current(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Settings.Current", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Settings.Current", r, params.Auditable())
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
		r.Post("/settings/{key}", h.Set)
		r.Get("/settings/current", h.Current)
	})
}
