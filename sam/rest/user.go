package rest

import (
	"context"

	"github.com/crusttech/crust/sam/rest/server"
	"github.com/crusttech/crust/sam/service"
	"github.com/crusttech/crust/sam/types"
	"github.com/pkg/errors"
)

var _ = errors.Wrap

type (
	User struct {
		service userService
	}

	userService interface {
		ValidateCredentials(context.Context, string, string) (*types.User, error)

		FindById(context.Context, uint64) (*types.User, error)
		Find(context.Context, *types.UserFilter) ([]*types.User, error)

		Create(context.Context, *types.User) (*types.User, error)
		Update(context.Context, *types.User) (*types.User, error)

		deleter
		suspender
	}
)

func (User) New() *User {
	return &User{service: service.User()}
}

// User lookup & login
func (self *User) Login(ctx context.Context, r *server.UserLoginRequest) (interface{}, error) {
	return self.service.ValidateCredentials(ctx, r.Username, r.Password)
}

// Searches the users table in the database to find users by matching (by-prefix) their.Username
func (self *User) Search(ctx context.Context, r *server.UserSearchRequest) (interface{}, error) {
	return self.service.Find(ctx, &types.UserFilter{Query: r.Query})
}
