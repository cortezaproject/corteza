package system

import (
	"github.com/crusttech/human/server/codegen/schema"
)

apigw_filter: {
	features: {
		labels: false
		projectScoped: true
	}

	types: {
		gen: true
	}

	model: {
		attributes: {
			id: schema.IdField & { json: "filterID,string" }
			tenant_id:  schema.TenantRefField
			project_id: schema.ProjectRefField
			route:  {
				sortable: true, goType: "uint64", storeIdent: "rel_route"
				dal: { type: "Ref", refModelResType: "corteza::system:apigw-route" }
				identAlias: ["route", "Route", "ApigwRouteID"]
				json: { field: "routeID", string: true }
				envoy: {
					store: {
						omitRefFilter: true
					}
				}
		  }
			weight: {
			  sortable: true,
			  goType: "uint64"
			  dal: { type: "Number", meta: { "rdbms:type": "integer" } }
			  json: { field: "weight", string: true }
			}
			kind: {
				sortable: true
				dal: { type: "Text", length: 64 }
				json: { field: "kind", omitEmpty: true }
			}
			ref: {
				dal: { type: "Text", length: 64 }
				json: { field: "ref", omitEmpty: true }
			}
			enabled: {
				sortable: true,
				goType: "bool"
				dal: { type: "Boolean" }
				json: { field: "enabled", omitEmpty: true }
			}
			params: {
				goType: "types.ApigwFilterParams"
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
			}

			created_at: schema.SortableTimestampNowField
			updated_at: schema.SortableTimestampNilField
			deleted_at: schema.SortableTimestampNilField
			created_by: schema.AttributeUserRef & { json: { field: "createdBy", string: true } }
			updated_by: schema.AttributeUserRef
			deleted_by: schema.AttributeUserRef
		}

		indexes: {
			"primary": { attribute: "id" }
		}
	}

	envoy: {
		yaml: {
			supportMappedInput: false
			omitEncoder: true
		}
		store: {
			handleField: ""
		}
	}

	filter: {
		struct: {
			apigw_filter_id: {goType: "[]uint64", ident: "apigwFilterID", storeIdent: "id"}
			tenant_id: schema.TenantFilterField
			project_id: schema.ProjectFilterField
			route_id: {goType: "uint64", ident: "routeID", storeIdent: "rel_route"}
			deleted:  {goType: "filter.State", storeIdent: "deleted_at"}
			disabled: {goType: "filter.State", storeIdent: "enabled"}
		}

		byValue: ["apigw_filter_id", "route_id"]
		byNilState: ["deleted"]
		byFalseState: ["disabled"]
	}

	service: {
		// Every CRUD op is bespoke: access control delegates to the parent
		// ApigwRoute (Can*ApigwRoute), each op loads the route, validates and
		// fires an endpoint reload side-effect. So all ops are custom-bodied and
		// custom-access. Undelete is left hand-written: the original records
		// ApigwFilterActionDelete (not Undelete) which the generated body cannot
		// reproduce, so it stays out of the generator.
		customBodyOps:   ["lookup", "search", "create", "update", "delete"]
		customAccessOps: ["lookup", "search", "create", "update", "delete"]

		// action-log props: resource prop is "filter", search prop is "search",
		// Create/Update reuse "filter" (no dedicated new/update field on props).
		actionProp:     "filter"
		filterProp:     "search"
		omitCreateProp: true
		updateProp:     "filter"
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
