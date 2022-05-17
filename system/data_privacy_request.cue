package system

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

data_privacy_request: schema.#Resource & {
	struct: {
		id: schema.IdField
		name: {}

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
