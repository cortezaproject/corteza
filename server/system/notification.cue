package system

import (
	"github.com/crusttech/human/server/codegen/schema"
)

_notificationDefs: {
			NotificationConfig: { name: "NotificationConfig", fields: [
				{ name: "Simple", type: _notificationDefs.SimpleNotificationConfig, json: "simple" },
				{ name: "Record", type: _notificationDefs.RecordNotificationConfig, json: "record" },
			]}
			SimpleNotificationConfig: { name: "SimpleNotificationConfig", fields: [
				{ name: "Title", type: "string", json: "title" },
				{ name: "Description", type: "string", json: "description" },
			]}
			RecordNotificationConfig: { name: "RecordNotificationConfig", fields: [
				{ name: "Title", type: "string", json: "title" },
				{ name: "Description", type: "string", json: "description" },
				{ name: "ModuleID", type: "uint64", json: "moduleID,string" },
				{ name: "NamespaceID", type: "uint64", json: "namespaceID,string" },
				{ name: "RecordID", type: "uint64", json: "recordID,string" },
				{ name: "OpenMode", type: _notificationDefs.OpenModeType, json: "openMode,omitempty" },
				{ name: "Edit", type: "bool", json: "edit,omitempty" },
			]}
			OpenModeType: { name: "OpenModeType", values: [
				{ ident: "OpenModeModal", value: "modal" },
				{ ident: "OpenModeNewTab", value: "newTab" },
				{ ident: "OpenModeSameTab", value: "sameTab" },
			]}
			NotificationKind: { name: "NotificationKind", values: [
				{ ident: "NotificationKindSimple", value: "simple" },
				{ ident: "NotificationKindRecord", value: "record" },
			]}
		}

notification: {
	features: {
		labels: false
	}

	types: {
		gen: true
		// NotificationConfig has a custom Scan/Value (kept hand-written)
		jsonTypesSkip: ["NotificationConfig"]

		defs: _notificationDefs
	}

	model: {
		omitGetterSetter: true

		attributes: {
			id:     schema.IdField
			tenant_id:  schema.TenantRefField
			project_id: schema.ProjectRefField
			kind: {
				sortable: true
				type: _notificationDefs.NotificationKind
				dal: { type: "Text", length: 32 }
			}
			config: {
				type: _notificationDefs.NotificationConfig
				dal: { type: "JSON", defaultEmptyObject: true }
			}
			recipient: schema.AttributeUserRef & {
				// hand-written tag is `recipient,string` (no omitempty);
				// the Ref convention would add omitempty.
				json: {field: "recipient", string: true}
			}
			created_by: schema.AttributeUserRef & {
				// hand-written tag is `createdBy,string` (no omitempty);
				// the Ref convention would add omitempty.
				json: {field: "createdBy", string: true}
			}
			read_at: schema.SortableTimestampNilField & {
				// hand-written tag is `readAt` (no omitempty);
				// the time convention would add omitempty.
				json: {field: "readAt"}
			}
			created_at: schema.SortableTimestampNowField
			updated_at: schema.SortableTimestampNilField
			deleted_at: schema.SortableTimestampNilField
		}

		indexes: {
			"primary": { attribute: "id" }
			"recipient": { attribute: "recipient" }
			"kind": { attribute: "kind" }
		}
	}

	envoy: {
		omit: true
	}

	filter: {
		struct: {
			notification_id: {goType: "[]uint64", ident: "notificationID", storeIdent: "id"}
			tenant_id: schema.TenantFilterField
			project_id: schema.ProjectFilterField
			kind: {goType: "[]types.NotificationKind"}
			recipient: {goType: "uint64"}
			read: {goType: "filter.State", storeIdent: "read_at"}
			deleted: {goType: "filter.State", storeIdent: "deleted_at"}
		}

		byValue: ["notification_id", "recipient"]
		byNilState: ["read", "deleted"]
	}

	service: {
		customBodyOps:   ["lookup", "search", "create", "update", "delete"]
		customFunctions: [
			{
				name:    "MarkAsRead"
				cap:     "write"
				action:  "MarkAsRead"
				args:    [{name: "ID", goType: "uint64"}]
				results: [{name: "err", goType: "error"}]
			},
			{
				name:    "MarkAsUnread"
				cap:     "write"
				action:  "MarkAsUnread"
				args:    [{name: "ID", goType: "uint64"}]
				results: [{name: "err", goType: "error"}]
			},
			{
				name:    "MarkAllAsRead"
				cap:     "write"
				action:  "MarkAllAsRead"
				results: [{name: "err", goType: "error"}]
			},
			{
				name:    "MarkAllAsUnread"
				cap:     "write"
				action:  "MarkAllAsUnread"
				results: [{name: "err", goType: "error"}]
			},
		]
		customAccessOps: ["search", "create"]

		updateProp: "updated"
	}

	store: {
		api: {
			lookups: [
				{ fields: ["id"] }
			]
		}
	}
}
