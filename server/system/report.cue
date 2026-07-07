package system

import (
	"github.com/crusttech/human/server/codegen/schema"
)

_reportDefs: {
			ReportMeta: { name: "ReportMeta", fields: [
				{ name: "Name", type: "string", json: "name" },
				{ name: "Description", type: "string", json: "description" },
			]}

			ReportScenarioSet: { name: "ReportScenarioSet", elem: _reportDefs.ReportScenario, elemPtr: true }
			ReportScenario: { name: "ReportScenario", fields: [
				{ name: "ScenarioID", type: "uint64", json: "scenarioID,string,omitempty" },
				{ name: "Label", type: "string", json: "label" },
				{ name: "Filters", type: _reportDefs.ScenarioFilterMap, json: "filters,omitempty" },
			]}
			ScenarioFilterMap: { name: "ScenarioFilterMap", key: "string", value: _reportDefs.ReportFilterExpr }

			ReportDataSourceSet: { name: "ReportDataSourceSet", elem: _reportDefs.ReportDataSource, elemPtr: true }
			ReportDataSource: { name: "ReportDataSource", fields: [
				{ name: "Meta", goType: "interface{}", json: "meta,omitempty" },
				{ name: "Step", type: _reportDefs.ReportStep, ptr: true, json: "step" },
			]}

			ReportBlockSet: { name: "ReportBlockSet", elem: _reportDefs.ReportBlock, elemPtr: true }
			ReportBlock: { name: "ReportBlock", fields: [
				{ name: "BlockID", type: "uint64", json: "blockID,string" },
				{ name: "Title", type: "string", json: "title" },
				{ name: "Description", type: "string", json: "description" },
				{ name: "Key", type: "string", json: "key" },
				{ name: "Kind", type: "string", json: "kind" },
				{ name: "Options", goType: "map[string]interface{}", json: "options,omitempty" },
				{ name: "Elements", goType: "[]interface{}", json: "elements" },
				{ name: "Sources", type: _reportDefs.ReportStepSet, json: "sources" },
				{ name: "XYWH", goType: "[4]int", json: "xywh" },
				{ name: "Layout", type: "string", json: "layout" },
			]}

			ReportStepSet: { name: "ReportStepSet", elem: _reportDefs.ReportStep, elemPtr: true }
			ReportStep: { name: "ReportStep", fields: [
				{ name: "Kind", type: "string", json: "kind,omitempty" },
				{ name: "Load", type: _reportDefs.ReportStepLoad, ptr: true, json: "load,omitempty" },
				{ name: "Join", type: _reportDefs.ReportStepJoin, ptr: true, json: "join,omitempty" },
				{ name: "Link", type: _reportDefs.ReportStepLink, ptr: true, json: "link,omitempty" },
				{ name: "Aggregate", type: _reportDefs.ReportStepAggregate, ptr: true, json: "aggregate,omitempty" },
				{ name: "Group_legacy", type: _reportDefs.ReportLegacyStepGroup, ptr: true, json: "group,omitempty" },
			]}
			ReportStepLoad: { name: "ReportStepLoad", fields: [
				{ name: "Name", type: "string", json: "name" },
				{ name: "Source", type: "string", json: "source" },
				{ name: "Definition", goType: "map[string]interface{}", json: "definition" },
				{ name: "Filter", type: _reportDefs.ReportFilterExpr, ptr: true, json: "filter,omitempty" },
			]}
			ReportStepJoin: { name: "ReportStepJoin", fields: [
				{ name: "Name", type: "string", json: "name" },
				{ name: "LocalSource", type: "string", json: "localSource" },
				{ name: "LocalColumn", type: "string", json: "localColumn" },
				{ name: "ForeignSource", type: "string", json: "foreignSource" },
				{ name: "ForeignColumn", type: "string", json: "foreignColumn" },
				{ name: "Filter", type: _reportDefs.ReportFilterExpr, ptr: true, json: "filter,omitempty" },
			]}
			ReportStepLink: { name: "ReportStepLink", fields: [
				{ name: "Name", type: "string", json: "name" },
				{ name: "LocalSource", type: "string", json: "localSource" },
				{ name: "LocalColumn", type: "string", json: "localColumn" },
				{ name: "ForeignSource", type: "string", json: "foreignSource" },
				{ name: "ForeignColumn", type: "string", json: "foreignColumn" },
				{ name: "Filter", type: _reportDefs.ReportFilterExpr, ptr: true, json: "filter,omitempty" },
			]}
			ReportLegacyStepGroup: { name: "ReportLegacyStepGroup", fields: [
				{ name: "Name", type: "string", json: "name" },
				{ name: "Source", type: "string", json: "source" },
				{ name: "Keys", type: _reportDefs.ReportAggregateColumnSet, json: "keys" },
				{ name: "Columns", type: _reportDefs.ReportAggregateColumnSet, json: "columns" },
				{ name: "Filter", type: _reportDefs.ReportFilterExpr, ptr: true, json: "filter,omitempty" },
			]}
			ReportStepAggregate: { name: "ReportStepAggregate", fields: [
				{ name: "Name", type: "string", json: "name" },
				{ name: "Source", type: "string", json: "source" },
				{ name: "Keys", type: _reportDefs.ReportAggregateColumnSet, json: "keys" },
				{ name: "Columns", type: _reportDefs.ReportAggregateColumnSet, json: "columns" },
				{ name: "Filter", type: _reportDefs.ReportFilterExpr, ptr: true, json: "filter,omitempty" },
			]}

			ReportAggregateColumnSet: { name: "ReportAggregateColumnSet", elem: _reportDefs.ReportAggregateColumn, elemPtr: true }
			ReportAggregateColumn: { name: "ReportAggregateColumn", fields: [
				{ name: "Name", type: "string", json: "name" },
				{ name: "Label", type: "string", json: "label" },
				{ name: "Def", type: _reportDefs.ReportFilterExpr, ptr: true, json: "def" },
			]}

			ReportFilterExpr: { name: "ReportFilterExpr", fields: [
				{ name: "ASTNode", goType: "*ast.ASTNode", json: ",omitempty" },
				{ name: "Error", type: "string", json: "error,omitempty" },
			]}
		}

