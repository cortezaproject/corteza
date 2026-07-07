package system

import (
	"github.com/crusttech/human/server/codegen/schema"
)

_connectionDefs: {
			ConnectionMeta: { name: "ConnectionMeta", fields: [
				{ name: "Short", type: "string", json: "short" },
				{ name: "Description", type: "string", json: "description" },
				{ name: "Icon", type: "string", json: "icon" },
				{ name: "Tags", slice: true, type: "string", json: "tags" },
			]}
			ConnectionService: { name: "ConnectionService", fields: [
				{ name: "BaseURL", type: _connectionDefs.ConnectionTemplate, json: "baseURL" },
				{ name: "Protocol", type: "string", json: "protocol" },
				{ name: "ContentType", type: "string", json: "contentType" },
				{ name: "Headers", goType: "map[string]ConnectionTemplate", json: "headers,omitempty" },
				{ name: "Auth", type: _connectionDefs.ConnectionAuth, json: "auth" },
				{ name: "Probe", ptr: true, type: _connectionDefs.ConnectionProbe, json: "probe,omitempty" },
				{ name: "Params", slice: true, type: _connectionDefs.ConnectionPlaceholder, json: "params,omitempty" },
			]}
			ConnectionTemplate: { name: "ConnectionTemplate", fields: [
				{ name: "Value", type: "string", json: "value" },
				{ name: "Placeholders", slice: true, type: _connectionDefs.ConnectionPlaceholder, json: "placeholders,omitempty" },
			]}
			ConnectionAuth: { name: "ConnectionAuth", fields: [
				{ name: "Method", type: "string", json: "method" },
				{ name: "Params", goType: "map[string]ConnectionTemplate", json: "params,omitempty" },
			]}
			ConnectionProbe: { name: "ConnectionProbe", fields: [
				{ name: "Path", type: _connectionDefs.ConnectionTemplate, json: "path" },
				{ name: "ExpectedStatus", type: "int", json: "expectedStatus,omitempty" },
			]}
			ConnectionPlaceholder: { name: "ConnectionPlaceholder", fields: [
				{ name: "Name", type: "string", json: "name" },
				{ name: "Label", type: "string", json: "label,omitempty" },
				{ name: "Type", type: "string", json: "type" },
				{ name: "Description", type: "string", json: "description" },
				{ name: "Required", type: "bool", json: "required" },
				{ name: "Default", type: "string", json: "default" },
				{ name: "Options", slice: true, type: "string", json: "options,omitempty" },
			]}
			ConnectionDerivedParam: { name: "ConnectionDerivedParam", fields: [
				{ name: "Name", type: "string", json: "name" },
				{ name: "Label", type: "string", json: "label" },
				{ name: "Scope", slice: true, type: "string", json: "scope" },
				{ name: "Type", type: "string", json: "type" },
				{ name: "Description", type: "string", json: "description" },
				{ name: "Required", type: "bool", json: "required" },
				{ name: "Default", type: "string", json: "default" },
				{ name: "Options", slice: true, type: "string", json: "options,omitempty" },
			]}
			ConnectionResources: { name: "ConnectionResources", elem: _connectionDefs.ConnectionResource }
			ConnectionOperations: { name: "ConnectionOperations", elem: _connectionDefs.ConnectionOperation }
			ConnectionResource: { name: "ConnectionResource", fields: [
				{ name: "Handle", type: "string", json: "handle" },
				{ name: "Meta", type: _connectionDefs.ConnectionResourceMeta, json: "meta" },
				{ name: "Endpoint", type: _connectionDefs.ConnectionTemplate, json: "endpoint" },
				{ name: "Operations", type: _connectionDefs.ConnectionResourceOperations, json: "operations,omitempty" },
				{ name: "Fields", slice: true, type: _connectionDefs.ConnectionResourceField, json: "fields" },
				{ name: "Webhooks", slice: true, type: _connectionDefs.ConnectionWebhook, json: "webhooks,omitempty" },
			]}
			ConnectionResourceMeta: { name: "ConnectionResourceMeta", fields: [
				{ name: "Short", type: "string", json: "short" },
				{ name: "Description", type: "string", json: "description" },
				{ name: "Icon", ptr: true, type: _connectionDefs.ConnectionIcon, json: "icon,omitempty" },
			]}
			ConnectionIcon: { name: "ConnectionIcon", fields: [
				{ name: "Type", type: "string", json: "type" },
				{ name: "Value", type: "string", json: "value" },
			]}
			ConnectionResourceField: { name: "ConnectionResourceField", fields: [
				{ name: "Name", type: "string", json: "name" },
				{ name: "Type", type: "string", json: "type" },
				{ name: "Selector", slice: true, type: "string", json: "selector,omitempty" },
				{ name: "Meta", goType: "map[string]any", json: "meta,omitempty" },
			]}
			ConnectionWebhook: { name: "ConnectionWebhook", fields: [
				{ name: "Event", type: "string", json: "event" },
				{ name: "Path", type: "string", json: "path" },
				{ name: "Payload", slice: true, type: _connectionDefs.ConnectionWebhookField, json: "payload,omitempty" },
				{ name: "Mapping", goType: "map[string]string", json: "mapping,omitempty" },
			]}
			ConnectionWebhookField: { name: "ConnectionWebhookField", fields: [
				{ name: "Name", type: "string", json: "name" },
				{ name: "Type", type: "string", json: "type" },
				{ name: "Selector", slice: true, type: "string", json: "selector,omitempty" },
				{ name: "Meta", goType: "map[string]any", json: "meta,omitempty" },
			]}
			ConnectionResourceOperations: { name: "ConnectionResourceOperations", fields: [
				{ name: "List", ptr: true, type: _connectionDefs.ConnectionHTTPAction, json: "list,omitempty" },
				{ name: "Read", ptr: true, type: _connectionDefs.ConnectionHTTPAction, json: "read,omitempty" },
				{ name: "Create", ptr: true, type: _connectionDefs.ConnectionHTTPAction, json: "create,omitempty" },
				{ name: "Update", ptr: true, type: _connectionDefs.ConnectionHTTPAction, json: "update,omitempty" },
				{ name: "Delete", ptr: true, type: _connectionDefs.ConnectionHTTPAction, json: "delete,omitempty" },
			]}
			ConnectionHTTPAction: { name: "ConnectionHTTPAction", fields: [
				{ name: "Method", type: "string", json: "method" },
				{ name: "Path", type: _connectionDefs.ConnectionTemplate, json: "path" },
				{ name: "Headers", goType: "map[string]ConnectionTemplate", json: "headers,omitempty" },
				{ name: "QueryParams", goType: "map[string]ConnectionTemplate", json: "queryParams,omitempty" },
				{ name: "BodyTemplate", type: _connectionDefs.ConnectionTemplate, json: "bodyTemplate,omitempty" },
				{ name: "ResponseMap", goType: "map[string]string", json: "responseMap,omitempty" },
			]}
			ConnectionOperation: { name: "ConnectionOperation", fields: [
				{ name: "Handle", type: "string", json: "handle" },
				{ name: "Meta", type: _connectionDefs.ConnectionResourceMeta, json: "meta" },
				{ name: "Input", slice: true, type: _connectionDefs.ConnectionOperationInputField, json: "input,omitempty" },
				{ name: "Output", slice: true, type: _connectionDefs.ConnectionOperationOutputField, json: "output,omitempty" },
				{ name: "Steps", slice: true, type: _connectionDefs.ConnectionOperationStep, json: "steps,omitempty" },
			]}
			ConnectionOperationInputField: { name: "ConnectionOperationInputField", fields: [
				{ name: "Name", type: "string", json: "name" },
				{ name: "Type", type: "string", json: "type" },
				{ name: "Required", type: "bool", json: "required,omitempty" },
				{ name: "Aggregate", type: "bool", json: "aggregate,omitempty" },
				{ name: "Meta", goType: "map[string]any", json: "meta,omitempty" },
			]}
			ConnectionOperationOutputField: { name: "ConnectionOperationOutputField", fields: [
				{ name: "Name", type: "string", json: "name" },
				{ name: "Type", type: "string", json: "type" },
				{ name: "Selector", slice: true, type: "string", json: "selector,omitempty" },
				{ name: "Meta", goType: "map[string]any", json: "meta,omitempty" },
			]}
			ConnectionOperationStep: { name: "ConnectionOperationStep", fields: [
				{ name: "Type", type: "string", json: "type" },
				{ name: "HTTP", ptr: true, type: _connectionDefs.ConnectionHTTPAction, json: "http,omitempty" },
				{ name: "MimeBuild", ptr: true, type: _connectionDefs.ConnectionMimeBuildAction, json: "mime_build,omitempty" },
			]}
			ConnectionMimeBuildAction: { name: "ConnectionMimeBuildAction", fields: [
				{ name: "To", type: "string", json: "to" },
				{ name: "Subject", type: "string", json: "subject" },
				{ name: "Body", type: "string", json: "body" },
				{ name: "From", type: "string", json: "from,omitempty" },
				{ name: "Output", type: "string", json: "output" },
			]}
		}

