package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"github.com/cortezaproject/corteza/server/pkg/cast2"
)

func (r Attachment) GetID() uint64 { return r.ID }

func (r *Attachment) GetValue(name string, pos uint) (any, error) {
	switch name {
	case "createdAt", "CreatedAt":
		return r.CreatedAt, nil
	case "deletedAt", "DeletedAt":
		return r.DeletedAt, nil
	case "id", "ID":
		return r.ID, nil
	case "kind", "Kind":
		return r.Kind, nil
	case "name", "Name":
		return r.Name, nil
	case "ownerID", "OwnerID":
		return r.OwnerID, nil
	case "previewUrl", "PreviewUrl":
		return r.PreviewUrl, nil
	case "updatedAt", "UpdatedAt":
		return r.UpdatedAt, nil
	case "url", "Url":
		return r.Url, nil

	}
	return nil, nil
}

func (r *Attachment) SetValue(name string, pos uint, value any) (err error) {
	switch name {
	case "createdAt", "CreatedAt":
		return cast2.Time(value, &r.CreatedAt)
	case "deletedAt", "DeletedAt":
		return cast2.TimePtr(value, &r.DeletedAt)
	case "id", "ID":
		return cast2.Uint64(value, &r.ID)
	case "kind", "Kind":
		return cast2.String(value, &r.Kind)
	case "name", "Name":
		return cast2.String(value, &r.Name)
	case "ownerID", "OwnerID":
		return cast2.Uint64(value, &r.OwnerID)
	case "previewUrl", "PreviewUrl":
		return cast2.String(value, &r.PreviewUrl)
	case "updatedAt", "UpdatedAt":
		return cast2.TimePtr(value, &r.UpdatedAt)
	case "url", "Url":
		return cast2.String(value, &r.Url)

	}
	return nil
}

func (r Application) GetID() uint64 { return r.ID }

func (r *Application) GetValue(name string, pos uint) (any, error) {
	switch name {
	case "createdAt", "CreatedAt":
		return r.CreatedAt, nil
	case "deletedAt", "DeletedAt":
		return r.DeletedAt, nil
	case "enabled", "Enabled":
		return r.Enabled, nil
	case "id", "ID":
		return r.ID, nil
	case "name", "Name":
		return r.Name, nil
	case "ownerID", "OwnerID":
		return r.OwnerID, nil
	case "updatedAt", "UpdatedAt":
		return r.UpdatedAt, nil
	case "weight", "Weight":
		return r.Weight, nil

	}
	return nil, nil
}

func (r *Application) SetValue(name string, pos uint, value any) (err error) {
	switch name {
	case "createdAt", "CreatedAt":
		return cast2.Time(value, &r.CreatedAt)
	case "deletedAt", "DeletedAt":
		return cast2.TimePtr(value, &r.DeletedAt)
	case "enabled", "Enabled":
		return cast2.Bool(value, &r.Enabled)
	case "id", "ID":
		return cast2.Uint64(value, &r.ID)
	case "name", "Name":
		return cast2.String(value, &r.Name)
	case "ownerID", "OwnerID":
		return cast2.Uint64(value, &r.OwnerID)
	case "updatedAt", "UpdatedAt":
		return cast2.TimePtr(value, &r.UpdatedAt)
	case "weight", "Weight":
		return cast2.Int(value, &r.Weight)

	}
	return nil
}

func (r ApigwRoute) GetID() uint64 { return r.ID }

func (r *ApigwRoute) GetValue(name string, pos uint) (any, error) {
	switch name {
	case "createdAt", "CreatedAt":
		return r.CreatedAt, nil
	case "createdBy", "CreatedBy":
		return r.CreatedBy, nil
	case "deletedAt", "DeletedAt":
		return r.DeletedAt, nil
	case "deletedBy", "DeletedBy":
		return r.DeletedBy, nil
	case "enabled", "Enabled":
		return r.Enabled, nil
	case "endpoint", "Endpoint":
		return r.Endpoint, nil
	case "group", "Group":
		return r.Group, nil
	case "id", "ID":
		return r.ID, nil
	case "method", "Method":
		return r.Method, nil
	case "updatedAt", "UpdatedAt":
		return r.UpdatedAt, nil
	case "updatedBy", "UpdatedBy":
		return r.UpdatedBy, nil

	}
	return nil, nil
}

