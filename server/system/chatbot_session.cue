package system

import (
	"github.com/crusttech/human/server/codegen/schema"
)

chatbot_session: {
	features: {
		labels: false
		projectScoped: true
	}

	types: {
		gen: true
	}

	model: {
		ident: "chatbot_sessions"
		attributes: {
			id: schema.IdField & { json: "id,string" }
			tenant_id:  schema.TenantRefField
			project_id: schema.ProjectRefField
			chatbot_id: {
				sortable: true,
				goType: "uint64"
				ident: "chatbotID"
				storeIdent: "rel_chatbot"
				dal: { type: "Ref", refModelResType: "corteza::system:chatbot" }
				json: "chatbotID,string"
			}
			status: {
				sortable: true,
				goType: "string"
				dal: { length: 32 }
			}
			current_step: {
				goType: "int"
				ident: "currentStep"
				storeIdent: "current_step"
				dal: { type: "Number", meta: { "rdbms:type": "integer" } }
			}
			state: {
				goType: "types.ChatbotSessionState"
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
				json: "state,omitempty"
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
			"idx_chatbot": {
				attribute: "chatbot_id"
			}
			"idx_status": {
				attribute: "status"
			}
		}
	}

	filter: {
		struct: {
			chatbot_session_id: {goType: "[]uint64", ident: "chatbotSessionID", storeIdent: "id"}
			tenant_id: schema.TenantFilterField
			project_id: schema.ProjectFilterField
			chatbot_id: {goType: "uint64", ident: "chatbotID", storeIdent: "rel_chatbot"}
			status: {goType: "[]string"}
			deleted: {goType: "filter.State", storeIdent: "deleted_at"}
		}

		query: ["status"]
		byValue: ["chatbot_session_id", "chatbot_id", "status"]
		byNilState: ["deleted"]
	}

	rbac: {
		operations: {
			read: description:   "Read chatbot session"
			update: description: "Update chatbot session"
			delete: description: "Delete chatbot session"
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
						searches for chatbot session by ID

						It also returns deleted sessions.
						"""
				}, {
					fields: ["chatbot_id"]
					description: """
						searches for chatbot sessions by chatbot ID
						"""
				},
			]
		}
	}
}

chatbot_session_step: {
	features: {
		labels: false
	}

	model: {
		ident: "chatbot_session_steps"
		attributes: {
			id: schema.IdField
			session_id: {
				sortable: true,
				goType: "uint64"
				ident: "sessionID"
				storeIdent: "rel_session"
				dal: { type: "Ref", refModelResType: "corteza::system:chatbot-session" }
			}
			scenario_index: {
				goType: "int"
				ident: "scenarioIndex"
				storeIdent: "scenario_index"
				dal: { type: "Number", meta: { "rdbms:type": "integer" } }
			}
			conversation_id: {
				sortable: true,
				goType: "uint64"
				ident: "conversationID"
				storeIdent: "rel_conversation"
				dal: { type: "Ref", refModelResType: "corteza::system:ai-conversation" }
			}
			status: {
				sortable: true,
				goType: "string"
				dal: { length: 32 }
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
			"idx_session": {
				attribute: "session_id"
			}
			"idx_conversation": {
				attribute: "conversation_id"
			}
		}
	}

	filter: {
		struct: {
			chatbot_session_step_id: {goType: "[]uint64", ident: "chatbotSessionStepID", storeIdent: "id"}
			session_id: {goType: "uint64", ident: "sessionID", storeIdent: "rel_session"}
			conversation_id: {goType: "uint64", ident: "conversationID", storeIdent: "rel_conversation"}
			status: {goType: "string"}
			deleted: {goType: "filter.State", storeIdent: "deleted_at"}
		}

		query: ["status"]
		byValue: ["chatbot_session_step_id", "session_id", "conversation_id", "status"]
		byNilState: ["deleted"]
	}

	rbac: {
		operations: {
			read: description:   "Read chatbot session step"
			update: description: "Update chatbot session step"
			delete: description: "Delete chatbot session step"
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
						searches for chatbot session step by ID

						It also returns deleted steps.
						"""
				}, {
					fields: ["session_id"]
					description: """
						searches for chatbot session steps by session ID
						"""
				}, {
					fields: ["conversation_id"]
					description: """
						searches for chatbot session step by conversation ID
						"""
				},
			]
		}
	}
}

chatbot_session_handoff: {
	features: {
		labels: false
	}

	model: {
		ident: "chatbot_session_handoffs"
		attributes: {
			id: schema.IdField
			session_id: {
				sortable: true,
				goType: "uint64"
				ident: "sessionID"
				storeIdent: "rel_session"
				dal: { type: "Ref", refModelResType: "corteza::system:chatbot-session" }
			}
			step_id: {
				sortable: true,
				goType: "uint64"
				ident: "stepID"
				storeIdent: "rel_step"
				dal: { type: "Ref", refModelResType: "corteza::system:chatbot-session-step" }
			}
			status: {
				sortable: true,
				goType: "string"
				dal: { length: 32 }
			}
			initiated_at: {
				sortable: true,
				goType: "time.Time"
				ident: "initiatedAt"
				storeIdent: "initiated_at"
				dal: { type: "Timestamp", nullable: false }
			}
			closed_at: {
				sortable: true,
				goType: "*time.Time"
				ident: "closedAt"
				storeIdent: "closed_at"
				dal: { type: "Timestamp", nullable: true }
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
			"idx_session": {
				attribute: "session_id"
			}
			"idx_status": {
				attribute: "status"
			}
		}
	}

	filter: {
		struct: {
			chatbot_session_handoff_id: {goType: "[]uint64", ident: "chatbotSessionHandoffID", storeIdent: "id"}
			session_id: {goType: "uint64", ident: "sessionID", storeIdent: "rel_session"}
			step_id: {goType: "uint64", ident: "stepID", storeIdent: "rel_step"}
			status: {goType: "string"}
			deleted: {goType: "filter.State", storeIdent: "deleted_at"}
		}

		query: ["status"]
		byValue: ["chatbot_session_handoff_id", "session_id", "status"]
		byNilState: ["deleted"]
	}

	rbac: {
		operations: {
			read: description:   "Read chatbot session handoff"
			update: description: "Update chatbot session handoff"
			delete: description: "Delete chatbot session handoff"
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
						searches for chatbot session handoff by ID

						It also returns deleted handoffs.
						"""
				}, {
					fields: ["session_id"]
					description: """
						searches for chatbot session handoff by session ID
						"""
				}, {
					fields: ["status"]
					description: """
						searches for chatbot session handoffs by status
						"""
				},
			]
		}
	}
}
