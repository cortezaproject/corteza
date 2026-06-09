package system

import (
	"github.com/crusttech/human/server/codegen/schema"
)

attachment: {
	features: {
		labels: false
		projectScoped: true
	}

	model: {
		attributes: {
			id: schema.IdField
			tenant_id:  schema.TenantRefField
			project_id: schema.ProjectRefField
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
				omitSetter: true
				omitGetter: true
			}
			created_at: schema.SortableTimestampNowField
			updated_at: schema.SortableTimestampNilField
			deleted_at: schema.SortableTimestampNilField
		}

		indexes: {
			"primary": { attribute: "id" }
		}
	}

	envoy: {
		omit: true
	}

	filter: {
		struct: {
			tenant_id: schema.TenantFilterField
			project_id: schema.ProjectFilterField
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
