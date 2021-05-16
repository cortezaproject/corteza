package rbac

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"strings"
	"sync"
)

type (
	ServiceAllowAll struct{}
	ServiceDenyAll  struct{}
	TestService     struct {
		service
	}
)

func (ServiceAllowAll) Can([]uint64, string, Resource) bool {
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

func (ServiceDenyAll) Can([]uint64, string, string) bool {
	return false
}

func (ServiceDenyAll) Check(string, string, ...uint64) (v Access) {
	return Deny
}

func (ServiceDenyAll) Grant(context.Context, ...*Rule) error {
	return nil
}

func (ServiceDenyAll) FindRulesByRoleID(uint64) (rr RuleSet) {
	return
}

func (svc *TestService) ClearGrants() {
	_ = svc.store.TruncateRbacRules(context.Background())
	svc.rules = RuleSet{}
}

func (svc *TestService) String() (out string) {
	tpl := "%20v\t%-30s\t%-30s\t%v\n"
	out = fmt.Sprintf(tpl, "role", "res", "op", "access")
	out += strings.Repeat("-", 120) + "\n"

	for _, r := range svc.rules {
		out += fmt.Sprintf(tpl, r.RoleID, r.Resource, r.Operation, r.Access)
	}

	out += strings.Repeat("-", 120) + "\n"

	return
}

func NewTestService(logger *zap.Logger, s rbacRulesStore) (svc *TestService) {
	svc = &TestService{
		service: service{
			l: &sync.Mutex{},
			f: make(chan bool),

			logger: logger.Named("rbac-test"),
			store:  s,
		},
	}

	return
}
