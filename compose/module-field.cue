package compose

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

moduleField: schema.#Resource & {
	parents: [
		{handle: "namespace"},
		{handle: "module"},
	]

	struct: {
		id: schema.IdField
		namespace_id: { ident: "namespaceID", goType: "uint64", store: false }
		module_id: { ident: "moduleID", goType: "uint64", storeIdent: "rel_module" }
		place: { goType: "int" }
		kind: { goType: "string" }
		name: {}
		label: {}
		options: { goType: "types.ModuleFieldOptions" }
		encoding_strategy: { goType: "types.EncodingStrategy" }
		private: { goType: "bool", storeIdent: "is_private" }
		required: { goType: "bool", storeIdent: "is_required" }
		visible: { goType: "bool", storeIdent: "is_visible" }
		multi: { goType: "bool", storeIdent: "is_multi" }
		default_value: { goType: "types.RecordValueSet" }
		expressions: { goType: "types.ModuleFieldExpr" }
		created_at: { goType: "time.Time" }
		updated_at: { goType: "*time.Time" }
		deleted_at: { goType: "*time.Time" }
	}

	filter: {
		struct: {
			module_id: { goType: "[]uint64" }
			deleted: { goType: "filter.State", storeIdent: "deleted_at" }
		}

		byNilState: ["deleted"]
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
				},
			]
		}
	}
}
