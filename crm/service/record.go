package service

import (
	"context"

	"github.com/pkg/errors"
	"github.com/titpetric/factory"

	"github.com/crusttech/crust/crm/repository"
	"github.com/crusttech/crust/crm/types"

	systemService "github.com/crusttech/crust/system/service"
)

type (
	record struct {
		db  *factory.DB
		ctx context.Context

		repository repository.RecordRepository
		pageRepo   repository.PageRepository
		moduleRepo repository.ModuleRepository

		userSvc systemService.UserService
	}

	RecordService interface {
		With(ctx context.Context) RecordService

		FindByID(moduleID uint64, recordID uint64) (*types.Record, error)

		Report(moduleID uint64, metrics, dimensions, filter string) (interface{}, error)
		Find(moduleID uint64, filter string, sort string, page int, perPage int) (*repository.FindResponse, error)

		Create(record *types.Record) (*types.Record, error)
		Update(record *types.Record) (*types.Record, error)
		DeleteByID(recordID uint64) error

		Fields(module *types.Module, record *types.Record) ([]*types.RecordValue, error)
	}
)

func Record() RecordService {
	return (&record{
		userSvc: systemService.DefaultUser,
	}).With(context.Background())
}

func (s *record) With(ctx context.Context) RecordService {
	db := repository.DB(ctx)
	return &record{
		db:         db,
		ctx:        ctx,
		repository: repository.Record(ctx, db),
		pageRepo:   repository.Page(ctx, db),
		moduleRepo: repository.Module(ctx, db),
		userSvc:    s.userSvc.With(ctx),
	}
}

func (s *record) FindByID(moduleID uint64, id uint64) (response *types.Record, err error) {
	var module *types.Module

	if module, err = s.moduleRepo.FindByID(moduleID); err != nil {
		return nil, err
	}

	if response, err = s.repository.FindByID(id); err != nil {
		return nil, err
	}
	return response, s.preload(module, response, "page", "user", "fields")
}

func (s *record) Report(moduleID uint64, metrics, dimensions, filter string) (interface{}, error) {
	return s.repository.Report(moduleID, metrics, dimensions, filter)
}

func (s *record) Find(moduleID uint64, filter string, sort string, page int, perPage int) (response *repository.FindResponse, err error) {
	var module *types.Module

	if module, err = s.moduleRepo.FindByID(moduleID); err != nil {
		return nil, err
	} else if response, err = s.repository.Find(module, filter, sort, page, perPage); err != nil {
		return nil, err
	} else if err := s.preloadAll(module, response.Records, "user", "fields"); err != nil {
		return nil, err
	}

	return response, nil
}

func (s *record) Create(mod *types.Record) (*types.Record, error) {
	response, err := s.repository.Create(mod)
	if err != nil {
		return nil, err
	}
	return response, s.preload(nil, response, "user", "fields")
}

func (s *record) Update(record *types.Record) (c *types.Record, err error) {
	validate := func() error {
		if record.ID == 0 {
			return errors.New("Error updating record: invalid ID")
		} else if c, err = s.repository.FindByID(record.ID); err != nil {
			return errors.Wrap(err, "Error while loading record for update")
		} else {
			record.CreatedAt = c.CreatedAt
		}

		return nil
	}

	if err = validate(); err != nil {
		return nil, err
	}

	return c, s.db.Transaction(func() (err error) {
		c, err = s.repository.Update(record)
		return
	})
}

func (s *record) Fields(module *types.Module, record *types.Record) ([]*types.RecordValue, error) {
	return s.repository.Fields(module, record)
}

func (s *record) DeleteByID(id uint64) error {
	return s.repository.DeleteByID(id)
}
