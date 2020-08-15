package options

import (
	"github.com/cortezaproject/corteza-server/pkg/rand"
)

type (
	HTTPServerOpt struct {
		Addr        string `env:"HTTP_ADDR"`
		LogRequest  bool   `env:"HTTP_LOG_REQUEST"`
		LogResponse bool   `env:"HTTP_LOG_RESPONSE"`
		Tracing     bool   `env:"HTTP_ERROR_TRACING"`

		EnableHealthcheckRoute bool `env:"HTTP_ENABLE_HEALTHCHECK_ROUTE"`
		EnableVersionRoute     bool `env:"HTTP_ENABLE_VERSION_ROUTE"`
		EnableDebugRoute       bool `env:"HTTP_ENABLE_DEBUG_ROUTE"`

		EnableMetrics       bool   `env:"HTTP_METRICS"`
		MetricsServiceLabel string `env:"HTTP_METRICS_NAME"`
		MetricsUsername     string `env:"HTTP_METRICS_USERNAME"`
		MetricsPassword     string `env:"HTTP_METRICS_PASSWORD"`

		EnablePanicReporting bool `env:"HTTP_REPORT_PANIC"`

		ApiEnabled bool   `env:"HTTP_API_ENABLED"`
		ApiBaseUrl string `env:"HTTP_API_BASE_URL"`

		WebappEnabled bool   `env:"HTTP_WEBAPP_ENABLED"`
		WebappBaseUrl string `env:"HTTP_WEBAPP_BASE_URL"`
		WebappBaseDir string `env:"HTTP_WEBAPP_BASE_DIR"`
		WebappList    string `env:"HTTP_WEBAPP_LIST"`
	}
)

func HTTP(pfix string) (o *HTTPServerOpt) {
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

		// Reports panics to Sentry through HTTP middleware
		EnablePanicReporting: true,

		// Setting metrics password to random string to prevent security accidents...
		MetricsPassword: string(rand.Bytes(5)),

		ApiEnabled: true,
		ApiBaseUrl: "",

		WebappEnabled: false,
		WebappBaseUrl: "/",
		WebappBaseDir: "/webapp",
		WebappList:    "admin,auth,messaging,compose",
	}

	fill(o)

	return
}
