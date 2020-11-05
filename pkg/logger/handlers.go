package logger

import (
	"net/http"

	"go.uber.org/zap"
)

// LogParamError for loggin invalid params
func LogParamError(name string, r *http.Request, err error) {
	ContextValue(r.Context()).Debug(
		"invalid params for REST controller "+name,
		zap.Error(err),
		zap.String("controller", name),
	)
}
