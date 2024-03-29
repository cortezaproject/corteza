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
	GatekeepAPI interface {
		Lock(context.Context, *request.GatekeepLock) (interface{}, error)
		Unlock(context.Context, *request.GatekeepUnlock) (interface{}, error)
		Check(context.Context, *request.GatekeepCheck) (interface{}, error)
	}

	// HTTP API interface
	Gatekeep struct {
		Lock   func(http.ResponseWriter, *http.Request)
		Unlock func(http.ResponseWriter, *http.Request)
		Check  func(http.ResponseWriter, *http.Request)
	}
)

func NewGatekeep(h GatekeepAPI) *Gatekeep {
	return &Gatekeep{
		Lock: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewGatekeepLock()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Lock(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Unlock: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewGatekeepUnlock()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Unlock(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Check: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewGatekeepCheck()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Check(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
	}
}

func (h Gatekeep) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Post("/gatekeep/lock", h.Lock)
		r.Post("/gatekeep/unlock", h.Unlock)
		r.Post("/gatekeep/{lockID}/check", h.Check)
	})
}
