package federation

moduleMapping: {
	parents: [
		{handle: "node"},
	]

	features: {
		labels: false
	}

	model: {
		ident: "federation_module_mapping"
		attributes: {
			node_id: {
				ident:   "nodeID"
				unique:  true
				goType:  "uint64"
				dal: { type: "ID" }
			}
			federation_module_id: {
				sortable: true,
				ident: "federationModuleID",
				goType: "uint64"
				dal: { type: "ID" }
			}
			compose_module_id: {
				sortable: true,
				ident: "composeModuleID",
				goType: "uint64"
				dal: { type: "ID" }
			}
			compose_namespace_id: {
				sortable: true,
				ident: "composeNamespaceID",
				goType: "uint64"
				dal: { type: "ID" }
			}
			field_mapping: {
				goType: "types.ModuleFieldMappingSet"
				dal: { type: "JSON", defaultEmptyObject: true }
			}
		}

		indexes: {
			"primary": { attribute: "node_id" }
			"unique_module_compose_module": {
				attributes: ["federation_module_id", "compose_module_id", "compose_namespace_id" ]
			}
		}
	}

	filter: {
		struct: {
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
