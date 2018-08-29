package rest

import (
	"context"

	"github.com/crusttech/crust/sam/rest/request"
	"github.com/crusttech/crust/sam/service"
	"github.com/crusttech/crust/sam/types"
	"github.com/pkg/errors"
)

var _ = errors.Wrap

type (
	User struct {
		svc struct {
			user    service.UserService
			message service.MessageService
		}
	}
)

func (User) New(user service.UserService, message service.MessageService) *User {
	ctrl := &User{}
	ctrl.svc.user = user
	ctrl.svc.message = message
	return ctrl
}

// Searches the users table in the database to find users by matching (by-prefix) their.Username
func (ctrl *User) Search(ctx context.Context, r *request.UserSearch) (interface{}, error) {
	return ctrl.svc.user.Find(ctx, &types.UserFilter{Query: r.Query})
}

func (ctrl *User) Message(ctx context.Context, r *request.UserMessage) (interface{}, error) {
	return ctrl.svc.message.Direct(ctx, r.UserID, &types.Message{Message: r.Message})
}
