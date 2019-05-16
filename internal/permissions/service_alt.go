package permissions

import (
	"context"
)

type (
	ServiceAllowAll struct{}
	ServiceDenyAll  struct{}
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
