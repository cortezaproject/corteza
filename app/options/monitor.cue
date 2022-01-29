package options

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

monitor: schema.#optionsGroup & {
	handle: "monitor"

	imports: [
		"\"time\"",
	]

	options: {
		interval: {
			type:        "time.Duration"
			default:     "300 * time.Second"
			description: "Output (log) interval for monitoring."
		}
	}
	title: "Monitoring"
}
