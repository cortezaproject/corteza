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
			sensitivity_level_id: {goType: "[]uint64", ident: "sensitivityLevelID", storeIdent: "id"}

			deleted: {goType: "filter.State", storeIdent: "deleted_at"}
		}

		byValue: ["sensitivity_level_id"]
		byNilState: ["deleted"]
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
