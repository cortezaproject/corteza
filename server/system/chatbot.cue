package system

import (
	"github.com/crusttech/human/server/codegen/schema"
)

_chatbotDefs: {
			ChatbotHandoff: { name: "ChatbotHandoff", fields: [
				{ name: "Enabled", type: "bool", json: "enabled" },
				{ name: "Automation", type: _chatbotDefs.ChatbotHandoffAutomation, json: "automation" },
			]}
			ChatbotHandoffAutomation: { name: "ChatbotHandoffAutomation", fields: [
				{ name: "OnRequested", type: _chatbotDefs.ChatbotAutomationHook, json: "onRequested" },
				{ name: "OnAccepted", type: _chatbotDefs.ChatbotAutomationHook, json: "onAccepted" },
			]}
			ChatbotAutomationHook: { name: "ChatbotAutomationHook", fields: [
				{ name: "Automation", type: "string", json: "automation" },
				{ name: "Async", type: "bool", json: "async" },
				{ name: "Mappings", slice: true, type: _chatbotDefs.ChatbotHookMapping, json: "mappings,omitempty" },
			]}
			ChatbotHookMapping: { name: "ChatbotHookMapping", fields: [
				{ name: "StateExpression", type: _chatbotDefs.ChatbotStateExpression, json: "stateExpression" },
				{ name: "TriggerParam", type: "string", json: "triggerParam" },
			]}
			ChatbotStateExpression: { name: "ChatbotStateExpression", fields: [
				{ name: "Source", type: "string", json: "source,omitempty" },
				{ name: "Expr", type: "string", json: "expr,omitempty" },
				{ name: "Value", type: "any", json: "value,omitempty" },
				{ name: "Type", type: "string", json: "type,omitempty" },
			]}

			ChatbotScenario: { name: "ChatbotScenario", fields: [
				{ name: "ID", type: "string", json: "id" },
				{ name: "Name", type: "string", json: "name,omitempty" },
				{ name: "Type", type: "string", json: "type" },
				{ name: "AgentID", type: "uint64", json: "agentID,string,omitempty" },
				{ name: "Config", goType: "json.RawMessage", json: "config,omitempty" },
				{ name: "Automation", type: _chatbotDefs.ChatbotScenarioAutomation, json: "automation" },
			]}
			ChatbotScenarioAutomation: { name: "ChatbotScenarioAutomation", fields: [
				{ name: "Before", type: _chatbotDefs.ChatbotAutomationHook, json: "before" },
				{ name: "After", type: _chatbotDefs.ChatbotAutomationHook, json: "after" },
			]}
			ChatbotScenarios: { name: "ChatbotScenarios", elem: _chatbotDefs.ChatbotScenario }
			ChatbotAllowedOrigins: { name: "ChatbotAllowedOrigins", elem: "string" }

			ChatbotStyling: { name: "ChatbotStyling", fields: [
				{ name: "LogoAttachmentID", type: "uint64", json: "logoAttachmentID,string,omitempty" },
				{ name: "LogoURL", type: "string", json: "logoURL,omitempty" },
				{ name: "FontFamily", type: "string", json: "fontFamily,omitempty" },
				{ name: "FontSizes", type: _chatbotDefs.ChatbotFontSizes, json: "fontSizes,omitempty" },
				{ name: "Colors", type: _chatbotDefs.ChatbotColors, json: "colors,omitempty" },
				{ name: "Launcher", type: _chatbotDefs.ChatbotLauncher, json: "launcher,omitempty" },
			]}
			ChatbotFontSizes: { name: "ChatbotFontSizes", fields: [
				{ name: "Base", type: "string", json: "base,omitempty" },
				{ name: "Small", type: "string", json: "small,omitempty" },
				{ name: "Heading", type: "string", json: "heading,omitempty" },
			]}
			ChatbotColors: { name: "ChatbotColors", fields: [
				{ name: "Primary", type: "string", json: "primary,omitempty" },
				{ name: "PrimaryText", type: "string", json: "primaryText,omitempty" },
				{ name: "Header", type: "string", json: "header,omitempty" },
				{ name: "HeaderText", type: "string", json: "headerText,omitempty" },
				{ name: "Background", type: "string", json: "background,omitempty" },
				{ name: "Text", type: "string", json: "text,omitempty" },
				{ name: "UserBubble", type: "string", json: "userBubble,omitempty" },
				{ name: "AgentBubble", type: "string", json: "agentBubble,omitempty" },
			]}
			ChatbotLauncher: { name: "ChatbotLauncher", fields: [
				{ name: "IconURL", type: "string", json: "iconURL,omitempty" },
				{ name: "IconAttachmentID", type: "uint64", json: "iconAttachmentID,string,omitempty" },
				{ name: "IconVisible", type: "bool", json: "iconVisible" },
				{ name: "Label", type: "string", json: "label,omitempty" },
				{ name: "ButtonLabel", type: "string", json: "buttonLabel,omitempty" },
				{ name: "Size", type: "string", json: "size,omitempty" },
				{ name: "Shape", type: "string", json: "shape,omitempty" },
				{ name: "Position", type: "string", json: "position,omitempty" },
				{ name: "StartOpen", type: "bool", json: "startOpen,omitempty" },
			]}
		}

chatbot: {
	features: {
		labels:            true
		labelResourceType: "chatbot"
	}

	types: {
		gen: true
		// ChatbotStyling has a custom Value (kept hand-written)
		// ChatbotStateExpression has a field named Value — Scan/Value would conflict
		jsonTypesSkip: ["ChatbotStyling", "ChatbotStateExpression"]

		defs: _chatbotDefs
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
				type: _chatbotDefs.ChatbotAllowedOrigins
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
				type: _chatbotDefs.ChatbotHandoff
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
			}
			styling: {
				type: _chatbotDefs.ChatbotStyling
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
			}
			scenarios: {
				type: _chatbotDefs.ChatbotScenarios
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
		extraServices:       false
		genAccessController: true

		undelete: true

		// widget_key is generated at create and rotated only via
		// RegenerateWidgetKey; a plain update must not wipe it.
		omitUpdateFields: ["widget_key"]

		customBodyOps: ["lookup", "search", "create", "update", "delete", "undelete"]

		customFunctions: [
			{
				name: "RegenerateWidgetKey"
				cap:  "write"
				args: [{name: "ID", goType: "uint64"}]
				results: [{name: "c", goType: "*types.Chatbot"}, {name: "err", goType: "error"}]
			},
		]

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
