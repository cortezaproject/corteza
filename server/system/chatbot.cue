package system

import (
	"github.com/crusttech/human/server/codegen/schema"
)

chatbot: {
	features: {
		labels: true
		projectScoped: true
	}

	types: {
		gen: true
		// ChatbotStyling has a custom Value (kept hand-written)
		jsonTypesSkip: ["ChatbotStyling"]
	}

	model: {
		attributes: {
			id:     schema.IdField
			tenant_id:  schema.TenantRefField
			project_id: schema.ProjectRefField
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
				json: { field: "widgetKey", omitEmpty: true }
			}
			allowed_origins: {
				goType: "types.ChatbotAllowedOrigins"
				storeIdent: "allowed_origins"
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
				json: { field: "allowedOrigins", omitEmpty: true }
			}
			session_ttl: {
				goType: "string"
				ident: "sessionTTL"
				storeIdent: "session_ttl"
				dal: { length: 32 }
				json: { field: "sessionTTL", omitEmpty: true }
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
			created_by: schema.AttributeUserRef & { json: { field: "createdBy", string: true } }
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
			tenant_id: schema.TenantFilterField
			project_id: schema.ProjectFilterField
			handle: {goType: "string"}
			widget_key: {goType: "string", ident: "widgetKey", storeIdent: "widget_key"}
			deleted: {goType: "filter.State", storeIdent: "deleted_at"}
		}

		query: ["handle", "name"]
		byValue: ["chatbot_id", "project_id", "handle", "widget_key"]
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

	refs: {
		extended: true
	}

	service: {
		genAccessController: true
		genConstructor:      true

		undelete: true

		customBodyOps: ["lookup", "search", "create", "update", "delete", "undelete"]

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
