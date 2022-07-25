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
	"github.com/cortezaproject/corteza-server/compose/rest/request"
	"github.com/cortezaproject/corteza-server/pkg/api"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type (
	// Internal API interface
	DataPrivacyAPI interface {
		RecordList(context.Context, *request.DataPrivacyRecordList) (interface{}, error)
		ModuleList(context.Context, *request.DataPrivacyModuleList) (interface{}, error)
	}

	// HTTP API interface
	DataPrivacy struct {
		RecordList func(http.ResponseWriter, *http.Request)
		ModuleList func(http.ResponseWriter, *http.Request)
	}
)

func NewDataPrivacy(h DataPrivacyAPI) *DataPrivacy {
	return &DataPrivacy{
		RecordList: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewDataPrivacyRecordList()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.RecordList(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		ModuleList: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewDataPrivacyModuleList()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.ModuleList(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
	}
}

func (h DataPrivacy) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/data-privacy/record", h.RecordList)
		r.Get("/data-privacy/module", h.ModuleList)
	})
}
