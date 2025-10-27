package handlers

import (
	"context"
	"net/http"

	"github.com/cortezaproject/corteza/server/pkg/api"
	"github.com/go-chi/chi/v5"
)

type (
	RagAPI interface {
		Query(context.Context, string, int) (string, error)
	}

	Rag struct {
		Query func(http.ResponseWriter, *http.Request)
	}
)

func NewRagQuery(h RagAPI) *Rag {
	return &Rag{
		Query: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()

			// GET params
			params := r.URL.Query()
			question := params.Get("question")

			value, err := h.Query(r.Context(), question, 10)
			if err != nil {
				api.Send(w, r, err)
				return
			}

			api.Send(w, r, value)
		},
	}
}

func (h Rag) MountRoutes(r chi.Router, middlewares ...func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares...)
		r.Get("/rag", h.Query)
	})
}
