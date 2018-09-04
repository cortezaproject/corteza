package rest

import (
	"context"
	"github.com/crusttech/crust/auth/rest/request"
	"github.com/crusttech/crust/auth/service"
	"github.com/crusttech/crust/auth/types"
	"github.com/pkg/errors"
)

var _ = errors.Wrap

type (
	Auth struct {
		user  service.UserService
		token types.TokenEncoder
	}
)

func (Auth) New(credValidator service.UserService, tknEncoder types.TokenEncoder) *Auth {
	return &Auth{
		credValidator,
		tknEncoder,
	}
}

func (ctrl *Auth) Login(ctx context.Context, r *request.AuthLogin) (interface{}, error) {
	return ctrl.tokenize(ctrl.user.With(ctx).ValidateCredentials(r.Username, r.Password))
}

func (ctrl *Auth) Create(ctx context.Context, r *request.AuthCreate) (interface{}, error) {
	user := &types.User{Username: r.Username}
	user.GeneratePassword(r.Password)
	return ctrl.tokenize(ctrl.user.With(ctx).Create(user))
}

// Wraps user return value and appends JWT
func (ctrl *Auth) tokenize(user *types.User, err error) (interface{}, error) {
	if err != nil {
		return nil, err
	}

	return struct {
		JWT  string
		User *types.User `json:"user"`
	}{
		JWT:  ctrl.token.Encode(user),
		User: user,
	}, nil
}
