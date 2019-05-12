package permissions

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/crusttech/crust/internal/auth"
)

type (
	service struct {
		l      sync.Mutex
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
		f: make(chan bool, 0),

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

	var roles = auth.GetIdentityFromContext(ctx).Roles()
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

	for _, r := range rules {
		if !wl.Check(r) {
			return errors.Errorf("invalid rule: '%s' on '%s'", r.Operation, r.Resource)
		}
	}

	if svc.rules, err = svc.rules.merge(rules...); err != nil {
		return
	}

	return svc.flush(ctx)
}

// Watches for changes
func (svc service) Watch(ctx context.Context) {
	go func() {
		var ticker = time.NewTicker(watchInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				break
			case <-ticker.C:
				svc.Reload(ctx)
			case <-svc.f:
				svc.Reload(ctx)
			}
		}
	}()

	svc.logger.Debug("watcher initialized")
}

func (svc *service) Reload(ctx context.Context) {
	svc.l.Lock()
	defer svc.l.Unlock()

	rr, err := svc.repository.With(ctx).Load()
	svc.logger.Info(
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
	d, u := svc.rules.split()
	err = svc.repository.With(ctx).Store(d, u)

	if err != nil {
		return
	}

	svc.rules = u
	svc.logger.Info("flushed rules",
		zap.Int("updated", len(u)),
		zap.Int("deleted", len(d)))

	return
}
