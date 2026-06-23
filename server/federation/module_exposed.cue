package federation

import (
	"github.com/crusttech/human/server/codegen/schema"
)

exposedModule: {
	parents: [
		{handle: "node"},
	]

	features: {
		labels: false
		projectScoped: true
	}

	types: {
		gen: true
	}

	model: {
		ident: "federation_module_exposed"
		attributes: {
			id: schema.IdField & {
				json: {field: "moduleID", string: true}
			}
			tenant_id:  schema.TenantRefField
			project_id: schema.ProjectRefField
			handle: schema.HandleField
			name: {
				sortable: true
				dal: {}
			}
			node_id: {
				sortable: true,
				ident: "nodeID",
				goType: "uint64",
				storeIdent: "rel_node"
				dal: { type: "Ref", refModelResType: "corteza::federation:node", default: 0 }
				json: {field: "nodeID", string: true}
			}
			compose_module_id: {
				ident: "composeModuleID",
				goType: "uint64",
				storeIdent: "rel_compose_module"
				dal: { type: "ID" }
				json: {field: "composeModuleID", string: true}
			}
			compose_namespace_id: {
				ident: "composeNamespaceID",
				goType: "uint64",
				storeIdent: "rel_compose_namespace"
				dal: { type: "ID" }
				json: {field: "composeNamespaceID", string: true}
			}
			fields: {
				goType: "types.ModuleFieldSet"
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
			node_id:              { goType: "uint64", ident: "nodeID",             storeIdent: "rel_node" }
			compose_module_id:    { goType: "uint64", ident: "composeModuleID",    storeIdent: "rel_compose_module" }
			compose_namespace_id: { goType: "uint64", ident: "composeNamespaceID", storeIdent: "rel_compose_namespace" }
		}

		byValue: ["compose_module_id", "compose_namespace_id", "node_id"]
	}

	envoy: {
		omit: true
	}

	rbac: {
		operations: {
			"manage": description: "Manage exposed module module"
		}
	}

	service: {
		scoped: true

		customBodyOps: ["lookup", "search", "create", "update", "delete"]

		customAccessOps: ["search", "create"]

		actionProp: "module"
		createProp: "create"
	}

	store: {
		ident: "federationExposedModule"

		api: {
			lookups: [
				{
					fields: ["id"]
					description: """
						searches for federation module by ID

						It returns federation module
						"""
				}
			]
		}
	}
}
