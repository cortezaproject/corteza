package compose

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

recordValue: schema.#Resource & {
	struct: {
		record_id: { primaryKey: true, goType: "uint64", ident: "recordID" }
		name:      { primaryKey: true, }
		place:     { primaryKey: true, goType: "uint" }
		ref:       { goType: "uint64" }
		value:     { }
		deleted_at: schema.SortableTimestampNilField
	}

	features: {
		labels: false
		paging: false
		sorting: false
		checkFn: false
	}

	filter: {
		struct: {
			record_id: { goType: "[]uint64", ident: "recordID" }
			deleted: { goType: "filter.State", storeIdent: "deleted_at" }
		}

		byValue: ["record_id"]
		byNilState: ["deleted"]
	}


	store: {
		ident: "composeRecordValue"

		settings: {
			rdbms: {
				table: "compose_record_value"
			}
		}

		api: {


			functions: [
				{
					expIdent: "ComposeRecordValueRefLookup"
					args: [
						{ ident: "mod", goType: "*types.Module" },
						{ ident: "field", goType: "string" },
						{ ident: "ref", goType: "uint64" }
					]
					return: [ "uint64" ]
				}
			]
		}
	}
}
