package system

import (
	"github.com/cortezaproject/corteza/server/codegen/schema"
)

user_group: {
	model: {
		attributes: {
			id: schema.IdField
			handle: schema.HandleField
			meta: {
				goType: "*types.UserGroupMeta"
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
			}

			self_id: {
				ident: "selfID",
				goType: "uint64",
				dal: { type: "Ref", refModelResType: "corteza::system:user-group" }
				sortable: true
				envoy: {
					store: {
						filterRefField: "ParentID"
					}
					yaml: {
						identKeyAlias: ["parent"]
					}
				}
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
