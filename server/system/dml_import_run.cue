package system

import (
	"github.com/crusttech/human/server/codegen/schema"
)

_dml_import_runDefs: {
			DmlImportMethod: {
				name: "DmlImportMethod"
				values: [
					{ident: "DmlImportMethodBackground", value: "background"},
				]
			}
		}

dml_import_run: {
	types: {
		gen: true
		defs: _dml_import_runDefs
	}

	features: {
		labels: false
	}

	model: {
		ident: "dml_import_runs"
		attributes: {
			id: schema.IdField & { json: { field: "runID", string: true } }
			connection_id: {
				goType: "uint64"
				ident: "connectionID"
				storeIdent: "rel_connection"
				dal: { type: "Ref", refModelResType: "corteza::system:dal-connection", default: 0 }
				json: "connectionID,string"
			}
			mapping_id: {
				goType: "uint64"
				ident: "mappingID"
				storeIdent: "rel_mapping"
				dal: { type: "Ref", refModelResType: "corteza::system:dml-mapping", default: 0 }
				json: "mappingID,string"
			}
			method: {
				type: _dml_import_runDefs.DmlImportMethod
				dal: { type: "Text", length: 32 }
				omitSetter: true
				omitGetter: true
				json: "method"
			}
			status: {
				sortable: true
				goType: "string"
				dal: { type: "Text", length: 32 }
				json: "status"
			}
			processed: {
				goType: "uint64"
				dal: { type: "Number", meta: { "rdbms:type": "bigint" } }
				json: "processed"
			}
			failed: {
				goType: "uint64"
				dal: { type: "Number", meta: { "rdbms:type": "bigint" } }
				json: "failed"
			}
			error: {
				goType: "string"
				dal: { type: "Text", length: 0 }
				omitSetter: true
				omitGetter: true
				json: { field: "error", omitEmpty: true }
			}
			// map[string]string; store: false — transient cursor state, not persisted in row
			cursor: {
				goType: "map[string]string"
				store: false
				omitSetter: true
				omitGetter: true
				json: { field: "cursor", omitEmpty: true }
			}
		}

		indexes: {
			"primary": { attribute: "id" }
			"idx_mapping": { attribute: "mapping_id" }
			"idx_status":  { attribute: "status" }
		}
	}

	filter: {
		struct: {
			import_run_id: { goType: "[]uint64", ident: "importRunID", storeIdent: "id" }
			mapping_id:    { goType: "uint64", ident: "mappingID", storeIdent: "rel_mapping" }
			status:        { goType: "[]string" }
		}
		byValue: ["import_run_id", "mapping_id", "status"]
	}

	envoy: { omit: true }

	store: {
		api: {
			lookups: [
				{
					fields: ["id"]
					description: "searches for DML import run by ID"
				},
				{
					fields: ["mapping_id"]
					description: "searches for DML import runs by mapping ID"
				},
			]
		}
	}
}
