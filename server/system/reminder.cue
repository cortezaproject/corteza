package system

import (
	"github.com/crusttech/human/server/codegen/schema"
)

reminder: {
	features: {
		labels: false
	}

	model: {
		omitGetterSetter: true

		attributes: {
			id:     schema.IdField
			tenant_id:  schema.TenantRefField
			project_id: schema.ProjectRefField
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

	envoy: {
		omit: true
	}

	service: {
		// Reminder access is entirely non-RBAC: it is gated by comparing the
		// reminder's assigned_to with the current user (checkAssignTo) plus the
		// custom CanAssignReminder check. There is no RBAC Can{Read,Search,Create,
		// Update,Delete}Reminder, so every op delegates its body to a hand-written
		// on<Op> handler and the AC checks the template would otherwise emit
		// (search/create) are suppressed via customAccessOps.
		//
		// The action-log Update prop is named "updated" (not the default "update").
		updateProp: "updated"

		// snooze_count and dismissed_by are managed by onSnooze/onDismiss (and the
		// onUpdate remind-at reset), not by a plain update payload; excluding them
		// keeps a generic update from zeroing them.
		omitUpdateFields: ["snooze_count", "dismissed_by"]

		undelete: false

		customBodyOps:   ["lookup", "search", "create", "update", "delete"]
		customAccessOps: ["search", "create"]

		customFunctions: [
			{
				name: "FindByIDs"
				cap:  "read"
				args: [{name: "IDs", goType: "[]uint64"}]
				results: [{name: "rr", goType: "types.ReminderSet"}, {name: "err", goType: "error"}]
			},
			{
				name:   "Dismiss"
				cap:    "write"
				action: "Dismiss"
				args: [{name: "ID", goType: "uint64"}]
				results: [{name: "err", goType: "error"}]
			},
			{
				name:   "Undismiss"
				cap:    "write"
				action: "Dismiss"
				args: [{name: "ID", goType: "uint64"}]
				results: [{name: "err", goType: "error"}]
			},
			{
				name:   "Snooze"
				cap:    "write"
				action: "Snooze"
				args: [{name: "ID", goType: "uint64"}, {name: "remindAt", goType: "*time.Time"}]
				results: [{name: "err", goType: "error"}]
			},
		]
		customFunctionImports: ["\"time\""]
	}

	filter: {
		struct: {
			reminder_id: {goType: "[]uint64", ident: "reminderID", storeIdent: "id"}
			tenant_id: schema.TenantFilterField
			project_id: schema.ProjectFilterField
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
