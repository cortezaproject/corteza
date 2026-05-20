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
	ChatbotSessionAPI interface {
		List(context.Context, *request.ChatbotSessionList) (interface{}, error)
		Read(context.Context, *request.ChatbotSessionRead) (interface{}, error)
		AdvanceStep(context.Context, *request.ChatbotSessionAdvanceStep) (interface{}, error)
		Close(context.Context, *request.ChatbotSessionClose) (interface{}, error)
		HandoffAccept(context.Context, *request.ChatbotSessionHandoffAccept) (interface{}, error)
		OperatorMessage(context.Context, *request.ChatbotSessionOperatorMessage) (interface{}, error)
		HandoffComplete(context.Context, *request.ChatbotSessionHandoffComplete) (interface{}, error)
		Stream(http.ResponseWriter, *http.Request)
	}

	// HTTP API interface
	ChatbotSession struct {
		List            func(http.ResponseWriter, *http.Request)
		Read            func(http.ResponseWriter, *http.Request)
		AdvanceStep     func(http.ResponseWriter, *http.Request)
		Close           func(http.ResponseWriter, *http.Request)
		HandoffAccept   func(http.ResponseWriter, *http.Request)
		OperatorMessage func(http.ResponseWriter, *http.Request)
		HandoffComplete func(http.ResponseWriter, *http.Request)
		Stream          func(http.ResponseWriter, *http.Request)
	}
)

func NewChatbotSession(h ChatbotSessionAPI) *ChatbotSession {
	return &ChatbotSession{
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChatbotSessionList()
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
		Read: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChatbotSessionRead()
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
		AdvanceStep: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChatbotSessionAdvanceStep()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.AdvanceStep(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Close: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChatbotSessionClose()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Close(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		HandoffAccept: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChatbotSessionHandoffAccept()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.HandoffAccept(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		OperatorMessage: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChatbotSessionOperatorMessage()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.OperatorMessage(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		HandoffComplete: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChatbotSessionHandoffComplete()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.HandoffComplete(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Stream: h.Stream,
	}
}

func (h ChatbotSession) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/chatbots/sessions/", h.List)
		r.Get("/chatbots/sessions/{sessionID}", h.Read)
		r.Post("/chatbots/sessions/{sessionID}/advance-step", h.AdvanceStep)
		r.Post("/chatbots/sessions/{sessionID}/close", h.Close)
		r.Post("/chatbots/sessions/{sessionID}/handoff-accept", h.HandoffAccept)
		r.Post("/chatbots/sessions/{sessionID}/operator-message", h.OperatorMessage)
		r.Post("/chatbots/sessions/{sessionID}/handoff-complete", h.HandoffComplete)
		r.Get("/chatbots/sessions/{sessionID}/stream", h.Stream)
	})
}
