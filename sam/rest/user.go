package rest

import (
	"context"
	"github.com/crusttech/crust/sam/rest/request"
	"github.com/pkg/errors"
)

var _ = errors.Wrap

type User struct{}

func (User) New() *User {
	return &User{}
}

func (ctrl *User) Search(ctx context.Context, r *request.UserSearch) (interface{}, error) {
	return nil, errors.New("Not implemented: User.search")
}

func (ctrl *User) Message(ctx context.Context, r *request.UserMessage) (interface{}, error) {
	return nil, errors.New("Not implemented: User.message")
}
