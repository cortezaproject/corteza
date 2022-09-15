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
			code: {
				dal: { type: "Text", length: 48 }
			}
			access: {
				dal: { type: "Text", length: 2048 }
			}
			refresh: {
				dal: { type: "Text", length: 48 }
			}
			data: {
				goType: "rawJson"
				dal: { type: "JSON", defaultEmptyObject: true }
			}
			remote_addr: {
				dal: { type: "Text", length: 64 }
			}
			user_agent: {
				dal: {}
			}
			client_id: {
				goType: "uint64",
				ident: "clientID",
				storeIdent: "rel_client"
				dal: { type: "Ref", refModelResType: "corteza::system:auth-client", default: 0 }
			}
			user_id: {
				goType: "uint64",
				ident: "userID",
				storeIdent: "rel_user"
				dal: { type: "Ref", refModelResType: "corteza::system:user", default: 0 }
			}
			created_at: schema.SortableTimestampNowField
			expires_at: schema.SortableTimestampField
		}

		indexes: {
			"primary": { attribute: "id" }
			"client_id": { attribute: "client_id" }
			"code": { attribute: "code" }
			"refresh": { attribute: "refresh" }
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
