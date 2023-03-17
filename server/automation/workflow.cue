package automation

import (
	"github.com/cortezaproject/corteza/server/codegen/schema"
)

workflow: {
	model: {
		ident: "automation_workflows"
		attributes: {
			id: schema.IdField
			handle: schema.HandleField
			meta: {
				goType: "*types.WorkflowMeta"
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
			}
			enabled: {
				sortable: true,
				goType: "bool"
				dal: { type: "Boolean", default: true }
			}
			trace: {
				goType: "bool"
				dal: { type: "Boolean", default: false }
			}
			keep_sessions: {
				goType: "int"
				dal: { type: "Number", default: 0, meta: { "rdbms:type": "integer" } }
		  }
			scope: {
				goType: "*expr.Vars"
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
			}
			steps: {
				goType: "types.WorkflowStepSet"
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
			}
			paths: {
				goType: "types.WorkflowPathSet"
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
			}
			issues: {
				goType: "types.WorkflowIssueSet"
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
				envoy: {
					yaml: {
						omitEncoder: true
					}
				}
			}

			run_as: schema.AttributeUserRef

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
			supportMappedInput: true
			mappedField: "Handle"
			identKeyAlias: ["workflows"]
			extendedResourceDecoders: [{
				ident: "triggers"
				expIdent: "Triggers"
				supportMappedInput: false
				identKeys: ["triggers"]
			}]
			extendedResourceEncoders: [{
				ident: "trigger"
				expIdent: "Trigger"
				identKey: "trigger"
			}]
		}
		store: {
			customFilterBuilder: true
			extendedDecoder: true
		}
	}

	filter: {
		struct: {
			workflow_id: { goType: "[]string", ident: "workflowID", storeIdent: "id" }
			handle: { goType: "string" }
			sub_workflow: { goType: "filter.State" }
			deleted: { goType: "filter.State", storeIdent: "deleted_at" }
			disabled: { goType: "filter.State", storeIdent: "enabled" }
		}

		query: ["handle"]
		byValue: ["workflow_id", "handle"]
		byNilState: ["deleted"]
		byFalseState: ["disabled"]
	}

	rbac: {
		operations: {
			"read": description:            "Read workflow"
			"update": description:          "Update workflow"
			"delete": description:          "Delete workflow"
			"undelete": description:        "Undelete workflow"
			"execute": description:         "Execute workflow"
			"triggers.manage": description: "Manage workflow triggers"
			"sessions.manage": description: "Manage workflow sessions"
		}
	}

	store: {
		ident: "automationWorkflow"

		api: {
			lookups: [
				{
					fields: ["id"]
					description: """
						searches for workflow by ID

						It returns workflow even if deleted
						"""
				}, {
					fields: ["handle"]
					nullConstraint: ["deleted_at"]
					constraintCheck: true
					description: """
						searches for workflow by their handle

						It returns only valid workflows
						"""
				}
			]
		}
	}
}
