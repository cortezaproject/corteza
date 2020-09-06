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
	PairAPI interface {
		ApprovePairing(context.Context, *request.PairApprovePairing) (interface{}, error)
		CompletePairing(context.Context, *request.PairCompletePairing) (interface{}, error)
	}

	// HTTP API interface
	Pair struct {
		ApprovePairing  func(http.ResponseWriter, *http.Request)
		CompletePairing func(http.ResponseWriter, *http.Request)
	}
)

func NewPair(h PairAPI) *Pair {
	return &Pair{
		ApprovePairing: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewPairApprovePairing()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Pair.ApprovePairing", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.ApprovePairing(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Pair.ApprovePairing", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Pair.ApprovePairing", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		CompletePairing: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewPairCompletePairing()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Pair.CompletePairing", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.CompletePairing(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Pair.CompletePairing", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Pair.CompletePairing", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
	}
}

func (h Pair) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Post("/node/pair/approve", h.ApprovePairing)
		r.Post("/node/pair/complete", h.CompletePairing)
	})
}
