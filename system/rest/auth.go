package rest

import (
	"context"
	"github.com/crusttech/crust/system/rest/request"
	"github.com/pkg/errors"
)

var _ = errors.Wrap

type Auth struct{}

func (Auth) New() *Auth {
	return &Auth{}
}

func (ctrl *Auth) Check(ctx context.Context, r *request.AuthCheck) (interface{}, error) {
	return nil, errors.New("Not implemented: Auth.check")
}
