package types

import (
	"fmt"
	"net/http"
)

type (
	HTTPHandler interface {
		Handler() HandlerFunc
	}

	HTTPErrorHandler interface {
		Handler() ErrorHandlerFunc
	}

	Handler interface {
		HTTPHandler
		fmt.Stringer

		New(Config) Handler
		Merge([]byte, Config) (Handler, error)
		Meta() FilterMeta
		Enabled() bool
	}

	HandlerFunc      func(rw http.ResponseWriter, r *http.Request) error
	ErrorHandlerFunc func(rw http.ResponseWriter, r *http.Request, err error)
)
