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
	"github.com/cortezaproject/corteza/server/pkg/api"
	"github.com/cortezaproject/corteza/server/system/rest/request"
	"github.com/go-chi/chi/v5"
	"net/http"
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
		Update: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewSettingsUpdate()
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
		Get: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewSettingsGet()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Get(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Set: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewSettingsSet()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Set(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Current: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewSettingsCurrent()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Current(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
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
