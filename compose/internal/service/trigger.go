package service

import (
	"context"

	"github.com/pkg/errors"
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
		moduleRepo  repository.ModuleRepository
	}

	TriggerService interface {
		With(ctx context.Context) TriggerService

		FindByID(triggerID uint64) (*types.Trigger, error)
		Find(filter types.TriggerFilter) (set types.TriggerSet, err error)

		Create(trigger *types.Trigger) (*types.Trigger, error)
		Update(trigger *types.Trigger) (*types.Trigger, error)
		DeleteByID(triggerID uint64) error
	}
)

func Trigger() TriggerService {
	return (&trigger{
		prmSvc: DefaultPermissions}).With(context.Background())
}

func (svc *trigger) With(ctx context.Context) TriggerService {
	db := repository.DB(ctx)
	return &trigger{
		db:  db,
		ctx: ctx,

		prmSvc: svc.prmSvc.With(ctx),

		triggerRepo: repository.Trigger(ctx, db),
		moduleRepo:  repository.Module(ctx, db),
	}
}

func (svc *trigger) FindByID(id uint64) (t *types.Trigger, err error) {
	if t, err = svc.triggerRepo.FindByID(id); err != nil {
		return
	} else if !svc.prmSvc.CanReadTrigger(t) {
		return nil, errors.New("not allowed to access this trigger")
	}

	return
}

func (svc *trigger) Find(filter types.TriggerFilter) (tt types.TriggerSet, err error) {
	if tt, err = svc.triggerRepo.Find(filter); err != nil {
		return nil, err
	} else {
		return tt.Filter(func(m *types.Trigger) (bool, error) {
			return svc.prmSvc.CanReadTrigger(m), nil
		})
	}
}

func (svc *trigger) Create(trigger *types.Trigger) (p *types.Trigger, err error) {
	if !svc.prmSvc.CanCreateTrigger(crmNamespace()) {
		return nil, errors.New("not allowed to create this trigger")
	}

	return p, svc.db.Transaction(func() (err error) {
		p, err = svc.triggerRepo.Create(trigger)
		return
	})
}

func (svc *trigger) Update(trigger *types.Trigger) (t *types.Trigger, err error) {
	validate := func() error {
		if trigger.ID == 0 {
			return errors.New("Error updating trigger: invalid ID")
		} else if t, err = svc.triggerRepo.FindByID(trigger.ID); err != nil {
			return errors.Wrap(err, "Error while loading trigger for update")
		} else {
			if !svc.prmSvc.CanUpdateModule(t) {
				return errors.New("not allowed to update this trigger")
			}

			trigger.CreatedAt = t.CreatedAt
		}

		return nil
	}

	if err := validate(); err != nil {
		return nil, err
	}

	return t, svc.db.Transaction(func() (err error) {
		t, err = svc.triggerRepo.Update(trigger)
		return
	})
}

func (svc *trigger) DeleteByID(ID uint64) error {
	if t, err := svc.triggerRepo.FindByID(ID); err != nil {
		return errors.Wrap(err, "could not delete trigger")
	} else if !svc.prmSvc.CanDeleteTrigger(t) {
		return errors.New("not allowed to delete this trigger")
	}

	return svc.triggerRepo.DeleteByID(ID)
}
