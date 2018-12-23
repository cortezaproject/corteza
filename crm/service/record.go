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

		userSvc systemService.UserService
	}

	RecordService interface {
		With(ctx context.Context) RecordService

		FindByID(recordID uint64) (*types.Record, error)

		Report(moduleID uint64, params *types.RecordReport) (interface{}, error)
		Find(moduleID uint64, query string, page int, perPage int, sort string) (*repository.FindResponse, error)

		Create(record *types.Record) (*types.Record, error)
		Update(record *types.Record) (*types.Record, error)
		DeleteByID(recordID uint64) error

		Fields(mod *types.Record) ([]*types.RecordColumn, error)
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
		userSvc:    s.userSvc.With(ctx),
	}
}

func (s *record) FindByID(id uint64) (*types.Record, error) {
	response, err := s.repository.FindByID(id)
	if err != nil {
		return nil, err
	}
	return response, s.preload(response, "page", "user", "fields")
}

func (s *record) Report(moduleID uint64, params *types.RecordReport) (interface{}, error) {
	return s.repository.Report(moduleID, params)
}

func (s *record) Find(moduleID uint64, query string, page int, perPage int, sort string) (*repository.FindResponse, error) {
	response, err := s.repository.Find(moduleID, query, page, perPage, sort)
	if err != nil {
		return nil, err
	}
	if err := s.preloadAll(response.Records, "user", "fields"); err != nil {
		return nil, err
	}
	return response, nil
}

func (s *record) Create(mod *types.Record) (*types.Record, error) {
	response, err := s.repository.Create(mod)
	if err != nil {
		return nil, err
	}
	return response, s.preload(response, "user", "fields")
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

func (s *record) Fields(mod *types.Record) ([]*types.RecordColumn, error) {
	return s.repository.Fields(mod)
}

func (s *record) DeleteByID(id uint64) error {
	return s.repository.DeleteByID(id)
}
