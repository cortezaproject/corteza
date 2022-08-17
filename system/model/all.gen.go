package model

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"github.com/cortezaproject/corteza-server/pkg/dal"
	"github.com/cortezaproject/corteza-server/system/types"
)

var (
	attachmentModel = &dal.Model{
		Ident: "attachmentModel",

		ResourceType: types.AttachmentResourceType,
		Attributes: dal.AttributeSet{

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "createdAt",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "deletedAt",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "id", PrimaryKey: true,
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "kind",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "meta",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "name",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "ownerID",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "previewUrl",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "updatedAt",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "url",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},
		},
	}

	applicationModel = &dal.Model{
		Ident: "applicationModel",

		ResourceType: types.ApplicationResourceType,
		Attributes: dal.AttributeSet{

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "createdAt",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "deletedAt",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "enabled",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "id", PrimaryKey: true,
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "name",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "ownerID",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "unify",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "updatedAt",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "weight",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},
		},
	}

	apigwRouteModel = &dal.Model{
		Ident: "apigwRouteModel",

		ResourceType: types.ApigwRouteResourceType,
		Attributes: dal.AttributeSet{

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "createdAt",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "createdBy",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "deletedAt",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "deletedBy",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "enabled",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "endpoint",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "group",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "id", PrimaryKey: true,
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "meta",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "method",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "updatedAt",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "updatedBy",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},
		},
	}

	apigwFilterModel = &dal.Model{
		Ident: "apigwFilterModel",

		ResourceType: types.ApigwFilterResourceType,
		Attributes: dal.AttributeSet{

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "createdAt",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "createdBy",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "deletedAt",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "deletedBy",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "enabled",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "id", PrimaryKey: true,
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "kind",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "params",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "ref",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "route",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "updatedAt",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "updatedBy",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "weight",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},
		},
	}

	authClientModel = &dal.Model{
		Ident: "authClientModel",

		ResourceType: types.AuthClientResourceType,
		Attributes: dal.AttributeSet{

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "createdAt",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "createdBy",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "deletedAt",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "deletedBy",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "enabled",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "expiresAt",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "handle",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "id", PrimaryKey: true,
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "meta",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "ownedBy",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "redirectURI",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "scope",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "secret",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "security",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "trusted",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "updatedAt",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "updatedBy",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "validFrom",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "validGrant",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},
		},
	}

	authConfirmedClientModel = &dal.Model{
		Ident: "authConfirmedClientModel",

		ResourceType: types.AuthConfirmedClientResourceType,
		Attributes: dal.AttributeSet{

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "clientID", PrimaryKey: true,
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "confirmedAt",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "userID", PrimaryKey: true,
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},
		},
	}

	authSessionModel = &dal.Model{
		Ident: "authSessionModel",

		ResourceType: types.AuthSessionResourceType,
		Attributes: dal.AttributeSet{

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "createdAt",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "data",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "expiresAt",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "id", PrimaryKey: true,
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "remoteAddr",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "userAgent",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "userID",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},
		},
	}

	authOa2tokenModel = &dal.Model{
		Ident: "authOa2tokenModel",

		ResourceType: types.AuthOa2tokenResourceType,
		Attributes: dal.AttributeSet{

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "access",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "clientID",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "code",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "createdAt",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "data",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "expiresAt",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "id", PrimaryKey: true,
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "refresh",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "remoteAddr",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "userAgent",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "userID",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},
		},
	}

	credentialModel = &dal.Model{
		Ident: "credentialModel",

		ResourceType: types.CredentialResourceType,
		Attributes: dal.AttributeSet{

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "createdAt",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "credentials",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "deletedAt",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "expiresAt",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "id", PrimaryKey: true,
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "kind",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "label",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "lastUsedAt",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "meta",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "ownerID",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "updatedAt",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},
		},
	}

	dataPrivacyRequestModel = &dal.Model{
		Ident: "dataPrivacyRequestModel",

		ResourceType: types.DataPrivacyRequestResourceType,
		Attributes: dal.AttributeSet{

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "completedAt",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "completedBy",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "createdAt",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "createdBy",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "deletedAt",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "deletedBy",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "id", PrimaryKey: true,
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "kind",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "payload",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "requestedAt",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "requestedBy",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "status",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "updatedAt",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "updatedBy",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},
		},
	}

	dataPrivacyRequestCommentModel = &dal.Model{
		Ident: "dataPrivacyRequestCommentModel",

		ResourceType: types.DataPrivacyRequestCommentResourceType,
		Attributes: dal.AttributeSet{

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "comment",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "createdAt",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "createdBy",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "deletedAt",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "deletedBy",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "id", PrimaryKey: true,
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "requestID",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "updatedAt",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "updatedBy",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},
		},
	}

	queueModel = &dal.Model{
		Ident: "queueModel",

		ResourceType: types.QueueResourceType,
		Attributes: dal.AttributeSet{

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "consumer",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "createdAt",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "createdBy",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "deletedAt",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "deletedBy",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "id", PrimaryKey: true,
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "meta",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "queue",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "updatedAt",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "updatedBy",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},
		},
	}

	queueMessageModel = &dal.Model{
		Ident: "queueMessageModel",

		ResourceType: types.QueueMessageResourceType,
		Attributes: dal.AttributeSet{

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "created",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "id", PrimaryKey: true,
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "payload",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "processed",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "queue",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},
		},
	}

	reminderModel = &dal.Model{
		Ident: "reminderModel",

		ResourceType: types.ReminderResourceType,
		Attributes: dal.AttributeSet{

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "assignedAt",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "assignedBy",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "assignedTo",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "createdAt",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "deletedAt",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "dismissedAt",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "dismissedBy",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "id", PrimaryKey: true,
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "payload",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "remindAt",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "resource",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "snoozeCount",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "updatedAt",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},
		},
	}

	reportModel = &dal.Model{
		Ident: "reportModel",

		ResourceType: types.ReportResourceType,
		Attributes: dal.AttributeSet{

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "blocks",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "createdAt",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "createdBy",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "deletedAt",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "deletedBy",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "handle",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "id", PrimaryKey: true,
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "meta",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "ownedBy",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "scenarios",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "sources",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "updatedAt",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "updatedBy",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},
		},
	}

	resourceTranslationModel = &dal.Model{
		Ident: "resourceTranslationModel",

		ResourceType: types.ResourceTranslationResourceType,
		Attributes: dal.AttributeSet{

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "createdAt",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "createdBy",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "deletedAt",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "deletedBy",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "id", PrimaryKey: true,
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "k",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "lang",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "message",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "ownedBy",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "resource",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "updatedAt",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "updatedBy",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},
		},
	}

	roleModel = &dal.Model{
		Ident: "roleModel",

		ResourceType: types.RoleResourceType,
		Attributes: dal.AttributeSet{

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "archivedAt",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "createdAt",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "deletedAt",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "handle",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "id", PrimaryKey: true,
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "meta",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "name",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "updatedAt",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},
		},
	}

	roleMemberModel = &dal.Model{
		Ident: "roleMemberModel",

		ResourceType: types.RoleMemberResourceType,
		Attributes: dal.AttributeSet{

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "roleID", PrimaryKey: true,
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "userID", PrimaryKey: true,
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},
		},
	}

	settingValueModel = &dal.Model{
		Ident: "settingValueModel",

		ResourceType: types.SettingValueResourceType,
		Attributes: dal.AttributeSet{

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "name", PrimaryKey: true,
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "ownedBy", PrimaryKey: true,
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "updatedAt",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "updatedBy",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "value",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},
		},
	}

	templateModel = &dal.Model{
		Ident: "templateModel",

		ResourceType: types.TemplateResourceType,
		Attributes: dal.AttributeSet{

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "createdAt",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "deletedAt",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "handle",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "id", PrimaryKey: true,
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "language",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "lastUsedAt",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "meta",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "ownerID",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "partial",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "template",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "type",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "updatedAt",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},
		},
	}

	userModel = &dal.Model{
		Ident: "userModel",

		ResourceType: types.UserResourceType,
		Attributes: dal.AttributeSet{

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "createdAt",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "deletedAt",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "email",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "emailConfirmed",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "handle",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "id", PrimaryKey: true,
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "kind",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "meta",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "name",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "suspendedAt",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "updatedAt",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "username",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},
		},
	}

	dalConnectionModel = &dal.Model{
		Ident: "dalConnectionModel",

		ResourceType: types.DalConnectionResourceType,
		Attributes: dal.AttributeSet{

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "config",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "createdAt",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "createdBy",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "deletedAt",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "deletedBy",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "handle",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "id", PrimaryKey: true,
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "meta",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "type",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "updatedAt",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "updatedBy",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},
		},
	}

	dalSensitivityLevelModel = &dal.Model{
		Ident: "dalSensitivityLevelModel",

		ResourceType: types.DalSensitivityLevelResourceType,
		Attributes: dal.AttributeSet{

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "createdAt",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "createdBy",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "deletedAt",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "deletedBy",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "handle",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "id", PrimaryKey: true,
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "level",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "meta",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "updatedAt",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},

			// &dal.Attribute{Ident: "id", PrimaryKey: true, Store: &dal.CodecPlain{}, Type: &dal.TypeID{}},
			&dal.Attribute{
				Ident: "updatedBy",
				Store: &dal.CodecPlain{},
				Type:  &dal.TypeID{},
			},
		},
	}
)

