package federation

import (
	"github.com/crusttech/human/server/codegen/schema"
)

node: {
	features: {
		labels: false
		paging: true
		sorting: true
	}

	types: {
		gen: true
	}

	model: {
		ident: "federation_nodes"
		attributes: {
      id: schema.IdField
      tenant_id: schema.TenantRefField
      shared_node_id: {
      	sortable: true,
      	ident: "sharedNodeID",
      	goType: "uint64"
				dal: { type: "ID" }
				json: { field: "sharedNodeID", string: true }
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
				json: "-"
			}
      auth_token: {
      	dal: {}
				json: "-"
			}
      created_at: schema.SortableTimestampNowField
      updated_at: schema.SortableTimestampNilField
      deleted_at: schema.SortableTimestampNilField
			created_by: schema.AttributeUserRef & { json: { field: "createdBy", string: true } }
			updated_by: schema.AttributeUserRef
			deleted_by: schema.AttributeUserRef
		}

		indexes: {
			"primary": { attribute: "id" }
		}
	}

	envoy: {
		omit: true
	}

	service: {
		lookup:   false
		undelete: true

		// status and the pairing tokens are managed by the pairing/handshake
		// flow, not by a plain node update; excluding them keeps a metadata
		// update (name/baseURL/contact) from wiping the node's auth to zero.
		omitUpdateFields: ["status", "pair_token", "auth_token"]

		actionProp:     "node"
		omitCreateProp: true
		updateProp:     "node"

		customBodyOps: ["search", "create", "update", "delete", "undelete"]

		customFunctions: [
			{
				name:   "Read"
				cap:    "read"
				action: "Create"
				args: [{name: "ID", goType: "uint64"}]
				results: [{name: "res", goType: "*types.Node"}, {name: "err", goType: "error"}]
			},
			{
				name:   "CreateFromPairingURI"
				cap:    "write"
				ac:     "CanPair"
				acErr:  "ErrNotAllowedToPair"
				action: "CreateFromPairingURI"
				args: [{name: "uri", goType: "string"}]
				results: [{name: "n", goType: "*types.Node"}, {name: "err", goType: "error"}]
			},
			{
				name: "RegenerateNodeURI"
				cap:  "write"
				args: [{name: "nodeID", goType: "uint64"}]
				results: [{name: "uri", goType: "string"}, {name: "err", goType: "error"}]
			},
			{
				name:  "Pair"
				cap:   "write"
				ac:    "CanPair"
				acErr: "ErrNotAllowedToPair"
				args: [{name: "nodeID", goType: "uint64"}]
				results: [{name: "err", goType: "error"}]
			},
			{
				name: "HandshakeInit"
				cap:  "write"
				args: [{name: "nodeID", goType: "uint64"}, {name: "pairToken", goType: "string"}, {name: "sharedNodeID", goType: "uint64"}, {name: "authToken", goType: "string"}]
				results: [{name: "err", goType: "error"}]
			},
			{
				name:  "HandshakeConfirm"
				cap:   "write"
				ac:    "CanPair"
				acErr: "ErrNotAllowedToPair"
				args: [{name: "nodeID", goType: "uint64"}]
				results: [{name: "err", goType: "error"}]
			},
			{
				name:  "HandshakeComplete"
				cap:   "write"
				ac:    "CanPair"
				acErr: "ErrNotAllowedToPair"
				args: [{name: "sharedNodeID", goType: "uint64"}, {name: "token", goType: "string"}]
				results: [{name: "err", goType: "error"}]
			},
		]
	}

	filter: {
		struct: {
			tenant_id: schema.TenantFilterField
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
