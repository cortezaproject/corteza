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
		svc struct {
			user    userService
			message userMessageService
		}
	}

	userService interface {
		Find(ctx context.Context, filter *types.UserFilter) ([]*types.User, error)
	}

	userMessageService interface {
		Direct(ctx context.Context, recipientID uint64, in *types.Message) (out *types.Message, err error)
	}
)

func (User) New(userSvc userService, msgSvc userMessageService) *User {
	var ctrl = &User{}
	ctrl.svc.user = userSvc
	ctrl.svc.message = msgSvc
	return ctrl
}

// Searches the users table in the database to find users by matching (by-prefix) their.Username
func (ctrl *User) Search(ctx context.Context, r *server.UserSearchRequest) (interface{}, error) {
	return ctrl.svc.user.Find(ctx, &types.UserFilter{Query: r.Query})
}

func (ctrl *User) Message(ctx context.Context, r *server.UserMessageRequest) (interface{}, error) {
	return ctrl.svc.message.Direct(ctx, r.UserID, &types.Message{Message: r.Message})
}
