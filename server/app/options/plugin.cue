package options

import (
	"github.com/cortezaproject/corteza/server/codegen/schema"
)

plugin: schema.#optionsGroup & {
	handle: "plugin"

	options: {
		enabled: {
			type:          "bool"
			defaultGoExpr: "true"
			description:   "Enable Corteza plugins"
		}
	}
	title: "Corteza Plugins"
}
