package flags

import (
	"github.com/spf13/cobra"

	"github.com/cortezaproject/corteza-server/internal/rand"
)

type (
	HTTPOpt struct {
		Addr    string
		Logging bool
		Pretty  bool
		Tracing bool

		EnableVersionRoute bool
		EnableDebugRoute   bool

		EnableMetrics       bool
		MetricsServiceLabel string
		MetricsUsername     string
		MetricsPassword     string
	}
)

func HTTP(cmd *cobra.Command, pfix string) (o *HTTPOpt) {
	o = &HTTPOpt{}

	bindString(cmd, &o.Addr,
		pFlag(pfix, "http-addr"), ":80",
		"Listen address for HTTP server")

	bindBool(cmd, &o.Logging,
		pFlag(pfix, "http-log"), true,
		"Enable/disable HTTP request log")

	bindBool(cmd, &o.Pretty,
		pFlag(pfix, "http-pretty-json"), false,
		"Prettify returned JSON output")

	bindBool(cmd, &o.Tracing,
		pFlag(pfix, "http-error-tracing"), false,
		"Return error stack frame")

	bindBool(cmd, &o.EnableVersionRoute,
		pFlag(pfix, "http-enable-version-route"), true,
		"Enable /version route")

	bindBool(cmd, &o.EnableDebugRoute,
		pFlag(pfix, "http-enable-debug-route"), false,
		"Enable /debug route with pprof data")

	bindBool(cmd, &o.EnableMetrics,
		pFlag(pfix, "http-metrics"), false,
		"Enable metrics")

	bindString(cmd, &o.MetricsServiceLabel,
		pFlag(pfix, "http-metrics-name"), "corteza",
		"Provide metrics service label for Prometheus")

	bindString(cmd, &o.MetricsUsername,
		pFlag(pfix, "http-metrics-username"), "metrics",
		"Provide metrics username for Prometheus")

	// Setting metrics password to random string to prevent security accidents...
	bindString(cmd, &o.MetricsPassword,
		pFlag(pfix, "http-metrics-password"), string(rand.Bytes(5)),
		"Provide metrics password for Prometheus")

	return
}
