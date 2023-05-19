package system

import (
	"github.com/cortezaproject/corteza/server/codegen/schema"
)

dal_schema_alteration: {
	features: {
		labels: false
		checkFn: false
	}

	model: {
		// lengths for the lang, resource fields are now a bit shorter
		// Reason for that is supported index length in MySQL
		attributes: {
			id: schema.IdField
			batchID: {
				goType: "uint64"
				dal: { type: "ID" }
			}
			dependsOn: {
				goType: "uint64"
				dal: { type: "Ref", refModelResType: "corteza::system:dal-schema-alteration" }
			}

			kind: {
				dal: { type: "Text", length: 256 }
			}
			params: {
				goType: "*types.DalSchemaAlterationParams"
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
			}

			created_at: schema.SortableTimestampNowField
			updated_at: schema.SortableTimestampNilField
			deleted_at: schema.SortableTimestampNilField
			completed_at: schema.SortableTimestampNilField
			created_by: schema.AttributeUserRef
			updated_by: schema.AttributeUserRef
			deleted_by: schema.AttributeUserRef
			completed_by: schema.AttributeUserRef
		}

		indexes: {
			"primary": { attribute: "id" }
			"unique_alteration": {
				 fields: [
				   { attribute: "id",      modifiers: [ "LOWERCASE" ] },
				   { attribute: "batchID", modifiers: [ "LOWERCASE" ] },
				 ]
		 	}
		}
	}

	envoy: {
		// Not needed for this resource (yet)
		omit: true
	}

	filter: {
		struct: {
			alteration_id: {goType: "[]uint64", ident: "alterationID" }
			batch_id: {goType: "[]uint64", ident: "batchID" }
			kind: {}
			deleted: {goType: "filter.State", storeIdent: "deleted_at"}
			completed: {goType: "filter.State", storeIdent: "completed_at"}
		}

		byValue: ["kind", "alteration_id"]
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
		}
	}
}
