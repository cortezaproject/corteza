package service

import (
	"context"

	"github.com/crusttech/crust/internal/rules"
	"github.com/crusttech/crust/system/repository"
	"github.com/crusttech/crust/system/types"
)

type (
	permissions struct {
		db  db
		ctx context.Context

		resources rules.ResourcesInterface
	}

	PermissionsService interface {
		With(ctx context.Context) PermissionsService

		List() (interface{}, error)
		Read(roleID uint64) (interface{}, error)
		Update(roleID uint64, rules []rules.Rule) (interface{}, error)
		Delete(roleID uint64) (interface{}, error)
	}
)

func Permission() PermissionsService {
	return (&permissions{}).With(context.Background())
}

func (p *permissions) With(ctx context.Context) PermissionsService {
	db := repository.DB(ctx)
	return &permissions{
		db:  db,
		ctx: ctx,

		resources: rules.NewResources(ctx, db),
	}
}

func (p *permissions) List() (interface{}, error) {
	perms := []types.Permission{}
	for resource, operations := range permissionList {
		for ops := range operations {
			perms = append(perms, types.Permission{Resource: resource, Operation: ops})
		}
	}
	return perms, nil
}

func (p *permissions) Read(roleID uint64) (interface{}, error) {
	return p.resources.Read(roleID)
}

func (p *permissions) Update(roleID uint64, rules []rules.Rule) (interface{}, error) {
	for _, rule := range rules {
		err := validatePermission(rule.Resource, rule.Operation)
		if err != nil {
			return nil, err
		}
	}
	err := p.resources.Grant(roleID, rules)
	if err != nil {
		return nil, err
	}
	return p.resources.Read(roleID)
}

func (p *permissions) Delete(roleID uint64) (interface{}, error) {
	return nil, p.resources.Delete(roleID)
}
