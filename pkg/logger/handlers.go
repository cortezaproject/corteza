package logger

import (
	"net/http"

	"go.uber.org/zap"
)

// LogParamError for logggin invalid params
func LogParamError(name string, r *http.Request, err error, params interface{}) {
	ContextValue(r.Context()).Error(
		"invalid params for REST controller "+name,
		zap.Error(err),
		zap.String("controller", name),
		zap.Any("params", params),
	)
}

// LogControllerError for logging errors inside REST controllers
func LogControllerError(name string, r *http.Request, err error, params interface{}) {
	ContextValue(r.Context()).Error(
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
	// @todo params should provide an (auditable?) interface that would return obfuscated data
	ContextValue(r.Context()).Debug(
		"REST controller "+name+" called",
		zap.String("controller", name),
		zap.Any("params", params),
	)
}
