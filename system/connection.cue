package system

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

connection: schema.#Resource & {
	struct: {
		id:     schema.IdField
		handle: schema.HandleField
		dsn: {expIdent: "DSN"}
		location: {}
		ownership: {}
		sensitive: {goType: "bool"}

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
			dsn: {expIdent: "DSN", goType: "string"}
			location: {goType: "string"}
			ownership: {goType: "string"}
			sensitive: {goType: "bool"}

			deleted: {goType: "filter.State", storeIdent: "deleted_at"}
		}

		byValue: ["connection_id", "handle", "dsn", "location", "ownership"]
		byNilState: ["deleted"]
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
