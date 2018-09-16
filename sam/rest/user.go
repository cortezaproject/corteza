package rest

import (
	"context"

	authService "github.com/crusttech/crust/auth/service"
	authTypes "github.com/crusttech/crust/auth/types"
	"github.com/crusttech/crust/sam/rest/request"
	"github.com/crusttech/crust/sam/service"
	"github.com/crusttech/crust/sam/types"
	"github.com/pkg/errors"
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

func (ctrl *User) Message(ctx context.Context, r *request.UserMessage) (interface{}, error) {
	return ctrl.svc.message.Direct(ctx, r.UserID, &types.Message{Message: r.Message})
}
