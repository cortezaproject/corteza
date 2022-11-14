package rbac

import (
	"context"
	"fmt"
)

type (
	// ServiceAllowAll constructs not-for-production RBAC service
	ServiceAllowAll struct{}
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

func (ServiceAllowAll) Trace(Session, string, Resource) *Trace {
	return nil
}

func (ServiceAllowAll) CloneRulesByRoleID(context.Context, uint64, ...uint64) error {
	return fmt.Errorf(" ServiceAllowAll does not support rule clonning")
}
