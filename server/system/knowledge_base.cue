package system

import (
	"github.com/crusttech/human/server/codegen/schema"
)

knowledge_base: {
	features: {
		labels: false
		projectScoped: true
	}

	types: { gen: true, idLists: ["KnowledgeBaseIDList"] }

	model: {
		attributes: {
			id:     schema.IdField
			tenant_id:  schema.TenantRefField
			project_id: schema.ProjectRefField
			handle: schema.HandleField
			title: {
				goType: "string"
				dal: { length: 512 }
			}
			description: {
				goType: "string"
				dal: { type: "Text" }
				// hand-written tag is `description,omitempty`; convention (plain
				// string) would drop omitempty
				json: { field: "description", omitEmpty: true }
			}
			context: {
				goType: "*types.KnowledgeBaseContext"
				dal: { type: "JSON", nullable: true }
				omitSetter: true
				omitGetter: true
				// hand-written tag is `context,omitempty`; convention (non
				// ID/ref/time) would drop omitempty
				json: { field: "context", omitEmpty: true }
			}

			created_at: schema.SortableTimestampNowField
			updated_at: schema.SortableTimestampNilField
			deleted_at: schema.SortableTimestampNilField
			// hand-written tag is `createdBy,string` (NO omitempty); the ref
			// convention would emit `createdBy,string,omitempty`
			created_by: schema.AttributeUserRef & { json: { field: "createdBy", "string": true } }
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
			knowledge_base_id: {goType: "[]uint64", ident: "knowledgeBaseID", storeIdent: "id"}
			tenant_id: schema.TenantFilterField
			project_id: schema.ProjectFilterField
			handle: {goType: "string"}
			deleted: {goType: "filter.State", storeIdent: "deleted_at"}
		}

		query: ["handle", "title"]
		byValue: ["knowledge_base_id", "handle"]
		byNilState: ["deleted"]
	}

	rbac: {
		operations: {
			read: description:   "Read knowledge base"
			update: description: "Update knowledge base"
			delete: description: "Delete knowledge base"
		}
	}

	refs: {
		extended: true
	}

	service: {
		genAccessController: true
		genConstructor:      true

		// no eventbus events: the original service does not emit any and the
		// struct carries no eventbus dependency
		events: false

		undelete: true

		// the action-log filter prop is named "search" (not the default "filter")
		filterProp: "search"

		// Update is bespoke: access is checked on the incoming `upd` (before the
		// existing record is loaded), the existing audit fields are copied ONTO
		// `upd`, and `upd` (not the loaded record) is what gets persisted.
		// Undelete is bespoke too: it sets UpdatedAt/UpdatedBy and there is no
		// dedicated notAllowedToUndelete error (it reuses notAllowedToDelete).
		customBodyOps: ["update", "undelete"]

		// audit-author bookkeeping (CreatedBy / DeletedBy) the standard scaffold
		// does not emit lives in these hooks
		hooks: {
			beforeCreate: true
			beforeDelete: true
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
						searches for knowledge base by ID

						It also returns deleted knowledge bases.
						"""
				}, {
					fields: ["handle"]
					nullConstraint: ["deleted_at"]
					constraintCheck: true
					description: """
						searches for knowledge base by handle

						It returns only valid knowledge bases (not deleted)
						"""
				},
			]
		}
	}
}
