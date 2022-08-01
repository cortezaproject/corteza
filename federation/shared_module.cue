package federation

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

sharedModule: {
	features: {
		labels: false
	}

	parents: [
		{handle: "node"},
	]

	model: {
		id:          schema.IdField
		handle:      schema.HandleField
		node_id: { sortable: true, ident: "nodeID", goType: "uint64", storeIdent: "rel_node" }
		name: {sortable: true}
		external_federation_module_id: { sortable: true, ident: "externalFederationModuleID", goType: "uint64", storeIdent: "xref_module",  }
		fields: { goType: "types.ModuleFieldSet" }

		created_at: schema.SortableTimestampField
		updated_at: schema.SortableTimestampNilField
		deleted_at: schema.SortableTimestampNilField
		created_by: { goType: "uint64" }
		updated_by: { goType: "uint64" }
		deleted_by: { goType: "uint64" }
	}

	filter: {
		model: {
			node_id:  { goType: "uint64", ident: "nodeID", storeIdent: "rel_node" }
			handle:   { goType: "string" }
			name:     { goType: "string" }
			external_federation_module_id: { goType: "uint64", storeIdent: "xref_module", ident: "externalFederationModuleID" }
		}

		query: ["name", "handle"]
		byValue: ["handle", "node_id", "name", "external_federation_module_id"]
	}

	rbac: {
		operations: {
			"map": description: "Map shared module"
		}
	}

	store: {
		ident: "federationSharedModule"

		settings: {
			rdbms: {
				table: "federation_module_shared"
			}
		}

		api: {
			lookups: [
				{
					fields: ["id"]
					description: """
						searches for shared federation module by ID

						It returns shared federation module
						"""
				}
			]
		}
	}
}
