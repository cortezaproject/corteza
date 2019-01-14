package service

import (
	"context"
	"strconv"

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

		// Fields(module *types.Module, record *types.Record) ([]*types.RecordValue, error)
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

func (s *record) Create(new *types.Record) (record *types.Record, err error) {
	var module *types.Module

	err = s.db.Transaction(func() (err error) {
		if module, err = s.moduleRepo.FindByID(new.ModuleID); err != nil {
			return
		}

		if err = s.sanitizeValues(module, new.Values); err != nil {
			return
		}

		if record, err = s.repository.Create(new); err != nil {
			return
		}

		if err = s.repository.UpdateValues(record.ID, new.Values); err != nil {
			return
		}

		return s.preload(module, record, "user")
	})

	return record, errors.Wrap(err, "unable to create record")
}

func (s *record) Update(updated *types.Record) (record *types.Record, err error) {
	var module *types.Module

	err = s.db.Transaction(func() (err error) {
		if updated.ID == 0 {
			return errors.New("invalid record ID")
		}

		if record, err = s.repository.FindByID(updated.ID); err != nil {
			return errors.Wrap(err, "unexisting record")
		}

		updated.CreatedAt = record.CreatedAt
		updated.UserID = record.UserID

		if module, err = s.moduleRepo.FindByID(updated.ModuleID); err != nil {
			return
		}

		if err = s.sanitizeValues(module, updated.Values); err != nil {
			return
		}

		if record, err = s.repository.Update(updated); err != nil {
			return
		}

		if err = s.repository.UpdateValues(record.ID, updated.Values); err != nil {
			return
		}

		return s.preload(module, record, "user")
	})

	return record, errors.Wrap(err, "unable to update record")
}

// func (s *record) Fields(module *types.Module, record *types.Record) ([]*types.RecordValue, error) {
// 	return s.repository.Fields(module, record)
// }

func (s *record) DeleteByID(id uint64) error {
	return s.repository.DeleteByID(id)
}

// Validates and filters record values
func (s *record) sanitizeValues(module *types.Module, values types.RecordValueSet) (err error) {
	// Make sure there are no multi values in a non-multi value fields
	err = module.Fields.Walk(func(field *types.ModuleField) error {
		if !field.Multi && len(values.FilterByName(field.Name)) > 1 {
			return errors.Errorf("more than one value for a single-value field %q", field.Name)
		}
		return nil
	})

	if err != nil {
		return
	}

	return values.Walk(func(value *types.RecordValue) (err error) {
		var field = module.Fields.FindByName(value.Name)
		if field == nil {
			return errors.Errorf("no such field %q", value.Name)
		}

		if field.IsRef() {
			if value.Ref, err = strconv.ParseUint(value.Value, 10, 64); err != nil {
				return err
			}
		}

		return nil
	})
}
