package permissions

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/sentry"
)

type (
	service struct {
		l      *sync.Mutex
		logger *zap.Logger

		//  service will flush values on TRUE or just reload on FALSE
		f chan bool

		rules      RuleSet
		repository *repository
	}
)

const (
	watchInterval = time.Hour
)

// Service initializes service{} struct
//
// service{} struct preloads, checks, grants and flushes privileges to and from repository
// It acts as a caching layer
func Service(ctx context.Context, logger *zap.Logger, repository *repository) (svc *service) {
	svc = &service{
		l: &sync.Mutex{},
		f: make(chan bool),

		logger:     logger.Named("permissions"),
		repository: repository,
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
func (svc service) Can(ctx context.Context, res Resource, op Operation, ff ...CheckAccessFunc) bool {
	{
		// @todo remove this ASAP
		//       for now, we need it because of complex init/setup relations under system
		ctxTestingVal := ctx.Value("testing")
		if t, ok := ctxTestingVal.(bool); ok && t {
			return true
		}
	}

	u := auth.GetIdentityFromContext(ctx)

	if auth.IsSuperUser(u) {
		return true
	}

	var roles = u.Roles()
	// Checking rules
	var v = svc.Check(res, op, roles...)
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

func (svc *service) Reload(ctx context.Context) {
	svc.l.Lock()
	defer svc.l.Unlock()

	rr, err := svc.repository.With(ctx).Load()
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

func (svc service) flush(ctx context.Context) (err error) {
	d, u := svc.rules.dirty()
	err = svc.repository.With(ctx).Store(d, u)

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
