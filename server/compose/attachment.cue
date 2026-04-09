package compose

import (
	"github.com/cortezaproject/corteza/server/codegen/schema"
)

attachment: {
	features: {
		labels: false
	}

	model: {
		ident: "compose_attachment"
		attributes: {
			id:       schema.IdField
			namespace_id: {
				ident: "namespaceID",
				goType: "uint64",
				storeIdent: "rel_namespace"
				dal: { type: "Ref", refModelResType: "corteza::compose:namespace" }
			}
			owner_id: {
				sortable: true,
				goType: "uint64",
				storeIdent: "rel_owner",
				ident: "ownerID"
				dal: { type: "Ref", refModelResType: "corteza::system:user" }
			}
			kind: {
				sortable: true
				dal: {}
			}
			url:  {
				dal: {}
			}
			preview_url: {
				dal: {}
			}
			name:        {
				sortable: true
				dal: {}
			}
			meta:        {
				goType: "types.AttachmentMeta"
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
			}
			hash: {
				dal: { type: "Text", nullable: true }
			}
			rel_module: {
				ident: "moduleID"
				goType: "uint64"
				storeIdent: "rel_module"
				dal: { type: "ID" }
			}
			created_at: schema.SortableTimestampNowField
			updated_at: schema.SortableTimestampNilField
			deleted_at: schema.SortableTimestampNilField
		}

		indexes: {
			"primary": { attribute: "id" }
			"namespace": { attribute: "namespace_id" },
			"hash_module": {
				fields: [{ attribute: "hash" }, { attribute: "rel_module" }]
			}
		}
	}

	filter: {
		struct: {
			kind: {}
			namespace_id: { goType: "uint64", ident: "namespaceID" }
			page_id: { goType: "uint64", ident: "pageID" }
			record_id: { goType: "uint64", ident: "recordID" }
			module_id: { goType: "uint64", ident: "moduleID" }
			field_name: { }
			hash: { }
		}

		byValue: ["kind", "namespace_id"]
	}

	envoy: {
		omit: true
	}

	store: {
		ident: "composeAttachment"

		api: {
			lookups: [
				{ fields: ["id"] },
			]
		}
	}
}
