package system

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

credential: {
	model: {
		attributes: {
			id:     schema.IdField
			owner_id: { schema.AttributeUserRef, storeIdent: "rel_owner", ident: "ownerID" }
			label: {
				dal: {}
			}
			kind: {
				dal: { type: "Text", length: 128 }
			}
			credentials: {
				dal: {}
			}
			meta: {
				goType: "rawJson"
				dal: { type: "JSON", defaultEmptyObject: true }
			}

			created_at: schema.SortableTimestampNowField
			updated_at: schema.SortableTimestampNilField
			deleted_at: schema.SortableTimestampNilField
			last_used_at: schema.SortableTimestampNilField
			expires_at: schema.SortableTimestampNilField
		}

		indexes: {
			"primary": { attribute: "id" }
			"owner_kind": {
				attributes: [ "owner_id", "kind" ]
				predicate: "deleted_at IS NULL"
			}
		}
	}

	filter: {
		struct: {
			owner_id: {goType: "uint64", ident: "ownerID", storeIdent: "rel_owner"}
			kind: {goType: "string"}
			credentials: {goType: "string"}
			deleted: {goType: "filter.State", storeIdent: "deleted_at"}
		}

		byValue: ["owner_id", "kind", "credentials"]
		byNilState: ["deleted"]
	}

	features: {
		labels: false
		paging: false
		sorting: false
		checkFn: false
	}

	store: {
		api: {
			lookups: [
				{
					fields: ["id"]
					description: """
						searches for credentials by ID

						It returns credentials even if deleted
						"""
				},
			]
		}
	}
}
