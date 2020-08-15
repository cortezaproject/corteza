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
	ActivityAPI interface {
		Send(context.Context, *request.ActivitySend) (interface{}, error)
	}

	// HTTP API interface
	Activity struct {
		Send func(http.ResponseWriter, *http.Request)
	}
)

func NewActivity(h ActivityAPI) *Activity {
	return &Activity{
		Send: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewActivitySend()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Activity.Send", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Send(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Activity.Send", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Activity.Send", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
	}
}

func (h Activity) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Post("/activity/", h.Send)
	})
}
