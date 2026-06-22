package system

import (
	"github.com/crusttech/human/server/codegen/schema"
)

configured_connection: {
	model: {
		attributes: {
			id: schema.IdField & {
				json: "configurationID,string"
			}
			tenant_id:  schema.TenantRefField
			project_id: schema.ProjectRefField
			connection_id: {
				sortable: true
				goType:   "uint64"
				ident:    "connectionID"
				storeIdent: "rel_connection"
				dal: { type: "Ref", refModelResType: "corteza::system:connection" }
				json: {field: "connectionID", string: true}
			}
			name: {
				sortable: true
				dal: {}
			}
			status: {
				sortable: true
				dal: { type: "Text", length: 32 }
			}
			connection: {
				goType: "types.Connection"
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
			}
			config: {
				goType: "types.ConfiguredConnectionConfig"
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
			}

			created_at: schema.SortableTimestampNowField
			updated_at: schema.SortableTimestampNilField
			deleted_at: schema.SortableTimestampNilField
			created_by: schema.AttributeUserRef & {
				json: {field: "createdBy", string: true}
			}
			updated_by: schema.AttributeUserRef
			deleted_by: schema.AttributeUserRef
		}

		indexes: {
			"primary": { attribute: "id" }
		}
	}

	filter: {
		struct: {
			tenant_id: schema.TenantFilterField
			project_id: schema.ProjectFilterField
			connection_id: {goType: "uint64", ident: "connectionID", storeIdent: "rel_connection"}
			status:       {goType: "[]string"}
			query:        {goType: "string"}
			deleted:      {goType: "filter.State", storeIdent: "deleted_at"}
		}
		byValue: ["connection_id", "status", "project_id"]
		byNilState: ["deleted"]
	}

	features: {
		labels: true
		projectScoped: true
	}

	types: {
		gen: true
	}

	envoy: {
		omit: true
	}

	rbac: {
		operations: {
			"read": description:   "Read connection"
			"delete": description: "Delete connection"
      "update": description: "Update connection"
		}
	}

	service: {
		lookup: false
		search: false
		update: false
		delete: false

		actionProp: "connection"

		hooks: {beforeCreate: true}
	}

	store: {
		api: {
			lookups: [
				{
					fields: ["id"]
					description: """
						searches for connection connection by ID

						It returns connection connection even if deleted
						"""
				},
			]
		}
	}
}
