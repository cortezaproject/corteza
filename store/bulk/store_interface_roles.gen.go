package bulk

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
//  - store/roles.yaml

import (
	"context"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	rolesStore interface {
		SearchRoles(ctx context.Context, f types.RoleFilter) (types.RoleSet, types.RoleFilter, error)
		LookupRoleByID(ctx context.Context, id uint64) (*types.Role, error)
		LookupRoleByHandle(ctx context.Context, handle string) (*types.Role, error)
		LookupRoleByName(ctx context.Context, name string) (*types.Role, error)
		CreateRole(ctx context.Context, rr ...*types.Role) error
		UpdateRole(ctx context.Context, rr ...*types.Role) error
		PartialUpdateRole(ctx context.Context, onlyColumns []string, rr ...*types.Role) error
		RemoveRole(ctx context.Context, rr ...*types.Role) error
		RemoveRoleByID(ctx context.Context, ID uint64) error

		TruncateRoles(ctx context.Context) error
	}
)
