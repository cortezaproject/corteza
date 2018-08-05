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
			user  authCredentialsValidator
			token authTokenEncoder
		}
	}

	authCredentialsValidator interface {
		ValidateCredentials(ctx context.Context, username, password string) (*types.User, error)
	}

	authTokenEncoder interface {
		Encode(identity auth.Identifiable) string
	}
)

func (Auth) New(credValidator authCredentialsValidator, tknEncoder authTokenEncoder) *Auth {
	auth := &Auth{}
	auth.svc.user = credValidator
	auth.svc.token = tknEncoder

	return auth
}

func (ctrl *Auth) Login(ctx context.Context, r *server.AuthLoginRequest) (interface{}, error) {
	user, err := ctrl.svc.user.ValidateCredentials(ctx, r.Username, r.Password)
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
