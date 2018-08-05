package rest

import (
	"context"
	"github.com/crusttech/crust/auth"
	"github.com/crusttech/crust/sam/rest/server"
	"github.com/crusttech/crust/sam/types"
	"github.com/pkg/errors"
)

var _ = errors.Wrap

type (
	Auth struct {
		svc struct {
			user  authUserBasics
			token authTokenEncoder
		}
	}

	authUserBasics interface {
		ValidateCredentials(ctx context.Context, username, password string) (*types.User, error)
		Create(ctx context.Context, input *types.User) (user *types.User, err error)
	}

	authTokenEncoder interface {
		Encode(identity auth.Identifiable) string
	}
)

func (Auth) New(credValidator authUserBasics, tknEncoder authTokenEncoder) *Auth {
	auth := &Auth{}
	auth.svc.user = credValidator
	auth.svc.token = tknEncoder

	return auth
}

func (ctrl *Auth) Login(ctx context.Context, r *server.AuthLoginRequest) (interface{}, error) {
	return ctrl.tokenize(ctrl.svc.user.ValidateCredentials(ctx, r.Username, r.Password))
}

func (ctrl *Auth) Create(ctx context.Context, r *server.AuthCreateRequest) (interface{}, error) {
	user := &types.User{Username: r.Username}
	user.GeneratePassword(r.Password)

	return ctrl.tokenize(ctrl.svc.user.Create(ctx, user))
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
		JWT:  ctrl.svc.token.Encode(user),
		User: user,
	}, nil
}
