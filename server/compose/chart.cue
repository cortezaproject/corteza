package compose

import (
	"github.com/cortezaproject/corteza/server/codegen/schema"
)

chart: schema.#resource & {
	parents: [
		{handle: "namespace"},
	]

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
}
