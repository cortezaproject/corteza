package compose

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

moduleField: {
	parents: [
		{handle: "namespace"},
		{handle: "module"},
	]

	model: {
		id: schema.IdField
		namespace_id: { ident: "namespaceID", goType: "uint64", store: false }
		module_id: { sortable: true, ident: "moduleID", goType: "uint64", storeIdent: "rel_module" }
		place: { sortable: true, goType: "int" }
		kind: { sortable: true, goType: "string" }
		name: {sortable: true}
		label: {sortable: true}
		options: { goType: "types.ModuleFieldOptions" }
		config: { goType: "types.ModuleFieldConfig" }
		required: { goType: "bool", storeIdent: "is_required" }
		multi: { goType: "bool", storeIdent: "is_multi" }
		default_value: { goType: "types.RecordValueSet" }
		expressions: { goType: "types.ModuleFieldExpr" }
		created_at: { goType: "time.Time" }
		updated_at: { goType: "*time.Time" }
		deleted_at: { goType: "*time.Time" }
	}

	filter: {
		model: {
			module_id: { goType: "[]uint64", ident: "moduleID", storeIdent: "rel_module" }
			deleted: { goType: "filter.State", storeIdent: "deleted_at" }
		}

		byNilState: ["deleted"]
		byValue: ["module_id"]
	}

	features: {
		labels: false
		paging: false
		sorting: false
		checkFn: false
	}

	rbac: {
		operations: {
			"record.value.read": description:   "Read field value on records"
			"record.value.update": description: "Update field value on records"
		}
	}

	locale: {
		skipSvc: true

		keys: {
			label: {}
			descriptionView: {
				path: ["meta", "description", "view"]
				customHandler: true
			}
			descriptionEdit: {
				path: ["meta", "description", "edit"]
				customHandler: true
			}
			hintView: {
				path: ["meta", "hint", "view"]
				customHandler: true
			}
			hintEdit: {
				path: ["meta", "hint", "edit"]
				customHandler: true
			}
			validatorError: {
				path: ["expression", "validator", {part: "validatorID", var: true}, "error"]
				customHandler: true
			}
			optionsOptionTexts: {
				path: ["meta", "options", {part: "value", var: true}, "text"]
				customHandler: true
			}
			optionsBoolLabels: {
				path: ["meta", "bool", {part: "value", var: true}, "label"]
				customHandler: true
			}
		}
	}

	store: {
		ident: "composeModuleField"

		settings: {
			rdbms: {
				table: "compose_module_field"
			}
		}

		api: {
			lookups: [
				{
					fields: ["module_id", "name"]
					constraintCheck: true
					nullConstraint: ["deleted_at"]
					description: """
						searches for compose module field by name (case-insensitive)
						"""
				}, {
					fields: ["id"]
					description: """
						searches for compose module field by ID
						"""
				}
			]
		}
	}
}
