package options

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

HTTPServer: schema.#optionsGroup & {
	handle: "HTTPServer"
	title: "HTTP Server"

	imports: [
		"\"github.com/cortezaproject/corteza-server/pkg/rand\"",
	]

	options: {
		addr: {
			default:     "\":80\""
			description: "IP and port for the HTTP server."
			env:         "HTTP_ADDR"
		}
		logRequest: {
			type:        "bool"
			default:     "false"
			description: "Log HTTP requests."
			env:         "HTTP_LOG_REQUEST"
		}
		logResponse: {
			type:        "bool"
			default:     "false"
			description: "Log HTTP responses."
			env:         "HTTP_LOG_RESPONSE"
		}
		tracing: {
			type:    "bool"
			default: "false"
			env:     "HTTP_ERROR_TRACING"
		}
		enableHealthcheckRoute: {
			type:    "bool"
			default: "true"
			env:     "HTTP_ENABLE_HEALTHCHECK_ROUTE"
		}
		enableVersionRoute: {
			type:        "bool"
			default:     "true"
			description: "Enable `/version` route."
			env:         "HTTP_ENABLE_VERSION_ROUTE"
		}
		enableDebugRoute: {
			type:        "bool"
			default:     "false"
			description: "Enable `/debug` route."
			env:         "HTTP_ENABLE_DEBUG_ROUTE"
		}
		enableMetrics: {
			type:        "bool"
			default:     "false"
			description: "Enable (prometheus) metrics."
			env:         "HTTP_METRICS"
		}
		metricsServiceLabel: {
			default:     "\"corteza\""
			description: "Name for metrics endpoint."
			env:         "HTTP_METRICS_NAME"
		}
		metricsUsername: {
			default:     "\"metrics\""
			description: "Username for the metrics endpoint."
			env:         "HTTP_METRICS_USERNAME"
		}
		metricsPassword: {
			default:     "string(rand.Bytes(5))"
			description: "Password for the metrics endpoint."
			env:         "HTTP_METRICS_PASSWORD"
		}
		enablePanicReporting: {
			type:        "bool"
			default:     "true"
			description: "Report HTTP panic to Sentry."
			env:         "HTTP_REPORT_PANIC"
		}
		baseUrl: {
			default:     "\"/\""
			description: "Base URL (prefix) for all routes (<baseUrl>/auth, <baseUrl>/api, ...)"
			env:         "HTTP_BASE_URL"
		}
		apiEnabled: {
			type:    "bool"
			default: "true"
			env:     "HTTP_API_ENABLED"
		}
		apiBaseUrl: {
			default: "\"/\""
			description: """
				When webapps are enabled (HTTP_WEBAPP_ENABLED) this is moved to '/api' if not explicitly set otherwise.
				API base URL is internaly prefixed with baseUrl
				"""
			env: "HTTP_API_BASE_URL"
		}
		webappEnabled: {
			type:    "bool"
			default: "false"
			env:     "HTTP_WEBAPP_ENABLED"
		}
		webappBaseUrl: {
			default:     "\"/\""
			description: "Webapp base URL is internaly prefixed with baseUrl"
			env:         "HTTP_WEBAPP_BASE_URL"
		}
		webappBaseDir: {
			default: "\"./webapp/public\""
			env:     "HTTP_WEBAPP_BASE_DIR"
		}
		webappList: {
			default: "\"admin,compose,workflow,reporter\""
			env:     "HTTP_WEBAPP_LIST"
		}
		sslTerminated: {
			type:    "bool"
			default: "isSecure()"
			description: """
				Is SSL termination enabled in ingres, proxy or load balancer that is in front of Corteza?
				By default, Corteza checks for presence of LETSENCRYPT_HOST environmental variable.
				This DOES NOT enable SSL termination in Cortreza!
				"""
			env: "HTTP_SSL_TERMINATED"
		}
	}
}
