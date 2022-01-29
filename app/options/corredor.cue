package options

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

corredor: schema.#optionsGroup & {
	handle: "corredor"

	imports: [
		"\"time\"",
	]

	options: {
		enabled: {
			type:        "bool"
			default:     "false"
			description: "Enable/disable Corredor integration"
		}
		addr: {
			default:     "\"localhost:50051\""
			description: "Hostname and port of the Corredor gRPC server."
		}
		max_backoff_delay: {
			type:        "time.Duration"
			default:     "time.Minute"
			description: "Max delay for backoff on connection."
		}
		max_receive_message_size: {
			type:        "int"
			default:     "2 << 23"
			description: "Max message size that can be recived."
		}
		default_exec_timeout: {
			type:    "time.Duration"
			default: "time.Minute"
		}
		list_timeout: {
			type:    "time.Duration"
			default: "time.Second * 2"
		}
		list_refresh: {
			type:    "time.Duration"
			default: "time.Second * 5"
		}
		run_as_enabled: {
			type:    "bool"
			default: "true"
		}
		tls_cert_enabled: {
			type:    "bool"
			default: "false"
			env:     "CORREDOR_CLIENT_CERTIFICATES_ENABLED"
		}
		tls_cert_path: {
			default: "\"/certs/corredor/client\""
			env:     "CORREDOR_CLIENT_CERTIFICATES_PATH"
		}
		tls_cert_cA: {
			default: "\"ca.crt\""
			env:     "CORREDOR_CLIENT_CERTIFICATES_CA"
		}
		tls_cert_private: {
			default: "\"private.key\""
			env:     "CORREDOR_CLIENT_CERTIFICATES_PRIVATE"
		}
		tls_cert_public: {
			default: "\"public.crt\""
			env:     "CORREDOR_CLIENT_CERTIFICATES_PUBLIC"
		}
		tls_server_name: {
			env: "CORREDOR_CLIENT_CERTIFICATES_SERVER_NAME"
		}
	}
	title: "Connection to Corredor"
}
