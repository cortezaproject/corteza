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
		Search(context.Context, *request.NodeSearch) (interface{}, error)
		Create(context.Context, *request.NodeCreate) (interface{}, error)
		GenerateURI(context.Context, *request.NodeGenerateURI) (interface{}, error)
		Update(context.Context, *request.NodeUpdate) (interface{}, error)
		Delete(context.Context, *request.NodeDelete) (interface{}, error)
		Undelete(context.Context, *request.NodeUndelete) (interface{}, error)
		Pair(context.Context, *request.NodePair) (interface{}, error)
		HandshakeConfirm(context.Context, *request.NodeHandshakeConfirm) (interface{}, error)
		HandshakeComplete(context.Context, *request.NodeHandshakeComplete) (interface{}, error)
	}

	// HTTP API interface
	Node struct {
		Search            func(http.ResponseWriter, *http.Request)
		Create            func(http.ResponseWriter, *http.Request)
		GenerateURI       func(http.ResponseWriter, *http.Request)
		Update            func(http.ResponseWriter, *http.Request)
		Delete            func(http.ResponseWriter, *http.Request)
		Undelete          func(http.ResponseWriter, *http.Request)
		Pair              func(http.ResponseWriter, *http.Request)
		HandshakeConfirm  func(http.ResponseWriter, *http.Request)
		HandshakeComplete func(http.ResponseWriter, *http.Request)
	}
)

func NewNode(h NodeAPI) *Node {
	return &Node{
		Search: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewNodeSearch()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Node.Search", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Search(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Node.Search", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Node.Search", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
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
		GenerateURI: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewNodeGenerateURI()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Node.GenerateURI", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.GenerateURI(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Node.GenerateURI", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Node.GenerateURI", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Update: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewNodeUpdate()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Node.Update", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Update(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Node.Update", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Node.Update", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Delete: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewNodeDelete()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Node.Delete", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Delete(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Node.Delete", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Node.Delete", r, params.Auditable())
			if !serveHTTP(value, w, r) {
				resputil.JSON(w, value)
			}
		},
		Undelete: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewNodeUndelete()
			if err := params.Fill(r); err != nil {
				logger.LogParamError("Node.Undelete", r, err)
				resputil.JSON(w, err)
				return
			}

			value, err := h.Undelete(r.Context(), params)
			if err != nil {
				logger.LogControllerError("Node.Undelete", r, err, params.Auditable())
				resputil.JSON(w, err)
				return
			}
			logger.LogControllerCall("Node.Undelete", r, params.Auditable())
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
		r.Get("/nodes", h.Search)
		r.Post("/nodes", h.Create)
		r.Post("/nodes/{nodeID}/uri", h.GenerateURI)
		r.Post("/nodes/{nodeID}", h.Update)
		r.Delete("/nodes/{nodeID}", h.Delete)
		r.Post("/nodes/{nodeID}/undelete", h.Undelete)
		r.Post("/nodes/{nodeID}/pair", h.Pair)
		r.Post("/nodes/{nodeID}/handshake-confirm", h.HandshakeConfirm)
		r.Post("/nodes/{nodeID}/handshake-complete", h.HandshakeComplete)
	})
}
