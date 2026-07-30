package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"github.com/crusttech/human/server/pkg/cast2"
)

func (r Attachment) GetID() uint64 { return r.ID }

func (r *Attachment) GetValue(name string, pos uint) (any, error) {
	if r == nil {
		return nil, nil
	}

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
	case "projectID", "ProjectID":
		return r.ProjectID, nil
	case "tenantID", "TenantID":
		return r.TenantID, nil
	case "updatedAt", "UpdatedAt":
		return r.UpdatedAt, nil
	case "url", "Url":
		return r.Url, nil

	}
	return nil, nil
}

func (r *Attachment) SetValue(name string, pos uint, value any) (err error) {
	if r == nil {
		r = &Attachment{}
	}

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
	case "projectID", "ProjectID":
		return cast2.Uint64(value, &r.ProjectID)
	case "tenantID", "TenantID":
		return cast2.Uint64(value, &r.TenantID)
	case "updatedAt", "UpdatedAt":
		return cast2.TimePtr(value, &r.UpdatedAt)
	case "url", "Url":
		return cast2.String(value, &r.Url)

	}
	return nil
}

func (r Application) GetID() uint64 { return r.ID }

func (r *Application) GetValue(name string, pos uint) (any, error) {
	if r == nil {
		return nil, nil
	}

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
	case "projectID", "ProjectID":
		return r.ProjectID, nil
	case "tenantID", "TenantID":
		return r.TenantID, nil
	case "updatedAt", "UpdatedAt":
		return r.UpdatedAt, nil
	case "weight", "Weight":
		return r.Weight, nil

	}
	return nil, nil
}

func (r *Application) SetValue(name string, pos uint, value any) (err error) {
	if r == nil {
		r = &Application{}
	}

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
	case "projectID", "ProjectID":
		return cast2.Uint64(value, &r.ProjectID)
	case "tenantID", "TenantID":
		return cast2.Uint64(value, &r.TenantID)
	case "updatedAt", "UpdatedAt":
		return cast2.TimePtr(value, &r.UpdatedAt)
	case "weight", "Weight":
		return cast2.Int(value, &r.Weight)

	}
	return nil
}

func (r ApigwRoute) GetID() uint64 { return r.ID }

func (r *ApigwRoute) GetValue(name string, pos uint) (any, error) {
	if r == nil {
		return nil, nil
	}

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
	case "projectID", "ProjectID":
		return r.ProjectID, nil
	case "tenantID", "TenantID":
		return r.TenantID, nil
	case "updatedAt", "UpdatedAt":
		return r.UpdatedAt, nil
	case "updatedBy", "UpdatedBy":
		return r.UpdatedBy, nil

	}
	return nil, nil
}

func (r *ApigwRoute) SetValue(name string, pos uint, value any) (err error) {
	if r == nil {
		r = &ApigwRoute{}
	}

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
	case "projectID", "ProjectID":
		return cast2.Uint64(value, &r.ProjectID)
	case "tenantID", "TenantID":
		return cast2.Uint64(value, &r.TenantID)
	case "updatedAt", "UpdatedAt":
		return cast2.TimePtr(value, &r.UpdatedAt)
	case "updatedBy", "UpdatedBy":
		return cast2.Uint64(value, &r.UpdatedBy)

	}
	return nil
}

func (r ApigwFilter) GetID() uint64 { return r.ID }

func (r *ApigwFilter) GetValue(name string, pos uint) (any, error) {
	if r == nil {
		return nil, nil
	}

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
	case "projectID", "ProjectID":
		return r.ProjectID, nil
	case "ref", "Ref":
		return r.Ref, nil
	case "route", "Route", "ApigwRouteID":
		return r.Route, nil
	case "tenantID", "TenantID":
		return r.TenantID, nil
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
	if r == nil {
		r = &ApigwFilter{}
	}

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
	case "projectID", "ProjectID":
		return cast2.Uint64(value, &r.ProjectID)
	case "ref", "Ref":
		return cast2.String(value, &r.Ref)
	case "route", "Route", "ApigwRouteID":
		return cast2.Uint64(value, &r.Route)
	case "tenantID", "TenantID":
		return cast2.Uint64(value, &r.TenantID)
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
	if r == nil {
		return nil, nil
	}

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
	case "projectID", "ProjectID":
		return r.ProjectID, nil
	case "redirectURI", "RedirectURI":
		return r.RedirectURI, nil
	case "scope", "Scope":
		return r.Scope, nil
	case "secret", "Secret":
		return r.Secret, nil
	case "tenantID", "TenantID":
		return r.TenantID, nil
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
	if r == nil {
		r = &AuthClient{}
	}

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
	case "projectID", "ProjectID":
		return cast2.Uint64(value, &r.ProjectID)
	case "redirectURI", "RedirectURI":
		return cast2.String(value, &r.RedirectURI)
	case "scope", "Scope":
		return cast2.String(value, &r.Scope)
	case "secret", "Secret":
		return cast2.String(value, &r.Secret)
	case "tenantID", "TenantID":
		return cast2.Uint64(value, &r.TenantID)
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
	if r == nil {
		return nil, nil
	}

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
	if r == nil {
		r = &DataPrivacyRequestComment{}
	}

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
	if r == nil {
		return nil, nil
	}

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
	case "projectID", "ProjectID":
		return r.ProjectID, nil
	case "queue", "Queue":
		return r.Queue, nil
	case "tenantID", "TenantID":
		return r.TenantID, nil
	case "updatedAt", "UpdatedAt":
		return r.UpdatedAt, nil
	case "updatedBy", "UpdatedBy":
		return r.UpdatedBy, nil

	}
	return nil, nil
}

func (r *Queue) SetValue(name string, pos uint, value any) (err error) {
	if r == nil {
		r = &Queue{}
	}

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
	case "projectID", "ProjectID":
		return cast2.Uint64(value, &r.ProjectID)
	case "queue", "Queue":
		return cast2.String(value, &r.Queue)
	case "tenantID", "TenantID":
		return cast2.Uint64(value, &r.TenantID)
	case "updatedAt", "UpdatedAt":
		return cast2.TimePtr(value, &r.UpdatedAt)
	case "updatedBy", "UpdatedBy":
		return cast2.Uint64(value, &r.UpdatedBy)

	}
	return nil
}

