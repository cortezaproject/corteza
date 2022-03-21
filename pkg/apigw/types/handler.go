package types

import (
	"fmt"
	"net/http"

	"github.com/cortezaproject/corteza-server/pkg/options"
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

		New(options.ApigwOpt) Handler
		Merge([]byte) (Handler, error)
		Meta() FilterMeta
		Enabled() bool
	}

	HandlerFunc      func(rw http.ResponseWriter, r *http.Request) error
	ErrorHandlerFunc func(rw http.ResponseWriter, r *http.Request, err error)
)
