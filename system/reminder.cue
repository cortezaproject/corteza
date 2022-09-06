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
				resource: { sortable: true }
				payload: { goType: "rawJson" }
				snooze_count: { goType: "uint" }
				assigned_to: { goType: "uint64" }
				assigned_by: { goType: "uint64" }
				assigned_at: schema.SortableTimestampField
				dismissed_by: { goType: "uint64" }
				dismissed_at: schema.SortableTimestampNilField
				remind_at: schema.SortableTimestampNilField
				created_at: schema.SortableTimestampNowField
				updated_at: schema.SortableTimestampNilField
				deleted_at: schema.SortableTimestampNilField
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
