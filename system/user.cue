package system

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

user: {
	model: {
		attributes: {
		  id:     schema.IdField
		  email: {
		  	sortable: true,
		  	unique: true,
		  	ignoreCase: true
				dal: { length: 254 }
			}
		  email_confirmed: {
		  	goType: "bool"
				dal: { type: "Boolean" }
			}
		  username: {
		  	sortable: true,
		  	unique: true,
		  	ignoreCase: true
				dal: {}

			}
		  name: {
		  	sortable: true
				dal: {}
			}
		  handle: schema.HandleField
		  kind: {
		  	sortable: true,
		  	goType: "types.UserKind"
				dal: { length: 8 }
			}
		  meta: {
		  	goType: "*types.UserMeta"
				dal: { type: "JSON", defaultEmptyObject: true }
			}
		  suspended_at: schema.SortableTimestampNilField
		  created_at: schema.SortableTimestampNowField
		  updated_at: schema.SortableTimestampNilField
		  deleted_at: schema.SortableTimestampNilField
		}

		indexes: {
			"primary": { attribute: "id" }
			"unique_email": {
				 fields: [{ attribute: "email", modifiers: ["LOWERCASE"] }]
				 predicate: "email != '' AND deleted_at IS NULL"
		 	}
			"unique_handle": {
				 fields: [{ attribute: "handle", modifiers: ["LOWERCASE"] }]
				 predicate: "handle != '' AND deleted_at IS NULL"
		 	}
			"unique_username": {
				 fields: [{ attribute: "username", modifiers: ["LOWERCASE"] }]
				 predicate: "username != '' AND deleted_at IS NULL"
		 	}
		}
	}

	filter: {
		struct: {
			user_id: {goType: "[]uint64", ident: "userID", storeIdent: "id"}
			role_id: {goType: "[]uint64", ident: "roleID"}
			email: {goType: "string"}
			name: {goType: "string"}
			username: {goType: "string"}
			handle: {goType: "string"}
			kind: {goType: "types.UserKind"}
			allKinds: {goType: "bool"}

			deleted: {goType: "filter.State", storeIdent: "deleted_at"}
			suspended: {goType: "filter.State", storeIdent: "suspended_at"}
		}

		query: ["email", "username", "handle", "name"]
		byValue: ["user_id", "email", "username", "handle"]
		byNilState: ["deleted", "suspended"]
	}

	rbac: {
		operations: {
			"read": description:         "Read user"
			"update": description:       "Update user"
			"delete": description:       "Delete user"
			"suspend": description:      "Suspend user"
			"unsuspend": description:    "Unsuspend user"
			"email.unmask": description: "Unmask email"
			"name.unmask": description:  "Unmask name"
			"impersonate": description:  "Impersonate user"
			"credentials.manage": description: "Manage user's credentials"
		}
	}

	store: {
		api: {
			lookups: [
				{
					fields: ["id"]
					description: """
						searches for user by ID

						It returns user even if deleted or suspended
						"""
				}, {
					fields: ["email"]
					nullConstraint: ["deleted_at"]
					constraintCheck: true
					description: """
						searches for user by email

						It returns only valid user (not deleted, not suspended)
						"""
				}, {
					fields: ["handle"]
					nullConstraint: ["deleted_at"]
					constraintCheck: true
					description: """
						searches for user by handle

						It returns only valid user (not deleted, not suspended)
						"""
				}, {
					fields: ["username"]
					nullConstraint: ["deleted_at"]
					constraintCheck: true
					description: """
						searches for user by username

						It returns only valid user (not deleted, not suspended)
						"""
				},
			]

			functions: [
				{
					expIdent: "CountUsers"
					args: [ {ident: "u", goType: "types.UserFilter"}]
					return: [ "uint"]
				}, {
					expIdent: "UserMetrics"
					return: [ "*types.UserMetrics"]
				},
			]
		}
	}
}
