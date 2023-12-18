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
	case "namespaceID", "NamespaceID":
		return r.NamespaceID, nil
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
	case "namespaceID", "NamespaceID":
		return cast2.Uint64(value, &r.NamespaceID)
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

func (r Chart) GetID() uint64 { return r.ID }

func (r *Chart) GetValue(name string, pos uint) (any, error) {
	switch name {
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
	case "namespaceID", "NamespaceID":
		return r.NamespaceID, nil
	case "updatedAt", "UpdatedAt":
		return r.UpdatedAt, nil

	}
	return nil, nil
}

func (r *Chart) SetValue(name string, pos uint, value any) (err error) {
	switch name {
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
	case "namespaceID", "NamespaceID":
		return cast2.Uint64(value, &r.NamespaceID)
	case "updatedAt", "UpdatedAt":
		return cast2.TimePtr(value, &r.UpdatedAt)

	default:
		return r.setValue(name, pos, value)

	}
	return nil
}

func (r Module) GetID() uint64 { return r.ID }

func (r *Module) GetValue(name string, pos uint) (any, error) {
	switch name {
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
	case "namespaceID", "NamespaceID":
		return r.NamespaceID, nil
	case "updatedAt", "UpdatedAt":
		return r.UpdatedAt, nil

	}
	return nil, nil
}

func (r *Module) SetValue(name string, pos uint, value any) (err error) {
	switch name {
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
	case "namespaceID", "NamespaceID":
		return cast2.Uint64(value, &r.NamespaceID)
	case "updatedAt", "UpdatedAt":
		return cast2.TimePtr(value, &r.UpdatedAt)

	}
	return nil
}

func (r ModuleField) GetID() uint64 { return r.ID }

func (r *ModuleField) GetValue(name string, pos uint) (any, error) {
	switch name {
	case "createdAt", "CreatedAt":
		return r.CreatedAt, nil
	case "deletedAt", "DeletedAt":
		return r.DeletedAt, nil
	case "id", "ID":
		return r.ID, nil
	case "kind", "Kind":
		return r.Kind, nil
	case "label", "Label":
		return r.Label, nil
	case "moduleID", "ModuleID":
		return r.ModuleID, nil
	case "multi", "Multi":
		return r.Multi, nil
	case "name", "Name":
		return r.Name, nil
	case "place", "Place":
		return r.Place, nil
	case "required", "Required":
		return r.Required, nil
	case "updatedAt", "UpdatedAt":
		return r.UpdatedAt, nil

	}
	return nil, nil
}

func (r *ModuleField) SetValue(name string, pos uint, value any) (err error) {
	switch name {
	case "createdAt", "CreatedAt":
		return cast2.Time(value, &r.CreatedAt)
	case "deletedAt", "DeletedAt":
		return cast2.TimePtr(value, &r.DeletedAt)
	case "id", "ID":
		return cast2.Uint64(value, &r.ID)
	case "kind", "Kind":
		return cast2.String(value, &r.Kind)
	case "label", "Label":
		return cast2.String(value, &r.Label)
	case "moduleID", "ModuleID":
		return cast2.Uint64(value, &r.ModuleID)
	case "multi", "Multi":
		return cast2.Bool(value, &r.Multi)
	case "name", "Name":
		return cast2.String(value, &r.Name)
	case "place", "Place":
		return cast2.Int(value, &r.Place)
	case "required", "Required":
		return cast2.Bool(value, &r.Required)
	case "updatedAt", "UpdatedAt":
		return cast2.TimePtr(value, &r.UpdatedAt)

	default:
		return r.setValue(name, pos, value)

	}
	return nil
}

func (r Namespace) GetID() uint64 { return r.ID }

func (r *Namespace) GetValue(name string, pos uint) (any, error) {
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
	case "slug", "Slug":
		return r.Slug, nil
	case "updatedAt", "UpdatedAt":
		return r.UpdatedAt, nil

	}
	return nil, nil
}

func (r *Namespace) SetValue(name string, pos uint, value any) (err error) {
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
	case "slug", "Slug":
		return cast2.String(value, &r.Slug)
	case "updatedAt", "UpdatedAt":
		return cast2.TimePtr(value, &r.UpdatedAt)

	}
	return nil
}

func (r Page) GetID() uint64 { return r.ID }

func (r *Page) GetValue(name string, pos uint) (any, error) {
	switch name {
	case "createdAt", "CreatedAt":
		return r.CreatedAt, nil
	case "deletedAt", "DeletedAt":
		return r.DeletedAt, nil
	case "description", "Description":
		return r.Description, nil
	case "handle", "Handle":
		return r.Handle, nil
	case "id", "ID":
		return r.ID, nil
	case "moduleID", "ModuleID":
		return r.ModuleID, nil
	case "namespaceID", "NamespaceID":
		return r.NamespaceID, nil
	case "selfID", "SelfID":
		return r.SelfID, nil
	case "title", "Title":
		return r.Title, nil
	case "updatedAt", "UpdatedAt":
		return r.UpdatedAt, nil
	case "visible", "Visible":
		return r.Visible, nil
	case "weight", "Weight":
		return r.Weight, nil

	}
	return nil, nil
}

