package options

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

actionLog: schema.#optionsGroup & {
	handle: "actionLog"
	options: {
		enabled: {
			type:    "bool"
			default: "true"
		}
		debug: {
			type:    "bool"
			default: "false"
		}
		workflow_functions_enabled: {
			type:    "bool"
			default: "false"
		}
	}
	title: "Actionlog"
}
