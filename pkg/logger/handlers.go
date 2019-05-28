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

// LogControllerError for logging errors inside REST controllers
func LogControllerError(name string, r *http.Request, err error, params interface{}) {
	ContextValue(r.Context()).Debug(
		"error in REST controller "+name,
		zap.Error(err),
		zap.String("controller", name),
		zap.Any("params", params),
	)
}

// LogRequest logs REST controller call and HTTP request
//
// Each log entry has it's own set of fields
func LogControllerCall(name string, r *http.Request, params interface{}) {
	ContextValue(r.Context()).Debug(
		"REST controller "+name+" called",
		zap.String("controller", name),
		zap.Any("params", params),
	)
}
