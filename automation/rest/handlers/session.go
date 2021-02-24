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
	"github.com/cortezaproject/corteza-server/automation/rest/request"
	"github.com/cortezaproject/corteza-server/pkg/api"
	"github.com/go-chi/chi"
	"net/http"
)

type (
	// Internal API interface
	SessionAPI interface {
		List(context.Context, *request.SessionList) (interface{}, error)
		Read(context.Context, *request.SessionRead) (interface{}, error)
		Trace(context.Context, *request.SessionTrace) (interface{}, error)
		Delete(context.Context, *request.SessionDelete) (interface{}, error)
		ListPrompts(context.Context, *request.SessionListPrompts) (interface{}, error)
		ResumeState(context.Context, *request.SessionResumeState) (interface{}, error)
		DeleteState(context.Context, *request.SessionDeleteState) (interface{}, error)
	}

	// HTTP API interface
	Session struct {
		List        func(http.ResponseWriter, *http.Request)
		Read        func(http.ResponseWriter, *http.Request)
		Trace       func(http.ResponseWriter, *http.Request)
		Delete      func(http.ResponseWriter, *http.Request)
		ListPrompts func(http.ResponseWriter, *http.Request)
		ResumeState func(http.ResponseWriter, *http.Request)
		DeleteState func(http.ResponseWriter, *http.Request)
	}
)

func NewSession(h SessionAPI) *Session {
	return &Session{
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewSessionList()
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
			params := request.NewSessionRead()
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
		Trace: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewSessionTrace()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Trace(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Delete: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewSessionDelete()
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
		ListPrompts: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewSessionListPrompts()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.ListPrompts(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		ResumeState: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewSessionResumeState()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.ResumeState(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		DeleteState: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewSessionDeleteState()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.DeleteState(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
	}
}

func (h Session) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/sessions/", h.List)
		r.Get("/sessions/{sessionID}", h.Read)
		r.Get("/sessions/{sessionID}/trace", h.Trace)
		r.Delete("/sessions/{sessionID}", h.Delete)
		r.Get("/sessions/prompts", h.ListPrompts)
		r.Post("/sessions/{sessionID}/state/{stateID}", h.ResumeState)
		r.Delete("/sessions/{sessionID}/state/{stateID}", h.DeleteState)
	})
}
