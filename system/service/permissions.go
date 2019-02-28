package service

import (
	"context"
	"strings"

	"github.com/pkg/errors"

	"github.com/crusttech/crust/internal/rules"
	"github.com/crusttech/crust/system/repository"
	"github.com/crusttech/crust/system/types"
)

const (
	delimiter = ":"
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

		Check(resource string, operation string, fallbacks ...rules.CheckAccessFunc) rules.Access

		Read(roleID uint64) (interface{}, error)
		Update(roleID uint64, rules []rules.Rule) (interface{}, error)
		Delete(roleID uint64) (interface{}, error)
	}
)

func Permissions() PermissionsService {
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
		err := p.checkServiceAccess(resource)
		if err == nil {
			for ops := range operations {
				perms = append(perms, types.Permission{Resource: resource, Operation: ops})
			}
		}
	}
	return perms, nil
}

func (p *permissions) Check(resource string, operation string, fallbacks ...rules.CheckAccessFunc) rules.Access {
	return p.resources.Check(resource, operation, fallbacks...)
}

func (p *permissions) Read(roleID uint64) (interface{}, error) {
	ret, err := p.resources.Read(roleID)
	if err != nil {
		return nil, err
	}

	// Only display rules under granted scopes.
	rules := []rules.Rule{}
	for _, rule := range ret {
		err = p.checkServiceAccess(rule.Resource)
		if err == nil {
			rules = append(rules, rule)
		}
	}
	return rules, nil
}

func (p *permissions) Update(roleID uint64, rules []rules.Rule) (interface{}, error) {
	for _, rule := range rules {
		err := validatePermission(rule.Resource, rule.Operation)
		if err != nil {
			return nil, err
		}
		err = p.checkServiceAccess(rule.Resource)
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

func (p *permissions) checkServiceAccess(resource string) error {
	service := strings.Split(resource, delimiter)[0]

	grant := p.resources.Check(service, "grant")
	if grant == rules.Allow {
		return nil
	}
	return errors.Errorf("No grant permissions for: %v", service)
}
