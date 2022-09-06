package system

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

apigw_filter: {
	features: {
		labels: false
	}

	model: {
		attributes: {
			id: schema.IdField
			route:  { sortable: true, goType: "uint64", storeIdent: "rel_route" }
			weight: { sortable: true, goType: "uint64" }
			ref: {}
			kind: {sortable: true}
			enabled: {sortable: true, goType: "bool"}
			params: {goType: "types.ApigwFilterParams"}

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
			route_id: {goType: "uint64", ident: "routeID", storeIdent: "rel_route"}
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
