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
	ProjectReportAPI interface {
		Report(context.Context, *request.ProjectReportReport) (interface{}, error)
	}

	// HTTP API interface
	ProjectReport struct {
		Report func(http.ResponseWriter, *http.Request)
	}
)

func NewProjectReport(h ProjectReportAPI) *ProjectReport {
	return &ProjectReport{
		Report: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewProjectReportReport()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Report(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
	}
}

func (h ProjectReport) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/project-report/", h.Report)
	})
}
