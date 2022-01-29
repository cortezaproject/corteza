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
			type:        "bool"
			default:     "true"
			description: "Enable API Gateway"
		}
		debug: {
			type:        "bool"
			default:     "false"
			description: "Enable API Gateway debugging info"
		}
		log_enabled: {
			type:        "bool"
			default:     "false"
			description: "Enable extra logging"
		}
		log_request_body: {
			type:        "bool"
			default:     "false"
			description: "Enable incoming request body output in logs"
		}
		proxy_enable_debug_log: {
			type:        "bool"
			default:     "false"
			description: "Enable full debug log on requests / responses - warning, includes sensitive data"
		}
		proxy_follow_redirects: {
			type:        "bool"
			default:     "true"
			description: "Follow redirects on proxy requests"
		}
		proxy_outbound_timeout: {
			type:        "time.Duration"
			default:     "time.Second * 30"
			description: "Outbound request timeout"
		}
	}
	title: "API Gateway"
}
