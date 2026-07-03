package compose

import (
	"github.com/crusttech/human/server/codegen/schema"
)

_pageLayoutDefs: {
			PageLayoutMeta: { name: "PageLayoutMeta", fields: [
				{ name: "Title", type: "string", json: "title" },
				{ name: "Description", type: "string", json: "description" },
				{ name: "Style", goType: "map[string]any", json: "style,omitempty" },
			]}

			PageLayoutConfig: { name: "PageLayoutConfig", fields: [
				{ name: "Visibility", type: _pageLayoutDefs.PageLayoutVisibility, json: "visibility" },
				{ name: "Buttons", type: _pageLayoutDefs.PageLayoutButtonConfig, json: "buttons" },
				{ name: "Actions", slice: true, type: _pageLayoutDefs.PageLayoutAction, json: "actions,omitempty" },
				{ name: "Validation", type: _pageLayoutDefs.PageLayoutValidation, json: "validation" },
				{ name: "UseTitle", type: "bool", json: "useTitle" },
			]}

			PageLayoutVisibility: { name: "PageLayoutVisibility", fields: [
				{ name: "Expression", type: "string", json: "expression" },
				{ name: "Roles", slice: true, type: "string", json: "roles,omitempty" },
			]}

			PageLayoutButtonConfig: { name: "PageLayoutButtonConfig", fields: [
				{ name: "New", type: _pageLayoutDefs.PageLayoutButton, json: "new" },
				{ name: "Edit", type: _pageLayoutDefs.PageLayoutButton, json: "edit" },
				{ name: "Submit", type: _pageLayoutDefs.PageLayoutButton, json: "submit" },
				{ name: "Delete", type: _pageLayoutDefs.PageLayoutButton, json: "delete" },
				{ name: "Clone", type: _pageLayoutDefs.PageLayoutButton, json: "clone" },
				{ name: "Back", type: _pageLayoutDefs.PageLayoutButton, json: "back" },
			]}

			PageLayoutButton: { name: "PageLayoutButton", fields: [
				{ name: "Enabled", type: "bool", json: "enabled" },
				{ name: "Label", type: "string", json: "label" },
			]}

			PageLayoutAction: { name: "PageLayoutAction", fields: [
				{ name: "ActionID", type: "uint64", json: "actionID,string" },
				{ name: "Placement", type: "string", json: "placement" },
				{ name: "Meta", type: _pageLayoutDefs.PageLayoutActionMeta, json: "meta" },
				{ name: "Enabled", type: "bool", json: "enabled" },
				{ name: "Kind", type: "string", json: "kind" },
				{ name: "Params", type: "any", json: "params" },
			]}

			PageLayoutActionMeta: { name: "PageLayoutActionMeta", fields: [
				{ name: "Label", type: "string", json: "label" },
				{ name: "Style", goType: "map[string]any", json: "style,omitempty" },
			]}

			PageLayoutValidation: { name: "PageLayoutValidation", fields: [
				{ name: "RequiredFields", slice: true, type: _pageLayoutDefs.PageLayoutRequiredField, json: "requiredFields,omitempty" },
			]}

			PageLayoutRequiredField: { name: "PageLayoutRequiredField", fields: [
				{ name: "Field", type: "string", json: "field" },
				{ name: "Condition", type: "string", json: "condition" },
			]}
		}

pageLayout: {
	features: {
	}

	types: {
		gen: true

		defs: _pageLayoutDefs
	}

	parents: [
		{handle: "namespace"},
		{handle: "page"},
	]

	model: {
		ident: "compose_page_layout"

		defaultGetter: true
		defaultSetter: true

		attributes: {
			id:         schema.IdField
			tenant_id:  schema.TenantRefField
			project_id: schema.ProjectRefField
			handle: schema.HandleField

			// struct-only field: `Primary bool json:"primary"` is not backed by a
			// stored/dal attribute. store:false + no dal block keeps it out of the
			// store/dal/getters codegen so it appears ONLY in the generated struct.
			// convention yields json:"primary", matching the hand-written tag.
			primary: {
				goType: "bool"
				store: false
				omitGetter: true
				omitSetter: true
			}

			page_id: {
				ident: "pageID",
				goType: "uint64",
				json: { field: "pageID", string: true }
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
				json: { field: "parentID", string: true }
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
				json: { field: "namespaceID", string: true }
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
				type: _pageLayoutDefs.PageLayoutMeta
				json: { field: "meta", omitEmpty: true }
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
			}

			config: {
				type: _pageLayoutDefs.PageLayoutConfig
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
			}
			blocks: {
				goType: "types.PageLayoutBlocks"
				json: { field: "blocks", omitEmpty: true }
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
			}

			owned_by: schema.AttributeUserRef & {json: {field: "ownedBy", string: true}}
			created_at: schema.SortableTimestampNowField
			updated_at: schema.SortableTimestampNilField
			deleted_at: schema.SortableTimestampNilField
			created_by_agent: schema.AttributeAgentRef
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
			tenant_id:      schema.TenantFilterField
			project_id:     schema.ProjectFilterField
			namespace_id:   { goType: "uint64", ident: "namespaceID", storeIdent: "rel_namespace" }
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

	service: {

		scoped: true

		undelete: true

		customBodyOps: ["lookup", "search", "create", "update", "delete", "undelete"]

		customFunctions: [
			{
				name: "FindByHandle"
				cap:  "read"
				action: "Lookup"
				args: [
					{name: "namespaceID", goType: "uint64"},
					{name: "h", goType: "string"},
				]
				results: [
					{name: "c", goType: "*types.PageLayout"},
					{name: "err", goType: "error"},
				]
			},
			{
				name: "FindByPageLayoutID"
				cap:  "read"
				action: "Lookup"
				args: [
					{name: "namespaceID", goType: "uint64"},
					{name: "pageLayoutID", goType: "uint64"},
				]
				results: [
					{name: "p", goType: "*types.PageLayout"},
					{name: "err", goType: "error"},
				]
			},
			{
				name: "Reorder"
				cap:  "write"
				action: "Reorder"
				args: [
					{name: "namespaceID", goType: "uint64"},
					{name: "pageID", goType: "uint64"},
					{name: "pageLayoutIDs", goType: "[]uint64"},
				]
				results: [
					{name: "err", goType: "error"},
				]
			},
		]

		customAccessOps: ["search", "create"]

		opParents: {
			lookup: ["namespace"]
		}

		omitCreateProp: true
		updateProp:     "changed"
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
