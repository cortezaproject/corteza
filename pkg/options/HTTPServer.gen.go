package options

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// pkg/options/HTTPServer.yaml

import (
	"github.com/cortezaproject/corteza-server/pkg/rand"
)

type (
	HTTPServerOpt struct {
		Addr                   string `env:"HTTP_ADDR"`
		LogRequest             bool   `env:"HTTP_LOG_REQUEST"`
		LogResponse            bool   `env:"HTTP_LOG_RESPONSE"`
		Tracing                bool   `env:"HTTP_ERROR_TRACING"`
		EnableHealthcheckRoute bool   `env:"HTTP_ENABLE_HEALTHCHECK_ROUTE"`
		EnableVersionRoute     bool   `env:"HTTP_ENABLE_VERSION_ROUTE"`
		EnableDebugRoute       bool   `env:"HTTP_ENABLE_DEBUG_ROUTE"`
		EnableMetrics          bool   `env:"HTTP_METRICS"`
		MetricsServiceLabel    string `env:"HTTP_METRICS_NAME"`
		MetricsUsername        string `env:"HTTP_METRICS_USERNAME"`
		MetricsPassword        string `env:"HTTP_METRICS_PASSWORD"`
		EnablePanicReporting   bool   `env:"HTTP_REPORT_PANIC"`
		ApiEnabled             bool   `env:"HTTP_API_ENABLED"`
		ApiBaseUrl             string `env:"HTTP_API_BASE_URL"`
		WebappEnabled          bool   `env:"HTTP_WEBAPP_ENABLED"`
		WebappBaseUrl          string `env:"HTTP_WEBAPP_BASE_URL"`
		WebappBaseDir          string `env:"HTTP_WEBAPP_BASE_DIR"`
		WebappList             string `env:"HTTP_WEBAPP_LIST"`
	}
)

// HTTPServer initializes and returns a HTTPServerOpt with default values
func HTTPServer() (o *HTTPServerOpt) {
	o = &HTTPServerOpt{
		Addr:                   ":80",
		LogRequest:             false,
		LogResponse:            false,
		Tracing:                false,
		EnableHealthcheckRoute: true,
		EnableVersionRoute:     true,
		EnableDebugRoute:       false,
		EnableMetrics:          false,
		MetricsServiceLabel:    "corteza",
		MetricsUsername:        "metrics",
		MetricsPassword:        string(rand.Bytes(5)),
		EnablePanicReporting:   true,
		ApiEnabled:             true,
		WebappEnabled:          false,
		WebappBaseUrl:          "/",
		WebappBaseDir:          "webapp/public",
		WebappList:             "admin,compose,workflow",
	}

	fill(o)

	// Function that allows access to custom logic inside the parent function.
	// The custom logic in the other file should be like:
	// func (o *HTTPServer) Defaults() {...}
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	return
}
