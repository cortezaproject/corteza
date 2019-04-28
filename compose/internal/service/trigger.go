package service

import (
	"context"

	"github.com/titpetric/factory"

	"github.com/crusttech/crust/compose/internal/repository"
	"github.com/crusttech/crust/compose/types"
)

type (
	trigger struct {
		db  *factory.DB
		ctx context.Context

		prmSvc PermissionsService

		triggerRepo repository.TriggerRepository
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
		prmSvc: DefaultPermissions,
	}).With(context.Background())
}

func (svc trigger) With(ctx context.Context) TriggerService {
	db := repository.DB(ctx)
	return &trigger{
		db:  db,
		ctx: ctx,

		prmSvc: svc.prmSvc.With(ctx),

		triggerRepo: repository.Trigger(ctx, db),
	}
}

func (svc trigger) FindByID(namespaceID, triggerID uint64) (c *types.Trigger, err error) {
	if namespaceID == 0 {
		return nil, ErrNamespaceRequired
	}

	if c, err = svc.triggerRepo.FindByID(namespaceID, triggerID); err != nil {
		return
	} else if !svc.prmSvc.CanReadTrigger(c) {
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
		return svc.prmSvc.CanReadTrigger(m), nil
	})

	return
}

func (svc trigger) Create(mod *types.Trigger) (c *types.Trigger, err error) {
	if !svc.prmSvc.CanCreateTrigger(crmNamespace()) {
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

	if !svc.prmSvc.CanUpdateTrigger(c) {
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
	if namespaceID == 0 {
		return ErrNamespaceRequired.withStack()
	}

	if c, err := svc.triggerRepo.FindByID(namespaceID, triggerID); err != nil {
		return err
	} else if !svc.prmSvc.CanDeleteTrigger(c) {
		return ErrNoDeletePermissions.withStack()
	}

	return svc.triggerRepo.DeleteByID(namespaceID, triggerID)
}
