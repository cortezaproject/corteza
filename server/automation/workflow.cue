package automation

import (
	"github.com/crusttech/human/server/codegen/schema"
)

_workflowDefs: {
	WorkflowMeta: {name: "WorkflowMeta", fields: [
				{name: "Name", type:        "string", json:                 "name"},
				{name: "Description", type: "string", json:                 "description"},
				{name: "Visual", goType:    "map[string]interface{}", json: "visual"},
				{name: "SubWorkflow", type: "bool", json:                   "subWorkflow,omitempty"},

				// Named input/output contract, read by the "Run Workflow" step
				{name: "Input", goType:  "[]WorkflowIODef", json: "input,omitempty"},
				{name: "Output", goType: "[]WorkflowIODef", json: "output,omitempty"},
	]}

	WorkflowStepSet: {name: "WorkflowStepSet", elem: _workflowDefs.WorkflowStep, elemPtr: true}
	WorkflowPathSet: {name: "WorkflowPathSet", elem: _workflowDefs.WorkflowPath, elemPtr: true}
	WorkflowIssueSet: {name: "WorkflowIssueSet", elem: _workflowDefs.WorkflowIssue, elemPtr: true}

	WorkflowStepKind: {name: "WorkflowStepKind", kind: "string", values: [
					{ident:                   "WorkflowStepKindExpressions", value:  "expressions"},
					{ident:                   "WorkflowStepKindGateway", value:      "gateway"},
					{ident:                   "WorkflowStepKindFunction", value:     "function"},
					{ident:                   "WorkflowStepKindIterator", value:     "iterator"},
					{ident:                   "WorkflowStepKindError", value:        "error"},
					{ident:                   "WorkflowStepKindTermination", value:  "termination"},
					{ident:                   "WorkflowStepKindPrompt", value:       "prompt"},
					{ident:                   "WorkflowStepKindDelay", value:        "delay"},
					{ident:                   "WorkflowStepKindErrHandler", value:   "error-handler"},
					{ident:                   "WorkflowStepKindVisual", value:       "visual"},
					{ident:                   "WorkflowStepKindDebug", value:        "debug"},
					{ident:                   "WorkflowStepKindBreak", value:        "break"},
					{ident:                   "WorkflowStepKindContinue", value:     "continue"},
					{ident:                   "WorkflowStepKindExecWorkflow", value: "exec-workflow"},
	]}

	WorkflowStep: {name: "WorkflowStep", fields: [
				{name: "ID", type:          "uint64", json:                       "stepID,string"},
				{name: "Kind", type:        _workflowDefs.WorkflowStepKind, json: "kind"},
				{name: "Ref", type:         "string", json:                       "ref"},
				{name: "Arguments", goType: "[]*Expr", json:                      "arguments"},
				{name: "Results", goType:   "[]*Expr", json:                      "results"},
				{name: "Meta", type:        _workflowDefs.WorkflowStepMeta, json: "meta,omitempty"},
				{name: "Labels", goType:    "map[string]string", json:            "labels,omitempty"},
	]}

	WorkflowStepMeta: {name: "WorkflowStepMeta", fields: [
					{name: "Name", type:        "string", json:                 "name"},
					{name: "Description", type: "string", json:                 "description"},
					{name: "Visual", goType:    "map[string]interface{}", json: "visual"},
	]}

	WorkflowPath: {name: "WorkflowPath", fields: [
				{name: "Expr", type:     "string", json:                       "expr,omitempty"},
				{name: "ParentID", type: "uint64", json:                       "parentID,string"},
				{name: "ChildID", type:  "uint64", json:                       "childID,string"},
				{name: "Meta", type:     _workflowDefs.WorkflowPathMeta, json: "meta,omitempty"},
	]}

	WorkflowPathMeta: {name: "WorkflowPathMeta", fields: [
					{name: "Name", type:        "string", json:                 "name"},
					{name: "Description", type: "string", json:                 "description"},
					{name: "Visual", goType:    "map[string]interface{}", json: "visual"},
	]}

	WorkflowIssue: {name: "WorkflowIssue", fields: [
				{name: "Culprit", goType:   "map[string]int", json: "culprit"},
				{name: "Description", type: "string", json:         "description"},
	]}

	WorkflowExecParams: {name: "WorkflowExecParams", fields: [
					{name: "CallerWorkflowID", type: "uint64"},
					{name: "CallerSessionID", type:  "uint64"},
					{name: "CallerStepID", type:     "uint64"},
					{name: "StepID", type:           "uint64"},
					{name: "EventType", type:        "string"},
					{name: "ResourceType", type:     "string"},
					{name: "Trace", type:            "bool"},
					{name: "Async", type:            "bool"},
					{name: "Wait", type:             "bool"},
					{name: "Input", goType:          "*expr.Vars"},
	]}
}

