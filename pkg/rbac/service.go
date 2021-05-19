package rbac

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/sentry"
	"go.uber.org/zap"
	"sync"
	"time"
)

type (
	service struct {
		l      *sync.Mutex
		logger *zap.Logger

		//  service will flush values on TRUE or just reload on FALSE
		f chan bool

		rules   RuleSet
		indexed OptRuleSet

		roles []*Role

		store rbacRulesStore
	}

	// RuleFilter is a dummy struct to satisfy store codegen
	RuleFilter struct{}

	Controller interface {
		Can(ses Session, op string, res Resource) bool
		Check(ses Session, op string, res Resource) (v Access)
		Grant(ctx context.Context, rules ...*Rule) (err error)
		Watch(ctx context.Context)
		FindRulesByRoleID(roleID uint64) (rr RuleSet)
		Rules() (rr RuleSet)
		Reload(ctx context.Context)
		UpdateRoles(rr ...*Role)
	}

	RoleSettings struct {
		Bypass        []uint64
		Authenticated []uint64
		Anonymous     []uint64
	}
)

var (
	// Global RBAC service
	gRBAC Controller
)

const (
	watchInterval = time.Hour
)

// Global returns global RBAC service
func Global() Controller {
	return gRBAC
}

func SetGlobal(svc Controller) {
	gRBAC = svc
}

func Initialize(logger *zap.Logger, s rbacRulesStore) error {
	if gRBAC != nil {
		// Prevent multiple initializations
		return nil
	}

	SetGlobal(NewService(logger, s))
	return nil
}

// NewService initializes service{} struct
//
// service{} struct preloads, checks, grants and flushes privileges to and from store
// It acts as a caching layer
func NewService(logger *zap.Logger, s rbacRulesStore) (svc *service) {
	svc = &service{
		l: &sync.Mutex{},
		f: make(chan bool),

		store: s,
	}

	if logger != nil {
		svc.logger = logger.Named("rbac")
	}

	return
}

// Can function performs permission check for roles in context
//
// First extracts roles from context, then
// use Check() to test against permission rules and
// iterate over all fallback functions
//
// System user is always allowed to do everything
//
// When not explicitly allowed through rules or fallbacks, function will return FALSE.
func (svc service) Can(ses Session, op string, res Resource) bool {
	return svc.Check(ses, op, res) == Allow
}

// Check verifies if role has access to perform an operation on a resource
//
// See RuleSet's Check() func for details
func (svc service) Check(ses Session, op string, res Resource) (v Access) {
	svc.l.Lock()
	defer svc.l.Unlock()

	return checkOptimised(
		svc.indexed,
		getContextRoles(ses.Roles(), res, svc.roles),
		op,
		res.RbacResource(),
	)
}

// Grant appends and/or overwrites internal rules slice
//
// All rules with Inherit are removed
func (svc *service) Grant(ctx context.Context, rules ...*Rule) (err error) {
	svc.l.Lock()
	defer svc.l.Unlock()

	svc.grant(rules...)

	return svc.flush(ctx)
}

func (svc *service) grant(rules ...*Rule) {
	svc.rules = merge(svc.rules, rules...)
	// @todo reindex
}

// Watches for changes
func (svc service) Watch(ctx context.Context) {
	go func() {
		defer sentry.Recover()

		var ticker = time.NewTicker(watchInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				svc.Reload(ctx)
			case <-svc.f:
				svc.Reload(ctx)
			}
		}
	}()

	svc.logger.Debug("watcher initialized")
}

func (svc service) FindRulesByRoleID(roleID uint64) (rr RuleSet) {
	svc.l.Lock()
	defer svc.l.Unlock()

	return ruleByRole(svc.rules, roleID)
}

func (svc service) Rules() (rr RuleSet) {
	svc.l.Lock()
	defer svc.l.Unlock()
	return svc.rules
}

func (svc *service) Reload(ctx context.Context) {
	svc.l.Lock()
	defer svc.l.Unlock()
	svc.reloadRules(ctx)
}

func (svc *service) reloadRules(ctx context.Context) {
	rr, _, err := svc.store.SearchRbacRules(ctx, RuleFilter{})
	svc.logger.Debug(
		"reloading rules",
		zap.Error(err),
		zap.Int("before", len(svc.rules)),
		zap.Int("after", len(rr)),
	)

	if err == nil {
		svc.rules = rr
	}
}

func (svc *service) UpdateRoles(rr ...*Role) {
	svc.l.Lock()
	defer svc.l.Unlock()
	svc.logger.Debug(
		"updating roles",
		zap.Int("before", len(svc.rules)),
		zap.Int("after", len(rr)),
	)
	svc.roles = rr
}

func (svc service) flush(ctx context.Context) (err error) {
	d, u := flushable(svc.rules)

	err = svc.store.DeleteRbacRule(ctx, d...)
	if err != nil {
		return
	}

	err = svc.store.UpsertRbacRule(ctx, u...)
	if err != nil {
		return
	}

	clear(u)
	svc.rules = u
	svc.logger.Debug("flushed rules",
		zap.Int("updated", len(u)),
		zap.Int("deleted", len(d)))

	return
}
