package system

import (
	"github.com/crusttech/human/server/codegen/schema"
)

_attachmentDefs: {
			AttachmentMeta: { name: "AttachmentMeta", fields: [
				{ name: "Original", type: _attachmentDefs.AttachmentFileMeta, json: "original" },
				{ name: "Preview", ptr: true, type: _attachmentDefs.AttachmentFileMeta, json: "preview,omitempty" },
				{ name: "Labels", goType: "map[string]string", json: "labels,omitempty" },
			]}
			AttachmentFileMeta: { name: "AttachmentFileMeta", fields: [
				{ name: "Size", type: "int64", json: "size" },
				{ name: "Extension", type: "string", json: "ext" },
				{ name: "Mimetype", type: "string", json: "mimetype" },
				{ name: "Image", ptr: true, type: _attachmentDefs.AttachmentImageMeta, json: "image,omitempty" },
			]}
			AttachmentImageMeta: { name: "AttachmentImageMeta", fields: [
				{ name: "Width", type: "int", json: "width,omitempty" },
				{ name: "Height", type: "int", json: "height,omitempty" },
				{ name: "Animated", type: "bool", json: "animated" },
				{ name: "Initial", type: "string", json: "initial,omitempty" },
				{ name: "InitialColor", type: "string", json: "initial-color,omitempty" },
				{ name: "BackgroundColor", type: "string", json: "background-color,omitempty" },
			]}
		}

attachment: {
	features: {
		labels: false
	}

	types: {
		gen: true, jsonTypesSkip: ["AttachmentMeta"]
		defs: _attachmentDefs
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
				json: { field: "ownerID", string: true }
			}
			kind: {
				sortable: true
				dal: {}
				json: "-"
			}
			url: {
				dal: {}
				json: { field: "url", omitEmpty: true }
			}
			preview_url: {
				dal: {}
				json: { field: "previewUrl", omitEmpty: true }
			}
			name: {
				sortable: true
				dal: {}
				json: { field: "name", omitEmpty: true }
			}
			meta: {
				type: _attachmentDefs.AttachmentMeta
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
