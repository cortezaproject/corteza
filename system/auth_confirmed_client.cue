package system

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

auth_confirmed_client: schema.#Resource & {
	features: {
		labels: false
		paging: false
		sorting: false
		checkFn: false
	}

	struct: {
		user_id:      { goType: "uint64", primaryKey: true, ident: "userID", storeIdent: "rel_user" }
		client_id:    { goType: "uint64", primaryKey: true, ident: "clientID", storeIdent: "rel_client" }
		confirmed_at: schema.SortableTimestampField
	}

	filter: {
		struct: {
			user_id:   { goType: "uint64", ident: "userID", storeIdent: "rel_client" }
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
