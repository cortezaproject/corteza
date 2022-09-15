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
      id: schema.IdField
      shared_node_id: {
      	sortable: true,
      	ident: "sharedNodeID",
      	goType: "uint64"
				dal: { type: "ID" }
			}
      name: {
      	sortable: true
      	dal: {}
			}
      base_url: {
      	sortable: true,
      	ident: "baseURL"
      	dal: {}
			}
      status: {
      	sortable: true,
      	dal: {}
			}
      contact: {
      	sortable: true,
      	dal: {}
			}
      pair_token: {
      	dal: {}
			}
      auth_token: {
      	dal: {}
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
