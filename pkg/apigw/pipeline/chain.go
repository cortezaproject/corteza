package pipeline

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type (
	chiHandlerChain struct {
		chain []func(http.Handler) http.Handler
	}
)

func (cc chiHandlerChain) Handler() (h http.Handler) {
	return chi.
		Chain(cc.chain...).
		HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {})
}
