package system

import (
	"github.com/crusttech/human/server/codegen/schema"
)

_llm_providerDefs: {
			LLMProviderMeta: { name: "LLMProviderMeta", fields: [
				{ name: "Short", type: "string", json: "short" },
				{ name: "Description", type: "string", json: "description" },
			]}
			LLMProviderConfig: { name: "LLMProviderConfig", fields: [
				{ name: "PromptURL", type: "string", json: "promptURL" },
				{ name: "Model", type: "string", json: "model" },
				{ name: "Temperature", type: "float64", ptr: true, json: "temperature,omitempty" },
				{ name: "Timeout", type: "string", json: "timeout" },
				{ name: "Guard", type: _llm_providerDefs.LLMProviderGuardConfig, ptr: true, json: "guard,omitempty" },
			]}
			LLMProviderGuardConfig: { name: "LLMProviderGuardConfig", fields: [
				{ name: "Enabled", type: "bool", json: "enabled" },
				{ name: "Provider", type: "string", json: "provider" },
				{ name: "Model", type: "string", json: "model" },
				{ name: "Thresholds", goType: "map[string]float64", json: "thresholds,omitempty" },
			]}
		}

llm_provider: {
	model: {
		attributes: {
			id:     schema.IdField
			tenant_id:  schema.TenantRefField
			project_id: schema.ProjectRefField
			handle: schema.HandleField
			status: {
				sortable: true
				dal: { type: "Text", length: 64 }
			}
			provider: {
				sortable: true
				dal: { type: "Text", length: 128 }
			}
			credential_id: {
				goType: "uint64"
				ident: "credentialID"
				storeIdent: "rel_credential"
				json: "credentialID,string"
				dal: { type: "Ref", refModelResType: "corteza::system:credential", default: 0 }
			}
			meta: {
				type: _llm_providerDefs.LLMProviderMeta
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
			}
			// Virtual sortable mapped onto meta->>'short'.
			name: {
				sortableJSON: { json: "meta.short", accessor: "Meta.Short", nullable: false }
				store: false
				goType: "string"
				omitSetter: true
				omitGetter: true
				envoy: { yaml: { omitEncoder: true } }
			}
			config: {
				type: _llm_providerDefs.LLMProviderConfig
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
			}

			created_at: schema.SortableTimestampNowField
			updated_at: schema.SortableTimestampNilField
			deleted_at: schema.SortableTimestampNilField
			created_by: schema.AttributeUserRef & { json: "createdBy,string" }
			updated_by: schema.AttributeUserRef
			deleted_by: schema.AttributeUserRef
		}

		indexes: {
			"primary": { attribute: "id" }
		}
	}

	filter: {
		struct: {
			llm_provider_id: {goType: "[]uint64", ident: "llmProviderID", storeIdent: "id"}
			tenant_id: schema.TenantFilterField
			project_id: schema.ProjectFilterField
			handle: {goType: "string"}
			status: {goType: "string"}
			provider: {goType: "string"}
			deleted: {goType: "filter.State", storeIdent: "deleted_at"}
		}

		byValue: ["llm_provider_id", "handle", "status", "provider"]
		byNilState: ["deleted"]
	}

	rbac: {
		operations: {
			read: description:   "Read LLM provider"
			update: description: "Update LLM provider"
			delete: description: "Delete LLM provider"
		}
	}

	features: {
		labels: false
	}

	types: {
		gen: true
		defs: _llm_providerDefs
	}

	envoy: {
		omit: true
	}

	store: {
		api: {
			lookups: [
				{
					fields: ["id"]
					description: """
						searches for LLM provider by ID

						It returns LLM provider even if deleted
						"""
				}, {
					fields: ["handle"]
					nullConstraint: ["deleted_at"]
					constraintCheck: true
					description: """
						searches for LLM provider by handle

						It returns only valid LLM provider (not deleted)
						"""
				},
			]
		}
	}
}