func (r *ApigwRoute) SetValue(name string, pos uint, value any) (err error) {
	switch name {
	case "createdAt", "CreatedAt":
		return cast2.Time(value, &r.CreatedAt)
	case "createdBy", "CreatedBy":
		return cast2.Uint64(value, &r.CreatedBy)
	case "deletedAt", "DeletedAt":
		return cast2.TimePtr(value, &r.DeletedAt)
	case "deletedBy", "DeletedBy":
		return cast2.Uint64(value, &r.DeletedBy)
	case "enabled", "Enabled":
		return cast2.Bool(value, &r.Enabled)
	case "endpoint", "Endpoint":
		return cast2.String(value, &r.Endpoint)
	case "group", "Group":
		return cast2.Uint64(value, &r.Group)
	case "id", "ID":
		return cast2.Uint64(value, &r.ID)
	case "method", "Method":
		return cast2.String(value, &r.Method)
	case "updatedAt", "UpdatedAt":
		return cast2.TimePtr(value, &r.UpdatedAt)
	case "updatedBy", "UpdatedBy":
		return cast2.Uint64(value, &r.UpdatedBy)

	}
	return nil
}

func (r ApigwFilter) GetID() uint64 { return r.ID }

func (r *ApigwFilter) GetValue(name string, pos uint) (any, error) {
	switch name {
	case "createdAt", "CreatedAt":
		return r.CreatedAt, nil
	case "createdBy", "CreatedBy":
		return r.CreatedBy, nil
	case "deletedAt", "DeletedAt":
		return r.DeletedAt, nil
	case "deletedBy", "DeletedBy":
		return r.DeletedBy, nil
	case "enabled", "Enabled":
		return r.Enabled, nil
	case "id", "ID":
		return r.ID, nil
	case "kind", "Kind":
		return r.Kind, nil
	case "ref", "Ref":
		return r.Ref, nil
	case "route", "Route", "ApigwRouteID":
		return r.Route, nil
	case "updatedAt", "UpdatedAt":
		return r.UpdatedAt, nil
	case "updatedBy", "UpdatedBy":
		return r.UpdatedBy, nil
	case "weight", "Weight":
		return r.Weight, nil

	}
	return nil, nil
}

func (r *ApigwFilter) SetValue(name string, pos uint, value any) (err error) {
	switch name {
	case "createdAt", "CreatedAt":
		return cast2.Time(value, &r.CreatedAt)
	case "createdBy", "CreatedBy":
		return cast2.Uint64(value, &r.CreatedBy)
	case "deletedAt", "DeletedAt":
		return cast2.TimePtr(value, &r.DeletedAt)
	case "deletedBy", "DeletedBy":
		return cast2.Uint64(value, &r.DeletedBy)
	case "enabled", "Enabled":
		return cast2.Bool(value, &r.Enabled)
	case "id", "ID":
		return cast2.Uint64(value, &r.ID)
	case "kind", "Kind":
		return cast2.String(value, &r.Kind)
	case "ref", "Ref":
		return cast2.String(value, &r.Ref)
	case "route", "Route", "ApigwRouteID":
		return cast2.Uint64(value, &r.Route)
	case "updatedAt", "UpdatedAt":
		return cast2.TimePtr(value, &r.UpdatedAt)
	case "updatedBy", "UpdatedBy":
		return cast2.Uint64(value, &r.UpdatedBy)
	case "weight", "Weight":
		return cast2.Uint64(value, &r.Weight)

	}
	return nil
}

func (r AuthClient) GetID() uint64 { return r.ID }

func (r *AuthClient) GetValue(name string, pos uint) (any, error) {
	switch name {
	case "createdAt", "CreatedAt":
		return r.CreatedAt, nil
	case "createdBy", "CreatedBy":
		return r.CreatedBy, nil
	case "deletedAt", "DeletedAt":
		return r.DeletedAt, nil
	case "deletedBy", "DeletedBy":
		return r.DeletedBy, nil
	case "enabled", "Enabled":
		return r.Enabled, nil
	case "expiresAt", "ExpiresAt":
		return r.ExpiresAt, nil
	case "handle", "Handle":
		return r.Handle, nil
	case "id", "ID":
		return r.ID, nil
	case "ownedBy", "OwnedBy":
		return r.OwnedBy, nil
	case "redirectURI", "RedirectURI":
		return r.RedirectURI, nil
	case "scope", "Scope":
		return r.Scope, nil
	case "secret", "Secret":
		return r.Secret, nil
	case "trusted", "Trusted":
		return r.Trusted, nil
	case "updatedAt", "UpdatedAt":
		return r.UpdatedAt, nil
	case "updatedBy", "UpdatedBy":
		return r.UpdatedBy, nil
	case "validFrom", "ValidFrom":
		return r.ValidFrom, nil
	case "validGrant", "ValidGrant":
		return r.ValidGrant, nil

	}
	return nil, nil
}

