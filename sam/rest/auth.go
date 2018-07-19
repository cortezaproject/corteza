package rest

import (
	"context"
	"github.com/crusttech/crust/sam/rest/server"
	"github.com/crusttech/crust/sam/service"
	"github.com/crusttech/crust/sam/types"
	"github.com/pkg/errors"
)

var _ = errors.Wrap

type (
	Auth struct {
		service authUserService
	}

	authUserService interface {
		ValidateCredentials(context.Context, string, string) (*types.User, error)
	}
)

func (Auth) New() *Auth {
	return &Auth{service: service.User()}
}

func (ctrl *Auth) Login(ctx context.Context, r *server.AuthLoginRequest) (interface{}, error) {
	return ctrl.service.ValidateCredentials(ctx, r.Username, r.Password)

}
