package compose

import (
	"github.com/crusttech/human/server/codegen/schema"
)

_pageDefs: {
			PageBlocks: {
				name: "PageBlocks"
				elem: _pageDefs.PageBlock
			}
			PageBlock: {
				name: "PageBlock"
				fields: [
					{name: "BlockID", type: "uint64", json: "blockID,string,omitempty"},
					{name: "Options", goType: "map[string]interface{}", json: "options,omitempty"},
					{name: "Style", type: _pageDefs.PageBlockStyle, json: "style,omitempty"},
					{name: "Kind", type: "string", json: "kind"},
					{name: "XYWH", goType: "[4]int", json: "xywh"},
					{name: "Meta", goType: "map[string]any", json: "meta,omitempty"},
					{name: "Title", type: "string", json: "title,omitempty"},
					{name: "Description", type: "string", json: "description,omitempty"},
				]
			}
			PageBlockStyle: {
				name: "PageBlockStyle"
				fields: [
					{name: "Variants", goType: "map[string]string", json: "variants,omitempty"},
					{name: "Wrap", goType: "map[string]string", json: "wrap,omitempty"},
					{name: "Border", goType: "map[string]interface{}", json: "border,omitempty"},
				]
			}
			PageChildrenDeleteStrategy: {
				name: "PageChildrenDeleteStrategy"
				doc: "how child pages are handled on delete"
				kind: "string"
				values: [
					{ident: "Abort", value: "abort"},
					{ident: "Rebase", value: "rebase", doc: "reattach children to parent"},
				]
			}
			PageMeta: {
				name: "PageMeta"
				fields: [
					{name: "AllowPersonalLayouts", type: "bool", json: "allowPersonalLayouts"},
					{name: "Notifications", goType: "map[string]any", json: "notifications,omitempty"},
				]
			}
			PageConfig: {
				name: "PageConfig"
				fields: [
					{name: "NavItem", type: _pageDefs.PageConfigNavItem, json: "navItem"},
				]
			}
			PageConfigNavItem: {
				name: "PageConfigNavItem"
				fields: [
					{name: "Expanded", type: "bool", json: "expanded"},
					{name: "Icon", type: _pageDefs.PageConfigIcon, ptr: true, json: "icon,omitempty"},
				]
			}
			PageConfigIcon: {
				name: "PageConfigIcon"
				fields: [
					{name: "Type", goType: "IconType", json: "type,omitempty"},
					{name: "Src", type: "string", json: "src"},
					{name: "Style", goType: "map[string]string", json: "style,omitempty"},
				]
			}
		}

