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
	NamespaceAPI interface {
		List(context.Context, *request.NamespaceList) (interface{}, error)
		Create(context.Context, *request.NamespaceCreate) (interface{}, error)
		Read(context.Context, *request.NamespaceRead) (interface{}, error)
		Update(context.Context, *request.NamespaceUpdate) (interface{}, error)
		Delete(context.Context, *request.NamespaceDelete) (interface{}, error)
		Upload(context.Context, *request.NamespaceUpload) (interface{}, error)
		Clone(context.Context, *request.NamespaceClone) (interface{}, error)
		Export(context.Context, *request.NamespaceExport) (interface{}, error)
		ImportInit(context.Context, *request.NamespaceImportInit) (interface{}, error)
		ImportRun(context.Context, *request.NamespaceImportRun) (interface{}, error)
		TriggerScript(context.Context, *request.NamespaceTriggerScript) (interface{}, error)
		ListTranslations(context.Context, *request.NamespaceListTranslations) (interface{}, error)
		UpdateTranslations(context.Context, *request.NamespaceUpdateTranslations) (interface{}, error)
	}

	// HTTP API interface
	Namespace struct {
		List               func(http.ResponseWriter, *http.Request)
		Create             func(http.ResponseWriter, *http.Request)
		Read               func(http.ResponseWriter, *http.Request)
		Update             func(http.ResponseWriter, *http.Request)
		Delete             func(http.ResponseWriter, *http.Request)
		Upload             func(http.ResponseWriter, *http.Request)
		Clone              func(http.ResponseWriter, *http.Request)
		Export             func(http.ResponseWriter, *http.Request)
		ImportInit         func(http.ResponseWriter, *http.Request)
		ImportRun          func(http.ResponseWriter, *http.Request)
		TriggerScript      func(http.ResponseWriter, *http.Request)
		ListTranslations   func(http.ResponseWriter, *http.Request)
		UpdateTranslations func(http.ResponseWriter, *http.Request)
	}
)

func NewNamespace(h NamespaceAPI) *Namespace {
	return &Namespace{
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewNamespaceList()
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
			params := request.NewNamespaceCreate()
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
			params := request.NewNamespaceRead()
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
			params := request.NewNamespaceUpdate()
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
			params := request.NewNamespaceDelete()
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
			params := request.NewNamespaceUpload()
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
		Clone: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewNamespaceClone()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Clone(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Export: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewNamespaceExport()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Export(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		ImportInit: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewNamespaceImportInit()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.ImportInit(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		ImportRun: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewNamespaceImportRun()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.ImportRun(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		TriggerScript: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewNamespaceTriggerScript()
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
			params := request.NewNamespaceListTranslations()
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
			params := request.NewNamespaceUpdateTranslations()
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

func (h Namespace) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/namespace/", h.List)
		r.Post("/namespace/", h.Create)
		r.Get("/namespace/{namespaceID}", h.Read)
		r.Post("/namespace/{namespaceID}", h.Update)
		r.Delete("/namespace/{namespaceID}", h.Delete)
		r.Post("/namespace/upload", h.Upload)
		r.Post("/namespace/{namespaceID}/clone", h.Clone)
		r.Get("/namespace/{namespaceID}/export/{filename}.zip", h.Export)
		r.Post("/namespace/import", h.ImportInit)
		r.Post("/namespace/import/{sessionID}", h.ImportRun)
		r.Post("/namespace/{namespaceID}/trigger", h.TriggerScript)
		r.Get("/namespace/{namespaceID}/translation", h.ListTranslations)
		r.Patch("/namespace/{namespaceID}/translation", h.UpdateTranslations)
	})
}
