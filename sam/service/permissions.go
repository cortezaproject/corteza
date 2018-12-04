package service

import (
	"context"

	"github.com/pkg/errors"

	"github.com/crusttech/crust/internal/rbac"
	"github.com/crusttech/crust/sam/repository"
)

type (
	permissions struct {
		db  db
		ctx context.Context
	}

	PermissionsService interface {
		With(ctx context.Context) PermissionsService

		List() (interface{}, error)
		Get(team string, scope string, resource string) (interface{}, error)
		Set(team string, permissions []rbac.Permission) (interface{}, error)
	}
)

func Permissions() PermissionsService {
	return (&permissions{}).With(context.Background())
}

func (svc *permissions) With(ctx context.Context) PermissionsService {
	db := repository.DB(ctx)
	return &permissions{
		db:  db,
		ctx: ctx,
	}
}

func (p *permissions) List() (interface{}, error) {
	return nil, errors.New("service.permissions.list: not implemented")
}

func (p *permissions) Get(team string, scope string, resource string) (interface{}, error) {
	return nil, errors.New("service.permissions.get: not implemented")
}

func (p *permissions) Set(team string, permissions []rbac.Permission) (interface{}, error) {
	return nil, errors.New("service.permissions.set: not implemented")
}
