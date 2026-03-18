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
	"github.com/cortezaproject/corteza/server/pkg/api"
	"github.com/cortezaproject/corteza/server/system/rest/request"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type (
	// Internal API interface
	LlmProviderAPI interface {
		List(context.Context, *request.LlmProviderList) (interface{}, error)
		Create(context.Context, *request.LlmProviderCreate) (interface{}, error)
		Read(context.Context, *request.LlmProviderRead) (interface{}, error)
		Update(context.Context, *request.LlmProviderUpdate) (interface{}, error)
		Delete(context.Context, *request.LlmProviderDelete) (interface{}, error)
		Models(context.Context, *request.LlmProviderModels) (interface{}, error)
		Validate(context.Context, *request.LlmProviderValidate) (interface{}, error)
	}

	// HTTP API interface
	LlmProvider struct {
		List     func(http.ResponseWriter, *http.Request)
		Create   func(http.ResponseWriter, *http.Request)
		Read     func(http.ResponseWriter, *http.Request)
		Update   func(http.ResponseWriter, *http.Request)
		Delete   func(http.ResponseWriter, *http.Request)
		Models   func(http.ResponseWriter, *http.Request)
		Validate func(http.ResponseWriter, *http.Request)
	}
)

func NewLlmProvider(h LlmProviderAPI) *LlmProvider {
	return &LlmProvider{
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewLlmProviderList()
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
			params := request.NewLlmProviderCreate()
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
			params := request.NewLlmProviderRead()
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
			params := request.NewLlmProviderUpdate()
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
			params := request.NewLlmProviderDelete()
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
		Models: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewLlmProviderModels()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Models(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Validate: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewLlmProviderValidate()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Validate(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
	}
}

func (h LlmProvider) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/llm-providers/", h.List)
		r.Post("/llm-providers/", h.Create)
		r.Get("/llm-providers/{llmProviderID}", h.Read)
		r.Put("/llm-providers/{llmProviderID}", h.Update)
		r.Delete("/llm-providers/{llmProviderID}", h.Delete)
		r.Get("/llm-providers/{llmProviderID}/models", h.Models)
		r.Post("/llm-providers/{llmProviderID}/validate", h.Validate)
	})
}
