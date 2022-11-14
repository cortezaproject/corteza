package system

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

reminder: {
	features: {
		labels: false
	}

	model: {
		attributes: {
			id:     schema.IdField
			resource: {
				sortable: true
				dal: { type: "Text", length: 512 }
			}
			payload: {
				goType: "rawJson"
				dal: { type: "JSON", defaultEmptyObject: true }
			}
			snooze_count: {
				goType: "uint"
				dal: { type: "Number", meta: { "rdbms:type": "integer" } }
			}
			assigned_to: schema.AttributeUserRef
			assigned_by: schema.AttributeUserRef
			assigned_at: schema.SortableTimestampField
			dismissed_by: schema.AttributeUserRef
			dismissed_at: schema.SortableTimestampNilField
			remind_at: schema.SortableTimestampNilField
			created_at: schema.SortableTimestampNowField
			updated_at: schema.SortableTimestampNilField
			deleted_at: schema.SortableTimestampNilField
		}

		indexes: {
			"primary": { attribute: "id" }
			"assigned_to": { attribute: "assigned_to" }
			"resource": { attribute: "resource" }
		}
	}

	filter: {
		struct: {
			reminder_id: {goType: "[]uint64", ident: "reminderID", storeIdent: "id"}
			resource: {}
			assigned_to: {goType: "uint64"}
			scheduled_from: {goType: "uint64"}
			scheduled_until: {goType: "uint64"}
			exclude_dismissed: { goType: "bool" }
			include_deleted: { goType: "bool" }
			scheduled_only: { goType: "bool" }
			deleted: {goType: "filter.State", storeIdent: "deleted_at"}
		}

		byValue: ["reminder_id", "assigned_to"]
	}

	store: {
		api: {
			lookups: [
				{ fields: ["id"] }
			]
		}
	}
}
