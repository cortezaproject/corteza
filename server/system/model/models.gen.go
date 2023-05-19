package model

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"github.com/cortezaproject/corteza/server/pkg/dal"
	"github.com/cortezaproject/corteza/server/system/types"
)

var ApigwFilter = &dal.Model{
	Ident:        "apigw_filters",
	ResourceType: types.ApigwFilterResourceType,

	Attributes: dal.AttributeSet{
		&dal.Attribute{
			Ident: "ID",
			Type:  &dal.TypeID{},
			Store: &dal.CodecAlias{Ident: "id"},
		},

		&dal.Attribute{
			Ident: "Route", Sortable: true,
			Type: &dal.TypeRef{
				RefAttribute: "id",
				RefModel: &dal.ModelRef{
					ResourceType: "corteza::system:apigw-route",
				},
			},
			Store: &dal.CodecAlias{Ident: "rel_route"},
		},

		&dal.Attribute{
			Ident: "Weight", Sortable: true,
			Type:  &dal.TypeNumber{Precision: -1, Scale: -1, Meta: map[string]interface{}{"rdbms:type": "integer"}},
			Store: &dal.CodecAlias{Ident: "weight"},
		},

		&dal.Attribute{
			Ident: "Kind", Sortable: true,
			Type:  &dal.TypeText{Length: 64},
			Store: &dal.CodecAlias{Ident: "kind"},
		},

		&dal.Attribute{
			Ident: "Ref",
			Type:  &dal.TypeText{Length: 64},
			Store: &dal.CodecAlias{Ident: "ref"},
		},

		&dal.Attribute{
			Ident: "Enabled", Sortable: true,
			Type:  &dal.TypeBoolean{},
			Store: &dal.CodecAlias{Ident: "enabled"},
		},

		&dal.Attribute{
			Ident: "Params",
			Type: &dal.TypeJSON{
				DefaultValue: "{}",
			},
			Store: &dal.CodecAlias{Ident: "params"},
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

var ApigwRoute = &dal.Model{
	Ident:        "apigw_routes",
	ResourceType: types.ApigwRouteResourceType,

	Attributes: dal.AttributeSet{
		&dal.Attribute{
			Ident: "ID",
			Type:  &dal.TypeID{},
			Store: &dal.CodecAlias{Ident: "id"},
		},

		&dal.Attribute{
			Ident: "Endpoint", Sortable: true,
			Type:  &dal.TypeText{},
			Store: &dal.CodecAlias{Ident: "endpoint"},
		},

		&dal.Attribute{
			Ident: "Method", Sortable: true,
			Type:  &dal.TypeText{},
			Store: &dal.CodecAlias{Ident: "method"},
		},

		&dal.Attribute{
			Ident: "Enabled", Sortable: true,
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
			Ident: "Group", Sortable: true,
			Type: &dal.TypeRef{
				RefAttribute: "id",
				RefModel: &dal.ModelRef{
					ResourceType: "corteza::system:apigw-group",
				},
			},
			Store: &dal.CodecAlias{Ident: "rel_group"},
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

var Application = &dal.Model{
	Ident:        "applications",
	ResourceType: types.ApplicationResourceType,

	Attributes: dal.AttributeSet{
		&dal.Attribute{
			Ident: "ID",
			Type:  &dal.TypeID{},
			Store: &dal.CodecAlias{Ident: "id"},
		},

		&dal.Attribute{
			Ident: "Name", Sortable: true,
			Type:  &dal.TypeText{},
			Store: &dal.CodecAlias{Ident: "name"},
		},

		&dal.Attribute{
			Ident: "Enabled", Sortable: true,
			Type: &dal.TypeBoolean{HasDefault: true,
				DefaultValue: true,
			},
			Store: &dal.CodecAlias{Ident: "enabled"},
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
			Ident: "Unify",
			Type: &dal.TypeJSON{
				DefaultValue: "{}",
			},
			Store: &dal.CodecAlias{Ident: "unify"},
		},

		&dal.Attribute{
			Ident: "OwnerID",
			Type: &dal.TypeRef{HasDefault: true,
				DefaultValue: 0,

				RefAttribute: "id",
				RefModel: &dal.ModelRef{
					ResourceType: "corteza::system:user",
				},
			},
			Store: &dal.CodecAlias{Ident: "rel_owner"},
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
	},
}

var Attachment = &dal.Model{
	Ident:        "attachments",
	ResourceType: types.AttachmentResourceType,

	Attributes: dal.AttributeSet{
		&dal.Attribute{
			Ident: "ID",
			Type:  &dal.TypeID{},
			Store: &dal.CodecAlias{Ident: "id"},
		},

		&dal.Attribute{
			Ident: "OwnerID",
			Type: &dal.TypeRef{HasDefault: true,
				DefaultValue: 0,

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

var AuthClient = &dal.Model{
	Ident:        "auth_clients",
	ResourceType: types.AuthClientResourceType,

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
			Ident: "Meta",
			Type: &dal.TypeJSON{
				DefaultValue: "{}",
			},
			Store: &dal.CodecAlias{Ident: "meta"},
		},

		&dal.Attribute{
			Ident: "Secret",
			Type:  &dal.TypeText{Length: 64},
			Store: &dal.CodecAlias{Ident: "secret"},
		},

		&dal.Attribute{
			Ident: "Scope",
			Type:  &dal.TypeText{Length: 512},
			Store: &dal.CodecAlias{Ident: "scope"},
		},

		&dal.Attribute{
			Ident: "ValidGrant",
			Type:  &dal.TypeText{Length: 32},
			Store: &dal.CodecAlias{Ident: "valid_grant"},
		},

		&dal.Attribute{
			Ident: "RedirectURI",
			Type:  &dal.TypeText{},
			Store: &dal.CodecAlias{Ident: "redirect_uri"},
		},

		&dal.Attribute{
			Ident: "Enabled", Sortable: true,
			Type: &dal.TypeBoolean{HasDefault: true,
				DefaultValue: false,
			},
			Store: &dal.CodecAlias{Ident: "enabled"},
		},

		&dal.Attribute{
			Ident: "Trusted", Sortable: true,
			Type: &dal.TypeBoolean{HasDefault: true,
				DefaultValue: false,
			},
			Store: &dal.CodecAlias{Ident: "trusted"},
		},

		&dal.Attribute{
			Ident: "ValidFrom", Sortable: true,
			Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true, Precision: -1},
			Store: &dal.CodecAlias{Ident: "valid_from"},
		},

		&dal.Attribute{
			Ident: "ExpiresAt", Sortable: true,
			Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true, Precision: -1},
			Store: &dal.CodecAlias{Ident: "expires_at"},
		},

		&dal.Attribute{
			Ident: "Security",
			Type: &dal.TypeJSON{
				DefaultValue: "{}",
			},
			Store: &dal.CodecAlias{Ident: "security"},
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

var AuthConfirmedClient = &dal.Model{
	Ident:        "auth_confirmed_clients",
	ResourceType: types.AuthConfirmedClientResourceType,

	Attributes: dal.AttributeSet{
		&dal.Attribute{
			Ident: "UserID",
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
			Ident: "ClientID",
			Type: &dal.TypeRef{HasDefault: true,
				DefaultValue: 0,

				RefAttribute: "id",
				RefModel: &dal.ModelRef{
					ResourceType: "corteza::system:auth-client",
				},
			},
			Store: &dal.CodecAlias{Ident: "rel_client"},
		},

		&dal.Attribute{
			Ident: "ConfirmedAt", Sortable: true,
			Type: &dal.TypeTimestamp{
				DefaultCurrentTimestamp: true, Timezone: true, Precision: -1,
			},
			Store: &dal.CodecAlias{Ident: "confirmed_at"},
		},
	},

	Indexes: dal.IndexSet{
		&dal.Index{
			Ident: "PRIMARY",
			Type:  "BTREE",

			Fields: []*dal.IndexField{
				{
					AttributeIdent: "UserID",
				},

				{
					AttributeIdent: "ClientID",
				},
			},
		},
	},
}

var AuthOa2token = &dal.Model{
	Ident:        "auth_oa2tokens",
	ResourceType: types.AuthOa2tokenResourceType,

	Attributes: dal.AttributeSet{
		&dal.Attribute{
			Ident: "ID",
			Type:  &dal.TypeID{},
			Store: &dal.CodecAlias{Ident: "id"},
		},

		&dal.Attribute{
			Ident: "Code",
			Type:  &dal.TypeText{Length: 48},
			Store: &dal.CodecAlias{Ident: "code"},
		},

		&dal.Attribute{
			Ident: "Access",
			Type:  &dal.TypeText{Length: 2048},
			Store: &dal.CodecAlias{Ident: "access"},
		},

		&dal.Attribute{
			Ident: "Refresh",
			Type:  &dal.TypeText{Length: 48},
			Store: &dal.CodecAlias{Ident: "refresh"},
		},

		&dal.Attribute{
			Ident: "Data",
			Type: &dal.TypeJSON{
				DefaultValue: "{}",
			},
			Store: &dal.CodecAlias{Ident: "data"},
		},

		&dal.Attribute{
			Ident: "RemoteAddr",
			Type:  &dal.TypeText{Length: 64},
			Store: &dal.CodecAlias{Ident: "remote_addr"},
		},

		&dal.Attribute{
			Ident: "UserAgent",
			Type:  &dal.TypeText{},
			Store: &dal.CodecAlias{Ident: "user_agent"},
		},

		&dal.Attribute{
			Ident: "ClientID",
			Type: &dal.TypeRef{HasDefault: true,
				DefaultValue: 0,

				RefAttribute: "id",
				RefModel: &dal.ModelRef{
					ResourceType: "corteza::system:auth-client",
				},
			},
			Store: &dal.CodecAlias{Ident: "rel_client"},
		},

		&dal.Attribute{
			Ident: "UserID",
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
			Ident: "auth_oa2tokens_clientId",
			Type:  "BTREE",

			Fields: []*dal.IndexField{
				{
					AttributeIdent: "ClientID",
				},
			},
		},

		&dal.Index{
			Ident: "auth_oa2tokens_code",
			Type:  "BTREE",

			Fields: []*dal.IndexField{
				{
					AttributeIdent: "Code",
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
			Ident: "auth_oa2tokens_refresh",
			Type:  "BTREE",

			Fields: []*dal.IndexField{
				{
					AttributeIdent: "Refresh",
				},
			},
		},
	},
}

var AuthSession = &dal.Model{
	Ident:        "auth_sessions",
	ResourceType: types.AuthSessionResourceType,

	Attributes: dal.AttributeSet{
		&dal.Attribute{
			Ident: "ID",
			Type:  &dal.TypeText{Length: 64},
			Store: &dal.CodecAlias{Ident: "id"},
		},

		&dal.Attribute{
			Ident: "Data",
			Type:  &dal.TypeBlob{},
			Store: &dal.CodecAlias{Ident: "data"},
		},

		&dal.Attribute{
			Ident: "UserID",
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
			Ident: "RemoteAddr",
			Type:  &dal.TypeText{},
			Store: &dal.CodecAlias{Ident: "remote_addr"},
		},

		&dal.Attribute{
			Ident: "UserAgent",
			Type:  &dal.TypeText{},
			Store: &dal.CodecAlias{Ident: "user_agent"},
		},

		&dal.Attribute{
			Ident: "ExpiresAt", Sortable: true,
			Type:  &dal.TypeTimestamp{Timezone: true, Precision: -1},
			Store: &dal.CodecAlias{Ident: "expires_at"},
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
			Ident: "auth_sessions_expiresAt",
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
	},
}

var Credential = &dal.Model{
	Ident:        "credentials",
	ResourceType: types.CredentialResourceType,

	Attributes: dal.AttributeSet{
		&dal.Attribute{
			Ident: "ID",
			Type:  &dal.TypeID{},
			Store: &dal.CodecAlias{Ident: "id"},
		},

		&dal.Attribute{
			Ident: "OwnerID",
			Type: &dal.TypeRef{HasDefault: true,
				DefaultValue: 0,

				RefAttribute: "id",
				RefModel: &dal.ModelRef{
					ResourceType: "corteza::system:user",
				},
			},
			Store: &dal.CodecAlias{Ident: "rel_owner"},
		},

		&dal.Attribute{
			Ident: "Label",
			Type:  &dal.TypeText{},
			Store: &dal.CodecAlias{Ident: "label"},
		},

		&dal.Attribute{
			Ident: "Kind",
			Type:  &dal.TypeText{Length: 128},
			Store: &dal.CodecAlias{Ident: "kind"},
		},

		&dal.Attribute{
			Ident: "Credentials",
			Type:  &dal.TypeText{},
			Store: &dal.CodecAlias{Ident: "credentials"},
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

		&dal.Attribute{
			Ident: "LastUsedAt", Sortable: true,
			Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true, Precision: -1},
			Store: &dal.CodecAlias{Ident: "last_used_at"},
		},

		&dal.Attribute{
			Ident: "ExpiresAt", Sortable: true,
			Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true, Precision: -1},
			Store: &dal.CodecAlias{Ident: "expires_at"},
		},
	},

	Indexes: dal.IndexSet{
		&dal.Index{
			Ident: "credentials_ownerKind",
			Type:  "BTREE",

			Predicate: "deleted_at IS NULL",
			Fields: []*dal.IndexField{
				{
					AttributeIdent: "OwnerID",
				},

				{
					AttributeIdent: "Kind",
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

var DalConnection = &dal.Model{
	Ident:        "dal_connections",
	ResourceType: types.DalConnectionResourceType,

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
			Ident: "Type", Sortable: true,
			Type:  &dal.TypeText{},
			Store: &dal.CodecAlias{Ident: "type"},
		},

		&dal.Attribute{
			Ident: "Config",
			Type: &dal.TypeJSON{
				DefaultValue: "{}",
			},
			Store: &dal.CodecAlias{Ident: "config"},
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

var DalSchemaAlteration = &dal.Model{
	Ident:        "dal_schema_alterations",
	ResourceType: types.DalSchemaAlterationResourceType,

	Attributes: dal.AttributeSet{
		&dal.Attribute{
			Ident: "ID",
			Type:  &dal.TypeID{},
			Store: &dal.CodecAlias{Ident: "id"},
		},

		&dal.Attribute{
			Ident: "BatchID",
			Type:  &dal.TypeID{},
			Store: &dal.CodecAlias{Ident: "batchID"},
		},

		&dal.Attribute{
			Ident: "DependsOn",
			Type: &dal.TypeRef{
				RefAttribute: "id",
				RefModel: &dal.ModelRef{
					ResourceType: "corteza::system:dal-schema-alteration",
				},
			},
			Store: &dal.CodecAlias{Ident: "dependsOn"},
		},

		&dal.Attribute{
			Ident: "Kind",
			Type:  &dal.TypeText{Length: 256},
			Store: &dal.CodecAlias{Ident: "kind"},
		},

		&dal.Attribute{
			Ident: "Params",
			Type: &dal.TypeJSON{
				DefaultValue: "{}",
			},
			Store: &dal.CodecAlias{Ident: "params"},
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
			Ident: "CompletedAt", Sortable: true,
			Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true, Precision: -1},
			Store: &dal.CodecAlias{Ident: "completed_at"},
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

		&dal.Attribute{
			Ident: "CompletedBy",
			Type: &dal.TypeRef{HasDefault: true,
				DefaultValue: 0,

				RefAttribute: "id",
				RefModel: &dal.ModelRef{
					ResourceType: "corteza::system:user",
				},
			},
			Store: &dal.CodecAlias{Ident: "completed_by"},
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
			Ident:  "dal_schema_alterations_uniqueAlteration",
			Type:   "BTREE",
			Unique: true,

			Fields: []*dal.IndexField{
				{
					AttributeIdent: "ID",
					Modifiers:      []dal.IndexFieldModifier{"LOWERCASE"},
				},

				{
					AttributeIdent: "BatchID",
					Modifiers:      []dal.IndexFieldModifier{"LOWERCASE"},
				},
			},
		},
	},
}

var DalSensitivityLevel = &dal.Model{
	Ident:        "dal_sensitivity_levels",
	ResourceType: types.DalSensitivityLevelResourceType,

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
			Ident: "Level", Sortable: true,
			Type:  &dal.TypeNumber{Precision: -1, Scale: -1, Meta: map[string]interface{}{"rdbms:type": "integer"}},
			Store: &dal.CodecAlias{Ident: "level"},
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

var DataPrivacyRequest = &dal.Model{
	Ident:        "data_privacy_requests",
	ResourceType: types.DataPrivacyRequestResourceType,

	Attributes: dal.AttributeSet{
		&dal.Attribute{
			Ident: "ID",
			Type:  &dal.TypeID{},
			Store: &dal.CodecAlias{Ident: "id"},
		},

		&dal.Attribute{
			Ident: "Kind", Sortable: true,
			Type:  &dal.TypeText{},
			Store: &dal.CodecAlias{Ident: "kind"},
		},

		&dal.Attribute{
			Ident: "Status", Sortable: true,
			Type:  &dal.TypeText{Length: 64},
			Store: &dal.CodecAlias{Ident: "status"},
		},

		&dal.Attribute{
			Ident: "Payload",
			Type:  &dal.TypeJSON{},
			Store: &dal.CodecAlias{Ident: "payload"},
		},

		&dal.Attribute{
			Ident: "RequestedAt", Sortable: true,
			Type:  &dal.TypeTimestamp{Timezone: true, Precision: -1},
			Store: &dal.CodecAlias{Ident: "requested_at"},
		},

		&dal.Attribute{
			Ident: "RequestedBy",
			Type: &dal.TypeRef{
				RefAttribute: "id",
				RefModel: &dal.ModelRef{
					ResourceType: "corteza::system:user",
				},
			},
			Store: &dal.CodecAlias{Ident: "requested_by"},
		},

		&dal.Attribute{
			Ident: "CompletedAt", Sortable: true,
			Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true, Precision: -1},
			Store: &dal.CodecAlias{Ident: "completed_at"},
		},

		&dal.Attribute{
			Ident: "CompletedBy",
			Type: &dal.TypeRef{
				RefAttribute: "id",
				RefModel: &dal.ModelRef{
					ResourceType: "corteza::system:user",
				},
			},
			Store: &dal.CodecAlias{Ident: "completed_by"},
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

var DataPrivacyRequestComment = &dal.Model{
	Ident:        "data_privacy_request_comments",
	ResourceType: types.DataPrivacyRequestCommentResourceType,

	Attributes: dal.AttributeSet{
		&dal.Attribute{
			Ident: "ID",
			Type:  &dal.TypeID{},
			Store: &dal.CodecAlias{Ident: "id"},
		},

		&dal.Attribute{
			Ident: "RequestID",
			Type: &dal.TypeRef{
				RefAttribute: "id",
				RefModel: &dal.ModelRef{
					ResourceType: "corteza::system:user",
				},
			},
			Store: &dal.CodecAlias{Ident: "rel_request"},
		},

		&dal.Attribute{
			Ident: "Comment",
			Type:  &dal.TypeText{},
			Store: &dal.CodecAlias{Ident: "comment"},
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

var Queue = &dal.Model{
	Ident:        "queue_settings",
	ResourceType: types.QueueResourceType,

	Attributes: dal.AttributeSet{
		&dal.Attribute{
			Ident: "ID",
			Type:  &dal.TypeID{},
			Store: &dal.CodecAlias{Ident: "id"},
		},

		&dal.Attribute{
			Ident: "Consumer", Sortable: true,
			Type:  &dal.TypeText{},
			Store: &dal.CodecAlias{Ident: "consumer"},
		},

		&dal.Attribute{
			Ident: "Queue", Sortable: true,
			Type:  &dal.TypeText{},
			Store: &dal.CodecAlias{Ident: "queue"},
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

var QueueMessage = &dal.Model{
	Ident:        "queue_messages",
	ResourceType: types.QueueMessageResourceType,

	Attributes: dal.AttributeSet{
		&dal.Attribute{
			Ident: "ID",
			Type:  &dal.TypeID{},
			Store: &dal.CodecAlias{Ident: "id"},
		},

		&dal.Attribute{
			Ident: "Queue", Sortable: true,
			Type:  &dal.TypeText{},
			Store: &dal.CodecAlias{Ident: "queue"},
		},

		&dal.Attribute{
			Ident: "Payload",
			Type:  &dal.TypeBlob{},
			Store: &dal.CodecAlias{Ident: "payload"},
		},

		&dal.Attribute{
			Ident: "Created", Sortable: true,
			Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true, Precision: -1},
			Store: &dal.CodecAlias{Ident: "created"},
		},

		&dal.Attribute{
			Ident: "Processed", Sortable: true,
			Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true, Precision: -1},
			Store: &dal.CodecAlias{Ident: "processed"},
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

var Reminder = &dal.Model{
	Ident:        "reminders",
	ResourceType: types.ReminderResourceType,

	Attributes: dal.AttributeSet{
		&dal.Attribute{
			Ident: "ID",
			Type:  &dal.TypeID{},
			Store: &dal.CodecAlias{Ident: "id"},
		},

		&dal.Attribute{
			Ident: "Resource", Sortable: true,
			Type:  &dal.TypeText{Length: 512},
			Store: &dal.CodecAlias{Ident: "resource"},
		},

		&dal.Attribute{
			Ident: "Payload",
			Type: &dal.TypeJSON{
				DefaultValue: "{}",
			},
			Store: &dal.CodecAlias{Ident: "payload"},
		},

		&dal.Attribute{
			Ident: "SnoozeCount",
			Type:  &dal.TypeNumber{Precision: -1, Scale: -1, Meta: map[string]interface{}{"rdbms:type": "integer"}},
			Store: &dal.CodecAlias{Ident: "snooze_count"},
		},

		&dal.Attribute{
			Ident: "AssignedTo",
			Type: &dal.TypeRef{HasDefault: true,
				DefaultValue: 0,

				RefAttribute: "id",
				RefModel: &dal.ModelRef{
					ResourceType: "corteza::system:user",
				},
			},
			Store: &dal.CodecAlias{Ident: "assigned_to"},
		},

		&dal.Attribute{
			Ident: "AssignedBy",
			Type: &dal.TypeRef{HasDefault: true,
				DefaultValue: 0,

				RefAttribute: "id",
				RefModel: &dal.ModelRef{
					ResourceType: "corteza::system:user",
				},
			},
			Store: &dal.CodecAlias{Ident: "assigned_by"},
		},

		&dal.Attribute{
			Ident: "AssignedAt", Sortable: true,
			Type:  &dal.TypeTimestamp{Timezone: true, Precision: -1},
			Store: &dal.CodecAlias{Ident: "assigned_at"},
		},

		&dal.Attribute{
			Ident: "DismissedBy",
			Type: &dal.TypeRef{HasDefault: true,
				DefaultValue: 0,

				RefAttribute: "id",
				RefModel: &dal.ModelRef{
					ResourceType: "corteza::system:user",
				},
			},
			Store: &dal.CodecAlias{Ident: "dismissed_by"},
		},

		&dal.Attribute{
			Ident: "DismissedAt", Sortable: true,
			Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true, Precision: -1},
			Store: &dal.CodecAlias{Ident: "dismissed_at"},
		},

		&dal.Attribute{
			Ident: "RemindAt", Sortable: true,
			Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true, Precision: -1},
			Store: &dal.CodecAlias{Ident: "remind_at"},
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
			Ident: "reminders_assignedTo",
			Type:  "BTREE",

			Fields: []*dal.IndexField{
				{
					AttributeIdent: "AssignedTo",
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
			Ident: "reminders_resource",
			Type:  "BTREE",

			Fields: []*dal.IndexField{
				{
					AttributeIdent: "Resource",
				},
			},
		},
	},
}

var Report = &dal.Model{
	Ident:        "reports",
	ResourceType: types.ReportResourceType,

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
			Ident: "Meta",
			Type: &dal.TypeJSON{
				DefaultValue: "{}",
			},
			Store: &dal.CodecAlias{Ident: "meta"},
		},

		&dal.Attribute{
			Ident: "Scenarios",
			Type: &dal.TypeJSON{
				DefaultValue: "{}",
			},
			Store: &dal.CodecAlias{Ident: "scenarios"},
		},

		&dal.Attribute{
			Ident: "Sources",
			Type: &dal.TypeJSON{
				DefaultValue: "{}",
			},
			Store: &dal.CodecAlias{Ident: "sources"},
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

var ResourceTranslation = &dal.Model{
	Ident:        "resource_translations",
	ResourceType: types.ResourceTranslationResourceType,

	Attributes: dal.AttributeSet{
		&dal.Attribute{
			Ident: "ID",
			Type:  &dal.TypeID{},
			Store: &dal.CodecAlias{Ident: "id"},
		},

		&dal.Attribute{
			Ident: "Lang",
			Type:  &dal.TypeText{Length: 32},
			Store: &dal.CodecAlias{Ident: "lang"},
		},

		&dal.Attribute{
			Ident: "Resource",
			Type:  &dal.TypeText{Length: 256},
			Store: &dal.CodecAlias{Ident: "resource"},
		},

		&dal.Attribute{
			Ident: "K",
			Type:  &dal.TypeText{Length: 256},
			Store: &dal.CodecAlias{Ident: "k"},
		},

		&dal.Attribute{
			Ident: "Message",
			Type:  &dal.TypeText{},
			Store: &dal.CodecAlias{Ident: "message"},
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
			Ident: "PRIMARY",
			Type:  "BTREE",

			Fields: []*dal.IndexField{
				{
					AttributeIdent: "ID",
				},
			},
		},

		&dal.Index{
			Ident:  "resource_translations_uniqueTranslation",
			Type:   "BTREE",
			Unique: true,

			Fields: []*dal.IndexField{
				{
					AttributeIdent: "Lang",
					Modifiers:      []dal.IndexFieldModifier{"LOWERCASE"},
				},

				{
					AttributeIdent: "Resource",
					Modifiers:      []dal.IndexFieldModifier{"LOWERCASE"},
				},

				{
					AttributeIdent: "K",
					Modifiers:      []dal.IndexFieldModifier{"LOWERCASE"},
				},
			},
		},
	},
}

var Role = &dal.Model{
	Ident:        "roles",
	ResourceType: types.RoleResourceType,

	Attributes: dal.AttributeSet{
		&dal.Attribute{
			Ident: "ID",
			Type:  &dal.TypeID{},
			Store: &dal.CodecAlias{Ident: "id"},
		},

		&dal.Attribute{
			Ident: "Name", Sortable: true,
			Type:  &dal.TypeText{},
			Store: &dal.CodecAlias{Ident: "name"},
		},

		&dal.Attribute{
			Ident: "Handle",
			Type:  &dal.TypeText{Length: 64},
			Store: &dal.CodecAlias{Ident: "handle"},
		},

		&dal.Attribute{
			Ident: "Meta",
			Type: &dal.TypeJSON{
				DefaultValue: "{}",
			},
			Store: &dal.CodecAlias{Ident: "meta"},
		},

		&dal.Attribute{
			Ident: "ArchivedAt", Sortable: true,
			Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true, Precision: -1},
			Store: &dal.CodecAlias{Ident: "archived_at"},
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
	},
}

var RoleMember = &dal.Model{
	Ident:        "role_members",
	ResourceType: types.RoleMemberResourceType,

	Attributes: dal.AttributeSet{
		&dal.Attribute{
			Ident: "UserID",
			Type: &dal.TypeRef{
				RefAttribute: "id",
				RefModel: &dal.ModelRef{
					ResourceType: "corteza::system:user",
				},
			},
			Store: &dal.CodecAlias{Ident: "rel_user"},
		},

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
	},

	Indexes: dal.IndexSet{
		&dal.Index{
			Ident: "PRIMARY",
			Type:  "BTREE",

			Fields: []*dal.IndexField{
				{
					AttributeIdent: "UserID",
				},

				{
					AttributeIdent: "RoleID",
				},
			},
		},
	},
}

var SettingValue = &dal.Model{
	Ident:        "settings",
	ResourceType: types.SettingValueResourceType,

	Attributes: dal.AttributeSet{
		&dal.Attribute{
			Ident: "OwnedBy",
			Type: &dal.TypeRef{
				RefAttribute: "id",
				RefModel: &dal.ModelRef{
					ResourceType: "corteza::system:user",
				},
			},
			Store: &dal.CodecAlias{Ident: "rel_owner"},
		},

		&dal.Attribute{
			Ident: "Name",
			Type:  &dal.TypeText{Length: 512},
			Store: &dal.CodecAlias{Ident: "name"},
		},

		&dal.Attribute{
			Ident: "Value",
			Type:  &dal.TypeJSON{},
			Store: &dal.CodecAlias{Ident: "value"},
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
			Ident: "UpdatedAt", Sortable: true,
			Type:  &dal.TypeTimestamp{Timezone: true, Precision: -1},
			Store: &dal.CodecAlias{Ident: "updated_at"},
		},
	},

	Indexes: dal.IndexSet{
		&dal.Index{
			Ident: "PRIMARY",
			Type:  "BTREE",

			Fields: []*dal.IndexField{
				{
					AttributeIdent: "OwnedBy",
				},

				{
					AttributeIdent: "Name",
				},
			},
		},
	},
}

var Template = &dal.Model{
	Ident:        "templates",
	ResourceType: types.TemplateResourceType,

	Attributes: dal.AttributeSet{
		&dal.Attribute{
			Ident: "ID",
			Type:  &dal.TypeID{},
			Store: &dal.CodecAlias{Ident: "id"},
		},

		&dal.Attribute{
			Ident: "OwnerID",
			Type: &dal.TypeRef{HasDefault: true,
				DefaultValue: 0,

				RefAttribute: "id",
				RefModel: &dal.ModelRef{
					ResourceType: "corteza::system:user",
				},
			},
			Store: &dal.CodecAlias{Ident: "rel_owner"},
		},

		&dal.Attribute{
			Ident: "Handle",
			Type:  &dal.TypeText{Length: 64},
			Store: &dal.CodecAlias{Ident: "handle"},
		},

		&dal.Attribute{
			Ident: "Language", Sortable: true,
			Type:  &dal.TypeText{Length: 32},
			Store: &dal.CodecAlias{Ident: "language"},
		},

		&dal.Attribute{
			Ident: "Type", Sortable: true,
			Type:  &dal.TypeText{},
			Store: &dal.CodecAlias{Ident: "type"},
		},

		&dal.Attribute{
			Ident: "Partial",
			Type:  &dal.TypeBoolean{},
			Store: &dal.CodecAlias{Ident: "partial"},
		},

		&dal.Attribute{
			Ident: "Meta",
			Type: &dal.TypeJSON{
				DefaultValue: "{}",
			},
			Store: &dal.CodecAlias{Ident: "meta"},
		},

		&dal.Attribute{
			Ident: "Template", Sortable: true,
			Type:  &dal.TypeText{},
			Store: &dal.CodecAlias{Ident: "template"},
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
			Ident: "LastUsedAt", Sortable: true,
			Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true, Precision: -1},
			Store: &dal.CodecAlias{Ident: "last_used_at"},
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
			Ident:     "templates_uniqueLanguageHandle",
			Type:      "BTREE",
			Unique:    true,
			Predicate: "handle != '' AND deleted_at IS NULL",
			Fields: []*dal.IndexField{
				{
					AttributeIdent: "Language",
				},

				{
					AttributeIdent: "Handle",
				},
			},
		},
	},
}

var User = &dal.Model{
	Ident:        "users",
	ResourceType: types.UserResourceType,

	Attributes: dal.AttributeSet{
		&dal.Attribute{
			Ident: "ID",
			Type:  &dal.TypeID{},
			Store: &dal.CodecAlias{Ident: "id"},
		},

		&dal.Attribute{
			Ident: "Email", Sortable: true,
			Type:  &dal.TypeText{Length: 254},
			Store: &dal.CodecAlias{Ident: "email"},
		},

		&dal.Attribute{
			Ident: "EmailConfirmed",
			Type:  &dal.TypeBoolean{},
			Store: &dal.CodecAlias{Ident: "email_confirmed"},
		},

		&dal.Attribute{
			Ident: "Username", Sortable: true,
			Type:  &dal.TypeText{},
			Store: &dal.CodecAlias{Ident: "username"},
		},

		&dal.Attribute{
			Ident: "Name", Sortable: true,
			Type:  &dal.TypeText{},
			Store: &dal.CodecAlias{Ident: "name"},
		},

		&dal.Attribute{
			Ident: "Handle",
			Type:  &dal.TypeText{Length: 64},
			Store: &dal.CodecAlias{Ident: "handle"},
		},

		&dal.Attribute{
			Ident: "Kind", Sortable: true,
			Type:  &dal.TypeText{Length: 8},
			Store: &dal.CodecAlias{Ident: "kind"},
		},

		&dal.Attribute{
			Ident: "Meta",
			Type: &dal.TypeJSON{
				DefaultValue: "{}",
			},
			Store: &dal.CodecAlias{Ident: "meta"},
		},

		&dal.Attribute{
			Ident: "SuspendedAt", Sortable: true,
			Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true, Precision: -1},
			Store: &dal.CodecAlias{Ident: "suspended_at"},
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
			Ident:     "users_uniqueEmail",
			Type:      "BTREE",
			Unique:    true,
			Predicate: "email != '' AND deleted_at IS NULL",
			Fields: []*dal.IndexField{
				{
					AttributeIdent: "Email",
					Modifiers:      []dal.IndexFieldModifier{"LOWERCASE"},
				},
			},
		},

		&dal.Index{
			Ident:     "users_uniqueHandle",
			Type:      "BTREE",
			Unique:    true,
			Predicate: "handle != '' AND deleted_at IS NULL",
			Fields: []*dal.IndexField{
				{
					AttributeIdent: "Handle",
					Modifiers:      []dal.IndexFieldModifier{"LOWERCASE"},
				},
			},
		},

		&dal.Index{
			Ident:     "users_uniqueUsername",
			Type:      "BTREE",
			Unique:    true,
			Predicate: "username != '' AND deleted_at IS NULL",
			Fields: []*dal.IndexField{
				{
					AttributeIdent: "Username",
					Modifiers:      []dal.IndexFieldModifier{"LOWERCASE"},
				},
			},
		},
	},
}

func init() {
	models = append(
		models,
		ApigwFilter,
		ApigwRoute,
		Application,
		Attachment,
		AuthClient,
		AuthConfirmedClient,
		AuthOa2token,
		AuthSession,
		Credential,
		DalConnection,
		DalSchemaAlteration,
		DalSensitivityLevel,
		DataPrivacyRequest,
		DataPrivacyRequestComment,
		Queue,
		QueueMessage,
		Reminder,
		Report,
		ResourceTranslation,
		Role,
		RoleMember,
		SettingValue,
		Template,
		User,
	)
}
