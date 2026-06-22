package system

import (
	"github.com/crusttech/human/server/codegen/schema"
)

report: {
	features: {
		projectScoped: true
	}
	types: {
		gen: true
		// ReportMeta keeps a hand-written pointer-returning ParseReportMeta (REST
		// request controllers depend on the *ReportMeta return), incompatible with
		// the value-returning generated version, so it is excluded from JSON helper
		// generation.
		jsonTypesPtr: ["ReportMeta"]
	}
	model: {
		attributes: {
			id:     schema.IdField
			tenant_id:  schema.TenantRefField
			project_id: schema.ProjectRefField
			handle: schema.HandleField
			meta: {
				goType: "*types.ReportMeta"
				dal: { type: "JSON", defaultEmptyObject: true }
				json: { field: "meta", omitEmpty: true }
				omitSetter: true
				omitGetter: true
			}
			// Virtual sortable mapped onto meta->>'name'.
			name: {
				sortableJSON: { json: "meta.name", accessor: "Meta.Name" }
				store: false
				goType: "string"
				omitSetter: true
				omitGetter: true
				envoy: { yaml: { omitEncoder: true } }
			}
			scenarios: {
				goType: "types.ReportScenarioSet"
				dal: { type: "JSON", defaultEmptyObject: true }
				json: { field: "scenarios", omitEmpty: true }
				omitSetter: true
				omitGetter: true
			}
			sources: {
				goType: "types.ReportDataSourceSet"
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
			}
			blocks: {
				goType: "types.ReportBlockSet"
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
			}

			owned_by:   schema.AttributeUserRef & { json: "ownedBy" }
			created_at: schema.SortableTimestampNowField & { json: "createdAt" }
			updated_at: schema.SortableTimestampNilField
			deleted_at: schema.SortableTimestampNilField
			created_by: schema.AttributeUserRef & { json: "createdBy" }
			updated_by: schema.AttributeUserRef & { json: { field: "updatedBy", omitEmpty: true } }
			deleted_by: schema.AttributeUserRef & { json: { field: "deletedBy", omitEmpty: true } }
		}

		indexes: {
			"primary": { attribute: "id" }
		}
	}

	envoy: {
		yaml: {
			supportMappedInput: true
			mappedField: "Handle"
			identKeyAlias: ["reports"]
		}
		store: {}
	}

	filter: {
		struct: {
			report_id: {goType: "[]uint64", storeIdent: "id", ident: "reportID" }
			tenant_id: schema.TenantFilterField
			project_id: schema.ProjectFilterField
			handle: {goType: "string"}
			deleted: {goType: "filter.State", storeIdent: "deleted_at"}
		}

		query: ["handle"]
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

	service: {
		undelete: true

		updateFields: ["Handle", "Meta", "Scenarios", "Sources", "Blocks"]

		hooks: {
			beforeCreate: true
			beforeUpdate: true
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
