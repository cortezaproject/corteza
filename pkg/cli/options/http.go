package options

import (
	"github.com/cortezaproject/corteza-server/internal/rand"
)

type (
	HTTPOpt struct {
		Addr    string `env:"HTTP_ADDR"`
		Logging bool   `env:"HTTP_LOG_REQUESTS"`
		Tracing bool   `env:"HTTP_ERROR_TRACING"`

		EnableVersionRoute bool `env:"HTTP_ENABLE_VERSION_ROUTE"`
		EnableDebugRoute   bool `env:"HTTP_ENABLE_DEBUG_ROUTE"`

		EnableMetrics       bool   `env:"HTTP_METRICS"`
		MetricsServiceLabel string `env:"HTTP_METRICS_NAME"`
		MetricsUsername     string `env:"HTTP_METRICS_USERNAME"`
		MetricsPassword     string `env:"HTTP_METRICS_PASSWORD"`
	}
)

func HTTP(pfix string) (o *HTTPOpt) {
	o = &HTTPOpt{
		Addr:                ":80",
		Logging:             true,
		Tracing:             false,
		EnableVersionRoute:  true,
		EnableDebugRoute:    false,
		EnableMetrics:       false,
		MetricsServiceLabel: "corteza",
		MetricsUsername:     "metrics",

		// Setting metrics password to random string to prevent security accidents...
		MetricsPassword: string(rand.Bytes(5)),
	}

	fill(o, pfix)

	return
}
