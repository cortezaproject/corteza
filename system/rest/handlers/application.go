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
	ApplicationAPI interface {
		List(context.Context, *request.ApplicationList) (interface{}, error)
		Create(context.Context, *request.ApplicationCreate) (interface{}, error)
		Update(context.Context, *request.ApplicationUpdate) (interface{}, error)
		Upload(context.Context, *request.ApplicationUpload) (interface{}, error)
		FlagCreate(context.Context, *request.ApplicationFlagCreate) (interface{}, error)
		FlagDelete(context.Context, *request.ApplicationFlagDelete) (interface{}, error)
		Read(context.Context, *request.ApplicationRead) (interface{}, error)
		Delete(context.Context, *request.ApplicationDelete) (interface{}, error)
		Undelete(context.Context, *request.ApplicationUndelete) (interface{}, error)
		TriggerScript(context.Context, *request.ApplicationTriggerScript) (interface{}, error)
		Reorder(context.Context, *request.ApplicationReorder) (interface{}, error)
	}

	// HTTP API interface
	Application struct {
		List          func(http.ResponseWriter, *http.Request)
		Create        func(http.ResponseWriter, *http.Request)
		Update        func(http.ResponseWriter, *http.Request)
		Upload        func(http.ResponseWriter, *http.Request)
		FlagCreate    func(http.ResponseWriter, *http.Request)
		FlagDelete    func(http.ResponseWriter, *http.Request)
		Read          func(http.ResponseWriter, *http.Request)
		Delete        func(http.ResponseWriter, *http.Request)
		Undelete      func(http.ResponseWriter, *http.Request)
		TriggerScript func(http.ResponseWriter, *http.Request)
		Reorder       func(http.ResponseWriter, *http.Request)
	}
)

func NewApplication(h ApplicationAPI) *Application {
	return &Application{
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewApplicationList()
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
			params := request.NewApplicationCreate()
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
			params := request.NewApplicationUpdate()
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
		Upload: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewApplicationUpload()
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
		FlagCreate: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewApplicationFlagCreate()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.FlagCreate(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		FlagDelete: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewApplicationFlagDelete()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.FlagDelete(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Read: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewApplicationRead()
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
			params := request.NewApplicationDelete()
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
			params := request.NewApplicationUndelete()
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
		TriggerScript: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewApplicationTriggerScript()
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
		Reorder: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewApplicationReorder()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Reorder(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
	}
}

func (h Application) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/application/", h.List)
		r.Post("/application/", h.Create)
		r.Put("/application/{applicationID}", h.Update)
		r.Post("/application/upload", h.Upload)
		r.Post("/application/{applicationID}/flag/{ownedBy}/{flag}", h.FlagCreate)
		r.Delete("/application/{applicationID}/flag/{ownedBy}/{flag}", h.FlagDelete)
		r.Get("/application/{applicationID}", h.Read)
		r.Delete("/application/{applicationID}", h.Delete)
		r.Post("/application/{applicationID}/undelete", h.Undelete)
		r.Post("/application/{applicationID}/trigger", h.TriggerScript)
		r.Post("/application/reorder", h.Reorder)
	})
}
