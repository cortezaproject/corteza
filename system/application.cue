package system

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

application: {
	model: {
		attributes: {
			id: schema.IdField
			name: {
				sortable: true
				dal: {}
			}
			enabled: {
				goType: "bool"
				sortable: true,
				dal: { type: "Boolean", default: true }
			}
			weight: {
				goType: "int",
				sortable: true
				dal: { type: "Number", default: 0, meta: { "rdbms:type": "integer" } }
			}
			unify: {
				goType: "*types.ApplicationUnify"
				dal: { type: "JSON", defaultEmptyObject: true }
			}
			owner_id:   {
				schema.AttributeUserRef,
				storeIdent: "rel_owner",
				ident: "ownerID"
			}
			created_at: schema.SortableTimestampNowField
			updated_at: schema.SortableTimestampNilField
			deleted_at: schema.SortableTimestampNilField
		}

		indexes: {
			"primary": { attribute: "id" }
		}
	}

	filter: {
		struct: {
			name: {goType: "string"}
			// not sure about the type of flagged_ids
			flagged_ids: {goType: "[]uint64"}
			flags: {goType: "[]string"}
			inc_flags: {goType: "uint"}
			deleted: {goType: "filter.State", storeIdent: "deleted_at"}
		}

		query: ["name"]
		byValue: ["name"]
		byNilState: ["deleted"]
	}

	features: {
		flags: true
	}

	rbac: {
		operations: {
			read:
				description: "Read application"
			update:
				description: "Update application"
			delete:
				description: "Delete application"
		}
	}

	store: {
		api: {
			lookups: [
				{
					fields: ["id"]
					description: """
						searches for role by ID

						It returns role even if deleted or suspended
						"""
				},
			]

			functions: [
				{
					expIdent: "ApplicationMetrics"
					return: [ "*types.ApplicationMetrics"]
				}, {
					expIdent: "ReorderApplications"
					// not sure about the ident
					args: [ {ident: "order", goType: "[]uint64"}]
				},
			]
		}
	}
}