func (r Report) GetID() uint64 { return r.ID }

func (r *Report) GetValue(name string, pos uint) (any, error) {
	if r == nil {
		return nil, nil
	}

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
	case "projectID", "ProjectID":
		return r.ProjectID, nil
	case "tenantID", "TenantID":
		return r.TenantID, nil
	case "updatedAt", "UpdatedAt":
		return r.UpdatedAt, nil
	case "updatedBy", "UpdatedBy":
		return r.UpdatedBy, nil

	}
	return nil, nil
}

func (r *Report) SetValue(name string, pos uint, value any) (err error) {
	if r == nil {
		r = &Report{}
	}

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
	case "projectID", "ProjectID":
		return cast2.Uint64(value, &r.ProjectID)
	case "tenantID", "TenantID":
		return cast2.Uint64(value, &r.TenantID)
	case "updatedAt", "UpdatedAt":
		return cast2.TimePtr(value, &r.UpdatedAt)
	case "updatedBy", "UpdatedBy":
		return cast2.Uint64(value, &r.UpdatedBy)

	}
	return nil
}

func (r ResourceTranslation) GetID() uint64 { return r.ID }

func (r *ResourceTranslation) GetValue(name string, pos uint) (any, error) {
	if r == nil {
		return nil, nil
	}

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
	case "projectID", "ProjectID":
		return r.ProjectID, nil
	case "resource", "Resource":
		return r.Resource, nil
	case "tenantID", "TenantID":
		return r.TenantID, nil
	case "updatedAt", "UpdatedAt":
		return r.UpdatedAt, nil
	case "updatedBy", "UpdatedBy":
		return r.UpdatedBy, nil

	}
	return nil, nil
}

func (r *ResourceTranslation) SetValue(name string, pos uint, value any) (err error) {
	if r == nil {
		r = &ResourceTranslation{}
	}

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
	case "projectID", "ProjectID":
		return cast2.Uint64(value, &r.ProjectID)
	case "resource", "Resource":
		return cast2.String(value, &r.Resource)
	case "tenantID", "TenantID":
		return cast2.Uint64(value, &r.TenantID)
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
	if r == nil {
		return nil, nil
	}

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
	case "projectID", "ProjectID":
		return r.ProjectID, nil
	case "tenantID", "TenantID":
		return r.TenantID, nil
	case "updatedAt", "UpdatedAt":
		return r.UpdatedAt, nil

	}
	return nil, nil
}

func (r *Role) SetValue(name string, pos uint, value any) (err error) {
	if r == nil {
		r = &Role{}
	}

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
	case "projectID", "ProjectID":
		return cast2.Uint64(value, &r.ProjectID)
	case "tenantID", "TenantID":
		return cast2.Uint64(value, &r.TenantID)
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
	if r == nil {
		return nil, nil
	}

	switch name {
	case "resource", "Resource":
		return r.Resource, nil
	case "roleID", "RoleID":
		return r.RoleID, nil

	}
	return nil, nil
}

func (r *RoleMember) SetValue(name string, pos uint, value any) (err error) {
	if r == nil {
		r = &RoleMember{}
	}

	switch name {
	case "resource", "Resource":
		return cast2.String(value, &r.Resource)
	case "roleID", "RoleID":
		return cast2.Uint64(value, &r.RoleID)

	}
	return nil
}

func (r UserGroup) GetID() uint64 { return r.ID }

func (r *UserGroup) GetValue(name string, pos uint) (any, error) {
	if r == nil {
		return nil, nil
	}

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
	case "projectID", "ProjectID":
		return r.ProjectID, nil
	case "tenantID", "TenantID":
		return r.TenantID, nil
	case "updatedAt", "UpdatedAt":
		return r.UpdatedAt, nil

	}
	return nil, nil
}

func (r *UserGroup) SetValue(name string, pos uint, value any) (err error) {
	if r == nil {
		r = &UserGroup{}
	}

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
	case "projectID", "ProjectID":
		return cast2.Uint64(value, &r.ProjectID)
	case "tenantID", "TenantID":
		return cast2.Uint64(value, &r.TenantID)
	case "updatedAt", "UpdatedAt":
		return cast2.TimePtr(value, &r.UpdatedAt)

	}
	return nil
}

func (r Template) GetID() uint64 { return r.ID }

func (r *Template) GetValue(name string, pos uint) (any, error) {
	if r == nil {
		return nil, nil
	}

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
	case "projectID", "ProjectID":
		return r.ProjectID, nil
	case "template", "Template":
		return r.Template, nil
	case "tenantID", "TenantID":
		return r.TenantID, nil
	case "updatedAt", "UpdatedAt":
		return r.UpdatedAt, nil

	}
	return nil, nil
}

func (r *Template) SetValue(name string, pos uint, value any) (err error) {
	if r == nil {
		r = &Template{}
	}

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
	case "projectID", "ProjectID":
		return cast2.Uint64(value, &r.ProjectID)
	case "template", "Template":
		return cast2.String(value, &r.Template)
	case "tenantID", "TenantID":
		return cast2.Uint64(value, &r.TenantID)
	case "updatedAt", "UpdatedAt":
		return cast2.TimePtr(value, &r.UpdatedAt)

	}
	return nil
}

func (r User) GetID() uint64 { return r.ID }

func (r *User) GetValue(name string, pos uint) (any, error) {
	if r == nil {
		return nil, nil
	}

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
	case "projectID", "ProjectID":
		return r.ProjectID, nil
	case "suspendedAt", "SuspendedAt":
		return r.SuspendedAt, nil
	case "tenantID", "TenantID":
		return r.TenantID, nil
	case "updatedAt", "UpdatedAt":
		return r.UpdatedAt, nil
	case "userGroupID", "UserGroupID":
		return r.UserGroupID, nil
	case "username", "Username":
		return r.Username, nil

	}
	return nil, nil
}

