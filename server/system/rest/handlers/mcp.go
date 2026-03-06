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
	McpAPI interface {
		ListTools(context.Context, *request.McpListTools) (interface{}, error)
	}

	// HTTP API interface
	Mcp struct {
		ListTools func(http.ResponseWriter, *http.Request)
	}
)

func NewMcp(h McpAPI) *Mcp {
	return &Mcp{
		ListTools: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewMcpListTools()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.ListTools(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
	}
}

func (h Mcp) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/mcp/tools", h.ListTools)
	})
}
