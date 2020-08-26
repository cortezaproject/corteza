package permissions

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/sentry"
	"github.com/pkg/errors"
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

		rules RuleSet

		store rbacRulesStore
	}

	// RuleFilter is a dummy struct to satisfy store codegen
	RuleFilter struct{}
)

const (
	watchInterval = time.Hour
)

// Service initializes service{} struct
//
// service{} struct preloads, checks, grants and flushes privileges to and from store
// It acts as a caching layer
func Service(ctx context.Context, logger *zap.Logger, s rbacRulesStore) (svc *service) {
	svc = &service{
		l: &sync.Mutex{},
		f: make(chan bool),

		logger: logger.Named("permissions"),
		store:  s,
	}

	svc.Reload(ctx)
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
func (svc service) Can(roles []uint64, res Resource, op Operation, ff ...CheckAccessFunc) bool {
	// Checking rules
	var v = svc.Check(res.PermissionResource(), op, roles...)
	if v != Inherit {
		return v == Allow
	}

	// Checking fallback functions
	for _, f := range ff {
		v = f()

		if v != Inherit {
			return v == Allow
		}
	}

	return false
}

// Check verifies if role has access to perform an operation on a resource
//
// See RuleSet's Check() func for details
func (svc service) Check(res Resource, op Operation, roles ...uint64) (v Access) {
	svc.l.Lock()
	defer svc.l.Unlock()

	return svc.rules.Check(res, op, roles...)
}

// Grant appends and/or overwrites internal rules slice
//
// All rules with Inherit are removed
func (svc *service) Grant(ctx context.Context, wl Whitelist, rules ...*Rule) (err error) {
	svc.l.Lock()
	defer svc.l.Unlock()

	if err = svc.checkRules(wl, rules...); err != nil {
		return err
	}

	svc.grant(rules...)

	return svc.flush(ctx)
}

func (svc service) checkRules(wl Whitelist, rules ...*Rule) error {
	for _, r := range rules {
		if !wl.Check(r) {
			return errors.Errorf("invalid rule: '%s' on '%s'", r.Operation, r.Resource)
		}
	}

	return nil
}

func (svc *service) grant(rules ...*Rule) {
	svc.rules = svc.rules.merge(rules...)
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

	rr, _ = svc.rules.Filter(func(rule *Rule) (b bool, e error) {
		return rule.RoleID == roleID, nil
	})

	return
}

func (svc service) Rules() (rr RuleSet) {
	svc.l.Lock()
	defer svc.l.Unlock()
	return svc.rules
}

func (svc *service) Reload(ctx context.Context) {
	svc.l.Lock()
	defer svc.l.Unlock()

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

// ResourceFilter is store helper that we use to filter resources directly in the database
//
// See ResourceFilter struct documentation for details
//
// @deprecated
func (svc *service) ResourceFilter(roles []uint64, r Resource, op Operation, fallback Access) *ResourceFilter {
	return &ResourceFilter{
		roles:     roles,
		resource:  r,
		operation: op,
		dbTable:   "rbac_rules",
		chk:       svc,
		fallback:  fallback,
		pkColName: "id",
	}
}

func (svc service) flush(ctx context.Context) (err error) {
	d, u := svc.rules.dirty()

	err = svc.store.DeleteRbacRule(ctx, d...)
	if err != nil {
		return
	}

	err = svc.store.UpsertRbacRule(ctx, u...)
	if err != nil {
		return
	}

	u.clear()
	svc.rules = u
	svc.logger.Debug("flushed rules",
		zap.Int("updated", len(u)),
		zap.Int("deleted", len(d)))

	return
}
