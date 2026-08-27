package model

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/dal"
)

var Attachment = &dal.Model{
	Ident:        "compose_attachment",
	ResourceType: types.AttachmentResourceType,

	Attributes: dal.AttributeSet{
		&dal.Attribute{
			Ident: "ID",
			Type:  &dal.TypeID{},
			Store: &dal.CodecAlias{Ident: "id"},
		},

		&dal.Attribute{
			Ident: "NamespaceID",
			Type: &dal.TypeRef{
				RefAttribute: "id",
				RefModel: &dal.ModelRef{
					ResourceType: "corteza::compose:namespace",
				},
			},
			Store: &dal.CodecAlias{Ident: "rel_namespace"},
		},

		&dal.Attribute{
			Ident: "OwnerID", Sortable: true,
			Type: &dal.TypeRef{
				RefAttribute: "id",
				RefModel: &dal.ModelRef{
					ResourceType: "corteza::system:user",
				},
			},
			Store: &dal.CodecAlias{Ident: "rel_owner"},
		},

		&dal.Attribute{
			Ident: "Kind", Sortable: true,
			Type:  &dal.TypeText{},
			Store: &dal.CodecAlias{Ident: "kind"},
		},

		&dal.Attribute{
			Ident: "Url",
			Type:  &dal.TypeText{},
			Store: &dal.CodecAlias{Ident: "url"},
		},

		&dal.Attribute{
			Ident: "PreviewUrl",
			Type:  &dal.TypeText{},
			Store: &dal.CodecAlias{Ident: "preview_url"},
		},

		&dal.Attribute{
			Ident: "Name", Sortable: true,
			Type:  &dal.TypeText{},
			Store: &dal.CodecAlias{Ident: "name"},
		},

		&dal.Attribute{
			Ident: "Meta",
			Type: &dal.TypeJSON{
				DefaultValue: "{}",
			},
			Store: &dal.CodecAlias{Ident: "meta"},
		},

		&dal.Attribute{
			Ident: "CreatedAt", Sortable: true,
			Type: &dal.TypeTimestamp{
				DefaultCurrentTimestamp: true, Timezone: true, Precision: -1,
			},
			Store: &dal.CodecAlias{Ident: "created_at"},
		},

		&dal.Attribute{
			Ident: "UpdatedAt", Sortable: true,
			Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true, Precision: -1},
			Store: &dal.CodecAlias{Ident: "updated_at"},
		},

		&dal.Attribute{
			Ident: "DeletedAt", Sortable: true,
			Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true, Precision: -1},
			Store: &dal.CodecAlias{Ident: "deleted_at"},
		},
	},

	Indexes: dal.IndexSet{
		&dal.Index{
			Ident: "compose_attachment_namespace",
			Type:  "BTREE",

			Fields: []*dal.IndexField{
				{
					AttributeIdent: "NamespaceID",
				},
			},
		},

		&dal.Index{
			Ident: "PRIMARY",
			Type:  "BTREE",

			Fields: []*dal.IndexField{
				{
					AttributeIdent: "ID",
				},
			},
		},
	},
}

var Chart = &dal.Model{
	Ident:        "compose_chart",
	ResourceType: types.ChartResourceType,

	Attributes: dal.AttributeSet{
		&dal.Attribute{
			Ident: "ID",
			Type:  &dal.TypeID{},
			Store: &dal.CodecAlias{Ident: "id"},
		},

		&dal.Attribute{
			Ident: "Handle",
			Type:  &dal.TypeText{Length: 64},
			Store: &dal.CodecAlias{Ident: "handle"},
		},

		&dal.Attribute{
			Ident: "NamespaceID",
			Type: &dal.TypeRef{
				RefAttribute: "id",
				RefModel: &dal.ModelRef{
					ResourceType: "corteza::compose:namespace",
				},
			},
			Store: &dal.CodecAlias{Ident: "rel_namespace"},
		},

		&dal.Attribute{
			Ident: "Name", Sortable: true,
			Type:  &dal.TypeText{},
			Store: &dal.CodecAlias{Ident: "name"},
		},

		&dal.Attribute{
			Ident: "Config",
			Type:  &dal.TypeText{},
			Store: &dal.CodecAlias{Ident: "config"},
		},

		&dal.Attribute{
			Ident: "CreatedAt", Sortable: true,
			Type: &dal.TypeTimestamp{
				DefaultCurrentTimestamp: true, Timezone: true, Precision: -1,
			},
			Store: &dal.CodecAlias{Ident: "created_at"},
		},

		&dal.Attribute{
			Ident: "UpdatedAt", Sortable: true,
			Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true, Precision: -1},
			Store: &dal.CodecAlias{Ident: "updated_at"},
		},

		&dal.Attribute{
			Ident: "DeletedAt", Sortable: true,
			Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true, Precision: -1},
			Store: &dal.CodecAlias{Ident: "deleted_at"},
		},
	},

	Indexes: dal.IndexSet{
		&dal.Index{
			Ident: "compose_chart_namespace",
			Type:  "BTREE",

			Fields: []*dal.IndexField{
				{
					AttributeIdent: "NamespaceID",
				},
			},
		},

		&dal.Index{
			Ident: "PRIMARY",
			Type:  "BTREE",

			Fields: []*dal.IndexField{
				{
					AttributeIdent: "ID",
				},
			},
		},

		&dal.Index{
			Ident:     "compose_chart_uniqueHandle",
			Type:      "BTREE",
			Unique:    true,
			Predicate: "handle != '' AND deleted_at IS NULL",
			Fields: []*dal.IndexField{
				{
					AttributeIdent: "Handle",
					Modifiers:      []dal.IndexFieldModifier{"LOWERCASE"},
				},

				{
					AttributeIdent: "NamespaceID",
				},
			},
		},
	},
}

