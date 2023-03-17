package system

import (
	"github.com/cortezaproject/corteza/server/codegen/schema"
)

template: {
	model: {
		// length for the lang is now a bit shorter
		// Reason for that is supported index length in MySQL
		attributes: {
			id:     schema.IdField
			owner_id:   {
				storeIdent: "rel_owner",
				ident: "ownerID"
				schema.AttributeUserRef,
			}
			handle: schema.HandleField
			language: {
				sortable: true,
				goType: "string"
				dal: { length: 32 }
			}
			type: {
				sortable: true,
				goType: "types.DocumentType"
				dal: {}
				omitSetter: true
				omitGetter: true
			}
			partial: {
				goType: "bool"
				dal: { type: "Boolean" }
			}
			meta: {
				goType: "types.TemplateMeta"
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
			}
			template: {
				sortable: true,
				goType: "string"
				dal: {}
			}

			created_at: schema.SortableTimestampNowField
			updated_at: schema.SortableTimestampNilField
			deleted_at: schema.SortableTimestampNilField
			last_used_at: schema.SortableTimestampNilField
		}

		indexes: {
			"primary": { attribute: "id" }
			"unique_language_handle": {
				fields: [
					{ attribute: "language" },
					{ attribute: "handle", modifier: [ "LOWERCASE" ] }
				]
				predicate: "handle != '' AND deleted_at IS NULL"
			}
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

	envoy: {
		yaml: {
			supportMappedInput: true
			mappedField: "Handle"
			identKeyAlias: ["templates"]
		}
		store: {}
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
