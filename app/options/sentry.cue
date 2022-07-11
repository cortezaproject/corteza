package options

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

sentry: schema.#optionsGroup & {
	handle: "sentry"

	imports: [
		"\"github.com/cortezaproject/corteza-server/pkg/version\"",
	]

	title: "Sentry monitoring"
	intro: """
		[NOTE]
		====
		These parameters help in the development and testing process.
		When you are deploying to production, these should be disabled to improve performance and reduce storage usage.

		You should configure external services such as Sentry or ELK to keep track of logs and error reports.
		====
		"""

	options: {
		DSN: {
			description: "Set to enable Sentry client."
		}
		debug: {
			type:        "bool"
			description: "Print out debugging information."
		}
		attach_stacktrace: {
			type:          "bool"
			defaultGoExpr: "true"
			description:   "Attach stacktraces"
		}
		sample_rate: {
			type:        "float64"
			description: "Sample rate for event submission (0.0 - 1.0. defaults to 1.0)"
		}
		max_breadcrumbs: {
			type:          "int"
			defaultGoExpr: "0"
			description:   "Maximum number of bredcrumbs."
		}
		server_name: {
			description: "Set reported Server name."
			env:         "SENTRY_SERVERNAME"
		}
		release: {
			defaultGoExpr: "version.Version"
			description:   "Set reported Release."
		}
		dist: {
			description: "Set reported distribution."
		}
		environment: {
			description: "Set reported environment."
		}
		webapp_DSN: {
			description: "Set to enable Sentry client for webapp."
		}
	}
}
