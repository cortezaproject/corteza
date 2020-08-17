package bulk

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
//  - store/role_members.yaml

import (
	"context"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	roleMembersStore interface {
		SearchRoleMembers(ctx context.Context, f types.RoleMemberFilter) (types.RoleMemberSet, types.RoleMemberFilter, error)
		CreateRoleMember(ctx context.Context, rr ...*types.RoleMember) error
		UpdateRoleMember(ctx context.Context, rr ...*types.RoleMember) error
		PartialUpdateRoleMember(ctx context.Context, onlyColumns []string, rr ...*types.RoleMember) error
		RemoveRoleMember(ctx context.Context, rr ...*types.RoleMember) error
		RemoveRoleMemberByUserIDRoleID(ctx context.Context, userID uint64, roleID uint64) error

		TruncateRoleMembers(ctx context.Context) error
	}
)
