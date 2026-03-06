package system

import (
	"github.com/cortezaproject/corteza/server/codegen/schema"
)

agent: {
	features: {
		labels: false
	}

	model: {
		attributes: {
			id:     schema.IdField
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
			created_by: schema.AttributeUserRef
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
