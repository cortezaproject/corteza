package system

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

resource_translation: {
	features: {
		labels: false
		checkFn: false
	}

	model: {
		id: schema.IdField
		lang:       { goType: "types.Lang" }
		resource:   {}
		k:          {}
		message:    {}

		created_at: schema.SortableTimestampField
		updated_at: schema.SortableTimestampNilField
		deleted_at: schema.SortableTimestampNilField
		owned_by: { goType: "uint64" }
		created_by: { goType: "uint64" }
		updated_by: { goType: "uint64" }
		deleted_by: { goType: "uint64" }
	}

	filter: {
		model: {
			translation_id: {goType: "[]uint64", ident: "translationID" }
			lang: {}
			resource: {}
			resourceType: {}
			owner_id: {goType: "uint64", ident: "ownerID", storeIdent: "rel_owner"}
			deleted: {goType: "filter.State", storeIdent: "deleted_at"}
		}

		byValue: ["resource", "lang", "translation_id"]
		byNilState: ["deleted"]
	}

	store: {
		api: {
		lookups: [
			{
				fields: ["id"]
				description: """
					searches for resource translation by ID
					It also returns deleted resource translations.
					"""
			},
		]

		functions: [
				{
					expIdent: "TransformResource"
					args: [
						{ident: "lang", goType: "language.Tag" },
					]
					return: [ "map[string]map[string]*locale.ResourceTranslation" ]
				},
			]
		}
	}
}
