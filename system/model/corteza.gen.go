package model

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	actionlogtype "github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/dal"
	discoverytype "github.com/cortezaproject/corteza-server/pkg/discovery/types"
	flagtype "github.com/cortezaproject/corteza-server/pkg/flag/types"
	labelstype "github.com/cortezaproject/corteza-server/pkg/label/types"
	rbactype "github.com/cortezaproject/corteza-server/pkg/rbac"
)

var Action = &dal.Model{
	Ident:        "actionlog",
	ResourceType: actionlogtype.ActionResourceType,

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
			Ident: "ActorIPAddr",
			Type:  &dal.TypeText{Length: 64},
			Store: &dal.CodecAlias{Ident: "actor_ip_addr"},
		},

		&dal.Attribute{
			Ident: "ActorID",
			Type: &dal.TypeRef{
				RefAttribute: "id",
				RefModel: &dal.ModelRef{
					ResourceType: "corteza::system:user",
				},
			},
			Store: &dal.CodecAlias{Ident: "actor_id"},
		},

		&dal.Attribute{
			Ident: "RequestOrigin",
			Type:  &dal.TypeText{Length: 32},
			Store: &dal.CodecAlias{Ident: "request_origin"},
		},

		&dal.Attribute{
			Ident: "RequestID",
			Type:  &dal.TypeText{Length: 256},
			Store: &dal.CodecAlias{Ident: "request_id"},
		},

		&dal.Attribute{
			Ident: "Resource",
			Type:  &dal.TypeText{Length: 512},
			Store: &dal.CodecAlias{Ident: "resource"},
		},

		&dal.Attribute{
			Ident: "Action",
			Type:  &dal.TypeText{Length: 64},
			Store: &dal.CodecAlias{Ident: "action"},
		},

		&dal.Attribute{
			Ident: "Error",
			Type:  &dal.TypeText{},
			Store: &dal.CodecAlias{Ident: "error"},
		},

		&dal.Attribute{
			Ident: "Severity",
			Type: &dal.TypeNumber{HasDefault: true,
				DefaultValue: 0,
				Precision:    -1, Scale: -1, Meta: map[string]interface{}{"rdbms:type": "integer"},
			},
			Store: &dal.CodecAlias{Ident: "severity"},
		},

		&dal.Attribute{
			Ident: "Description",
			Type:  &dal.TypeText{},
			Store: &dal.CodecAlias{Ident: "description"},
		},

		&dal.Attribute{
			Ident: "Meta",
			Type: &dal.TypeJSON{
				DefaultValue: "{}",
			},
			Store: &dal.CodecAlias{Ident: "meta"},
		},
	},

	Indexes: dal.IndexSet{
		&dal.Index{
			Ident: "action",
			Type:  "BTREE",

			Fields: []*dal.IndexField{
				{
					AttributeIdent: "Action",
				},
			},
		},

		&dal.Index{
			Ident: "actorId",
			Type:  "BTREE",

			Fields: []*dal.IndexField{
				{
					AttributeIdent: "ActorID",
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
			Ident: "relResource",
			Type:  "BTREE",

			Fields: []*dal.IndexField{
				{
					AttributeIdent: "Resource",
				},
			},
		},

		&dal.Index{
			Ident: "ts",
			Type:  "BTREE",

			Fields: []*dal.IndexField{
				{
					AttributeIdent: "Timestamp",
				},
			},
		},
	},
}

var Flag = &dal.Model{
	Ident:        "flags",
	ResourceType: flagtype.FlagResourceType,

	Attributes: dal.AttributeSet{
		&dal.Attribute{
			Ident: "Kind",
			Type:  &dal.TypeText{},
			Store: &dal.CodecAlias{Ident: "kind"},
		},

		&dal.Attribute{
			Ident: "ResourceID",
			Type:  &dal.TypeID{},
			Store: &dal.CodecAlias{Ident: "rel_resource"},
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
			Ident: "Name",
			Type:  &dal.TypeText{},
			Store: &dal.CodecAlias{Ident: "name"},
		},

		&dal.Attribute{
			Ident: "Active",
			Type:  &dal.TypeBoolean{},
			Store: &dal.CodecAlias{Ident: "active"},
		},
	},

	Indexes: dal.IndexSet{
		&dal.Index{
			Ident: "uniqueKindResOwnerName",
			Type:  "BTREE",

			Fields: []*dal.IndexField{
				{
					AttributeIdent: "Kind",
				},

				{
					AttributeIdent: "ResourceID",
				},

				{
					AttributeIdent: "OwnedBy",
				},

				{
					AttributeIdent: "Name",
					Modifiers:      []dal.IndexFieldModifier{"LOWERCASE"},
				},
			},
		},
	},
}

var Label = &dal.Model{
	Ident:        "labels",
	ResourceType: labelstype.LabelResourceType,

	Attributes: dal.AttributeSet{
		&dal.Attribute{
			Ident: "Kind",
			Type:  &dal.TypeText{Length: 64},
			Store: &dal.CodecAlias{Ident: "kind"},
		},

		&dal.Attribute{
			Ident: "ResourceID",
			Type:  &dal.TypeID{},
			Store: &dal.CodecAlias{Ident: "rel_resource"},
		},

		&dal.Attribute{
			Ident: "Name",
			Type:  &dal.TypeText{Length: 512},
			Store: &dal.CodecAlias{Ident: "name"},
		},

		&dal.Attribute{
			Ident: "Value",
			Type:  &dal.TypeText{},
			Store: &dal.CodecAlias{Ident: "value"},
		},
	},

	Indexes: dal.IndexSet{
		&dal.Index{
			Ident: "uniqueKindResName",
			Type:  "BTREE",

			Fields: []*dal.IndexField{
				{
					AttributeIdent: "Kind",
				},

				{
					AttributeIdent: "ResourceID",
				},

				{
					AttributeIdent: "Name",
					Modifiers:      []dal.IndexFieldModifier{"LOWERCASE"},
				},
			},
		},
	},
}

var ResourceActivity = &dal.Model{
	Ident:        "resource_activity_log",
	ResourceType: discoverytype.ResourceActivityResourceType,

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
			Ident: "ResourceType",
			Type:  &dal.TypeText{},
			Store: &dal.CodecAlias{Ident: "resource_type"},
		},

		&dal.Attribute{
			Ident: "ResourceAction",
			Type:  &dal.TypeText{},
			Store: &dal.CodecAlias{Ident: "resource_action"},
		},

		&dal.Attribute{
			Ident: "ResourceID",
			Type:  &dal.TypeID{},
			Store: &dal.CodecAlias{Ident: "rel_resource"},
		},

		&dal.Attribute{
			Ident: "Meta",
			Type: &dal.TypeJSON{
				DefaultValue: "{}",
			},
			Store: &dal.CodecAlias{Ident: "meta"},
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
			Ident: "resource",
			Type:  "BTREE",

			Fields: []*dal.IndexField{
				{
					AttributeIdent: "ResourceID",
				},
			},
		},
	},
}

