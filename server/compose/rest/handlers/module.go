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
	ModuleAPI interface {
		List(context.Context, *request.ModuleList) (interface{}, error)
		Create(context.Context, *request.ModuleCreate) (interface{}, error)
		Read(context.Context, *request.ModuleRead) (interface{}, error)
		Update(context.Context, *request.ModuleUpdate) (interface{}, error)
		Delete(context.Context, *request.ModuleDelete) (interface{}, error)
		TriggerScript(context.Context, *request.ModuleTriggerScript) (interface{}, error)
		ListTranslations(context.Context, *request.ModuleListTranslations) (interface{}, error)
		UpdateTranslations(context.Context, *request.ModuleUpdateTranslations) (interface{}, error)
	}

	// HTTP API interface
	Module struct {
		List               func(http.ResponseWriter, *http.Request)
		Create             func(http.ResponseWriter, *http.Request)
		Read               func(http.ResponseWriter, *http.Request)
		Update             func(http.ResponseWriter, *http.Request)
		Delete             func(http.ResponseWriter, *http.Request)
		TriggerScript      func(http.ResponseWriter, *http.Request)
		ListTranslations   func(http.ResponseWriter, *http.Request)
		UpdateTranslations func(http.ResponseWriter, *http.Request)
	}
)

func NewModule(h ModuleAPI) *Module {
	return &Module{
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewModuleList()
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
			params := request.NewModuleCreate()
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
			params := request.NewModuleRead()
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
			params := request.NewModuleUpdate()
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
			params := request.NewModuleDelete()
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
		TriggerScript: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewModuleTriggerScript()
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
			params := request.NewModuleListTranslations()
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
			params := request.NewModuleUpdateTranslations()
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

func (h Module) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/namespace/{namespaceID}/module/", h.List)
		r.Post("/namespace/{namespaceID}/module/", h.Create)
		r.Get("/namespace/{namespaceID}/module/{moduleID}", h.Read)
		r.Post("/namespace/{namespaceID}/module/{moduleID}", h.Update)
		r.Delete("/namespace/{namespaceID}/module/{moduleID}", h.Delete)
		r.Post("/namespace/{namespaceID}/module/{moduleID}/trigger", h.TriggerScript)
		r.Get("/namespace/{namespaceID}/module/{moduleID}/translation", h.ListTranslations)
		r.Patch("/namespace/{namespaceID}/module/{moduleID}/translation", h.UpdateTranslations)
	})
}
