package system

import (
	"github.com/crusttech/human/server/codegen/schema"
)

queue: {
	features: {
		labels: false
		projectScoped: true
	}

	types: { gen: true }

	model: {
		ident: "queue_settings"
		attributes: {
			id: schema.IdField
			tenant_id:  schema.TenantRefField
			project_id: schema.ProjectRefField
			consumer: {
				sortable: true,
				goType: "string"
				dal: {}
			}
			queue: {
				sortable: true,
				goType: "string"
				dal: {}
				envoy: {
					identifier: true
				}
			}
			meta: {
				goType: "types.QueueMeta"
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
			}

			created_at: schema.SortableTimestampNowField
			updated_at: schema.SortableTimestampNilField
			deleted_at: schema.SortableTimestampNilField
			created_by: schema.AttributeUserRef & {
				json: {field: "createdBy", string: true}
			}
			updated_by: schema.AttributeUserRef
			deleted_by: schema.AttributeUserRef
		}

		indexes: {
			"primary": { attribute: "id" }
		}
	}

	filter: {
		struct: {
			queue_id: {goType: "[]uint64", ident: "queueID", storeIdent: "id"}
			tenant_id: schema.TenantFilterField
			project_id: schema.ProjectFilterField
			query: {goType: "string"}
			deleted: {goType: "filter.State", storeIdent: "deleted_at"}
		}

		query: ["queue", "consumer"]
		byValue: ["queue_id"]
		byNilState: ["deleted"]
	}

	envoy: {
		yaml: {
			supportMappedInput: true
			mappedField: "Queue"
			identKeyAlias: []
		}
		store: {
			handleField: ""
		}
	}

	rbac: {
		operations: {
			"read": description:        "Read queue"
			"update": description:      "Update queue"
			"delete": description:      "Delete queue"
			"queue.read": description:  "Read from queue"
			"queue.write": description: "Write to queue"
		}
	}

	service: {
		events:   false
		undelete: true

		genAccessController: true
		genConstructor:      true

		// the action-log prop for the filter is named "search" (not the default
		// "filter"); the standard Search body is otherwise generated.
		filterProp: "search"

		updateFields: ["Consumer", "Queue", "Meta"]

		hooks: {
			beforeCreate:  true
			afterCreate:   true
			beforeUpdate:  true
			afterUpdate:   true
			afterDelete:   true
			afterUndelete: true
		}
	}

	store: {

		api: {
			lookups: [
				{
					fields: ["id"]
					description: """
						searches for queue by ID
						"""
				}, {
					fields: ["queue"]
					description: """
						searches for queue by queue name
						"""
				},
			]
		}
	}
}
