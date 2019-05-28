package handlers

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `commands.go`, `commands.util.go` or `commands_test.go` to
	implement your API calls, helper functions and tests. The file `commands.go`
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
type CommandsAPI interface {
	List(context.Context, *request.CommandsList) (interface{}, error)
}

// HTTP API interface
type Commands struct {
	List func(http.ResponseWriter, *http.Request)
}

func NewCommands(h CommandsAPI) *Commands {
	return &Commands{
		List: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewCommandsList()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Commands.List", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.List(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Commands.List", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Commands.List", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
	}
}

func (h Commands) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/commands/", h.List)
	})
}
