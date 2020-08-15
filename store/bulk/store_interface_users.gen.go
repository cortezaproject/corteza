package bulk

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
//  - store/users.yaml

import (
	"context"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	usersStore interface {
		SearchUsers(ctx context.Context, f types.UserFilter) (types.UserSet, types.UserFilter, error)
		LookupUserByID(ctx context.Context, id uint64) (*types.User, error)
		LookupUserByEmail(ctx context.Context, email string) (*types.User, error)
		LookupUserByHandle(ctx context.Context, handle string) (*types.User, error)
		LookupUserByUsername(ctx context.Context, username string) (*types.User, error)
		CreateUser(ctx context.Context, rr ...*types.User) error
		UpdateUser(ctx context.Context, rr ...*types.User) error
		PartialUpdateUser(ctx context.Context, onlyColumns []string, rr ...*types.User) error
		RemoveUser(ctx context.Context, rr ...*types.User) error
		RemoveUserByID(ctx context.Context, ID uint64) error

		TruncateUsers(ctx context.Context) error
	}
)
