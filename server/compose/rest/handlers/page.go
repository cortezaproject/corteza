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
	PageAPI interface {
		List(context.Context, *request.PageList) (interface{}, error)
		Create(context.Context, *request.PageCreate) (interface{}, error)
		Read(context.Context, *request.PageRead) (interface{}, error)
		Tree(context.Context, *request.PageTree) (interface{}, error)
		Update(context.Context, *request.PageUpdate) (interface{}, error)
		Reorder(context.Context, *request.PageReorder) (interface{}, error)
		Delete(context.Context, *request.PageDelete) (interface{}, error)
		Upload(context.Context, *request.PageUpload) (interface{}, error)
		TriggerScript(context.Context, *request.PageTriggerScript) (interface{}, error)
		ListTranslations(context.Context, *request.PageListTranslations) (interface{}, error)
		UpdateTranslations(context.Context, *request.PageUpdateTranslations) (interface{}, error)
		UpdateIcon(context.Context, *request.PageUpdateIcon) (interface{}, error)
	}

	// HTTP API interface
	Page struct {
		List               func(http.ResponseWriter, *http.Request)
		Create             func(http.ResponseWriter, *http.Request)
		Read               func(http.ResponseWriter, *http.Request)
		Tree               func(http.ResponseWriter, *http.Request)
		Update             func(http.ResponseWriter, *http.Request)
		Reorder            func(http.ResponseWriter, *http.Request)
		Delete             func(http.ResponseWriter, *http.Request)
		Upload             func(http.ResponseWriter, *http.Request)
		TriggerScript      func(http.ResponseWriter, *http.Request)
		ListTranslations   func(http.ResponseWriter, *http.Request)
		UpdateTranslations func(http.ResponseWriter, *http.Request)
		UpdateIcon         func(http.ResponseWriter, *http.Request)
	}
)

func NewPage(h PageAPI) *Page {
	return &Page{
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewPageList()
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
			params := request.NewPageCreate()
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
			params := request.NewPageRead()
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
		Tree: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewPageTree()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Tree(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Update: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewPageUpdate()
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
			params := request.NewPageReorder()
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
			params := request.NewPageDelete()
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
		Upload: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewPageUpload()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Upload(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		TriggerScript: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewPageTriggerScript()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.TriggerScript(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		ListTranslations: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewPageListTranslations()
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
			params := request.NewPageUpdateTranslations()
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
		UpdateIcon: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewPageUpdateIcon()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.UpdateIcon(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
	}
}

func (h Page) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/namespace/{namespaceID}/page/", h.List)
		r.Post("/namespace/{namespaceID}/page/", h.Create)
		r.Get("/namespace/{namespaceID}/page/{pageID}", h.Read)
		r.Get("/namespace/{namespaceID}/page/tree", h.Tree)
		r.Post("/namespace/{namespaceID}/page/{pageID}", h.Update)
		r.Post("/namespace/{namespaceID}/page/{selfID}/reorder", h.Reorder)
		r.Delete("/namespace/{namespaceID}/page/{pageID}", h.Delete)
		r.Post("/namespace/{namespaceID}/page/{pageID}/attachment", h.Upload)
		r.Post("/namespace/{namespaceID}/page/{pageID}/trigger", h.TriggerScript)
		r.Get("/namespace/{namespaceID}/page/{pageID}/translation", h.ListTranslations)
		r.Patch("/namespace/{namespaceID}/page/{pageID}/translation", h.UpdateTranslations)
		r.Patch("/namespace/{namespaceID}/page/{pageID}/icon", h.UpdateIcon)
	})
}