var City311ActorProfile = &dal.Model{
	Ident:        "compose_city311_actor_profile",
	ResourceType: types.City311ActorProfileResourceType,

	Attributes: dal.AttributeSet{
		&dal.Attribute{
			Ident: "ID",
			Type:  &dal.TypeID{},
			Store: &dal.CodecAlias{Ident: "id"},
		},

		&dal.Attribute{
			Ident: "ApplicationRoles",
			Type:  &dal.TypeJSON{},
			Store: &dal.CodecAlias{Ident: "application_roles"},
		},

		&dal.Attribute{
			Ident: "Department", Sortable: true,
			Type:  &dal.TypeText{Length: 64},
			Store: &dal.CodecAlias{Ident: "department"},
		},

		&dal.Attribute{
			Ident: "Districts",
			Type:  &dal.TypeJSON{},
			Store: &dal.CodecAlias{Ident: "districts"},
		},

		&dal.Attribute{
			Ident: "CreatedAt", Sortable: true,
			Type: &dal.TypeTimestamp{
				DefaultCurrentTimestamp: true, Timezone: true, Precision: -1,
			},
			Store: &dal.CodecAlias{Ident: "created_at"},
		},

		&dal.Attribute{
			Ident: "UpdatedAt", Sortable: true,
			Type: &dal.TypeTimestamp{
				DefaultCurrentTimestamp: true, Timezone: true, Precision: -1,
			},
			Store: &dal.CodecAlias{Ident: "updated_at"},
		},
	},

	Indexes: dal.IndexSet{
		&dal.Index{
			Ident: "compose_city311_actor_profile_department",
			Type:  "BTREE",

			Fields: []*dal.IndexField{
				{
					AttributeIdent: "Department",
				},
			},
		},

		&dal.Index{
			Ident: "PRIMARY",
			Type:  "BTREE",

			Fields: []*dal.IndexField{
				{
					AttributeIdent: "ID",
				},
			},
		},
	},
}

var City311AuditEvent = &dal.Model{
	Ident:        "compose_city311_audit_event",
	ResourceType: types.City311AuditEventResourceType,

	Attributes: dal.AttributeSet{
		&dal.Attribute{
			Ident: "ID",
			Type:  &dal.TypeID{},
			Store: &dal.CodecAlias{Ident: "id"},
		},

		&dal.Attribute{
			Ident: "RequestID", Sortable: true,
			Type:  &dal.TypeID{},
			Store: &dal.CodecAlias{Ident: "request_id"},
		},

		&dal.Attribute{
			Ident: "EventType", Sortable: true,
			Type:  &dal.TypeText{Length: 96},
			Store: &dal.CodecAlias{Ident: "event_type"},
		},

		&dal.Attribute{
			Ident: "ActorType",
			Type:  &dal.TypeText{Length: 32},
			Store: &dal.CodecAlias{Ident: "actor_type"},
		},

		&dal.Attribute{
			Ident: "ActorID",
			Type:  &dal.TypeID{},
			Store: &dal.CodecAlias{Ident: "actor_id"},
		},

		&dal.Attribute{
			Ident: "SourceChannel",
			Type:  &dal.TypeText{Length: 32},
			Store: &dal.CodecAlias{Ident: "source_channel"},
		},

		&dal.Attribute{
			Ident: "Before",
			Type: &dal.TypeJSON{
				DefaultValue: "{}",
			},
			Store: &dal.CodecAlias{Ident: "before"},
		},

		&dal.Attribute{
			Ident: "After",
			Type: &dal.TypeJSON{
				DefaultValue: "{}",
			},
			Store: &dal.CodecAlias{Ident: "after"},
		},

		&dal.Attribute{
			Ident: "CreatedAt", Sortable: true,
			Type: &dal.TypeTimestamp{
				DefaultCurrentTimestamp: true, Timezone: true, Precision: -1,
			},
			Store: &dal.CodecAlias{Ident: "created_at"},
		},
	},

	Indexes: dal.IndexSet{
		&dal.Index{
			Ident: "PRIMARY",
			Type:  "BTREE",

			Fields: []*dal.IndexField{
				{
					AttributeIdent: "ID",
				},
			},
		},

		&dal.Index{
			Ident: "compose_city311_audit_event_request",
			Type:  "BTREE",

			Fields: []*dal.IndexField{
				{
					AttributeIdent: "RequestID",
				},

				{
					AttributeIdent: "CreatedAt",
				},
			},
		},
	},
}

var City311IdempotencyRecord = &dal.Model{
	Ident:        "compose_city311_idempotency",
	ResourceType: types.City311IdempotencyRecordResourceType,

	Attributes: dal.AttributeSet{
		&dal.Attribute{
			Ident: "ID",
			Type:  &dal.TypeID{},
			Store: &dal.CodecAlias{Ident: "id"},
		},

		&dal.Attribute{
			Ident: "Operation",
			Type:  &dal.TypeText{Length: 96},
			Store: &dal.CodecAlias{Ident: "operation"},
		},

		&dal.Attribute{
			Ident: "KeyHash",
			Type:  &dal.TypeText{Length: 64},
			Store: &dal.CodecAlias{Ident: "key_hash"},
		},

		&dal.Attribute{
			Ident: "RequestHash",
			Type:  &dal.TypeText{Length: 64},
			Store: &dal.CodecAlias{Ident: "request_hash"},
		},

		&dal.Attribute{
			Ident: "ResponseStatus",
			Type:  &dal.TypeNumber{Precision: -1, Scale: -1, Meta: map[string]interface{}{"rdbms:type": "integer"}},
			Store: &dal.CodecAlias{Ident: "response_status"},
		},

		&dal.Attribute{
			Ident: "ResponseBody",
			Type: &dal.TypeJSON{
				DefaultValue: "{}",
			},
			Store: &dal.CodecAlias{Ident: "response_body"},
		},

		&dal.Attribute{
			Ident: "RequestID",
			Type:  &dal.TypeID{},
			Store: &dal.CodecAlias{Ident: "request_id"},
		},

		&dal.Attribute{
			Ident: "CreatedAt", Sortable: true,
			Type: &dal.TypeTimestamp{
				DefaultCurrentTimestamp: true, Timezone: true, Precision: -1,
			},
			Store: &dal.CodecAlias{Ident: "created_at"},
		},

		&dal.Attribute{
			Ident: "ExpiresAt", Sortable: true,
			Type:  &dal.TypeTimestamp{Timezone: true, Precision: -1},
			Store: &dal.CodecAlias{Ident: "expires_at"},
		},
	},

	Indexes: dal.IndexSet{
		&dal.Index{
			Ident: "compose_city311_idempotency_expires",
			Type:  "BTREE",

			Fields: []*dal.IndexField{
				{
					AttributeIdent: "ExpiresAt",
				},
			},
		},

		&dal.Index{
			Ident: "PRIMARY",
			Type:  "BTREE",

			Fields: []*dal.IndexField{
				{
					AttributeIdent: "ID",
				},
			},
		},

		&dal.Index{
			Ident:  "compose_city311_idempotency_uniqueOperationKey",
			Type:   "BTREE",
			Unique: true,

			Fields: []*dal.IndexField{
				{
					AttributeIdent: "Operation",
				},

				{
					AttributeIdent: "KeyHash",
				},
			},
		},
	},
}

