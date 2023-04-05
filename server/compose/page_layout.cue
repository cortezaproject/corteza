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

		defaultGetter: true
		defaultSetter: true

		attributes: {
			id: schema.IdField
			handle: schema.HandleField
			page_id: {
				ident: "pageID",
				goType: "uint64",
				dal: { type: "Ref", refModelResType: "corteza::compose:page" }
				sortable: true
				envoy: {
					yaml: {
						identKeyAlias: ["page"]
					}
				}
			}
			parent_id: {
				ident: "parentID",
				goType: "uint64",
				dal: { type: "Ref", refModelResType: "corteza::compose:page-layout" }
				sortable: true
				envoy: {
					yaml: {
						identKeyAlias: ["parent"]
					}
				}
			}

			namespace_id: {
				ident: "namespaceID",
				goType: "uint64",
				storeIdent: "rel_namespace"
				dal: { type: "Ref", refModelResType: "corteza::compose:namespace" }
				envoy: {
					yaml: {
						identKeyAlias: ["namespace"]
					}
				}
			}
			weight: {
				goType: "int", sortable: true
				dal: { type: "Number", default: 0, meta: { "rdbms:type": "integer" } }
			}

			meta: {
				goType: "types.PageLayoutMeta"
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
			}

			config: {
				goType: "types.PageLayoutConfig"
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
			}
			blocks: {
				goType: "types.PageLayoutBlocks"
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
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
				fields: [{ attribute: "handle", modifiers: ["LOWERCASE"] }, { attribute: "page_id" }, { attribute: "namespace_id" }]
				predicate: "handle != '' AND deleted_at IS NULL"
			}
		}
	}

	filter: {
		struct: {
			page_layout_id: { goType: "[]uint64", ident: "pageLayoutID", storeIdent: "id" }
			namespace_id: { goType: "uint64", ident: "namespaceID", storeIdent: "rel_namespace" }
			page_id: { goType: "uint64", ident: "pageID", storeIdent: "page_id" }
			parent_id: { goType: "uint64", ident: "parentID", storeIdent: "parent_id" }
			default: { goType: "bool", ident: "default" }
			handle: { goType: "string" }
			deleted: { goType: "filter.State", storeIdent: "deleted_at" }
		}

		query: ["handle"]
		byValue: ["handle", "parent_id", "namespace_id", "page_id", "page_layout_id"]
		byNilState: ["deleted"]
	}

	envoy: {
		scoped: true
		yaml: {
			supportMappedInput: true
			mappedField: "Handle"
			identKeyAlias: ["page_layouts", "pagelayouts", "layouts"]
		}
		store: {
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
			title: {
				path: ["meta", "title"]
			}
			description: {
				path: ["meta", "description"]
			}

			recordToolbarButtonNewLabel: {
				path: ["config", "buttons", "new", "label"]
				customHandler: true
			}
			recordToolbarButtonEditLabel: {
				path: ["config", "buttons", "edit", "label"]
				customHandler: true
			}
			recordToolbarButtonSubmitLabel: {
				path: ["config", "buttons", "submit", "label"]
				customHandler: true
			}
			recordToolbarButtonDeleteLabel: {
				path: ["config", "buttons", "delete", "label"]
				customHandler: true
			}
			recordToolbarButtonCloneLabel: {
				path: ["config", "buttons", "clone", "label"]
				customHandler: true
			}
			recordToolbarButtonBackLabel: {
				path: ["config", "buttons", "back", "label"]
				customHandler: true
			}
			actionLabel: {
				path: ["config", "actions", {part: "actionID", var: true}, "meta", "label"]
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
					fields: ["namespace_id", "page_id", "handle"]
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
