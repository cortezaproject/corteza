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

		db *factory.DB
	}

	ScriptsProvider interface {
		FilterByEvent(event, resource string, cc ...TriggerConditionChecker) ScriptSet
	}

	WatcherService interface {
		Watch(ctx context.Context)
	}

	AutomationServiceConfig struct {
		Logger        *zap.Logger
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
func Service(c AutomationServiceConfig) (svc *service) {
	svc = &service{
		logger: c.Logger.Named("automation"),

		c: c,

		srepo: ScriptRepository(c.DbTablePrefix),
		trepo: TriggerRepository(c.DbTablePrefix),

		db: c.DB,
	}

	// Reload ASAP
	svc.Reload()
	return
}

// Watch watches for changes
func (svc service) Watch(ctx context.Context) {
	svc.f = make(chan bool)
	go func() {
		defer sentry.Recover()
		defer close(svc.f)

		var ticker = time.NewTicker(watchInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				svc.reload(ctx)
			case <-svc.f:
				for len(svc.f) > 0 {
					// Drain just before we reload
					<-svc.f
				}
				svc.reload(ctx)
			}
		}
	}()

	svc.logger.Debug("watcher initialized")
}

func (svc *service) Reload() {
	select {
	case svc.f <- true:
		return
	default:
		// that's ok too..
	}
}

func (svc *service) reload(ctx context.Context) {
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

	ss, err = svc.srepo.findRunnable(svc.db)
	svc.logger.Info("scripts loaded", zap.Error(err), zap.Int("count", len(tt)))
	if err != nil {
		return
	}

	// Only interested in valid scritps
	ss, _ = ss.Filter(func(s *Script) (b bool, e error) {
		return s.IsValid(), nil
	})

	tt, err = svc.trepo.findRunnable(svc.db)
	svc.logger.Info("triggers loaded", zap.Error(err), zap.Int("count", len(tt)))
	if err != nil {
		return
	}

	_ = tt.Walk(func(t *Trigger) error {
		s := ss.FindByID(t.ScriptID)
		if s != nil && t.IsValid() && s.CheckCompatibility(t) != nil {
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

func (svc service) FindScriptByID(ctx context.Context, scriptID uint64) (*Script, error) {
	return svc.srepo.findByID(svc.db, scriptID)
}

func (svc service) FindScripts(ctx context.Context, f ScriptFilter) (ScriptSet, ScriptFilter, error) {
	return svc.srepo.find(svc.db, f)
}

// CreateScript - modifies script's props, pushes to repo & updates scripts cache
func (svc service) CreateScript(ctx context.Context, s *Script) error {
	s.ID = factory.Sonyflake.NextID()
	s.CreatedAt = time.Now()
	s.CreatedBy = auth.GetIdentityFromContext(ctx).Identity()

	return svc.db.Transaction(func() (err error) {
		if err = svc.srepo.create(svc.db, s); err != nil {
			return
		}

		err = s.triggers.Walk(func(t *Trigger) error {
			return svc.setNewTriggerInfo(ctx, s, t)
		})

		if err != nil {
			return
		}

		// Force no-pre-check
		if err = svc.trepo.mergeSet(svc.db, STMS_FRESH, s.ID, s.triggers); err != nil {
			return
		}

		svc.Reload()
		return
	})
}

// UpdateScript - modifies script's props, pushes to repo & updates scripts cache
func (svc service) UpdateScript(ctx context.Context, s *Script) error {
	// Ensure sanity
	s.UpdatedAt, s.UpdatedBy = &time.Time{}, auth.GetIdentityFromContext(ctx).Identity()
	*s.UpdatedAt = time.Now()
	s.DeletedAt, s.DeletedBy = nil, 0

	return svc.db.Transaction(func() (err error) {
		if err = svc.srepo.update(svc.db, s); err != nil {
			return
		}

		err = s.triggers.Walk(func(t *Trigger) error {
			if t.ID == 0 {
				return svc.setNewTriggerInfo(ctx, s, t)
			} else {
				return svc.setUpdatedTriggerInfo(ctx, s, t)
			}
		})

		if err != nil {
			return
		}

		if err = svc.trepo.mergeSet(svc.db, s.tms, s.ID, s.triggers); err != nil {
			return
		}

		svc.Reload()
		return
	})
}

// DeleteScript - modifies script's props, pushes to repo & updates scripts cache
func (svc service) DeleteScript(ctx context.Context, s *Script) (err error) {
	s.DeletedAt = &time.Time{}
	*s.DeletedAt = time.Now()
	s.DeletedBy = auth.GetIdentityFromContext(ctx).Identity()

	// We're doing soft delete in the repo
	if err = svc.srepo.update(svc.db, s); err != nil {
		return err
	}

	if err = svc.trepo.deleteByScriptID(svc.db, s.ID); err != nil {
		return err
	}

	svc.Reload()
	return
}

func (svc service) FindTriggerByID(ctx context.Context, scriptID uint64) (*Trigger, error) {
	return svc.trepo.findByID(svc.db, scriptID)
}

func (svc service) FindTriggers(ctx context.Context, f TriggerFilter) (TriggerSet, TriggerFilter, error) {
	return svc.trepo.find(svc.db, f)
}

// CreateTrigger - modifies script's props, pushes to repo & updates scripts cache
func (svc service) CreateTrigger(ctx context.Context, s *Script, t *Trigger) (err error) {
	if err = svc.setNewTriggerInfo(ctx, s, t); err != nil {
		return err
	}

	if err = svc.trepo.replace(svc.db, t); err != nil {
		return err
	}

	svc.Reload()
	return
}

func (svc service) setNewTriggerInfo(ctx context.Context, s *Script, t *Trigger) (err error) {
	if err = s.CheckCompatibility(t); err != nil {
		return err
	}

	t.ID = factory.Sonyflake.NextID()
	t.CreatedAt = time.Now()
	t.CreatedBy = auth.GetIdentityFromContext(ctx).Identity()
	t.ScriptID = s.ID
	return nil
}

// UpdateTrigger - modifies script's props, pushes to repo & updates scripts cache
func (svc service) UpdateTrigger(ctx context.Context, s *Script, t *Trigger) (err error) {
	if err = svc.setUpdatedTriggerInfo(ctx, s, t); err != nil {
		return err
	}

	if err = svc.trepo.replace(svc.db, t); err != nil {
		return err
	}

	svc.Reload()
	return
}

func (svc service) setUpdatedTriggerInfo(ctx context.Context, s *Script, t *Trigger) (err error) {
	if s.ID != t.ScriptID {
		return errors.New("invalid script-trigger reference")
	}

	if err = s.CheckCompatibility(t); err != nil {
		return err
	}

	t.UpdatedAt = &time.Time{}
	*t.UpdatedAt = time.Now()
	t.UpdatedBy = auth.GetIdentityFromContext(ctx).Identity()
	return nil
}

// DeleteTrigger - modifies script's props, pushes to repo & updates scripts cache
func (svc service) DeleteTrigger(ctx context.Context, t *Trigger) (err error) {
	t.DeletedAt = &time.Time{}
	*t.DeletedAt = time.Now()
	t.DeletedBy = auth.GetIdentityFromContext(ctx).Identity()

	// We're doing soft delete in the repo
	if err = svc.trepo.replace(svc.db, t); err != nil {
		return err
	}

	svc.Reload()
	return
}
