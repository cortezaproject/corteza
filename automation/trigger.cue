package automation

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

trigger: {
	model: {
		ident: "automation_triggers"
		attributes: {
			id:          schema.IdField
			workflow_id: { sortable: true, ident: "workflowID", goType: "uint64", storeIdent: "rel_workflow" }
			step_id: { ident: "stepID", goType: "uint64", storeIdent: "rel_step" }
			enabled: { sortable: true, goType: "bool" }
			resource_type: { sortable: true, goType: "string" }
			event_type: { sortable: true, goType: "string" }
			meta: { goType: "*types.TriggerMeta" }
			constraints: { goType: "types.TriggerConstraintSet" }
			input: { goType: "*expr.Vars" }

			created_at: schema.SortableTimestampNowField
			updated_at: schema.SortableTimestampNilField
			deleted_at: schema.SortableTimestampNilField
			owned_by:   schema.AttributeUserRef
			created_by: schema.AttributeUserRef
			updated_by: schema.AttributeUserRef
			deleted_by: schema.AttributeUserRef
		}
	}

	filter: {
		struct: {
			deleted: { goType: "filter.State", storeIdent: "deleted_at" }
			disabled: { goType: "filter.State", storeIdent: "enabled" }
			trigger_id: { goType: "[]uint64", ident: "triggerID", storeIdent: "id" }
			workflow_id: { goType: "[]uint64", ident: "workflowID", storeIdent: "rel_workflow" }
			event_type: { goType: "string" }
			resource_type: { goType: "string" }
		}

		byValue: ["trigger_id", "workflow_id", "event_type", "resource_type"]
		byNilState: ["deleted"]
		byFalseState: ["disabled"]
	}

	store: {
		ident: "automationTrigger"

		api: {
			lookups: [
				{
					fields: ["id"]
					description: """
						searches for trigger by ID

						It returns trigger even if deleted
						"""
				}
			]
		}
	}
}
