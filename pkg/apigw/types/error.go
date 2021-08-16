package types

import (
	"context"
	"net/http"
)

type (
	DefaultErrorHandler struct{}
)

func (h DefaultErrorHandler) Exec(ctx context.Context, scope *Scp, err error) {
	// set http status code
	scope.Writer().WriteHeader(http.StatusInternalServerError)

	// set body
	scope.Writer().Write([]byte(err.Error()))
}
