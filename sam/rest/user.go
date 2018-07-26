package rest

import (
	"context"

	"github.com/crusttech/crust/sam/rest/server"
	"github.com/crusttech/crust/sam/types"
	"github.com/pkg/errors"
)

var _ = errors.Wrap

type (
	User struct {
		svc userService
	}

	userService interface {
		Find(context.Context, *types.UserFilter) ([]*types.User, error)
	}
)

func (User) New(userSvc userService) *User {
	var ctrl = &User{}
	ctrl.svc = userSvc
	return ctrl
}

// Searches the users table in the database to find users by matching (by-prefix) their.Username
func (ctrl *User) Search(ctx context.Context, r *server.UserSearchRequest) (interface{}, error) {
	return ctrl.svc.Find(ctx, &types.UserFilter{Query: r.Query})
}
