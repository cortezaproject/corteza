package system

import (
	"github.com/crusttech/human/server/codegen/schema"
)

_dal_schema_alterationDefs: {
			DalSchemaAlterationParams: { name: "DalSchemaAlterationParams", fields: [
				{ name: "AttributeAdd", goType: "*dal.AttributeAdd", json: "attributeAdd,omitempty" },
				{ name: "AttributeDelete", goType: "*dal.AttributeDelete", json: "attributeDelete,omitempty" },
				{ name: "AttributeReType", goType: "*dal.AttributeReType", json: "attributeReType,omitempty" },
				{ name: "AttributeReEncode", goType: "*dal.AttributeReEncode", json: "attributeReEncode,omitempty" },
				{ name: "ModelAdd", goType: "*dal.ModelAdd", json: "modelAdd,omitempty" },
				{ name: "ModelDelete", goType: "*dal.ModelDelete", json: "modelDelete,omitempty" },
			]}
		}

dal_schema_alteration: {
	features: {
		labels: false
		checkFn: false
	}

	types: {
		gen: true
		imports: ["github.com/crusttech/human/server/pkg/dal"]
		defs: _dal_schema_alterationDefs
	}

	model: {
		// lengths for the lang, resource fields are now a bit shorter
		// Reason for that is supported index length in MySQL
		attributes: {
			id: schema.IdField & {
				json: "alterationID,string"
			}
			batchID: {
				goType: "uint64"
				storeIdent: "batch_id"
				dal: { type: "ID" }
				json: { field: "batchID", string: true }
			}
			dependsOn: {
				goType: "uint64"
				storeIdent: "depends_on"
				dal: { type: "Ref", refModelResType: "corteza::system:dal-schema-alteration" }
			}
			resource: {
				storeIdent: "resource"
				dal: { type: "Text", length: 256 }
			}
			resourceType: {
				storeIdent: "resource_type"
				dal: { type: "Text", length: 256 }
			}
			connectionID: {
				goType: "uint64"
				storeIdent: "connection_id"
				dal: { type: "Ref", refModelResType: "corteza::system:dal-connection" }
				json: { field: "connectionID", string: true }
			}

			kind: {
				dal: { type: "Text", length: 256 }
			}
			params: {
				type: _dal_schema_alterationDefs.DalSchemaAlterationParams
				ptr: true
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
			}

			error: {
				dal: { type: "Text" }
				json: { field: "error", omitEmpty: true }
			}

			created_at: schema.SortableTimestampNowField
			updated_at: schema.SortableTimestampNilField
			deleted_at: schema.SortableTimestampNilField
			completed_at: schema.SortableTimestampNilField
			dismissed_at: schema.SortableTimestampNilField
			created_by: schema.AttributeUserRef & {
				json: { field: "createdBy", string: true }
			}
			updated_by: schema.AttributeUserRef
			deleted_by: schema.AttributeUserRef
			completed_by: schema.AttributeUserRef
			dismissed_by: schema.AttributeUserRef
		}

		indexes: {
			"primary": { attribute: "id" }
			"unique_alteration": {
				 fields: [
				   { attribute: "id" },
				   { attribute: "batchID" },
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
			resource: {goType: "[]string", ident: "resource" }
			resourceType: {goType: "string", ident: "resourceType", storeIdent: "resource_type" }
			alteration_id: {goType: "[]uint64", ident: "alterationID", storeIdent: "id" }
			batch_id: {goType: "[]uint64", ident: "batchID" }
			kind: {}
			deleted: {goType: "filter.State", storeIdent: "deleted_at"}
			completed: {goType: "filter.State", storeIdent: "completed_at"}
			dismissed: {goType: "filter.State", storeIdent: "dismissed_at"}
		}

		byValue: ["kind", "resource", "resourceType", "alteration_id", "batch_id"]
		byNilState: ["deleted", "completed", "dismissed"]
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
