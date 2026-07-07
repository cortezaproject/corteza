package system

import (
	"github.com/crusttech/human/server/codegen/schema"
)

queue_message: {
	features: {
		labels: false
		checkFn: false
		noIdField: true
	}

	types: {
		gen: true
	}

	model: {
		omitGetterSetter: true

		attributes: {
		  id:        schema.IdField & { json: "messageID" }
		  queue:     {
		  	sortable: true
		  	dal: {}
		  }
		  payload:   {
		  	goType: "[]byte"
		  	dal: { type: "Blob" }
		  }
		  created:   schema.SortableTimestampNilField & { json: "created" }
		  processed: schema.SortableTimestampNilField & { json: "processed" }
		}

		indexes: {
			"primary": { attribute: "id" }
		}
	}

	envoy: {
		omit: true
	}

	filter: {
		struct: {
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
