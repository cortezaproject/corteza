package system

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

data_privacy_request: schema.#Resource & {
	features: {
		labels: false
	}

	struct: {
		id: schema.IdField
		name: {}
		kind: { goType: "types.RequestKind" }
		status: { goType: "types.RequestStatus" }

		requested_at: { schema.SortableTimestampField, ident: "requestedAt", storeIdent: "requested_at" }
		requested_by: { goType: "uint64", ident: "requestedBy", storeIdent: "requested_by" }
		completed_at: { schema.SortableTimestampNilField, ident: "completedAt", storeIdent: "completed_at" }
		completed_by: { goType: "uint64", ident: "completedBy", storeIdent: "completed_by" }

		created_at: schema.SortableTimestampField
		updated_at: schema.SortableTimestampNilField
		deleted_at: schema.SortableTimestampNilField
		created_by: { goType: "uint64" }
		updated_by: { goType: "uint64" }
		deleted_by: { goType: "uint64" }
	}

	filter: {
		struct: {
			request_id: {goType: "[]uint64", ident: "requestID", storeIdent: "id" }
			name: {goType: "string"}
			status: {goType: "[]types.RequestStatus"}

			deleted: {goType: "filter.State", storeIdent: "deleted_at"}
		}
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
				}, {
					fields: ["name"]
					nullConstraint: ["deleted_at"]
					constraintCheck: true
					description: """
						searches for data privacy request by name

						It returns only valid data privacy request (not deleted)
						"""
				},
			]
			functions: []
		}
	}
}
