package compose

import (
	"github.com/crusttech/human/server/codegen/schema"
)

chart: {
	features: {
		projectScoped: true
	}

	types: {
		gen: true
	}

	parents: [
		{handle: "namespace"},
	]

	model: {
		defaultSetter: true

		ident: "compose_chart"
		attributes: {
			id: schema.IdField & {
				envoy: {
					yaml: {
						identKeyEncode: "chartID"
					}
				}
			}
			handle:     schema.HandleField
			tenant_id:  schema.TenantRefField
			project_id: schema.ProjectRefField
			namespace_id: {
			  ident: "namespaceID",
				goType: "uint64",
				storeIdent: "rel_namespace"
				// parent ref carries no omitempty (hand-written json:"namespaceID,string")
				json: {field: "namespaceID", string: true}
				dal: { type: "Ref", refModelResType: "corteza::compose:namespace" }
				envoy: {
					yaml: {
						identKeyAlias: ["namespace", "namespace_id", "ns"]
					}
				}
			}
			name: {
				sortable: true
				dal: {}
		  }
			config: {
				goType: "types.ChartConfig"
				dal: {}
				omitSetter: true
				omitGetter: true
				envoy: {
					yaml: {
						customDecoder: true
						customEncoder: true
					}
				}
		  }
			created_at: schema.SortableTimestampNowField
			updated_at: schema.SortableTimestampNilField
			deleted_at: schema.SortableTimestampNilField
		}

		indexes: {
			"primary": { attribute: "id" }
			"namespace": { attribute: "namespace_id" },
			"unique_handle": {
				fields: [{ attribute: "handle", modifiers: ["LOWERCASE"] }, { attribute: "namespace_id" }]
				predicate: "handle != '' AND deleted_at IS NULL"
			}
		}

		indexes: {
			"primary": { attribute: "id" }
		}
	}

	filter: {
		struct: {
				chart_id:     { goType: "[]uint64", ident: "chartID", storeIdent: "id" }
				tenant_id:    schema.TenantFilterField
				project_id:   schema.ProjectFilterField
				namespace_id: { goType: "uint64", ident: "namespaceID", storeIdent: "rel_namespace" }
				handle: { goType: "string" }
				name: { goType: "string" }
				deleted: { goType: "filter.State", storeIdent: "deleted_at" }
		}

		query: ["handle", "name"]
		byValue: ["handle", "chart_id", "namespace_id"]
		byNilState: ["deleted"]
	}

	envoy: {
		scoped: true
		yaml: {
			supportMappedInput: true
			mappedField: "Handle"
			identKeyAlias: ["charts", "chrt"]
		}
		store: {
			extendedRefDecoder: true
		}
	}

	refs: {
		extended: true
	}

	service: {
		scoped: true

		undelete: true

		customBodyOps: ["lookup", "search", "create", "update", "delete", "undelete"]

		customAccessOps: ["search", "create"]

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
