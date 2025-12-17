package model

import (
	"github.com/cortezaproject/corteza/server/pkg/dal"
	"github.com/cortezaproject/corteza/server/pkg/revisions"
)

var Revision = &dal.Model{
	Ident:        "system_revisions",
	ResourceType: revisions.RevisionResourceType,

	Attributes: dal.AttributeSet{
		&dal.Attribute{
			Ident:      "ID",
			PrimaryKey: true,
			Type:       &dal.TypeID{},
			Store:      &dal.CodecAlias{Ident: "id"},
		},
		&dal.Attribute{
			Ident:    "Timestamp",
			Sortable: true,
			Type:     &dal.TypeTimestamp{Timezone: true, Precision: -1},
			Store:    &dal.CodecAlias{Ident: "ts"},
		},
		&dal.Attribute{
			Ident: "ResourceID",
			Type:  &dal.TypeID{},
			Store: &dal.CodecAlias{Ident: "rel_resource"},
		},
		&dal.Attribute{
			Ident: "ResourceType",
			Type:  &dal.TypeText{},
			Store: &dal.CodecAlias{Ident: "resource_type"},
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
			Ident: "Status",
			Type:  &dal.TypeText{},
			Store: &dal.CodecAlias{Ident: "status"},
		},
		&dal.Attribute{
			Ident: "UserID",
			Type:  &dal.TypeID{},
			Store: &dal.CodecAlias{Ident: "rel_user"},
		},
		&dal.Attribute{
			Ident: "Changes",
			Type:  &dal.TypeJSON{DefaultValue: "[]"},
			Store: &dal.CodecAlias{Ident: "delta"},
		},
		&dal.Attribute{
			Ident: "Comment",
			Type:  &dal.TypeText{},
			Store: &dal.CodecAlias{Ident: "comment"},
		},
		&dal.Attribute{
			Ident:    "DeletedAt",
			Sortable: true,
			Type:     &dal.TypeTimestamp{Nullable: true, Timezone: true, Precision: -1},
			Store:    &dal.CodecAlias{Ident: "deleted_at"},
		},
		&dal.Attribute{
			Ident: "DeletedBy",
			Type:  &dal.TypeID{},
			Store: &dal.CodecAlias{Ident: "deleted_by"},
		},
	},

	Indexes: dal.IndexSet{
		&dal.Index{
			Ident: "PRIMARY",
			Type:  "BTREE",
			Fields: []*dal.IndexField{
				{AttributeIdent: "ID"},
			},
		},
	},
}

func ModelRef() dal.ModelRef {
	return dal.ModelRef{ResourceType: revisions.RevisionResourceType}
}
