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
	"github.com/cortezaproject/corteza/server/federation/rest/request"
	"github.com/cortezaproject/corteza/server/pkg/api"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type (
	// Internal API interface
	NodeAPI interface {
		Search(context.Context, *request.NodeSearch) (interface{}, error)
		Create(context.Context, *request.NodeCreate) (interface{}, error)
		Read(context.Context, *request.NodeRead) (interface{}, error)
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
		Read              func(http.ResponseWriter, *http.Request)
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
				api.Send(w, r, err)
				return
			}

			value, err := h.Search(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Create: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewNodeCreate()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Create(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Read: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewNodeRead()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Read(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		GenerateURI: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewNodeGenerateURI()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.GenerateURI(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Update: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewNodeUpdate()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Update(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Delete: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewNodeDelete()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Delete(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Undelete: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewNodeUndelete()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Undelete(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		Pair: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewNodePair()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.Pair(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		HandshakeConfirm: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewNodeHandshakeConfirm()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.HandshakeConfirm(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
		HandshakeComplete: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			params := request.NewNodeHandshakeComplete()
			if err := params.Fill(r); err != nil {
				api.Send(w, r, err)
				return
			}

			value, err := h.HandshakeComplete(r.Context(), params)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
	}
}

func (h Node) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/nodes/", h.Search)
		r.Post("/nodes/", h.Create)
		r.Get("/nodes/{nodeID}", h.Read)
		r.Post("/nodes/{nodeID}/uri", h.GenerateURI)
		r.Post("/nodes/{nodeID}", h.Update)
		r.Delete("/nodes/{nodeID}", h.Delete)
		r.Post("/nodes/{nodeID}/undelete", h.Undelete)
		r.Post("/nodes/{nodeID}/pair", h.Pair)
		r.Post("/nodes/{nodeID}/handshake-confirm", h.HandshakeConfirm)
		r.Post("/nodes/{nodeID}/handshake-complete", h.HandshakeComplete)
	})
}
