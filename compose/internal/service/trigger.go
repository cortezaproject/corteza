package service

import (
	"context"

	"github.com/titpetric/factory"
	"go.uber.org/zap"

	"github.com/cortezaproject/corteza-server/compose/internal/repository"
	"github.com/cortezaproject/corteza-server/compose/types"
)

type (
	trigger struct {
		db     *factory.DB
		ctx    context.Context
		logger *zap.Logger

		ac triggerAccessController

		triggerRepo repository.TriggerRepository
		nsRepo      repository.NamespaceRepository
	}

	triggerAccessController interface {
		CanReadNamespace(context.Context, *types.Namespace) bool
		CanCreateTrigger(context.Context, *types.Namespace) bool
		CanReadTrigger(context.Context, *types.Trigger) bool
		CanUpdateTrigger(context.Context, *types.Trigger) bool
		CanDeleteTrigger(context.Context, *types.Trigger) bool
	}

	TriggerService interface {
		With(ctx context.Context) TriggerService

		FindByID(namespaceID, triggerID uint64) (*types.Trigger, error)
		Find(filter types.TriggerFilter) (set types.TriggerSet, f types.TriggerFilter, err error)

		Create(trigger *types.Trigger) (*types.Trigger, error)
		Update(trigger *types.Trigger) (*types.Trigger, error)
		DeleteByID(namespaceID, triggerID uint64) error
	}
)

func Trigger() TriggerService {
	return (&trigger{
		logger: DefaultLogger.Named("trigger"),
		ac:     DefaultAccessControl,
	}).With(context.Background())
}

func (svc trigger) With(ctx context.Context) TriggerService {
	db := repository.DB(ctx)
	return &trigger{
		db:     db,
		ctx:    ctx,
		logger: svc.logger,

		ac: svc.ac,

		triggerRepo: repository.Trigger(ctx, db),
		nsRepo:      repository.Namespace(ctx, db),
	}
}

// log() returns zap's logger with requestID from current context and fields.
// func (svc trigger) log(fields ...zapcore.Field) *zap.Logger {
// 	return logger.AddRequestID(svc.ctx, svc.logger).With(fields...)
// }

func (svc trigger) FindByID(namespaceID, triggerID uint64) (t *types.Trigger, err error) {
	if namespaceID == 0 {
		return nil, ErrNamespaceRequired
	}

	if t, err = svc.triggerRepo.FindByID(namespaceID, triggerID); err != nil {
		return
	} else if !svc.ac.CanReadTrigger(svc.ctx, t) {
		return nil, ErrNoReadPermissions.withStack()
	}

	return
}

func (svc trigger) Find(filter types.TriggerFilter) (set types.TriggerSet, f types.TriggerFilter, err error) {
	set, f, err = svc.triggerRepo.Find(filter)
	if err != nil {
		return
	}

	set, _ = set.Filter(func(m *types.Trigger) (bool, error) {
		return svc.ac.CanReadTrigger(svc.ctx, m), nil
	})

	return
}

func (svc trigger) Create(mod *types.Trigger) (t *types.Trigger, err error) {
	if mod.NamespaceID == 0 {
		return nil, ErrNamespaceRequired.withStack()
	}

	if ns, err := svc.loadNamespace(mod.NamespaceID); err != nil {
		return nil, err
	} else if !svc.ac.CanCreateTrigger(svc.ctx, ns) {
		return nil, ErrNoCreatePermissions.withStack()
	}

	return svc.triggerRepo.Create(mod)
}

func (svc trigger) Update(mod *types.Trigger) (c *types.Trigger, err error) {
	if mod.ID == 0 {
		return nil, ErrInvalidID.withStack()
	}

	if c, err = svc.triggerRepo.FindByID(mod.NamespaceID, mod.ID); err != nil {
		return
	}

	if isStale(mod.UpdatedAt, c.UpdatedAt, c.CreatedAt) {
		return nil, ErrStaleData.withStack()
	}

	if !svc.ac.CanUpdateTrigger(svc.ctx, c) {
		return nil, ErrNoUpdatePermissions.withStack()
	}

	c.Name = mod.Name
	c.ModuleID = mod.ModuleID
	c.Source = mod.Source
	c.Actions = mod.Actions
	c.Enabled = mod.Enabled

	return svc.triggerRepo.Update(c)
}

func (svc trigger) DeleteByID(namespaceID, triggerID uint64) error {
	if triggerID == 0 {
		return ErrInvalidID.withStack()
	}

	if namespaceID == 0 {
		return ErrNamespaceRequired.withStack()
	}

	if c, err := svc.triggerRepo.FindByID(namespaceID, triggerID); err != nil {
		return err
	} else if !svc.ac.CanDeleteTrigger(svc.ctx, c) {
		return ErrNoDeletePermissions.withStack()
	}

	return svc.triggerRepo.DeleteByID(namespaceID, triggerID)
}

func (svc trigger) loadNamespace(namespaceID uint64) (ns *types.Namespace, err error) {
	if namespaceID == 0 {
		return nil, ErrNamespaceRequired.withStack()
	}

	if ns, err = svc.nsRepo.FindByID(namespaceID); err != nil {
		return
	}

	if !svc.ac.CanReadNamespace(svc.ctx, ns) {
		return nil, ErrNoReadPermissions.withStack()
	}

	return
}
