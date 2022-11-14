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
	LocaleAPI interface {
		ListResource(context.Context, *request.LocaleListResource) (interface{}, error)
		CreateResource(context.Context, *request.LocaleCreateResource) (interface{}, error)
		UpdateResource(context.Context, *request.LocaleUpdateResource) (interface{}, error)
		ReadResource(context.Context, *request.LocaleReadResource) (interface{}, error)
		DeleteResource(context.Context, *request.LocaleDeleteResource) (interface{}, error)
		UndeleteResource(context.Context, *request.LocaleUndeleteResource) (interface{}, error)
		List(context.Context, *request.LocaleList) (interface{}, error)
		Get(context.Context, *request.LocaleGet) (interface{}, error)
	}

	// HTTP API interface
	Locale struct {
		ListResource     func(http.ResponseWriter, *http.Request)
		CreateResource   func(http.ResponseWriter, *http.Request)
		UpdateResource   func(http.ResponseWriter, *http.Request)
		ReadResource     func(http.ResponseWriter, *http.Request)
		DeleteResource   func(http.ResponseWriter, *http.Request)
		UndeleteResource func(http.ResponseWriter, *http.Request)
		List             func(http.ResponseWriter, *http.Request)
		Get              func(http.ResponseWriter, *http.Request)
	}
)

func NewLocale(h LocaleAPI) *Locale {
	return &Locale{
		ListResource: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewLocaleListResource()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.ListResource(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		CreateResource: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewLocaleCreateResource()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.CreateResource(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		UpdateResource: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewLocaleUpdateResource()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.UpdateResource(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		ReadResource: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewLocaleReadResource()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.ReadResource(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		DeleteResource: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewLocaleDeleteResource()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.DeleteResource(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		UndeleteResource: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewLocaleUndeleteResource()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.UndeleteResource(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewLocaleList()
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
		Get: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewLocaleGet()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Get(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
	}
}

func (h Locale) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/locale/resource", h.ListResource)
		r.Post("/locale/resource", h.CreateResource)
		r.Put("/locale/resource/{translationID}", h.UpdateResource)
		r.Get("/locale/resource/{translationID}", h.ReadResource)
		r.Delete("/locale/resource/{translationID}", h.DeleteResource)
		r.Post("/locale/resource/{translationID}/undelete", h.UndeleteResource)
		r.Get("/locale/", h.List)
		r.Get("/locale/{lang}/{application}", h.Get)
	})
}
