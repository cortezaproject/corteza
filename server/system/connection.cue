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
			source: {
				sortable: true
				dal: { type: "Text", length: 16 }
				json: { field: "source", omitEmpty: true }
			}

			meta: {
				goType: "types.ConnectionMeta"
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
			}
			// Virtual sortable mapped onto meta->>'short'.
			name: {
				sortableJSON: { json: "meta.short", accessor: "Meta.Short", nullable: false }
				store: false
				goType: "string"
				omitSetter: true
				omitGetter: true
				envoy: { yaml: { omitEncoder: true } }
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

			// Struct-only fields (not persisted via store/dal).
			derived_params: {
				goType: "[]types.ConnectionDerivedParam"
				store: false
				omitGetter: true
				omitSetter: true
				json: { field: "derivedParams", omitEmpty: true }
			}
			catalog_id: {
				expIdent: "CatalogID"
				goType: "string"
				store: false
				omitGetter: true
				omitSetter: true
				json: { field: "catalogID", omitEmpty: true }
			}
			installed_count: {
				goType: "int"
				store: false
				omitGetter: true
				omitSetter: true
				json: { field: "installedCount", omitEmpty: true }
			}

			created_at: schema.SortableTimestampNowField
			updated_at: schema.SortableTimestampNilField
			deleted_at: schema.SortableTimestampNilField
			created_by: schema.AttributeUserRef & {
				json: { field: "createdBy", string: true }
			}
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
			source: {goType: "string"}
			query:  {goType: "string"}
			deleted: {goType: "filter.State", storeIdent: "deleted_at"}
		}
		byValue: ["handle", "status", "source"]
		byNilState: ["deleted"]
	}

	features: {
		labels: true
	}

	types: {
		gen: true
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

	service: {
		events:   false

		lookup:   false
		search:   false
		update:   false
		undelete: true

		hooks: {
			beforeCreate:   true
			afterCreate:    true
			beforeDelete:   true
			beforeUndelete: true
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
