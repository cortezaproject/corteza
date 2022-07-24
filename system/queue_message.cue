package system

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

queue_message: schema.#Resource & {
	features: {
		labels: false
		checkFn: false
	}

	struct: {
		id:        schema.IdField
		queue:     { sortable: true }
		payload:   { goType: "[]byte" }
		processed: schema.SortableTimestampNilField
		created:   schema.SortableTimestampNilField
	}

	filter: {
		struct: {
			queue: {}
			processed: {goType: "filter.State", storeIdent: "processed"}
		}

		byValue: ["queue"]
		byNilState: ["processed"]
	}

	rbac: false

	store: {
		api: {
			lookups: []
		}
	}
}
