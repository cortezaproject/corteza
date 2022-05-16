package system

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

queue: schema.#Resource & {
	features: {
		labels: false
	}

	struct: {
		id: schema.IdField
		consumer: {goType: "string"}
		queue: {goType: "string"}
		meta: {goType: "types.QueueMeta"}

		created_at: schema.SortableTimestampField
		updated_at: schema.SortableTimestampNilField
		deleted_at: schema.SortableTimestampNilField
		created_by: { goType: "uint64" }
		updated_by: { goType: "uint64" }
		deleted_by: { goType: "uint64" }
	}

	filter: {
		struct: {
			query: {goType: "string"}
			deleted: {goType: "filter.State", storeIdent: "deleted_at"}
		}

		query: ["queue", "consumer"]
		byNilState: ["deleted"]
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

		settings: {
			rdbms: {
				table: "queue_settings"
			}
		}

		api: {
			lookups: [
				{
					fields: ["id"]
					description: """
						searches for queue by ID
						"""
				}, {
					fields: ["queue"]
					description: """
						searches for queue by queue name
						"""
				},
			]
		}
	}
}
