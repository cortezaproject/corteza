package federation

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

nodeSync: schema.#Resource & {
	features: {
		labels: false
	}

	struct: {
		node_id: { ident: "nodeID", goType: "uint64", primaryKey: true }
		module_id: { ident: "moduleID", goType: "uint64" }
		sync_type: { goType: "string" }
		sync_status: { goType: "string" }
	} & {
		time_of_action: schema.SortableTimestampField
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

		settings: {
			rdbms: {
				table: "federation_nodes_sync"
			}
		}

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
