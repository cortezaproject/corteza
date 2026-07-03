package automation

import (
	"github.com/crusttech/human/server/codegen/schema"
)

_ng_automationDefs: {
			NgAutomationMeta: { name: "NgAutomationMeta", fields: [
				{name: "Short", type: "string", json: "short"},
				{name: "Description", type: "string", json: "description"},
				{name: "Visual", type: _ng_automationDefs.NgAutomationVisual, json: "visual"},
				{name: "Icon", type: _ng_automationDefs.NgAutomationIcon, ptr: true, json: "icon,omitempty"},
			]}
			NgAutomationVisual: { name: "NgAutomationVisual", fields: []}
			NgAutomationIcon: { name: "NgAutomationIcon", fields: [
				{name: "Type", type: "string", json: "type"},
				{name: "Value", type: "string", json: "value"},
			]}

			NgAutomationTriggerSet: { name: "NgAutomationTriggerSet", elem: _ng_automationDefs.NgAutomationTrigger, elemPtr: true }
			NgAutomationStepSet: { name: "NgAutomationStepSet", elem: _ng_automationDefs.NgAutomationStep, elemPtr: true }
			NgAutomationPathSet: { name: "NgAutomationPathSet", elem: _ng_automationDefs.NgAutomationPath, elemPtr: true }
			NgAutomationIssueSet: { name: "NgAutomationIssueSet", elem: _ng_automationDefs.NgAutomationIssue, elemPtr: true }

			NgAutomationTriggerSchema: { name: "NgAutomationTriggerSchema", elem: _ng_automationDefs.NgAutomationTriggerParam, elemPtr: false }

			NgAutomationTrigger: { name: "NgAutomationTrigger", fields: [
				{name: "ID", type: "uint64", json: "triggerID,string"},
				{name: "Labels", goType: "map[string]labelTypes.LabelValue", json: "labels,omitempty"},
				{name: "Handle", type: "string", json: "handle"},
				{name: "Meta", type: _ng_automationDefs.NgTriggerMeta, ptr: true, json: "meta,omitempty"},
				{name: "Enabled", type: "bool", json: "enabled"},
				{name: "ResourceType", type: "string", json: "resourceType"},
				{name: "EventType", type: "string", json: "eventType"},
				{name: "Constraints", type: _ng_automationDefs.NgTriggerConstraint, slice: true, json: "constraints"},
				{name: "Input", goType: "*expr.Vars", json: "input"},
				{name: "InputSchema", type: _ng_automationDefs.NgAutomationTriggerSchema, json: "inputSchema,omitempty"},
			]}

			NgAutomationTriggerParam: { name: "NgAutomationTriggerParam", fields: [
				{name: "Name", type: "string", json: "name"},
				{name: "Type", type: "string", json: "type"},
				{name: "Required", type: "bool", json: "required"},
				{name: "Description", type: "string", json: "description"},
			]}

			NgTriggerConstraint: { name: "NgTriggerConstraint", fields: [
				{name: "Name", type: "string", json: "name"},
				{name: "Op", type: "string", json: "op,omitempty"},
				{name: "Values", type: _ng_automationDefs.NgTriggerConstraintValue, slice: true, json: "values,omitempty"},
			]}

			NgTriggerConstraintValue: { name: "NgTriggerConstraintValue", fields: [
				{name: "Type", type: "string", json: "@type"},
				{name: "Value", type: "string", json: "@value"},
			]}

			NgTriggerMeta: { name: "NgTriggerMeta", fields: [
				{name: "Short", type: "string", json: "short"},
				{name: "Description", type: "string", json: "description"},
				{name: "Visual", type: _ng_automationDefs.NgAutomationVisual, json: "visual"},
				{name: "Icon", type: _ng_automationDefs.NgAutomationIcon, ptr: true, json: "icon,omitempty"},
			]}

			NgAutomationStep: { name: "NgAutomationStep", fields: [
				{name: "ID", type: "uint64", json: "stepID,string"},
				{name: "Handle", type: "string", json: "handle"},
				{name: "Meta", type: _ng_automationDefs.NgAutomationStepMeta, json: "meta"},
				{name: "Kind", type: "string", json: "kind"},
				{name: "Ref", type: "string", json: "ref"},
				{name: "Arguments", goType: "[]*Expr", json: "arguments"},
				{name: "Results", goType: "[]*Expr", json: "results"},
				{name: "Recoverable", type: "bool", json: "recoverable,omitempty"},
				{name: "MaxRetries", type: "int", json: "maxRetries,omitempty"},
			]}

			NgAutomationStepMeta: { name: "NgAutomationStepMeta", fields: [
				{name: "Short", type: "string", json: "short"},
				{name: "Description", type: "string", json: "description"},
				{name: "Visual", type: _ng_automationDefs.NgAutomationVisual, json: "visual"},
				{name: "Icon", type: _ng_automationDefs.NgAutomationIcon, ptr: true, json: "icon,omitempty"},
				{name: "Extra", goType: "map[string]any", json: "extra,omitempty"},
			]}

			NgAutomationPath: { name: "NgAutomationPath", fields: [
				{name: "ParentID", type: "uint64", json: "parentID,string"},
				{name: "ChildID", type: "uint64", json: "childID,string"},
				{name: "Condition", goType: "*ast.ASTNode", json: "condition,omitempty"},
				{name: "Kind", type: "string", json: "kind,omitempty"},
				{name: "Meta", type: _ng_automationDefs.NgAutomationPathMeta, json: "meta"},
			]}

			NgAutomationPathMeta: { name: "NgAutomationPathMeta", fields: [
				{name: "Short", type: "string", json: "short"},
				{name: "Description", type: "string", json: "description"},
				{name: "Visual", type: _ng_automationDefs.NgAutomationVisual, json: "visual"},
				{name: "Icon", type: _ng_automationDefs.NgAutomationIcon, ptr: true, json: "icon,omitempty"},
			]}

			NgAutomationIssue: { name: "NgAutomationIssue", fields: [
				{name: "Code", type: "string", json: "code"},
				{name: "Severity", type: "string", json: "severity"},
				{name: "Message", type: "string", json: "message"},
				{name: "Details", goType: "[]*NgAutomationIssueDetail", json: "details,omitempty"},
			]}

			NgAutomationIssueDetail: { name: "NgAutomationIssueDetail", fields: [
				{name: "MissingReference", type: _ng_automationDefs.DetailMissingReference, ptr: true},
				{name: "InvalidType", type: _ng_automationDefs.DetailInvalidType, ptr: true},
				{name: "DuplicateID", type: _ng_automationDefs.DetailDuplicateID, ptr: true},
				{name: "Cycle", type: _ng_automationDefs.DetailCycle, ptr: true},
				{name: "ResourceRef", type: _ng_automationDefs.DetailResourceRef, ptr: true},
				{name: "GatewayPaths", type: _ng_automationDefs.DetailGatewayPaths, ptr: true},
				{name: "EmptyField", type: _ng_automationDefs.DetailEmptyField, ptr: true},
				{name: "Details", goType: "[]*NgAutomationIssueDetail"},
			]}

			DetailMissingReference: { name: "DetailMissingReference", fields: [
				{name: "RefKind", type: "string", json: "refKind"},
				{name: "Ref", type: "string", json: "ref"},
				{name: "StepID", type: "uint64", json: "stepID,string,omitempty"},
				{name: "Field", type: "string", json: "field,omitempty"},
				{name: "FieldIndex", type: "int", json: "fieldIndex,omitempty"},
			]}

			DetailInvalidType: { name: "DetailInvalidType", fields: [
				{name: "StepID", type: "uint64", json: "stepID,string"},
				{name: "Field", type: "string", json: "field,omitempty"},
				{name: "FieldIndex", type: "int", json: "fieldIndex,omitempty"},
				{name: "Target", type: "string", json: "target,omitempty"},
				{name: "Expected", type: "string", json: "expected"},
				{name: "Actual", type: "string", json: "actual"},
			]}

			DetailDuplicateID: { name: "DetailDuplicateID", fields: [
				{name: "Resource", type: "string", json: "resource"},
				{name: "ID", type: "uint64", json: "id,string"},
				{name: "Indices", slice: true, type: "int", json: "indices,omitempty"},
			]}

			DetailCycle: { name: "DetailCycle", fields: [
				{name: "StepIDs", goType: "IssueIDs", json: "stepIDs"},
			]}

			DetailResourceRef: { name: "DetailResourceRef", fields: [
				{name: "Resource", type: "string", json: "resource"},
				{name: "ID", type: "uint64", json: "id,string,omitempty"},
				{name: "Index", type: "int", json: "index,omitempty"},
				{name: "Field", type: "string", json: "field,omitempty"},
				{name: "FieldIndex", type: "int", json: "fieldIndex,omitempty"},
			]}

			DetailGatewayPaths: { name: "DetailGatewayPaths", fields: [
				{name: "StepID", type: "uint64", json: "stepID,string"},
				{name: "Violation", type: "string", json: "violation"},
				{name: "Got", type: "int", json: "got"},
				{name: "Want", type: "int", json: "want,omitempty"},
			]}

			DetailEmptyField: { name: "DetailEmptyField", fields: [
				{name: "Resource", type: "string", json: "resource"},
				{name: "Index", type: "int", json: "index"},
				{name: "Field", type: "string", json: "field,omitempty"},
			]}
		}

