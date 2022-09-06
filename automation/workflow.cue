package automation

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

workflow: {
	model: {
		ident: "automation_workflows"
		attributes: {
			id:          schema.IdField
			handle:      schema.HandleField
			meta: { goType: "*types.WorkflowMeta" }
			enabled: { sortable: true, goType: "bool" }
			trace: { goType: "bool" }
			keep_sessions: { goType: "int" }
			scope: { goType: "*expr.Vars" }
			steps: { goType: "types.WorkflowStepSet" }
			paths: { goType: "types.WorkflowPathSet" }
			issues: { goType: "types.WorkflowIssueSet" }
			run_as: { goType: "uint64" }

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
