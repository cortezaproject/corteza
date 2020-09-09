package store

// This file is auto-generated.
//
// Template:    pkg/codegen/assets/store_base.gen.go.tpl
// Definitions: store/users.yaml
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.

import (
	"context"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	Users interface {
		SearchUsers(ctx context.Context, f types.UserFilter) (types.UserSet, types.UserFilter, error)
		LookupUserByID(ctx context.Context, id uint64) (*types.User, error)
		LookupUserByEmail(ctx context.Context, email string) (*types.User, error)
		LookupUserByHandle(ctx context.Context, handle string) (*types.User, error)
		LookupUserByUsername(ctx context.Context, username string) (*types.User, error)

		CreateUser(ctx context.Context, rr ...*types.User) error

		UpdateUser(ctx context.Context, rr ...*types.User) error

		UpsertUser(ctx context.Context, rr ...*types.User) error

		DeleteUser(ctx context.Context, rr ...*types.User) error
		DeleteUserByID(ctx context.Context, ID uint64) error

		TruncateUsers(ctx context.Context) error

		// Additional custom functions

		// CountUsers (custom function)
		CountUsers(ctx context.Context, _f types.UserFilter) (uint, error)

		// UserMetrics (custom function)
		UserMetrics(ctx context.Context) (*types.UserMetrics, error)
	}
)

var _ *types.User
var _ context.Context

// SearchUsers returns all matching Users from store
func SearchUsers(ctx context.Context, s Users, f types.UserFilter) (types.UserSet, types.UserFilter, error) {
	return s.SearchUsers(ctx, f)
}

// LookupUserByID searches for user by ID
//
// It returns user even if deleted or suspended
func LookupUserByID(ctx context.Context, s Users, id uint64) (*types.User, error) {
	return s.LookupUserByID(ctx, id)
}

// LookupUserByEmail searches for user by their email
//
// It returns only valid users (not deleted, not suspended)
func LookupUserByEmail(ctx context.Context, s Users, email string) (*types.User, error) {
	return s.LookupUserByEmail(ctx, email)
}

// LookupUserByHandle searches for user by their email
//
// It returns only valid users (not deleted, not suspended)
func LookupUserByHandle(ctx context.Context, s Users, handle string) (*types.User, error) {
	return s.LookupUserByHandle(ctx, handle)
}

// LookupUserByUsername searches for user by their username
//
// It returns only valid users (not deleted, not suspended)
func LookupUserByUsername(ctx context.Context, s Users, username string) (*types.User, error) {
	return s.LookupUserByUsername(ctx, username)
}

// CreateUser creates one or more Users in store
func CreateUser(ctx context.Context, s Users, rr ...*types.User) error {
	return s.CreateUser(ctx, rr...)
}

// UpdateUser updates one or more (existing) Users in store
func UpdateUser(ctx context.Context, s Users, rr ...*types.User) error {
	return s.UpdateUser(ctx, rr...)
}

// UpsertUser creates new or updates existing one or more Users in store
func UpsertUser(ctx context.Context, s Users, rr ...*types.User) error {
	return s.UpsertUser(ctx, rr...)
}

// DeleteUser Deletes one or more Users from store
func DeleteUser(ctx context.Context, s Users, rr ...*types.User) error {
	return s.DeleteUser(ctx, rr...)
}

// DeleteUserByID Deletes User from store
func DeleteUserByID(ctx context.Context, s Users, ID uint64) error {
	return s.DeleteUserByID(ctx, ID)
}

// TruncateUsers Deletes all Users from store
func TruncateUsers(ctx context.Context, s Users) error {
	return s.TruncateUsers(ctx)
}

func CountUsers(ctx context.Context, s Users, _f types.UserFilter) (uint, error) {
	return s.CountUsers(ctx, _f)
}

func UserMetrics(ctx context.Context, s Users) (*types.UserMetrics, error) {
	return s.UserMetrics(ctx)
}
