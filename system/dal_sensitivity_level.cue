package system

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

dal_sensitivity_level: {
	model: {
		attributes: {
			id:     schema.IdField
			handle: schema.HandleField
			level: { sortable: true, goType: "int" }
			meta: {goType: "types.DalSensitivityLevelMeta"}

			created_at: schema.SortableTimestampNowField
			updated_at: schema.SortableTimestampNilField
			deleted_at: schema.SortableTimestampNilField
			created_by: schema.AttributeUserRef
			updated_by: schema.AttributeUserRef
			deleted_by: schema.AttributeUserRef
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
		paging: false
		sorting: false
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
