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
	"github.com/go-chi/chi/v5"
	"net/http"
)

type (
	// Internal API interface
	DalSensitivityLevelAPI interface {
		List(context.Context, *request.DalSensitivityLevelList) (interface{}, error)
		Create(context.Context, *request.DalSensitivityLevelCreate) (interface{}, error)
		Update(context.Context, *request.DalSensitivityLevelUpdate) (interface{}, error)
		Read(context.Context, *request.DalSensitivityLevelRead) (interface{}, error)
		Delete(context.Context, *request.DalSensitivityLevelDelete) (interface{}, error)
		Undelete(context.Context, *request.DalSensitivityLevelUndelete) (interface{}, error)
	}

	// HTTP API interface
	DalSensitivityLevel struct {
		List     func(http.ResponseWriter, *http.Request)
		Create   func(http.ResponseWriter, *http.Request)
		Update   func(http.ResponseWriter, *http.Request)
		Read     func(http.ResponseWriter, *http.Request)
		Delete   func(http.ResponseWriter, *http.Request)
		Undelete func(http.ResponseWriter, *http.Request)
	}
)

func NewDalSensitivityLevel(h DalSensitivityLevelAPI) *DalSensitivityLevel {
	return &DalSensitivityLevel{
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewDalSensitivityLevelList()
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
			params := request.NewDalSensitivityLevelCreate()
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
			params := request.NewDalSensitivityLevelUpdate()
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
			params := request.NewDalSensitivityLevelRead()
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
			params := request.NewDalSensitivityLevelDelete()
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
			params := request.NewDalSensitivityLevelUndelete()
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
	}
}

func (h DalSensitivityLevel) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/dal/sensitivity-levels/", h.List)
		r.Post("/dal/sensitivity-levels/", h.Create)
		r.Put("/dal/sensitivity-levels/{sensitivityLevelID}", h.Update)
		r.Get("/dal/sensitivity-levels/{sensitivityLevelID}", h.Read)
		r.Delete("/dal/sensitivity-levels/{sensitivityLevelID}", h.Delete)
		r.Post("/dal/sensitivity-levels/{sensitivityLevelID}/undelete", h.Undelete)
	})
}
