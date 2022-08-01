package federation

moduleMapping: {
	parents: [
		{handle: "node"},
	]

	features: {
		labels: false
	}

	model: {
		node_id: { ident: "nodeID", goType: "uint64", primaryKey: true, unique: true }
		federation_module_id: { sortable: true, ident: "federationModuleID", goType: "uint64" }
		compose_module_id: { sortable: true, ident: "composeModuleID", goType: "uint64" }
		compose_namespace_id: { sortable: true, ident: "composeNamespaceID", goType: "uint64" }
		field_mapping: { goType: "types.ModuleFieldMappingSet" }
	}

	filter: {
		model: {
			compose_module_id:    { goType: "uint64", ident: "composeModuleID", storeIdent: "rel_compose_module" }
			compose_namespace_id: { goType: "uint64", ident: "composeNamespaceID", storeIdent: "rel_compose_namespace" }
			federation_module_id: { goType: "uint64", ident: "federationModuleID", storeIdent: "rel_federation_module" }
		}

		byValue: ["compose_module_id", "compose_namespace_id", "federation_module_id"]
	}

	store: {
		ident: "federationModuleMapping"

		settings: {
			rdbms: {
				table: "federation_module_mapping"
			}
		}

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