func (r *AuthClient) SetValue(name string, pos uint, value any) (err error) {
	switch name {
	case "createdAt", "CreatedAt":
		return cast2.Time(value, &r.CreatedAt)
	case "createdBy", "CreatedBy":
		return cast2.Uint64(value, &r.CreatedBy)
	case "deletedAt", "DeletedAt":
		return cast2.TimePtr(value, &r.DeletedAt)
	case "deletedBy", "DeletedBy":
		return cast2.Uint64(value, &r.DeletedBy)
	case "enabled", "Enabled":
		return cast2.Bool(value, &r.Enabled)
	case "expiresAt", "ExpiresAt":
		return cast2.TimePtr(value, &r.ExpiresAt)
	case "handle", "Handle":
		return cast2.String(value, &r.Handle)
	case "id", "ID":
		return cast2.Uint64(value, &r.ID)
	case "ownedBy", "OwnedBy":
		return cast2.Uint64(value, &r.OwnedBy)
	case "redirectURI", "RedirectURI":
		return cast2.String(value, &r.RedirectURI)
	case "scope", "Scope":
		return cast2.String(value, &r.Scope)
	case "secret", "Secret":
		return cast2.String(value, &r.Secret)
	case "trusted", "Trusted":
		return cast2.Bool(value, &r.Trusted)
	case "updatedAt", "UpdatedAt":
		return cast2.TimePtr(value, &r.UpdatedAt)
	case "updatedBy", "UpdatedBy":
		return cast2.Uint64(value, &r.UpdatedBy)
	case "validFrom", "ValidFrom":
		return cast2.TimePtr(value, &r.ValidFrom)
	case "validGrant", "ValidGrant":
		return cast2.String(value, &r.ValidGrant)

	}
	return nil
}

func (r DataPrivacyRequestComment) GetID() uint64 { return r.ID }

func (r *DataPrivacyRequestComment) GetValue(name string, pos uint) (any, error) {
	switch name {
	case "comment", "Comment":
		return r.Comment, nil
	case "createdAt", "CreatedAt":
		return r.CreatedAt, nil
	case "createdBy", "CreatedBy":
		return r.CreatedBy, nil
	case "deletedAt", "DeletedAt":
		return r.DeletedAt, nil
	case "deletedBy", "DeletedBy":
		return r.DeletedBy, nil
	case "id", "ID":
		return r.ID, nil
	case "requestID", "RequestID":
		return r.RequestID, nil
	case "updatedAt", "UpdatedAt":
		return r.UpdatedAt, nil
	case "updatedBy", "UpdatedBy":
		return r.UpdatedBy, nil

	}
	return nil, nil
}

func (r *DataPrivacyRequestComment) SetValue(name string, pos uint, value any) (err error) {
	switch name {
	case "comment", "Comment":
		return cast2.String(value, &r.Comment)
	case "createdAt", "CreatedAt":
		return cast2.Time(value, &r.CreatedAt)
	case "createdBy", "CreatedBy":
		return cast2.Uint64(value, &r.CreatedBy)
	case "deletedAt", "DeletedAt":
		return cast2.TimePtr(value, &r.DeletedAt)
	case "deletedBy", "DeletedBy":
		return cast2.Uint64(value, &r.DeletedBy)
	case "id", "ID":
		return cast2.Uint64(value, &r.ID)
	case "requestID", "RequestID":
		return cast2.Uint64(value, &r.RequestID)
	case "updatedAt", "UpdatedAt":
		return cast2.TimePtr(value, &r.UpdatedAt)
	case "updatedBy", "UpdatedBy":
		return cast2.Uint64(value, &r.UpdatedBy)

	}
	return nil
}

func (r Queue) GetID() uint64 { return r.ID }

