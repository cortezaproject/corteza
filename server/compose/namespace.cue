package compose

import (
	"github.com/crusttech/human/server/codegen/schema"
)

_namespaceDefs: {
			NamespaceMeta: { name: "NamespaceMeta", fields: [
				{ name: "Icon", type: "string", json: "icon,omitempty" },
				{ name: "IconID", type: "uint64", json: "iconID,string" },
				{ name: "Logo", type: "string", json: "logo,omitempty" },
				{ name: "LogoID", type: "uint64", json: "logoID,string" },
				{ name: "LogoEnabled", type: "bool", json: "logoEnabled,omitempty" },
				{ name: "HideSidebar", type: "bool", json: "hideSidebar" },
				{ name: "Subtitle", type: "string", json: "subtitle,omitempty" },
				{ name: "Description", type: "string", json: "description,omitempty" },
			]}
		}

namespace: {
	features: {
	}

	types: {
		gen: true
		defs: _namespaceDefs
	}

	model: {
		ident: "compose_namespace"
		attributes: {
			id:         schema.IdField
			tenant_id:  schema.TenantRefField
			project_id: schema.ProjectRefField
			slug: {
				sortable: true,
				goType: "string"
				dal: {}
				envoy: {
					identifier: true
				}
			}
			enabled: {
				goType: "bool"
				dal: { type: "Boolean" }
			}
			meta: {
				type: _namespaceDefs.NamespaceMeta
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
			}
			name: {
				sortable: true
				dal: {}
			}

			created_at: schema.SortableTimestampNowField
			updated_at: schema.SortableTimestampNilField
			deleted_at: schema.SortableTimestampNilField
			created_by_agent: schema.AttributeAgentRef
		}

		indexes: {
			"primary": { attribute: "id" }
			"unique_handle": {
				fields: [{ attribute: "slug", modifiers: ["LOWERCASE"] }]
				predicate: "slug != '' AND deleted_at IS NULL"
			}
		}
	}

	filter: {
		struct: {
			namespace_id: { goType: "[]uint64", ident: "namespaceID", storeIdent: "id" }
			tenant_id:    schema.TenantFilterField
			project_id:   schema.ProjectFilterField
			slug: { goType: "string" }
			name: { goType: "string" }
			deleted: { goType: "filter.State", storeIdent: "deleted_at" }
		}

		query: ["name", "slug"]
		byValue: ["namespace_id", "name", "slug"]
		byNilState: ["deleted"]
	}

	envoy: {
		scoped: true
		yaml: {
			supportMappedInput: true
			mappedField: "Slug"
			identKeyAlias: ["namespaces", "ns"]
		}
		store: {
			handleField: "Slug"
			extendedFilterBuilder: true
		}
	}

	service: {

		customBodyOps: ["lookup", "search", "create", "update", "delete"]

		customAccessOps: ["create"]

		omitCreateProp: true
		updateProp:     "changed"
	}

	rbac: {
		operations: {
			"read": {}
			"update": {}
			"delete": {}
			"export": description:         "Access to export the entire namespace"
			"manage": description:         "Access to namespace admin panel"
			"module.create": description:  "Create module on namespace"
			"modules.search": description: "List, search or filter module on namespace"
			"modules.export": description: "Export modules on namespace"
			"chart.create": description:   "Create chart on namespace"
			"charts.search": description:  "List, search or filter chart on namespace"
			"charts.export": description:  "Export charts on namespace"
			"page.create": description:    "Create page on namespace"
			"pages.search": description:   "List, search or filter pages on namespace"
		}
	}

	locale: {
		keys: {
			name: {}
			metaSubtitle: {
				path: ["meta", "subtitle"]
			}
			metaDescription: {
				path: ["meta", "description"]
			}
		}
	}

	store: {
		ident: "composeNamespace"

		api: {
			lookups: [
				{
					fields: ["slug"]
					constraintCheck: true
					nullConstraint: ["deleted_at"]
					description: """
						searches for namespace by slug (case-insensitive)
						"""
				}, {
					fields: ["id"]
					description: """
						searches for compose namespace by ID

						It returns compose namespace even if deleted
						"""
				},
			]
		}
	}
}
