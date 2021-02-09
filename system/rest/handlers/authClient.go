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
	"github.com/go-chi/chi"
	"net/http"
)

type (
	// Internal API interface
	AuthClientAPI interface {
		List(context.Context, *request.AuthClientList) (interface{}, error)
		Create(context.Context, *request.AuthClientCreate) (interface{}, error)
		Update(context.Context, *request.AuthClientUpdate) (interface{}, error)
		Read(context.Context, *request.AuthClientRead) (interface{}, error)
		Delete(context.Context, *request.AuthClientDelete) (interface{}, error)
		Undelete(context.Context, *request.AuthClientUndelete) (interface{}, error)
		RegenerateSecret(context.Context, *request.AuthClientRegenerateSecret) (interface{}, error)
		ExposeSecret(context.Context, *request.AuthClientExposeSecret) (interface{}, error)
	}

	// HTTP API interface
	AuthClient struct {
		List             func(http.ResponseWriter, *http.Request)
		Create           func(http.ResponseWriter, *http.Request)
		Update           func(http.ResponseWriter, *http.Request)
		Read             func(http.ResponseWriter, *http.Request)
		Delete           func(http.ResponseWriter, *http.Request)
		Undelete         func(http.ResponseWriter, *http.Request)
		RegenerateSecret func(http.ResponseWriter, *http.Request)
		ExposeSecret     func(http.ResponseWriter, *http.Request)
	}
)

func NewAuthClient(h AuthClientAPI) *AuthClient {
	return &AuthClient{
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAuthClientList()
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
			params := request.NewAuthClientCreate()
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
		Update: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAuthClientUpdate()
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
		Read: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAuthClientRead()
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
		Delete: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAuthClientDelete()
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
			params := request.NewAuthClientUndelete()
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
		RegenerateSecret: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAuthClientRegenerateSecret()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.RegenerateSecret(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		ExposeSecret: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewAuthClientExposeSecret()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.ExposeSecret(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
	}
}

func (h AuthClient) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/auth/clients/", h.List)
		r.Post("/auth/clients/", h.Create)
		r.Put("/auth/clients/{clientID}", h.Update)
		r.Get("/auth/clients/{clientID}", h.Read)
		r.Delete("/auth/clients/{clientID}", h.Delete)
		r.Post("/auth/clients/{clientID}/undelete", h.Undelete)
		r.Post("/auth/clients/{clientID}/secret", h.RegenerateSecret)
		r.Get("/auth/clients/{clientID}/secret", h.ExposeSecret)
	})
}
