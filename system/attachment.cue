package system

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

attachment: schema.#Resource & {
	features: {
		labels: false
	}

	struct: {
		id:       schema.IdField
		owner_id: { goType: "uint64", storeIdent: "rel_owner", ident: "ownerID"}
		kind: {}
		url:  {}
		preview_url: {}
		name:        {}
		meta:        { goType: "types.AttachmentMeta" }
		created_at: schema.SortableTimestampField
		updated_at: schema.SortableTimestampNilField
		deleted_at: schema.SortableTimestampNilField
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
