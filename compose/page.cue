package compose

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

page: {
	parents: [
		{handle: "namespace"},
	]

	model: {
		ident: "compose_page"
		attributes: {
			id: schema.IdField
			title: {
				goType: "string",
				sortable: true
				dal: {}
			}
			handle: schema.HandleField
			self_id: {
				ident: "selfID",
				goType: "uint64",
				dal: { type: "Ref", refModelResType: "corteza::compose:page" }
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

			config: {
				goType: "types.PageConfig"
				dal: { type: "JSON", defaultEmptyObject: true }
			}
			blocks: {
				goType: "types.PageBlocks"
				dal: { type: "JSON", defaultEmptyObject: true }
			}
			children: {
				goType: "types.PageSet", store: false
			}
			visible: {
				goType: "bool"
				dal: { type: "Boolean", default: true }
			}
			weight: {
				goType: "int", sortable: true
				dal: { type: "Number", default: 0, meta: { "rdbms:type": "integer" } }
			}
			description: {
				goType: "string"
				dal: {}
			}

			created_at: schema.SortableTimestampNowField
			updated_at: schema.SortableTimestampNilField
			deleted_at: schema.SortableTimestampNilField
		}

		indexes: {
			"primary": { attribute: "id" }
			"namespace": { attribute: "namespace_id" },
			"module": { attribute: "module_id" },
			"self_id": { attribute: "self_id" },
			"unique_handle": {
				fields: [{ attribute: "handle", modifiers: ["LOWERCASE"] }, { attribute: "namespace_id" }]
				predicate: "handle != '' AND deleted_at IS NULL"
			}
		}
	}

	filter: {
		struct: {
			namespace_id: { goType: "uint64", ident: "namespaceID", storeIdent: "rel_namespace" }
			parent_id: { goType: "uint64", ident: "parentID" }
			module_id: { goType: "uint64", ident: "moduleID", storeIdent: "rel_module" }
			root: { goType: "bool" }
			handle: { goType: "string" }
			title: { goType: "string" }
			deleted: { goType: "filter.State", storeIdent: "deleted_at" }
		}

		query: ["handle", "title", "description"]
		byValue: ["handle", "namespace_id", "module_id"]
		byNilState: ["deleted"]
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
			title: {}
			description: {}
			recordToolbarButtonNewLabel: {
				path: ["recordToolbar", "new", "label"]
				customHandler: true
			}
			recordToolbarButtonEditLabel: {
				path: ["recordToolbar", "edit", "label"]
				customHandler: true
			}
			recordToolbarButtonSubmitLabel: {
				path: ["recordToolbar", "submit", "label"]
				customHandler: true
			}
			recordToolbarButtonDeleteLabel: {
				path: ["recordToolbar", "delete", "label"]
				customHandler: true
			}
			recordToolbarButtonCloneLabel: {
				path: ["recordToolbar", "clone", "label"]
				customHandler: true
			}
			recordToolbarButtonBackLabel: {
				path: ["recordToolbar", "back", "label"]
				customHandler: true
			}
			blockTitle: {
				path: ["pageBlock", {part: "blockID", var: true}, "title"]
				customHandler: true
			}
			blockDescription: {
				path: ["pageBlock", {part: "blockID", var: true}, "description"]
				customHandler: true
			}
			blockAutomationButtonLabel: {
				path: ["pageBlock", {part: "blockID", var: true}, "button", {part: "buttonID", var: true}, "label"]
				customHandler: true
			}
			blockContentBody: {
				path: ["pageBlock", {part: "blockID", var: true}, "content", "body"]
				customHandler: true
			}
		}
	}

	store: {
		ident: "composePage"

		api: {
			lookups: [
				{
					fields: ["namespace_id", "handle"]
					nullConstraint: ["deleted_at"]
					description: """
						searches for page by handle (case-insensitive)
						"""
				}, {
					fields: ["namespace_id", "module_id"]
					nullConstraint: ["deleted_at"]
					description: """
						searches for page by moduleID
						"""
				}, {
					fields: ["id"]
					description: """
						searches for compose page by ID

						It returns compose page even if deleted
						"""
				},
			]

			functions: [
				{
					expIdent: "ReorderComposePages"
					args: [
						{ ident: "namespace_id", goType: "uint64" },
						{ ident: "parent_id", goType: "uint64" },
						{ ident: "page_ids", goType: "[]uint64" }
					]
				}
			]
		}
	}
}
