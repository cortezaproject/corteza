package compose

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

attachment: {
	features: {
		labels: false
	}

	model: {
		id:       schema.IdField
		owner_id: { sortable: true, goType: "uint64", storeIdent: "rel_owner", ident: "ownerID" }
		namespace_id: { sortable: true, goType: "uint64", storeIdent: "rel_namespace", ident: "namespaceID" }
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
			namespace_id: { goType: "uint64", ident: "namespaceID" }
			page_id: { goType: "uint64", ident: "pageID" }
			record_id: { goType: "uint64", ident: "recordID" }
			module_id: { goType: "uint64", ident: "moduleID" }
			field_name: { }
		}

		byValue: ["kind", "namespace_id"]
	}

	store: {
		ident: "composeAttachment"

		settings: {
			rdbms: {
				table: "compose_attachment"
			}
		}
				api: {
			lookups: [
				{ fields: ["id"] },
			]
		}
	}
}