connection: {
	model: {
		attributes: {
			id:       schema.IdField
			handle:   schema.HandleField
			revision: {
				sortable: true
				goType:   "int"
				dal: { type: "Number", meta: { "rdbms:type": "integer" } }
			}
			status: {
				sortable: true
				dal: { type: "Text", length: 32 }
			}
			source: {
				sortable: true
				dal: { type: "Text", length: 16 }
				json: { field: "source", omitEmpty: true }
			}

			meta: {
				type: _connectionDefs.ConnectionMeta
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
			service: {
				type: _connectionDefs.ConnectionService
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
			}
			resources: {
				type: _connectionDefs.ConnectionResources
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
			}
			operations: {
				type: _connectionDefs.ConnectionOperations
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
			}

			// Struct-only fields (not persisted via store/dal).
			derived_params: {
				type: _connectionDefs.ConnectionDerivedParam
				slice: true
				store: false
				omitGetter: true
				omitSetter: true
				json: { field: "derivedParams", omitEmpty: true }
			}
			catalog_id: {
				expIdent: "CatalogID"
				goType: "string"
				store: false
				omitGetter: true
				omitSetter: true
				json: { field: "catalogID", omitEmpty: true }
			}
			installed_count: {
				goType: "int"
				store: false
				omitGetter: true
				omitSetter: true
				json: { field: "installedCount", omitEmpty: true }
			}

			created_at: schema.SortableTimestampNowField
			updated_at: schema.SortableTimestampNilField
			deleted_at: schema.SortableTimestampNilField
			created_by: schema.AttributeUserRef & {
				json: { field: "createdBy", string: true }
			}
			updated_by: schema.AttributeUserRef
			deleted_by: schema.AttributeUserRef
		}

		indexes: {
			"primary": { attribute: "id" }
		}
	}

	filter: {
		struct: {
			handle: {goType: "string"}
			status: {goType: "[]string"}
			source: {goType: "string"}
			query:  {goType: "string"}
			deleted: {goType: "filter.State", storeIdent: "deleted_at"}
		}
		byValue: ["handle", "status", "source"]
		byNilState: ["deleted"]
	}

	features: {
		labels:            true
		labelResourceType: "connection"
	}

	types: {
		gen: true
		// ConnectionTemplate and ConnectionIcon have a field named Value — Scan/Value would conflict
		jsonTypesSkip: ["ConnectionTemplate", "ConnectionIcon"]
		defs: _connectionDefs
	}

	envoy: {
		omit: true
	}

	rbac: {
		operations: {
			"read": description:   "Read connection"
			"update": description: "Update connection"
			"delete": description: "Delete connection"
			"install": description: "Install connection"
		}
	}

	service: {
		lookup:   false
		search:   false
		update:   false
		undelete: true

		hooks: {
			beforeCreate:   true
			afterCreate:    true
			beforeDelete:   true
			beforeUndelete: true
		}
	}

	store: {
		api: {
			lookups: [
				{
					fields: ["id"]
					description: """
						searches for connection by ID

						It returns connection even if deleted
						"""
				}, {
					fields: ["handle"]
					nullConstraint: ["deleted_at"]
					constraintCheck: true
					description: """
						searches for connection by handle

						It returns only valid connection (not deleted)
						"""
				},
			]
		}
	}
}
