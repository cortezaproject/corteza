package model

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/dal"
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
			Ident: "namespace",
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
			Ident: "namespace",
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
			Ident: "uniqueHandle",
			Type:  "BTREE",

			Predicate: "handle != '' AND deleted_at IS NULL",
			Fields: []*dal.IndexField{
				{
					AttributeIdent: "Handle",
					Modifiers:      []dal.IndexFieldModifier{"LOWERCASE"},
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
			Ident: "namespace",
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
			Ident: "uniqueHandle",
			Type:  "BTREE",

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
			Type:  &dal.TypeNumber{Precision: -1, Scale: -1},
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
			Ident: "module",
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
			Ident: "uniqueName",
			Type:  "BTREE",

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
			Ident: "uniqueHandle",
			Type:  "BTREE",

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
				Precision:    -1, Scale: -1,
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
			Ident: "module",
			Type:  "BTREE",

			Fields: []*dal.IndexField{
				{
					AttributeIdent: "ModuleID",
				},
			},
		},

		&dal.Index{
			Ident: "namespace",
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
			Ident: "selfId",
			Type:  "BTREE",

			Fields: []*dal.IndexField{
				{
					AttributeIdent: "SelfID",
				},
			},
		},

		&dal.Index{
			Ident: "uniqueHandle",
			Type:  "BTREE",

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
			Type:  &dal.TypeNumber{Precision: -1, Scale: -1},
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
			Ident: "idxComposeRecordBase",
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
			Type:  &dal.TypeNumber{Precision: -1, Scale: -1},
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
		Module,
		ModuleField,
		Namespace,
		Page,
		Record,
		RecordRevision,
	)
}
