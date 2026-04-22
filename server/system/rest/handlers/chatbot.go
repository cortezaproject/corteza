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
	"github.com/crusttech/human/server/pkg/api"
	"github.com/crusttech/human/server/system/rest/request"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type (
	// Internal API interface
	ChatbotAPI interface {
		List(context.Context, *request.ChatbotList) (interface{}, error)
		Create(context.Context, *request.ChatbotCreate) (interface{}, error)
		Read(context.Context, *request.ChatbotRead) (interface{}, error)
		Update(context.Context, *request.ChatbotUpdate) (interface{}, error)
		Delete(context.Context, *request.ChatbotDelete) (interface{}, error)
		Undelete(context.Context, *request.ChatbotUndelete) (interface{}, error)
		RegenerateWidgetKey(context.Context, *request.ChatbotRegenerateWidgetKey) (interface{}, error)
	}

	// HTTP API interface
	Chatbot struct {
		List                func(http.ResponseWriter, *http.Request)
		Create              func(http.ResponseWriter, *http.Request)
		Read                func(http.ResponseWriter, *http.Request)
		Update              func(http.ResponseWriter, *http.Request)
		Delete              func(http.ResponseWriter, *http.Request)
		Undelete            func(http.ResponseWriter, *http.Request)
		RegenerateWidgetKey func(http.ResponseWriter, *http.Request)
	}
)

func NewChatbot(h ChatbotAPI) *Chatbot {
	return &Chatbot{
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChatbotList()
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
			params := request.NewChatbotCreate()
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
			params := request.NewChatbotRead()
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
			params := request.NewChatbotUpdate()
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
			params := request.NewChatbotDelete()
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
			params := request.NewChatbotUndelete()
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
		RegenerateWidgetKey: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChatbotRegenerateWidgetKey()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.RegenerateWidgetKey(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
	}
}

func (h Chatbot) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/chatbots/", h.List)
		r.Post("/chatbots", h.Create)
		r.Get("/chatbots/{chatbotID}", h.Read)
		r.Put("/chatbots/{chatbotID}", h.Update)
		r.Delete("/chatbots/{chatbotID}", h.Delete)
		r.Post("/chatbots/{chatbotID}/undelete", h.Undelete)
		r.Post("/chatbots/{chatbotID}/regenerate-key", h.RegenerateWidgetKey)
	})
}
