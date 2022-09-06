package system

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

auth_oa2token: {
	features: {
		labels: false
		paging: false
		sorting: false
		checkFn: false
	}

	model: {
		attributes: {
				id:     schema.IdField
				code: {}
				access: {}
				refresh: {}
				expires_at: schema.SortableTimestampField
				created_at: schema.SortableTimestampNowField
				data: { goType: "rawJson" }
				client_id: { goType: "uint64", ident: "clientID", storeIdent: "rel_client" }
				user_id: { goType: "uint64", ident: "userID", storeIdent: "rel_user" }
				remote_addr: {}
				user_agent: {}
		}
	}

	filter: {
		struct: {
			user_id: { goType: "uint64", ident: "userID" }
		}

		byValue: ["user_id"]
	}

	store: {
		api: {
			lookups: [
				{ fields: ["id"] },
				{ fields: ["code"] },
				{ fields: ["access"] },
				{ fields: ["refresh"] },
			]

			functions: [
				{ expIdent: "DeleteExpiredAuthOA2Tokens" },
				{ expIdent: "DeleteAuthOA2TokenByCode",    args: [ { ident: "code",    goType: "string" } ] },
				{ expIdent: "DeleteAuthOA2TokenByAccess",  args: [ { ident: "access",  goType: "string" } ] },
				{ expIdent: "DeleteAuthOA2TokenByRefresh", args: [ { ident: "refresh", goType: "string" } ] },
				{ expIdent: "DeleteAuthOA2TokenByUserID",  args: [ { ident: "userID",  goType: "uint64" } ] },
			]
		}
	}
}
