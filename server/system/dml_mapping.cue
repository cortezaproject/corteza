package system

import (
	"github.com/crusttech/human/server/codegen/schema"
)

_dml_mappingDefs: {
			DmlColumnMapSet: { name: "DmlColumnMapSet", elem: _dml_mappingDefs.DmlColumnMap, elemPtr: true }
			DmlColumnMap: { name: "DmlColumnMap", fields: [
				{ name: "SourceIdent", type: "string", json: "sourceIdent" },
				{ name: "FieldName", type: "string", json: "fieldName" },
				{ name: "Label", type: "string", json: "label,omitempty" },
				{ name: "FieldKind", type: "string", json: "fieldKind" },
				{ name: "Skip", type: "bool", json: "skip" },
			]}
		}

dml_mapping: {
	types: {
		gen: true
		defs: _dml_mappingDefs
	}

	features: {
		labels: false
	}

	model: {
		ident: "dml_mappings"
		attributes: {
			id: schema.IdField & { json: { field: "mappingID", string: true } }
			connection_id: {
				goType: "uint64"
				ident: "connectionID"
				storeIdent: "rel_connection"
				dal: { type: "Ref", refModelResType: "corteza::system:dml-connection", default: 0 }
				json: "connectionID,string"
			}
			namespace_handle: {
				goType: "string"
				ident: "namespaceHandle"
				storeIdent: "namespace_handle"
				dal: { type: "Text", length: 64 }
				json: "namespaceHandle"
			}
			source_ident: {
				sortable: true
				goType: "string"
				ident: "sourceIdent"
				storeIdent: "source_ident"
				dal: { type: "Text", length: 256 }
				json: "sourceIdent"
			}
			module_handle: {
				sortable: true
				goType: "string"
				ident: "moduleHandle"
				storeIdent: "module_handle"
				dal: { type: "Text", length: 64 }
				json: "moduleHandle"
			}
			module_name: {
				goType: "string"
				ident: "moduleName"
				storeIdent: "module_name"
				dal: { type: "Text", length: 256 }
				json: "moduleName"
			}
			skip: {
				goType: "bool"
				dal: { type: "Boolean", default: false }
				json: "skip"
			}
			identifier: {
				goType: "string"
				dal: { type: "Text", length: 256 }
				json: "identifier"
			}
			columns: {
				type: _dml_mappingDefs.DmlColumnMapSet
				dal: { type: "JSON" }
				omitSetter: true
				omitGetter: true
				json: "columns"
			}

			created_at: schema.SortableTimestampNowField
			updated_at: schema.SortableTimestampNilField
			deleted_at: schema.SortableTimestampNilField
		}

		indexes: {
			"primary": { attribute: "id" }
			"idx_connection": { attribute: "connection_id" }
			"unique_source_per_connection": {
				fields: [{attribute: "connection_id"}, {attribute: "source_ident"}]
				predicate: "deleted_at IS NULL"
			}
		}
	}

	filter: {
		struct: {
			mapping_id:    { goType: "[]uint64", ident: "mappingID", storeIdent: "id" }
			connection_id: { goType: "uint64", ident: "connectionID", storeIdent: "rel_connection" }
			source_ident:  { goType: "string", ident: "sourceIdent", storeIdent: "source_ident" }
			deleted:       { goType: "filter.State", storeIdent: "deleted_at" }
		}
		byValue: ["mapping_id", "connection_id", "source_ident"]
		byNilState: ["deleted"]
	}

	envoy: { omit: true }

	store: {
		api: {
			lookups: [
				{
					fields: ["id"]
					description: "searches for DML mapping by ID"
				},
				{
					fields: ["connection_id", "source_ident"]
					nullConstraint: ["deleted_at"]
					description: "searches for DML mapping by connection and source table"
				},
			]
		}
	}
}
