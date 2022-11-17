package options

import (
	"github.com/cortezaproject/corteza/server/codegen/schema"
)

template: schema.#optionsGroup & {
	handle: "template"
	title:  "Rendering engine"

	options: {
		renderer_gotenberg_address: {
			defaultGoExpr: ""
			description:   "Gotenberg rendering container address."
		}

		renderer_gotenberg_enabled: {
			type:        "bool"
			description: "Is Gotenberg rendering container enabled."
		}
	}
}
