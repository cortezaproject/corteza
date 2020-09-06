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

	"github.com/cortezaproject/corteza-server/federation/rest/request"
	"github.com/cortezaproject/corteza-server/pkg/logger"
)

type (
	// Internal API interface
	IdentityAPI interface {
		GenerateNodeIdentity(context.Context, *request.IdentityGenerateNodeIdentity) (interface{}, error)
		RegisterOriginNode(context.Context, *request.IdentityRegisterOriginNode) (interface{}, error)
	}

	// HTTP API interface
	Identity struct {
		GenerateNodeIdentity func(http.ResponseWriter, *http.Request)
		RegisterOriginNode   func(http.ResponseWriter, *http.Request)
	}
)

func NewIdentity(h IdentityAPI) *Identity {
	return &Identity{
		GenerateNodeIdentity: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewIdentityGenerateNodeIdentity()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Identity.GenerateNodeIdentity", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.GenerateNodeIdentity(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Identity.GenerateNodeIdentity", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Identity.GenerateNodeIdentity", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		RegisterOriginNode: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewIdentityRegisterOriginNode()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Identity.RegisterOriginNode", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.RegisterOriginNode(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Identity.RegisterOriginNode", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Identity.RegisterOriginNode", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
	}
}

func (h Identity) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Post("/node/identity/generate", h.GenerateNodeIdentity)
		r.Post("/node/identity/register", h.RegisterOriginNode)
	})
}
