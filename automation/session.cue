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
			id:          schema.IdField
			workflow_id: { sortable: true, ident: "workflowID", goType: "uint64", storeIdent: "rel_workflow" }
			event_type: { sortable: true, goType: "string" }
			resource_type: { sortable: true, goType: "string" }
			status: { sortable: true, goType: "types.SessionStatus" }
			input: { goType: "*expr.Vars" }
			output: { goType: "*expr.Vars" }
			stacktrace: { goType: "types.Stacktrace" }

			created_by: { goType: "uint64" }
			created_at: schema.SortableTimestampNowField
			purge_at: schema.SortableTimestampNilField
			completed_at: schema.SortableTimestampNilField
			suspended_at: schema.SortableTimestampNilField
			error: {}
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
