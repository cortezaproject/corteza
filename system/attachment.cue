package system

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

attachment: {
	features: {
		labels: false
	}

	model: {
		attributes: {
			id: schema.IdField
			owner_id:   {
				storeIdent: "rel_owner",
				ident: "ownerID"
				schema.AttributeUserRef,
			}
			kind: {
				sortable: true
				dal: {}
			}
			url: {
				dal: {}
			}
			preview_url: {
				dal: {}
			}
			name: {
				sortable: true
				dal: {}
			}
			meta: {
				goType: "types.AttachmentMeta"
				dal: { type: "JSON", defaultEmptyObject: true }
			}
			created_at: schema.SortableTimestampNowField
			updated_at: schema.SortableTimestampNilField
			deleted_at: schema.SortableTimestampNilField
		}

		indexes: {
			"primary": { attribute: "id" }
		}
	}

	filter: {
		struct: {
			kind: {}
		}

		byValue: ["kind"]
	}

	store: {
		api: {
			lookups: [
				{ fields: ["id"] },
			]
		}
	}
}
