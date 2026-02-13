package system

import (
	"github.com/cortezaproject/corteza/server/codegen/schema"
)

ai_conversation: {
	features: {
		labels: false
	}

	model: {
		ident: "ai_conversations"
		attributes: {
			id: schema.IdField
			agentID: {
				sortable: true,
				goType: "uint64"
				storeIdent: "rel_agent"
				dal: { type: "Ref", refModelResType: "corteza::system:agent" }
			}
			messages: {
				goType: "types.AiConversationMessages"
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
			}
			tokenCount: {
				goType: "int"
				storeIdent: "token_count"
				dal: { type: "Number", meta: { "rdbms:type": "integer" } }
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
			"idx_agent": {
				attribute: "agentID"
			}

		}
	}

	filter: {
		struct: {
			ai_conversation_id: {goType: "[]uint64", ident: "aiConversationID", storeIdent: "id"}
			agent_id: {goType: "uint64", ident: "agentID", storeIdent: "rel_agent"}
			deleted: {goType: "filter.State", storeIdent: "deleted_at"}
		}

		byValue: ["ai_conversation_id", "agent_id"]
		byNilState: ["deleted"]
	}

	rbac: {
		operations: {
			read: description:   "Read AI conversation"
			update: description: "Update AI conversation"
			delete: description: "Delete AI conversation"
		}
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
						searches for AI conversation by ID

						It also returns deleted conversations.
						"""
				},
			]
		}
	}
}
