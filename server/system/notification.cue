package system

import (
	"github.com/crusttech/human/server/codegen/schema"
)

notification: {
	features: {
		labels: false
		projectScoped: true
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
			recipient: schema.AttributeUserRef
			created_by: schema.AttributeUserRef
			read_at: schema.SortableTimestampNilField
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

	store: {
		api: {
			lookups: [
				{ fields: ["id"] }
			]
		}
	}
}
