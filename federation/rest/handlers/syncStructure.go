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
		ReadExposedAll(context.Context, *request.SyncStructureReadExposedAll) (interface{}, error)
	}

	// HTTP API interface
	SyncStructure struct {
		ReadExposedAll func(http.ResponseWriter, *http.Request)
	}
)

func NewSyncStructure(h SyncStructureAPI) *SyncStructure {
	return &SyncStructure{
		ReadExposedAll: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewSyncStructureReadExposedAll()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("SyncStructure.ReadExposedAll", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.ReadExposedAll(r.Context(), params)
			if err != nil {
				logger.LogControllerError("SyncStructure.ReadExposedAll", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("SyncStructure.ReadExposedAll", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
	}
}

func (h SyncStructure) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/nodes/{nodeID}/modules/exposed/", h.ReadExposedAll)
	})
}
