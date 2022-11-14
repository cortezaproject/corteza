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
			type:          "bool"
			defaultGoExpr: "true"
			description:   "Enable eventbus scheduler."
		}
		scheduler_interval: {
			type:        "time.Duration"
			description: "Set time interval for `eventbus` scheduler."

			defaultGoExpr: "time.Minute"
			defaultValue:  "60s"
		}
	}
	title: "Events and scheduler"
}
