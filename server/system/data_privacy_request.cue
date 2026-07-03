package system

import (
	"github.com/crusttech/human/server/codegen/schema"
)

_data_privacy_requestDefs: {
			RequestKind: { name: "RequestKind", values: [
				{ident: "RequestKindCorrect", value: "correct", doc: "RequestKindCorrect to correct module fields"},
				{ident: "RequestKindDelete", value: "delete", doc: "RequestKindDelete to delete module fields"},
				{ident: "RequestKindExport", value: "export", doc: "RequestKindExport to export module fields"},
			]}
			RequestStatus: { name: "RequestStatus", values: [
				{ident: "RequestStatusPending", value: "pending", doc: "RequestStatusPending initially request will be in pending status"},
				{ident: "RequestStatusCanceled", value: "canceled", doc: "RequestStatusCanceled owner of request has cancelled the request"},
				{ident: "RequestStatusApproved", value: "approved", doc: "RequestStatusApproved data officer has of request has cancelled the request"},
				{ident: "RequestStatusRejected", value: "rejected", doc: "RequestStatusRejected data officer has denied the request"},
			]}
			DataPrivacyRequestPayloadSet: { name: "DataPrivacyRequestPayloadSet", elemGoType: "map[string]any" }
		}

data_privacy_request: {
	features: {
		labels: false
	}

	types: {
		gen: true
		defs: _data_privacy_requestDefs
	}

	model: {
		omitGetterSetter: true
		attributes: {
			id: schema.IdField & { json: "requestID,string" }
			tenant_id:  schema.TenantRefField
			project_id: schema.ProjectRefField

			kind: {
				type: _data_privacy_requestDefs.RequestKind,
				sortable: true
				dal: {}
			}

			status: {
				type: _data_privacy_requestDefs.RequestStatus,
				sortable: true
				dal: { type: "Text", length: 64 }

			}
			payload: {
				type: _data_privacy_requestDefs.DataPrivacyRequestPayloadSet
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
