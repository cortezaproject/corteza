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
	PairRequestAPI interface {
		RequestPairing(context.Context, *request.PairRequestRequestPairing) (interface{}, error)
	}

	// HTTP API interface
	PairRequest struct {
		RequestPairing func(http.ResponseWriter, *http.Request)
	}
)

func NewPairRequest(h PairRequestAPI) *PairRequest {
	return &PairRequest{
		RequestPairing: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewPairRequestRequestPairing()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("PairRequest.RequestPairing", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.RequestPairing(r.Context(), params)
			if err != nil {
				logger.LogControllerError("PairRequest.RequestPairing", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("PairRequest.RequestPairing", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
	}
}

func (h PairRequest) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Post("/node/pair/request/", h.RequestPairing)
	})
}
