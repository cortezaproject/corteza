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
	"github.com/cortezaproject/corteza-server/federation/rest/request"
	"github.com/cortezaproject/corteza-server/pkg/api"
	"github.com/go-chi/chi"
	"net/http"
)

type (
	// Internal API interface
	ManageStructureAPI interface {
		ReadExposed(context.Context, *request.ManageStructureReadExposed) (interface{}, error)
		CreateExposed(context.Context, *request.ManageStructureCreateExposed) (interface{}, error)
		UpdateExposed(context.Context, *request.ManageStructureUpdateExposed) (interface{}, error)
		RemoveExposed(context.Context, *request.ManageStructureRemoveExposed) (interface{}, error)
		ReadShared(context.Context, *request.ManageStructureReadShared) (interface{}, error)
		CreateMappings(context.Context, *request.ManageStructureCreateMappings) (interface{}, error)
		ReadMappings(context.Context, *request.ManageStructureReadMappings) (interface{}, error)
		ListAll(context.Context, *request.ManageStructureListAll) (interface{}, error)
	}

	// HTTP API interface
	ManageStructure struct {
		ReadExposed    func(http.ResponseWriter, *http.Request)
		CreateExposed  func(http.ResponseWriter, *http.Request)
		UpdateExposed  func(http.ResponseWriter, *http.Request)
		RemoveExposed  func(http.ResponseWriter, *http.Request)
		ReadShared     func(http.ResponseWriter, *http.Request)
		CreateMappings func(http.ResponseWriter, *http.Request)
		ReadMappings   func(http.ResponseWriter, *http.Request)
		ListAll        func(http.ResponseWriter, *http.Request)
	}
)

func NewManageStructure(h ManageStructureAPI) *ManageStructure {
	return &ManageStructure{
		ReadExposed: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewManageStructureReadExposed()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.ReadExposed(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		CreateExposed: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewManageStructureCreateExposed()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.CreateExposed(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		UpdateExposed: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewManageStructureUpdateExposed()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.UpdateExposed(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		RemoveExposed: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewManageStructureRemoveExposed()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.RemoveExposed(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		ReadShared: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewManageStructureReadShared()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.ReadShared(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		CreateMappings: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewManageStructureCreateMappings()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.CreateMappings(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		ReadMappings: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewManageStructureReadMappings()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.ReadMappings(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		ListAll: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewManageStructureListAll()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.ListAll(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
	}
}

func (h ManageStructure) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/nodes/{nodeID}/modules/{moduleID}/exposed", h.ReadExposed)
		r.Put("/nodes/{nodeID}/modules/", h.CreateExposed)
		r.Post("/nodes/{nodeID}/modules/{moduleID}/exposed", h.UpdateExposed)
		r.Delete("/nodes/{nodeID}/modules/{moduleID}/exposed", h.RemoveExposed)
		r.Get("/nodes/{nodeID}/modules/{moduleID}/shared", h.ReadShared)
		r.Put("/nodes/{nodeID}/modules/{moduleID}/mapped", h.CreateMappings)
		r.Get("/nodes/{nodeID}/modules/{moduleID}/mapped", h.ReadMappings)
		r.Get("/nodes/{nodeID}/modules/", h.ListAll)
	})
}
