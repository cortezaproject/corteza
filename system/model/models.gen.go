package model

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/dal"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	modelReplacer interface {
		ReplaceModel(ctx context.Context, model *dal.Model) (err error)
	}
)

var (
	Attachment = &dal.Model{
		Ident:        "attachments",
		ResourceType: types.AttachmentResourceType,

		Attributes: dal.AttributeSet{
			&dal.Attribute{
				Ident: "ID", PrimaryKey: true,
				Type:  &dal.TypeID{},
				Store: &dal.CodecAlias{Ident: "id"},
			},

			&dal.Attribute{
				Ident: "OwnerID", Sortable: true,
				Type:  &dal.TypeText{},
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
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "meta"},
			},

			&dal.Attribute{
				Ident: "CreatedAt", Sortable: true,
				Type: &dal.TypeTimestamp{
					DefaultCurrentTimestamp: true,
					Timezone:                true,
				},
				Store: &dal.CodecAlias{Ident: "created_at"},
			},

			&dal.Attribute{
				Ident: "UpdatedAt", Sortable: true,
				Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true},
				Store: &dal.CodecAlias{Ident: "updated_at"},
			},

			&dal.Attribute{
				Ident: "DeletedAt", Sortable: true,
				Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true},
				Store: &dal.CodecAlias{Ident: "deleted_at"},
			},
		},
	}

	Application = &dal.Model{
		Ident:        "applications",
		ResourceType: types.ApplicationResourceType,

		Attributes: dal.AttributeSet{
			&dal.Attribute{
				Ident: "ID", PrimaryKey: true,
				Type:  &dal.TypeID{},
				Store: &dal.CodecAlias{Ident: "id"},
			},

			&dal.Attribute{
				Ident: "Name", Sortable: true,
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "name"},
			},

			&dal.Attribute{
				Ident: "OwnerID", Sortable: true,
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "rel_owner"},
			},

			&dal.Attribute{
				Ident: "Enabled", Sortable: true,
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "enabled"},
			},

			&dal.Attribute{
				Ident: "Weight", Sortable: true,
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "weight"},
			},

			&dal.Attribute{
				Ident: "Unify",
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "unify"},
			},

			&dal.Attribute{
				Ident: "CreatedAt", Sortable: true,
				Type: &dal.TypeTimestamp{
					DefaultCurrentTimestamp: true,
					Timezone:                true,
				},
				Store: &dal.CodecAlias{Ident: "created_at"},
			},

			&dal.Attribute{
				Ident: "UpdatedAt", Sortable: true,
				Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true},
				Store: &dal.CodecAlias{Ident: "updated_at"},
			},

			&dal.Attribute{
				Ident: "DeletedAt", Sortable: true,
				Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true},
				Store: &dal.CodecAlias{Ident: "deleted_at"},
			},
		},
	}

	ApigwRoute = &dal.Model{
		Ident:        "apigw_routes",
		ResourceType: types.ApigwRouteResourceType,

		Attributes: dal.AttributeSet{
			&dal.Attribute{
				Ident: "ID", PrimaryKey: true,
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
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "enabled"},
			},

			&dal.Attribute{
				Ident: "Group", Sortable: true,
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "rel_group"},
			},

			&dal.Attribute{
				Ident: "Meta",
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "meta"},
			},

			&dal.Attribute{
				Ident: "CreatedAt", Sortable: true,
				Type: &dal.TypeTimestamp{
					DefaultCurrentTimestamp: true,
					Timezone:                true,
				},
				Store: &dal.CodecAlias{Ident: "created_at"},
			},

			&dal.Attribute{
				Ident: "UpdatedAt", Sortable: true,
				Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true},
				Store: &dal.CodecAlias{Ident: "updated_at"},
			},

			&dal.Attribute{
				Ident: "DeletedAt", Sortable: true,
				Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true},
				Store: &dal.CodecAlias{Ident: "deleted_at"},
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

	ApigwFilter = &dal.Model{
		Ident:        "apigw_filters",
		ResourceType: types.ApigwFilterResourceType,

		Attributes: dal.AttributeSet{
			&dal.Attribute{
				Ident: "ID", PrimaryKey: true,
				Type:  &dal.TypeID{},
				Store: &dal.CodecAlias{Ident: "id"},
			},

			&dal.Attribute{
				Ident: "Route", Sortable: true,
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "rel_route"},
			},

			&dal.Attribute{
				Ident: "Weight", Sortable: true,
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "weight"},
			},

			&dal.Attribute{
				Ident: "Ref",
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "ref"},
			},

			&dal.Attribute{
				Ident: "Kind", Sortable: true,
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "kind"},
			},

			&dal.Attribute{
				Ident: "Enabled", Sortable: true,
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "enabled"},
			},

			&dal.Attribute{
				Ident: "Params",
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "params"},
			},

			&dal.Attribute{
				Ident: "CreatedAt", Sortable: true,
				Type: &dal.TypeTimestamp{
					DefaultCurrentTimestamp: true,
					Timezone:                true,
				},
				Store: &dal.CodecAlias{Ident: "created_at"},
			},

			&dal.Attribute{
				Ident: "UpdatedAt", Sortable: true,
				Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true},
				Store: &dal.CodecAlias{Ident: "updated_at"},
			},

			&dal.Attribute{
				Ident: "DeletedAt", Sortable: true,
				Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true},
				Store: &dal.CodecAlias{Ident: "deleted_at"},
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

	AuthClient = &dal.Model{
		Ident:        "auth_clients",
		ResourceType: types.AuthClientResourceType,

		Attributes: dal.AttributeSet{
			&dal.Attribute{
				Ident: "ID", PrimaryKey: true,
				Type:  &dal.TypeID{},
				Store: &dal.CodecAlias{Ident: "id"},
			},

			&dal.Attribute{
				Ident: "Handle",
				Type:  &dal.TypeText{Length: 255},
				Store: &dal.CodecAlias{Ident: "handle"},
			},

			&dal.Attribute{
				Ident: "Meta",
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "meta"},
			},

			&dal.Attribute{
				Ident: "Secret",
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "secret"},
			},

			&dal.Attribute{
				Ident: "Scope",
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "scope"},
			},

			&dal.Attribute{
				Ident: "ValidGrant",
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "valid_grant"},
			},

			&dal.Attribute{
				Ident: "RedirectURI",
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "redirect_uri"},
			},

			&dal.Attribute{
				Ident: "Enabled", Sortable: true,
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "enabled"},
			},

			&dal.Attribute{
				Ident: "Trusted", Sortable: true,
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "trusted"},
			},

			&dal.Attribute{
				Ident: "ValidFrom",
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "valid_from"},
			},

			&dal.Attribute{
				Ident: "ExpiresAt", Sortable: true,
				Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true},
				Store: &dal.CodecAlias{Ident: "expires_at"},
			},

			&dal.Attribute{
				Ident: "Security",
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "security"},
			},

			&dal.Attribute{
				Ident: "CreatedAt", Sortable: true,
				Type: &dal.TypeTimestamp{
					DefaultCurrentTimestamp: true,
					Timezone:                true,
				},
				Store: &dal.CodecAlias{Ident: "created_at"},
			},

			&dal.Attribute{
				Ident: "UpdatedAt", Sortable: true,
				Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true},
				Store: &dal.CodecAlias{Ident: "updated_at"},
			},

			&dal.Attribute{
				Ident: "DeletedAt", Sortable: true,
				Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true},
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

	AuthConfirmedClient = &dal.Model{
		Ident:        "auth_confirmed_clients",
		ResourceType: types.AuthConfirmedClientResourceType,

		Attributes: dal.AttributeSet{
			&dal.Attribute{
				Ident: "UserID", PrimaryKey: true,
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "rel_user"},
			},

			&dal.Attribute{
				Ident: "ClientID", PrimaryKey: true,
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "rel_client"},
			},

			&dal.Attribute{
				Ident: "ConfirmedAt", Sortable: true,
				Type:  &dal.TypeTimestamp{Timezone: true},
				Store: &dal.CodecAlias{Ident: "confirmed_at"},
			},
		},
	}

	AuthSession = &dal.Model{
		Ident:        "auth_sessions",
		ResourceType: types.AuthSessionResourceType,

		Attributes: dal.AttributeSet{
			&dal.Attribute{
				Ident: "ID", PrimaryKey: true,
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "id"},
			},

			&dal.Attribute{
				Ident: "Data",
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "data"},
			},

			&dal.Attribute{
				Ident: "UserID",
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "rel_user"},
			},

			&dal.Attribute{
				Ident: "ExpiresAt", Sortable: true,
				Type:  &dal.TypeTimestamp{Timezone: true},
				Store: &dal.CodecAlias{Ident: "expires_at"},
			},

			&dal.Attribute{
				Ident: "CreatedAt", Sortable: true,
				Type: &dal.TypeTimestamp{
					DefaultCurrentTimestamp: true,
					Timezone:                true,
				},
				Store: &dal.CodecAlias{Ident: "created_at"},
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
		},
	}

	AuthOa2token = &dal.Model{
		Ident:        "auth_oa2tokens",
		ResourceType: types.AuthOa2tokenResourceType,

		Attributes: dal.AttributeSet{
			&dal.Attribute{
				Ident: "ID", PrimaryKey: true,
				Type:  &dal.TypeID{},
				Store: &dal.CodecAlias{Ident: "id"},
			},

			&dal.Attribute{
				Ident: "Code",
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "code"},
			},

			&dal.Attribute{
				Ident: "Access",
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "access"},
			},

			&dal.Attribute{
				Ident: "Refresh",
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "refresh"},
			},

			&dal.Attribute{
				Ident: "ExpiresAt", Sortable: true,
				Type:  &dal.TypeTimestamp{Timezone: true},
				Store: &dal.CodecAlias{Ident: "expires_at"},
			},

			&dal.Attribute{
				Ident: "CreatedAt", Sortable: true,
				Type: &dal.TypeTimestamp{
					DefaultCurrentTimestamp: true,
					Timezone:                true,
				},
				Store: &dal.CodecAlias{Ident: "created_at"},
			},

			&dal.Attribute{
				Ident: "Data",
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "data"},
			},

			&dal.Attribute{
				Ident: "ClientID",
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "rel_client"},
			},

			&dal.Attribute{
				Ident: "UserID",
				Type:  &dal.TypeText{},
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
		},
	}

	Credential = &dal.Model{
		Ident:        "credentials",
		ResourceType: types.CredentialResourceType,

		Attributes: dal.AttributeSet{
			&dal.Attribute{
				Ident: "ID", PrimaryKey: true,
				Type:  &dal.TypeID{},
				Store: &dal.CodecAlias{Ident: "id"},
			},

			&dal.Attribute{
				Ident: "OwnerID",
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "rel_owner"},
			},

			&dal.Attribute{
				Ident: "Kind",
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "kind"},
			},

			&dal.Attribute{
				Ident: "Label",
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "label"},
			},

			&dal.Attribute{
				Ident: "Credentials",
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "credentials"},
			},

			&dal.Attribute{
				Ident: "Meta",
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "meta"},
			},

			&dal.Attribute{
				Ident: "CreatedAt", Sortable: true,
				Type: &dal.TypeTimestamp{
					DefaultCurrentTimestamp: true,
					Timezone:                true,
				},
				Store: &dal.CodecAlias{Ident: "created_at"},
			},

			&dal.Attribute{
				Ident: "UpdatedAt", Sortable: true,
				Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true},
				Store: &dal.CodecAlias{Ident: "updated_at"},
			},

			&dal.Attribute{
				Ident: "DeletedAt", Sortable: true,
				Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true},
				Store: &dal.CodecAlias{Ident: "deleted_at"},
			},

			&dal.Attribute{
				Ident: "LastUsedAt", Sortable: true,
				Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true},
				Store: &dal.CodecAlias{Ident: "last_used_at"},
			},

			&dal.Attribute{
				Ident: "ExpiresAt", Sortable: true,
				Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true},
				Store: &dal.CodecAlias{Ident: "expires_at"},
			},
		},
	}

	DataPrivacyRequest = &dal.Model{
		Ident:        "data_privacy_requests",
		ResourceType: types.DataPrivacyRequestResourceType,

		Attributes: dal.AttributeSet{
			&dal.Attribute{
				Ident: "ID", PrimaryKey: true,
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
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "status"},
			},

			&dal.Attribute{
				Ident: "Payload",
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "payload"},
			},

			&dal.Attribute{
				Ident: "RequestedAt", Sortable: true,
				Type:  &dal.TypeTimestamp{Timezone: true},
				Store: &dal.CodecAlias{Ident: "requested_at"},
			},

			&dal.Attribute{
				Ident: "RequestedBy",
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "requested_by"},
			},

			&dal.Attribute{
				Ident: "CompletedAt", Sortable: true,
				Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true},
				Store: &dal.CodecAlias{Ident: "completed_at"},
			},

			&dal.Attribute{
				Ident: "CompletedBy",
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "completed_by"},
			},

			&dal.Attribute{
				Ident: "CreatedAt", Sortable: true,
				Type: &dal.TypeTimestamp{
					DefaultCurrentTimestamp: true,
					Timezone:                true,
				},
				Store: &dal.CodecAlias{Ident: "created_at"},
			},

			&dal.Attribute{
				Ident: "UpdatedAt", Sortable: true,
				Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true},
				Store: &dal.CodecAlias{Ident: "updated_at"},
			},

			&dal.Attribute{
				Ident: "DeletedAt", Sortable: true,
				Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true},
				Store: &dal.CodecAlias{Ident: "deleted_at"},
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

	DataPrivacyRequestComment = &dal.Model{
		Ident:        "data_privacy_request_comments",
		ResourceType: types.DataPrivacyRequestCommentResourceType,

		Attributes: dal.AttributeSet{
			&dal.Attribute{
				Ident: "ID", PrimaryKey: true,
				Type:  &dal.TypeID{},
				Store: &dal.CodecAlias{Ident: "id"},
			},

			&dal.Attribute{
				Ident: "RequestID",
				Type:  &dal.TypeText{},
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
					DefaultCurrentTimestamp: true,
					Timezone:                true,
				},
				Store: &dal.CodecAlias{Ident: "created_at"},
			},

			&dal.Attribute{
				Ident: "UpdatedAt", Sortable: true,
				Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true},
				Store: &dal.CodecAlias{Ident: "updated_at"},
			},

			&dal.Attribute{
				Ident: "DeletedAt", Sortable: true,
				Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true},
				Store: &dal.CodecAlias{Ident: "deleted_at"},
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

	Queue = &dal.Model{
		Ident:        "queue_settings",
		ResourceType: types.QueueResourceType,

		Attributes: dal.AttributeSet{
			&dal.Attribute{
				Ident: "ID", PrimaryKey: true,
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
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "meta"},
			},

			&dal.Attribute{
				Ident: "CreatedAt", Sortable: true,
				Type: &dal.TypeTimestamp{
					DefaultCurrentTimestamp: true,
					Timezone:                true,
				},
				Store: &dal.CodecAlias{Ident: "created_at"},
			},

			&dal.Attribute{
				Ident: "UpdatedAt", Sortable: true,
				Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true},
				Store: &dal.CodecAlias{Ident: "updated_at"},
			},

			&dal.Attribute{
				Ident: "DeletedAt", Sortable: true,
				Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true},
				Store: &dal.CodecAlias{Ident: "deleted_at"},
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

	QueueMessage = &dal.Model{
		Ident:        "queue_messages",
		ResourceType: types.QueueMessageResourceType,

		Attributes: dal.AttributeSet{
			&dal.Attribute{
				Ident: "ID", PrimaryKey: true,
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
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "payload"},
			},

			&dal.Attribute{
				Ident: "Processed", Sortable: true,
				Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true},
				Store: &dal.CodecAlias{Ident: "processed"},
			},

			&dal.Attribute{
				Ident: "Created", Sortable: true,
				Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true},
				Store: &dal.CodecAlias{Ident: "created"},
			},
		},
	}

	Reminder = &dal.Model{
		Ident:        "reminders",
		ResourceType: types.ReminderResourceType,

		Attributes: dal.AttributeSet{
			&dal.Attribute{
				Ident: "ID", PrimaryKey: true,
				Type:  &dal.TypeID{},
				Store: &dal.CodecAlias{Ident: "id"},
			},

			&dal.Attribute{
				Ident: "Resource", Sortable: true,
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "resource"},
			},

			&dal.Attribute{
				Ident: "Payload",
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "payload"},
			},

			&dal.Attribute{
				Ident: "SnoozeCount",
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "snooze_count"},
			},

			&dal.Attribute{
				Ident: "AssignedTo",
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "assigned_to"},
			},

			&dal.Attribute{
				Ident: "AssignedBy",
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "assigned_by"},
			},

			&dal.Attribute{
				Ident: "AssignedAt", Sortable: true,
				Type:  &dal.TypeTimestamp{Timezone: true},
				Store: &dal.CodecAlias{Ident: "assigned_at"},
			},

			&dal.Attribute{
				Ident: "DismissedBy",
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "dismissed_by"},
			},

			&dal.Attribute{
				Ident: "DismissedAt", Sortable: true,
				Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true},
				Store: &dal.CodecAlias{Ident: "dismissed_at"},
			},

			&dal.Attribute{
				Ident: "RemindAt", Sortable: true,
				Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true},
				Store: &dal.CodecAlias{Ident: "remind_at"},
			},

			&dal.Attribute{
				Ident: "CreatedAt", Sortable: true,
				Type: &dal.TypeTimestamp{
					DefaultCurrentTimestamp: true,
					Timezone:                true,
				},
				Store: &dal.CodecAlias{Ident: "created_at"},
			},

			&dal.Attribute{
				Ident: "UpdatedAt", Sortable: true,
				Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true},
				Store: &dal.CodecAlias{Ident: "updated_at"},
			},

			&dal.Attribute{
				Ident: "DeletedAt", Sortable: true,
				Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true},
				Store: &dal.CodecAlias{Ident: "deleted_at"},
			},
		},
	}

	Report = &dal.Model{
		Ident:        "reports",
		ResourceType: types.ReportResourceType,

		Attributes: dal.AttributeSet{
			&dal.Attribute{
				Ident: "ID", PrimaryKey: true,
				Type:  &dal.TypeID{},
				Store: &dal.CodecAlias{Ident: "id"},
			},

			&dal.Attribute{
				Ident: "Handle",
				Type:  &dal.TypeText{Length: 255},
				Store: &dal.CodecAlias{Ident: "handle"},
			},

			&dal.Attribute{
				Ident: "Meta",
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "meta"},
			},

			&dal.Attribute{
				Ident: "Scenarios",
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "scenarios"},
			},

			&dal.Attribute{
				Ident: "Sources",
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "sources"},
			},

			&dal.Attribute{
				Ident: "Blocks",
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "blocks"},
			},

			&dal.Attribute{
				Ident: "CreatedAt", Sortable: true,
				Type: &dal.TypeTimestamp{
					DefaultCurrentTimestamp: true,
					Timezone:                true,
				},
				Store: &dal.CodecAlias{Ident: "created_at"},
			},

			&dal.Attribute{
				Ident: "UpdatedAt", Sortable: true,
				Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true},
				Store: &dal.CodecAlias{Ident: "updated_at"},
			},

			&dal.Attribute{
				Ident: "DeletedAt", Sortable: true,
				Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true},
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

	ResourceTranslation = &dal.Model{
		Ident:        "resource_translations",
		ResourceType: types.ResourceTranslationResourceType,

		Attributes: dal.AttributeSet{
			&dal.Attribute{
				Ident: "ID", PrimaryKey: true,
				Type:  &dal.TypeID{},
				Store: &dal.CodecAlias{Ident: "id"},
			},

			&dal.Attribute{
				Ident: "Lang",
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "lang"},
			},

			&dal.Attribute{
				Ident: "Resource",
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "resource"},
			},

			&dal.Attribute{
				Ident: "K",
				Type:  &dal.TypeText{},
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
					DefaultCurrentTimestamp: true,
					Timezone:                true,
				},
				Store: &dal.CodecAlias{Ident: "created_at"},
			},

			&dal.Attribute{
				Ident: "UpdatedAt", Sortable: true,
				Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true},
				Store: &dal.CodecAlias{Ident: "updated_at"},
			},

			&dal.Attribute{
				Ident: "DeletedAt", Sortable: true,
				Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true},
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

	Role = &dal.Model{
		Ident:        "roles",
		ResourceType: types.RoleResourceType,

		Attributes: dal.AttributeSet{
			&dal.Attribute{
				Ident: "ID", PrimaryKey: true,
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
				Type:  &dal.TypeText{Length: 255},
				Store: &dal.CodecAlias{Ident: "handle"},
			},

			&dal.Attribute{
				Ident: "Meta",
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "meta"},
			},

			&dal.Attribute{
				Ident: "CreatedAt", Sortable: true,
				Type: &dal.TypeTimestamp{
					DefaultCurrentTimestamp: true,
					Timezone:                true,
				},
				Store: &dal.CodecAlias{Ident: "created_at"},
			},

			&dal.Attribute{
				Ident: "UpdatedAt", Sortable: true,
				Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true},
				Store: &dal.CodecAlias{Ident: "updated_at"},
			},

			&dal.Attribute{
				Ident: "DeletedAt", Sortable: true,
				Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true},
				Store: &dal.CodecAlias{Ident: "deleted_at"},
			},

			&dal.Attribute{
				Ident: "ArchivedAt", Sortable: true,
				Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true},
				Store: &dal.CodecAlias{Ident: "archived_at"},
			},
		},
	}

	RoleMember = &dal.Model{
		Ident:        "role_members",
		ResourceType: types.RoleMemberResourceType,

		Attributes: dal.AttributeSet{
			&dal.Attribute{
				Ident: "UserID", PrimaryKey: true,
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "rel_user"},
			},

			&dal.Attribute{
				Ident: "RoleID", PrimaryKey: true,
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "rel_role"},
			},
		},
	}

	SettingValue = &dal.Model{
		Ident:        "settings",
		ResourceType: types.SettingValueResourceType,

		Attributes: dal.AttributeSet{
			&dal.Attribute{
				Ident: "Name", PrimaryKey: true,
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "name"},
			},

			&dal.Attribute{
				Ident: "OwnedBy", PrimaryKey: true,
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "rel_owner"},
			},

			&dal.Attribute{
				Ident: "Value",
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "value"},
			},

			&dal.Attribute{
				Ident: "UpdatedBy",
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "updated_by"},
			},

			&dal.Attribute{
				Ident: "UpdatedAt", Sortable: true,
				Type:  &dal.TypeTimestamp{Timezone: true},
				Store: &dal.CodecAlias{Ident: "updated_at"},
			},
		},
	}

	Template = &dal.Model{
		Ident:        "templates",
		ResourceType: types.TemplateResourceType,

		Attributes: dal.AttributeSet{
			&dal.Attribute{
				Ident: "ID", PrimaryKey: true,
				Type:  &dal.TypeID{},
				Store: &dal.CodecAlias{Ident: "id"},
			},

			&dal.Attribute{
				Ident: "Handle",
				Type:  &dal.TypeText{Length: 255},
				Store: &dal.CodecAlias{Ident: "handle"},
			},

			&dal.Attribute{
				Ident: "Language", Sortable: true,
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "language"},
			},

			&dal.Attribute{
				Ident: "Type", Sortable: true,
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "type"},
			},

			&dal.Attribute{
				Ident: "Partial",
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "partial"},
			},

			&dal.Attribute{
				Ident: "Meta",
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "meta"},
			},

			&dal.Attribute{
				Ident: "Template", Sortable: true,
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "template"},
			},

			&dal.Attribute{
				Ident: "OwnerID",
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "rel_owner"},
			},

			&dal.Attribute{
				Ident: "CreatedAt", Sortable: true,
				Type: &dal.TypeTimestamp{
					DefaultCurrentTimestamp: true,
					Timezone:                true,
				},
				Store: &dal.CodecAlias{Ident: "created_at"},
			},

			&dal.Attribute{
				Ident: "UpdatedAt", Sortable: true,
				Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true},
				Store: &dal.CodecAlias{Ident: "updated_at"},
			},

			&dal.Attribute{
				Ident: "DeletedAt", Sortable: true,
				Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true},
				Store: &dal.CodecAlias{Ident: "deleted_at"},
			},

			&dal.Attribute{
				Ident: "LastUsedAt", Sortable: true,
				Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true},
				Store: &dal.CodecAlias{Ident: "last_used_at"},
			},
		},
	}

	User = &dal.Model{
		Ident:        "users",
		ResourceType: types.UserResourceType,

		Attributes: dal.AttributeSet{
			&dal.Attribute{
				Ident: "ID", PrimaryKey: true,
				Type:  &dal.TypeID{},
				Store: &dal.CodecAlias{Ident: "id"},
			},

			&dal.Attribute{
				Ident: "Handle",
				Type:  &dal.TypeText{Length: 255},
				Store: &dal.CodecAlias{Ident: "handle"},
			},

			&dal.Attribute{
				Ident: "Email", Sortable: true,
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "email"},
			},

			&dal.Attribute{
				Ident: "EmailConfirmed",
				Type:  &dal.TypeText{},
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
				Ident: "Kind", Sortable: true,
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "kind"},
			},

			&dal.Attribute{
				Ident: "Meta",
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "meta"},
			},

			&dal.Attribute{
				Ident: "CreatedAt", Sortable: true,
				Type: &dal.TypeTimestamp{
					DefaultCurrentTimestamp: true,
					Timezone:                true,
				},
				Store: &dal.CodecAlias{Ident: "created_at"},
			},

			&dal.Attribute{
				Ident: "UpdatedAt", Sortable: true,
				Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true},
				Store: &dal.CodecAlias{Ident: "updated_at"},
			},

			&dal.Attribute{
				Ident: "DeletedAt", Sortable: true,
				Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true},
				Store: &dal.CodecAlias{Ident: "deleted_at"},
			},

			&dal.Attribute{
				Ident: "SuspendedAt", Sortable: true,
				Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true},
				Store: &dal.CodecAlias{Ident: "suspended_at"},
			},
		},
	}

	DalConnection = &dal.Model{
		Ident:        "dal_connections",
		ResourceType: types.DalConnectionResourceType,

		Attributes: dal.AttributeSet{
			&dal.Attribute{
				Ident: "ID", PrimaryKey: true,
				Type:  &dal.TypeID{},
				Store: &dal.CodecAlias{Ident: "id"},
			},

			&dal.Attribute{
				Ident: "Handle",
				Type:  &dal.TypeText{Length: 255},
				Store: &dal.CodecAlias{Ident: "handle"},
			},

			&dal.Attribute{
				Ident: "Type", Sortable: true,
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "type"},
			},

			&dal.Attribute{
				Ident: "Meta",
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "meta"},
			},

			&dal.Attribute{
				Ident: "Config",
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "config"},
			},

			&dal.Attribute{
				Ident: "CreatedAt", Sortable: true,
				Type: &dal.TypeTimestamp{
					DefaultCurrentTimestamp: true,
					Timezone:                true,
				},
				Store: &dal.CodecAlias{Ident: "created_at"},
			},

			&dal.Attribute{
				Ident: "UpdatedAt", Sortable: true,
				Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true},
				Store: &dal.CodecAlias{Ident: "updated_at"},
			},

			&dal.Attribute{
				Ident: "DeletedAt", Sortable: true,
				Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true},
				Store: &dal.CodecAlias{Ident: "deleted_at"},
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

	DalSensitivityLevel = &dal.Model{
		Ident:        "dal_sensitivity_levels",
		ResourceType: types.DalSensitivityLevelResourceType,

		Attributes: dal.AttributeSet{
			&dal.Attribute{
				Ident: "ID", PrimaryKey: true,
				Type:  &dal.TypeID{},
				Store: &dal.CodecAlias{Ident: "id"},
			},

			&dal.Attribute{
				Ident: "Handle",
				Type:  &dal.TypeText{Length: 255},
				Store: &dal.CodecAlias{Ident: "handle"},
			},

			&dal.Attribute{
				Ident: "Level", Sortable: true,
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "level"},
			},

			&dal.Attribute{
				Ident: "Meta",
				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "meta"},
			},

			&dal.Attribute{
				Ident: "CreatedAt", Sortable: true,
				Type: &dal.TypeTimestamp{
					DefaultCurrentTimestamp: true,
					Timezone:                true,
				},
				Store: &dal.CodecAlias{Ident: "created_at"},
			},

			&dal.Attribute{
				Ident: "UpdatedAt", Sortable: true,
				Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true},
				Store: &dal.CodecAlias{Ident: "updated_at"},
			},

			&dal.Attribute{
				Ident: "DeletedAt", Sortable: true,
				Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true},
				Store: &dal.CodecAlias{Ident: "deleted_at"},
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
		Application,
		ApigwRoute,
		ApigwFilter,
		AuthClient,
		AuthConfirmedClient,
		AuthSession,
		AuthOa2token,
		Credential,
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
		DalConnection,
		DalSensitivityLevel,
	}
}

