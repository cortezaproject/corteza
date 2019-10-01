package automation

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/titpetric/factory"
	"go.uber.org/zap"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/handle"
	"github.com/cortezaproject/corteza-server/pkg/sentry"
)

type (
	service struct {
		l      *sync.Mutex
		logger *zap.Logger

		c AutomationServiceConfig

		//  service will flush values on TRUE or just reload on FALSE
		f chan bool

		// internal list of runnable scripts (and their accompanying triggers)
		runnables ScriptSet

		// internal list of scheduled scripts
		scheduled scheduledSet

		// turns user-id (rel_runner / runAs) into valid credentials (JWT)
		makeToken TokenMaker

		srepo *scriptRepository
		trepo *triggerRepository

		db *factory.DB
	}

	ScriptsProvider interface {
		FilterByTrigger(event, resource string, cc ...TriggerConditionChecker) ScriptSet
	}

	DeferredAutomationRunner interface {
		RecordDeferred(ctx context.Context, script *Script, ns *types.Namespace, m *types.Module, r *types.Record) (err error)
	}

	TokenMaker func(context.Context, uint64) (string, error)

	WatcherService interface {
		Watch(ctx context.Context)
	}

	AutomationServiceConfig struct {
		Logger        *zap.Logger
		DB            *factory.DB
		TokenMaker    TokenMaker
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
		l: &sync.Mutex{},

		makeToken: c.TokenMaker,

		srepo: ScriptRepository(c.DbTablePrefix),
		trepo: TriggerRepository(c.DbTablePrefix),

		db: c.DB,

		f: make(chan bool, 64),
	}

	// Reload ASAP
	svc.Reload()
	return
}

// Watch watches for changes
func (svc *service) Watch(ctx context.Context) {
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

// RunScheduled runs scheduled scripts periodically
func (svc *service) WatchScheduled(ctx context.Context, r DeferredAutomationRunner) {
	go func() {
		var ticker = time.NewTicker(time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				for _, scriptID := range svc.scheduled.pick() {
					// @todo parallelize this
					svc.logger.Debug(
						"Running scheduled script",
						zap.Uint64("scriptID", scriptID),
					)
					err := r.RecordDeferred(ctx, svc.runnables.FindByID(scriptID), nil, nil, nil)
					if err != nil {
						svc.logger.Error(
							"Script failed to run",
							zap.Uint64("scriptID", scriptID),
							zap.Error(err),
						)
					} else {
						svc.logger.Debug(
							"Ran scheduled script",
							zap.Uint64("scriptID", scriptID),
						)
					}
				}
			}
		}
	}()

	svc.logger.Debug("scheduled runner initialized")
}

func (svc *service) Reload() {
	go func() {
		select {
		case svc.f <- true:
			return
		default:
			// that's ok too..
		}
	}()
}

func (svc *service) reload(ctx context.Context) {
	svc.l.Lock()
	defer svc.l.Unlock()

	var (
		ss ScriptSet
		tt TriggerSet
		db = svc.db.With(ctx)
	)

	_ = db.Transaction(func() (err error) {
		ss, err = svc.srepo.findRunnable(db)
		svc.logger.Info("scripts loaded", zap.Error(err), zap.Int("count", len(ss)))
		if err != nil {
			return
		}

		// Remove all invalid scripts
		svc.runnables, _ = ss.Filter(func(s *Script) (b bool, e error) {
			return s.IsValid(), nil
		})

		tt, err = svc.trepo.findRunnable(db)
		svc.logger.Info("triggers loaded", zap.Error(err), zap.Int("count", len(tt)))
		if err != nil {
			return
		}

		_ = svc.runnables.Walk(func(script *Script) (err error) {
			if script.RunAsDefined() {
				script.credentials, err = svc.makeToken(ctx, script.RunAs)

				if err != nil {
					script.Enabled = false

					svc.logger.Info(
						"could not make token, disabling script",
						zap.Uint64("runAs", script.RunAs),
						zap.Error(err),
					)
				}
			}

			return nil
		})

		_ = tt.Walk(func(t *Trigger) error {
			s := svc.runnables.FindByID(t.ScriptID)
			if s != nil && t.IsValid() && s.CheckCompatibility(t) == nil {
				// Add only compatible triggers
				s.triggers = append(s.triggers, t)
			}

			return nil
		})

		// update scheduled list
		// @todo enable when deferred scripts can be execured

		svc.scheduled = buildScheduleList(svc.runnables)
		svc.logger.Info("deferred scripts scheduled", zap.Int("count", len(svc.scheduled)))

		return
	})
}

