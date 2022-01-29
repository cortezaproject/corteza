package options

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

HTTPClient: schema.#optionsGroup & {
	title:    "HTTP Client"
	// Explicitly define all variants to be 100% compaltible with old name
	handle:   "http-client"

	// @todo remove explcitly defined expIdent and adjust the code
	expIdent: "HTTPClient"

	imports: [
		"\"time\"",
	]



	options: {
		tls_insecure: {
			type:    "bool"
			default: "false"
			description: """
				Allow insecure (invalid, expired TLS/SSL certificates) connections.

				[IMPORTANT]
				====
				We strongly recommend keeping this value set to false except for local development or demos.
				====
				"""
		}
		timeout: {
			type:        "time.Duration"
			default:     "30 * time.Second"
			description: "Default timeout for clients."
		}
	}
}
