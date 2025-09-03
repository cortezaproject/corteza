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
			resource: {
				goType: "string",
				storeIdent: "rel_resource",
				ident: "resource"
				dal: {}
			}
			role_id: {
				goType: "uint64",
				storeIdent: "rel_role",
				ident: "roleID"
				dal: { type: "Ref", refModelResType: "corteza::system:role" }
			}
		}

		indexes: {
			"primary": { attributes: ["resource", "role_id"] }
		}
	}

	filter: {
		struct: {
			resource: {goType: "string", ident: "resource", storeIdent: "rel_resource" }
			role_id: {goType: "uint64", ident: "roleID", storeIdent: "rel_role" }
		}

		byValue: [ "resource", "role_id"]
	}

	envoy: {
		omit: true
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
