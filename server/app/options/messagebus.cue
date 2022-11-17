package options

import (
	"github.com/cortezaproject/corteza/server/codegen/schema"
)

messagebus: schema.#optionsGroup & {
	handle: "messagebus"
	options: {
		Enabled: {
			type:          "bool"
			defaultGoExpr: "true"
			description:   "Enable messagebus"
		}
		log_enabled: {
			type:        "bool"
			description: "Enable extra logging for messagebus watchers"
		}
	}
	title: "Messaging queue"
}
