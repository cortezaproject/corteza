package system

import (
	"github.com/crusttech/human/server/codegen/schema"
)

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
				goType: "types.LLMProviderMeta"
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
				goType: "types.LLMProviderConfig"
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
		projectScoped: true
	}

	types: {
		gen: true
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
