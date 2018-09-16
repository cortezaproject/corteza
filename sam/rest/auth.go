package rest

import (
	"context"
	"github.com/crusttech/crust/auth/service"
	"github.com/crusttech/crust/auth/types"
	"github.com/crusttech/crust/internal/auth"
	"github.com/crusttech/crust/sam/rest/request"
	"github.com/pkg/errors"
)

var _ = errors.Wrap

type (
	Auth struct {
		svc struct {
			user  service.UserService
			token auth.TokenEncoder
		}
	}

	authPayload struct {
		JWT  string
		User *types.User `json:"user"`
	}

	authUserBasics interface {
		ValidateCredentials(ctx context.Context, username, password string) (*types.User, error)
		Create(ctx context.Context, input *types.User) (user *types.User, err error)
	}
)

func (Auth) New(tknEncoder auth.TokenEncoder) *Auth {
	ctrl := &Auth{}
	ctrl.svc.user = service.DefaultUser
	ctrl.svc.token = tknEncoder

	return ctrl
}

func (ctrl *Auth) Login(ctx context.Context, r *request.AuthLogin) (interface{}, error) {
	return ctrl.tokenize(ctrl.svc.user.ValidateCredentials(r.Username, r.Password))
}

func (ctrl *Auth) Create(ctx context.Context, r *request.AuthCreate) (interface{}, error) {
	user := &types.User{Username: r.Username}
	user.GeneratePassword(r.Password)

	return ctrl.tokenize(ctrl.svc.user.With(ctx).Create(user))
}

// Wraps user return value and appends JWT
func (ctrl *Auth) tokenize(user *types.User, err error) (interface{}, error) {
	if err != nil {
		return nil, err
	}

	return &authPayload{
		JWT:  ctrl.svc.token.Encode(user),
		User: user,
	}, nil
}

func (ap authPayload) Token() string {
	return ap.JWT
}
