package handlers

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `chart.go`, `chart.util.go` or `chart_test.go` to
	implement your API calls, helper functions and tests. The file `chart.go`
	is only generated the first time, and will not be overwritten if it exists.
*/

import (
	"context"

	"net/http"

	"github.com/go-chi/chi"
	"github.com/titpetric/factory/resputil"

	"github.com/cortezaproject/corteza-server/compose/rest/request"
	"github.com/cortezaproject/corteza-server/pkg/logger"
)

// Internal API interface
type ChartAPI interface {
	List(context.Context, *request.ChartList) (interface{}, error)
	Create(context.Context, *request.ChartCreate) (interface{}, error)
	Read(context.Context, *request.ChartRead) (interface{}, error)
	Update(context.Context, *request.ChartUpdate) (interface{}, error)
	Delete(context.Context, *request.ChartDelete) (interface{}, error)
}

// HTTP API interface
type Chart struct {
	List   func(http.ResponseWriter, *http.Request)
	Create func(http.ResponseWriter, *http.Request)
	Read   func(http.ResponseWriter, *http.Request)
	Update func(http.ResponseWriter, *http.Request)
	Delete func(http.ResponseWriter, *http.Request)
}

func NewChart(h ChartAPI) *Chart {
	return &Chart{
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChartList()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Chart.List", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.List(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Chart.List", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Chart.List", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Create: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChartCreate()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Chart.Create", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Create(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Chart.Create", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Chart.Create", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Read: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChartRead()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Chart.Read", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Read(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Chart.Read", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Chart.Read", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Update: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChartUpdate()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Chart.Update", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Update(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Chart.Update", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Chart.Update", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Delete: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChartDelete()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Chart.Delete", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Delete(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Chart.Delete", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Chart.Delete", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
	}
}

func (h Chart) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/namespace/{namespaceID}/chart/", h.List)
		r.Post("/namespace/{namespaceID}/chart/", h.Create)
		r.Get("/namespace/{namespaceID}/chart/{chartID}", h.Read)
		r.Post("/namespace/{namespaceID}/chart/{chartID}", h.Update)
		r.Delete("/namespace/{namespaceID}/chart/{chartID}", h.Delete)
	})
}
