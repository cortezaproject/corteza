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
	NodeHandshakeAPI interface {
		Initialize(context.Context, *request.NodeHandshakeInitialize) (interface{}, error)
	}

	// HTTP API interface
	NodeHandshake struct {
		Initialize func(http.ResponseWriter, *http.Request)
	}
)

func NewNodeHandshake(h NodeHandshakeAPI) *NodeHandshake {
	return &NodeHandshake{
		Initialize: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewNodeHandshakeInitialize()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("NodeHandshake.Initialize", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Initialize(r.Context(), params)
			if err != nil {
				logger.LogControllerError("NodeHandshake.Initialize", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("NodeHandshake.Initialize", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
	}
}

func (h NodeHandshake) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Post("/nodes/{nodeID}/handshake", h.Initialize)
	})
}
