package federation

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

nodeSync: {
	features: {
		labels: false
	}

	model: {
		ident: "federation_nodes_sync"
		attributes: {
			node_id: {
				sortable: true,
				ident: "nodeID",
				goType: "uint64",
				dal: { type: "ID" }
			}
			module_id: {
				sortable: true,
				ident: "moduleID",
				goType: "uint64"
				dal: { type: "ID" }
			}
			sync_type: {
				sortable: true,
				goType: "string"
				dal: {}
			}
			sync_status: {
				sortable: true,
				goType: "string"
				dal: {}
			}
			time_of_action: schema.SortableTimestampField
		}
	}

	filter: {
		struct: {
			node_id:     { goType: "uint64", storeIdent: "rel_node",   ident: "nodeID" }
			module_id:   { goType: "uint64", storeIdent: "rel_module", ident: "moduleID" }
			sync_status: { goType: "string", storeIdent: "sync_status" }
			sync_type:   { goType: "string", storeIdent: "sync_type"   }
		}

		byValue: ["node_id", "module_id", "sync_status", "sync_type"]
	}

	store: {
		ident: "federationNodeSync"

		api: {
			lookups: [
				{
					fields: ["node_id"]
					description: """
						searches for sync activity by node ID

						It returns sync activity
						"""
				}, {
					fields: ["node_id", "module_id", "sync_type", "sync_status"]
					description: """
						searches for activity by node, type and status

						It returns sync activity
						"""
				}
			]
		}
	}
}
