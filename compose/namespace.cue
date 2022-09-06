package compose

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

namespace: {
	model: {
		ident: "compose_namespace"
		attributes: {
			id: schema.IdField
			slug: { sortable: true, goType: "string" }
			enabled: { goType: "bool" }
			meta: { goType: "types.NamespaceMeta" }
			name: { sortable: true }

			created_at: schema.SortableTimestampNowField
			updated_at: schema.SortableTimestampNilField
			deleted_at: schema.SortableTimestampNilField
		}
	}

	filter: {
		struct: {
			namespace_id: { goType: "[]uint64", ident: "namespaceID" }
			slug: { goType: "string" }
			name: { goType: "string" }
			deleted: { goType: "filter.State", storeIdent: "deleted_at" }
		}

		query: ["name", "slug"]
		byValue: ["namespace_id", "name", "slug"]
		byNilState: ["deleted"]
	}

	rbac: {
		operations: {
			"read": {}
			"update": {}
			"delete": {}
			"manage": description:         "Access to namespace admin panel"
			"module.create": description:  "Create module on namespace"
			"modules.search": description: "List, search or filter module on namespace"
			"chart.create": description:   "Create chart on namespace"
			"charts.search": description:  "List, search or filter chart on namespace"
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
