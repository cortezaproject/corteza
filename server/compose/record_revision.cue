package compose

import (
	"github.com/crusttech/human/server/codegen/schema"
)

_record_revisionDefs: {
			RecordValueSet: { name: "RecordValueSet", elem: _record_revisionDefs.RecordValue, elemPtr: true }
			RecordValue: { name: "RecordValue", fields: [
				{ name: "RecordID", type: "uint64", json: "-" },
				{ name: "Name", type: "string", json: "name" },
				{ name: "Value", type: "string", json: "value,omitempty" },
				{ name: "Ref", type: "uint64", json: "-" },
				{ name: "Place", type: "uint", json: "place,omitempty" },
				{ name: "DeletedAt", type: "time.Time", ptr: true, json: "deletedAt,omitempty" },
				{ name: "Updated", type: "bool", json: "-" },
				{ name: "OldValue", type: "string", json: "-" },
			]}
		}

record_revision: {
	features: {
		noTypeSet: true
	}

	model: {
		ident: "compose_record_revisions"
		omitGetterSetter: true

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
				type: _record_revisionDefs.RecordValueSet,
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

	types: {
		defs: _record_revisionDefs
	}

	envoy: {
		omit: true
	}
}
