package system

import (
	"github.com/crusttech/human/server/codegen/schema"
)

connection: {
	model: {
		attributes: {
			id:       schema.IdField
			handle:   schema.HandleField
			revision: {
				sortable: true
				goType:   "int"
				dal: { type: "Number", meta: { "rdbms:type": "integer" } }
			}
			status: {
				sortable: true
				dal: { type: "Text", length: 32 }
			}

			meta: {
				goType: "types.ConnectionMeta"
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
			}
			service: {
				goType: "types.ConnectionService"
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
			}
			resources: {
				goType: "types.ConnectionResources"
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
			}
			operations: {
				goType: "types.ConnectionOperations"
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
			}

			created_at: schema.SortableTimestampNowField
			updated_at: schema.SortableTimestampNilField
			deleted_at: schema.SortableTimestampNilField
			created_by: schema.AttributeUserRef
			updated_by: schema.AttributeUserRef
			deleted_by: schema.AttributeUserRef
		}

		indexes: {
			"primary": { attribute: "id" }
		}
	}

	filter: {
		struct: {
			handle: {goType: "string"}
			status: {goType: "[]string"}
			query:  {goType: "string"}
			deleted: {goType: "filter.State", storeIdent: "deleted_at"}
		}
		byValue: ["handle", "status"]
		byNilState: ["deleted"]
	}

	features: {
		labels: true
	}

	envoy: {
		omit: true
	}

	rbac: {
		operations: {
			"read": description:   "Read connection"
			"update": description: "Update connection"
			"delete": description: "Delete connection"
			"install": description: "Install connection"
		}
	}

	store: {
		api: {
			lookups: [
				{
					fields: ["id"]
					description: """
						searches for connection by ID

						It returns connection even if deleted
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