func Register(ctx context.Context, mr modelReplacer) (err error) {
	if err = mr.ReplaceModel(ctx, Attachment); err != nil {
		return
	}

	if err = mr.ReplaceModel(ctx, Application); err != nil {
		return
	}

	if err = mr.ReplaceModel(ctx, ApigwRoute); err != nil {
		return
	}

	if err = mr.ReplaceModel(ctx, ApigwFilter); err != nil {
		return
	}

	if err = mr.ReplaceModel(ctx, AuthClient); err != nil {
		return
	}

	if err = mr.ReplaceModel(ctx, AuthConfirmedClient); err != nil {
		return
	}

	if err = mr.ReplaceModel(ctx, AuthSession); err != nil {
		return
	}

	if err = mr.ReplaceModel(ctx, AuthOa2token); err != nil {
		return
	}

	if err = mr.ReplaceModel(ctx, Credential); err != nil {
		return
	}

	if err = mr.ReplaceModel(ctx, DataPrivacyRequest); err != nil {
		return
	}

	if err = mr.ReplaceModel(ctx, DataPrivacyRequestComment); err != nil {
		return
	}

	if err = mr.ReplaceModel(ctx, Queue); err != nil {
		return
	}

	if err = mr.ReplaceModel(ctx, QueueMessage); err != nil {
		return
	}

	if err = mr.ReplaceModel(ctx, Reminder); err != nil {
		return
	}

	if err = mr.ReplaceModel(ctx, Report); err != nil {
		return
	}

	if err = mr.ReplaceModel(ctx, ResourceTranslation); err != nil {
		return
	}

	if err = mr.ReplaceModel(ctx, Role); err != nil {
		return
	}

	if err = mr.ReplaceModel(ctx, RoleMember); err != nil {
		return
	}

	if err = mr.ReplaceModel(ctx, SettingValue); err != nil {
		return
	}

	if err = mr.ReplaceModel(ctx, Template); err != nil {
		return
	}

	if err = mr.ReplaceModel(ctx, User); err != nil {
		return
	}

	if err = mr.ReplaceModel(ctx, DalConnection); err != nil {
		return
	}

	if err = mr.ReplaceModel(ctx, DalSensitivityLevel); err != nil {
		return
	}

	return
}
