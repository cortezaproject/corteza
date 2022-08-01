package system

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

report: {
	struct: {
		id:     schema.IdField
		handle: schema.HandleField
		meta: {goType: "*types.ReportMeta"}
		scenarios: {goType: "types.ReportScenarioSet"}
		sources: {goType: "types.ReportDataSourceSet"}
		blocks: {goType: "types.ReportBlockSet"}


		created_at: schema.SortableTimestampField
		updated_at: schema.SortableTimestampNilField
		deleted_at: schema.SortableTimestampNilField
		owned_by: { goType: "uint64" }
		created_by: { goType: "uint64" }
		updated_by: { goType: "uint64" }
		deleted_by: { goType: "uint64" }
	}

	filter: {
		struct: {
			report_id: {goType: "[]uint64", storeIdent: "id", ident: "reportID" }
			handle: {goType: "string"}
			deleted: {goType: "filter.State", storeIdent: "deleted_at"}
		}

		byValue: ["handle", "report_id"]
		byNilState: ["deleted"]
	}

	rbac: {
		operations: {
			read: description:   "Read report"
			update: description: "Update report"
			delete: description: "Delete report"
			run: description:    "Run report"
		}
	}

	store: {
		api: {
			lookups: [
				{
					fields: ["id"]
					description: """
						searches for report by ID

						It returns report even if deleted
						"""
				}, {
					fields: ["handle"]
					nullConstraint: ["deleted_at"]
					constraintCheck: true
					description: """
						searches for report by handle

						It returns report if deleted
						"""
				},
			]
		}
	}
	// locale:
	//   extended: true
	//   keys:
	//     - { path: name,    field: "Meta.Name" }
	//     - { path: description, field: "Meta.Description" }
	//     - { name: block title, path: "block.{{blockID}}.title", custom: true }
	//     - { name: block description, path: "block.{{blockID}}.description", custom: true }
}
