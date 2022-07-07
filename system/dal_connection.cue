package system

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

dal_connection: schema.#Resource & {
	struct: {
		id:     schema.IdField
		name: { sortable: true, goType: "string" }
		handle: schema.HandleField

		// omitting isPrimary and replacing with a special type
		type: { sortable: true }

		location: { goType: "geolocation.Full" }
		ownership: {}
		sensitivity_level: { goType: "uint64" }

		config: {goType: "types.ConnectionConfig"}
		capabilities: {goType: "types.ConnectionCapabilities"}

		created_at: schema.SortableTimestampField
		updated_at: schema.SortableTimestampNilField
		deleted_at: schema.SortableTimestampNilField
		created_by: { goType: "uint64" }
		updated_by: { goType: "uint64" }
		deleted_by: { goType: "uint64" }
	}

	filter: {
		struct: {
			connection_id: {goType: "[]uint64", ident: "connectionID", storeIdent: "id"}
			handle: {goType: "string"}
			type: {goType: "string"}

			deleted: {goType: "filter.State", storeIdent: "deleted_at"}
		}

		byValue: ["connection_id", "handle", "type"]
		byNilState: ["deleted"]
	}

	features: {
		labels: false
		paging: false
		sorting: false
	}

	rbac: {
		operations: {
			"read": description:         "Read connection"
			"update": description:       "Update connection"
			"delete": description:       "Delete connection"
		}
	}

	store: {
		api: {
			lookups: [
				{
					fields: ["id"]
					description: """
						searches for connection by ID

						It returns connection even if deleted or suspended
						"""
				}, {
					fields: ["handle"]
					nullConstraint: ["deleted_at"]
					constraintCheck: true
					description: """
						searches for connection by handle

						It returns only valid connection (not deleted)
						"""
				},
			]
		}
	}
}
