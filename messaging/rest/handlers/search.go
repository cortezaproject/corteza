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
	"github.com/go-chi/chi"
	"github.com/titpetric/factory/resputil"
	"net/http"

	"github.com/cortezaproject/corteza-server/messaging/rest/request"
	"github.com/cortezaproject/corteza-server/pkg/logger"
)

type (
	// Internal API interface
	SearchAPI interface {
		Messages(context.Context, *request.SearchMessages) (interface{}, error)
		Threads(context.Context, *request.SearchThreads) (interface{}, error)
	}

	// HTTP API interface
	Search struct {
		Messages func(http.ResponseWriter, *http.Request)
		Threads  func(http.ResponseWriter, *http.Request)
	}
)

func NewSearch(h SearchAPI) *Search {
	return &Search{
		Messages: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewSearchMessages()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Search.Messages", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Messages(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Search.Messages", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Search.Messages", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Threads: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewSearchThreads()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Search.Threads", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Threads(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Search.Threads", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Search.Threads", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
	}
}

func (h Search) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/search/messages", h.Messages)
		r.Get("/search/threads", h.Threads)
	})
}
