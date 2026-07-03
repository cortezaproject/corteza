package federation

import (
	"github.com/crusttech/human/server/codegen/schema"
)

nodeSync: {
	features: {
		labels: false
	}

	types: {
		gen: true
	}

	model: {
		ident: "federation_nodes_sync"
		attributes: {
			rel_node: {
			  sortable: true,
				ident: "nodeID",
				storeIdent: "rel_node",
				goType: "uint64",
				dal: { type: "ID" }
				// hand-written tag omits the conventional ",omitempty" for refs
				json: { field: "nodeID", string: true }
			}
			rel_module: {
				sortable: true,
				ident: "moduleID",
				storeIdent: "rel_compose_module",
				goType: "uint64"
				dal: { type: "ID" }
				// hand-written tag omits the conventional ",omitempty" for refs
				json: { field: "moduleID", string: true }
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
			time_of_action: schema.SortableTimestampField & {
				// hand-written tag omits the conventional ",omitempty" for timestamps
				json: "timeOfAction"
			}
		}

		indexes: {
			"idx_rel_node": { attribute: "rel_node" }
		}
	}

	filter: {
		struct: {
			rel_node:     { goType: "uint64", storeIdent: "rel_node",   ident: "nodeID" }
			rel_module:   { goType: "uint64", storeIdent: "rel_module", ident: "moduleID" }
			sync_status: { goType: "string", storeIdent: "sync_status" }
			sync_type:   { goType: "string", storeIdent: "sync_type"   }
		}

		byValue: ["rel_node", "rel_module", "sync_status", "sync_type"]
	}

	envoy: {
		omit: true
	}

	service: {
		lookup:   false
		update:   false
		delete:   false
		undelete: false

		customBodyOps:   ["create", "search"]
		customAccessOps: ["create", "search"]

		customFunctions: []

		omitCreateProp: true
		filterProp:     "nodeSyncFilter"
	}

	store: {
		ident: "federationNodeSync"

		api: {
			lookups: [
				{
					fields: ["rel_node"]
					description: """
						searches for sync activity by node ID

						It returns sync activity
						"""
				}, {
					fields: ["rel_node", "rel_module", "sync_type", "sync_status"]
					description: """
						searches for activity by node, type and status

						It returns sync activity
						"""
				}
			]
		}
	}
}
