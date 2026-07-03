package system

import (
	"github.com/crusttech/human/server/codegen/schema"
)

_ai_conversationDefs: {
	AiConversationMessages: { name: "AiConversationMessages", elem: _ai_conversationDefs.AiConversationMessage }
	AiConversationMessage: { name: "AiConversationMessage", fields: [
		{ name: "Role", type: "string", json: "role" },
		{ name: "Content", type: "string", json: "content" },
		{ name: "Operator", type: "string", json: "operator,omitempty" },
		{ name: "ToolCalls", slice: true, type: _ai_conversationDefs.AiConversationToolCall, json: "toolCalls,omitempty" },
		{ name: "ToolResults", slice: true, type: _ai_conversationDefs.AiConversationToolResult, json: "toolResults,omitempty" },
	]}
	AiConversationToolCall: { name: "AiConversationToolCall", fields: [
		{ name: "CallID", type: "string", json: "callID" },
		{ name: "Name", type: "string", json: "name" },
		{ name: "Data", type: "string", json: "data" },
	]}
	AiConversationToolResult: { name: "AiConversationToolResult", fields: [
		{ name: "CallID", type: "string", json: "callID" },
		{ name: "Data", type: "string", json: "data" },
		{ name: "Error", type: "string", json: "error,omitempty" },
	]}
}

ai_conversation: {
	features: {
		labels: false
	}

	types: {
		gen: true
		defs: _ai_conversationDefs
	}

	model: {
		ident: "ai_conversations"
		attributes: {
			id: schema.IdField
			tenant_id:  schema.TenantRefField
			project_id: schema.ProjectRefField
			agentID: {
				sortable: true,
				goType: "uint64"
				storeIdent: "rel_agent"
				dal: { type: "Ref", refModelResType: "corteza::system:agent" }
				json: {field: "agentID", string: true}
			}
			messages: {
				type: _ai_conversationDefs.AiConversationMessages
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
			created_by: schema.AttributeUserRef & { json: {field: "createdBy", string: true} }
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
			tenant_id: schema.TenantFilterField
			project_id: schema.ProjectFilterField
			agent_id: {goType: "uint64", ident: "agentID", storeIdent: "rel_agent"}
			deleted: {goType: "filter.State", storeIdent: "deleted_at"}
		}

		byValue: ["ai_conversation_id", "agent_id"]
		byNilState: ["deleted"]
	}

	service: {
		genAccessController: true
		genConstructor:      true

		undelete: true

		customBodyOps: ["lookup", "search", "create", "update", "delete", "undelete"]

		customAccessOps: ["search"]
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
