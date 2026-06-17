package system

import (
	"github.com/crusttech/human/server/codegen/schema"
)

user_group: {
	features: {
		projectScoped: true
	}

	types: {
		gen: true
		// UserGroupMeta/UserGroupConfig keep hand-written pointer-returning
		// ParseUserGroupMeta/ParseUserGroupConfig (REST request controllers depend
		// on the *T return), incompatible with the value-returning generated
		// versions, so they are excluded from JSON helper generation.
		jsonTypesPtr: ["UserGroupMeta", "UserGroupConfig"]
	}

	model: {
		attributes: {
			id: schema.IdField
			tenant_id:  schema.TenantRefField
			project_id: schema.ProjectRefField
			handle: schema.HandleField
			meta: {
				goType: "*types.UserGroupMeta"
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
			}
			// Virtual sortable mapped onto meta->>'short'.
			name: {
				sortableJSON: { json: "meta.short", accessor: "Meta.Short" }
				store: false
				goType: "string"
				omitSetter: true
				omitGetter: true
				envoy: { yaml: { omitEncoder: true } }
			}

			config: {
				goType: "*types.UserGroupConfig"
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
			}

			// Struct-only computed flag; not persisted, not a getter/setter target.
			is_root: {
				goType: "bool"
				store: false
				omitGetter: true
				omitSetter: true
				json: "isRoot"
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
			user_group_id: {goType: "[]uint64", ident: "userGroupID", storeIdent: "id" }
			tenant_id: schema.TenantFilterField
			project_id: schema.ProjectFilterField
			member_id: {goType: "uint64" }
			handle: {goType: "string"}

			deleted: {goType: "filter.State", storeIdent: "deleted_at"}
			archived: {goType: "filter.State", storeIdent: "archived_at"}
		}

		query: ["handle"]
		byValue: ["user_group_id", "handle"]
		byNilState: ["deleted", "archived"]
	}

	envoy: {
		omit: true
	}

	rbac: {
		operations: {
			read: description:             "Read user group"
			update: description:           "Update user group"
			delete: description:           "Delete user group"
			"members.manage": description: "Manage members"
		}
	}

	service: {
		events:   false

		lookup: false

		search: false

		update: false

		delete: false

		undelete: false

		hooks: {
			validate:     true
			beforeCreate: true
			afterCreate:  true
		}
	}

	store: {
		api: {
			lookups: [
				{
					fields: ["id"]
					description: """
						searches for user group by ID

						It returns user group even if deleted or suspended
						"""
				}, {
					fields: ["handle"]
					nullConstraint: ["deleted_at"]
					constraintCheck: true
					description: """
						searches for user group by handle

						It returns only valid user group (not deleted, not suspended)
						"""
				}
			]
		}
	}
}
