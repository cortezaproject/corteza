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
	Role struct {
		set         types.RoleSet
		dirty       map[uint64]bool
		permissions importer.PermissionImporter
	}

	roleKeeper interface {
		Update(*types.Role) (*types.Role, error)
		Create(*types.Role) (*types.Role, error)
	}
)

func NewRoleImport(permissions importer.PermissionImporter, set types.RoleSet) *Role {
	if set == nil {
		set = types.RoleSet{}
	}

	out := &Role{
		set:         set,
		dirty:       make(map[uint64]bool),
		permissions: permissions,
	}

	return out
}

func (rImp *Role) CastSet(set interface{}) error {
	return deinterfacer.Each(set, func(index int, handle string, def interface{}) error {
		if index > -1 {
			// Roles defined as collection
			deinterfacer.KVsetString(&handle, "handle", def)
		}

		return rImp.Cast(handle, def)
	})
}

func (rImp *Role) Cast(handle string, def interface{}) (err error) {
	var role *types.Role

	if !importer.IsValidHandle(handle) {
		return errors.New("invalid role handle")
	}

	handle = importer.NormalizeHandle(handle)
	if role, err = rImp.Get(handle); err != nil {
		return err
	} else if role == nil {
		role = &types.Role{
			Handle: handle,
		}

		rImp.set = append(rImp.set, role)
	} else if role.ID == 0 {
		return errors.Errorf("role handle %q already defined in this import session", role.Handle)
	} else {
		rImp.dirty[role.ID] = true
	}

	if name, ok := def.(string); ok && name != "" {
		role.Name = name
		return nil
	}

	return deinterfacer.Each(def, func(_ int, key string, val interface{}) (err error) {
		switch key {
		case "handle":
			// already handled
		case "name":
			name := deinterfacer.ToString(val)
			role.Name = name

		case "allow", "deny":
			return rImp.permissions.CastSet(types.RolePermissionResource.String()+role.Handle, key, val)

		default:
			return fmt.Errorf("unexpected key %q for role %q", key, role.Handle)
		}

		return err
	})
}

func (rImp *Role) Get(handle string) (*types.Role, error) {
	handle = importer.NormalizeHandle(handle)

	if !importer.IsValidHandle(handle) {
		return nil, errors.New("invalid role handle")
	}

	return rImp.set.FindByHandle(handle), nil
}

func (rImp *Role) Store(ctx context.Context, k roleKeeper) error {
	return rImp.set.Walk(func(role *types.Role) (err error) {
		var handle = role.Handle

		if role.ID == 0 {
			role, err = k.Create(role)
		} else if rImp.dirty[role.ID] {
			role, err = k.Update(role)
		}

		if err != nil {
			return
		}

		rImp.permissions.UpdateResources(types.RolePermissionResource.String(), handle, role.ID)
		rImp.permissions.UpdateRoles(role.Handle, role.ID)

		return
	})
}
