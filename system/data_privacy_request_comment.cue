package system

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

data_privacy_request_comment: schema.#Resource & {
	features: {
		labels: false
	}

	struct: {
		id: schema.IdField
		request_id: { ident: "requestID", goType: "uint64", storeIdent: "rel_request" }
		comment: { goType: "string" }

		created_at: schema.SortableTimestampField
		updated_at: schema.SortableTimestampNilField
		deleted_at: schema.SortableTimestampNilField
		created_by: { goType: "uint64" }
		updated_by: { goType: "uint64" }
		deleted_by: { goType: "uint64" }
	}

	filter: {
		struct: {
			request_id: {goType: "[]uint64", ident: "requestID", storeIdent: "rel_request" }
		}

		query: []
		byValue: ["request_id"]
	}

	rbac: false

	store: {
		api: {
			functions: []
		}
	}
}
