package compose

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

page: schema.#Resource & {
	parents: [
		{handle: "namespace"},
	]

	struct: {
		id: schema.IdField
		self_id: { ident: "selfID", goType: "uint64" }
		module_id: { ident: "moduleID", goType: "uint64", storeIdent: "rel_module" }
		namespace_id: { ident: "namespaceID", goType: "uint64", storeIdent: "rel_namespace" }
		handle: schema.HandleField
		config: { goType: "types.PageConfig" }
		blocks: { goType: "types.PageBlocks" }
		children: { goType: "types.PageSet", store: false }
		visible: { goType: "bool" }
		weight: { goType: "int" }
		title: { goType: "string" }
		description: { goType: "string" }

		created_at: schema.SortableTimestampField
		updated_at: schema.SortableTimestampNilField
		deleted_at: schema.SortableTimestampNilField
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

		settings: {
			rdbms: {
				table: "compose_page"
			}
		}

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
