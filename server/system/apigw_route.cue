package system

import (
	"github.com/cortezaproject/corteza/server/codegen/schema"
)

apigw_route: {
	features: {
		labels: false
	}

	model: {
		attributes: {
			id:       schema.IdField
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
			created_by: schema.AttributeUserRef
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
