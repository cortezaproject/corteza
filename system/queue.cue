package system

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

queue: {
	features: {
		labels: false
	}

	model: {
		ident: "queue_settings"
		attributes: {
			id: schema.IdField
			consumer: {sortable: true, goType: "string"}
			queue: {sortable: true, goType: "string"}
			meta: {goType: "types.QueueMeta"}

			created_at: schema.SortableTimestampNowField
			updated_at: schema.SortableTimestampNilField
			deleted_at: schema.SortableTimestampNilField
			created_by: schema.AttributeUserRef
			updated_by: schema.AttributeUserRef
			deleted_by: schema.AttributeUserRef
		}
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