// Attachment returns attachmentModel
//
// This function is auto-generated
func Attachment() *dal.Model {
	return attachmentModel
}

// Application returns applicationModel
//
// This function is auto-generated
func Application() *dal.Model {
	return applicationModel
}

// ApigwRoute returns apigwRouteModel
//
// This function is auto-generated
func ApigwRoute() *dal.Model {
	return apigwRouteModel
}

// ApigwFilter returns apigwFilterModel
//
// This function is auto-generated
func ApigwFilter() *dal.Model {
	return apigwFilterModel
}

// AuthClient returns authClientModel
//
// This function is auto-generated
func AuthClient() *dal.Model {
	return authClientModel
}

// AuthConfirmedClient returns authConfirmedClientModel
//
// This function is auto-generated
func AuthConfirmedClient() *dal.Model {
	return authConfirmedClientModel
}

// AuthSession returns authSessionModel
//
// This function is auto-generated
func AuthSession() *dal.Model {
	return authSessionModel
}

// AuthOa2token returns authOa2tokenModel
//
// This function is auto-generated
func AuthOa2token() *dal.Model {
	return authOa2tokenModel
}

// Credential returns credentialModel
//
// This function is auto-generated
func Credential() *dal.Model {
	return credentialModel
}

// DataPrivacyRequest returns dataPrivacyRequestModel
//
// This function is auto-generated
func DataPrivacyRequest() *dal.Model {
	return dataPrivacyRequestModel
}

