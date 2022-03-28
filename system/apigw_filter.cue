package system

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

apigw_filter: schema.#Resource & {
	features: {
		labels: false
	}

	struct: {
		id: schema.IdField
		route:  { goType: "uint64" }
		weight: { goType: "uint64" }
		ref: {}
		kind: {}
		enabled: {goType: "bool"}
		params: {goType: "types.ApigwFilterParams"}

		created_at: schema.SortableTimestampField
		updated_at: schema.SortableTimestampNilField
		deleted_at: schema.SortableTimestampNilField
		created_by: { goType: "uint64" }
		updated_by: { goType: "uint64" }
		deleted_by: { goType: "uint64" }
	}

	filter: {
		struct: {
			route_id: {goType: "uint64", ident: "routeID"}
			deleted:  {goType: "filter.State", storeIdent: "deleted_at"}
			disabled: {goType: "filter.State", storeIdent: "enabled"}
		}

		byValue: ["route_id"]
		byNilState: ["deleted"]
		byFalseState: ["disabled"]
	}

	store: {
		api: {
			lookups: [
				{
					fields: ["id"]
					description: """
						searches for filter by ID
						"""
				}, {
					fields: ["route"]
					description: """
						searches for filter by route
						"""
				},
			]
		}
	}
}
