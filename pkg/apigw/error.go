package apigw

import (
	"context"
	"net/http"
)

type (
	defaultErrorHandler struct{}
)

func (h defaultErrorHandler) Exec(ctx context.Context, scope *scp, err error) {
	// set http status code
	scope.Writer().WriteHeader(http.StatusInternalServerError)

	// set body
	scope.Writer().Write([]byte(err.Error()))
}