func (r *User) SetValue(name string, pos uint, value any) (err error) {
	if r == nil {
		r = &User{}
	}

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
	case "projectID", "ProjectID":
		return cast2.Uint64(value, &r.ProjectID)
	case "suspendedAt", "SuspendedAt":
		return cast2.TimePtr(value, &r.SuspendedAt)
	case "tenantID", "TenantID":
		return cast2.Uint64(value, &r.TenantID)
	case "updatedAt", "UpdatedAt":
		return cast2.TimePtr(value, &r.UpdatedAt)
	case "userGroupID", "UserGroupID":
		return cast2.Uint64(value, &r.UserGroupID)
	case "username", "Username":
		return cast2.String(value, &r.Username)

	}
	return nil
}

func (r DalConnection) GetID() uint64 { return r.ID }

func (r *DalConnection) GetValue(name string, pos uint) (any, error) {
	if r == nil {
		return nil, nil
	}

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
	if r == nil {
		r = &DalConnection{}
	}

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
	if r == nil {
		return nil, nil
	}

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
	case "projectID", "ProjectID":
		return r.ProjectID, nil
	case "tenantID", "TenantID":
		return r.TenantID, nil
	case "updatedAt", "UpdatedAt":
		return r.UpdatedAt, nil
	case "updatedBy", "UpdatedBy":
		return r.UpdatedBy, nil

	}
	return nil, nil
}

func (r *DalSensitivityLevel) SetValue(name string, pos uint, value any) (err error) {
	if r == nil {
		r = &DalSensitivityLevel{}
	}

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
	case "projectID", "ProjectID":
		return cast2.Uint64(value, &r.ProjectID)
	case "tenantID", "TenantID":
		return cast2.Uint64(value, &r.TenantID)
	case "updatedAt", "UpdatedAt":
		return cast2.TimePtr(value, &r.UpdatedAt)
	case "updatedBy", "UpdatedBy":
		return cast2.Uint64(value, &r.UpdatedBy)

	}
	return nil
}

func (r DalSchemaAlteration) GetID() uint64 { return r.ID }

func (r *DalSchemaAlteration) GetValue(name string, pos uint) (any, error) {
	if r == nil {
		return nil, nil
	}

	switch name {
	case "batchID", "BatchID":
		return r.BatchID, nil
	case "completedAt", "CompletedAt":
		return r.CompletedAt, nil
	case "completedBy", "CompletedBy":
		return r.CompletedBy, nil
	case "connectionID", "ConnectionID":
		return r.ConnectionID, nil
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
	case "dismissedAt", "DismissedAt":
		return r.DismissedAt, nil
	case "dismissedBy", "DismissedBy":
		return r.DismissedBy, nil
	case "error", "Error":
		return r.Error, nil
	case "id", "ID":
		return r.ID, nil
	case "kind", "Kind":
		return r.Kind, nil
	case "resource", "Resource":
		return r.Resource, nil
	case "resourceType", "ResourceType":
		return r.ResourceType, nil
	case "updatedAt", "UpdatedAt":
		return r.UpdatedAt, nil
	case "updatedBy", "UpdatedBy":
		return r.UpdatedBy, nil

	}
	return nil, nil
}

func (r *DalSchemaAlteration) SetValue(name string, pos uint, value any) (err error) {
	if r == nil {
		r = &DalSchemaAlteration{}
	}

	switch name {
	case "batchID", "BatchID":
		return cast2.Uint64(value, &r.BatchID)
	case "completedAt", "CompletedAt":
		return cast2.TimePtr(value, &r.CompletedAt)
	case "completedBy", "CompletedBy":
		return cast2.Uint64(value, &r.CompletedBy)
	case "connectionID", "ConnectionID":
		return cast2.Uint64(value, &r.ConnectionID)
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
	case "dismissedAt", "DismissedAt":
		return cast2.TimePtr(value, &r.DismissedAt)
	case "dismissedBy", "DismissedBy":
		return cast2.Uint64(value, &r.DismissedBy)
	case "error", "Error":
		return cast2.String(value, &r.Error)
	case "id", "ID":
		return cast2.Uint64(value, &r.ID)
	case "kind", "Kind":
		return cast2.String(value, &r.Kind)
	case "resource", "Resource":
		return cast2.String(value, &r.Resource)
	case "resourceType", "ResourceType":
		return cast2.String(value, &r.ResourceType)
	case "updatedAt", "UpdatedAt":
		return cast2.TimePtr(value, &r.UpdatedAt)
	case "updatedBy", "UpdatedBy":
		return cast2.Uint64(value, &r.UpdatedBy)

	}
	return nil
}

func (r Connection) GetID() uint64 { return r.ID }

func (r *Connection) GetValue(name string, pos uint) (any, error) {
	if r == nil {
		return nil, nil
	}

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
	case "revision", "Revision":
		return r.Revision, nil
	case "source", "Source":
		return r.Source, nil
	case "status", "Status":
		return r.Status, nil
	case "updatedAt", "UpdatedAt":
		return r.UpdatedAt, nil
	case "updatedBy", "UpdatedBy":
		return r.UpdatedBy, nil

	}
	return nil, nil
}

func (r *Connection) SetValue(name string, pos uint, value any) (err error) {
	if r == nil {
		r = &Connection{}
	}

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
	case "revision", "Revision":
		return cast2.Int(value, &r.Revision)
	case "source", "Source":
		return cast2.String(value, &r.Source)
	case "status", "Status":
		return cast2.String(value, &r.Status)
	case "updatedAt", "UpdatedAt":
		return cast2.TimePtr(value, &r.UpdatedAt)
	case "updatedBy", "UpdatedBy":
		return cast2.Uint64(value, &r.UpdatedBy)

	}
	return nil
}

func (r ConfiguredConnection) GetID() uint64 { return r.ID }

func (r *ConfiguredConnection) GetValue(name string, pos uint) (any, error) {
	if r == nil {
		return nil, nil
	}

	switch name {
	case "connectionID", "ConnectionID":
		return r.ConnectionID, nil
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
	case "name", "Name":
		return r.Name, nil
	case "projectID", "ProjectID":
		return r.ProjectID, nil
	case "status", "Status":
		return r.Status, nil
	case "tenantID", "TenantID":
		return r.TenantID, nil
	case "updatedAt", "UpdatedAt":
		return r.UpdatedAt, nil
	case "updatedBy", "UpdatedBy":
		return r.UpdatedBy, nil

	}
	return nil, nil
}

