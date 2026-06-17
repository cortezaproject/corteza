package compose

import (
	"github.com/crusttech/human/server/codegen/schema"
)

namespace: {
	features: {
		projectScoped: true
	}

	types: { gen: true }

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
		// single-id (non-compound) resource: FindByID/Search/Create/Update/
		// DeleteByID. UndeleteByID and the import/clone surface stay hand-written
		// in the companion file (custom methods), so undelete generation is off.

		// namespace CRUD is bespoke (label load on lookup, i18n decode + label load
		// on search, tx + translations + uniqueCheck + setChanged on create, the
		// shared updater for update/delete) so every op delegates its body to an
		// on<Op> handler.
		customBodyOps: ["lookup", "search", "create", "update", "delete"]

		// create's access check (CanCreateNamespace) runs INSIDE the tx, after the
		// handle validation, so the generated scaffold must not emit a standard
		// check for it. search keeps the standard CanSearchNamespaces check (emitted
		// by the generated wrapper, like system application), so it is not listed.
		customAccessOps: ["create"]

		// namespace's action props expose `namespace`/`changed`, not the default
		// `new`/`update`. actionProp already sets `namespace` in Create, so omit the
		// second prop; map Update's prop onto `changed`.
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
