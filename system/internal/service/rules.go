package service

import (
	"context"

	"github.com/pkg/errors"

	internalRules "github.com/crusttech/crust/internal/rules"
	"github.com/crusttech/crust/system/internal/repository"
	"github.com/crusttech/crust/system/types"
)

const (
	delimiter = ":"
)

type (
	rules struct {
		db  db
		ctx context.Context

		resources internalRules.ResourcesInterface
	}

	EffectiveRules struct {
		Resource  string `json:"resource"`
		Operation string `json:"operation"`
		Allow     bool   `json:"allow"`
	}

	RulesService interface {
		With(ctx context.Context) RulesService

		List() (interface{}, error)
		Effective(filter string) ([]EffectiveRules, error)

		Check(resource internalRules.Resource, operation string, fallbacks ...internalRules.CheckAccessFunc) internalRules.Access

		Read(roleID uint64) (interface{}, error)
		Update(roleID uint64, rules []internalRules.Rule) (interface{}, error)
		Delete(roleID uint64) (interface{}, error)
	}
)

func Rules(ctx context.Context) RulesService {
	return (&rules{}).With(ctx)
}

func (p *rules) With(ctx context.Context) RulesService {
	db := repository.DB(ctx)
	return &rules{
		db:  db,
		ctx: ctx,

		resources: internalRules.NewResources(ctx, db),
	}
}

func (p *rules) List() (interface{}, error) {
	perms := []types.Permission{}
	for resource, operations := range permissionList {
		for ops := range operations {
			perms = append(perms, types.Permission{Resource: resource, Operation: ops})
		}
	}
	return perms, nil
}

func (p *rules) Effective(filter string) (eff []EffectiveRules, err error) {
	eff = []EffectiveRules{}
	for resource, operations := range permissionList {
		// err := p.checkServiceAccess(resource)
		if err == nil {
			for ops := range operations {
				eff = append(eff, EffectiveRules{
					Resource:  resource,
					Operation: ops,
					Allow:     false,
				})
			}
		}
	}
	return
}

func (p *rules) Check(resource internalRules.Resource, operation string, fallbacks ...internalRules.CheckAccessFunc) internalRules.Access {
	return p.resources.Check(resource, operation, fallbacks...)
}

func (p *rules) Read(roleID uint64) (interface{}, error) {
	ret, err := p.resources.Read(roleID)
	if err != nil {
		return nil, err
	}

	// Only display rules under granted scopes.
	rules := []internalRules.Rule{}
	for _, rule := range ret {
		err = p.checkServiceAccess(rule.Resource)
		if err == nil {
			rules = append(rules, rule)
		}
	}
	return rules, nil
}

func (p *rules) Update(roleID uint64, rules []internalRules.Rule) (interface{}, error) {
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

func (p *rules) Delete(roleID uint64) (interface{}, error) {
	return nil, p.resources.Delete(roleID)
}

func (p *rules) checkServiceAccess(resource internalRules.Resource) error {
	// First, we trim off resource to get service - only service resource have grant operation
	grant := p.resources.Check(resource.GetService(), "grant")
	if grant == internalRules.Allow {
		return nil
	}
	return errors.Errorf("No grant permissions for: %q", resource)
}
