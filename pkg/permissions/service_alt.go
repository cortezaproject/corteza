package permissions

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/titpetric/factory"
	"go.uber.org/zap"
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

func (ServiceAllowAll) Check(res Resource, op Operation, roles ...uint64) (v Access) {
	return Allow
}

func (ServiceAllowAll) Grant(ctx context.Context, wl Whitelist, rules ...*Rule) (err error) {
	return nil
}

func (ServiceAllowAll) FindRulesByRoleID(roleID uint64) (rr RuleSet) {
	return
}

func (ServiceAllowAll) ResourceFilter(context.Context, Resource, Operation, Access) *ResourceFilter {
	return &ResourceFilter{superuser: true}
}

func (ServiceDenyAll) Can(ctx context.Context, res Resource, op Operation, ff ...CheckAccessFunc) bool {
	return false
}

func (ServiceDenyAll) Check(res Resource, op Operation, roles ...uint64) (v Access) {
	return Deny
}

func (ServiceDenyAll) Grant(ctx context.Context, wl Whitelist, rules ...*Rule) (err error) {
	return nil
}

func (ServiceDenyAll) FindRulesByRoleID(roleID uint64) (rr RuleSet) {
	return
}

func (svc *TestService) ClearGrants() {
	svc.repository.Purge()
	svc.rules = RuleSet{}
}

func (svc *TestService) String() (out string) {
	tpl := "%20v\t%-30s\t%-30s\t%v\n"
	out = fmt.Sprintf(tpl, "role", "res", "op", "access")
	out += strings.Repeat("-", 120) + "\n"

	_ = svc.rules.Walk(func(r *Rule) error {
		out += fmt.Sprintf(tpl, r.RoleID, r.Resource, r.Operation, r.Access)
		return nil
	})

	out += strings.Repeat("-", 120) + "\n"

	return
}

func NewTestService(ctx context.Context, logger *zap.Logger, db *factory.DB, tbl string) (svc *TestService) {
	svc = &TestService{
		service: service{
			l: &sync.Mutex{},
			f: make(chan bool),

			logger:     logger.Named("permissions"),
			repository: Repository(db, tbl),
			dbTable:    tbl,
		},
	}

	return
}
