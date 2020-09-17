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
	SyncStructureAPI interface {
		ReadExposed(context.Context, *request.SyncStructureReadExposed) (interface{}, error)
		Remove(context.Context, *request.SyncStructureRemove) (interface{}, error)
	}

	// HTTP API interface
	SyncStructure struct {
		ReadExposed func(http.ResponseWriter, *http.Request)
		Remove      func(http.ResponseWriter, *http.Request)
	}
)

func NewSyncStructure(h SyncStructureAPI) *SyncStructure {
	return &SyncStructure{
		ReadExposed: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewSyncStructureReadExposed()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("SyncStructure.ReadExposed", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.ReadExposed(r.Context(), params)
			if err != nil {
				logger.LogControllerError("SyncStructure.ReadExposed", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("SyncStructure.ReadExposed", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Remove: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewSyncStructureRemove()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("SyncStructure.Remove", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Remove(r.Context(), params)
			if err != nil {
				logger.LogControllerError("SyncStructure.Remove", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("SyncStructure.Remove", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
	}
}

func (h SyncStructure) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/nodes/{nodeID}/modules/{moduleID}/exposed", h.ReadExposed)
		r.Delete("/nodes/{nodeID}/modules/{moduleID}/exposed", h.Remove)
	})
}
