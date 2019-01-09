package service

import (
	"context"
	"errors"

	"github.com/titpetric/factory"

	"github.com/crusttech/crust/crm/repository"
	"github.com/crusttech/crust/crm/types"
)

type (
	trigger struct {
		db  *factory.DB
		ctx context.Context

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
	return (&trigger{}).With(context.Background())
}

func (s *trigger) With(ctx context.Context) TriggerService {
	db := repository.DB(ctx)
	return &trigger{
		db:          db,
		ctx:         ctx,
		triggerRepo: repository.Trigger(ctx, db),
		moduleRepo:  repository.Module(ctx, db),
	}
}

func (s *trigger) FindByID(id uint64) (*types.Trigger, error) {
	return s.triggerRepo.FindByID(id)
}

func (s *trigger) Find(filter types.TriggerFilter) (types.TriggerSet, error) {
	return s.triggerRepo.Find(filter)
}

func (s *trigger) Create(trigger *types.Trigger) (p *types.Trigger, err error) {
	validate := func() error {
		return nil
	}
	if err := validate(); err != nil {
		return nil, err
	}
	return p, s.db.Transaction(func() (err error) {
		p, err = s.triggerRepo.Create(trigger)
		return
	})
}

func (s *trigger) Update(trigger *types.Trigger) (p *types.Trigger, err error) {
	validate := func() error {
		if trigger.ID == 0 {
			return errors.New("could not update trigger, invalid ID")
		}
		return nil
	}
	if err := validate(); err != nil {
		return nil, err
	}
	return p, s.db.Transaction(func() (err error) {
		p, err = s.triggerRepo.Update(trigger)
		return
	})
}

func (s *trigger) DeleteByID(id uint64) error {
	return s.triggerRepo.DeleteByID(id)
}