// FindRunnableScripts finds runnable scripts in internal list
//
// It uses resource, event and extra condition checkers to filter out all scripts
// that have matching triggers
func (svc service) FindRunnableScripts(resource, event string, cc ...TriggerConditionChecker) ScriptSet {
	svc.l.Lock()
	defer svc.l.Unlock()

	return svc.runnables.FilterByTrigger(
		event,
		resource,
		cc...,
	)
}

func (svc service) FindScriptByID(ctx context.Context, scriptID uint64) (*Script, error) {
	return svc.srepo.findByID(svc.db.With(ctx), scriptID)
}

func (svc service) FindScripts(ctx context.Context, f ScriptFilter) (ScriptSet, ScriptFilter, error) {
	return svc.srepo.find(svc.db.With(ctx), f)
}

// CreateScript - modifies script's props, pushes to repo & updates scripts cache
func (svc service) CreateScript(ctx context.Context, s *Script) error {
	if !handle.IsValid(s.Name) {
		return errors.New("invalid script name")
	}

	s.CreatedAt = time.Now()
	s.CreatedBy = auth.GetIdentityFromContext(ctx).Identity()

	db := svc.db.With(ctx)

	// Reloading scripts at the end (after the transaction completes)
	defer svc.Reload()

	return db.Transaction(func() (err error) {
		if err := svc.UniqueCheck(db, s); err != nil {
			return err
		}

		if err = svc.srepo.create(db, s); err != nil {
			return
		}

		err = s.triggers.Walk(func(t *Trigger) error {
			return svc.setNewTriggerInfo(ctx, s, t)
		})

		if err != nil {
			return
		}

		// Force no-pre-check
		if err = svc.trepo.mergeSet(db, STMS_FRESH, s.ID, s.triggers); err != nil {
			return
		}

		return
	})
}

// UpdateScript - modifies script's props, pushes to repo & updates scripts cache
func (svc service) UpdateScript(ctx context.Context, s *Script) error {
	if !handle.IsValid(s.Name) {
		return errors.New("invalid script name")
	}

	// Ensure sanity
	s.UpdatedAt, s.UpdatedBy = &time.Time{}, auth.GetIdentityFromContext(ctx).Identity()
	*s.UpdatedAt = time.Now()
	s.DeletedAt, s.DeletedBy = nil, 0

	db := svc.db.With(ctx)

	// Reloading scripts at the end (after the transaction completes)
	defer svc.Reload()

	return db.Transaction(func() (err error) {
		if err := svc.UniqueCheck(db, s); err != nil {
			return err
		}

		if err = svc.srepo.update(db, s); err != nil {
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

		if err = svc.trepo.mergeSet(db, s.tms, s.ID, s.triggers); err != nil {
			return
		}

		return
	})
}

// DeleteScript - modifies script's props, pushes to repo & updates scripts cache
func (svc service) DeleteScript(ctx context.Context, s *Script) (err error) {
	s.DeletedAt = &time.Time{}
	*s.DeletedAt = time.Now()
	s.DeletedBy = auth.GetIdentityFromContext(ctx).Identity()

	db := svc.db.With(ctx)

	// Reloading scripts at the end (after the transaction completes)
	defer svc.Reload()

	return db.Transaction(func() error {
		// We're doing soft delete in the repo
		if err = svc.srepo.update(db, s); err != nil {
			return err
		}

		if err = svc.trepo.deleteByScriptID(db, s.ID); err != nil {
			return err
		}

		return nil
	})
}

func (svc service) UniqueCheck(db *factory.DB, s *Script) (err error) {
	if s.Name != "" {
		f := ScriptFilter{
			NamespaceID: s.NamespaceID,
			Name:        s.Name,
			IncDeleted:  false,
		}

		ss, _, err := svc.srepo.find(db, f)
		if err != nil || len(ss) == 0 {
			return err
		}

		if len(ss) > 1 || ss.FindByID(s.ID) == nil {
			return errors.New("script name not unique")
		}
	}

	return nil
}

func (svc service) FindTriggerByID(ctx context.Context, scriptID uint64) (*Trigger, error) {
	return svc.trepo.findByID(svc.db.With(ctx), scriptID)
}

func (svc service) FindTriggers(ctx context.Context, f TriggerFilter) (TriggerSet, TriggerFilter, error) {
	return svc.trepo.find(svc.db.With(ctx), f)
}

// CreateTrigger - modifies script's props, pushes to repo & updates scripts cache
func (svc service) CreateTrigger(ctx context.Context, s *Script, t *Trigger) (err error) {
	if err = svc.setNewTriggerInfo(ctx, s, t); err != nil {
		return err
	}

	if err = svc.trepo.replace(svc.db.With(ctx), t); err != nil {
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

	if err = svc.trepo.replace(svc.db.With(ctx), t); err != nil {
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
	if err = svc.trepo.replace(svc.db.With(ctx), t); err != nil {
		return err
	}

	svc.Reload()
	return
}
