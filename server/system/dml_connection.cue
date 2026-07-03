package system

import (
	"github.com/crusttech/human/server/codegen/schema"
)

_dml_connectionDefs: {
			DmlConnectionParams: { name: "DmlConnectionParams", fields: [
				{ name: "Type", type: "string", json: "type" },
				{ name: "Params", goType: "map[string]any", json: "params" },
				{ name: "ModelIdent", type: "string", json: "modelIdent" },
				{ name: "ModelIdentCheck", slice: true, type: "string", json: "modelIdentCheck" },
			]}
		}

dml_connection: {
	types: {
		gen: true
		defs: _dml_connectionDefs
	}

	features: {
		labels: false
	}

	model: {
		ident: "dml_connections"
		attributes: {
			id:     schema.IdField & { json: { field: "connectionID", string: true } }
			handle: schema.HandleField
			label: {
				goType: "string"
				dal: { type: "Text", length: 256 }
				json: "label"
			}
			params: {
				type: _dml_connectionDefs.DmlConnectionParams
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
				json: "params"
			}

			created_at: schema.SortableTimestampNowField
			updated_at: schema.SortableTimestampNilField
			deleted_at: schema.SortableTimestampNilField
		}

		indexes: {
			"primary": { attribute: "id" }
			"unique_handle": {
				attribute: "handle"
				predicate: "handle != '' AND deleted_at IS NULL"
			}
		}
	}

	filter: {
		struct: {
			connection_id: { goType: "[]uint64", ident: "connectionID", storeIdent: "id" }
			handle:        { goType: "string" }
			deleted:       { goType: "filter.State", storeIdent: "deleted_at" }
		}
		byValue: ["connection_id", "handle"]
		byNilState: ["deleted"]
	}

	envoy: { omit: true }

	store: {
		api: {
			lookups: [
				{
					fields: ["id"]
					description: "searches for DML connection by ID"
				},
				{
					fields: ["handle"]
					nullConstraint: ["deleted_at"]
					constraintCheck: true
					description: "searches for DML connection by handle"
				},
			]
		}
	}
}
