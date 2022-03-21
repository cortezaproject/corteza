package chain

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type (
	middleware []func(http.Handler) http.Handler

	ChainHandler interface {
		Handler() http.Handler
		Chain(middleware)
	}

	chiHandlerChainDefault struct {
		chain middleware
	}
)

func NewDefault() *chiHandlerChainDefault {
	return &chiHandlerChainDefault{}
}

func (cc *chiHandlerChainDefault) Chain(m middleware) {
	cc.chain = m
}

func (cc *chiHandlerChainDefault) Handler() http.Handler {
	return chiHandler(cc.chain)
}

func chiHandler(cc middleware) http.Handler {
	return chi.
		Chain(cc...).
		HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {})
}
