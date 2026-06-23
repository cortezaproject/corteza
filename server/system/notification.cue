package system

import (
	"github.com/crusttech/human/server/codegen/schema"
)

notification: {
	features: {
		labels: false
		projectScoped: true
	}

	types: {
		gen: true
		// NotificationConfig has a custom Scan/Value (kept hand-written)
		jsonTypesSkip: ["NotificationConfig"]
	}

	model: {
		omitGetterSetter: true

		attributes: {
			id:     schema.IdField
			tenant_id:  schema.TenantRefField
			project_id: schema.ProjectRefField
			kind: {
				sortable: true
				goType: "types.NotificationKind"
				dal: { type: "Text", length: 32 }
			}
			config: {
				goType: "types.NotificationConfig"
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
