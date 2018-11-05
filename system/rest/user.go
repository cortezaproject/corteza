package rest

import (
	"context"

	"github.com/pkg/errors"

	"github.com/crusttech/crust/system/rest/request"
	"github.com/crusttech/crust/system/service"
	"github.com/crusttech/crust/system/types"
)

var _ = errors.Wrap

type (
	User struct {
		user service.UserService
	}
)

func (User) New() *User {
	ctrl := &User{}
	ctrl.user = service.DefaultUser
	return ctrl
}

// Searches the users table in the database to find users by matching (by-prefix) their.Username
func (ctrl *User) Search(ctx context.Context, r *request.UserSearch) (interface{}, error) {
	return ctrl.user.With(ctx).Find(&types.UserFilter{Query: r.Query})
}