func (r *Queue) GetValue(name string, pos uint) (any, error) {
	switch name {
	case "consumer", "Consumer":
		return r.Consumer, nil
	case "createdAt", "CreatedAt":
		return r.CreatedAt, nil
	case "createdBy", "CreatedBy":
		return r.CreatedBy, nil
	case "deletedAt", "DeletedAt":
		return r.DeletedAt, nil
	case "deletedBy", "DeletedBy":
		return r.DeletedBy, nil
	case "id", "ID":
		return r.ID, nil
	case "queue", "Queue":
		return r.Queue, nil
	case "updatedAt", "UpdatedAt":
		return r.UpdatedAt, nil
	case "updatedBy", "UpdatedBy":
		return r.UpdatedBy, nil

	}
	return nil, nil
}

func (r *Queue) SetValue(name string, pos uint, value any) (err error) {
	switch name {
	case "consumer", "Consumer":
		return cast2.String(value, &r.Consumer)
	case "createdAt", "CreatedAt":
		return cast2.Time(value, &r.CreatedAt)
	case "createdBy", "CreatedBy":
		return cast2.Uint64(value, &r.CreatedBy)
	case "deletedAt", "DeletedAt":
		return cast2.TimePtr(value, &r.DeletedAt)
	case "deletedBy", "DeletedBy":
		return cast2.Uint64(value, &r.DeletedBy)
	case "id", "ID":
		return cast2.Uint64(value, &r.ID)
	case "queue", "Queue":
		return cast2.String(value, &r.Queue)
	case "updatedAt", "UpdatedAt":
		return cast2.TimePtr(value, &r.UpdatedAt)
	case "updatedBy", "UpdatedBy":
		return cast2.Uint64(value, &r.UpdatedBy)

	}
	return nil
}

func (r Report) GetID() uint64 { return r.ID }

func (r *Report) GetValue(name string, pos uint) (any, error) {
	switch name {
	case "createdAt", "CreatedAt":
		return r.CreatedAt, nil
	case "createdBy", "CreatedBy":
		return r.CreatedBy, nil
	case "deletedAt", "DeletedAt":
		return r.DeletedAt, nil
	case "deletedBy", "DeletedBy":
		return r.DeletedBy, nil
	case "handle", "Handle":
		return r.Handle, nil
	case "id", "ID":
		return r.ID, nil
	case "ownedBy", "OwnedBy":
		return r.OwnedBy, nil
	case "updatedAt", "UpdatedAt":
		return r.UpdatedAt, nil
	case "updatedBy", "UpdatedBy":
		return r.UpdatedBy, nil

	}
	return nil, nil
}

func (r *Report) SetValue(name string, pos uint, value any) (err error) {
	switch name {
	case "createdAt", "CreatedAt":
		return cast2.Time(value, &r.CreatedAt)
	case "createdBy", "CreatedBy":
		return cast2.Uint64(value, &r.CreatedBy)
	case "deletedAt", "DeletedAt":
		return cast2.TimePtr(value, &r.DeletedAt)
	case "deletedBy", "DeletedBy":
		return cast2.Uint64(value, &r.DeletedBy)
	case "handle", "Handle":
		return cast2.String(value, &r.Handle)
	case "id", "ID":
		return cast2.Uint64(value, &r.ID)
	case "ownedBy", "OwnedBy":
		return cast2.Uint64(value, &r.OwnedBy)
	case "updatedAt", "UpdatedAt":
		return cast2.TimePtr(value, &r.UpdatedAt)
	case "updatedBy", "UpdatedBy":
		return cast2.Uint64(value, &r.UpdatedBy)

	}
	return nil
}

func (r ResourceTranslation) GetID() uint64 { return r.ID }

func (r *ResourceTranslation) GetValue(name string, pos uint) (any, error) {
	switch name {
	case "createdAt", "CreatedAt":
		return r.CreatedAt, nil
	case "createdBy", "CreatedBy":
		return r.CreatedBy, nil
	case "deletedAt", "DeletedAt":
		return r.DeletedAt, nil
	case "deletedBy", "DeletedBy":
		return r.DeletedBy, nil
	case "id", "ID":
		return r.ID, nil
	case "k", "K":
		return r.K, nil
	case "message", "Message":
		return r.Message, nil
	case "ownedBy", "OwnedBy":
		return r.OwnedBy, nil
	case "resource", "Resource":
		return r.Resource, nil
	case "updatedAt", "UpdatedAt":
		return r.UpdatedAt, nil
	case "updatedBy", "UpdatedBy":
		return r.UpdatedBy, nil

	}
	return nil, nil
}

