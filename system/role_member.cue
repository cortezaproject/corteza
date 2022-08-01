package system

role_member: {
	features: {
		labels: false
		paging: false
		sorting: false
		checkFn: false
	}

	model: {
		user_id: { goType: "uint64", primaryKey: true, storeIdent: "rel_user", ident: "userID" }
		role_id: { goType: "uint64", primaryKey: true, storeIdent: "rel_role", ident: "roleID" }
	}

	filter: {
		model: {
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
