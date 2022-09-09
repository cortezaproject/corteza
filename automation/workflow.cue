package automation

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
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
			}
			steps: {
				goType: "types.WorkflowStepSet"
				dal: { type: "JSON", defaultEmptyObject: true }
			}
			paths: {
				goType: "types.WorkflowPathSet"
				dal: { type: "JSON", defaultEmptyObject: true }
			}
			issues: {
				goType: "types.WorkflowIssueSet"
				dal: { type: "JSON", defaultEmptyObject: true }
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

	filter: {
		struct: {
			deleted: { goType: "filter.State", storeIdent: "deleted_at" }
			disabled: { goType: "filter.State", storeIdent: "enabled" }
		}

		query: ["handle"]
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
