package options

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

HTTPServer: schema.#optionsGroup & {
	handle: "http-server"
	title:  "HTTP Server"

	imports: [
		"\"github.com/cortezaproject/corteza-server/pkg/rand\"",
	]

	options: {
		addr: {
			defaultValue: ":80"
			description:  "IP and port for the HTTP server."
			env:          "HTTP_ADDR"
		}
		logRequest: {
			type:        "bool"
			description: "Log HTTP requests."
			env:         "HTTP_LOG_REQUEST"
		}
		logResponse: {
			type:        "bool"
			description: "Log HTTP responses."
			env:         "HTTP_LOG_RESPONSE"
		}
		tracing: {
			type: "bool"
			env:  "HTTP_ERROR_TRACING"
		}
		enableHealthcheckRoute: {
			type:          "bool"
			defaultGoExpr: "true"
			env:           "HTTP_ENABLE_HEALTHCHECK_ROUTE"
		}
		enableVersionRoute: {
			type:          "bool"
			defaultGoExpr: "true"
			description:   "Enable `/version` route."
			env:           "HTTP_ENABLE_VERSION_ROUTE"
		}
		enableDebugRoute: {
			type:        "bool"
			description: "Enable `/debug` route."
			env:         "HTTP_ENABLE_DEBUG_ROUTE"
		}
		enableMetrics: {
			type:        "bool"
			description: "Enable (prometheus) metrics."
			env:         "HTTP_METRICS"
		}
		metricsServiceLabel: {
			defaultValue: "corteza"
			description:  "Name for metrics endpoint."
			env:          "HTTP_METRICS_NAME"
		}
		metricsUsername: {
			defaultValue: "metrics"
			description:  "Username for the metrics endpoint."
			env:          "HTTP_METRICS_USERNAME"
		}
		metricsPassword: {
			defaultGoExpr: "string(rand.Bytes(5))"
			description:   "Password for the metrics endpoint."
			env:           "HTTP_METRICS_PASSWORD"
		}
		enablePanicReporting: {
			type:          "bool"
			defaultGoExpr: "true"
			description:   "Report HTTP panic to Sentry."
			env:           "HTTP_REPORT_PANIC"
		}
		baseUrl: {
			defaultValue: "/"
			description:  "Base URL (prefix) for all routes (<baseUrl>/auth, <baseUrl>/api, ...)"
			env:          "HTTP_BASE_URL"
		}
		apiEnabled: {
			type:          "bool"
			defaultGoExpr: "true"
			env:           "HTTP_API_ENABLED"
		}
		apiBaseUrl: {
			defaultValue: "/"
			description: """
				When webapps are enabled (HTTP_WEBAPP_ENABLED) this is moved to '/api' if not explicitly set otherwise.
				API base URL is internaly prefixed with baseUrl
				"""
			env: "HTTP_API_BASE_URL"
		}
		webappEnabled: {
			type: "bool"
			env:  "HTTP_WEBAPP_ENABLED"
		}
		webappBaseUrl: {
			defaultValue: "/"
			description:  "Webapp base URL is internaly prefixed with baseUrl"
			env:          "HTTP_WEBAPP_BASE_URL"
		}
		webappBaseDir: {
			defaultValue: "./webapp/public"
			env:          "HTTP_WEBAPP_BASE_DIR"
		}
		webappList: {
			defaultValue: "admin,compose,workflow,reporter"
			env:          "HTTP_WEBAPP_LIST"
		}
		sslTerminated: {
			type:          "bool"
			defaultGoExpr: "isSecure()"
			description: """
				Is SSL termination enabled in ingres, proxy or load balancer that is in front of Corteza?
				By default, Corteza checks for presence of LETSENCRYPT_HOST environmental variable.
				This DOES NOT enable SSL termination in Cortreza!
				"""
			env: "HTTP_SSL_TERMINATED"
		}
	}
}
