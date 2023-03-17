package system

import (
	"github.com/cortezaproject/corteza/server/codegen/schema"
)

dal_sensitivity_level: {
	model: {
		attributes: {
			id:     schema.IdField
			handle: schema.HandleField
			level: {
				sortable: true,
				goType: "int"
				dal: { type: "Number", meta: { "rdbms:type": "integer" } }
			}

			meta: {
				goType: "types.DalSensitivityLevelMeta"
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
			}
			created_at: schema.SortableTimestampNowField
			updated_at: schema.SortableTimestampNilField
			deleted_at: schema.SortableTimestampNilField
			created_by: schema.AttributeUserRef
			updated_by: schema.AttributeUserRef
			deleted_by: schema.AttributeUserRef
		}

		indexes: {
			"primary": { attribute: "id" }
		}
	}

	filter: {
		struct: {
			dal_sensitivity_level_id: {goType: "[]uint64", ident: "dalSensitivityLevelID", storeIdent: "id"}
			handle: { goType: "string" }

			deleted: {goType: "filter.State", storeIdent: "deleted_at"}
		}

		byValue: ["dal_sensitivity_level_id", "handle"]
		byNilState: ["deleted"]
	}

	envoy: {
		yaml: {
			supportMappedInput: true
			mappedField: "Handle"
			identKeyAlias: ["sensitivity_level"]
		}
		store: {}
	}

	features: {
		labels: false
	}

	store: {
		api: {
			lookups: [
				{
					fields: ["id"]
					description: """
						searches for user by ID

						It returns user even if deleted or suspended
						"""
				}
			]
		}
	}
}
