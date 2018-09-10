package rest

import (
	"context"
	auth "github.com/crusttech/crust/auth/types"
	"github.com/crusttech/crust/sam/rest/request"
	"github.com/crusttech/crust/sam/types"
	"github.com/pkg/errors"
)

var _ = errors.Wrap

type (
	Auth struct {
		user  authUserBasics
		token auth.TokenEncoder
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

func (Auth) New(credValidator authUserBasics, tknEncoder auth.TokenEncoder) *Auth {
	return &Auth{
		credValidator,
		tknEncoder,
	}
}

func (ctrl *Auth) Login(ctx context.Context, r *request.AuthLogin) (interface{}, error) {
	return ctrl.tokenize(ctrl.user.ValidateCredentials(ctx, r.Username, r.Password))
}

func (ctrl *Auth) Create(ctx context.Context, r *request.AuthCreate) (interface{}, error) {
	user := &types.User{Username: r.Username}
	user.GeneratePassword(r.Password)

	return ctrl.tokenize(ctrl.user.Create(ctx, user))
}

// Wraps user return value and appends JWT
func (ctrl *Auth) tokenize(user *types.User, err error) (interface{}, error) {
	if err != nil {
		return nil, err
	}

	return &authPayload{
		JWT:  ctrl.token.Encode(user),
		User: user,
	}, nil
}

func (ap authPayload) Token() string {
	return ap.JWT
}
