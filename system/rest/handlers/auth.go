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
	AuthAPI interface {
		Settings(context.Context, *request.AuthSettings) (interface{}, error)
		Check(context.Context, *request.AuthCheck) (interface{}, error)
		Impersonate(context.Context, *request.AuthImpersonate) (interface{}, error)
		ExchangeAuthToken(context.Context, *request.AuthExchangeAuthToken) (interface{}, error)
		Logout(context.Context, *request.AuthLogout) (interface{}, error)
	}

	// HTTP API interface
	Auth struct {
		Settings          func(http.ResponseWriter, *http.Request)
		Check             func(http.ResponseWriter, *http.Request)
		Impersonate       func(http.ResponseWriter, *http.Request)
		ExchangeAuthToken func(http.ResponseWriter, *http.Request)
		Logout            func(http.ResponseWriter, *http.Request)
	}
)

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
		Impersonate: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAuthImpersonate()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Auth.Impersonate", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Impersonate(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Auth.Impersonate", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Auth.Impersonate", r, params.Auditable())
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
		r.Post("/auth/impersonate", h.Impersonate)
		r.Post("/auth/exchange", h.ExchangeAuthToken)
		r.Get("/auth/logout", h.Logout)
	})
}