var City311PublicHistoryItem = &dal.Model{
	Ident:        "compose_city311_public_history",
	ResourceType: types.City311PublicHistoryItemResourceType,

	Attributes: dal.AttributeSet{
		&dal.Attribute{
			Ident: "ID",
			Type:  &dal.TypeID{},
			Store: &dal.CodecAlias{Ident: "id"},
		},

		&dal.Attribute{
			Ident: "RequestID", Sortable: true,
			Type:  &dal.TypeID{},
			Store: &dal.CodecAlias{Ident: "request_id"},
		},

		&dal.Attribute{
			Ident: "Action",
			Type:  &dal.TypeText{Length: 64},
			Store: &dal.CodecAlias{Ident: "action"},
		},

		&dal.Attribute{
			Ident: "ResponsibleDepartment",
			Type:  &dal.TypeText{Length: 64},
			Store: &dal.CodecAlias{Ident: "responsible_department"},
		},

		&dal.Attribute{
			Ident: "OccurredAt", Sortable: true,
			Type: &dal.TypeTimestamp{
				DefaultCurrentTimestamp: true, Timezone: true, Precision: -1,
			},
			Store: &dal.CodecAlias{Ident: "occurred_at"},
		},
	},

	Indexes: dal.IndexSet{
		&dal.Index{
			Ident: "PRIMARY",
			Type:  "BTREE",

			Fields: []*dal.IndexField{
				{
					AttributeIdent: "ID",
				},
			},
		},

		&dal.Index{
			Ident: "compose_city311_public_history_request",
			Type:  "BTREE",

			Fields: []*dal.IndexField{
				{
					AttributeIdent: "RequestID",
				},

				{
					AttributeIdent: "OccurredAt",
				},
			},
		},
	},
}

var City311RequestAttachment = &dal.Model{
	Ident:        "compose_city311_request_attachment",
	ResourceType: types.City311RequestAttachmentResourceType,

	Attributes: dal.AttributeSet{
		&dal.Attribute{
			Ident: "ID",
			Type:  &dal.TypeID{},
			Store: &dal.CodecAlias{Ident: "id"},
		},

		&dal.Attribute{
			Ident: "RequestID", Sortable: true,
			Type:  &dal.TypeID{},
			Store: &dal.CodecAlias{Ident: "request_id"},
		},

		&dal.Attribute{
			Ident: "Filename",
			Type:  &dal.TypeText{Length: 120},
			Store: &dal.CodecAlias{Ident: "filename"},
		},

		&dal.Attribute{
			Ident: "MediaType",
			Type:  &dal.TypeText{Length: 128},
			Store: &dal.CodecAlias{Ident: "media_type"},
		},

		&dal.Attribute{
			Ident: "Size",
			Type:  &dal.TypeNumber{Precision: -1, Scale: -1, Meta: map[string]interface{}{"rdbms:type": "bigint"}},
			Store: &dal.CodecAlias{Ident: "size"},
		},

		&dal.Attribute{
			Ident: "Content",
			Type:  &dal.TypeBlob{},
			Store: &dal.CodecAlias{Ident: "content"},
		},

		&dal.Attribute{
			Ident: "CreatedAt", Sortable: true,
			Type: &dal.TypeTimestamp{
				DefaultCurrentTimestamp: true, Timezone: true, Precision: -1,
			},
			Store: &dal.CodecAlias{Ident: "created_at"},
		},
	},

	Indexes: dal.IndexSet{
		&dal.Index{
			Ident: "PRIMARY",
			Type:  "BTREE",

			Fields: []*dal.IndexField{
				{
					AttributeIdent: "ID",
				},
			},
		},

		&dal.Index{
			Ident: "compose_city311_request_attachment_request",
			Type:  "BTREE",

			Fields: []*dal.IndexField{
				{
					AttributeIdent: "RequestID",
				},

				{
					AttributeIdent: "CreatedAt",
				},
			},
		},
	},
}

var City311RequestSequence = &dal.Model{
	Ident:        "compose_city311_request_sequence",
	ResourceType: types.City311RequestSequenceResourceType,

	Attributes: dal.AttributeSet{
		&dal.Attribute{
			Ident: "ID",
			Type:  &dal.TypeID{},
			Store: &dal.CodecAlias{Ident: "id"},
		},

		&dal.Attribute{
			Ident: "NextNumber",
			Type:  &dal.TypeNumber{Precision: -1, Scale: -1, Meta: map[string]interface{}{"rdbms:type": "bigint"}},
			Store: &dal.CodecAlias{Ident: "next_number"},
		},
	},

	Indexes: dal.IndexSet{
		&dal.Index{
			Ident: "PRIMARY",
			Type:  "BTREE",

			Fields: []*dal.IndexField{
				{
					AttributeIdent: "ID",
				},
			},
		},
	},
}

