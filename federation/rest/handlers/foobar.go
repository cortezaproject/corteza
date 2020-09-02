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

	"github.com/cortezaproject/corteza-server/federation/rest/request"
	"github.com/cortezaproject/corteza-server/pkg/logger"
)

type (
	// Internal API interface
	FoobarAPI interface {
		Foobar(context.Context, *request.FoobarFoobar) (interface{}, error)
	}

	// HTTP API interface
	Foobar struct {
		Foobar func(http.ResponseWriter, *http.Request)
	}
)

func NewFoobar(h FoobarAPI) *Foobar {
	return &Foobar{
		Foobar: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewFoobarFoobar()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Foobar.Foobar", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Foobar(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Foobar.Foobar", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Foobar.Foobar", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
	}
}

func (h Foobar) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/foobar/", h.Foobar)
	})
}
