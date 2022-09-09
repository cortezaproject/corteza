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
		ident: "federation_module_shared"
		attributes: {
			id: schema.IdField
			handle: schema.HandleField
			node_id: {
				sortable: true,
				ident: "nodeID",
				goType: "uint64",
				storeIdent: "rel_node"
				dal: { type: "ID" }
			}
			name: {
				sortable: true
				dal: {}
			}
			external_federation_module_id: {
				sortable: true,
				ident: "externalFederationModuleID",
				goType: "uint64",
				storeIdent: "xref_module",
				dal: { type: "ID" }
			}
			fields: {
				goType: "types.ModuleFieldSet"
				dal: { type: "JSON", defaultEmptyObject: true }
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

	filter: {
		struct: {
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
