package system

import (
	"github.com/crusttech/human/server/codegen/schema"
)

data_privacy_request: {
	features: {
		labels: false
		projectScoped: true
	}

	types: {
		gen: true
	}

	model: {
		omitGetterSetter: true
		attributes: {
			id: schema.IdField & { json: "requestID,string" }
			tenant_id:  schema.TenantRefField
			project_id: schema.ProjectRefField

			kind: {
				goType: "types.RequestKind",
				sortable: true
				dal: {}
			}

			status: {
				goType: "types.RequestStatus",
				sortable: true
				dal: { type: "Text", length: 64 }

			}
			payload: {
				goType: "types.DataPrivacyRequestPayloadSet"
				json: "payload,omitempty"
				dal: { type: "JSON" }
			}

			requested_at: schema.SortableTimestampField
			requested_by: {
				goType: "uint64"
				json: "requestedBy,string"
				dal: { type: "Ref", refModelResType: "corteza::system:user" }
			}

			completed_at: schema.SortableTimestampNilField
			completed_by: {
				goType: "uint64"
				dal: { type: "Ref", refModelResType: "corteza::system:user" }
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
			request_id: {goType: "[]uint64", ident: "requestID", storeIdent: "id" }
			tenant_id: schema.TenantFilterField
			project_id: schema.ProjectFilterField
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

	envoy: {
		omit: true
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