func (r *ResourceTranslation) SetValue(name string, pos uint, value any) (err error) {
	switch name {
	case "createdAt", "CreatedAt":
		return cast2.Time(value, &r.CreatedAt)
	case "createdBy", "CreatedBy":
		return cast2.Uint64(value, &r.CreatedBy)
	case "deletedAt", "DeletedAt":
		return cast2.TimePtr(value, &r.DeletedAt)
	case "deletedBy", "DeletedBy":
		return cast2.Uint64(value, &r.DeletedBy)
	case "id", "ID":
		return cast2.Uint64(value, &r.ID)
	case "k", "K":
		return cast2.String(value, &r.K)
	case "message", "Message":
		return cast2.String(value, &r.Message)
	case "ownedBy", "OwnedBy":
		return cast2.Uint64(value, &r.OwnedBy)
	case "resource", "Resource":
		return cast2.String(value, &r.Resource)
	case "updatedAt", "UpdatedAt":
		return cast2.TimePtr(value, &r.UpdatedAt)
	case "updatedBy", "UpdatedBy":
		return cast2.Uint64(value, &r.UpdatedBy)

	default:
		return r.setValue(name, pos, value)

	}
	return nil
}

func (r Role) GetID() uint64 { return r.ID }

func (r *Role) GetValue(name string, pos uint) (any, error) {
	switch name {
	case "archivedAt", "ArchivedAt":
		return r.ArchivedAt, nil
	case "createdAt", "CreatedAt":
		return r.CreatedAt, nil
	case "deletedAt", "DeletedAt":
		return r.DeletedAt, nil
	case "handle", "Handle":
		return r.Handle, nil
	case "id", "ID":
		return r.ID, nil
	case "name", "Name":
		return r.Name, nil
	case "updatedAt", "UpdatedAt":
		return r.UpdatedAt, nil

	}
	return nil, nil
}

func (r *Role) SetValue(name string, pos uint, value any) (err error) {
	switch name {
	case "archivedAt", "ArchivedAt":
		return cast2.TimePtr(value, &r.ArchivedAt)
	case "createdAt", "CreatedAt":
		return cast2.Time(value, &r.CreatedAt)
	case "deletedAt", "DeletedAt":
		return cast2.TimePtr(value, &r.DeletedAt)
	case "handle", "Handle":
		return cast2.String(value, &r.Handle)
	case "id", "ID":
		return cast2.Uint64(value, &r.ID)
	case "name", "Name":
		return cast2.String(value, &r.Name)
	case "updatedAt", "UpdatedAt":
		return cast2.TimePtr(value, &r.UpdatedAt)

	}
	return nil
}

func (r RoleMember) GetID() uint64 {
	// The resource does not define an ID field
	return 0
}

func (r *RoleMember) GetValue(name string, pos uint) (any, error) {
	switch name {
	case "roleID", "RoleID":
		return r.RoleID, nil
	case "userID", "UserID":
		return r.UserID, nil

	}
	return nil, nil
}

func (r *RoleMember) SetValue(name string, pos uint, value any) (err error) {
	switch name {
	case "roleID", "RoleID":
		return cast2.Uint64(value, &r.RoleID)
	case "userID", "UserID":
		return cast2.Uint64(value, &r.UserID)

	}
	return nil
}

func (r Template) GetID() uint64 { return r.ID }

func (r *Template) GetValue(name string, pos uint) (any, error) {
	switch name {
	case "createdAt", "CreatedAt":
		return r.CreatedAt, nil
	case "deletedAt", "DeletedAt":
		return r.DeletedAt, nil
	case "handle", "Handle":
		return r.Handle, nil
	case "id", "ID":
		return r.ID, nil
	case "language", "Language":
		return r.Language, nil
	case "lastUsedAt", "LastUsedAt":
		return r.LastUsedAt, nil
	case "ownerID", "OwnerID":
		return r.OwnerID, nil
	case "partial", "Partial":
		return r.Partial, nil
	case "template", "Template":
		return r.Template, nil
	case "updatedAt", "UpdatedAt":
		return r.UpdatedAt, nil

	}
	return nil, nil
}

