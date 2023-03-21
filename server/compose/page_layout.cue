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

			namespace_id: {
				ident: "namespaceID",
				goType: "uint64",
				storeIdent: "rel_namespace"
				dal: { type: "Ref", refModelResType: "corteza::compose:namespace" }
			}
			weight: {
				goType: "int", sortable: true
				dal: { type: "Number", default: 0, meta: { "rdbms:type": "integer" } }
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
			page_id: { goType: "uint64", ident: "pageID", storeIdent: "page_id" }
			default: { goType: "bool", ident: "default" }
			handle: { goType: "string" }
			deleted: { goType: "filter.State", storeIdent: "deleted_at" }
		}

		query: ["handle"]
		byValue: ["handle", "namespace_id", "page_id"]
		byNilState: ["deleted"]
	}

	rbac: {
		operations: {
			// @todo not sure how RBAC should work here
			// we'll probably use the page's RBAC whith some bool flags for user-defined ones
		}
	}

	rbac: {
		operations: {
			"read": {}
			"update": {}
			"delete": {}
		}
	}

	locale: {
		extended: true

		keys: {
			name: {
				path: ["meta", "name"]
			}
			description: {
				path: ["meta", "description"]
			}
			actionLabel: {
				path: ["config", "actions", {part: "actionID", var: true}, "label"]
				customHandler: true
			}
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
					fields: ["id"]
					description: """
						searches for compose page layour by ID

						It returns compose page layour even if deleted
						"""
				},
			]

			functions: [
				{
					expIdent: "ReorderComposePageLayouts"
					args: [
						{ ident: "namespace_id", goType: "uint64" },
						{ ident: "page_id", goType: "uint64" },
						{ ident: "page_layout_ids", goType: "[]uint64" }
					]
				}
			]
		}
	}
}
