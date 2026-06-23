package system

import (
	"github.com/crusttech/human/server/codegen/schema"
)

agent: {
	types: {
		gen:     true
		idLists: ["AgentAccessIDList"]
	}

	features: {
		labels: true
		projectScoped: true
	}

	model: {
		attributes: {
			id:     schema.IdField
			tenant_id:  schema.TenantRefField
			project_id: schema.ProjectRefField
			handle: schema.HandleField
			status: {
				sortable: true,
				goType: "string"
				dal: { length: 32 }
			}
			revision: {
				goType: "int"
				dal: { type: "Number", meta: { "rdbms:type": "integer" } }
			}
			meta: {
				goType: "types.AgentMeta"
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
			}
			// Virtual sortable mapped onto meta->>'short'.
			name: {
				sortableJSON: { json: "meta.short", accessor: "Meta.Short", nullable: false }
				store: false
				goType: "string"
				omitSetter: true
				omitGetter: true
				envoy: { yaml: { omitEncoder: true } }
			}
			behavior: {
				goType: "types.AgentBehavior"
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
			}
			execution: {
				goType: "types.AgentExecution"
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
			}
			access: {
				goType: "types.AgentAccess"
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
			}
			invocation: {
				goType: "types.AgentInvocation"
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
			}

			created_at: schema.SortableTimestampNowField
			updated_at: schema.SortableTimestampNilField
			deleted_at: schema.SortableTimestampNilField
			created_by: schema.AttributeUserRef & {
				// hand-written tag has no omitempty (record is always created)
				json: { field: "createdBy", string: true }
			}
			updated_by: schema.AttributeUserRef
			deleted_by: schema.AttributeUserRef
		}

		indexes: {
			"primary": { attribute: "id" }
			"unique_handle": {
			 fields: [{ attribute: "handle", modifiers: ["LOWERCASE"] }]
			 predicate: "handle != '' AND deleted_at IS NULL"
		 }
		}
	}

	filter: {
		struct: {
			agent_id: {goType: "[]uint64", ident: "agentID", storeIdent: "id"}
			tenant_id: schema.TenantFilterField
			project_id: schema.ProjectFilterField
			handle: {goType: "string"}
			status: {goType: "string"}
			deleted: {goType: "filter.State", storeIdent: "deleted_at"}
		}

		query: ["handle", "status"]
		byValue: ["agent_id", "handle", "status"]
		byNilState: ["deleted"]
	}

	rbac: {
		operations: {
			read: description:   "Read agent"
			update: description: "Update agent"
			delete: description: "Delete agent"
		}
	}

	service: {
		genAccessController: true

		events: false

		undelete: false

		customBodyOps: ["lookup", "search", "create", "update"]
	}

	refs: {
		items: [
			{path: "Behavior.KnowledgeBases",        kind: "KindKnowledgeBase",      reason: "ReasonAgentKnowledgeBase", iter: "sliceID"},
			{path: "Execution.Model.LLMProviderID",  kind: "KindLlmProvider",         reason: "ReasonAgentLlmProvider"},
			{path: "Access.TAQs",                    kind: "KindNgAutomation",         reason: "ReasonAgentAutomation",    iter: "sliceField"},
			{path: "Access.Workflows",               kind: "KindAutomationWorkflow",   reason: "ReasonAgentAutomation",    iter: "sliceField"},
		]
		extended: true
	}

	envoy: {
		omit: true
	}

	store: {
		api: {
			lookups: [
				{
					fields: ["id"]
					description: """
						searches for agent by ID

						It also returns deleted agents.
						"""
				}, {
					fields: ["handle"]
					nullConstraint: ["deleted_at"]
					constraintCheck: true
					description: """
						searches for agent by handle

						It returns only valid agents (not deleted)
						"""
				},
			]
		}
	}
}