var City311ServiceRequest = &dal.Model{
	Ident:        "compose_city311_service_request",
	ResourceType: types.City311ServiceRequestResourceType,

	Attributes: dal.AttributeSet{
		&dal.Attribute{
			Ident: "ID",
			Type:  &dal.TypeID{},
			Store: &dal.CodecAlias{Ident: "id"},
		},

		&dal.Attribute{
			Ident: "RequestNumber", Sortable: true,
			Type:  &dal.TypeText{Length: 16},
			Store: &dal.CodecAlias{Ident: "request_number"},
		},

		&dal.Attribute{
			Ident: "Summary", Sortable: true,
			Type:  &dal.TypeText{Length: 160},
			Store: &dal.CodecAlias{Ident: "summary"},
		},

		&dal.Attribute{
			Ident: "Description",
			Type:  &dal.TypeText{},
			Store: &dal.CodecAlias{Ident: "description"},
		},

		&dal.Attribute{
			Ident: "ServiceType", Sortable: true,
			Type:  &dal.TypeText{Length: 64},
			Store: &dal.CodecAlias{Ident: "service_type"},
		},

		&dal.Attribute{
			Ident: "OwningDepartment", Sortable: true,
			Type:  &dal.TypeText{Length: 64},
			Store: &dal.CodecAlias{Ident: "owning_department"},
		},

		&dal.Attribute{
			Ident: "CouncilDistrict", Sortable: true,
			Type:  &dal.TypeText{Length: 32},
			Store: &dal.CodecAlias{Ident: "council_district"},
		},

		&dal.Attribute{
			Ident: "SourceChannel", Sortable: true,
			Type:  &dal.TypeText{Length: 32},
			Store: &dal.CodecAlias{Ident: "source_channel"},
		},

		&dal.Attribute{
			Ident: "OriginClass", Sortable: true,
			Type:  &dal.TypeText{Length: 16},
			Store: &dal.CodecAlias{Ident: "origin_class"},
		},

		&dal.Attribute{
			Ident: "Status", Sortable: true,
			Type:  &dal.TypeText{Length: 32},
			Store: &dal.CodecAlias{Ident: "status"},
		},

		&dal.Attribute{
			Ident: "PrimaryRequester",
			Type: &dal.TypeJSON{
				DefaultValue: "{}",
			},
			Store: &dal.CodecAlias{Ident: "primary_requester"},
		},

		&dal.Attribute{
			Ident: "Location",
			Type: &dal.TypeJSON{
				DefaultValue: "{}",
			},
			Store: &dal.CodecAlias{Ident: "location"},
		},

		&dal.Attribute{
			Ident: "CustomFields",
			Type: &dal.TypeJSON{
				DefaultValue: "{}",
			},
			Store: &dal.CodecAlias{Ident: "custom_fields"},
		},

		&dal.Attribute{
			Ident: "PrimaryAssigneeID", Sortable: true,
			Type: &dal.TypeID{HasDefault: true,
				DefaultValue: 0,
			},
			Store: &dal.CodecAlias{Ident: "primary_assignee_id"},
		},

		&dal.Attribute{
			Ident: "CollaboratorIDs",
			Type:  &dal.TypeJSON{},
			Store: &dal.CodecAlias{Ident: "collaborator_ids"},
		},

		&dal.Attribute{
			Ident: "DuplicateGroupID", Sortable: true,
			Type:  &dal.TypeText{Length: 64},
			Store: &dal.CodecAlias{Ident: "duplicate_group_id"},
		},

		&dal.Attribute{
			Ident: "Version", Sortable: true,
			Type: &dal.TypeNumber{HasDefault: true,
				DefaultValue: 1,
				Precision:    -1, Scale: -1, Meta: map[string]interface{}{"rdbms:type": "integer"},
			},
			Store: &dal.CodecAlias{Ident: "version"},
		},

		&dal.Attribute{
			Ident: "CreatedAt", Sortable: true,
			Type: &dal.TypeTimestamp{
				DefaultCurrentTimestamp: true, Timezone: true, Precision: -1,
			},
			Store: &dal.CodecAlias{Ident: "created_at"},
		},

		&dal.Attribute{
			Ident: "UpdatedAt", Sortable: true,
			Type: &dal.TypeTimestamp{
				DefaultCurrentTimestamp: true, Timezone: true, Precision: -1,
			},
			Store: &dal.CodecAlias{Ident: "updated_at"},
		},
	},

	Indexes: dal.IndexSet{
		&dal.Index{
			Ident: "PRIMARY",
			Type:  "BTREE",

			Fields: []*dal.IndexField{
				{
					AttributeIdent: "ID",
				},
			},
		},

		&dal.Index{
			Ident: "compose_city311_service_request_requestScope",
			Type:  "BTREE",

			Fields: []*dal.IndexField{
				{
					AttributeIdent: "OwningDepartment",
				},

				{
					AttributeIdent: "CouncilDistrict",
				},

				{
					AttributeIdent: "Status",
				},
			},
		},

		&dal.Index{
			Ident:  "compose_city311_service_request_uniqueRequestNumber",
			Type:   "BTREE",
			Unique: true,

			Fields: []*dal.IndexField{
				{
					AttributeIdent: "RequestNumber",
				},
			},
		},
	},
}

var Module = &dal.Model{
	Ident:        "compose_module",
	ResourceType: types.ModuleResourceType,

	Attributes: dal.AttributeSet{
		&dal.Attribute{
			Ident: "ID",
			Type:  &dal.TypeID{},
			Store: &dal.CodecAlias{Ident: "id"},
		},

		&dal.Attribute{
			Ident: "NamespaceID",
			Type: &dal.TypeRef{
				RefAttribute: "id",
				RefModel: &dal.ModelRef{
					ResourceType: "corteza::compose:namespace",
				},
			},
			Store: &dal.CodecAlias{Ident: "rel_namespace"},
		},

		&dal.Attribute{
			Ident: "Handle",
			Type:  &dal.TypeText{Length: 64},
			Store: &dal.CodecAlias{Ident: "handle"},
		},

		&dal.Attribute{
			Ident: "Name", Sortable: true,
			Type:  &dal.TypeText{},
			Store: &dal.CodecAlias{Ident: "name"},
		},

		&dal.Attribute{
			Ident: "Meta",
			Type: &dal.TypeJSON{
				DefaultValue: "{}",
			},
			Store: &dal.CodecAlias{Ident: "meta"},
		},

		&dal.Attribute{
			Ident: "Config",
			Type: &dal.TypeJSON{
				DefaultValue: "{}",
			},
			Store: &dal.CodecAlias{Ident: "config"},
		},

		&dal.Attribute{
			Ident: "CreatedAt", Sortable: true,
			Type: &dal.TypeTimestamp{
				DefaultCurrentTimestamp: true, Timezone: true, Precision: -1,
			},
			Store: &dal.CodecAlias{Ident: "created_at"},
		},

		&dal.Attribute{
			Ident: "UpdatedAt", Sortable: true,
			Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true, Precision: -1},
			Store: &dal.CodecAlias{Ident: "updated_at"},
		},

		&dal.Attribute{
			Ident: "DeletedAt", Sortable: true,
			Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true, Precision: -1},
			Store: &dal.CodecAlias{Ident: "deleted_at"},
		},
	},

	Indexes: dal.IndexSet{
		&dal.Index{
			Ident: "compose_module_namespace",
			Type:  "BTREE",

			Fields: []*dal.IndexField{
				{
					AttributeIdent: "NamespaceID",
				},
			},
		},

		&dal.Index{
			Ident: "PRIMARY",
			Type:  "BTREE",

			Fields: []*dal.IndexField{
				{
					AttributeIdent: "ID",
				},
			},
		},

		&dal.Index{
			Ident:     "compose_module_uniqueHandle",
			Type:      "BTREE",
			Unique:    true,
			Predicate: "handle != '' AND deleted_at IS NULL",
			Fields: []*dal.IndexField{
				{
					AttributeIdent: "Handle",
					Modifiers:      []dal.IndexFieldModifier{"LOWERCASE"},
				},

				{
					AttributeIdent: "NamespaceID",
				},
			},
		},
	},
}

