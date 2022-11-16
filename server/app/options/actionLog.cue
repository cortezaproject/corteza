package options

import (
	"github.com/cortezaproject/corteza/server/codegen/schema"
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
		compose_record_enabled: {
			type: "bool"
			defaultGoExpr: "false"
			description: """
				Enables actionlog for compose record create, update, and delete. which is disabled by default.

				[IMPORTANT]
				====
				This is temp fix for now, it will be removed completely in future release.
				Once new env var will be introduced for actionlog policy, which will enable more control over action log policies.
				====

			"""
		}
	}
	title: "Actionlog"
}
