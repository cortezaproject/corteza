package automation

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/titpetric/factory"
	"go.uber.org/zap"

	"github.com/cortezaproject/corteza-server/internal/auth"
	"github.com/cortezaproject/corteza-server/pkg/sentry"
)

type (
	service struct {
		l      sync.Mutex
		logger *zap.Logger

		c AutomationServiceConfig

		//  service will flush values on TRUE or just reload on FALSE
		f chan bool

		// internal list of runnable scripts (and their accompanying triggers)
		runnables ScriptSet

		srepo *scriptRepository
		trepo *triggerRepository
	}

	ScriptsProvider interface {
		FilterByEvent(event, resource string, cc ...TriggerConditionChecker) ScriptSet
	}

	WatcherService interface {
		Watch(ctx context.Context)
	}

	AutomationServiceConfig struct {
		DB            *factory.DB
		DbTablePrefix string
	}
)

const (
	watchInterval = time.Hour
)

// Service initializes service{} struct
//
// service{} struct handles scripts & triggers. It acts as a caching layer and
// proxy to repository where it verifies and enriches payloads
//
func Service(ctx context.Context, logger *zap.Logger, c AutomationServiceConfig) (svc *service) {
	svc = &service{
		logger: logger.Named("automation"),

		c: c,

		f: make(chan bool),
	}

	if c.DB != nil {
		svc.srepo = ScriptRepository(c.DB, c.DbTablePrefix)
		svc.trepo = TriggerRepository(c.DB, c.DbTablePrefix)
	}

	svc.Reload(ctx)
	return
}

// Watch() Watches for changes
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

func (svc *service) Reload(ctx context.Context) {
	svc.l.Lock()
	defer svc.l.Unlock()

	if svc.c.DB == nil {
		return
	}

	var (
		err error
		ss  ScriptSet
		tt  TriggerSet
	)

	ss, err = svc.srepo.With(ctx).FindAllRunnable()
	svc.logger.Info("scripts loaded", zap.Error(err), zap.Int("count", len(tt)))
	if err != nil {
		return
	}

	// Only interested in valid scritps
	ss, _ = ss.Filter(func(s *Script) (b bool, e error) {
		return s.IsValid(), nil
	})

	tt, err = svc.trepo.With(ctx).FindAllRunnable()
	svc.logger.Info("triggers loaded", zap.Error(err), zap.Int("count", len(tt)))
	if err != nil {
		return
	}

	_ = tt.Walk(func(t *Trigger) error {
		s := ss.FindByID(t.ScriptID)
		if t.IsValid() && s.CheckCompatibility(t) != nil {
			// Add only compatible triggers
			s.triggers = append(s.triggers, t)
		}

		return nil
	})
}

// FindRunnableScripts scans internal list of runnable scripts and filters them by (trigger's) event and origin
func (svc service) FindRunnableScripts(event, origin string, cc ...TriggerConditionChecker) ScriptSet {
	return svc.runnables.FilterByEvent(event, origin, cc...)
}

// updateRunnableScripts - updates script set (internal runnable scripts list)
func (svc service) updateRunnableScripts(n *Script) {
	svc.l.Lock()
	defer svc.l.Unlock()

	ss := svc.runnables

	for i := range svc.runnables {
		if ss[i].ID != n.ID {
			continue
		}

		if n.IsValid() {
			// Valid, replace
			ss[i] = n
		}

		// Invalid, remove
		ss = append(ss[:i], ss[i+1:]...)
		return
	}

	if n.IsValid() {
		ss = append(ss, n)
	}

	return
}

// updateScriptsWithTrigger - finds the referenced script and updates its trigger set
func (svc service) updateScriptWithTrigger(n *Trigger) {
	svc.l.Lock()
	defer svc.l.Unlock()

	ss := svc.runnables

	for i := range ss {
		if ss[i].ID != n.ScriptID {
			continue
		}

		tt := ss[i].triggers

		for i = range tt {
			if n.IsValid() {
				// Valid, replace
				tt[i] = n
			}

			// Invalid, remove
			tt = append(tt[:i], tt[i+i:]...)
			return
		}

		if n.IsValid() {
			tt = append(tt, n)
		}

		return
	}
}