func (r *ConfiguredConnection) SetValue(name string, pos uint, value any) (err error) {
	if r == nil {
		r = &ConfiguredConnection{}
	}

	switch name {
	case "connectionID", "ConnectionID":
		return cast2.Uint64(value, &r.ConnectionID)
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
	case "name", "Name":
		return cast2.String(value, &r.Name)
	case "projectID", "ProjectID":
		return cast2.Uint64(value, &r.ProjectID)
	case "status", "Status":
		return cast2.String(value, &r.Status)
	case "tenantID", "TenantID":
		return cast2.Uint64(value, &r.TenantID)
	case "updatedAt", "UpdatedAt":
		return cast2.TimePtr(value, &r.UpdatedAt)
	case "updatedBy", "UpdatedBy":
		return cast2.Uint64(value, &r.UpdatedBy)

	}
	return nil
}

func (r LlmProvider) GetID() uint64 { return r.ID }

func (r *LlmProvider) GetValue(name string, pos uint) (any, error) {
	if r == nil {
		return nil, nil
	}

	switch name {
	case "createdAt", "CreatedAt":
		return r.CreatedAt, nil
	case "createdBy", "CreatedBy":
		return r.CreatedBy, nil
	case "credentialID", "CredentialID":
		return r.CredentialID, nil
	case "deletedAt", "DeletedAt":
		return r.DeletedAt, nil
	case "deletedBy", "DeletedBy":
		return r.DeletedBy, nil
	case "handle", "Handle":
		return r.Handle, nil
	case "id", "ID":
		return r.ID, nil
	case "projectID", "ProjectID":
		return r.ProjectID, nil
	case "provider", "Provider":
		return r.Provider, nil
	case "status", "Status":
		return r.Status, nil
	case "tenantID", "TenantID":
		return r.TenantID, nil
	case "updatedAt", "UpdatedAt":
		return r.UpdatedAt, nil
	case "updatedBy", "UpdatedBy":
		return r.UpdatedBy, nil

	}
	return nil, nil
}

func (r *LlmProvider) SetValue(name string, pos uint, value any) (err error) {
	if r == nil {
		r = &LlmProvider{}
	}

	switch name {
	case "createdAt", "CreatedAt":
		return cast2.Time(value, &r.CreatedAt)
	case "createdBy", "CreatedBy":
		return cast2.Uint64(value, &r.CreatedBy)
	case "credentialID", "CredentialID":
		return cast2.Uint64(value, &r.CredentialID)
	case "deletedAt", "DeletedAt":
		return cast2.TimePtr(value, &r.DeletedAt)
	case "deletedBy", "DeletedBy":
		return cast2.Uint64(value, &r.DeletedBy)
	case "handle", "Handle":
		return cast2.String(value, &r.Handle)
	case "id", "ID":
		return cast2.Uint64(value, &r.ID)
	case "projectID", "ProjectID":
		return cast2.Uint64(value, &r.ProjectID)
	case "provider", "Provider":
		return cast2.String(value, &r.Provider)
	case "status", "Status":
		return cast2.String(value, &r.Status)
	case "tenantID", "TenantID":
		return cast2.Uint64(value, &r.TenantID)
	case "updatedAt", "UpdatedAt":
		return cast2.TimePtr(value, &r.UpdatedAt)
	case "updatedBy", "UpdatedBy":
		return cast2.Uint64(value, &r.UpdatedBy)

	}
	return nil
}

func (r Agent) GetID() uint64 { return r.ID }

func (r *Agent) GetValue(name string, pos uint) (any, error) {
	if r == nil {
		return nil, nil
	}

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
	case "projectID", "ProjectID":
		return r.ProjectID, nil
	case "revision", "Revision":
		return r.Revision, nil
	case "status", "Status":
		return r.Status, nil
	case "tenantID", "TenantID":
		return r.TenantID, nil
	case "updatedAt", "UpdatedAt":
		return r.UpdatedAt, nil
	case "updatedBy", "UpdatedBy":
		return r.UpdatedBy, nil

	}
	return nil, nil
}

func (r *Agent) SetValue(name string, pos uint, value any) (err error) {
	if r == nil {
		r = &Agent{}
	}

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
	case "projectID", "ProjectID":
		return cast2.Uint64(value, &r.ProjectID)
	case "revision", "Revision":
		return cast2.Int(value, &r.Revision)
	case "status", "Status":
		return cast2.String(value, &r.Status)
	case "tenantID", "TenantID":
		return cast2.Uint64(value, &r.TenantID)
	case "updatedAt", "UpdatedAt":
		return cast2.TimePtr(value, &r.UpdatedAt)
	case "updatedBy", "UpdatedBy":
		return cast2.Uint64(value, &r.UpdatedBy)

	}
	return nil
}

func (r AiConversation) GetID() uint64 { return r.ID }

func (r *AiConversation) GetValue(name string, pos uint) (any, error) {
	if r == nil {
		return nil, nil
	}

	switch name {
	case "agentID", "AgentID":
		return r.AgentID, nil
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
	case "projectID", "ProjectID":
		return r.ProjectID, nil
	case "tenantID", "TenantID":
		return r.TenantID, nil
	case "tokenCount", "TokenCount":
		return r.TokenCount, nil
	case "updatedAt", "UpdatedAt":
		return r.UpdatedAt, nil
	case "updatedBy", "UpdatedBy":
		return r.UpdatedBy, nil

	}
	return nil, nil
}

func (r *AiConversation) SetValue(name string, pos uint, value any) (err error) {
	if r == nil {
		r = &AiConversation{}
	}

	switch name {
	case "agentID", "AgentID":
		return cast2.Uint64(value, &r.AgentID)
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
	case "projectID", "ProjectID":
		return cast2.Uint64(value, &r.ProjectID)
	case "tenantID", "TenantID":
		return cast2.Uint64(value, &r.TenantID)
	case "tokenCount", "TokenCount":
		return cast2.Int(value, &r.TokenCount)
	case "updatedAt", "UpdatedAt":
		return cast2.TimePtr(value, &r.UpdatedAt)
	case "updatedBy", "UpdatedBy":
		return cast2.Uint64(value, &r.UpdatedBy)

	}
	return nil
}

