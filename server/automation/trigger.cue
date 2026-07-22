package automation

import (
	"github.com/crusttech/human/server/codegen/schema"
)

_triggerDefs: {
			TriggerMeta: { name: "TriggerMeta", fields: [
				{ name: "Description", type: "string", json: "description" },
				{ name: "Visual", goType: "map[string]interface{}", json: "visual" },
			]}
			TriggerConstraint: { name: "TriggerConstraint", fields: [
				{ name: "Name", type: "string", json: "name" },
				{ name: "Op", type: "string", json: "op,omitempty" },
				{ name: "Values", type: "string", slice: true, json: "values,omitempty" },
			]}
			TriggerConstraintSet: { name: "TriggerConstraintSet", elem: _triggerDefs.TriggerConstraint, elemPtr: true }
		}

trigger: {
	features: {
		labelResourceType: "trigger"
	}
	types: {
		// generate the Trigger struct from the model into types/trigger.gen.go
		gen: true
		// Input *expr.Vars needs the expr package qualifier
		imports: ["github.com/crusttech/human/server/pkg/expr"]
		// TriggerMeta keeps a hand-written pointer-returning ParseTriggerMeta
		// (REST request controllers depend on the *TriggerMeta return), which is
		// incompatible with the value-returning generated version, so it is
		// excluded from JSON helper generation.
		jsonTypesPtr: ["TriggerMeta"]
		defs: _triggerDefs
	}
	model: {
		ident: "automation_triggers"
		attributes: {
			id:  schema.IdField
			tenant_id:  schema.TenantRefField
			project_id: schema.ProjectRefField
			workflow_id: {
				sortable: true,
				ident: "workflowID",
				goType: "uint64",
				storeIdent: "rel_workflow"
				dal: { type: "Ref", refModelResType: "corteza::automation:workflow" }
				// hand-written tag has no omitempty (convention would add it for refs)
				json: { field: "workflowID", string: true }
			}
			step_id: {
				ident: "stepID",
				goType: "uint64",
				storeIdent: "rel_step"
				dal: { type: "ID" }
				// hand-written tag has no omitempty (convention would add it for refs)
				json: { field: "stepID", string: true }
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
				type: _triggerDefs.TriggerMeta
				ptr: true
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
				// hand-written tag is meta,omitempty (convention has no omitempty)
				json: { omitEmpty: true }
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
				type: _triggerDefs.TriggerConstraintSet
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

			// hand-written tag is ownedBy,string (no omitempty); convention would add it
			owned_by:   schema.AttributeUserRef & { json: { field: "ownedBy", string: true } }
			created_at: schema.SortableTimestampNowField
			updated_at: schema.SortableTimestampNilField
			deleted_at: schema.SortableTimestampNilField
			// hand-written tag is createdBy,string (no omitempty); convention would add it
			created_by: schema.AttributeUserRef & { json: { field: "createdBy", string: true } }
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
			tenant_id: schema.TenantFilterField
			project_id: schema.ProjectFilterField
			workflow_id: { goType: "[]uint64", ident: "workflowID", storeIdent: "rel_workflow" }
			event_type: { goType: "string" }
			resource_type: { goType: "string" }
		}

		byValue: ["trigger_id", "workflow_id", "event_type", "resource_type"]
		byNilState: ["deleted"]
		byFalseState: ["disabled"]
	}

	refs: {
		items: [
			{path: "WorkflowID", kind: "KindAutomationWorkflow", reason: "ReasonTriggerWorkflow"},
		]
		extended: true
	}

	service: {
		// project-state lock: block writes unless the owning project is a draft
		guard: true

		// trigger exposes a bespoke LookupByID (custom name + label load), kept in
		// the companion file -- do not generate FindByID.
		lookup: false

		search:   true
		create:   true
		update:   true
		delete:   true
		undelete: true

		customBodyOps: ["search", "create", "update", "delete", "undelete"]

		customFunctions: [
			{
				name: "SearchOnManual"
				cap:  "read"
				args: [
					{name: "workflowID", goType: "uint64"},
					{name: "stepID", goType: "uint64"},
				]
				results: [
					{name: "res", goType: "*types.Trigger"},
					{name: "err", goType: "error"},
				]
			},
		]

		customAccessOps: ["create", "update", "delete", "undelete"]
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
