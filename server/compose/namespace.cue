package compose

import (
	"github.com/cortezaproject/corteza/server/codegen/schema"
)

namespace: {
	model: {
		ident: "compose_namespace"
		attributes: {
			id: schema.IdField
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
				goType: "types.NamespaceMeta"
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
			"pages.export": description:   "Export pages on namespace"
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
