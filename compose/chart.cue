package compose

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

chart: schema.#Resource & {
	parents: [
		{handle: "namespace"},
	]

	struct: {
		id: schema.IdField
		handle: schema.HandleField
		name: {sortable: true}
		config: { goType: "types.ChartConfig" }
		namespace_id: { ident: "namespaceID", goType: "uint64", storeIdent: "rel_namespace" }

		created_at: schema.SortableTimestampField
		updated_at: schema.SortableTimestampNilField
		deleted_at: schema.SortableTimestampNilField
	}

	filter: {
		struct: {
			chart_id: { goType: "[]uint64", ident: "chartID", storeIdent: "id" }
			namespace_id: { goType: "uint64", ident: "namespaceID", storeIdent: "rel_namespace" }
			handle: { goType: "string" }
			name: { goType: "string" }
			deleted: { goType: "filter.State", storeIdent: "deleted_at" }
		}

		query: ["handle", "name"]
		byValue: ["handle", "chart_id", "namespace_id"]
		byNilState: ["deleted"]
	}

	rbac: {
		operations: {
			"read": {}
			"update": {}
			"delete": {}
		}
	}

	store: {
		ident: "composeChart"

		settings: {
			rdbms: {
				table: "compose_chart"
			}
		}

		api: {
			lookups: [
				{
					fields: ["id"]
					description: """
						searches for compose chart by ID

						It returns compose chart even if deleted
						"""
				}, {
					fields: ["namespace_id", "handle"]
					nullConstraint: ["deleted_at"]
					description: """
						searches for compose chart by handle (case-insensitive)
						"""
				},
			]
		}
	}
}
