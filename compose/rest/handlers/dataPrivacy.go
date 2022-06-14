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
		ListSensitiveData(context.Context, *request.DataPrivacyListSensitiveData) (interface{}, error)
	}

	// HTTP API interface
	DataPrivacy struct {
		ListSensitiveData func(http.ResponseWriter, *http.Request)
	}
)

func NewDataPrivacy(h DataPrivacyAPI) *DataPrivacy {
	return &DataPrivacy{
		ListSensitiveData: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewDataPrivacyListSensitiveData()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.ListSensitiveData(r.Context(), params)
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
		r.Get("/data-privacy/sensitive-data", h.ListSensitiveData)
	})
}
