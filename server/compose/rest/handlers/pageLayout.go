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
	"github.com/cortezaproject/corteza/server/compose/rest/request"
	"github.com/cortezaproject/corteza/server/pkg/api"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type (
	// Internal API interface
	PageLayoutAPI interface {
		ListNamespace(context.Context, *request.PageLayoutListNamespace) (interface{}, error)
		List(context.Context, *request.PageLayoutList) (interface{}, error)
		Create(context.Context, *request.PageLayoutCreate) (interface{}, error)
		Read(context.Context, *request.PageLayoutRead) (interface{}, error)
		Update(context.Context, *request.PageLayoutUpdate) (interface{}, error)
		Reorder(context.Context, *request.PageLayoutReorder) (interface{}, error)
		Delete(context.Context, *request.PageLayoutDelete) (interface{}, error)
		Undelete(context.Context, *request.PageLayoutUndelete) (interface{}, error)
		ListTranslations(context.Context, *request.PageLayoutListTranslations) (interface{}, error)
		UpdateTranslations(context.Context, *request.PageLayoutUpdateTranslations) (interface{}, error)
	}

	// HTTP API interface
	PageLayout struct {
		ListNamespace      func(http.ResponseWriter, *http.Request)
		List               func(http.ResponseWriter, *http.Request)
		Create             func(http.ResponseWriter, *http.Request)
		Read               func(http.ResponseWriter, *http.Request)
		Update             func(http.ResponseWriter, *http.Request)
		Reorder            func(http.ResponseWriter, *http.Request)
		Delete             func(http.ResponseWriter, *http.Request)
		Undelete           func(http.ResponseWriter, *http.Request)
		ListTranslations   func(http.ResponseWriter, *http.Request)
		UpdateTranslations func(http.ResponseWriter, *http.Request)
	}
)

func NewPageLayout(h PageLayoutAPI) *PageLayout {
	return &PageLayout{
		ListNamespace: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewPageLayoutListNamespace()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.ListNamespace(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewPageLayoutList()
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
			params := request.NewPageLayoutCreate()
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
			params := request.NewPageLayoutRead()
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
			params := request.NewPageLayoutUpdate()
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
		Reorder: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewPageLayoutReorder()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Reorder(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Delete: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewPageLayoutDelete()
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
		Undelete: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewPageLayoutUndelete()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Undelete(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		ListTranslations: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewPageLayoutListTranslations()
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
			params := request.NewPageLayoutUpdateTranslations()
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

func (h PageLayout) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/namespace/{namespaceID}/page-layout", h.ListNamespace)
		r.Get("/namespace/{namespaceID}/page/{pageID}/layout/", h.List)
		r.Post("/namespace/{namespaceID}/page/{pageID}/layout/", h.Create)
		r.Get("/namespace/{namespaceID}/page/{pageID}/layout/{pageLayoutID}", h.Read)
		r.Post("/namespace/{namespaceID}/page/{pageID}/layout/{pageLayoutID}", h.Update)
		r.Post("/namespace/{namespaceID}/page/{pageID}/layout/reorder", h.Reorder)
		r.Delete("/namespace/{namespaceID}/page/{pageID}/layout/{pageLayoutID}", h.Delete)
		r.Post("/namespace/{namespaceID}/page/{pageID}/layout/{pageLayoutID}/undelete", h.Undelete)
		r.Get("/namespace/{namespaceID}/page/{pageID}/layout/{pageLayoutID}/translation", h.ListTranslations)
		r.Patch("/namespace/{namespaceID}/page/{pageID}/layout/{pageLayoutID}/translation", h.UpdateTranslations)
	})
}
