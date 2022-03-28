package system

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

role_member: schema.#Resource & {
	features: {
		labels: false
		paging: false
		sorting: false
		checkFn: false
	}

	struct: {
		user_id: { goType: "uint64", primaryKey: true, storeIdent: "rel_user", ident: "userID" }
		role_id: { goType: "uint64", primaryKey: true, storeIdent: "rel_role", ident: "roleID" }
	}

	filter: {
		struct: {
			user_id: {goType: "uint64", ident: "userID", storeIdent: "rel_user" }
			role_id: {goType: "uint64", ident: "roleID", storeIdent: "rel_role" }
		}

		byValue: [ "user_id", "role_id"]
	}

	rbac: {
		operations: {
			read: description:             "Read role"
			update: description:           "Update role"
			delete: description:           "Delete role"
			"members.manage": description: "Manage members"
		}}

	store: {
		api: {
			lookups: []
			functions: [
				{
					expIdent: "TransferRoleMembers"
					args: [
						{ident: "src", goType: "uint64"},
						{ident: "dst", goType: "uint64"},
					]
					return: []
				},
			]
		}
	}
}
