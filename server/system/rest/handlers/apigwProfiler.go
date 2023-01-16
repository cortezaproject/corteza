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
	ApigwProfilerAPI interface {
		Aggregation(context.Context, *request.ApigwProfilerAggregation) (interface{}, error)
		Route(context.Context, *request.ApigwProfilerRoute) (interface{}, error)
		Hit(context.Context, *request.ApigwProfilerHit) (interface{}, error)
		PurgeAll(context.Context, *request.ApigwProfilerPurgeAll) (interface{}, error)
		Purge(context.Context, *request.ApigwProfilerPurge) (interface{}, error)
	}

	// HTTP API interface
	ApigwProfiler struct {
		Aggregation func(http.ResponseWriter, *http.Request)
		Route       func(http.ResponseWriter, *http.Request)
		Hit         func(http.ResponseWriter, *http.Request)
		PurgeAll    func(http.ResponseWriter, *http.Request)
		Purge       func(http.ResponseWriter, *http.Request)
	}
)

func NewApigwProfiler(h ApigwProfilerAPI) *ApigwProfiler {
	return &ApigwProfiler{
		Aggregation: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewApigwProfilerAggregation()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Aggregation(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Route: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewApigwProfilerRoute()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Route(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Hit: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewApigwProfilerHit()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Hit(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		PurgeAll: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewApigwProfilerPurgeAll()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.PurgeAll(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Purge: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewApigwProfilerPurge()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Purge(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
	}
}

func (h ApigwProfiler) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/apigw/profiler/", h.Aggregation)
		r.Get("/apigw/profiler/route/{routeID}", h.Route)
		r.Get("/apigw/profiler/hit/{hitID}", h.Hit)
		r.Post("/apigw/profiler/purge", h.PurgeAll)
		r.Post("/apigw/profiler/purge/{routeID}", h.Purge)
	})
}
