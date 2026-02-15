package system

import (
	"github.com/cortezaproject/corteza/server/codegen/schema"
)

llm_provider: {
	model: {
		attributes: {
			id:     schema.IdField
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
				dal: { type: "Ref", refModelResType: "corteza::system:credential", default: 0 }
			}
			meta: {
				goType: "types.LLMProviderMeta"
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
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
			llm_provider_id: {goType: "[]uint64", ident: "llmProviderID", storeIdent: "id"}
			handle: {goType: "string"}
			status: {goType: "string"}
			provider: {goType: "string"}
			deleted: {goType: "filter.State", storeIdent: "deleted_at"}
		}

		byValue: ["llm_provider_id", "handle", "status", "provider"]
		byNilState: ["deleted"]
	}

	features: {
		labels: false
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
