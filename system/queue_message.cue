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
		attributes: {
		  id:        schema.IdField
		  queue:     {
		  	sortable: true
		  	dal: {}
		  }
		  payload:   {
		  	goType: "[]byte"
		  	dal: { type: "Blob" }
		  }
		  created:   schema.SortableTimestampNilField
		  processed: schema.SortableTimestampNilField
		}

		indexes: {
			"primary": { attribute: "id" }
		}
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
