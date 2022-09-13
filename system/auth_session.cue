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
			id: {
				expIdent: "ID",
				goType: string
				dal: { length: 64 }
			}
			data:    {
				goType: "[]byte"
				dal: { type: "Blob" }
			}
			user_id: {
				goType: "uint64",
				ident: "userID",
				storeIdent: "rel_user"
				dal: { type: "Ref", refModelResType: "corteza::system:user", default: 0 }
			}
			remote_addr: {
				dal: {}
			}
			user_agent: {
				dal: {}
			}
			expires_at: schema.SortableTimestampField
			created_at: schema.SortableTimestampNowField
		}

		indexes: {
			"primary": { attribute: "id" }
			"expires_at": { attribute: "expires_at" }
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
