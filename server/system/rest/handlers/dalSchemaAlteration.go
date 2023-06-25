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
	DalSchemaAlterationAPI interface {
		List(context.Context, *request.DalSchemaAlterationList) (interface{}, error)
		Read(context.Context, *request.DalSchemaAlterationRead) (interface{}, error)
		Apply(context.Context, *request.DalSchemaAlterationApply) (interface{}, error)
		Dismiss(context.Context, *request.DalSchemaAlterationDismiss) (interface{}, error)
	}

	// HTTP API interface
	DalSchemaAlteration struct {
		List    func(http.ResponseWriter, *http.Request)
		Read    func(http.ResponseWriter, *http.Request)
		Apply   func(http.ResponseWriter, *http.Request)
		Dismiss func(http.ResponseWriter, *http.Request)
	}
)

func NewDalSchemaAlteration(h DalSchemaAlterationAPI) *DalSchemaAlteration {
	return &DalSchemaAlteration{
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewDalSchemaAlterationList()
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
			params := request.NewDalSchemaAlterationRead()
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
		Apply: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewDalSchemaAlterationApply()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Apply(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Dismiss: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewDalSchemaAlterationDismiss()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Dismiss(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
	}
}

func (h DalSchemaAlteration) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/dal/schema/alterations/", h.List)
		r.Get("/dal/schema/alterations/{alterationID}", h.Read)
		r.Post("/dal/schema/alterations/apply", h.Apply)
		r.Post("/dal/schema/alterations/dismiss", h.Dismiss)
	})
}
