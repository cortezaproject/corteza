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
	"github.com/go-chi/chi"
	"net/http"

	"github.com/titpetric/factory/resputil"

	"github.com/crusttech/crust/crm/rest/request"
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

func NewChart(ch ChartAPI) *Chart {
	return &Chart{
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChartList()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return ch.List(r.Context(), params)
			})
		},
		Create: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChartCreate()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return ch.Create(r.Context(), params)
			})
		},
		Read: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChartRead()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return ch.Read(r.Context(), params)
			})
		},
		Update: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChartUpdate()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return ch.Update(r.Context(), params)
			})
		},
		Delete: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChartDelete()
			resputil.JSON(w, params.Fill(r), func() (interface{}, error) {
				return ch.Delete(r.Context(), params)
			})
		},
	}
}

func (ch *Chart) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Route("/chart", func(r chi.Router) {
			r.Get("/", ch.List)
			r.Post("/", ch.Create)
			r.Get("/{chartID}", ch.Read)
			r.Post("/{chartID}", ch.Update)
			r.Delete("/{chartID}", ch.Delete)
		})
	})
}
