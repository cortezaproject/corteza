package compose

import (
	"github.com/crusttech/human/server/codegen/schema"
)

attachment: {
	features: {
		projectScoped: true
		labels: false
	}

	model: {
		ident: "compose_attachment"
		attributes: {
			id:         schema.IdField
			tenant_id:  schema.TenantRefField
			project_id: schema.ProjectRefField
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
			created_at: schema.SortableTimestampNowField
			updated_at: schema.SortableTimestampNilField
			deleted_at: schema.SortableTimestampNilField
		}

		indexes: {
			"primary": { attribute: "id" }
			"namespace": { attribute: "namespace_id" },
		}
	}

	filter: {
		struct: {
			kind:         {}
			tenant_id:    schema.TenantFilterField
			project_id:   schema.ProjectFilterField
			namespace_id: { goType: "uint64", ident: "namespaceID" }
			page_id: { goType: "uint64", ident: "pageID" }
			record_id: { goType: "uint64", ident: "recordID" }
			module_id: { goType: "uint64", ident: "moduleID" }
			field_name: { }
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
