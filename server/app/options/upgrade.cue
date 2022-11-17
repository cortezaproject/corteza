package options

import (
	"github.com/cortezaproject/corteza/server/codegen/schema"
)

upgrade: schema.#optionsGroup & {
	handle: "upgrade"
	title:  "Data store (database) upgrade"

	options: {
		debug: {
			type: "bool"
			description: """
				Enable/disable debug logging.
				    To enable debug logging set `UPGRADE_DEBUG=true`.
				"""
		}
		always: {
			type:          "bool"
			defaultGoExpr: "true"
			description:   "Controls if the upgradable systems should be upgraded when the server starts."
		}
	}
}
