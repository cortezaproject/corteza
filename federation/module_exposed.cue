package federation

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

exposedModule: {
	parents: [
		{handle: "node"},
	]

	features: {
		labels: false
	}

	model: {
		ident: "federation_module_exposed"
		attributes: {
			id:          schema.IdField
			handle:      schema.HandleField
			name: { sortable: true }
			node_id: { sortable: true, ident: "nodeID", goType: "uint64", storeIdent: "rel_node" }
			compose_module_id: { ident: "composeModuleID", goType: "uint64", storeIdent: "rel_compose_module" }
			compose_namespace_id: { ident: "composeNamespaceID", goType: "uint64", storeIdent: "rel_compose_namespace" }
			fields: { goType: "types.ModuleFieldSet" }

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
			node_id:              { goType: "uint64", ident: "nodeID",             storeIdent: "rel_node" }
			compose_module_id:    { goType: "uint64", ident: "composeModuleID",    storeIdent: "rel_compose_module" }
			compose_namespace_id: { goType: "uint64", ident: "composeNamespaceID", storeIdent: "rel_compose_namespace" }
		}

		byValue: ["compose_module_id", "compose_namespace_id", "node_id"]
	}


	rbac: {
		operations: {
			"manage": description: "Manage exposed module module"
		}
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
