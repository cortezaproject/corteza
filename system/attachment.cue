package system

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

attachment: {
	features: {
		labels: false
	}

	model: {
		id:       schema.IdField
		owner_id: { sortable: true, goType: "uint64", storeIdent: "rel_owner", ident: "ownerID"}
		kind: {sortable: true}
		url:  {}
		preview_url: {}
		name:        {sortable: true}
		meta:        { goType: "types.AttachmentMeta" }
		created_at: schema.SortableTimestampField
		updated_at: schema.SortableTimestampNilField
		deleted_at: schema.SortableTimestampNilField
	}

	filter: {
		model: {
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
