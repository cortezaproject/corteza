package options

import (
	"github.com/crusttech/human/server/codegen/schema"
)

observability: schema.#optionsGroup & {
	handle: "observability"
	title:  "Agent observability"

	options: {
		langfuse_host: {
			description: "Langfuse server host. When set, agent traces are exported to Langfuse."
		}
		langfuse_public_key: {
			description: "Langfuse public key."
		}
		langfuse_secret_key: {
			description: "Langfuse secret key."
		}
		otel_exporter_otlp_endpoint: {
			description: "OpenTelemetry OTLP endpoint. When set, agent traces are exported via OTLP HTTP."
		}
	}
}