page: {
	features: {
		labelResourceType: "compose:page"
	}

	types: {
		gen: true
		defs: _pageDefs
	}

	parents: [
		{handle: "namespace"},
	]

	model: {
		defaultSetter: true

		ident: "compose_page"
		attributes: {
			id:         schema.IdField
			tenant_id:  schema.TenantRefField
			project_id: schema.ProjectRefField
			title: {
				goType: "string",
				sortable: true
				dal: {}
			}
			handle: schema.HandleField
			self_id: {
				ident: "selfID",
				goType: "uint64",
				json: { field: "selfID", string: true }
				dal: { type: "Ref", refModelResType: "corteza::compose:page" }
				sortable: true
				envoy: {
					store: {
						filterRefField: "ParentID"
					}
					yaml: {
						identKeyAlias: ["parent"]
					}
				}
			}
			module_id: {
				ident: "moduleID",
				goType: "uint64",
				json: { field: "moduleID", string: true }
				storeIdent: "rel_module"
				dal: { type: "Ref", refModelResType: "corteza::compose:module" }
				envoy: {
					yaml: {
						identKeyAlias: ["module"]
					}
				}
			}
			namespace_id: {
				ident: "namespaceID",
				goType: "uint64",
				json: { field: "namespaceID", string: true }
				storeIdent: "rel_namespace"
				dal: { type: "Ref", refModelResType: "corteza::compose:namespace" }
				envoy: {
					yaml: {
						identKeyAlias: ["namespace"]
					}
				}
			}

			meta: {
				type: _pageDefs.PageMeta
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
			}
			config: {
				type: _pageDefs.PageConfig
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
			}
			blocks: {
				type: _pageDefs.PageBlocks
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
				envoy: {
					yaml: {
						customDecoder: true
						customEncoder: true
					}
				}
			}
			children: {
				goType: "types.PageSet", store: false
				json: { field: "children", omitEmpty: true }
				omitSetter: true
				omitGetter: true
			}
			visible: {
				goType: "bool"
				dal: { type: "Boolean", default: true }
			}
			weight: {
				goType: "int", sortable: true
				dal: { type: "Number", default: 0, meta: { "rdbms:type": "integer" } }
				envoy: {
					yaml: {
						identKeyAlias: ["order"]
					}
				}
			}
			description: {
				goType: "string"
				dal: {}
			}

			created_at: schema.SortableTimestampNowField
			updated_at: schema.SortableTimestampNilField
			deleted_at: schema.SortableTimestampNilField
			created_by_agent: schema.AttributeAgentRef
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
			page_id:      { goType: "[]uint64", ident: "pageID", storeIdent: "id" }
			tenant_id:    schema.TenantFilterField
			project_id:   schema.ProjectFilterField
			namespace_id: { goType: "uint64", ident: "namespaceID", storeIdent: "rel_namespace" }
			parent_id: { goType: "uint64", ident: "parentID" }
			module_id: { goType: "uint64", ident: "moduleID", storeIdent: "rel_module" }
			root: { goType: "bool" }
			handle: { goType: "string" }
			title: { goType: "string" }
			deleted: { goType: "filter.State", storeIdent: "deleted_at" }
		}

		query: ["handle", "title", "description"]
		byValue: ["page_id", "handle", "namespace_id", "module_id", "project_id"]
		byNilState: ["deleted"]
	}

	envoy: {
		scoped: true
		yaml: {
			supportMappedInput: true
			mappedField: "Handle"
			identKeyAlias: ["pages", "pg"]

			extendedResourceDecoders: [{
				ident: "pages"
				expIdent: "Pages"
				identKeys: ["children", "pages"]
				supportMappedInput: true
				mappedField: "Handle"
			}]
			extendedResourceRefIdent: "SelfID"
		}
		store: {
			extendedFilterBuilder: true
			extendedRefDecoder: true
		}
	}

	refs: {
		items: [
			{path: "ModuleID", kind: "KindComposeModule", reason: "ReasonPageModule"},
		]
		extended: true
	}

	service: {
		scoped: true

		undelete: true

		customBodyOps: ["lookup", "search", "create", "update", "delete", "undelete"]

		customFunctions: [
			{
				name: "Tree"
				cap:  "read"
				args: [
					{name: "namespaceID", goType: "uint64"},
				]
				results: [
					{name: "tree", goType: "types.PageSet"},
					{name: "err", goType: "error"},
				]
			},
			{
				name:   "Reorder"
				cap:    "write"
				action: "Reorder"
				args: [
					{name: "namespaceID", goType: "uint64"},
					{name: "parentID", goType: "uint64"},
					{name: "pageIDs", goType: "[]uint64"},
				]
				results: [
					{name: "err", goType: "error"},
				]
			},
			{
				name: "UpdateIcon"
				cap:  "write"
				args: [
					{name: "namespaceID", goType: "uint64"},
					{name: "pageID", goType: "uint64"},
					{name: "icon", goType: "*types.PageConfigIcon"},
				]
				results: [
					{name: "out", goType: "*types.PageConfigIcon"},
					{name: "err", goType: "error"},
				]
			},
		]

		customAccessOps: ["search", "create"]

		omitCreateProp: true
		updateProp:     "changed"

		deleteExtraArgs: [
			{name: "strategy", goType: "types.PageChildrenDeleteStrategy"},
		]
	}

	rbac: {
		operations: {
			"read": {}
			"update": {}
			"delete": {}
			"page-layout.create": description:    "Create page layout on namespace"
			"page-layouts.search": description:   "List, search or filter page layouts on namespace"
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