func (r KnowledgeBase) GetID() uint64 { return r.ID }

func (r *KnowledgeBase) GetValue(name string, pos uint) (any, error) {
	if r == nil {
		return nil, nil
	}

	switch name {
	case "createdAt", "CreatedAt":
		return r.CreatedAt, nil
	case "createdBy", "CreatedBy":
		return r.CreatedBy, nil
	case "deletedAt", "DeletedAt":
		return r.DeletedAt, nil
	case "deletedBy", "DeletedBy":
		return r.DeletedBy, nil
	case "description", "Description":
		return r.Description, nil
	case "handle", "Handle":
		return r.Handle, nil
	case "id", "ID":
		return r.ID, nil
	case "projectID", "ProjectID":
		return r.ProjectID, nil
	case "tenantID", "TenantID":
		return r.TenantID, nil
	case "title", "Title":
		return r.Title, nil
	case "updatedAt", "UpdatedAt":
		return r.UpdatedAt, nil
	case "updatedBy", "UpdatedBy":
		return r.UpdatedBy, nil

	}
	return nil, nil
}

func (r *KnowledgeBase) SetValue(name string, pos uint, value any) (err error) {
	if r == nil {
		r = &KnowledgeBase{}
	}

	switch name {
	case "createdAt", "CreatedAt":
		return cast2.Time(value, &r.CreatedAt)
	case "createdBy", "CreatedBy":
		return cast2.Uint64(value, &r.CreatedBy)
	case "deletedAt", "DeletedAt":
		return cast2.TimePtr(value, &r.DeletedAt)
	case "deletedBy", "DeletedBy":
		return cast2.Uint64(value, &r.DeletedBy)
	case "description", "Description":
		return cast2.String(value, &r.Description)
	case "handle", "Handle":
		return cast2.String(value, &r.Handle)
	case "id", "ID":
		return cast2.Uint64(value, &r.ID)
	case "projectID", "ProjectID":
		return cast2.Uint64(value, &r.ProjectID)
	case "tenantID", "TenantID":
		return cast2.Uint64(value, &r.TenantID)
	case "title", "Title":
		return cast2.String(value, &r.Title)
	case "updatedAt", "UpdatedAt":
		return cast2.TimePtr(value, &r.UpdatedAt)
	case "updatedBy", "UpdatedBy":
		return cast2.Uint64(value, &r.UpdatedBy)

	}
	return nil
}

func (r Chatbot) GetID() uint64 { return r.ID }

func (r *Chatbot) GetValue(name string, pos uint) (any, error) {
	if r == nil {
		return nil, nil
	}

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
	case "handle", "Handle":
		return r.Handle, nil
	case "id", "ID":
		return r.ID, nil
	case "name", "Name":
		return r.Name, nil
	case "projectID", "ProjectID":
		return r.ProjectID, nil
	case "sessionTTL", "SessionTTL":
		return r.SessionTTL, nil
	case "tenantID", "TenantID":
		return r.TenantID, nil
	case "updatedAt", "UpdatedAt":
		return r.UpdatedAt, nil
	case "updatedBy", "UpdatedBy":
		return r.UpdatedBy, nil
	case "widgetKey", "WidgetKey":
		return r.WidgetKey, nil

	}
	return nil, nil
}

func (r *Chatbot) SetValue(name string, pos uint, value any) (err error) {
	if r == nil {
		r = &Chatbot{}
	}

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
	case "handle", "Handle":
		return cast2.String(value, &r.Handle)
	case "id", "ID":
		return cast2.Uint64(value, &r.ID)
	case "name", "Name":
		return cast2.String(value, &r.Name)
	case "projectID", "ProjectID":
		return cast2.Uint64(value, &r.ProjectID)
	case "sessionTTL", "SessionTTL":
		return cast2.String(value, &r.SessionTTL)
	case "tenantID", "TenantID":
		return cast2.Uint64(value, &r.TenantID)
	case "updatedAt", "UpdatedAt":
		return cast2.TimePtr(value, &r.UpdatedAt)
	case "updatedBy", "UpdatedBy":
		return cast2.Uint64(value, &r.UpdatedBy)
	case "widgetKey", "WidgetKey":
		return cast2.String(value, &r.WidgetKey)

	}
	return nil
}

func (r ChatbotSession) GetID() uint64 { return r.ID }

func (r *ChatbotSession) GetValue(name string, pos uint) (any, error) {
	if r == nil {
		return nil, nil
	}

	switch name {
	case "chatbotID", "ChatbotID":
		return r.ChatbotID, nil
	case "createdAt", "CreatedAt":
		return r.CreatedAt, nil
	case "createdBy", "CreatedBy":
		return r.CreatedBy, nil
	case "currentStep", "CurrentStep":
		return r.CurrentStep, nil
	case "deletedAt", "DeletedAt":
		return r.DeletedAt, nil
	case "deletedBy", "DeletedBy":
		return r.DeletedBy, nil
	case "id", "ID":
		return r.ID, nil
	case "projectID", "ProjectID":
		return r.ProjectID, nil
	case "status", "Status":
		return r.Status, nil
	case "tenantID", "TenantID":
		return r.TenantID, nil
	case "updatedAt", "UpdatedAt":
		return r.UpdatedAt, nil
	case "updatedBy", "UpdatedBy":
		return r.UpdatedBy, nil

	}
	return nil, nil
}

func (r *ChatbotSession) SetValue(name string, pos uint, value any) (err error) {
	if r == nil {
		r = &ChatbotSession{}
	}

	switch name {
	case "chatbotID", "ChatbotID":
		return cast2.Uint64(value, &r.ChatbotID)
	case "createdAt", "CreatedAt":
		return cast2.Time(value, &r.CreatedAt)
	case "createdBy", "CreatedBy":
		return cast2.Uint64(value, &r.CreatedBy)
	case "currentStep", "CurrentStep":
		return cast2.Int(value, &r.CurrentStep)
	case "deletedAt", "DeletedAt":
		return cast2.TimePtr(value, &r.DeletedAt)
	case "deletedBy", "DeletedBy":
		return cast2.Uint64(value, &r.DeletedBy)
	case "id", "ID":
		return cast2.Uint64(value, &r.ID)
	case "projectID", "ProjectID":
		return cast2.Uint64(value, &r.ProjectID)
	case "status", "Status":
		return cast2.String(value, &r.Status)
	case "tenantID", "TenantID":
		return cast2.Uint64(value, &r.TenantID)
	case "updatedAt", "UpdatedAt":
		return cast2.TimePtr(value, &r.UpdatedAt)
	case "updatedBy", "UpdatedBy":
		return cast2.Uint64(value, &r.UpdatedBy)

	}
	return nil
}

