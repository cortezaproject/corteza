package federation

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

node: schema.#Resource & {
	features: {
		labels: false
		paging: false
		sorting: false
	}

	struct: {
		id:          schema.IdField
		name: {}
		shared_node_id: { ident: "sharedNodeID", goType: "uint64" }
		base_url: { goType: "string", ident: "baseURL" }
		status: { goType: "string" }
		contact: { goType: "string" }
		pair_token: { goType: "string" }
		auth_token: { goType: "string" }

		created_at: schema.SortableTimestampField
		updated_at: schema.SortableTimestampNilField
		deleted_at: schema.SortableTimestampNilField
		created_by: { goType: "uint64" }
		updated_by: { goType: "uint64" }
		deleted_by: { goType: "uint64" }
	}

	filter: {
		struct: {
			name: { goType: "string" }
			base_url: { goType: "string", ident: "baseURL" }
			status: { goType: "string" }
			deleted: { goType: "filter.State", storeIdent: "deleted_at" }
		}

		query: ["name", "base_url"]
		byQuery: ["status"]
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

		settings: {
			rdbms: {
				table: "federation_nodes"
			}
		}

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
