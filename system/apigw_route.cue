package system

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

apigw_route: {
	features: {
		labels: false
	}

	model: {
		attributes: {
			id:       schema.IdField
			endpoint: {sortable: true}
			method:   {sortable: true}
			enabled:  {sortable: true, goType: "bool"}
			group:    {sortable: true, goType: "uint64", storeIdent: "rel_group"}
			meta:     {goType: "types.ApigwRouteMeta"}

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
			route: {}
			query: {goType: "string"}

			deleted: {goType: "filter.State", storeIdent: "deleted_at"}
			disabled: {goType: "filter.State", storeIdent: "enabled"}
		}

		query: ["endpoint"]
		byNilState: ["deleted"]
		byFalseState: ["disabled"]
	}


	rbac: {
		operations: {
			read: description:   "Read API Gateway route"
			update: description: "Update API Gateway route"
			delete: description: "Delete API Gateway route"
		}
	}

	store: {
		api: {
			lookups: [
				{
					fields: ["id"]
					description: """
						searches for route by ID

						It returns route even if deleted or suspended
						"""
				}, {
					fields: ["endpoint"]
					description: """
						searches for route by endpoint

						It returns route even if deleted or suspended
						"""
				},
			]
		}
	}
}
