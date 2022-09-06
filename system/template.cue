package system

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

template: {
	model: {
		attributes: {
				id:     schema.IdField
				handle: schema.HandleField
				language: {sortable: true, goType: "string"}
				type: {sortable: true, goType: "types.DocumentType"}
				partial: {goType: "bool"}
				meta: {goType: "types.TemplateMeta"}
				template: {sortable: true, goType: "string"}

				owner_id: { ident: "ownerID", goType: "uint64", storeIdent: "rel_owner" }
				created_at: schema.SortableTimestampNowField
				updated_at: schema.SortableTimestampNilField
				deleted_at: schema.SortableTimestampNilField
				last_used_at: schema.SortableTimestampNilField
		}
	}

	filter: {
		struct: {
			template_id: {goType: "[]uint64", ident: "templateID", storeIdent: "id"}
			handle: {goType: "string"}
			type: {goType: "string"}
			owner_id: {goType: "uint64", storeIdent: "rel_owner", ident: "ownerID" }
			partial: {goType: "bool"}
			deleted: {goType: "filter.State", storeIdent: "deleted_at"}
		}

		query: ["handle", "type"]
		byValue: ["template_id", "handle", "partial", "type", "owner_id"]
		byNilState: ["deleted"]
	}

	rbac: {
		operations: {
			read: description:   "Read template"
			update: description: "Update template"
			delete: description: "Delete template"
			render: description: "Render template"
		}
	}

	store: {
		api: {
			lookups: [
				{
					fields: ["id"]
					description: """
						searches for template by ID

						It also returns deleted templates.
						"""
				}, {
					fields: ["handle"]
					nullConstraint: ["deleted_at"]
					constraintCheck: true
					description: """
						searches for template by handle

						It returns only valid templates (not deleted)
						"""
				},
			]
		}
	}
}
