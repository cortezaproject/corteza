package store

// This file is auto-generated.
//
// Template:    pkg/codegen/assets/store_base.gen.go.tpl
// Definitions: store/roles.yaml
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.

import (
	"context"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	Roles interface {
		SearchRoles(ctx context.Context, f types.RoleFilter) (types.RoleSet, types.RoleFilter, error)
		LookupRoleByID(ctx context.Context, id uint64) (*types.Role, error)
		LookupRoleByHandle(ctx context.Context, handle string) (*types.Role, error)
		LookupRoleByName(ctx context.Context, name string) (*types.Role, error)

		CreateRole(ctx context.Context, rr ...*types.Role) error

		UpdateRole(ctx context.Context, rr ...*types.Role) error

		UpsertRole(ctx context.Context, rr ...*types.Role) error

		DeleteRole(ctx context.Context, rr ...*types.Role) error
		DeleteRoleByID(ctx context.Context, ID uint64) error

		TruncateRoles(ctx context.Context) error

		// Additional custom functions

		// RoleMetrics (custom function)
		RoleMetrics(ctx context.Context) (*types.RoleMetrics, error)
	}
)

var _ *types.Role
var _ context.Context

// SearchRoles returns all matching Roles from store
func SearchRoles(ctx context.Context, s Roles, f types.RoleFilter) (types.RoleSet, types.RoleFilter, error) {
	return s.SearchRoles(ctx, f)
}

// LookupRoleByID searches for role by ID
//
// It returns role even if deleted or suspended
func LookupRoleByID(ctx context.Context, s Roles, id uint64) (*types.Role, error) {
	return s.LookupRoleByID(ctx, id)
}

// LookupRoleByHandle searches for role by its handle
//
// It returns only valid roles (not deleted, not archived)
func LookupRoleByHandle(ctx context.Context, s Roles, handle string) (*types.Role, error) {
	return s.LookupRoleByHandle(ctx, handle)
}

// LookupRoleByName searches for role by its name
//
// It returns only valid roles (not deleted, not archived)
func LookupRoleByName(ctx context.Context, s Roles, name string) (*types.Role, error) {
	return s.LookupRoleByName(ctx, name)
}

// CreateRole creates one or more Roles in store
func CreateRole(ctx context.Context, s Roles, rr ...*types.Role) error {
	return s.CreateRole(ctx, rr...)
}

// UpdateRole updates one or more (existing) Roles in store
func UpdateRole(ctx context.Context, s Roles, rr ...*types.Role) error {
	return s.UpdateRole(ctx, rr...)
}

// UpsertRole creates new or updates existing one or more Roles in store
func UpsertRole(ctx context.Context, s Roles, rr ...*types.Role) error {
	return s.UpsertRole(ctx, rr...)
}

// DeleteRole Deletes one or more Roles from store
func DeleteRole(ctx context.Context, s Roles, rr ...*types.Role) error {
	return s.DeleteRole(ctx, rr...)
}

// DeleteRoleByID Deletes Role from store
func DeleteRoleByID(ctx context.Context, s Roles, ID uint64) error {
	return s.DeleteRoleByID(ctx, ID)
}

// TruncateRoles Deletes all Roles from store
func TruncateRoles(ctx context.Context, s Roles) error {
	return s.TruncateRoles(ctx)
}

func RoleMetrics(ctx context.Context, s Roles) (*types.RoleMetrics, error) {
	return s.RoleMetrics(ctx)
}
