package compose

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

record_revision: {
	model: {
		ident: "compose_record_revisions"

		attributes: {
			id: schema.IdField
			timestamp: schema.SortableTimestampField & { storeIdent: "ts" }
			rel_resource: {
			 	ident: "resourceID",
				goType: "uint64",
				dal: { type: "ID" }
			}
			revision: {
				goType: "uint"
				dal: { type: "Number", meta: { "rdbms:type": "integer" } }
			}
			operation: {
				dal: {}
			}
			rel_user:   schema.AttributeUserRef
			delta: {
				goType: "types.RecordValueSet",
				dal: { type: "JSON", defaultEmptyObject: true }
			}
			comment: {
				dal: {}
			}
		}

		indexes: {
			"primary": { attribute: "id" }
		}
	}
}
