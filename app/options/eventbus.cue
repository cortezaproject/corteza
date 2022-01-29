package options

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

eventbus: schema.#optionsGroup & {
	handle: "eventbus"

	imports: [
		"\"time\"",
	]

	options: {
		scheduler_enabled: {
			type:        "bool"
			default:     "true"
			description: "Enable eventbus sheduler."
		}
		scheduler_interval: {
			type:        "time.Duration"
			default:     "time.Minute"
			description: "Set time interval for `eventbus` scheduler."
		}
	}
	title: "Events and scheduler"
}