func (r *Template) SetValue(name string, pos uint, value any) (err error) {
	switch name {
	case "createdAt", "CreatedAt":
		return cast2.Time(value, &r.CreatedAt)
	case "deletedAt", "DeletedAt":
		return cast2.TimePtr(value, &r.DeletedAt)
	case "handle", "Handle":
		return cast2.String(value, &r.Handle)
	case "id", "ID":
		return cast2.Uint64(value, &r.ID)
	case "language", "Language":
		return cast2.String(value, &r.Language)
	case "lastUsedAt", "LastUsedAt":
		return cast2.TimePtr(value, &r.LastUsedAt)
	case "ownerID", "OwnerID":
		return cast2.Uint64(value, &r.OwnerID)
	case "partial", "Partial":
		return cast2.Bool(value, &r.Partial)
	case "template", "Template":
		return cast2.String(value, &r.Template)
	case "updatedAt", "UpdatedAt":
		return cast2.TimePtr(value, &r.UpdatedAt)

	}
	return nil
}

func (r User) GetID() uint64 { return r.ID }

func (r *User) GetValue(name string, pos uint) (any, error) {
	switch name {
	case "createdAt", "CreatedAt":
		return r.CreatedAt, nil
	case "deletedAt", "DeletedAt":
		return r.DeletedAt, nil
	case "email", "Email":
		return r.Email, nil
	case "emailConfirmed", "EmailConfirmed":
		return r.EmailConfirmed, nil
	case "handle", "Handle":
		return r.Handle, nil
	case "id", "ID":
		return r.ID, nil
	case "name", "Name":
		return r.Name, nil
	case "suspendedAt", "SuspendedAt":
		return r.SuspendedAt, nil
	case "updatedAt", "UpdatedAt":
		return r.UpdatedAt, nil
	case "username", "Username":
		return r.Username, nil

	}
	return nil, nil
}

func (r *User) SetValue(name string, pos uint, value any) (err error) {
	switch name {
	case "createdAt", "CreatedAt":
		return cast2.Time(value, &r.CreatedAt)
	case "deletedAt", "DeletedAt":
		return cast2.TimePtr(value, &r.DeletedAt)
	case "email", "Email":
		return cast2.String(value, &r.Email)
	case "emailConfirmed", "EmailConfirmed":
		return cast2.Bool(value, &r.EmailConfirmed)
	case "handle", "Handle":
		return cast2.String(value, &r.Handle)
	case "id", "ID":
		return cast2.Uint64(value, &r.ID)
	case "name", "Name":
		return cast2.String(value, &r.Name)
	case "suspendedAt", "SuspendedAt":
		return cast2.TimePtr(value, &r.SuspendedAt)
	case "updatedAt", "UpdatedAt":
		return cast2.TimePtr(value, &r.UpdatedAt)
	case "username", "Username":
		return cast2.String(value, &r.Username)

	}
	return nil
}

func (r DalConnection) GetID() uint64 { return r.ID }

func (r *DalConnection) GetValue(name string, pos uint) (any, error) {
	switch name {
	case "createdAt", "CreatedAt":
		return r.CreatedAt, nil
	case "createdBy", "CreatedBy":
		return r.CreatedBy, nil
	case "deletedAt", "DeletedAt":
		return r.DeletedAt, nil
	case "deletedBy", "DeletedBy":
		return r.DeletedBy, nil
	case "handle", "Handle":
		return r.Handle, nil
	case "id", "ID":
		return r.ID, nil
	case "type", "Type":
		return r.Type, nil
	case "updatedAt", "UpdatedAt":
		return r.UpdatedAt, nil
	case "updatedBy", "UpdatedBy":
		return r.UpdatedBy, nil

	}
	return nil, nil
}

func (r *DalConnection) SetValue(name string, pos uint, value any) (err error) {
	switch name {
	case "createdAt", "CreatedAt":
		return cast2.Time(value, &r.CreatedAt)
	case "createdBy", "CreatedBy":
		return cast2.Uint64(value, &r.CreatedBy)
	case "deletedAt", "DeletedAt":
		return cast2.TimePtr(value, &r.DeletedAt)
	case "deletedBy", "DeletedBy":
		return cast2.Uint64(value, &r.DeletedBy)
	case "handle", "Handle":
		return cast2.String(value, &r.Handle)
	case "id", "ID":
		return cast2.Uint64(value, &r.ID)
	case "type", "Type":
		return cast2.String(value, &r.Type)
	case "updatedAt", "UpdatedAt":
		return cast2.TimePtr(value, &r.UpdatedAt)
	case "updatedBy", "UpdatedBy":
		return cast2.Uint64(value, &r.UpdatedBy)

	}
	return nil
}

