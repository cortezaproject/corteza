package compose

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

chart: {
	parents: [
		{handle: "namespace"},
	]

	model: {
		ident: "compose_chart"
		attributes: {
			id: schema.IdField
			handle: schema.HandleField
			name: {sortable: true}
			config: { goType: "types.ChartConfig" }
			  namespace_id: {
			  ident: "namespaceID",
				goType: "uint64",
				storeIdent: "rel_namespace"
				dal: { type: "Ref", refModelResType: "corteza::compose:namespace" }
			}
			created_at: schema.SortableTimestampNowField
			updated_at: schema.SortableTimestampNilField
			deleted_at: schema.SortableTimestampNilField
		}
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

	locale: {
		extended: true

		keys: {
			reportsYaxisLabel: {
				path: ["yAxis", "label"]
				customHandler: true
			}
			reportsMetricLabel: {
				path: ["metrics", {part: "metricID", var: true}, "label"]
				customHandler: true
			}
			reportsDimensionStepLabel: {
				path: ["dimensions", {part: "dimensionID", var: true}, "meta", "steps", {part: "stepID", var: true}, "label"]
				customHandler: true
			}
		}
	}

	store: {
		ident: "composeChart"

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