var ModuleField = &dal.Model{
	Ident:        "compose_module_field",
	ResourceType: types.ModuleFieldResourceType,

	Attributes: dal.AttributeSet{
		&dal.Attribute{
			Ident: "ID",
			Type:  &dal.TypeID{},
			Store: &dal.CodecAlias{Ident: "id"},
		},

		&dal.Attribute{
			Ident: "ModuleID",
			Type: &dal.TypeRef{
				RefAttribute: "id",
				RefModel: &dal.ModelRef{
					ResourceType: "corteza::compose:module",
				},
			},
			Store: &dal.CodecAlias{Ident: "rel_module"},
		},

		&dal.Attribute{
			Ident: "Place", Sortable: true,
			Type:  &dal.TypeNumber{Precision: -1, Scale: -1, Meta: map[string]interface{}{"rdbms:type": "integer"}},
			Store: &dal.CodecAlias{Ident: "place"},
		},

		&dal.Attribute{
			Ident: "Kind", Sortable: true,
			Type:  &dal.TypeText{},
			Store: &dal.CodecAlias{Ident: "kind"},
		},

		&dal.Attribute{
			Ident: "Options",
			Type: &dal.TypeJSON{
				DefaultValue: "{}",
			},
			Store: &dal.CodecAlias{Ident: "options"},
		},

		&dal.Attribute{
			Ident: "Name", Sortable: true,
			Type:  &dal.TypeText{},
			Store: &dal.CodecAlias{Ident: "name"},
		},

		&dal.Attribute{
			Ident: "Label", Sortable: true,
			Type:  &dal.TypeText{},
			Store: &dal.CodecAlias{Ident: "label"},
		},

		&dal.Attribute{
			Ident: "Config",
			Type: &dal.TypeJSON{
				DefaultValue: "{}",
			},
			Store: &dal.CodecAlias{Ident: "config"},
		},

		&dal.Attribute{
			Ident: "Required",
			Type:  &dal.TypeBoolean{},
			Store: &dal.CodecAlias{Ident: "is_required"},
		},

		&dal.Attribute{
			Ident: "Multi",
			Type:  &dal.TypeBoolean{},
			Store: &dal.CodecAlias{Ident: "is_multi"},
		},

		&dal.Attribute{
			Ident: "DefaultValue",
			Type: &dal.TypeJSON{
				DefaultValue: "{}",
			},
			Store: &dal.CodecAlias{Ident: "default_value"},
		},

		&dal.Attribute{
			Ident: "Expressions",
			Type: &dal.TypeJSON{
				DefaultValue: "{}",
			},
			Store: &dal.CodecAlias{Ident: "expressions"},
		},

		&dal.Attribute{
			Ident: "CreatedAt", Sortable: true,
			Type: &dal.TypeTimestamp{
				DefaultCurrentTimestamp: true, Timezone: true, Precision: -1,
			},
			Store: &dal.CodecAlias{Ident: "created_at"},
		},

		&dal.Attribute{
			Ident: "UpdatedAt", Sortable: true,
			Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true, Precision: -1},
			Store: &dal.CodecAlias{Ident: "updated_at"},
		},

		&dal.Attribute{
			Ident: "DeletedAt", Sortable: true,
			Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true, Precision: -1},
			Store: &dal.CodecAlias{Ident: "deleted_at"},
		},
	},

	Indexes: dal.IndexSet{
		&dal.Index{
			Ident: "compose_module_field_module",
			Type:  "BTREE",

			Fields: []*dal.IndexField{
				{
					AttributeIdent: "ModuleID",
				},
			},
		},

		&dal.Index{
			Ident: "PRIMARY",
			Type:  "BTREE",

			Fields: []*dal.IndexField{
				{
					AttributeIdent: "ID",
				},
			},
		},

		&dal.Index{
			Ident:     "compose_module_field_uniqueName",
			Type:      "BTREE",
			Unique:    true,
			Predicate: "name != '' AND deleted_at IS NULL",
			Fields: []*dal.IndexField{
				{
					AttributeIdent: "Name",
					Modifiers:      []dal.IndexFieldModifier{"LOWERCASE"},
				},

				{
					AttributeIdent: "ModuleID",
				},
			},
		},
	},
}

