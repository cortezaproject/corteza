package automation

import (
	"github.com/crusttech/human/server/codegen/schema"
)

_sessionDefs: {
	SessionStatus: {
		name: "SessionStatus"
		kind: "int"
		values: [
			{ident: "SessionStarted", value:   0},
			{ident: "SessionPrompted", value:  1},
			{ident: "SessionSuspended", value: 2},
			{ident: "SessionFailed", value:    3},
			{ident: "SessionCompleted", value: 4},
			{ident: "SessionCanceled", value:  5},
		]
	}
	Stacktrace: {name: "Stacktrace", elemGoType: "wfexec.Frame", elemPtr: true}
	SessionStartParams: {name: "SessionStartParams", fields: [
					{name: "Invoker", goType:    "auth.Identifiable"},
					{name: "Runner", goType:     "auth.Identifiable"},
					{name: "WorkflowID", type:   "uint64"},
					{name: "KeepFor", type:      "int"},
					{name: "Trace", type:        "bool"},
					{name: "Input", goType:      "*expr.Vars"},
					{name: "StepID", type:       "uint64"},
					{name: "EventType", type:    "string"},
					{name: "ResourceType", type: "string"},
					{name: "CallStack", slice:   true, type: "uint64"},
	]}
}

session: {
	features: {
		labels: false
	}

	model: {
		// length for the lang is now a bit shorter
		// Reason for that is supported index length in MySQL
		ident: "automation_sessions"
		attributes: {
			id:         schema.IdField
			tenant_id:  schema.TenantRefField
			project_id: schema.ProjectRefField
			workflow_id: {
				sortable:   true
				ident:      "workflowID"
				goType:     "uint64"
				storeIdent: "rel_workflow"
				dal: {type: "Ref", refModelResType: "corteza::automation:workflow"}
			}
			status: {
				sortable: true
				type:     _sessionDefs.SessionStatus
				dal: {type: "Number", default: 0, meta: {"rdbms:type": "integer"}}
				omitSetter: true
				omitGetter: true
			}
			event_type: {
				sortable: true
				goType:   "string"
				dal: {length: 32}
			}
			resource_type: {
				sortable: true
				goType:   "string"
				dal: {length: 64}
			}
			input: {
				goType: "*expr.Vars"
				dal: {type: "JSON", defaultEmptyObject: true}
				omitSetter: true
				omitGetter: true
			}
			output: {
				goType: "*expr.Vars"
				dal: {type: "JSON", defaultEmptyObject: true}
				omitSetter: true
				omitGetter: true
			}
			stacktrace: {
				type: _sessionDefs.Stacktrace
				dal: {type: "JSON", defaultEmptyObject: true}
				omitSetter: true
				omitGetter: true
			}

			created_by:   schema.AttributeUserRef
			created_at:   schema.SortableTimestampNowField
			purge_at:     schema.SortableTimestampNilField
			suspended_at: schema.SortableTimestampNilField
			completed_at: schema.SortableTimestampNilField
			error: {
				dal: {}
			}
		}

		indexes: {
			"primary": {attribute: "id"}
			"completed_at": {attribute: "completed_at"}
			"created_at": {attribute: "created_at"}
			"event_type": {attribute: "event_type"}
			"resource_type": {attribute: "resource_type"}
			"status": {attribute: "status"}
			"suspended_at": {attribute: "suspended_at"}
			"resource_type": {attribute: "resource_type"}
		}
	}

	types: {
		defs: _sessionDefs
		imports: [
			"github.com/crusttech/human/server/pkg/auth",
			"github.com/crusttech/human/server/pkg/expr",
		]
	}

	envoy: {
		omit: true
	}

	filter: {
		struct: {
			session_id: {goType: "[]string", storeIdent: "id", ident: "sessionID"}
			tenant_id:  schema.TenantFilterField
			project_id: schema.ProjectFilterField
			completed: {schema.SortableTimestampNilField, storeIdent: "completed_at"}
			created_by: {goType: "[]string"}
			status: {goType: "[]uint"}
			workflow_id: {goType: "[]string", storeIdent: "rel_workflow", ident: "workflowID"}
			event_type: {goType: "string"}
			resource_type: {goType: "string"}
		}

		byValue: ["status", "session_id", "workflow_id", "event_type", "resource_type", "created_by"]
		byNilState: ["completed"]
	}

	store: {
		ident: "automationSession"

		api: {
			lookups: [
				{
					fields: ["id"]
					description: """
						searches for session by ID

						It returns session even if deleted
						"""
				},
			]
		}
	}
}
