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
		queue:     {}
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

	rbac: {
		operations: {
			"read": description:        "Read queue"
			"update": description:      "Update queue"
			"delete": description:      "Delete queue"
			"queue.read": description:  "Read from queue"
			"queue.write": description: "Write to queue"
		}
	}

	store: {
		api: {
			lookups: []
		}
	}
}