ng_automation: {
	features: {
	}

	types: {
		gen: true
		// scope is *expr.Vars; the `expr` qualifier matches the path's last segment.
		imports: [
			"github.com/crusttech/human/server/pkg/expr",
			"github.com/crusttech/human/server/pkg/ast",
			"labelTypes \"github.com/crusttech/human/server/pkg/label/types\"",
		]
		// NgAutomationMeta has a custom ParseNgAutomationMeta (returns *NgAutomationMeta,
		// required by generated REST request decoders); keep hand-written Scan/Value/Parse.
		jsonTypesPtr: ["NgAutomationMeta"]

		defs: _ng_automationDefs
	}

	model: {
		attributes: {
			// convention derives "ngAutomationID" from res.ident; hand-written uses "automationID"
			id:         schema.IdField & {json: "automationID,string"}
			tenant_id:  schema.TenantRefField
			project_id: schema.ProjectRefField
			handle:     schema.HandleField
			meta: {
				type: _ng_automationDefs.NgAutomationMeta
				ptr:  true
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
				type: _ng_automationDefs.NgAutomationTriggerSet
				dal: {type: "JSON", defaultEmptyObject: true}
				omitSetter: true
				omitGetter: true
			}
			steps: {
				type: _ng_automationDefs.NgAutomationStepSet
				dal: {type: "JSON", defaultEmptyObject: true}
				omitSetter: true
				omitGetter: true
			}
			paths: {
				type: _ng_automationDefs.NgAutomationPathSet
				dal: {type: "JSON", defaultEmptyObject: true}
				omitSetter: true
				omitGetter: true
			}
			issues: {
				type: _ng_automationDefs.NgAutomationIssueSet
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
		customFunctions: [
			{
				name: "Exec"
				cap:  "write"
				args: [
					{name: "automationID", goType: "uint64"},
					{name: "p", goType: "types.NgAutomationExecParams"},
				]
				results: [
					{name: "executionID", goType: "id.ID"},
					{name: "err", goType: "error"},
				]
			},
			{
				name: "ExecAndWait"
				cap:  "write"
				args: [
					{name: "automationID", goType: "uint64"},
					{name: "p", goType: "types.NgAutomationExecParams"},
				]
				results: [
					{name: "out", goType: "*execTypes.ExecutionResult"},
					{name: "err", goType: "error"},
				]
			},
			{
				name: "GetExecutions"
				cap:  "read"
				args: [
					{name: "automationID", goType: "uint64"},
				]
				results: [
					{name: "out", goType: "[]*execTypes.ExecutionResult"},
					{name: "err", goType: "error"},
				]
			},
			{
				name: "GetExecutionTrace"
				cap:  "read"
				args: [
					{name: "exeID", goType: "uint64"},
					{name: "executionID", goType: "uint64"},
					{name: "rev", goType: "int"},
				]
				results: [
					{name: "out", goType: "[]execTypes.StackFrame"},
					{name: "err", goType: "error"},
				]
			},
			{
				name: "GetAllExecutions"
				cap:  "read"
				args: [
					{name: "f", goType: "execTypes.ExecutionFilter"},
				]
				results: [
					{name: "out", goType: "[]*execTypes.ExecutionResult"},
					{name: "err", goType: "error"},
				]
			},
			{
				name: "Load"
				cap:  "write"
			},
		]
		customFunctionImports: [
			"\"github.com/crusttech/human/server/pkg/id\"",
			"execTypes \"github.com/crusttech/human/server/pkg/automation_exec/types\"",
		]
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
