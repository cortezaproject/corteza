package options

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

plugins: schema.#optionsGroup & {
	handle: "plugins"
	options: {
		Enabled: {
			type:          "bool"
			defaultGoExpr: "true"
			description:   "Enable plugins"
		}
		Paths: {
			description: "List of colon seperated paths or patterns where plugins could be found"
		}
	}
	title: "Plugins"
}
