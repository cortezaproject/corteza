package system

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

role: {
	model: {
		attributes: {
			id: schema.IdField
			name: {
				sortable: true
				dal: {}
			}
			handle: schema.HandleField
			meta: {
				goType: "*types.RoleMeta"
				dal: { type: "JSON", defaultEmptyObject: true }
			}

			archived_at: schema.SortableTimestampNilField
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
			role_id: {goType: "[]uint64", ident: "roleID", storeIdent: "id" }
			member_id: {goType: "uint64" }
			handle: {goType: "string"}
			name: {goType: "string"}

			deleted: {goType: "filter.State", storeIdent: "deleted_at"}
			archived: {goType: "filter.State", storeIdent: "archived_at"}
		}

		query: ["handle", "name"]
		byValue: ["role_id", "name", "handle"]
		byNilState: ["deleted", "archived"]
	}

	rbac: {
		operations: {
			read: description:             "Read role"
			update: description:           "Update role"
			delete: description:           "Delete role"
			"members.manage": description: "Manage members"
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
				}, {
					fields: ["handle"]
					nullConstraint: ["deleted_at"]
					constraintCheck: true
					description: """
						searches for role by handle

						It returns only valid role (not deleted, not suspended)
						"""
				}, {
					fields: ["name"]
					nullConstraint: ["deleted_at"]
					constraintCheck: true
					description: """
						searches for role by name

						It returns only valid role (not deleted, not suspended)
						"""
				},
			]
			functions: [
				{
					expIdent: "RoleMetrics"
					return: [ "*types.RoleMetrics"]
				},
			]
		}
	}
}
