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
		ident: "compose_module_field"
		attributes: {
			id: schema.IdField
			module_id: {
			  ident: "moduleID",
				goType: "uint64",
				storeIdent: "rel_module"
				dal: { type: "Ref", refModelResType: "corteza::compose:module" }
			}
			place: {
				sortable: true,
				goType: "int"
				dal: { type: "Number", meta: { "rdbms:type": "integer" } }
			}
			kind: {
				sortable: true,
				goType: "string"
				dal: {}
			}
			options: {
				goType: "types.ModuleFieldOptions"
				dal: { type: "JSON", defaultEmptyObject: true }
			}
			name: {
				sortable: true
				dal: {}
			}
			label: {
				sortable: true
				dal: {}
			}
			config: {
				goType: "types.ModuleFieldConfig"
				dal: { type: "JSON", defaultEmptyObject: true }
			}
			required: {
				goType: "bool",
				storeIdent: "is_required"
				dal: { type: "Boolean" }
			}
			multi: {
				goType: "bool",
				storeIdent: "is_multi"
				dal: { type: "Boolean" }
			}
			default_value: {
				goType: "types.RecordValueSet"
				dal: { type: "JSON", defaultEmptyObject: true }
			}
			expressions: {
				goType: "types.ModuleFieldExpr"
				dal: { type: "JSON", defaultEmptyObject: true }
			}
			created_at: schema.SortableTimestampNowField
			updated_at: schema.SortableTimestampNilField
			deleted_at: schema.SortableTimestampNilField
		}

		indexes: {
			"primary": { attribute: "id" }
			"module": { attribute: "module_id" },
			"unique_name": {
				fields: [{ attribute: "name", modifiers: ["LOWERCASE"] }, { attribute: "module_id" }]
				predicate: "name != '' AND deleted_at IS NULL"
			}
		}
	}

	filter: {
		struct: {
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