func (r DalSensitivityLevel) GetID() uint64 { return r.ID }

func (r *DalSensitivityLevel) GetValue(name string, pos uint) (any, error) {
	switch name {
	case "createdAt", "CreatedAt":
		return r.CreatedAt, nil
	case "createdBy", "CreatedBy":
		return r.CreatedBy, nil
	case "deletedAt", "DeletedAt":
		return r.DeletedAt, nil
	case "deletedBy", "DeletedBy":
		return r.DeletedBy, nil
	case "handle", "Handle":
		return r.Handle, nil
	case "id", "ID":
		return r.ID, nil
	case "level", "Level":
		return r.Level, nil
	case "updatedAt", "UpdatedAt":
		return r.UpdatedAt, nil
	case "updatedBy", "UpdatedBy":
		return r.UpdatedBy, nil

	}
	return nil, nil
}

func (r *DalSensitivityLevel) SetValue(name string, pos uint, value any) (err error) {
	switch name {
	case "createdAt", "CreatedAt":
		return cast2.Time(value, &r.CreatedAt)
	case "createdBy", "CreatedBy":
		return cast2.Uint64(value, &r.CreatedBy)
	case "deletedAt", "DeletedAt":
		return cast2.TimePtr(value, &r.DeletedAt)
	case "deletedBy", "DeletedBy":
		return cast2.Uint64(value, &r.DeletedBy)
	case "handle", "Handle":
		return cast2.String(value, &r.Handle)
	case "id", "ID":
		return cast2.Uint64(value, &r.ID)
	case "level", "Level":
		return cast2.Int(value, &r.Level)
	case "updatedAt", "UpdatedAt":
		return cast2.TimePtr(value, &r.UpdatedAt)
	case "updatedBy", "UpdatedBy":
		return cast2.Uint64(value, &r.UpdatedBy)

	}
	return nil
}

func (r DalSchemaAlteration) GetID() uint64 { return r.ID }

func (r *DalSchemaAlteration) GetValue(name string, pos uint) (any, error) {
	switch name {
	case "batchID", "BatchID":
		return r.BatchID, nil
	case "completedAt", "CompletedAt":
		return r.CompletedAt, nil
	case "completedBy", "CompletedBy":
		return r.CompletedBy, nil
	case "createdAt", "CreatedAt":
		return r.CreatedAt, nil
	case "createdBy", "CreatedBy":
		return r.CreatedBy, nil
	case "deletedAt", "DeletedAt":
		return r.DeletedAt, nil
	case "deletedBy", "DeletedBy":
		return r.DeletedBy, nil
	case "dependsOn", "DependsOn":
		return r.DependsOn, nil
	case "id", "ID":
		return r.ID, nil
	case "kind", "Kind":
		return r.Kind, nil
	case "updatedAt", "UpdatedAt":
		return r.UpdatedAt, nil
	case "updatedBy", "UpdatedBy":
		return r.UpdatedBy, nil

	}
	return nil, nil
}

func (r *DalSchemaAlteration) SetValue(name string, pos uint, value any) (err error) {
	switch name {
	case "batchID", "BatchID":
		return cast2.Uint64(value, &r.BatchID)
	case "completedAt", "CompletedAt":
		return cast2.TimePtr(value, &r.CompletedAt)
	case "completedBy", "CompletedBy":
		return cast2.Uint64(value, &r.CompletedBy)
	case "createdAt", "CreatedAt":
		return cast2.Time(value, &r.CreatedAt)
	case "createdBy", "CreatedBy":
		return cast2.Uint64(value, &r.CreatedBy)
	case "deletedAt", "DeletedAt":
		return cast2.TimePtr(value, &r.DeletedAt)
	case "deletedBy", "DeletedBy":
		return cast2.Uint64(value, &r.DeletedBy)
	case "dependsOn", "DependsOn":
		return cast2.Uint64(value, &r.DependsOn)
	case "id", "ID":
		return cast2.Uint64(value, &r.ID)
	case "kind", "Kind":
		return cast2.String(value, &r.Kind)
	case "updatedAt", "UpdatedAt":
		return cast2.TimePtr(value, &r.UpdatedAt)
	case "updatedBy", "UpdatedBy":
		return cast2.Uint64(value, &r.UpdatedBy)

	}
	return nil
}
