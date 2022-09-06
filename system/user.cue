package system

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

user: {
	model: {
		attributes: {
				id:     schema.IdField
				handle: schema.HandleField
				email: {sortable: true, unique: true, ignoreCase: true}
				email_confirmed: {goType: "bool"}
				username: {sortable: true, unique: true, ignoreCase: true}
				name: {sortable: true}
				kind: {sortable: true, goType: "types.UserKind"}
				meta: {goType: "*types.UserMeta"}
				created_at: schema.SortableTimestampNowField
				updated_at: schema.SortableTimestampNilField
				deleted_at: schema.SortableTimestampNilField
				suspended_at: schema.SortableTimestampNilField
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
