package system

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

auth_confirmed_client: {
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
				ident: "userID",
				storeIdent: "rel_user"
				dal: { type: "Ref", refModelResType: "corteza::system:user", default: 0 }
		  }
			client_id: {
				goType: "uint64",
				ident: "clientID",
				storeIdent: "rel_client"
				dal: { type: "Ref", refModelResType: "corteza::system:auth-client", default: 0 }
			}
			confirmed_at: schema.SortableTimestampNowField
		}

		indexes: {
			"primary": { attributes: ["user_id", "client_id"] }
		}
	}

	filter: {
		struct: {
			user_id:   { goType: "uint64", ident: "userID", storeIdent: "rel_user" }
		}

		byValue: ["user_id"]
	}


	store: {
		api: {
			lookups: [
				{ fields: ["user_id", "client_id"] },
			]
		}
	}
}
