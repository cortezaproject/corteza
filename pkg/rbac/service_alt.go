package rbac

import (
	"context"
)

type (
	ServiceAllowAll struct{ *service }
)

func (ServiceAllowAll) Can(Session, string, Resource) bool {
	return true
}

func (ServiceAllowAll) Check([]uint64, string, Resource) (v Access) {
	return Allow
}

func (ServiceAllowAll) FindRulesByRoleID(uint64) (rr RuleSet) {
	return
}
func (ServiceAllowAll) Grant(context.Context, ...*Rule) error {
	return nil
}
