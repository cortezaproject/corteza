package federation

import (
	"github.com/crusttech/human/server/codegen/schema"
)

_sharedModuleDefs: {
	ModuleFieldSet: { name: "ModuleFieldSet", elem: _exposedModuleDefs.ModuleField, elemPtr: true }
}

sharedModule: {
	features: {
		labels: false
	}

	types: {
		gen: true
		defs: _sharedModuleDefs
		// ModuleFieldSet's Scan/Value/Parse are generated via exposed_module (same pkg)
		jsonTypesSkip: ["ModuleFieldSet"]
	}

	parents: [
		{handle: "node"},
	]

	model: {
		ident: "federation_module_shared"
		attributes: {
			id: schema.IdField & {
				// hand-written tag is `moduleID,string` (resource ident is
				// `sharedModule`, so the convention would emit `sharedModuleID`).
				json: { field: "moduleID", string: true }
			}
			tenant_id:  schema.TenantRefField
			project_id: schema.ProjectRefField
			handle: schema.HandleField
			node_id: {
				sortable: true,
				ident: "nodeID",
				goType: "uint64",
				storeIdent: "rel_node"
				dal: { type: "ID" }
				// parent ref: hand-written tag has no omitempty.
				json: { field: "nodeID", string: true }
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
				// hand-written tag has no omitempty.
				json: { field: "externalFederationModuleID", string: true }
			}
			fields: {
				type: _sharedModuleDefs.ModuleFieldSet
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
			}

			created_at: schema.SortableTimestampNowField
			updated_at: schema.SortableTimestampNilField
			deleted_at: schema.SortableTimestampNilField
			created_by: schema.AttributeUserRef & {
				// hand-written tag has no omitempty (unlike updated_by/deleted_by).
				json: { field: "createdBy", string: true }
			}
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

	filter: {
		struct: {
			tenant_id: schema.TenantFilterField
			project_id: schema.ProjectFilterField
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

	service: {
		scoped: true

		delete:   false
		undelete: false

		customBodyOps: ["lookup", "search", "create", "update"]

		customAccessOps: ["search", "create"]

		actionProp: "module"
		createProp: "changed"
		updateProp: "module"
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
