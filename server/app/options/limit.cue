package options

import (
	"github.com/cortezaproject/corteza/server/codegen/schema"
)

limit: schema.#optionsGroup & {
	handle: "limit"
	options: {
		system_users: {
			type:        "int"
			description: "Maximum number of valid (not deleted, not suspended) users"
		}
		record_count_per_namespace: {
			type: "int"
			description:  "Maximum number of records per namespace"
		}
	}
	title: "Limits"
}
