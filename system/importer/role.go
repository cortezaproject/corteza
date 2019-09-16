package importer

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/cortezaproject/corteza-server/pkg/deinterfacer"
	"github.com/cortezaproject/corteza-server/pkg/importer"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	RoleImport struct {
		set   types.RoleSet
		dirty map[string]bool

		permissions importer.PermissionImporter

		finder roleFinder
	}

	roleFinder interface {
		FindByHandle(string) (*types.Role, error)
	}

	roleKeeper interface {
		Update(*types.Role) (*types.Role, error)
		Create(*types.Role) (*types.Role, error)
	}
)

func NewRoleImporter(finder roleFinder, permissions importer.PermissionImporter) *RoleImport {
	return &RoleImport{
		set:         types.RoleSet{},
		dirty:       map[string]bool{},
		permissions: permissions,
		finder:      finder,
	}
}

// Resolves permission rules:
// { <role-handle>: { role } } or [ { role }, ... ]
func (imp *RoleImport) CastSet(set interface{}) error {
	return deinterfacer.Each(set, func(index int, handle string, def interface{}) error {
		if index > -1 {
			// Roles defined as collection
			deinterfacer.KVsetString(&handle, "handle", def)
		}

		return imp.Cast(handle, def)
	})
}

// Resolves permission rules:
// { <role-handle>: { role } } or [ { role }, ... ]
func (imp *RoleImport) Cast(handle string, def interface{}) (err error) {
	var role *types.Role

	if !importer.IsValidHandle(handle) {
		return errors.New("invalid role handle")
	}

	handle = importer.NormalizeHandle(handle)
	if role, err = imp.Get(handle); err != nil {
		return err
	}

	if name, ok := def.(string); ok && name != "" {
		imp.dirty[handle] = role.Name != name
		role.Name = name
		return nil
	}

	return deinterfacer.Each(def, func(_ int, key string, val interface{}) (err error) {
		switch key {
		case "handle":
			// already handled
		case "name":
			name := deinterfacer.ToString(val)
			imp.dirty[handle] = role.Name != name
			role.Name = name

		case "allow", "deny":
			return imp.permissions.CastSet(types.RolePermissionResource.String()+role.Handle, key, val)

		default:
			return fmt.Errorf("unexpected key %q for role %q", key, role.Handle)
		}

		return err
	})
}

// Exists returns true if role exists in the buffer or
// can be loaded from the storage
func (imp *RoleImport) Exists(handle string) bool {
	handle = importer.NormalizeHandle(handle)
	role := imp.set.FindByHandle(handle)
	if role != nil {
		return true
	}

	if imp.finder != nil {
		role, err := imp.finder.FindByHandle(handle)
		if err == nil && role != nil {
			imp.set = append(imp.set, role)
			return true
		}
	}

	return false
}

// finds or makes new role
func (imp *RoleImport) Get(handle string) (*types.Role, error) {
	handle = importer.NormalizeHandle(handle)

	if !importer.IsValidHandle(handle) {
		return nil, errors.New("invalid role handle")
	}

	if !imp.Exists(handle) {
		imp.set = append(imp.set, &types.Role{
			Handle: handle,
			Name:   handle,
		})
	}

	return imp.set.FindByHandle(handle), nil
}

func (imp *RoleImport) Store(ctx context.Context, k roleKeeper) error {
	return imp.set.Walk(func(role *types.Role) (err error) {
		var handle = role.Handle

		if role.ID == 0 {
			role, err = k.Create(role)
		} else {
			if imp.dirty[handle] {
				role, err = k.Update(role)
			}
		}

		if err != nil {
			return
		}

		imp.permissions.UpdateResources(types.RolePermissionResource.String(), handle, role.ID)
		imp.permissions.UpdateRoles(role.Handle, role.ID)

		return
	})
}
