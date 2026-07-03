package compose

import (
	"github.com/crusttech/human/server/codegen/schema"
)

_moduleFieldDefs: {
			ModuleFieldConfig: { name: "ModuleFieldConfig", fields: [
				{ name: "DAL", type: _moduleFieldDefs.ModuleFieldConfigDAL, json: "dal" },
				{ name: "Privacy", type: _moduleFieldDefs.ModuleFieldConfigDataPrivacy, json: "privacy" },
				{ name: "RecordRevisions", type: _moduleFieldDefs.ModuleFieldConfigRecordRevisions, json: "recordRevisions" },
			]}
			ModuleFieldConfigDAL: { name: "ModuleFieldConfigDAL", fields: [
				{ name: "EncodingStrategy", goType: "*EncodingStrategy", json: "encodingStrategy" },
			]}
			ModuleFieldConfigDataPrivacy: { name: "ModuleFieldConfigDataPrivacy", fields: [
				{ name: "SensitivityLevelID", type: "uint64", json: "sensitivityLevelID,string,omitempty" },
				{ name: "UsageDisclosure", type: "string", json: "usageDisclosure" },
			]}
			ModuleFieldConfigRecordRevisions: { name: "ModuleFieldConfigRecordRevisions", fields: [
				{ name: "Skip", type: "bool", json: "enabled" },
			]}
			ModuleFieldExpr: { name: "ModuleFieldExpr", fields: [
				{ name: "ValueExpr", type: "string", json: "value,omitempty" },
				{ name: "Sanitizers", slice: true, type: "string", json: "sanitizers,omitempty" },
				{ name: "Validators", slice: true, type: _moduleFieldDefs.ModuleFieldValidator, json: "validators,omitempty" },
				{ name: "DisableDefaultValidators", type: "bool", json: "disableDefaultValidators,omitempty" },
				{ name: "Formatters", slice: true, type: "string", json: "formatters,omitempty" },
				{ name: "DisableDefaultFormatters", type: "bool", json: "disableDefaultFormatters,omitempty" },
			]}
			ModuleFieldValidator: { name: "ModuleFieldValidator", fields: [
				{ name: "ValidatorID", type: "uint64", json: "validatorID,string,omitempty" },
				{ name: "Test", type: "string", json: "test,omitempty" },
				{ name: "Error", type: "string", json: "error,omitempty" },
			]}
			ModuleFieldOptions: { name: "ModuleFieldOptions", key: "string", valueGoType: "interface{}" }
			RecordValueSet: { name: "RecordValueSet", elem: _moduleFieldDefs.RecordValue, elemPtr: true }
			RecordValue: { name: "RecordValue", fields: [
				{ name: "RecordID", type: "uint64", json: "-" },
				{ name: "Name", type: "string", json: "name" },
				{ name: "Value", type: "string", json: "value,omitempty" },
				{ name: "Ref", type: "uint64", json: "-" },
				{ name: "Place", goType: "uint", json: "place,omitempty" },
				{ name: "DeletedAt", type: "time.Time", ptr: true, json: "deletedAt,omitempty" },
				{ name: "Updated", type: "bool", json: "-" },
				{ name: "OldValue", type: "string", json: "-" },
			]}
		}

moduleField: {
	parents: [
		{handle: "namespace"},
		{handle: "module"},
	]

	model: {
		defaultSetter: true

		ident: "compose_module_field"
		attributes: {
			id: schema.IdField & {
				envoy: {
					identifier: true
				}
			}
			tenant_id:  schema.TenantRefField
			project_id: schema.ProjectRefField
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
				type: _moduleFieldDefs.ModuleFieldOptions
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
				envoy: {
					yaml: {
						customDecoder: true
						customEncoder: true
					}
				}
			}
			name: {
				sortable: true
				dal: {}
			} & {
				envoy: {
					identifier: true
				}
			}
			label: {
				sortable: true
				dal: {}
			}
			config: {
				type: _moduleFieldDefs.ModuleFieldConfig
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
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
				type: _moduleFieldDefs.RecordValueSet
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
				envoy: {
					yaml: {
						customDecoder: true
					}
				}
			}
			expressions: {
				type: _moduleFieldDefs.ModuleFieldExpr
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
				envoy: {
					yaml: {
						customDecoder: true
					}
				}
			}
			created_at: schema.SortableTimestampNowField
			updated_at: schema.SortableTimestampNilField
			deleted_at: schema.SortableTimestampNilField
			created_by_agent: schema.AttributeAgentRef
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

	types: {
		defs: _moduleFieldDefs
	}

	refs: {
		extended: true
	}

	filter: {
		struct: {
			tenant_id:  schema.TenantFilterField
			project_id: schema.ProjectFilterField
			module_id:  { goType: "[]uint64", ident: "moduleID", storeIdent: "rel_module" }
			deleted:    { goType: "filter.State", storeIdent: "deleted_at" }
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

	envoy: {
		scoped: true
		yaml: {
			supportMappedInput: true
			mappedField: "Name"
			identKeyAlias: ["module_fields", "modulefields", "fields"]
		}
		store: {
			handleField: ""
			customFilterBuilder: true
			extendedRefDecoder: true
			sanitizeBeforeSave: true
		}
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
