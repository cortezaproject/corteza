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
	"github.com/cortezaproject/corteza-server/pkg/api"
	"github.com/cortezaproject/corteza-server/system/rest/request"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type (
	// Internal API interface
	DataPrivacyAPI interface {
		ListRequests(context.Context, *request.DataPrivacyListRequests) (interface{}, error)
		CreateRequest(context.Context, *request.DataPrivacyCreateRequest) (interface{}, error)
		ReadRequest(context.Context, *request.DataPrivacyReadRequest) (interface{}, error)
		ListResponsesOfRequest(context.Context, *request.DataPrivacyListResponsesOfRequest) (interface{}, error)
		CreateResponseForRequest(context.Context, *request.DataPrivacyCreateResponseForRequest) (interface{}, error)
	}

	// HTTP API interface
	DataPrivacy struct {
		ListRequests             func(http.ResponseWriter, *http.Request)
		CreateRequest            func(http.ResponseWriter, *http.Request)
		ReadRequest              func(http.ResponseWriter, *http.Request)
		ListResponsesOfRequest   func(http.ResponseWriter, *http.Request)
		CreateResponseForRequest func(http.ResponseWriter, *http.Request)
	}
)

func NewDataPrivacy(h DataPrivacyAPI) *DataPrivacy {
	return &DataPrivacy{
		ListRequests: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewDataPrivacyListRequests()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.ListRequests(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		CreateRequest: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewDataPrivacyCreateRequest()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.CreateRequest(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		ReadRequest: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewDataPrivacyReadRequest()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.ReadRequest(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		ListResponsesOfRequest: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewDataPrivacyListResponsesOfRequest()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.ListResponsesOfRequest(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		CreateResponseForRequest: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewDataPrivacyCreateResponseForRequest()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.CreateResponseForRequest(r.Context(), params)
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
		r.Get("/data-privacy/requests/", h.ListRequests)
		r.Post("/data-privacy/requests", h.CreateRequest)
		r.Get("/data-privacy/requests/{requestID}", h.ReadRequest)
		r.Get("/data-privacy/requests/{requestID}/responses", h.ListResponsesOfRequest)
		r.Get("/data-privacy/requests/{requestID}/responses", h.CreateResponseForRequest)
	})
}
