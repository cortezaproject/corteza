package options

import (
	"github.com/crusttech/human/server/codegen/schema"
)

corredor: schema.#optionsGroup & {
	handle: "corredor"

	imports: [
		"\"time\"",
	]

	options: {
		enabled: {
			type:        "bool"
			description: "Enable/disable Corredor integration"
		}
		addr: {
			defaultValue: "localhost:50051"
			description:  "Hostname and port of the Corredor gRPC server."
		}
		max_backoff_delay: {
			type:          "time.Duration"
			description:   "Max delay for backoff on connection."
			defaultGoExpr: "time.Minute"
			defaultValue:  "1m"
		}
		max_receive_message_size: {
			type:          "int"
			defaultGoExpr: "2 << 23"
			defaultValue:  "16777216"
			description:   "Max message size that can be received."
		}
		default_exec_timeout: {
			type:          "time.Duration"
			defaultGoExpr: "time.Minute"
			defaultValue:  "1m"
		}
		list_timeout: {
			type:          "time.Duration"
			defaultGoExpr: "time.Second * 2"
			defaultValue:  "2s"
		}
		list_refresh: {
			type:          "time.Duration"
			defaultGoExpr: "time.Second * 5"
			defaultValue:  "5s"
		}
		run_as_enabled: {
			type:          "bool"
			defaultGoExpr: "true"
			defaultValue:  "true"
		}
		tls_cert_enabled: {
			type: "bool"
			env:  "CORREDOR_CLIENT_CERTIFICATES_ENABLED"
		}
		tls_cert_path: {
			defaultValue: "/certs/corredor/client"
			env:          "CORREDOR_CLIENT_CERTIFICATES_PATH"
		}
		tls_cert_cA: {
			defaultValue: "ca.crt"
			env:          "CORREDOR_CLIENT_CERTIFICATES_CA"
		}
		tls_cert_private: {
			defaultValue: "private.key"
			env:          "CORREDOR_CLIENT_CERTIFICATES_PRIVATE"
		}
		tls_cert_public: {
			defaultValue: "public.crt"
			env:          "CORREDOR_CLIENT_CERTIFICATES_PUBLIC"
		}
		tls_server_name: {
			env: "CORREDOR_CLIENT_CERTIFICATES_SERVER_NAME"
		}
	}
	title: "Connection to Corredor"
}
