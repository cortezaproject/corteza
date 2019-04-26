package service

import (
	"context"
	"io"
	"testing"

	"github.com/crusttech/crust/system/internal/service"
	"github.com/crusttech/crust/system/types"
)

type (
	UserService interface {
		Create(input *types.User) (*types.User, error)
		Update(mod *types.User) (*types.User, error)

		CreateWithAvatar(input *types.User, avatar io.Reader) (*types.User, error)
		UpdateWithAvatar(mod *types.User, avatar io.Reader) (*types.User, error)

		Delete(id uint64) error

		FindByUsername(username string) (*types.User, error)
		FindByEmail(email string) (*types.User, error)
		FindByID(id uint64) (*types.User, error)
		FindByIDs(id ...uint64) (types.UserSet, error)
		Find(filter *types.UserFilter) (types.UserSet, error)
	}
)

var DefaultUser = service.DefaultUser

func User(ctx context.Context) UserService {
	return DefaultUser.With(ctx)
}

// Expose the full User API for testing
func TestUser(_ *testing.T, ctx context.Context) service.UserService {
	return DefaultUser.With(ctx)
}
