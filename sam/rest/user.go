package rest

import (
	"context"

	"github.com/pkg/errors"

	authService "github.com/crusttech/crust/auth/service"
	authTypes "github.com/crusttech/crust/auth/types"
	"github.com/crusttech/crust/sam/rest/request"
	"github.com/crusttech/crust/sam/service"
)

var _ = errors.Wrap

type (
	User struct {
		svc struct {
			user    authService.UserService
			message service.MessageService
		}
	}
)

func (User) New() *User {
	ctrl := &User{}
	ctrl.svc.user = authService.DefaultUser
	ctrl.svc.message = service.DefaultMessage
	return ctrl
}

// Searches the users table in the database to find users by matching (by-prefix) their.Username
func (ctrl *User) Search(ctx context.Context, r *request.UserSearch) (interface{}, error) {
	return ctrl.svc.user.With(ctx).Find(&authTypes.UserFilter{Query: r.Query})
}
