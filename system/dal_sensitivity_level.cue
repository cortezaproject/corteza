package system

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

dal_sensitivity_level: schema.#Resource & {
	struct: {
		id:     schema.IdField
		handle: schema.HandleField
		level: { sortable: true, goType: "int" }
		meta: {goType: "types.DalSensitivityLevelMeta"}

		created_at: schema.SortableTimestampField
		updated_at: schema.SortableTimestampNilField
		deleted_at: schema.SortableTimestampNilField
		created_by: { goType: "uint64" }
		updated_by: { goType: "uint64" }
		deleted_by: { goType: "uint64" }
	}

	filter: {
		struct: {
			sensitivity_level_id: {goType: "[]uint64", ident: "sensitivityLevelID", storeIdent: "id"}

			deleted: {goType: "filter.State", storeIdent: "deleted_at"}
		}

		byValue: ["sensitivity_level_id"]
		byNilState: ["deleted"]
	}

	rbac: false

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