func (svc service) FindScriptByID(ctx context.Context, scriptID uint64) (*Script, error) {
	return svc.srepo.FindByID(ctx, scriptID)
}

func (svc service) FindScripts(ctx context.Context, f ScriptFilter) (ScriptSet, ScriptFilter, error) {
	return svc.srepo.Find(ctx, f)
}

// CreateScript - modifies script's props, pushes to repo & updates scripts cache
func (svc service) CreateScript(ctx context.Context, s *Script) (err error) {
	s.ID = factory.Sonyflake.NextID()
	s.CreatedAt = time.Now()
	s.CreatedBy = auth.GetIdentityFromContext(ctx).Identity()

	if err = svc.srepo.Create(s); err != nil {
		return err
	}

	svc.updateRunnableScripts(s)
	return
}

// UpdateScript - modifies script's props, pushes to repo & updates scripts cache
func (svc service) UpdateScript(ctx context.Context, s *Script) (err error) {
	s.UpdatedAt = &time.Time{}
	*s.UpdatedAt = time.Now()
	s.UpdatedBy = auth.GetIdentityFromContext(ctx).Identity()

	// Ensure sanity
	s.UpdatedAt, s.UpdatedBy = nil, 0
	s.DeletedAt, s.DeletedBy = nil, 0

	if err = svc.srepo.Update(s); err != nil {
		return err
	}

	svc.updateRunnableScripts(s)
	return
}

// DeleteScript - modifies script's props, pushes to repo & updates scripts cache
func (svc service) DeleteScript(ctx context.Context, s *Script) (err error) {
	s.DeletedAt = &time.Time{}
	*s.DeletedAt = time.Now()
	s.DeletedBy = auth.GetIdentityFromContext(ctx).Identity()

	// We're doing soft delete in the repo
	if err = svc.srepo.Update(s); err != nil {
		return err
	}

	svc.updateRunnableScripts(s)
	return
}

func (svc service) FindTriggerByID(ctx context.Context, scriptID uint64) (*Trigger, error) {
	return svc.trepo.FindByID(ctx, scriptID)
}

func (svc service) FindTriggers(ctx context.Context, f TriggerFilter) (TriggerSet, TriggerFilter, error) {
	return svc.trepo.Find(ctx, f)
}

// CreateScript - modifies script's props, pushes to repo & updates scripts cache
func (svc service) CreateTrigger(ctx context.Context, s *Script, t *Trigger) (err error) {
	if err = s.CheckCompatibility(t); err != nil {
		return err
	}

	t.ID = factory.Sonyflake.NextID()
	t.CreatedAt = time.Now()
	t.CreatedBy = auth.GetIdentityFromContext(ctx).Identity()
	t.ScriptID = s.ID

	if err = svc.trepo.Create(t); err != nil {
		return err
	}

	svc.updateScriptWithTrigger(t)
	return
}

// UpdateTrigger - modifies script's props, pushes to repo & updates scripts cache
func (svc service) UpdateTrigger(ctx context.Context, s *Script, t *Trigger) (err error) {
	if s.ID != t.ScriptID {
		return errors.New("invalid script-trigger reference")
	}

	if err = s.CheckCompatibility(t); err != nil {
		return err
	}

	t.UpdatedAt = &time.Time{}
	*t.UpdatedAt = time.Now()
	t.UpdatedBy = auth.GetIdentityFromContext(ctx).Identity()

	if err = svc.trepo.Update(t); err != nil {
		return err
	}

	svc.updateScriptWithTrigger(t)
	return
}

// DeleteTrigger - modifies script's props, pushes to repo & updates scripts cache
func (svc service) DeleteTrigger(ctx context.Context, t *Trigger) (err error) {
	t.DeletedAt = &time.Time{}
	*t.DeletedAt = time.Now()
	t.DeletedBy = auth.GetIdentityFromContext(ctx).Identity()

	// We're doing soft delete in the repo
	if err = svc.trepo.Update(t); err != nil {
		return err
	}

	svc.updateScriptWithTrigger(t)
	return
}
