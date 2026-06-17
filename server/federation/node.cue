package federation

import (
	"github.com/crusttech/human/server/codegen/schema"
)

node: {
	features: {
		labels: false
		paging: true
		sorting: true
		tenantScoped: true
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
		// Node CRUD is heavily bespoke: every mutating op funnels through the
		// hand-written svc.updater (status=Failed-on-error bookkeeping, UpdatedBy/
		// DeletedBy stamping) and uses CanManageNode rather than the unified
		// CanRead/Update/Delete checks, so the standard generated bodies can't be
		// used. The generator owns only the recordAction scaffold + the search/
		// create access checks; the bodies are delegated to on<Op> handlers.
		//
		// lookup is disabled: FindByID (and the controller-facing Read alias) are
		// kept fully custom -- they record no action / a different action and use
		// CanManageNode.
		lookup:   false
		undelete: true

		// action-log prop is named "node" (no dedicated "new"/"update" props exist
		// in node_actions.yaml), so omit the create prop and point updateProp at it.
		actionProp:     "node"
		omitCreateProp: true
		updateProp:     "node"

		customBodyOps: ["search", "create", "update", "delete", "undelete"]
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