func (r ChatbotSessionStep) GetID() uint64 { return r.ID }

func (r *ChatbotSessionStep) GetValue(name string, pos uint) (any, error) {
	if r == nil {
		return nil, nil
	}

	switch name {
	case "conversationID", "ConversationID":
		return r.ConversationID, nil
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
	case "scenarioIndex", "ScenarioIndex":
		return r.ScenarioIndex, nil
	case "sessionID", "SessionID":
		return r.SessionID, nil
	case "status", "Status":
		return r.Status, nil
	case "updatedAt", "UpdatedAt":
		return r.UpdatedAt, nil
	case "updatedBy", "UpdatedBy":
		return r.UpdatedBy, nil

	}
	return nil, nil
}

func (r *ChatbotSessionStep) SetValue(name string, pos uint, value any) (err error) {
	if r == nil {
		r = &ChatbotSessionStep{}
	}

	switch name {
	case "conversationID", "ConversationID":
		return cast2.Uint64(value, &r.ConversationID)
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
	case "scenarioIndex", "ScenarioIndex":
		return cast2.Int(value, &r.ScenarioIndex)
	case "sessionID", "SessionID":
		return cast2.Uint64(value, &r.SessionID)
	case "status", "Status":
		return cast2.String(value, &r.Status)
	case "updatedAt", "UpdatedAt":
		return cast2.TimePtr(value, &r.UpdatedAt)
	case "updatedBy", "UpdatedBy":
		return cast2.Uint64(value, &r.UpdatedBy)

	}
	return nil
}

func (r ChatbotSessionHandoff) GetID() uint64 { return r.ID }

func (r *ChatbotSessionHandoff) GetValue(name string, pos uint) (any, error) {
	if r == nil {
		return nil, nil
	}

	switch name {
	case "closedAt", "ClosedAt":
		return r.ClosedAt, nil
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
	case "initiatedAt", "InitiatedAt":
		return r.InitiatedAt, nil
	case "sessionID", "SessionID":
		return r.SessionID, nil
	case "status", "Status":
		return r.Status, nil
	case "stepID", "StepID":
		return r.StepID, nil
	case "updatedAt", "UpdatedAt":
		return r.UpdatedAt, nil
	case "updatedBy", "UpdatedBy":
		return r.UpdatedBy, nil

	}
	return nil, nil
}

func (r *ChatbotSessionHandoff) SetValue(name string, pos uint, value any) (err error) {
	if r == nil {
		r = &ChatbotSessionHandoff{}
	}

	switch name {
	case "closedAt", "ClosedAt":
		return cast2.TimePtr(value, &r.ClosedAt)
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
	case "initiatedAt", "InitiatedAt":
		return cast2.Time(value, &r.InitiatedAt)
	case "sessionID", "SessionID":
		return cast2.Uint64(value, &r.SessionID)
	case "status", "Status":
		return cast2.String(value, &r.Status)
	case "stepID", "StepID":
		return cast2.Uint64(value, &r.StepID)
	case "updatedAt", "UpdatedAt":
		return cast2.TimePtr(value, &r.UpdatedAt)
	case "updatedBy", "UpdatedBy":
		return cast2.Uint64(value, &r.UpdatedBy)

	}
	return nil
}

func (r Tenant) GetID() uint64 { return r.ID }

func (r *Tenant) GetValue(name string, pos uint) (any, error) {
	if r == nil {
		return nil, nil
	}

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
	case "suspendedAt", "SuspendedAt":
		return r.SuspendedAt, nil
	case "updatedAt", "UpdatedAt":
		return r.UpdatedAt, nil
	case "updatedBy", "UpdatedBy":
		return r.UpdatedBy, nil

	}
	return nil, nil
}

func (r *Tenant) SetValue(name string, pos uint, value any) (err error) {
	if r == nil {
		r = &Tenant{}
	}

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
	case "suspendedAt", "SuspendedAt":
		return cast2.TimePtr(value, &r.SuspendedAt)
	case "updatedAt", "UpdatedAt":
		return cast2.TimePtr(value, &r.UpdatedAt)
	case "updatedBy", "UpdatedBy":
		return cast2.Uint64(value, &r.UpdatedBy)

	}
	return nil
}

func (r TenantMembership) GetID() uint64 { return r.ID }

func (r *TenantMembership) GetValue(name string, pos uint) (any, error) {
	if r == nil {
		return nil, nil
	}

	switch name {
	case "createdAt", "CreatedAt":
		return r.CreatedAt, nil
	case "id", "ID":
		return r.ID, nil
	case "invitedBy", "InvitedBy":
		return r.InvitedBy, nil
	case "projectID", "ProjectID":
		return r.ProjectID, nil
	case "tenantID", "TenantID":
		return r.TenantID, nil
	case "updatedAt", "UpdatedAt":
		return r.UpdatedAt, nil
	case "userID", "UserID":
		return r.UserID, nil

	}
	return nil, nil
}

func (r *TenantMembership) SetValue(name string, pos uint, value any) (err error) {
	if r == nil {
		r = &TenantMembership{}
	}

	switch name {
	case "createdAt", "CreatedAt":
		return cast2.Time(value, &r.CreatedAt)
	case "id", "ID":
		return cast2.Uint64(value, &r.ID)
	case "invitedBy", "InvitedBy":
		return cast2.Uint64(value, &r.InvitedBy)
	case "projectID", "ProjectID":
		return cast2.Uint64(value, &r.ProjectID)
	case "tenantID", "TenantID":
		return cast2.Uint64(value, &r.TenantID)
	case "updatedAt", "UpdatedAt":
		return cast2.TimePtr(value, &r.UpdatedAt)
	case "userID", "UserID":
		return cast2.Uint64(value, &r.UserID)

	}
	return nil
}

