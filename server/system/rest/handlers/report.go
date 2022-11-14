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
	ReportAPI interface {
		List(context.Context, *request.ReportList) (interface{}, error)
		Create(context.Context, *request.ReportCreate) (interface{}, error)
		Update(context.Context, *request.ReportUpdate) (interface{}, error)
		Read(context.Context, *request.ReportRead) (interface{}, error)
		Delete(context.Context, *request.ReportDelete) (interface{}, error)
		Undelete(context.Context, *request.ReportUndelete) (interface{}, error)
		Describe(context.Context, *request.ReportDescribe) (interface{}, error)
		Run(context.Context, *request.ReportRun) (interface{}, error)
	}

	// HTTP API interface
	Report struct {
		List     func(http.ResponseWriter, *http.Request)
		Create   func(http.ResponseWriter, *http.Request)
		Update   func(http.ResponseWriter, *http.Request)
		Read     func(http.ResponseWriter, *http.Request)
		Delete   func(http.ResponseWriter, *http.Request)
		Undelete func(http.ResponseWriter, *http.Request)
		Describe func(http.ResponseWriter, *http.Request)
		Run      func(http.ResponseWriter, *http.Request)
	}
)

func NewReport(h ReportAPI) *Report {
	return &Report{
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewReportList()
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
			params := request.NewReportCreate()
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
			params := request.NewReportUpdate()
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
			params := request.NewReportRead()
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
			params := request.NewReportDelete()
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
			params := request.NewReportUndelete()
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
		Describe: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewReportDescribe()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Describe(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Run: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewReportRun()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Run(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
	}
}

func (h Report) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/reports/", h.List)
		r.Post("/reports/", h.Create)
		r.Put("/reports/{reportID}", h.Update)
		r.Get("/reports/{reportID}", h.Read)
		r.Delete("/reports/{reportID}", h.Delete)
		r.Post("/reports/{reportID}/undelete", h.Undelete)
		r.Post("/reports/describe", h.Describe)
		r.Post("/reports/{reportID}/run", h.Run)
	})
}
