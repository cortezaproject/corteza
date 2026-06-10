package automation

import (
	"github.com/crusttech/human/server/codegen/schema"
)

ng_automation: {
	features: {
		projectScoped: true
	}

	model: {
		attributes: {
			id:         schema.IdField
			tenant_id:  schema.TenantRefField
			project_id: schema.ProjectRefField
			handle:     schema.HandleField
			meta: {
				goType: "*types.NgAutomationMeta"
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
			}
			// Virtual sortable mapped onto meta->>'short'.
			name: {
				sortableJSON: { json: "meta.short", accessor: "Meta.Short" }
				store: false
				goType: "string"
				omitSetter: true
				omitGetter: true
				envoy: { yaml: { omitEncoder: true } }
			}

			enabled: {
				sortable: true,
				goType: "bool"
				dal: { type: "Boolean", default: true }
			}

			scope: {
				goType: "*expr.Vars"
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
			}

			triggers: {
				goType: "types.NgAutomationTriggerSet"
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
			}
			steps: {
				goType: "types.NgAutomationStepSet"
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
			}
			paths: {
				goType: "types.NgAutomationPathSet"
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
			}
			issues: {
				goType: "types.NgAutomationIssueSet"
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

	// @todo
	envoy: {
		omit: true
	}

	filter: {
		struct: {
			automation_id: { goType: "[]string", ident: "automationID", storeIdent: "id" }
			tenant_id:     schema.TenantFilterField
			project_id:    schema.ProjectFilterField
			handle:        { goType: "string" }
			deleted: { goType: "filter.State", storeIdent: "deleted_at" }
			disabled: { goType: "filter.State", storeIdent: "enabled" }
		}

		query: ["handle"]
		byValue: ["automation_id", "handle"]
		byNilState: ["deleted"]
		byFalseState: ["disabled"]
	}

	rbac: {
		operations: {
			"read": description:            "Read automation"
			"update": description:          "Update automation"
			"delete": description:          "Delete automation"
			"undelete": description:        "Undelete automation"
			"execute": description:         "Execute automation"
		}
	}

	store: {
		ident: "automationNgAutomation"

		api: {
			lookups: [
				{
					fields: ["id"]
					description: """
						searches for automation by ID

						It returns automation even if deleted
						"""
				}, {
					fields: ["handle"]
					nullConstraint: ["deleted_at"]
					constraintCheck: true
					description: """
						searches for automation by their handle

						It returns only valid automations
						"""
				}
			]
		}
	}
}
