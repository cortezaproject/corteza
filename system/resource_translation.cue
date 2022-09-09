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
		attributes: {
			id: schema.IdField
			lang: {
		 		goType: "types.Lang"
				dal: { type: "Text", length: 128 }
		 	}
			resource: {
				dal: { type: "Text", length: 512 }
			}
			k: {
				dal: { type: "Text", length: 256 }
			}
			message: {
				dal: {}
			}

			created_at: schema.SortableTimestampNowField
			updated_at: schema.SortableTimestampNilField
			deleted_at: schema.SortableTimestampNilField
			owned_by:   schema.AttributeUserRef
			created_by: schema.AttributeUserRef
			updated_by: schema.AttributeUserRef
			deleted_by: schema.AttributeUserRef
		}

		indexes: {
			"primary": { attribute: "id" }
			"unique_translation": {
				 fields: [
				   { attribute: "lang",     modifiers: [ "LOWERCASE" ] },
				   { attribute: "resource", modifiers: [ "LOWERCASE" ] },
				 	 { attribute: "k",        modifiers: [ "LOWERCASE" ] },
				 ]
		 	}
		}
	}

	filter: {
		struct: {
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
