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

	"github.com/crusttech/crust/system/rest/request"
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

func NewSettings(sh SettingsAPI) *Settings {
	return &Settings{
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewSettingsList()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return sh.List(r.Context(), params)
			})
		},
		Update: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewSettingsUpdate()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return sh.Update(r.Context(), params)
			})
		},
		Get: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewSettingsGet()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return sh.Get(r.Context(), params)
			})
		},
		Set: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewSettingsSet()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return sh.Set(r.Context(), params)
			})
		},
	}
}

func (sh *Settings) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/settings/", sh.List)
		r.Patch("/settings/", sh.Update)
		r.Get("/settings/{key}", sh.Get)
		r.Put("/settings/{key}", sh.Set)
	})
}