workflow: {
	features: {
		labelResourceType: "workflow"
	}
	types: {
		gen: true
		// WorkflowMeta keeps a hand-written ParseWorkflowMeta that returns
		// *WorkflowMeta (REST request controllers depend on the pointer return),
		// which is incompatible with the value-returning generated version, so
		// it is excluded from JSON helper generation.
		jsonTypesPtr: ["WorkflowMeta"]
		// WorkflowPath has an unexported eval field — hand-written; skip struct gen.
		structTypesSkip: ["WorkflowPath"]
		imports: [
			"github.com/crusttech/human/server/pkg/expr",
		]
		defs: _workflowDefs
	}
	model: {
		ident: "automation_workflows"
		attributes: {
			id:         schema.IdField
			tenant_id:  schema.TenantRefField
			project_id: schema.ProjectRefField
			handle:     schema.HandleField
			meta: {
				type: _workflowDefs.WorkflowMeta
				ptr:  true
				dal: {type: "JSON", defaultEmptyObject: true}
				json: {field: "meta", omitEmpty: true}
				omitSetter: true
				omitGetter: true
			}
			// Virtual sortable mapped onto meta->>'name'.
			name: {
				sortableJSON: {json: "meta.name", accessor: "Meta.Name"}
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
			trace: {
				goType: "bool"
				dal: {type: "Boolean", default: false}
			}
			keep_sessions: {
				goType: "int"
				dal: {type: "Number", default: 0, meta: {"rdbms:type": "integer"}}
			}
			scope: {
				goType: "*expr.Vars"
				dal: {type: "JSON", defaultEmptyObject: true}
				omitSetter: true
				omitGetter: true
			}
			steps: {
				type: _workflowDefs.WorkflowStepSet
				dal: {type: "JSON", defaultEmptyObject: true}
				omitSetter: true
				omitGetter: true
			}
			paths: {
				type: _workflowDefs.WorkflowPathSet
				dal: {type: "JSON", defaultEmptyObject: true}
				omitSetter: true
				omitGetter: true
			}
			issues: {
				type: _workflowDefs.WorkflowIssueSet
				dal: {type: "JSON", defaultEmptyObject: true}
				json: {field: "issues", omitEmpty: true}
				omitSetter: true
				omitGetter: true
				envoy: {
					yaml: {
						omitEncoder: true
					}
				}
			}

			run_as: schema.AttributeUserRef & {json: {field: "runAs", "string": true}}

			owned_by:   schema.AttributeUserRef & {json: {field: "ownedBy", "string": true}}
			created_at: schema.SortableTimestampNowField
			updated_at: schema.SortableTimestampNilField
			deleted_at: schema.SortableTimestampNilField
			created_by: schema.AttributeUserRef & {json: {field: "createdBy", "string": true}}
			updated_by: schema.AttributeUserRef
			deleted_by: schema.AttributeUserRef
		}

		indexes: {
			"primary": {attribute: "id"}
		}
	}

	envoy: {
		yaml: {
			supportMappedInput: true
			mappedField:        "Handle"
			identKeyAlias: ["workflows"]
			extendedResourceDecoders: [{
				ident:              "triggers"
				expIdent:           "Triggers"
				supportMappedInput: false
				identKeys: ["triggers"]
			}]
			extendedResourceEncoders: [{
				ident:    "trigger"
				expIdent: "Trigger"
				identKey: "trigger"
			}]
		}
		store: {
			customFilterBuilder: true
			extendedDecoder:     true
		}
	}

	filter: {
		struct: {
			workflow_id: {goType: "[]string", ident: "workflowID", storeIdent: "id"}
			tenant_id:  schema.TenantFilterField
			project_id: schema.ProjectFilterField
			handle: {goType: "string"}
			sub_workflow: {goType: "filter.State"}
			deleted: {goType: "filter.State", storeIdent: "deleted_at"}
			disabled: {goType: "filter.State", storeIdent: "enabled"}
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

	refs: {
		extended: true
	}

	service: {
		customFunctions: [
			{
				name:   "Exec"
				cap:    "write"
				action: "Execute"
				args: [
					{name: "workflowID", goType: "uint64"},
					{name: "p", goType:          "types.WorkflowExecParams"},
				]
				results: [
					{name: "results", goType:    "*expr.Vars"},
					{name: "sessionID", goType:  "uint64"},
					{name: "stacktrace", goType: "types.Stacktrace"},
					{name: "err", goType:        "error"},
				]
			},
		]
		customFunctionImports: [
			"\"github.com/crusttech/human/server/pkg/expr\"",
		]
		lookup:   false
		create:   false
		update:   false
		delete:   false
		undelete: false
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
				},
			]
		}
	}
}
