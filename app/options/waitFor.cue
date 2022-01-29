package options

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

waitFor: schema.#optionsGroup & {
	handle: "waitFor"

	imports: [
		"\"time\"",
	]
	title: "Delay system startup"
	intro: """
		You can configure these options to defer API execution until another external (HTTP) service is up and running.

		[ TIP ]
		====
		Delaying API execution can come in handy in complex setups where execution order is important.
		====
		"""

	options: {
		delay: {
			type:    "time.Duration"
			default: "0"
			description: """
				Delays API startup for the amount of time specified (10s, 2m...).
				    This delay happens before service (`WAIT_FOR_SERVICES`) probing.
				"""
			env: "WAIT_FOR"
		}
		status_page: {
			type:        "bool"
			default:     "true"
			description: "Show temporary status web page."
			env:         "WAIT_FOR_STATUS_PAGE"
		}

		services: {
			description: """
				Space delimited list of hosts and/or URLs to probe.
				    Host format: `host` or `host:443` (port will default to 80).

				[NOTE]
				====
				Services are probed in parallel.
				====
				"""
			env: "WAIT_FOR_SERVICES"
		}

		services_timeout: {
			type:        "time.Duration"
			default:     "time.Minute"
			description: "Max time for each service probe."
			env:         "WAIT_FOR_SERVICES_TIMEOUT"
		}

		services_probe_timeout: {
			type:        "time.Duration"
			default:     "time.Second * 30"
			description: "Timeout for each service probe."
			env:         "WAIT_FOR_SERVICES_PROBE_TIMEOUT"
		}

		services_probe_interval: {
			type:        "time.Duration"
			default:     "time.Second * 5"
			description: "Interval between service probes."
			env:         "WAIT_FOR_SERVICES_PROBE_INTERVAL"
		}
	}
}
