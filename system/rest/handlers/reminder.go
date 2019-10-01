package handlers

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `reminder.go`, `reminder.util.go` or `reminder_test.go` to
	implement your API calls, helper functions and tests. The file `reminder.go`
	is only generated the first time, and will not be overwritten if it exists.
*/

import (
	"context"

	"net/http"

	"github.com/go-chi/chi"
	"github.com/titpetric/factory/resputil"

	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/system/rest/request"
)

// Internal API interface
type ReminderAPI interface {
	List(context.Context, *request.ReminderList) (interface{}, error)
	Create(context.Context, *request.ReminderCreate) (interface{}, error)
	Update(context.Context, *request.ReminderUpdate) (interface{}, error)
	Read(context.Context, *request.ReminderRead) (interface{}, error)
	Delete(context.Context, *request.ReminderDelete) (interface{}, error)
	Dismiss(context.Context, *request.ReminderDismiss) (interface{}, error)
	Snooze(context.Context, *request.ReminderSnooze) (interface{}, error)
}

// HTTP API interface
type Reminder struct {
	List    func(http.ResponseWriter, *http.Request)
	Create  func(http.ResponseWriter, *http.Request)
	Update  func(http.ResponseWriter, *http.Request)
	Read    func(http.ResponseWriter, *http.Request)
	Delete  func(http.ResponseWriter, *http.Request)
	Dismiss func(http.ResponseWriter, *http.Request)
	Snooze  func(http.ResponseWriter, *http.Request)
}

func NewReminder(h ReminderAPI) *Reminder {
	return &Reminder{
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewReminderList()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Reminder.List", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.List(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Reminder.List", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Reminder.List", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Create: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewReminderCreate()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Reminder.Create", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Create(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Reminder.Create", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Reminder.Create", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Update: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewReminderUpdate()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Reminder.Update", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Update(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Reminder.Update", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Reminder.Update", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Read: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewReminderRead()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Reminder.Read", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Read(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Reminder.Read", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Reminder.Read", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Delete: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewReminderDelete()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Reminder.Delete", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Delete(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Reminder.Delete", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Reminder.Delete", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Dismiss: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewReminderDismiss()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Reminder.Dismiss", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Dismiss(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Reminder.Dismiss", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Reminder.Dismiss", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Snooze: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewReminderSnooze()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Reminder.Snooze", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Snooze(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Reminder.Snooze", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Reminder.Snooze", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
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
