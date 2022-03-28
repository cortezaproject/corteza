package compose

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

record: schema.#Resource & {
	parents: [
		{handle: "namespace"},
		{handle: "module"},
	]

	struct: {
		id: schema.IdField
		module_id: { ident: "moduleID", goType: "uint64" }
		module: { goType: "*types.Module", store: false }
		values: { goType: "types.RecordValueSet", store: false }
		namespace_id: { ident: "namespaceID", goType: "uint64", storeIdent: "rel_namespace" }

		owned_by: { goType: "uint64" }
		created_at: schema.SortableTimestampField
		updated_at: schema.SortableTimestampNilField
		deleted_at: schema.SortableTimestampNilField
		created_by: { goType: "uint64" }
		updated_by: { goType: "uint64" }
		deleted_by: { goType: "uint64" }
	}

	filter: {
		struct: {
			module_id: { goType: "uint64" }
			namespace_id: { goType: "uint64" }
			query: { goType: "string" }
			deleted: { goType: "filter.State", storeIdent: "deleted_at" }
		}

		byNilState: ["deleted"]
	}

	rbac: {
		operations: {
			"read": {}
			"update": {}
			"delete": {}
		}
	}

	store: {
		ident: "composeRecord"

		settings: {
			rdbms: {
				table: "compose_record"
			}
		}

		api: {
			lookups: [
				{
					fields: ["id"]
					// export: false
					description: """
						searches for compose record by ID

						It returns compose record even if deleted
						"""
				},
			]

			functions: [
				{
					expIdent: "ComposeRecordReport"
					args: [
						{ ident: "mod", goType: "*types.Module" },
						{ ident: "metrics", goType: "string" },
						{ ident: "dimensions", goType: "string" },
						{ ident: "filters", goType: "string" }
					]
					return: [ "[]map[string]interface{}" ]
				}, {
					expIdent: "ComposeRecordDatasource"
					args: [
						{ ident: "mod", goType: "*types.Module" },
						{ ident: "ld", goType: "*report.LoadStepDefinition" }
					]
					return: [ "report.Datasource" ]
				}, {
					expIdent: "PartialComposeRecordValueUpdate"
					args: [
						{ ident: "mod", goType: "*types.Module" },
						{ ident: "values", goType: "*types.RecordValue", spread: true }
					]
				}
			]
		}
	}
}
