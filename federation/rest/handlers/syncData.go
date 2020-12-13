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
	SyncDataAPI interface {
		ReadExposedAll(context.Context, *request.SyncDataReadExposedAll) (interface{}, error)
		ReadExposed(context.Context, *request.SyncDataReadExposed) (interface{}, error)
	}

	// HTTP API interface
	SyncData struct {
		ReadExposedAll func(http.ResponseWriter, *http.Request)
		ReadExposed    func(http.ResponseWriter, *http.Request)
	}
)

func NewSyncData(h SyncDataAPI) *SyncData {
	return &SyncData{
		ReadExposedAll: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewSyncDataReadExposedAll()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.ReadExposedAll(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		ReadExposed: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewSyncDataReadExposed()
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
	}
}

func (h SyncData) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/nodes/{nodeID}/modules/exposed/records/", h.ReadExposedAll)
		r.Get("/nodes/{nodeID}/modules/{moduleID}/records/", h.ReadExposed)
	})
}
