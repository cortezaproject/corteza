package automation

import (
	"github.com/crusttech/human/server/codegen/schema"
)

ng_automation: {
	features: {
		projectScoped: true
	}

	types: {
		gen: true
		// scope is *expr.Vars; the `expr` qualifier matches the path's last segment.
		imports: ["github.com/crusttech/human/server/pkg/expr"]
		// NgAutomationMeta has a custom ParseNgAutomationMeta (returns *NgAutomationMeta,
		// required by generated REST request decoders); keep hand-written Scan/Value/Parse.
		jsonTypesPtr: ["NgAutomationMeta"]
	}

	model: {
		attributes: {
			// convention derives "ngAutomationID" from res.ident; hand-written uses "automationID"
			id:         schema.IdField & {json: "automationID,string"}
			tenant_id:  schema.TenantRefField
			project_id: schema.ProjectRefField
			handle:     schema.HandleField
			meta: {
				goType: "*types.NgAutomationMeta"
				dal: {type: "JSON", defaultEmptyObject: true}
				omitSetter: true
				omitGetter: true
				json: {field: "meta", omitEmpty: true}
			}
			// Virtual sortable mapped onto meta->>'short'.
			name: {
				sortableJSON: {json: "meta.short", accessor: "Meta.Short"}
				store:      false
				goType:     "string"
				omitSetter: true
				omitGetter: true
				envoy: {yaml: {omitEncoder: true}}
			}

			enabled: {
				sortable: true
				goType:   "bool"
				dal: {type: "Boolean", default: true}
			}

			scope: {
				goType: "*expr.Vars"
				dal: {type: "JSON", defaultEmptyObject: true}
				omitSetter: true
				omitGetter: true
			}

			triggers: {
				goType: "types.NgAutomationTriggerSet"
				dal: {type: "JSON", defaultEmptyObject: true}
				omitSetter: true
				omitGetter: true
			}
			steps: {
				goType: "types.NgAutomationStepSet"
				dal: {type: "JSON", defaultEmptyObject: true}
				omitSetter: true
				omitGetter: true
			}
			paths: {
				goType: "types.NgAutomationPathSet"
				dal: {type: "JSON", defaultEmptyObject: true}
				omitSetter: true
				omitGetter: true
			}
			issues: {
				goType: "types.NgAutomationIssueSet"
				dal: {type: "JSON", defaultEmptyObject: true}
				omitSetter: true
				omitGetter: true
				json: {field: "issues", omitEmpty: true}
				envoy: {
					yaml: {
						omitEncoder: true
					}
				}
			}

			// user refs: convention adds ",omitempty"; hand-written omits it for runAs/ownedBy/createdBy
			run_as: schema.AttributeUserRef & {json: {field: "runAs", string: true}}

			owned_by:   schema.AttributeUserRef & {json: {field: "ownedBy", string: true}}
			created_at: schema.SortableTimestampNowField
			updated_at: schema.SortableTimestampNilField
			deleted_at: schema.SortableTimestampNilField
			created_by: schema.AttributeUserRef & {json: {field: "createdBy", string: true}}
			updated_by: schema.AttributeUserRef
			deleted_by: schema.AttributeUserRef
		}

		indexes: {
			"primary": {attribute: "id"}
		}
	}

	// @todo
	envoy: {
		omit: true
	}

	filter: {
		struct: {
			automation_id: {goType: "[]string", ident: "automationID", storeIdent: "id"}
			tenant_id:  schema.TenantFilterField
			project_id: schema.ProjectFilterField
			handle: {goType: "string"}
			deleted: {goType: "filter.State", storeIdent: "deleted_at"}
			disabled: {goType: "filter.State", storeIdent: "enabled"}
		}

		query: ["handle"]
		byValue: ["automation_id", "project_id", "handle"]
		byNilState: ["deleted"]
		byFalseState: ["disabled"]
	}

	rbac: {
		operations: {
			"read": description:     "Read automation"
			"update": description:   "Update automation"
			"delete": description:   "Delete automation"
			"undelete": description: "Undelete automation"
			"execute": description:  "Execute automation"
		}
	}

	refs: {
		extended: true
	}

	service: {
		events: false

		lookup:   false
		create:   false
		update:   false
		delete:   false
		undelete: false
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
				},
			]
		}
	}
}
