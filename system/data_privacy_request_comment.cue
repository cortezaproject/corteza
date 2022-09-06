package system

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

data_privacy_request_comment: {
	features: {
		labels: false
	}

	model: {
		attributes: {
			id: schema.IdField
			request_id: { ident: "requestID", goType: "uint64", storeIdent: "rel_request" }
			comment: { goType: "string" }

			created_at: schema.SortableTimestampNowField
			updated_at: schema.SortableTimestampNilField
			deleted_at: schema.SortableTimestampNilField
			created_by: schema.AttributeUserRef
			updated_by: schema.AttributeUserRef
			deleted_by: schema.AttributeUserRef
		}
	}

	filter: {
		struct: {
			request_id: {goType: "[]uint64", ident: "requestID", storeIdent: "rel_request" }
		}

		query: []
		byValue: ["request_id"]
	}

	store: {
		api: {
			functions: []
		}
	}
}
