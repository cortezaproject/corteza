package system

import (
	"github.com/crusttech/human/server/codegen/schema"
)

apigw_route: {
	features: {
		labels: false
		projectScoped: true
	}

	types: {
		gen: true
	}

	model: {
		attributes: {
			id:       schema.IdField & {
				// hand-written tag uses custom "routeID" name (not the apigwRoute resource ident)
				json: "routeID,string"
			}
			tenant_id:  schema.TenantRefField
			project_id: schema.ProjectRefField
			endpoint: {
				sortable: true
				dal: {}
			}
			method:   {
				sortable: true
				dal: {}
			}
			enabled: {
				sortable: true,
				goType: "bool"
				dal: { type: "Boolean" }
			}
			meta: {
				goType: "types.ApigwRouteMeta"
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
			}
			group:    {
				sortable: true,
				goType: "uint64",
				storeIdent: "rel_group"
				// hand-written tag is `group,string` (no omitempty)
				json: {field: "group", string: true}
			  dal: {
			  	type: "Ref",
			  	// @todo what does this do?
			  	refModelResType: "corteza::system:apigw-group"
				}
				envoy: {
					store: {
						omitRefFilter: true
					}
				}
			}

			created_at: schema.SortableTimestampNowField
			updated_at: schema.SortableTimestampNilField
			deleted_at: schema.SortableTimestampNilField
			created_by: schema.AttributeUserRef & {
				// hand-written tag is `createdBy,string` (no omitempty)
				json: {field: "createdBy", string: true}
			}
			updated_by: schema.AttributeUserRef
			deleted_by: schema.AttributeUserRef
		}

		indexes: {
			"primary": { attribute: "id" }
		}
	}

	envoy: {
		yaml: {
			supportMappedInput: true
			mappedField: "Endpoint"
			identKeyAlias: ["endpoints"]
			extendedResourceDecoders: [{
				ident: "filters"
				expIdent: "Filters"
				identKeys: ["filters"]
				supportMappedInput: false
			}]
			extendedResourceEncoders: [{
				ident: "apigwFilter"
				expIdent: "ApigwFilter"
				identKey: "filters"
			}]
		}
		store: {
			handleField: "Endpoint"
			extendedDecoder: true
		}
	}

	filter: {
		struct: {
			apigw_route_id: { goType: "[]uint64", ident: "apigwRouteID", storeIdent: "id" }
			tenant_id: schema.TenantFilterField
			project_id: schema.ProjectFilterField
			route: {goType: "string", storeIdent: "id"}
			endpoint: {goType: "string"}
			method: {goType: "string"}

			deleted: {goType: "filter.State", storeIdent: "deleted_at"}
			disabled: {goType: "filter.State", storeIdent: "enabled"}
		}

		byValue: ["apigw_route_id", "route", "method"]
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

	service: {
		// action-log resource prop is "route" (not the "apigwRoute" ident)
		actionProp: "route"
		// Search records under the "search" action-log prop
		filterProp: "search"

		undelete: true

		// lookup + search use the standard generated bodies (FindByID/Search);
		// create/update/delete/undelete carry bespoke logic (CreatedBy/Group
		// defaulting, endpoint-moved 404 handling, apigw reload/not-found
		// signalling, soft-delete via UpdateApigwRoute) handled by on<Op>.
		customBodyOps: ["create", "update", "delete", "undelete"]
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
