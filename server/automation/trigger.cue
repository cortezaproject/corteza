package automation

import (
	"github.com/cortezaproject/corteza/server/codegen/schema"
)

trigger: {
	model: {
		ident: "automation_triggers"
		attributes: {
			id:  schema.IdField
			workflow_id: {
				sortable: true,
				ident: "workflowID",
				goType: "uint64",
				storeIdent: "rel_workflow"
				dal: { type: "Ref", refModelResType: "corteza::automation:workflow" }
			}
			step_id: {
				ident: "stepID",
				goType: "uint64",
				storeIdent: "rel_step"
				dal: { type: "ID" }
				envoy: {
					yaml: {
						identKeyAlias: ["stepID", "step_id"]
					}
				}
			}
			enabled: {
				sortable: true,
				goType: "bool"
				dal: { type: "Boolean", default: true }
			}
			meta: {
				goType: "*types.TriggerMeta"
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
			}
			resource_type: {
				sortable: true,
				goType: "string"
				dal: { length: 64 }
				envoy: {
					yaml: {
						identKeyAlias: ["resourceType", "resource_type"]
					}
				}
			}
			event_type: {
				sortable: true,
				goType: "string"
				dal: {}
				envoy: {
					yaml: {
						identKeyAlias: ["eventType", "event_type"]
					}
				}
			}
			constraints: {
				goType: "types.TriggerConstraintSet"
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
			}
			input: {
				goType: "*expr.Vars"
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
			}

			owned_by:   schema.AttributeUserRef
			created_at: schema.SortableTimestampNowField
			updated_at: schema.SortableTimestampNilField
			deleted_at: schema.SortableTimestampNilField
			created_by: schema.AttributeUserRef
			updated_by: schema.AttributeUserRef
			deleted_by: schema.AttributeUserRef
		}

		indexes: {
			"primary": { attribute: "id" }
		}
	}

	envoy: {
		yaml: {
			supportMappedInput: false
			identKeyAlias: []
		}
		store: {
			customFilterBuilder: true
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
