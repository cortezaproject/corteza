package system

import (
	"github.com/crusttech/human/server/codegen/schema"
)

_templateDefs: {
			TemplateMeta: { name: "TemplateMeta", fields: [
				{ name: "Short", type: "string", json: "short" },
				{ name: "Description", type: "string", json: "description,omitempty" },
			]}
			DocumentType: { name: "DocumentType", values: [
				{ ident: "DocumentTypePlain", value: "text/plain" },
				{ ident: "DocumentTypeHTML", value: "text/html" },
				{ ident: "DocumentTypePDF", value: "application/pdf" },
			]}
		}

template: {
	features: {
	}

	types: {
		gen: true
		defs: _templateDefs
	}

	model: {
		// length for the lang is now a bit shorter
		// Reason for that is supported index length in MySQL
		attributes: {
			id:     schema.IdField
			tenant_id:  schema.TenantRefField
			project_id: schema.ProjectRefField
			owner_id:   {
				storeIdent: "rel_owner",
				ident: "ownerID"
				schema.AttributeUserRef,
				// hand-written tag is "ownerID,string"; Ref convention would add ",omitempty"
				json: "ownerID,string"
			}
			handle: schema.HandleField
			language: {
				sortable: true,
				goType: "string"
				dal: { length: 32 }
			}
			type: {
				sortable: true,
				type: _templateDefs.DocumentType
				dal: {}
				omitSetter: true
				omitGetter: true
			}
			partial: {
				goType: "bool"
				dal: { type: "Boolean" }
			}
			meta: {
				type: _templateDefs.TemplateMeta
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
			tenant_id: schema.TenantFilterField
			project_id: schema.ProjectFilterField
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

	service: {

		customFunctions: [
			{
				name:   "Render"
				cap:    "read"
				action: "Render"
				args: [
					{name: "templateID", goType: "uint64"},
					{name: "dstType", goType: "string"},
					{name: "variables", goType: "map[string]interface{}"},
					{name: "options", goType: "map[string]string"},
				]
				results: [
					{name: "document", goType: "io.ReadSeeker"},
					{name: "err", goType: "error"},
				]
			},
		]
		customFunctionImports: ["\"io\""]

		undelete: true

		updateFields: ["Handle", "Language", "Type", "Partial", "Meta", "Template", "OwnerID"]

		hooks: {validate: true}
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
