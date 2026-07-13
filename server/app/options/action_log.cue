package options

import (
	"github.com/crusttech/human/server/codegen/schema"
)

actionLog: schema.#optionsGroup & {
	handle: "action-log"
	env: "ACTIONLOG"
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
		DB_DSN: {
			description: "Database connection string for the action log; when empty the primary DB connection is used."
		}
	}
	title: "Actionlog"
}
