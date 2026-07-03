package federation

import (
	"github.com/crusttech/human/server/codegen/schema"
)

_moduleMappingDefs: {
			ModuleFieldMappingSet: { name: "ModuleFieldMappingSet", elem: _moduleMappingDefs.ModuleFieldMapping, elemPtr: true }
			ModuleFieldMapping: { name: "ModuleFieldMapping", fields: [
				{ name: "Origin", type: _moduleMappingDefs.ModuleField, json: "origin" },
				{ name: "Destination", type: _moduleMappingDefs.ModuleField, json: "destination" },
			]}
			ModuleField: { name: "ModuleField", fields: [
				{ name: "Kind", type: "string", json: "kind" },
				{ name: "Name", type: "string", json: "name" },
				{ name: "Label", type: "string", json: "label" },
				{ name: "IsMulti", type: "bool", json: "isMulti" },
			]}
		}

moduleMapping: {
	parents: [
		{handle: "node"},
	]

	features: {
		labels: false
	}

	types: {
		gen: true
		defs: _moduleMappingDefs
	}

	model: {
		ident: "federation_module_mapping"
		attributes: {
			tenant_id:  schema.TenantRefField
			project_id: schema.ProjectRefField
			node_id: {
				ident:   "nodeID"
				unique:  true
				goType:  "uint64"
				dal: { type: "ID" }
				json: { field: "nodeID", string: true }
			}
			federation_module_id: {
				sortable: true,
				ident: "federationModuleID",
				goType: "uint64"
				storeIdent: "rel_federation_module"
				dal: { type: "ID" }
				json: { field: "federationModuleID", string: true }
			}
			compose_module_id: {
				sortable: true,
				ident: "composeModuleID",
				goType: "uint64"
				storeIdent: "rel_compose_module"
				dal: { type: "ID" }
				json: { field: "composeModuleID", string: true }
			}
			compose_namespace_id: {
				sortable: true,
				ident: "composeNamespaceID",
				goType: "uint64"
				storeIdent: "rel_compose_namespace"
				dal: { type: "ID" }
				json: { field: "composeNamespaceID", string: true }
			}
			field_mapping: {
				type: _moduleMappingDefs.ModuleFieldMappingSet
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
				json: { field: "fields" }
			}
		}

		indexes: {
			"unique_module_compose_module": {
				attributes: ["federation_module_id", "compose_module_id", "compose_namespace_id" ]
			}
		}
	}

	service: {
		delete:   false
		undelete: false

		customBodyOps: ["lookup", "search", "create", "update"]

		customAccessOps: ["search", "create"]

		actionProp: "mapping"
		createProp: "created"
		updateProp: "changed"
	}

	envoy: {
		omit: true
	}

	filter: {
		struct: {
			tenant_id: schema.TenantFilterField
			project_id: schema.ProjectFilterField
			compose_module_id:    { goType: "uint64", ident: "composeModuleID", storeIdent: "rel_compose_module" }
			compose_namespace_id: { goType: "uint64", ident: "composeNamespaceID", storeIdent: "rel_compose_namespace" }
			federation_module_id: { goType: "uint64", ident: "federationModuleID", storeIdent: "rel_federation_module" }
		}

		byValue: ["compose_module_id", "compose_namespace_id", "federation_module_id"]
	}

	store: {
		ident: "federationModuleMapping"

		api: {
			lookups: [
				{
					fields: ["federation_module_id", "compose_module_id", "compose_namespace_id"]
					description: """
						searches for module mapping by federation module id and compose module id

						It returns module mapping
						"""
				}, {
					fields: ["federation_module_id"]
					description: """
						searches for module mapping by federation module id

						It returns module mapping
						"""
				}
			]
		}
	}
}
