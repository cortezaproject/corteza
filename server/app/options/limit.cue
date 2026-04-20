package options

import (
	"github.com/crusttech/human/server/codegen/schema"
)

limit: schema.#optionsGroup & {
	handle: "limit"
	options: {
		system_users: {
			type:        "int"
			description: "Maximum number of valid (not deleted, not suspended) users"
		}
		record_count_per_module: {
			type: "int"
			description:  "Maximum number of records per module"
		}
	}
	title: "Limits"
}
