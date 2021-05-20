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
		SearchRoleMembers(ctx context.Context, f types.RoleMemberFilter) (types.RoleMemberSet, types.RoleMemberFilter, error)

		CreateRoleMember(ctx context.Context, rr ...*types.RoleMember) error

		UpdateRoleMember(ctx context.Context, rr ...*types.RoleMember) error

		UpsertRoleMember(ctx context.Context, rr ...*types.RoleMember) error

		DeleteRoleMember(ctx context.Context, rr ...*types.RoleMember) error
		DeleteRoleMemberByUserIDRoleID(ctx context.Context, userID uint64, roleID uint64) error

		TruncateRoleMembers(ctx context.Context) error

		// Additional custom functions

		// TransferRoleMembers (custom function)
		TransferRoleMembers(ctx context.Context, _srcRole uint64, _dstRole uint64) error
	}
)

var _ *types.RoleMember
var _ context.Context

// SearchRoleMembers returns all matching RoleMembers from store
func SearchRoleMembers(ctx context.Context, s RoleMembers, f types.RoleMemberFilter) (types.RoleMemberSet, types.RoleMemberFilter, error) {
	return s.SearchRoleMembers(ctx, f)
}

// CreateRoleMember creates one or more RoleMembers in store
func CreateRoleMember(ctx context.Context, s RoleMembers, rr ...*types.RoleMember) error {
	return s.CreateRoleMember(ctx, rr...)
}

// UpdateRoleMember updates one or more (existing) RoleMembers in store
func UpdateRoleMember(ctx context.Context, s RoleMembers, rr ...*types.RoleMember) error {
	return s.UpdateRoleMember(ctx, rr...)
}

// UpsertRoleMember creates new or updates existing one or more RoleMembers in store
func UpsertRoleMember(ctx context.Context, s RoleMembers, rr ...*types.RoleMember) error {
	return s.UpsertRoleMember(ctx, rr...)
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

func TransferRoleMembers(ctx context.Context, s RoleMembers, _srcRole uint64, _dstRole uint64) error {
	return s.TransferRoleMembers(ctx, _srcRole, _dstRole)
}
