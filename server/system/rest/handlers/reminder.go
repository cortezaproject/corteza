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
	ReminderAPI interface {
		List(context.Context, *request.ReminderList) (interface{}, error)
		Create(context.Context, *request.ReminderCreate) (interface{}, error)
		Update(context.Context, *request.ReminderUpdate) (interface{}, error)
		Read(context.Context, *request.ReminderRead) (interface{}, error)
		Delete(context.Context, *request.ReminderDelete) (interface{}, error)
		Dismiss(context.Context, *request.ReminderDismiss) (interface{}, error)
		Snooze(context.Context, *request.ReminderSnooze) (interface{}, error)
	}

	// HTTP API interface
	Reminder struct {
		List    func(http.ResponseWriter, *http.Request)
		Create  func(http.ResponseWriter, *http.Request)
		Update  func(http.ResponseWriter, *http.Request)
		Read    func(http.ResponseWriter, *http.Request)
		Delete  func(http.ResponseWriter, *http.Request)
		Dismiss func(http.ResponseWriter, *http.Request)
		Snooze  func(http.ResponseWriter, *http.Request)
	}
)

func NewReminder(h ReminderAPI) *Reminder {
	return &Reminder{
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewReminderList()
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
			params := request.NewReminderCreate()
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
			params := request.NewReminderUpdate()
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
			params := request.NewReminderRead()
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
			params := request.NewReminderDelete()
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
		Dismiss: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewReminderDismiss()
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
		Snooze: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewReminderSnooze()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Snooze(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
	}
}

func (h Reminder) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/reminder/", h.List)
		r.Post("/reminder/", h.Create)
		r.Put("/reminder/{reminderID}", h.Update)
		r.Get("/reminder/{reminderID}", h.Read)
		r.Delete("/reminder/{reminderID}", h.Delete)
		r.Patch("/reminder/{reminderID}/dismiss", h.Dismiss)
		r.Patch("/reminder/{reminderID}/snooze", h.Snooze)
	})
}
