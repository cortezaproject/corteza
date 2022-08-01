package system

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

queue_message: {
	features: {
		labels: false
		checkFn: false
	}

	model: {
		id:        schema.IdField
		queue:     { sortable: true }
		payload:   { goType: "[]byte" }
		processed: schema.SortableTimestampNilField
		created:   schema.SortableTimestampNilField
	}

	filter: {
		model: {
			queue: {}
			processed: {goType: "filter.State", storeIdent: "processed"}
		}

		byValue: ["queue"]
		byNilState: ["processed"]
	}

	store: {
		api: {
			lookups: []
		}
	}
}
