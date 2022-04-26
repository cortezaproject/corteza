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
	ChartAPI interface {
		List(context.Context, *request.ChartList) (interface{}, error)
		Create(context.Context, *request.ChartCreate) (interface{}, error)
		Read(context.Context, *request.ChartRead) (interface{}, error)
		Update(context.Context, *request.ChartUpdate) (interface{}, error)
		Delete(context.Context, *request.ChartDelete) (interface{}, error)
		ListTranslations(context.Context, *request.ChartListTranslations) (interface{}, error)
		UpdateTranslations(context.Context, *request.ChartUpdateTranslations) (interface{}, error)
	}

	// HTTP API interface
	Chart struct {
		List               func(http.ResponseWriter, *http.Request)
		Create             func(http.ResponseWriter, *http.Request)
		Read               func(http.ResponseWriter, *http.Request)
		Update             func(http.ResponseWriter, *http.Request)
		Delete             func(http.ResponseWriter, *http.Request)
		ListTranslations   func(http.ResponseWriter, *http.Request)
		UpdateTranslations func(http.ResponseWriter, *http.Request)
	}
)

func NewChart(h ChartAPI) *Chart {
	return &Chart{
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChartList()
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
			params := request.NewChartCreate()
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
		Read: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChartRead()
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
		Update: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChartUpdate()
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
		Delete: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChartDelete()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Delete(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		ListTranslations: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChartListTranslations()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.ListTranslations(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		UpdateTranslations: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChartUpdateTranslations()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.UpdateTranslations(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
	}
}

func (h Chart) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/namespace/{namespaceID}/chart/", h.List)
		r.Post("/namespace/{namespaceID}/chart/", h.Create)
		r.Get("/namespace/{namespaceID}/chart/{chartID}", h.Read)
		r.Post("/namespace/{namespaceID}/chart/{chartID}", h.Update)
		r.Delete("/namespace/{namespaceID}/chart/{chartID}", h.Delete)
		r.Get("/namespace/{namespaceID}/chart/{chartID}/translation", h.ListTranslations)
		r.Patch("/namespace/{namespaceID}/chart/{chartID}/translation", h.UpdateTranslations)
	})
}
