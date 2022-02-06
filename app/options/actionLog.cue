package options

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

actionLog: schema.#optionsGroup & {
	handle: "actionLog"
	options: {
		enabled: {
			type:          "bool"
			defaultGoExpr: "true"
		}
		debug: {
			type: "bool"
		}
		workflow_functions_enabled: {
			type: "bool"
		}
	}
	title: "Actionlog"
}
