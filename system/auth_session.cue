package system

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

auth_session: {
	features: {
		labels: false
		paging: false
		sorting: false
		checkFn: false
	}

	model: {
		attributes: {
				id:      { primaryKey: true, expIdent: "ID" }
				data:    { goType: "[]byte" }
				user_id: { goType: "uint64", storeIdent: "rel_user", ident: "userID"}
				expires_at: schema.SortableTimestampField
				created_at: schema.SortableTimestampNowField
				remote_addr: {}
				user_agent: {}
		}
	}

	filter: {
		struct: {
			user_id: { goType: "uint64", ident: "userID", storeIdent: "rel_user" }
		}

		byValue: ["user_id"]
	}

	store: {
		api: {
			lookups: [
				{ fields: ["id"] },
			]

			functions: [
				{ expIdent: "DeleteExpiredAuthSessions" },
				{ expIdent: "DeleteAuthSessionsByUserID",  args: [ { ident: "userID",  goType: "uint64" } ] },
			]
		}
	}
}