var Namespace = &dal.Model{
	Ident:        "compose_namespace",
	ResourceType: types.NamespaceResourceType,

	Attributes: dal.AttributeSet{
		&dal.Attribute{
			Ident: "ID",
			Type:  &dal.TypeID{},
			Store: &dal.CodecAlias{Ident: "id"},
		},

		&dal.Attribute{
			Ident: "Slug", Sortable: true,
			Type:  &dal.TypeText{},
			Store: &dal.CodecAlias{Ident: "slug"},
		},

		&dal.Attribute{
			Ident: "Enabled",
			Type:  &dal.TypeBoolean{},
			Store: &dal.CodecAlias{Ident: "enabled"},
		},

		&dal.Attribute{
			Ident: "Meta",
			Type: &dal.TypeJSON{
				DefaultValue: "{}",
			},
			Store: &dal.CodecAlias{Ident: "meta"},
		},

		&dal.Attribute{
			Ident: "Name", Sortable: true,
			Type:  &dal.TypeText{},
			Store: &dal.CodecAlias{Ident: "name"},
		},

		&dal.Attribute{
			Ident: "CreatedAt", Sortable: true,
			Type: &dal.TypeTimestamp{
				DefaultCurrentTimestamp: true, Timezone: true, Precision: -1,
			},
			Store: &dal.CodecAlias{Ident: "created_at"},
		},

		&dal.Attribute{
			Ident: "UpdatedAt", Sortable: true,
			Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true, Precision: -1},
			Store: &dal.CodecAlias{Ident: "updated_at"},
		},

		&dal.Attribute{
			Ident: "DeletedAt", Sortable: true,
			Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true, Precision: -1},
			Store: &dal.CodecAlias{Ident: "deleted_at"},
		},
	},

	Indexes: dal.IndexSet{
		&dal.Index{
			Ident: "PRIMARY",
			Type:  "BTREE",

			Fields: []*dal.IndexField{
				{
					AttributeIdent: "ID",
				},
			},
		},

		&dal.Index{
			Ident:     "compose_namespace_uniqueHandle",
			Type:      "BTREE",
			Unique:    true,
			Predicate: "slug != '' AND deleted_at IS NULL",
			Fields: []*dal.IndexField{
				{
					AttributeIdent: "Slug",
					Modifiers:      []dal.IndexFieldModifier{"LOWERCASE"},
				},
			},
		},
	},
}

var Page = &dal.Model{
	Ident:        "compose_page",
	ResourceType: types.PageResourceType,

	Attributes: dal.AttributeSet{
		&dal.Attribute{
			Ident: "ID",
			Type:  &dal.TypeID{},
			Store: &dal.CodecAlias{Ident: "id"},
		},

		&dal.Attribute{
			Ident: "Title", Sortable: true,
			Type:  &dal.TypeText{},
			Store: &dal.CodecAlias{Ident: "title"},
		},

		&dal.Attribute{
			Ident: "Handle",
			Type:  &dal.TypeText{Length: 64},
			Store: &dal.CodecAlias{Ident: "handle"},
		},

		&dal.Attribute{
			Ident: "SelfID", Sortable: true,
			Type: &dal.TypeRef{
				RefAttribute: "id",
				RefModel: &dal.ModelRef{
					ResourceType: "corteza::compose:page",
				},
			},
			Store: &dal.CodecAlias{Ident: "self_id"},
		},

		&dal.Attribute{
			Ident: "ModuleID",
			Type: &dal.TypeRef{
				RefAttribute: "id",
				RefModel: &dal.ModelRef{
					ResourceType: "corteza::compose:module",
				},
			},
			Store: &dal.CodecAlias{Ident: "rel_module"},
		},

		&dal.Attribute{
			Ident: "NamespaceID",
			Type: &dal.TypeRef{
				RefAttribute: "id",
				RefModel: &dal.ModelRef{
					ResourceType: "corteza::compose:namespace",
				},
			},
			Store: &dal.CodecAlias{Ident: "rel_namespace"},
		},

		&dal.Attribute{
			Ident: "Meta",
			Type: &dal.TypeJSON{
				DefaultValue: "{}",
			},
			Store: &dal.CodecAlias{Ident: "meta"},
		},

		&dal.Attribute{
			Ident: "Config",
			Type: &dal.TypeJSON{
				DefaultValue: "{}",
			},
			Store: &dal.CodecAlias{Ident: "config"},
		},

		&dal.Attribute{
			Ident: "Blocks",
			Type: &dal.TypeJSON{
				DefaultValue: "{}",
			},
			Store: &dal.CodecAlias{Ident: "blocks"},
		},

		&dal.Attribute{
			Ident: "Visible",
			Type: &dal.TypeBoolean{HasDefault: true,
				DefaultValue: true,
			},
			Store: &dal.CodecAlias{Ident: "visible"},
		},

		&dal.Attribute{
			Ident: "Weight", Sortable: true,
			Type: &dal.TypeNumber{HasDefault: true,
				DefaultValue: 0,
				Precision:    -1, Scale: -1, Meta: map[string]interface{}{"rdbms:type": "integer"},
			},
			Store: &dal.CodecAlias{Ident: "weight"},
		},

		&dal.Attribute{
			Ident: "Description",
			Type:  &dal.TypeText{},
			Store: &dal.CodecAlias{Ident: "description"},
		},

		&dal.Attribute{
			Ident: "CreatedAt", Sortable: true,
			Type: &dal.TypeTimestamp{
				DefaultCurrentTimestamp: true, Timezone: true, Precision: -1,
			},
			Store: &dal.CodecAlias{Ident: "created_at"},
		},

		&dal.Attribute{
			Ident: "UpdatedAt", Sortable: true,
			Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true, Precision: -1},
			Store: &dal.CodecAlias{Ident: "updated_at"},
		},

		&dal.Attribute{
			Ident: "DeletedAt", Sortable: true,
			Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true, Precision: -1},
			Store: &dal.CodecAlias{Ident: "deleted_at"},
		},
	},

	Indexes: dal.IndexSet{
		&dal.Index{
			Ident: "compose_page_module",
			Type:  "BTREE",

			Fields: []*dal.IndexField{
				{
					AttributeIdent: "ModuleID",
				},
			},
		},

		&dal.Index{
			Ident: "compose_page_namespace",
			Type:  "BTREE",

			Fields: []*dal.IndexField{
				{
					AttributeIdent: "NamespaceID",
				},
			},
		},

		&dal.Index{
			Ident: "PRIMARY",
			Type:  "BTREE",

			Fields: []*dal.IndexField{
				{
					AttributeIdent: "ID",
				},
			},
		},

		&dal.Index{
			Ident: "compose_page_selfId",
			Type:  "BTREE",

			Fields: []*dal.IndexField{
				{
					AttributeIdent: "SelfID",
				},
			},
		},

		&dal.Index{
			Ident:     "compose_page_uniqueHandle",
			Type:      "BTREE",
			Unique:    true,
			Predicate: "handle != '' AND deleted_at IS NULL",
			Fields: []*dal.IndexField{
				{
					AttributeIdent: "Handle",
					Modifiers:      []dal.IndexFieldModifier{"LOWERCASE"},
				},

				{
					AttributeIdent: "NamespaceID",
				},
			},
		},
	},
}

