package options

import (
	"github.com/cortezaproject/corteza/server/codegen/schema"
)

HTTPClient: schema.#optionsGroup & {
	title: "HTTP Client"
	// Explicitly define all variants to be 100% compaltible with old name
	handle: "http-client"

	// @todo remove explcitly defined expIdent and adjust the code
	expIdent: "HTTPClient"

	imports: [
		"\"time\"",
	]

	options: {
		tls_insecure: {
			type: "bool"
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
			description: "Default timeout for clients."

			defaultGoExpr: "30 * time.Second"
			defaultValue:  "30s"
		}
		restricted_networks: {
			defaultValue: "link-local"
			description: """
				Networks that can not be reached by HTTP requests to destinations controlled by users
				(workflow HTTP request and OAuth2 functions, API gateway proxy).

				Supported values:
				 - link-local: blocks link-local addresses and well known cloud metadata endpoints (default)
				 - private: additionally blocks loopback, private and other non-public addresses
				 - none: requests can reach any address

				[NOTE]
				====
				Restrictions are checked against the resolved IP address when connecting.
				When outbound HTTP proxy is used, restrictions need to be enforced on the proxy.
				====
				"""
			env: "HTTP_CLIENT_RESTRICTED_NETWORKS"
		}
	}
}
