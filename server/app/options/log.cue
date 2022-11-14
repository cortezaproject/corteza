package options

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

log: schema.#optionsGroup & {
	handle: "log"
	options: {
		debug: {
			type: "bool"
			description: """
				Disables json format for logging and enables more human-readable output with colors.

				Disable for production.

				"""
		}
		level: {
			defaultValue: "warn"
			description: """
				Minimum logging level. If set to "warn",
				Levels warn, error, dpanic panic and fatal will be logged.

				Recommended value for production: warn

				Possible values: debug, info, warn, error, dpanic, panic, fatal

				"""
		}
		filter: {
			description: """
				Log filtering rules by level and name (log-level:log-namespace).
				Please note that level (LOG_LEVEL) is applied before filter and it affects the final output!

				Leave unset for production.

				Example:
				`warn+:* *:auth,workflow.*`
				Log warnings, errors, panic, fatals. Everything from auth and workflow is logged.


				See more examples and documentation here: https://github.com/moul/zapfilter

				"""
		}
		include_caller: {
			type: "bool"
			description: """
				Set to true to see where the logging was called from.

				Disable for production.

				"""
		}
		stacktrace_level: {
			defaultValue: "dpanic"
			description: """
				Include stack-trace when logging at a specified level or below.
				Disable for production.

				Possible values: debug, info, warn, error, dpanic, panic, fatal

				"""
		}
	}
}
