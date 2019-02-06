package service

import (
	"context"

	"github.com/pkg/errors"

	"github.com/crusttech/crust/internal/rules"
	"github.com/crusttech/crust/messaging/repository"
)

type (
	permissions struct {
		db  db
		ctx context.Context

		scopes rules.ScopeInterface
	}

	PermissionsService interface {
		With(ctx context.Context) PermissionsService

		List() (interface{}, error)
		Get(team string, scope string, resource string) (interface{}, error)
		Set(team string, permissions []rules.Permission) (interface{}, error)
	}
)

func Permissions(scopes rules.ScopeInterface) PermissionsService {
	return (&permissions{
		scopes: scopes,
	}).With(context.Background())
}

func (svc *permissions) With(ctx context.Context) PermissionsService {
	db := repository.DB(ctx)
	return &permissions{
		db:     db,
		ctx:    ctx,
		scopes: svc.scopes,
	}
}

func (p *permissions) List() (interface{}, error) {
	return p.scopes.List(), nil
}

func (p *permissions) Get(team string, scope string, resource string) (interface{}, error) {
	return nil, errors.New("service.permissions.get: not implemented")
}

func (p *permissions) Set(team string, permissions []rules.Permission) (interface{}, error) {
	return nil, errors.New("service.permissions.set: not implemented")
}
