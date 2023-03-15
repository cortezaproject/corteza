package compose

import (
	"github.com/cortezaproject/corteza/server/codegen/schema"
)

pageLayout: {
	parents: [
		{handle: "namespace"},
		{handle: "page"},
	]

	model: {
		ident: "compose_page_layout"
		attributes: {
			id: schema.IdField
			handle: schema.HandleField
			page_id: {
				ident: "pageID",
				goType: "uint64",
				dal: { type: "Ref", refModelResType: "corteza::compose:page" }
				sortable: true
			}
			parent_id: {
				ident: "parentID",
				goType: "uint64",
				dal: { type: "Ref", refModelResType: "corteza::compose:page-layout" }
				sortable: true
			}

			module_id: {
				ident: "moduleID",
				goType: "uint64",
				storeIdent: "rel_module"
				dal: { type: "Ref", refModelResType: "corteza::compose:module" }
			}
			namespace_id: {
				ident: "namespaceID",
				goType: "uint64",
				storeIdent: "rel_namespace"
				dal: { type: "Ref", refModelResType: "corteza::compose:namespace" }
			}

			meta: {
				goType: "*types.PageLayoutMeta"
				dal: { type: "JSON", defaultEmptyObject: true }
			}

			primary: {
				goType: "bool"
				dal: { type: "Boolean", default: false }
			}

			config: {
				goType: "types.PageLayoutConfig"
				dal: { type: "JSON", defaultEmptyObject: true }
			}
			blocks: {
				goType: "types.PageBlocks"
				dal: { type: "JSON", defaultEmptyObject: true }
			}

			owned_by:   schema.AttributeUserRef
			created_at: schema.SortableTimestampNowField
			updated_at: schema.SortableTimestampNilField
			deleted_at: schema.SortableTimestampNilField
		}

		indexes: {
			"primary": { attribute: "id" }
			"namespace": { attribute: "namespace_id" },
			"module": { attribute: "module_id" },
			"page_id": { attribute: "page_id" },
			"parent_id": { attribute: "parent_id" },
			"unique_handle": {
				fields: [{ attribute: "handle", modifiers: ["LOWERCASE"] }, { attribute: "namespace_id" }]
				predicate: "handle != '' AND deleted_at IS NULL"
			}
		}
	}

	filter: {
		struct: {
			namespace_id: { goType: "uint64", ident: "namespaceID", storeIdent: "rel_namespace" }
			page_id: { goType: "uint64", ident: "pageID", storeIdent: "rel_page" }
			module_id: { goType: "uint64", ident: "moduleID", storeIdent: "rel_module" }
			default: { goType: "bool", ident: "default" }
			handle: { goType: "string" }
			deleted: { goType: "filter.State", storeIdent: "deleted_at" }
		}

		query: ["handle"]
		byValue: ["handle", "namespace_id", "module_id", "page_id", "module_id"]
		byNilState: ["deleted"]
	}

	rbac: {
		operations: {
			// @todo not sure how RBAC should work here
			// we'll probably use the page's RBAC whith some bool flags for user-defined ones
		}
	}

	locale: {
		skipSvc: true
		// extended: true

		keys: {
			// @todo
		}
	}

	store: {
		ident: "composePageLayout"

		api: {
			lookups: [
				{
					fields: ["namespace_id", "handle"]
					nullConstraint: ["deleted_at"]
					description: """
						searches for page layour by handle (case-insensitive)
						"""
				}, {
					fields: ["namespace_id", "module_id"]
					nullConstraint: ["deleted_at"]
					description: """
						searches for page layour by moduleID
						"""
				}, {
					fields: ["id"]
					description: """
						searches for compose page layour by ID

						It returns compose page layour even if deleted
						"""
				},
			]
		}
	}
}
