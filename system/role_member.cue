package system

role_member: {
	features: {
		labels: false
		paging: false
		sorting: false
		checkFn: false
	}

	model: {
		attributes: {
			user_id: {
				goType: "uint64",
				storeIdent: "rel_user",
				ident: "userID"
				dal: { type: "Ref", refModelResType: "corteza::system:user" }
			}
			role_id: {
				goType: "uint64",
				storeIdent: "rel_role",
				ident: "roleID"
				dal: { type: "Ref", refModelResType: "corteza::system:role" }
			}
		}

		indexes: {
			"primary": { attributes: ["user_id", "role_id"] }
		}
	}

	filter: {
		struct: {
			user_id: {goType: "uint64", ident: "userID", storeIdent: "rel_user" }
			role_id: {goType: "uint64", ident: "roleID", storeIdent: "rel_role" }
		}

		byValue: [ "user_id", "role_id"]
	}

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
