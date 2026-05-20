package system

import (
	"github.com/crusttech/human/server/codegen/schema"
)

chatbot: {
	features: {
		labels: true
	}

	model: {
		attributes: {
			id:     schema.IdField
			handle: schema.HandleField
			name: {
				sortable: true,
				goType: "string"
				dal: { length: 512 }
			}
			enabled: {
				goType: "bool"
				dal: { type: "Boolean" }
			}
			widget_key: {
				goType: "string"
				storeIdent: "widget_key"
				dal: { length: 128 }
			}
			allowed_origins: {
				goType: "types.ChatbotAllowedOrigins"
				storeIdent: "allowed_origins"
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
			}
			session_ttl: {
				goType: "string"
				ident: "sessionTTL"
				storeIdent: "session_ttl"
				dal: { length: 32 }
			}
			handoff: {
				goType: "types.ChatbotHandoff"
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
			}
			styling: {
				goType: "types.ChatbotStyling"
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
			}
			scenarios: {
				goType: "types.ChatbotScenarios"
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
			"unique_handle": {
				fields: [{ attribute: "handle", modifiers: ["LOWERCASE"] }]
				predicate: "handle != '' AND deleted_at IS NULL"
			}
			"unique_widget_key": {
				fields: [{ attribute: "widget_key" }]
				predicate: "widget_key != '' AND deleted_at IS NULL"
			}
		}
	}

	filter: {
		struct: {
			chatbot_id: {goType: "[]uint64", ident: "chatbotID", storeIdent: "id"}
			handle: {goType: "string"}
			widget_key: {goType: "string", ident: "widgetKey", storeIdent: "widget_key"}
			deleted: {goType: "filter.State", storeIdent: "deleted_at"}
		}

		query: ["handle", "name"]
		byValue: ["chatbot_id", "handle", "widget_key"]
		byNilState: ["deleted"]
	}

	rbac: {
		operations: {
			read: description:                      "Read chatbot"
			update: description:                    "Update chatbot"
			delete: description:                    "Delete chatbot"
			"sessions.view": description:           "View chatbot sessions"
			"sessions.manage": description:         "Manage chatbot sessions (advance, close)"
			"sessions.handoff.manage": description: "Manage chatbot session handoffs"
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
						searches for chatbot by ID

						It also returns deleted chatbots.
						"""
				}, {
					fields: ["handle"]
					nullConstraint: ["deleted_at"]
					constraintCheck: true
					description: """
						searches for chatbot by handle

						It returns only valid chatbots (not deleted)
						"""
				}, {
					fields: ["widget_key"]
					nullConstraint: ["deleted_at"]
					constraintCheck: true
					description: """
						searches for chatbot by widget key

						It returns only valid chatbots (not deleted)
						"""
				},
			]
		}
	}
}
