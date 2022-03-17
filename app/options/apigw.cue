package options

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

apigw: schema.#optionsGroup & {
	handle: "apigw"

	imports: [
		"\"time\"",
	]

	options: {
		enabled: {
			type:          "bool"
			defaultGoExpr: "true"
			description:   "Enable API Gateway"
		}
		debug: {
			type:        "bool"
			description: "Enable API Gateway debugging info"
		}
		log_enabled: {
			type:        "bool"
			description: "Enable extra logging"
		}
		profiler_enabled: {
			type:        "bool"
			description: "Enable profiler"
		}
		profiler_global: {
			type:        "bool"
			description: "Profiler enabled for all routes"
		}
		log_request_body: {
			type:        "bool"
			description: "Enable incoming request body output in logs"
		}
		proxy_enable_debug_log: {
			type:        "bool"
			description: "Enable full debug log on requests / responses - warning, includes sensitive data"
		}
		proxy_follow_redirects: {
			type:          "bool"
			defaultGoExpr: "true"
			description:   "Follow redirects on proxy requests"
		}
		proxy_outbound_timeout: {
			type:          "time.Duration"
			description:   "Outbound request timeout"
			defaultGoExpr: "time.Second * 30"
			defaultValue:  "30s"
		}
	}
	title: "API Gateway"
}