var PageLayout = &dal.Model{
	Ident:        "compose_page_layout",
	ResourceType: types.PageLayoutResourceType,

	Attributes: dal.AttributeSet{
		&dal.Attribute{
			Ident: "ID",
			Type:  &dal.TypeID{},
			Store: &dal.CodecAlias{Ident: "id"},
		},

		&dal.Attribute{
			Ident: "Handle",
			Type:  &dal.TypeText{Length: 64},
			Store: &dal.CodecAlias{Ident: "handle"},
		},

		&dal.Attribute{
			Ident: "PageID", Sortable: true,
			Type: &dal.TypeRef{
				RefAttribute: "id",
				RefModel: &dal.ModelRef{
					ResourceType: "corteza::compose:page",
				},
			},
			Store: &dal.CodecAlias{Ident: "page_id"},
		},

		&dal.Attribute{
			Ident: "ParentID", Sortable: true,
			Type: &dal.TypeRef{
				RefAttribute: "id",
				RefModel: &dal.ModelRef{
					ResourceType: "corteza::compose:page-layout",
				},
			},
			Store: &dal.CodecAlias{Ident: "parent_id"},
		},

		&dal.Attribute{
			Ident: "NamespaceID",
			Type: &dal.TypeRef{
				RefAttribute: "id",
				RefModel: &dal.ModelRef{
					ResourceType: "corteza::compose:namespace",
				},
			},
			Store: &dal.CodecAlias{Ident: "rel_namespace"},
		},

		&dal.Attribute{
			Ident: "Weight", Sortable: true,
			Type: &dal.TypeNumber{HasDefault: true,
				DefaultValue: 0,
				Precision:    -1, Scale: -1, Meta: map[string]interface{}{"rdbms:type": "integer"},
			},
			Store: &dal.CodecAlias{Ident: "weight"},
		},

		&dal.Attribute{
			Ident: "Meta",
			Type: &dal.TypeJSON{
				DefaultValue: "{}",
			},
			Store: &dal.CodecAlias{Ident: "meta"},
		},

		&dal.Attribute{
			Ident: "Config",
			Type: &dal.TypeJSON{
				DefaultValue: "{}",
			},
			Store: &dal.CodecAlias{Ident: "config"},
		},

		&dal.Attribute{
			Ident: "Blocks",
			Type: &dal.TypeJSON{
				DefaultValue: "{}",
			},
			Store: &dal.CodecAlias{Ident: "blocks"},
		},

		&dal.Attribute{
			Ident: "OwnedBy",
			Type: &dal.TypeRef{HasDefault: true,
				DefaultValue: 0,

				RefAttribute: "id",
				RefModel: &dal.ModelRef{
					ResourceType: "corteza::system:user",
				},
			},
			Store: &dal.CodecAlias{Ident: "owned_by"},
		},

		&dal.Attribute{
			Ident: "CreatedAt", Sortable: true,
			Type: &dal.TypeTimestamp{
				DefaultCurrentTimestamp: true, Timezone: true, Precision: -1,
			},
			Store: &dal.CodecAlias{Ident: "created_at"},
		},

		&dal.Attribute{
			Ident: "UpdatedAt", Sortable: true,
			Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true, Precision: -1},
			Store: &dal.CodecAlias{Ident: "updated_at"},
		},

		&dal.Attribute{
			Ident: "DeletedAt", Sortable: true,
			Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true, Precision: -1},
			Store: &dal.CodecAlias{Ident: "deleted_at"},
		},
	},

	Indexes: dal.IndexSet{
		&dal.Index{
			Ident: "compose_page_layout_namespace",
			Type:  "BTREE",

			Fields: []*dal.IndexField{
				{
					AttributeIdent: "NamespaceID",
				},
			},
		},

		&dal.Index{
			Ident: "compose_page_layout_pageId",
			Type:  "BTREE",

			Fields: []*dal.IndexField{
				{
					AttributeIdent: "PageID",
				},
			},
		},

		&dal.Index{
			Ident: "compose_page_layout_parentId",
			Type:  "BTREE",

			Fields: []*dal.IndexField{
				{
					AttributeIdent: "ParentID",
				},
			},
		},

		&dal.Index{
			Ident: "PRIMARY",
			Type:  "BTREE",

			Fields: []*dal.IndexField{
				{
					AttributeIdent: "ID",
				},
			},
		},

		&dal.Index{
			Ident:     "compose_page_layout_uniqueHandle",
			Type:      "BTREE",
			Unique:    true,
			Predicate: "handle != '' AND deleted_at IS NULL",
			Fields: []*dal.IndexField{
				{
					AttributeIdent: "Handle",
					Modifiers:      []dal.IndexFieldModifier{"LOWERCASE"},
				},

				{
					AttributeIdent: "PageID",
				},

				{
					AttributeIdent: "NamespaceID",
				},
			},
		},
	},
}

