package federation

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

node: {
	features: {
		labels: false
		paging: false
		sorting: false
	}

	model: {
		ident: "federation_nodes"
		attributes: {
				id:          schema.IdField
				name: {sortable: true}
				shared_node_id: { sortable: true, ident: "sharedNodeID", goType: "uint64" }
				base_url: { sortable: true, goType: "string", ident: "baseURL" }
				status: { sortable: true, goType: "string" }
				contact: { sortable: true, goType: "string" }
				pair_token: { goType: "string" }
				auth_token: { goType: "string" }

				created_at: schema.SortableTimestampNowField
				updated_at: schema.SortableTimestampNilField
				deleted_at: schema.SortableTimestampNilField
				created_by: { goType: "uint64" }
				updated_by: { goType: "uint64" }
				deleted_by: { goType: "uint64" }
		}
	}

	filter: {
		struct: {
			name: { goType: "string" }
			base_url: { goType: "string", ident: "baseURL" }
			status: { goType: "string" }
			deleted: { goType: "filter.State", storeIdent: "deleted_at" }
		}

		query: ["name", "base_url"]
		byNilState: ["deleted"]
	}

	rbac: {
		operations: {
			"manage": description:        "Manage federation node"
			"module.create": description: "Create shared module"
		}
	}

	store: {
		ident: "federationNode"

		api: {
			lookups: [
				{
					fields: ["id"]
					description: """
						searches for federation node by ID

						It returns federation node
						"""
				}, {
					fields: ["base_url", "shared_node_id"]
					description: """
						searches for node by shared-node-id and base-url
						"""
				}, {
					fields: ["shared_node_id"]
					description: """
						searches for node by shared-node-id
						"""
				}
			]
		}
	}
}
