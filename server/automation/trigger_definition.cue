package automation

import (
	"github.com/cortezaproject/corteza/server/codegen/schema"
)

trigger_definition: {
	model: {
		attributes: {
			id: schema.IdField
			handle: schema.HandleField
			meta: {
				goType: "*types.TriggerDefinitionMeta"
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
			}

			input_schema: {
				goType: "types.TriggerDefinitionSchema"
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
			}

			output_schema: {
				goType: "types.TriggerDefinitionSchema"
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
			}

			skip_event_bus: {
				goType: "bool"
				dal: { type: "Boolean", default: false }
			}

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

	envoy: {
		omit: true
	}

	filter: {
		struct: {
			trigger_definition_id: { goType: "[]string", ident: "triggerDefinitionID", storeIdent: "id" }
			handle: { goType: "string" }
			query: { goType: "string" }
			deleted: { goType: "filter.State", storeIdent: "deleted_at" }
		}

		query: ["handle"]
		byValue: ["trigger_definition_id", "handle"]
		byNilState: ["deleted"]
	}

	rbac: {
		operations: {
			"read": description:            "Read trigger definition"
			"update": description:          "Update trigger definition"
			"delete": description:          "Delete trigger definition"
			"undelete": description:        "Undelete trigger definition"
		}
	}

	store: {
		ident: "automationTriggerDefinition"

		api: {
			lookups: [
				{
					fields: ["id"]
					description: """
						searches for trigger definition by ID

						It returns trigger definition even if deleted
						"""
				}, {
					fields: ["handle"]
					nullConstraint: ["deleted_at"]
					constraintCheck: true
					description: """
						searches for trigger definition by their handle

						It returns only valid trigger definitions
						"""
				}
			]
		}
	}
}