var Rule = &dal.Model{
	Ident:        "rbac_rules",
	ResourceType: rbactype.RuleResourceType,

	Attributes: dal.AttributeSet{
		&dal.Attribute{
			Ident: "RoleID",
			Type: &dal.TypeRef{
				RefAttribute: "id",
				RefModel: &dal.ModelRef{
					ResourceType: "corteza::system:role",
				},
			},
			Store: &dal.CodecAlias{Ident: "rel_role"},
		},

		&dal.Attribute{
			Ident: "Resource",
			Type:  &dal.TypeText{Length: 512},
			Store: &dal.CodecAlias{Ident: "resource"},
		},

		&dal.Attribute{
			Ident: "Operation",
			Type:  &dal.TypeText{Length: 50},
			Store: &dal.CodecAlias{Ident: "operation"},
		},

		&dal.Attribute{
			Ident: "Access",
			Type:  &dal.TypeNumber{Precision: -1, Scale: -1, Meta: map[string]interface{}{"rdbms:type": "integer"}},
			Store: &dal.CodecAlias{Ident: "access"},
		},
	},

	Indexes: dal.IndexSet{
		&dal.Index{
			Ident: "PRIMARY",
			Type:  "BTREE",

			Fields: []*dal.IndexField{
				{
					AttributeIdent: "RoleID",
				},

				{
					AttributeIdent: "Resource",
				},

				{
					AttributeIdent: "Operation",
				},
			},
		},
	},
}

func init() {
	models = append(
		models,
		Action,
		Flag,
		Label,
		ResourceActivity,
		Rule,
	)
}
