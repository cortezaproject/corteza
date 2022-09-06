package system

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

data_privacy_request: {
	features: {
		labels: false
	}

	model: {
		attributes: {
			id: schema.IdField
			kind: { goType: "types.RequestKind", sortable: true }
			status: { goType: "types.RequestStatus", sortable: true }
			payload: { goType: "types.DataPrivacyRequestPayloadSet" }

			requested_at: schema.SortableTimestampField
			requested_by: { goType: "uint64" }
			completed_at: schema.SortableTimestampNilField
			completed_by: { goType: "uint64" }

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
			request_id: {goType: "[]uint64", ident: "requestID", storeIdent: "id" }
			requested_by: {goType: "[]uint64", ident: "requestedBy" }
			kind: {goType: "[]types.RequestKind"}
			status: {goType: "[]types.RequestStatus"}
		}

		query: ["kind", "status"]
		byValue: ["kind", "status", "requested_by"]
	}

	rbac: {
		operations: {
			read:
				description: "Read data privacy request"
			approve:
				description: "Approve/Reject data privacy request"
		}
	}

	store: {
		api: {
			lookups: [
				{
					fields: ["id"]
					description: """
						searches for data privacy request by ID

						It returns data privacy request even if deleted
						"""
				}
			]
			functions: []
		}
	}
}