report: {
	features: {
		labelResourceType: "report"
	}
	types: {
		gen: true
		// ReportMeta keeps a hand-written pointer-returning ParseReportMeta (REST
		// request controllers depend on the *ReportMeta return), incompatible with
		// the value-returning generated version, so it is excluded from JSON helper
		// generation.
		jsonTypesPtr: ["ReportMeta"]
		// ReportFilterExpr uses *ast.ASTNode embedding — hand-written; skip struct gen.
		structTypesSkip: ["ReportFilterExpr"]
		defs: _reportDefs
	}
	model: {
		attributes: {
			id:     schema.IdField
			tenant_id:  schema.TenantRefField
			project_id: schema.ProjectRefField
			handle: schema.HandleField
			meta: {
				type: _reportDefs.ReportMeta
				ptr: true
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
				type: _reportDefs.ReportScenarioSet
				dal: { type: "JSON", defaultEmptyObject: true }
				json: { field: "scenarios", omitEmpty: true }
				omitSetter: true
				omitGetter: true
			}
			sources: {
				type: _reportDefs.ReportDataSourceSet
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
			}
			blocks: {
				type: _reportDefs.ReportBlockSet
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
		customFunctions: [
			{
				name:  "Describe"
				cap:   "read"
				ac:    "CanCreateReport"
				acErr: "ErrNotAllowedToCreate"
				args: [
					{name: "src", goType: "types.ReportDataSourceSet"},
					{name: "st", goType: "types.ReportStepSet"},
					{name: "sources", goType: "...string"},
				]
				results: [
					{name: "out", goType: "[]reporting.FrameDescription"},
					{name: "err", goType: "error"},
				]
			},
			{
				name:   "Run"
				cap:    "read"
				action: "Run"
				args: [
					{name: "reportID", goType: "uint64"},
					{name: "dd", goType: "reporting.FrameDefinitionSet"},
				]
				results: [
					{name: "out", goType: "[]*reporting.Frame"},
					{name: "err", goType: "error"},
				]
			},
		]
		customFunctionImports: ["\"github.com/crusttech/human/server/system/reporting\""]

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
