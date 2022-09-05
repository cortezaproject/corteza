package model

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"context"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/dal"
)

type (
	modelReplacer interface {
		ReplaceModel(ctx context.Context, model *dal.Model) (err error)
	}
)

var (
	Attachment = &dal.Model{
		Ident:        "compose_attachment",
		ResourceType: types.AttachmentResourceType,

		Attributes: dal.AttributeSet{
			&dal.Attribute{
				Ident:      "ID",
				PrimaryKey: true,
				Type:       &dal.TypeID{},
				Store:      &dal.CodecAlias{Ident: "id"},
			},

			&dal.Attribute{
				Ident:    "OwnerID",
				Sortable: true,
				Type: &dal.TypeRef{

					RefAttribute: "id",
					RefModel: &dal.ModelRef{
						ResourceType: "corteza::system:user",
					},
				},
				Store: &dal.CodecAlias{Ident: "rel_owner"},
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
				Ident:    "Kind",
				Sortable: true,
				Type:     &dal.TypeText{},
				Store:    &dal.CodecAlias{Ident: "kind"},
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
				Ident:    "Name",
				Sortable: true,
				Type:     &dal.TypeText{},
				Store:    &dal.CodecAlias{Ident: "name"},
			},

			&dal.Attribute{
				Ident: "Meta",

				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "meta"},
			},

			&dal.Attribute{
				Ident:    "CreatedAt",
				Sortable: true,
				Type:     &dal.TypeTimestamp{},
				Store:    &dal.CodecAlias{Ident: "created_at"},
			},

			&dal.Attribute{
				Ident:    "UpdatedAt",
				Sortable: true,
				Type: &dal.TypeTimestamp{
					Nullable: true},
				Store: &dal.CodecAlias{Ident: "updated_at"},
			},

			&dal.Attribute{
				Ident:    "DeletedAt",
				Sortable: true,
				Type: &dal.TypeTimestamp{
					Nullable: true},
				Store: &dal.CodecAlias{Ident: "deleted_at"},
			},
		},
	}

	Chart = &dal.Model{
		Ident:        "compose_chart",
		ResourceType: types.ChartResourceType,

		Attributes: dal.AttributeSet{
			&dal.Attribute{
				Ident:      "ID",
				PrimaryKey: true,
				Type:       &dal.TypeID{},
				Store:      &dal.CodecAlias{Ident: "id"},
			},

			&dal.Attribute{
				Ident: "Handle",

				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "handle"},
			},

			&dal.Attribute{
				Ident:    "Name",
				Sortable: true,
				Type:     &dal.TypeText{},
				Store:    &dal.CodecAlias{Ident: "name"},
			},

			&dal.Attribute{
				Ident: "Config",

				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "config"},
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
				Ident:    "CreatedAt",
				Sortable: true,
				Type:     &dal.TypeTimestamp{},
				Store:    &dal.CodecAlias{Ident: "created_at"},
			},

			&dal.Attribute{
				Ident:    "UpdatedAt",
				Sortable: true,
				Type: &dal.TypeTimestamp{
					Nullable: true},
				Store: &dal.CodecAlias{Ident: "updated_at"},
			},

			&dal.Attribute{
				Ident:    "DeletedAt",
				Sortable: true,
				Type: &dal.TypeTimestamp{
					Nullable: true},
				Store: &dal.CodecAlias{Ident: "deleted_at"},
			},
		},
	}

	Module = &dal.Model{
		Ident:        "compose_module",
		ResourceType: types.ModuleResourceType,

		Attributes: dal.AttributeSet{
			&dal.Attribute{
				Ident:      "ID",
				PrimaryKey: true,
				Type:       &dal.TypeID{},
				Store:      &dal.CodecAlias{Ident: "id"},
			},

			&dal.Attribute{
				Ident: "Handle",

				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "handle"},
			},

			&dal.Attribute{
				Ident: "Meta",

				Type:  &dal.TypeJSON{},
				Store: &dal.CodecAlias{Ident: "meta"},
			},

			&dal.Attribute{
				Ident: "Config",

				Type:  &dal.TypeJSON{},
				Store: &dal.CodecAlias{Ident: "config"},
			},

			&dal.Attribute{
				Ident: "Fields",

				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "fields"},
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
				Ident:    "Name",
				Sortable: true,
				Type:     &dal.TypeText{},
				Store:    &dal.CodecAlias{Ident: "name"},
			},

			&dal.Attribute{
				Ident:    "CreatedAt",
				Sortable: true,
				Type:     &dal.TypeTimestamp{},
				Store:    &dal.CodecAlias{Ident: "created_at"},
			},

			&dal.Attribute{
				Ident:    "UpdatedAt",
				Sortable: true,
				Type: &dal.TypeTimestamp{
					Nullable: true},
				Store: &dal.CodecAlias{Ident: "updated_at"},
			},

			&dal.Attribute{
				Ident:    "DeletedAt",
				Sortable: true,
				Type: &dal.TypeTimestamp{
					Nullable: true},
				Store: &dal.CodecAlias{Ident: "deleted_at"},
			},
		},
	}

	ModuleField = &dal.Model{
		Ident:        "compose_module_field",
		ResourceType: types.ModuleFieldResourceType,

		Attributes: dal.AttributeSet{
			&dal.Attribute{
				Ident:      "ID",
				PrimaryKey: true,
				Type:       &dal.TypeID{},
				Store:      &dal.CodecAlias{Ident: "id"},
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

				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "rel_namespace"},
			},

			&dal.Attribute{
				Ident:    "Place",
				Sortable: true,
				Type:     &dal.TypeText{},
				Store:    &dal.CodecAlias{Ident: "place"},
			},

			&dal.Attribute{
				Ident:    "Kind",
				Sortable: true,
				Type:     &dal.TypeText{},
				Store:    &dal.CodecAlias{Ident: "kind"},
			},

			&dal.Attribute{
				Ident:    "Name",
				Sortable: true,
				Type:     &dal.TypeText{},
				Store:    &dal.CodecAlias{Ident: "name"},
			},

			&dal.Attribute{
				Ident:    "Label",
				Sortable: true,
				Type:     &dal.TypeText{},
				Store:    &dal.CodecAlias{Ident: "label"},
			},

			&dal.Attribute{
				Ident: "Options",

				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "options"},
			},

			&dal.Attribute{
				Ident: "Config",

				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "config"},
			},

			&dal.Attribute{
				Ident: "Required",

				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "is_required"},
			},

			&dal.Attribute{
				Ident: "Multi",

				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "is_multi"},
			},

			&dal.Attribute{
				Ident: "DefaultValue",

				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "default_value"},
			},

			&dal.Attribute{
				Ident: "Expressions",

				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "expressions"},
			},

			&dal.Attribute{
				Ident: "CreatedAt",

				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "created_at"},
			},

			&dal.Attribute{
				Ident: "UpdatedAt",

				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "updated_at"},
			},

			&dal.Attribute{
				Ident: "DeletedAt",

				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "deleted_at"},
			},
		},
	}

	Namespace = &dal.Model{
		Ident:        "compose_namespace",
		ResourceType: types.NamespaceResourceType,

		Attributes: dal.AttributeSet{
			&dal.Attribute{
				Ident:      "ID",
				PrimaryKey: true,
				Type:       &dal.TypeID{},
				Store:      &dal.CodecAlias{Ident: "id"},
			},

			&dal.Attribute{
				Ident:    "Slug",
				Sortable: true,
				Type:     &dal.TypeText{},
				Store:    &dal.CodecAlias{Ident: "slug"},
			},

			&dal.Attribute{
				Ident: "Enabled",

				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "enabled"},
			},

			&dal.Attribute{
				Ident: "Meta",

				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "meta"},
			},

			&dal.Attribute{
				Ident:    "Name",
				Sortable: true,
				Type:     &dal.TypeText{},
				Store:    &dal.CodecAlias{Ident: "name"},
			},

			&dal.Attribute{
				Ident:    "CreatedAt",
				Sortable: true,
				Type:     &dal.TypeTimestamp{},
				Store:    &dal.CodecAlias{Ident: "created_at"},
			},

			&dal.Attribute{
				Ident:    "UpdatedAt",
				Sortable: true,
				Type: &dal.TypeTimestamp{
					Nullable: true},
				Store: &dal.CodecAlias{Ident: "updated_at"},
			},

			&dal.Attribute{
				Ident:    "DeletedAt",
				Sortable: true,
				Type: &dal.TypeTimestamp{
					Nullable: true},
				Store: &dal.CodecAlias{Ident: "deleted_at"},
			},
		},
	}

	Page = &dal.Model{
		Ident:        "compose_page",
		ResourceType: types.PageResourceType,

		Attributes: dal.AttributeSet{
			&dal.Attribute{
				Ident:      "ID",
				PrimaryKey: true,
				Type:       &dal.TypeID{},
				Store:      &dal.CodecAlias{Ident: "id"},
			},

			&dal.Attribute{
				Ident:    "SelfID",
				Sortable: true,
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
				Ident: "Handle",

				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "handle"},
			},

			&dal.Attribute{
				Ident: "Config",

				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "config"},
			},

			&dal.Attribute{
				Ident: "Blocks",

				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "blocks"},
			},

			&dal.Attribute{
				Ident: "Children",

				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "children"},
			},

			&dal.Attribute{
				Ident: "Visible",

				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "visible"},
			},

			&dal.Attribute{
				Ident:    "Weight",
				Sortable: true,
				Type:     &dal.TypeText{},
				Store:    &dal.CodecAlias{Ident: "weight"},
			},

			&dal.Attribute{
				Ident:    "Title",
				Sortable: true,
				Type:     &dal.TypeText{},
				Store:    &dal.CodecAlias{Ident: "title"},
			},

			&dal.Attribute{
				Ident: "Description",

				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "description"},
			},

			&dal.Attribute{
				Ident:    "CreatedAt",
				Sortable: true,
				Type:     &dal.TypeTimestamp{},
				Store:    &dal.CodecAlias{Ident: "created_at"},
			},

			&dal.Attribute{
				Ident:    "UpdatedAt",
				Sortable: true,
				Type: &dal.TypeTimestamp{
					Nullable: true},
				Store: &dal.CodecAlias{Ident: "updated_at"},
			},

			&dal.Attribute{
				Ident:    "DeletedAt",
				Sortable: true,
				Type: &dal.TypeTimestamp{
					Nullable: true},
				Store: &dal.CodecAlias{Ident: "deleted_at"},
			},
		},
	}

	Record = &dal.Model{
		Ident:        "records",
		ResourceType: types.RecordResourceType,

		Attributes: dal.AttributeSet{
			&dal.Attribute{
				Ident:      "ID",
				PrimaryKey: true,
				Type:       &dal.TypeID{},
				Store:      &dal.CodecAlias{Ident: "id"},
			},

			&dal.Attribute{
				Ident: "ModuleID",

				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "rel_module"},
			},

			&dal.Attribute{
				Ident: "Module",

				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "module"},
			},

			&dal.Attribute{
				Ident: "Values",

				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "values"},
			},

			&dal.Attribute{
				Ident: "NamespaceID",

				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "rel_namespace"},
			},

			&dal.Attribute{
				Ident:    "CreatedAt",
				Sortable: true,
				Type:     &dal.TypeTimestamp{},
				Store:    &dal.CodecAlias{Ident: "created_at"},
			},

			&dal.Attribute{
				Ident:    "UpdatedAt",
				Sortable: true,
				Type: &dal.TypeTimestamp{
					Nullable: true},
				Store: &dal.CodecAlias{Ident: "updated_at"},
			},

			&dal.Attribute{
				Ident:    "DeletedAt",
				Sortable: true,
				Type: &dal.TypeTimestamp{
					Nullable: true},
				Store: &dal.CodecAlias{Ident: "deleted_at"},
			},

			&dal.Attribute{
				Ident: "OwnedBy",

				Type: &dal.TypeRef{

					RefAttribute: "id",
					RefModel: &dal.ModelRef{
						ResourceType: "corteza::system:user",
					},
				},
				Store: &dal.CodecAlias{Ident: "owned_by"},
			},

			&dal.Attribute{
				Ident: "CreatedBy",

				Type: &dal.TypeRef{

					RefAttribute: "id",
					RefModel: &dal.ModelRef{
						ResourceType: "corteza::system:user",
					},
				},
				Store: &dal.CodecAlias{Ident: "created_by"},
			},

			&dal.Attribute{
				Ident: "UpdatedBy",

				Type: &dal.TypeRef{

					RefAttribute: "id",
					RefModel: &dal.ModelRef{
						ResourceType: "corteza::system:user",
					},
				},
				Store: &dal.CodecAlias{Ident: "updated_by"},
			},

			&dal.Attribute{
				Ident: "DeletedBy",

				Type: &dal.TypeRef{

					RefAttribute: "id",
					RefModel: &dal.ModelRef{
						ResourceType: "corteza::system:user",
					},
				},
				Store: &dal.CodecAlias{Ident: "deleted_by"},
			},
		},
	}
)

func All() dal.ModelSet {
	return dal.ModelSet{
		Attachment,
		Chart,
		Module,
		ModuleField,
		Namespace,
		Page,
		Record,
	}
}

func Register(ctx context.Context, mr modelReplacer) (err error) {
	if err = mr.ReplaceModel(ctx, Attachment); err != nil {
		return
	}

	if err = mr.ReplaceModel(ctx, Chart); err != nil {
		return
	}

	if err = mr.ReplaceModel(ctx, Module); err != nil {
		return
	}

	if err = mr.ReplaceModel(ctx, ModuleField); err != nil {
		return
	}

	if err = mr.ReplaceModel(ctx, Namespace); err != nil {
		return
	}

	if err = mr.ReplaceModel(ctx, Page); err != nil {
		return
	}

	if err = mr.ReplaceModel(ctx, Record); err != nil {
		return
	}

	return
}
