package automation

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

session: {
	features: {
		labels: false
	}

	model: {
		ident: "automation_sessions"
 		attributes: {
			id: schema.IdField
			workflow_id: {
				sortable: true,
				ident: "workflowID",
				goType: "uint64",
				storeIdent: "rel_workflow"
				dal: { type: "Ref", refModelResType: "corteza::automation:workflow" }
			}
			status: {
				sortable: true,
				goType: "types.SessionStatus"
				dal: { type: "Number", default: 0 }
			}
			event_type: {
				sortable: true,
				goType: "string"
				dal: {}
			}
			resource_type: {
				sortable: true,
				goType: "string"
				dal: {}
			}
			input: {
				goType: "*expr.Vars"
				dal: { type: "JSON", defaultEmptyObject: true }
			}
			output: {
				goType: "*expr.Vars"
				dal: { type: "JSON", defaultEmptyObject: true }
			}
			stacktrace: {
				goType: "types.Stacktrace"
				dal: { type: "JSON", defaultEmptyObject: true }
			}

			created_by: schema.AttributeUserRef
			created_at: schema.SortableTimestampNowField
			purge_at: schema.SortableTimestampNilField
			suspended_at: schema.SortableTimestampNilField
			completed_at: schema.SortableTimestampNilField
			error: {
				dal: {}
			}
		}

		indexes: {
			"primary": { attribute: "id" }
			"completed_at": { attribute: "completed_at" }
			"created_at": { attribute: "created_at" }
			"event_type": { attribute: "event_type" }
			"resource_type": { attribute: "resource_type" }
			"status": { attribute: "status" }
			"suspended_at": { attribute: "suspended_at" }
			"resource_type": { attribute: "resource_type" }
		}
	}

	filter: {
		struct: {
			session_id: { goType: "[]uint64", storeIdent: "id", ident: "sessionID" }
			completed: { schema.SortableTimestampNilField, storeIdent: "completed_at" }
			created_by: { goType: "[]uint64" }
			status: { goType: "[]uint" }
			workflow_id: { goType: "[]uint64", storeIdent: "rel_workflow", ident: "workflowID" }
			event_type: { goType: "string" }
			resource_type: { goType: "string" }
		}

		byValue: ["status", "workflow_id", "event_type", "resource_type", "created_by"]
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
				}
			]
		}
	}
}
