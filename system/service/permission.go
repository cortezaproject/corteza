package service

import (
	"context"

	"github.com/crusttech/crust/internal/rules"
	"github.com/crusttech/crust/system/repository"
)

type (
	permission struct {
		db  db
		ctx context.Context

		resources rules.ResourcesInterface
	}

	PermissionService interface {
		With(ctx context.Context) PermissionService

		Get(roleID uint64) (interface{}, error)
		Update(roleID uint64, rules []rules.Rule) (interface{}, error)
		Delete(roleID uint64) (interface{}, error)
	}
)

func Permission() PermissionService {
	return (&permission{}).With(context.Background())
}

func (p *permission) With(ctx context.Context) PermissionService {
	db := repository.DB(ctx)
	return &permission{
		db:  db,
		ctx: ctx,

		resources: rules.NewResources(ctx, db),
	}
}

func (p *permission) Get(roleID uint64) (interface{}, error) {
	return p.resources.List(roleID)
}

func (p *permission) Update(roleID uint64, rules []rules.Rule) (interface{}, error) {
	for _, rule := range rules {
		err := validatePermission(rule.Resource, rule.Operation)
		if err != nil {
			return nil, err
		}
	}
	return nil, p.resources.Grant(roleID, rules)
}

func (p *permission) Delete(roleID uint64) (interface{}, error) {
	return nil, p.resources.Delete(roleID)
}
