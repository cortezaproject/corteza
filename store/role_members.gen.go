package store

// This file is auto-generated.
//
// Template:    pkg/codegen/assets/store_base.gen.go.tpl
// Definitions: store/role_members.yaml
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.

import (
	"context"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	RoleMembers interface {
		CreateRoleMember(ctx context.Context, rr ...*types.RoleMember) error

		DeleteRoleMember(ctx context.Context, rr ...*types.RoleMember) error
		DeleteRoleMemberByUserIDRoleID(ctx context.Context, userID uint64, roleID uint64) error

		TruncateRoleMembers(ctx context.Context) error

		// Additional custom functions

		// SearchRoleMembers (custom function)
		SearchRoleMembers(ctx context.Context, _roleID uint64) ([]uint64, error)

		// SearchUserMemberships (custom function)
		SearchUserMemberships(ctx context.Context, _userID uint64) ([]uint64, error)
	}
)

var _ *types.RoleMember
var _ context.Context

// CreateRoleMember creates one or more RoleMembers in store
func CreateRoleMember(ctx context.Context, s RoleMembers, rr ...*types.RoleMember) error {
	return s.CreateRoleMember(ctx, rr...)
}

// DeleteRoleMember Deletes one or more RoleMembers from store
func DeleteRoleMember(ctx context.Context, s RoleMembers, rr ...*types.RoleMember) error {
	return s.DeleteRoleMember(ctx, rr...)
}

// DeleteRoleMemberByUserIDRoleID Deletes RoleMember from store
func DeleteRoleMemberByUserIDRoleID(ctx context.Context, s RoleMembers, userID uint64, roleID uint64) error {
	return s.DeleteRoleMemberByUserIDRoleID(ctx, userID, roleID)
}

// TruncateRoleMembers Deletes all RoleMembers from store
func TruncateRoleMembers(ctx context.Context, s RoleMembers) error {
	return s.TruncateRoleMembers(ctx)
}

func SearchRoleMembers(ctx context.Context, s RoleMembers, _roleID uint64) ([]uint64, error) {
	return s.SearchRoleMembers(ctx, _roleID)
}

func SearchUserMemberships(ctx context.Context, s RoleMembers, _userID uint64) ([]uint64, error) {
	return s.SearchUserMemberships(ctx, _userID)
}
