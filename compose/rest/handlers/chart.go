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

	"github.com/crusttech/crust/compose/rest/request"
	"github.com/crusttech/crust/internal/logger"
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
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Chart.List", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			if value, err := ch.List(r.Context(), params); err != nil {
				logger.LogControllerError("Chart.List", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			} else {
				logger.LogControllerCall("Chart.List", r, params.Auditable())
				switch fn := value.(type) {
				case func(http.ResponseWriter, *http.Request):
					fn(w, r)
					return
				}
				resputil.JSON(w, value)
				return
			}
		},
		Create: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChartCreate()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Chart.Create", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			if value, err := ch.Create(r.Context(), params); err != nil {
				logger.LogControllerError("Chart.Create", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			} else {
				logger.LogControllerCall("Chart.Create", r, params.Auditable())
				switch fn := value.(type) {
				case func(http.ResponseWriter, *http.Request):
					fn(w, r)
					return
				}
				resputil.JSON(w, value)
				return
			}
		},
		Read: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChartRead()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Chart.Read", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			if value, err := ch.Read(r.Context(), params); err != nil {
				logger.LogControllerError("Chart.Read", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			} else {
				logger.LogControllerCall("Chart.Read", r, params.Auditable())
				switch fn := value.(type) {
				case func(http.ResponseWriter, *http.Request):
					fn(w, r)
					return
				}
				resputil.JSON(w, value)
				return
			}
		},
		Update: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChartUpdate()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Chart.Update", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			if value, err := ch.Update(r.Context(), params); err != nil {
				logger.LogControllerError("Chart.Update", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			} else {
				logger.LogControllerCall("Chart.Update", r, params.Auditable())
				switch fn := value.(type) {
				case func(http.ResponseWriter, *http.Request):
					fn(w, r)
					return
				}
				resputil.JSON(w, value)
				return
			}
		},
		Delete: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewChartDelete()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Chart.Delete", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			if value, err := ch.Delete(r.Context(), params); err != nil {
				logger.LogControllerError("Chart.Delete", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			} else {
				logger.LogControllerCall("Chart.Delete", r, params.Auditable())
				switch fn := value.(type) {
				case func(http.ResponseWriter, *http.Request):
					fn(w, r)
					return
				}
				resputil.JSON(w, value)
				return
			}
		},
	}
}

func (ch *Chart) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/namespace/{namespaceID}/chart/", ch.List)
		r.Post("/namespace/{namespaceID}/chart/", ch.Create)
		r.Get("/namespace/{namespaceID}/chart/{chartID}", ch.Read)
		r.Post("/namespace/{namespaceID}/chart/{chartID}", ch.Update)
		r.Delete("/namespace/{namespaceID}/chart/{chartID}", ch.Delete)
	})
}