var Record = &dal.Model{
	Ident:        "compose_record",
	ResourceType: types.RecordResourceType,

	Attributes: dal.AttributeSet{
		&dal.Attribute{
			Ident: "ID",
			Type:  &dal.TypeID{},
			Store: &dal.CodecAlias{Ident: "id"},
		},

		&dal.Attribute{
			Ident: "Revision",
			Type: &dal.TypeNumber{HasDefault: true,
				DefaultValue: 0,
				Precision:    -1, Scale: -1, Meta: map[string]interface{}{"rdbms:type": "integer"},
			},
			Store: &dal.CodecAlias{Ident: "revision"},
		},

		&dal.Attribute{
			Ident: "ModuleID",
			Type: &dal.TypeRef{
				RefAttribute: "id",
				RefModel: &dal.ModelRef{
					ResourceType: "corteza::compose:module",
				},
			},
			Store: &dal.CodecAlias{Ident: "rel_module"},
		},

		&dal.Attribute{
			Ident: "Values",
			Type: &dal.TypeJSON{
				DefaultValue: "{}",
			},
			Store: &dal.CodecAlias{Ident: "values"},
		},

		&dal.Attribute{
			Ident: "Meta",
			Type: &dal.TypeJSON{
				DefaultValue: "{}",
			},
			Store: &dal.CodecAlias{Ident: "meta"},
		},

		&dal.Attribute{
			Ident: "NamespaceID",
			Type: &dal.TypeRef{
				RefAttribute: "id",
				RefModel: &dal.ModelRef{
					ResourceType: "corteza::compose:namespace",
				},
			},
			Store: &dal.CodecAlias{Ident: "rel_namespace"},
		},

		&dal.Attribute{
			Ident: "CreatedAt", Sortable: true,
			Type: &dal.TypeTimestamp{
				DefaultCurrentTimestamp: true, Timezone: true, Precision: -1,
			},
			Store: &dal.CodecAlias{Ident: "created_at"},
		},

		&dal.Attribute{
			Ident: "UpdatedAt", Sortable: true,
			Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true, Precision: -1},
			Store: &dal.CodecAlias{Ident: "updated_at"},
		},

		&dal.Attribute{
			Ident: "DeletedAt", Sortable: true,
			Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true, Precision: -1},
			Store: &dal.CodecAlias{Ident: "deleted_at"},
		},

		&dal.Attribute{
			Ident: "OwnedBy",
			Type: &dal.TypeRef{HasDefault: true,
				DefaultValue: 0,

				RefAttribute: "id",
				RefModel: &dal.ModelRef{
					ResourceType: "corteza::system:user",
				},
			},
			Store: &dal.CodecAlias{Ident: "owned_by"},
		},

		&dal.Attribute{
			Ident: "CreatedBy",
			Type: &dal.TypeRef{HasDefault: true,
				DefaultValue: 0,

				RefAttribute: "id",
				RefModel: &dal.ModelRef{
					ResourceType: "corteza::system:user",
				},
			},
			Store: &dal.CodecAlias{Ident: "created_by"},
		},

		&dal.Attribute{
			Ident: "UpdatedBy",
			Type: &dal.TypeRef{HasDefault: true,
				DefaultValue: 0,

				RefAttribute: "id",
				RefModel: &dal.ModelRef{
					ResourceType: "corteza::system:user",
				},
			},
			Store: &dal.CodecAlias{Ident: "updated_by"},
		},

		&dal.Attribute{
			Ident: "DeletedBy",
			Type: &dal.TypeRef{HasDefault: true,
				DefaultValue: 0,

				RefAttribute: "id",
				RefModel: &dal.ModelRef{
					ResourceType: "corteza::system:user",
				},
			},
			Store: &dal.CodecAlias{Ident: "deleted_by"},
		},
	},

	Indexes: dal.IndexSet{
		&dal.Index{
			Ident: "compose_record_idxComposeRecordBase",
			Type:  "BTREE",

			Predicate: "deleted_at IS NULL",
			Fields: []*dal.IndexField{
				{
					AttributeIdent: "ModuleID",
				},

				{
					AttributeIdent: "NamespaceID",
				},
			},
		},

		&dal.Index{
			Ident: "PRIMARY",
			Type:  "BTREE",

			Fields: []*dal.IndexField{
				{
					AttributeIdent: "ID",
				},
			},
		},
	},
}

var RecordRevision = &dal.Model{
	Ident:        "compose_record_revisions",
	ResourceType: types.RecordRevisionResourceType,

	Attributes: dal.AttributeSet{
		&dal.Attribute{
			Ident: "ID",
			Type:  &dal.TypeID{},
			Store: &dal.CodecAlias{Ident: "id"},
		},

		&dal.Attribute{
			Ident: "Timestamp", Sortable: true,
			Type:  &dal.TypeTimestamp{Timezone: true, Precision: -1},
			Store: &dal.CodecAlias{Ident: "ts"},
		},

		&dal.Attribute{
			Ident: "ResourceID",
			Type:  &dal.TypeID{},
			Store: &dal.CodecAlias{Ident: "rel_resource"},
		},

		&dal.Attribute{
			Ident: "Revision",
			Type:  &dal.TypeNumber{Precision: -1, Scale: -1, Meta: map[string]interface{}{"rdbms:type": "integer"}},
			Store: &dal.CodecAlias{Ident: "revision"},
		},

		&dal.Attribute{
			Ident: "Operation",
			Type:  &dal.TypeText{},
			Store: &dal.CodecAlias{Ident: "operation"},
		},

		&dal.Attribute{
			Ident: "RelUser",
			Type: &dal.TypeRef{HasDefault: true,
				DefaultValue: 0,

				RefAttribute: "id",
				RefModel: &dal.ModelRef{
					ResourceType: "corteza::system:user",
				},
			},
			Store: &dal.CodecAlias{Ident: "rel_user"},
		},

		&dal.Attribute{
			Ident: "Delta",
			Type: &dal.TypeJSON{
				DefaultValue: "{}",
			},
			Store: &dal.CodecAlias{Ident: "delta"},
		},

		&dal.Attribute{
			Ident: "Comment",
			Type:  &dal.TypeText{},
			Store: &dal.CodecAlias{Ident: "comment"},
		},
	},

	Indexes: dal.IndexSet{
		&dal.Index{
			Ident: "PRIMARY",
			Type:  "BTREE",

			Fields: []*dal.IndexField{
				{
					AttributeIdent: "ID",
				},
			},
		},
	},
}

func init() {
	models = append(
		models,
		Attachment,
		Chart,
		City311ActorProfile,
		City311AuditEvent,
		City311IdempotencyRecord,
		City311PublicHistoryItem,
		City311RequestAttachment,
		City311RequestSequence,
		City311ServiceRequest,
		Module,
		ModuleField,
		Namespace,
		Page,
		PageLayout,
		Record,
		RecordRevision,
	)
}