func (r *Page) SetValue(name string, pos uint, value any) (err error) {
	switch name {
	case "createdAt", "CreatedAt":
		return cast2.Time(value, &r.CreatedAt)
	case "deletedAt", "DeletedAt":
		return cast2.TimePtr(value, &r.DeletedAt)
	case "description", "Description":
		return cast2.String(value, &r.Description)
	case "handle", "Handle":
		return cast2.String(value, &r.Handle)
	case "id", "ID":
		return cast2.Uint64(value, &r.ID)
	case "moduleID", "ModuleID":
		return cast2.Uint64(value, &r.ModuleID)
	case "namespaceID", "NamespaceID":
		return cast2.Uint64(value, &r.NamespaceID)
	case "selfID", "SelfID":
		return cast2.Uint64(value, &r.SelfID)
	case "title", "Title":
		return cast2.String(value, &r.Title)
	case "updatedAt", "UpdatedAt":
		return cast2.TimePtr(value, &r.UpdatedAt)
	case "visible", "Visible":
		return cast2.Bool(value, &r.Visible)
	case "weight", "Weight":
		return cast2.Int(value, &r.Weight)

	default:
		return r.setValue(name, pos, value)

	}
	return nil
}

func (r PageLayout) GetID() uint64 { return r.ID }

func (r *PageLayout) GetValue(name string, pos uint) (any, error) {
	switch name {
	case "createdAt", "CreatedAt":
		return r.CreatedAt, nil
	case "deletedAt", "DeletedAt":
		return r.DeletedAt, nil
	case "handle", "Handle":
		return r.Handle, nil
	case "id", "ID":
		return r.ID, nil
	case "namespaceID", "NamespaceID":
		return r.NamespaceID, nil
	case "ownedBy", "OwnedBy":
		return r.OwnedBy, nil
	case "pageID", "PageID":
		return r.PageID, nil
	case "parentID", "ParentID":
		return r.ParentID, nil
	case "updatedAt", "UpdatedAt":
		return r.UpdatedAt, nil
	case "weight", "Weight":
		return r.Weight, nil

	default:
		return r.getValue(name, pos)

	}
	return nil, nil
}

func (r *PageLayout) SetValue(name string, pos uint, value any) (err error) {
	switch name {
	case "createdAt", "CreatedAt":
		return cast2.Time(value, &r.CreatedAt)
	case "deletedAt", "DeletedAt":
		return cast2.TimePtr(value, &r.DeletedAt)
	case "handle", "Handle":
		return cast2.String(value, &r.Handle)
	case "id", "ID":
		return cast2.Uint64(value, &r.ID)
	case "namespaceID", "NamespaceID":
		return cast2.Uint64(value, &r.NamespaceID)
	case "ownedBy", "OwnedBy":
		return cast2.Uint64(value, &r.OwnedBy)
	case "pageID", "PageID":
		return cast2.Uint64(value, &r.PageID)
	case "parentID", "ParentID":
		return cast2.Uint64(value, &r.ParentID)
	case "updatedAt", "UpdatedAt":
		return cast2.TimePtr(value, &r.UpdatedAt)
	case "weight", "Weight":
		return cast2.Int(value, &r.Weight)

	default:
		return r.setValue(name, pos, value)

	}
	return nil
}

func (r Record) GetID() uint64 { return r.ID }

func (r *Record) GetValue(name string, pos uint) (any, error) {
	switch name {
	case "createdAt", "CreatedAt":
		return r.CreatedAt, nil
	case "createdBy", "CreatedBy", "created_by":
		return r.CreatedBy, nil
	case "deletedAt", "DeletedAt":
		return r.DeletedAt, nil
	case "deletedBy", "DeletedBy", "deleted_by":
		return r.DeletedBy, nil
	case "id", "ID":
		return r.ID, nil
	case "meta", "Meta":
		return r.Meta, nil
	case "moduleID", "ModuleID":
		return r.ModuleID, nil
	case "namespaceID", "NamespaceID":
		return r.NamespaceID, nil
	case "ownedBy", "OwnedBy", "owned_by":
		return r.OwnedBy, nil
	case "revision", "Revision":
		return r.Revision, nil
	case "updatedAt", "UpdatedAt":
		return r.UpdatedAt, nil
	case "updatedBy", "UpdatedBy", "updated_by":
		return r.UpdatedBy, nil

	default:
		return r.getValue(name, pos)

	}
	return nil, nil
}

func (r *Record) SetValue(name string, pos uint, value any) (err error) {
	switch name {
	case "createdAt", "CreatedAt":
		return cast2.Time(value, &r.CreatedAt)
	case "createdBy", "CreatedBy", "created_by":
		return cast2.Uint64(value, &r.CreatedBy)
	case "deletedAt", "DeletedAt":
		return cast2.TimePtr(value, &r.DeletedAt)
	case "deletedBy", "DeletedBy", "deleted_by":
		return cast2.Uint64(value, &r.DeletedBy)
	case "id", "ID":
		return cast2.Uint64(value, &r.ID)
	case "meta", "Meta":
		return cast2.Meta(value, &r.Meta)
	case "moduleID", "ModuleID":
		return cast2.Uint64(value, &r.ModuleID)
	case "namespaceID", "NamespaceID":
		return cast2.Uint64(value, &r.NamespaceID)
	case "ownedBy", "OwnedBy", "owned_by":
		return cast2.Uint64(value, &r.OwnedBy)
	case "revision", "Revision":
		return cast2.Int(value, &r.Revision)
	case "updatedAt", "UpdatedAt":
		return cast2.TimePtr(value, &r.UpdatedAt)
	case "updatedBy", "UpdatedBy", "updated_by":
		return cast2.Uint64(value, &r.UpdatedBy)

	default:
		return r.setValue(name, pos, value)

	}
	return nil
}