func (r Project) GetID() uint64 { return r.ID }

func (r *Project) GetValue(name string, pos uint) (any, error) {
	if r == nil {
		return nil, nil
	}

	switch name {
	case "archivedAt", "ArchivedAt":
		return r.ArchivedAt, nil
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
	case "tenantID", "TenantID":
		return r.TenantID, nil
	case "updatedAt", "UpdatedAt":
		return r.UpdatedAt, nil
	case "updatedBy", "UpdatedBy":
		return r.UpdatedBy, nil

	}
	return nil, nil
}

func (r *Project) SetValue(name string, pos uint, value any) (err error) {
	if r == nil {
		r = &Project{}
	}

	switch name {
	case "archivedAt", "ArchivedAt":
		return cast2.TimePtr(value, &r.ArchivedAt)
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
	case "tenantID", "TenantID":
		return cast2.Uint64(value, &r.TenantID)
	case "updatedAt", "UpdatedAt":
		return cast2.TimePtr(value, &r.UpdatedAt)
	case "updatedBy", "UpdatedBy":
		return cast2.Uint64(value, &r.UpdatedBy)

	}
	return nil
}

func (r ProjectMember) GetID() uint64 { return r.ID }

func (r *ProjectMember) GetValue(name string, pos uint) (any, error) {
	if r == nil {
		return nil, nil
	}

	switch name {
	case "createdAt", "CreatedAt":
		return r.CreatedAt, nil
	case "deletedAt", "DeletedAt":
		return r.DeletedAt, nil
	case "id", "ID":
		return r.ID, nil
	case "invitedBy", "InvitedBy":
		return r.InvitedBy, nil
	case "projectID", "ProjectID":
		return r.ProjectID, nil
	case "tenantID", "TenantID":
		return r.TenantID, nil
	case "updatedAt", "UpdatedAt":
		return r.UpdatedAt, nil
	case "userID", "UserID":
		return r.UserID, nil

	}
	return nil, nil
}

func (r *ProjectMember) SetValue(name string, pos uint, value any) (err error) {
	if r == nil {
		r = &ProjectMember{}
	}

	switch name {
	case "createdAt", "CreatedAt":
		return cast2.Time(value, &r.CreatedAt)
	case "deletedAt", "DeletedAt":
		return cast2.TimePtr(value, &r.DeletedAt)
	case "id", "ID":
		return cast2.Uint64(value, &r.ID)
	case "invitedBy", "InvitedBy":
		return cast2.Uint64(value, &r.InvitedBy)
	case "projectID", "ProjectID":
		return cast2.Uint64(value, &r.ProjectID)
	case "tenantID", "TenantID":
		return cast2.Uint64(value, &r.TenantID)
	case "updatedAt", "UpdatedAt":
		return cast2.TimePtr(value, &r.UpdatedAt)
	case "userID", "UserID":
		return cast2.Uint64(value, &r.UserID)

	}
	return nil
}

func (r ProjectAiSystem) GetID() uint64 { return r.ID }

func (r *ProjectAiSystem) GetValue(name string, pos uint) (any, error) {
	if r == nil {
		return nil, nil
	}

	switch name {
	case "createdAt", "CreatedAt":
		return r.CreatedAt, nil
	case "deletedAt", "DeletedAt":
		return r.DeletedAt, nil
	case "handle", "Handle":
		return r.Handle, nil
	case "id", "ID":
		return r.ID, nil
	case "projectID", "ProjectID":
		return r.ProjectID, nil
	case "riskClass", "RiskClass":
		return r.RiskClass, nil
	case "tenantID", "TenantID":
		return r.TenantID, nil
	case "updatedAt", "UpdatedAt":
		return r.UpdatedAt, nil

	}
	return nil, nil
}

func (r *ProjectAiSystem) SetValue(name string, pos uint, value any) (err error) {
	if r == nil {
		r = &ProjectAiSystem{}
	}

	switch name {
	case "createdAt", "CreatedAt":
		return cast2.Time(value, &r.CreatedAt)
	case "deletedAt", "DeletedAt":
		return cast2.TimePtr(value, &r.DeletedAt)
	case "handle", "Handle":
		return cast2.String(value, &r.Handle)
	case "id", "ID":
		return cast2.Uint64(value, &r.ID)
	case "projectID", "ProjectID":
		return cast2.Uint64(value, &r.ProjectID)
	case "riskClass", "RiskClass":
		return cast2.String(value, &r.RiskClass)
	case "tenantID", "TenantID":
		return cast2.Uint64(value, &r.TenantID)
	case "updatedAt", "UpdatedAt":
		return cast2.TimePtr(value, &r.UpdatedAt)

	}
	return nil
}

func (r ProjectAiSystemEntry) GetID() uint64 { return r.ID }

func (r *ProjectAiSystemEntry) GetValue(name string, pos uint) (any, error) {
	if r == nil {
		return nil, nil
	}

	switch name {
	case "createdAt", "CreatedAt":
		return r.CreatedAt, nil
	case "projectAiSystemID", "ProjectAiSystemID":
		return r.ProjectAiSystemID, nil
	case "resourceRef", "ResourceRef":
		return r.ResourceRef, nil

	}
	return nil, nil
}

func (r *ProjectAiSystemEntry) SetValue(name string, pos uint, value any) (err error) {
	if r == nil {
		r = &ProjectAiSystemEntry{}
	}

	switch name {
	case "createdAt", "CreatedAt":
		return cast2.Time(value, &r.CreatedAt)
	case "projectAiSystemID", "ProjectAiSystemID":
		return cast2.Uint64(value, &r.ProjectAiSystemID)
	case "resourceRef", "ResourceRef":
		return cast2.String(value, &r.ResourceRef)

	}
	return nil
}

func (r ProjectFriaScenario) GetID() uint64 { return r.ID }

