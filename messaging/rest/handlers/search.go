package handlers

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `search.go`, `search.util.go` or `search_test.go` to
	implement your API calls, helper functions and tests. The file `search.go`
	is only generated the first time, and will not be overwritten if it exists.
*/

import (
	"context"

	"net/http"

	"github.com/go-chi/chi"
	"github.com/titpetric/factory/resputil"

	"github.com/cortezaproject/corteza-server/messaging/rest/request"
	"github.com/cortezaproject/corteza-server/pkg/logger"
)

// Internal API interface
type SearchAPI interface {
	Messages(context.Context, *request.SearchMessages) (interface{}, error)
	Threads(context.Context, *request.SearchThreads) (interface{}, error)
}

// HTTP API interface
type Search struct {
	Messages func(http.ResponseWriter, *http.Request)
	Threads  func(http.ResponseWriter, *http.Request)
}

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
