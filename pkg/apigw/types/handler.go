package types

import (
	"net/http"
)

type (
	Stringer interface {
		String() string
	}

	HTTPHandler interface {
		Handler() HandlerFunc
	}

	HTTPErrorHandler interface {
		Handler() ErrorHandlerFunc
	}

	Handler interface {
		HTTPHandler
		Stringer

		Merge([]byte) (Handler, error)
		Meta() FilterMeta
	}

	HandlerFunc      func(rw http.ResponseWriter, r *http.Request) error
	ErrorHandlerFunc func(rw http.ResponseWriter, r *http.Request, err error)
)
