package handlers

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `auth.go`, `auth.util.go` or `auth_test.go` to
	implement your API calls, helper functions and tests. The file `auth.go`
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
type AuthAPI interface {
	Settings(context.Context, *request.AuthSettings) (interface{}, error)
	Check(context.Context, *request.AuthCheck) (interface{}, error)
	ExchangeAuthToken(context.Context, *request.AuthExchangeAuthToken) (interface{}, error)
	Logout(context.Context, *request.AuthLogout) (interface{}, error)
}

// HTTP API interface
type Auth struct {
	Settings          func(http.ResponseWriter, *http.Request)
	Check             func(http.ResponseWriter, *http.Request)
	ExchangeAuthToken func(http.ResponseWriter, *http.Request)
	Logout            func(http.ResponseWriter, *http.Request)
}

func NewAuth(h AuthAPI) *Auth {
	return &Auth{
		Settings: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAuthSettings()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Auth.Settings", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Settings(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Auth.Settings", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Auth.Settings", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Check: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAuthCheck()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Auth.Check", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Check(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Auth.Check", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Auth.Check", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		ExchangeAuthToken: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAuthExchangeAuthToken()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Auth.ExchangeAuthToken", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.ExchangeAuthToken(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Auth.ExchangeAuthToken", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Auth.ExchangeAuthToken", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Logout: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAuthLogout()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Auth.Logout", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Logout(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Auth.Logout", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Auth.Logout", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
	}
}

func (h Auth) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/auth/", h.Settings)
		r.Get("/auth/check", h.Check)
		r.Post("/auth/exchange", h.ExchangeAuthToken)
		r.Get("/auth/logout", h.Logout)
	})
}