// DataPrivacyRequestComment returns dataPrivacyRequestCommentModel
//
// This function is auto-generated
func DataPrivacyRequestComment() *dal.Model {
	return dataPrivacyRequestCommentModel
}

// Queue returns queueModel
//
// This function is auto-generated
func Queue() *dal.Model {
	return queueModel
}

// QueueMessage returns queueMessageModel
//
// This function is auto-generated
func QueueMessage() *dal.Model {
	return queueMessageModel
}

// Reminder returns reminderModel
//
// This function is auto-generated
func Reminder() *dal.Model {
	return reminderModel
}

// Report returns reportModel
//
// This function is auto-generated
func Report() *dal.Model {
	return reportModel
}

// ResourceTranslation returns resourceTranslationModel
//
// This function is auto-generated
func ResourceTranslation() *dal.Model {
	return resourceTranslationModel
}

// Role returns roleModel
//
// This function is auto-generated
func Role() *dal.Model {
	return roleModel
}

// RoleMember returns roleMemberModel
//
// This function is auto-generated
func RoleMember() *dal.Model {
	return roleMemberModel
}

// SettingValue returns settingValueModel
//
// This function is auto-generated
func SettingValue() *dal.Model {
	return settingValueModel
}

// Template returns templateModel
//
// This function is auto-generated
func Template() *dal.Model {
	return templateModel
}

// User returns userModel
//
// This function is auto-generated
func User() *dal.Model {
	return userModel
}

// DalConnection returns dalConnectionModel
//
// This function is auto-generated
func DalConnection() *dal.Model {
	return dalConnectionModel
}

// DalSensitivityLevel returns dalSensitivityLevelModel
//
// This function is auto-generated
func DalSensitivityLevel() *dal.Model {
	return dalSensitivityLevelModel
}
