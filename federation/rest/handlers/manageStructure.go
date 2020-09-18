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
	ManageStructureAPI interface {
		ReadExposed(context.Context, *request.ManageStructureReadExposed) (interface{}, error)
		CreateExposed(context.Context, *request.ManageStructureCreateExposed) (interface{}, error)
		Remove(context.Context, *request.ManageStructureRemove) (interface{}, error)
		ReadShared(context.Context, *request.ManageStructureReadShared) (interface{}, error)
		ListAll(context.Context, *request.ManageStructureListAll) (interface{}, error)
	}

	// HTTP API interface
	ManageStructure struct {
		ReadExposed   func(http.ResponseWriter, *http.Request)
		CreateExposed func(http.ResponseWriter, *http.Request)
		Remove        func(http.ResponseWriter, *http.Request)
		ReadShared    func(http.ResponseWriter, *http.Request)
		ListAll       func(http.ResponseWriter, *http.Request)
	}
)

func NewManageStructure(h ManageStructureAPI) *ManageStructure {
	return &ManageStructure{
		ReadExposed: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewManageStructureReadExposed()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("ManageStructure.ReadExposed", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.ReadExposed(r.Context(), params)
			if err != nil {
				logger.LogControllerError("ManageStructure.ReadExposed", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("ManageStructure.ReadExposed", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		CreateExposed: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewManageStructureCreateExposed()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("ManageStructure.CreateExposed", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.CreateExposed(r.Context(), params)
			if err != nil {
				logger.LogControllerError("ManageStructure.CreateExposed", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("ManageStructure.CreateExposed", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Remove: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewManageStructureRemove()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("ManageStructure.Remove", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Remove(r.Context(), params)
			if err != nil {
				logger.LogControllerError("ManageStructure.Remove", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("ManageStructure.Remove", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		ReadShared: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewManageStructureReadShared()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("ManageStructure.ReadShared", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.ReadShared(r.Context(), params)
			if err != nil {
				logger.LogControllerError("ManageStructure.ReadShared", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("ManageStructure.ReadShared", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		ListAll: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewManageStructureListAll()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("ManageStructure.ListAll", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.ListAll(r.Context(), params)
			if err != nil {
				logger.LogControllerError("ManageStructure.ListAll", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("ManageStructure.ListAll", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
	}
}

func (h ManageStructure) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/nodes/{nodeID}/modules/{moduleID}/exposed", h.ReadExposed)
		r.Put("/nodes/{nodeID}/modules/{moduleID}/exposed", h.CreateExposed)
		r.Delete("/nodes/{nodeID}/modules/{moduleID}/exposed", h.Remove)
		r.Get("/nodes/{nodeID}/modules/{moduleID}/shared", h.ReadShared)
		r.Get("/nodes/{nodeID}/modules", h.ListAll)
	})
}
