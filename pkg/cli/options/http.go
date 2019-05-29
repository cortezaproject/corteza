package options

import (
	"github.com/cortezaproject/corteza-server/internal/rand"
)

type (
	HTTPOpt struct {
		Addr    string
		Logging bool
		Tracing bool

		EnableVersionRoute bool
		EnableDebugRoute   bool

		EnableMetrics       bool
		MetricsServiceLabel string
		MetricsUsername     string
		MetricsPassword     string
	}
)

func HTTP(pfix string) (o *HTTPOpt) {
	o = &HTTPOpt{
		Addr:                EnvString(pfix, "HTTP_ADDR", ":80"),
		Logging:             EnvBool(pfix, "HTTP_LOG_REQUESTS", true),
		Tracing:             EnvBool(pfix, "HTTP_ERROR_TRACING", false),
		EnableVersionRoute:  EnvBool(pfix, "HTTP_ENABLE_VERSION_ROUTE", true),
		EnableDebugRoute:    EnvBool(pfix, "HTTP_ENABLE_DEBUG_ROUTE", false),
		EnableMetrics:       EnvBool(pfix, "HTTP_METRICS", false),
		MetricsServiceLabel: EnvString(pfix, "HTTP_METRICS_NAME", "corteza"),
		MetricsUsername:     EnvString(pfix, "HTTP_METRICS_USERNAME", "metrics"),

		// Setting metrics password to random string to prevent security accidents...
		MetricsPassword: EnvString(pfix, "HTTP_METRICS_PASSWORD", string(rand.Bytes(5))),
	}

	return
}
