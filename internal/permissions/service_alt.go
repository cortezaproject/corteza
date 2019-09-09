package permissions

import (
	"context"
)

type (
	ServiceAllowAll struct{}
	ServiceDenyAll  struct{}
	TestService     struct {
		service
	}
)

func (ServiceAllowAll) Can(ctx context.Context, res Resource, op Operation, ff ...CheckAccessFunc) bool {
	return true
}

func (ServiceAllowAll) Grant(ctx context.Context, wl Whitelist, rules ...*Rule) (err error) {
	return nil
}

func (ServiceAllowAll) FindRulesByRoleID(roleID uint64) (rr RuleSet) {
	return
}

func (ServiceDenyAll) Can(ctx context.Context, res Resource, op Operation, ff ...CheckAccessFunc) bool {
	return false
}

func (ServiceDenyAll) Grant(ctx context.Context, wl Whitelist, rules ...*Rule) (err error) {
	return nil
}

func (ServiceDenyAll) FindRulesByRoleID(roleID uint64) (rr RuleSet) {
	return
}

func (svc *TestService) Grant(ctx context.Context, wl Whitelist, rules ...*Rule) (err error) {
	if err = svc.checkRules(wl, rules...); err != nil {
		return err
	}

	svc.grant(rules...)
	return nil
}

func (svc *TestService) ClearGrants() {
	svc.rules = RuleSet{}
}

func NewTestService() *TestService {
	return &TestService{
		service: service{},
	}
}
