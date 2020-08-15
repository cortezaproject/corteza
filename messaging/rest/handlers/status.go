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

	"github.com/cortezaproject/corteza-server/messaging/rest/request"
	"github.com/cortezaproject/corteza-server/pkg/logger"
)

type (
	// Internal API interface
	StatusAPI interface {
		List(context.Context, *request.StatusList) (interface{}, error)
		Set(context.Context, *request.StatusSet) (interface{}, error)
		Delete(context.Context, *request.StatusDelete) (interface{}, error)
	}

	// HTTP API interface
	Status struct {
		List   func(http.ResponseWriter, *http.Request)
		Set    func(http.ResponseWriter, *http.Request)
		Delete func(http.ResponseWriter, *http.Request)
	}
)

func NewStatus(h StatusAPI) *Status {
	return &Status{
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewStatusList()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Status.List", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.List(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Status.List", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Status.List", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Set: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewStatusSet()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Status.Set", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Set(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Status.Set", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Status.Set", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Delete: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewStatusDelete()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Status.Delete", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Delete(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Status.Delete", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Status.Delete", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
	}
}

func (h Status) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/status/", h.List)
		r.Post("/status/", h.Set)
		r.Delete("/status/", h.Delete)
	})
}
