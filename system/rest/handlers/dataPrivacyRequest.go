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
	DataPrivacyRequestAPI interface {
		List(context.Context, *request.DataPrivacyRequestList) (interface{}, error)
		Create(context.Context, *request.DataPrivacyRequestCreate) (interface{}, error)
		Update(context.Context, *request.DataPrivacyRequestUpdate) (interface{}, error)
		UpdateStatus(context.Context, *request.DataPrivacyRequestUpdateStatus) (interface{}, error)
		Read(context.Context, *request.DataPrivacyRequestRead) (interface{}, error)
		ListResponses(context.Context, *request.DataPrivacyRequestListResponses) (interface{}, error)
		CreateResponse(context.Context, *request.DataPrivacyRequestCreateResponse) (interface{}, error)
	}

	// HTTP API interface
	DataPrivacyRequest struct {
		List           func(http.ResponseWriter, *http.Request)
		Create         func(http.ResponseWriter, *http.Request)
		Update         func(http.ResponseWriter, *http.Request)
		UpdateStatus   func(http.ResponseWriter, *http.Request)
		Read           func(http.ResponseWriter, *http.Request)
		ListResponses  func(http.ResponseWriter, *http.Request)
		CreateResponse func(http.ResponseWriter, *http.Request)
	}
)

func NewDataPrivacyRequest(h DataPrivacyRequestAPI) *DataPrivacyRequest {
	return &DataPrivacyRequest{
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewDataPrivacyRequestList()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.List(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Create: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewDataPrivacyRequestCreate()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Create(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Update: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewDataPrivacyRequestUpdate()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Update(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		UpdateStatus: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewDataPrivacyRequestUpdateStatus()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.UpdateStatus(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Read: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewDataPrivacyRequestRead()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Read(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		ListResponses: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewDataPrivacyRequestListResponses()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.ListResponses(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		CreateResponse: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewDataPrivacyRequestCreateResponse()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.CreateResponse(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
	}
}

func (h DataPrivacyRequest) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/data-privacy/requests/", h.List)
		r.Post("/data-privacy/requests/", h.Create)
		r.Put("/data-privacy/requests/{requestID}", h.Update)
		r.Put("/data-privacy/requests/{requestID}/status/{status}", h.UpdateStatus)
		r.Get("/data-privacy/requests/{requestID}", h.Read)
		r.Get("/data-privacy/requests/{requestID}/responses", h.ListResponses)
		r.Get("/data-privacy/requests/{requestID}/responses", h.CreateResponse)
	})
}