func (r *ProjectFriaScenario) GetValue(name string, pos uint) (any, error) {
	if r == nil {
		return nil, nil
	}

	switch name {
	case "aiSystemID", "AiSystemID":
		return r.AiSystemID, nil
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
	case "projectID", "ProjectID":
		return r.ProjectID, nil
	case "severity", "Severity":
		return r.Severity, nil
	case "tenantID", "TenantID":
		return r.TenantID, nil
	case "title", "Title":
		return r.Title, nil
	case "updatedAt", "UpdatedAt":
		return r.UpdatedAt, nil
	case "updatedBy", "UpdatedBy":
		return r.UpdatedBy, nil

	}
	return nil, nil
}

func (r *ProjectFriaScenario) SetValue(name string, pos uint, value any) (err error) {
	if r == nil {
		r = &ProjectFriaScenario{}
	}

	switch name {
	case "aiSystemID", "AiSystemID":
		return cast2.Uint64(value, &r.AiSystemID)
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
	case "projectID", "ProjectID":
		return cast2.Uint64(value, &r.ProjectID)
	case "severity", "Severity":
		return cast2.String(value, &r.Severity)
	case "tenantID", "TenantID":
		return cast2.Uint64(value, &r.TenantID)
	case "title", "Title":
		return cast2.String(value, &r.Title)
	case "updatedAt", "UpdatedAt":
		return cast2.TimePtr(value, &r.UpdatedAt)
	case "updatedBy", "UpdatedBy":
		return cast2.Uint64(value, &r.UpdatedBy)

	}
	return nil
}

func (r DmlConnection) GetID() uint64 { return r.ID }

func (r *DmlConnection) GetValue(name string, pos uint) (any, error) {
	if r == nil {
		return nil, nil
	}

	switch name {
	case "createdAt", "CreatedAt":
		return r.CreatedAt, nil
	case "deletedAt", "DeletedAt":
		return r.DeletedAt, nil
	case "handle", "Handle":
		return r.Handle, nil
	case "id", "ID":
		return r.ID, nil
	case "label", "Label":
		return r.Label, nil
	case "updatedAt", "UpdatedAt":
		return r.UpdatedAt, nil

	}
	return nil, nil
}

func (r *DmlConnection) SetValue(name string, pos uint, value any) (err error) {
	if r == nil {
		r = &DmlConnection{}
	}

	switch name {
	case "createdAt", "CreatedAt":
		return cast2.Time(value, &r.CreatedAt)
	case "deletedAt", "DeletedAt":
		return cast2.TimePtr(value, &r.DeletedAt)
	case "handle", "Handle":
		return cast2.String(value, &r.Handle)
	case "id", "ID":
		return cast2.Uint64(value, &r.ID)
	case "label", "Label":
		return cast2.String(value, &r.Label)
	case "updatedAt", "UpdatedAt":
		return cast2.TimePtr(value, &r.UpdatedAt)

	}
	return nil
}

func (r DmlMapping) GetID() uint64 { return r.ID }

func (r *DmlMapping) GetValue(name string, pos uint) (any, error) {
	if r == nil {
		return nil, nil
	}

	switch name {
	case "connectionID", "ConnectionID":
		return r.ConnectionID, nil
	case "createdAt", "CreatedAt":
		return r.CreatedAt, nil
	case "deletedAt", "DeletedAt":
		return r.DeletedAt, nil
	case "id", "ID":
		return r.ID, nil
	case "identifier", "Identifier":
		return r.Identifier, nil
	case "moduleHandle", "ModuleHandle":
		return r.ModuleHandle, nil
	case "moduleName", "ModuleName":
		return r.ModuleName, nil
	case "namespaceHandle", "NamespaceHandle":
		return r.NamespaceHandle, nil
	case "skip", "Skip":
		return r.Skip, nil
	case "sourceIdent", "SourceIdent":
		return r.SourceIdent, nil
	case "updatedAt", "UpdatedAt":
		return r.UpdatedAt, nil

	}
	return nil, nil
}

func (r *DmlMapping) SetValue(name string, pos uint, value any) (err error) {
	if r == nil {
		r = &DmlMapping{}
	}

	switch name {
	case "connectionID", "ConnectionID":
		return cast2.Uint64(value, &r.ConnectionID)
	case "createdAt", "CreatedAt":
		return cast2.Time(value, &r.CreatedAt)
	case "deletedAt", "DeletedAt":
		return cast2.TimePtr(value, &r.DeletedAt)
	case "id", "ID":
		return cast2.Uint64(value, &r.ID)
	case "identifier", "Identifier":
		return cast2.String(value, &r.Identifier)
	case "moduleHandle", "ModuleHandle":
		return cast2.String(value, &r.ModuleHandle)
	case "moduleName", "ModuleName":
		return cast2.String(value, &r.ModuleName)
	case "namespaceHandle", "NamespaceHandle":
		return cast2.String(value, &r.NamespaceHandle)
	case "skip", "Skip":
		return cast2.Bool(value, &r.Skip)
	case "sourceIdent", "SourceIdent":
		return cast2.String(value, &r.SourceIdent)
	case "updatedAt", "UpdatedAt":
		return cast2.TimePtr(value, &r.UpdatedAt)

	}
	return nil
}

func (r DmlImportRun) GetID() uint64 { return r.ID }

func (r *DmlImportRun) GetValue(name string, pos uint) (any, error) {
	if r == nil {
		return nil, nil
	}

	switch name {
	case "connectionID", "ConnectionID":
		return r.ConnectionID, nil
	case "failed", "Failed":
		return r.Failed, nil
	case "id", "ID":
		return r.ID, nil
	case "mappingID", "MappingID":
		return r.MappingID, nil
	case "processed", "Processed":
		return r.Processed, nil
	case "status", "Status":
		return r.Status, nil

	}
	return nil, nil
}

func (r *DmlImportRun) SetValue(name string, pos uint, value any) (err error) {
	if r == nil {
		r = &DmlImportRun{}
	}

	switch name {
	case "connectionID", "ConnectionID":
		return cast2.Uint64(value, &r.ConnectionID)
	case "failed", "Failed":
		return cast2.Uint64(value, &r.Failed)
	case "id", "ID":
		return cast2.Uint64(value, &r.ID)
	case "mappingID", "MappingID":
		return cast2.Uint64(value, &r.MappingID)
	case "processed", "Processed":
		return cast2.Uint64(value, &r.Processed)
	case "status", "Status":
		return cast2.String(value, &r.Status)

	}
	return nil
}
