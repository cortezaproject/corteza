package system

import (
	"github.com/crusttech/human/server/codegen/schema"
)

data_privacy_request_comment: {
	features: {
		labels: false
	}

	model: {
		attributes: {
			id: schema.IdField & { json: "commentID,string" }
			request_id: {
				ident: "requestID",
				goType: "uint64",
				storeIdent: "rel_request"
				dal: { type: "Ref", refModelResType: "corteza::system:user" }
				json: "requestID,string"
			}
			comment: {
				goType: "string"
				dal: {}
			}

			created_at: schema.SortableTimestampNowField
			updated_at: schema.SortableTimestampNilField
			deleted_at: schema.SortableTimestampNilField
			created_by: schema.AttributeUserRef & { json: "createdBy,string" }
			updated_by: schema.AttributeUserRef
			deleted_by: schema.AttributeUserRef
		}

		indexes: {
			"primary": { attribute: "id" }
		}
	}

	filter: {
		struct: {
			request_id: {goType: "[]uint64", ident: "requestID", storeIdent: "rel_request" }
		}

		query: []
		byValue: ["request_id"]
	}

	types: {
		gen: true
	}

	envoy: {
		omit: true
	}

	store: {
		api: {
			functions: []
		}
	}
}
