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
	"github.com/cortezaproject/corteza/server/federation/rest/request"
	"github.com/cortezaproject/corteza/server/pkg/api"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type (
	// Internal API interface
	SyncStructureAPI interface {
		ReadExposedInternal(context.Context, *request.SyncStructureReadExposedInternal) (interface{}, error)
		ReadExposedSocial(context.Context, *request.SyncStructureReadExposedSocial) (interface{}, error)
	}

	// HTTP API interface
	SyncStructure struct {
		ReadExposedInternal func(http.ResponseWriter, *http.Request)
		ReadExposedSocial   func(http.ResponseWriter, *http.Request)
	}
)

func NewSyncStructure(h SyncStructureAPI) *SyncStructure {
	return &SyncStructure{
		ReadExposedInternal: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewSyncStructureReadExposedInternal()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.ReadExposedInternal(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		ReadExposedSocial: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewSyncStructureReadExposedSocial()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.ReadExposedSocial(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
	}
}

func (h SyncStructure) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/nodes/{nodeID}/modules/exposed/", h.ReadExposedInternal)
		r.Get("/nodes/{nodeID}/modules/exposed/activity-stream", h.ReadExposedSocial)
	})
}
