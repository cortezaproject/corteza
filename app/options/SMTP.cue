package options

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

SMTP: schema.#optionsGroup & {
	handle: "smtp"
	title:  "Email sending"

	// @todo remove explicitly defined expIdent and adjust the code
	expIdent: "SMTP"

	intro: """
		Configure your local SMTP server or use one of the available providers.

		These values are copied to settings when the server starts and can be managed from the administration console.
		We recommend you remove these values after they are copied to settings.
		If server detects difference between these options and settings, it shows a warning in the log on server start.
		"""

	options: {
		host: {
			defaultValue: "localhost"
			description:  "The SMTP server hostname."
		}
		port: {
			type:          "int"
			defaultGoExpr: "25"
			description:   "The SMTP post."
		}
		user: {
			description: "The SMTP username."
		}
		pass: {
			description: "The SMTP password."
		}
		from: {
			description: "The SMTP `from` email parameter"
		}
		tls_insecure: {
			type:        "bool"
			description: "Allow insecure (invalid, expired TLS certificates) connections."
		}
		tls_server_name: {}
	}
}
