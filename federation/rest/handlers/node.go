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
	NodeAPI interface {
		Create(context.Context, *request.NodeCreate) (interface{}, error)
		Pair(context.Context, *request.NodePair) (interface{}, error)
		HandshakeConfirm(context.Context, *request.NodeHandshakeConfirm) (interface{}, error)
		HandshakeComplete(context.Context, *request.NodeHandshakeComplete) (interface{}, error)
	}

	// HTTP API interface
	Node struct {
		Create            func(http.ResponseWriter, *http.Request)
		Pair              func(http.ResponseWriter, *http.Request)
		HandshakeConfirm  func(http.ResponseWriter, *http.Request)
		HandshakeComplete func(http.ResponseWriter, *http.Request)
	}
)

func NewNode(h NodeAPI) *Node {
	return &Node{
		Create: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewNodeCreate()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Node.Create", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Create(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Node.Create", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Node.Create", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Pair: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewNodePair()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Node.Pair", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Pair(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Node.Pair", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Node.Pair", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		HandshakeConfirm: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewNodeHandshakeConfirm()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Node.HandshakeConfirm", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.HandshakeConfirm(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Node.HandshakeConfirm", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Node.HandshakeConfirm", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		HandshakeComplete: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewNodeHandshakeComplete()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Node.HandshakeComplete", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.HandshakeComplete(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Node.HandshakeComplete", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Node.HandshakeComplete", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
	}
}

func (h Node) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Post("/nodes", h.Create)
		r.Post("/nodes/{nodeID}/pair", h.Pair)
		r.Post("/nodes/{nodeID}/handshake-confirm", h.HandshakeConfirm)
		r.Post("/nodes/{nodeID}/handshake-complete", h.HandshakeComplete)
	})
}
