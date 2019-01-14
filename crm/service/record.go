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
		moduleRepo repository.ModuleRepository

		userSvc systemService.UserService
	}

	RecordService interface {
		With(ctx context.Context) RecordService

		FindByID(recordID uint64) (*types.Record, error)

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

func (svc *record) With(ctx context.Context) RecordService {
	db := repository.DB(ctx)
	return &record{
		db:         db,
		ctx:        ctx,
		repository: repository.Record(ctx, db),
		moduleRepo: repository.Module(ctx, db),
		userSvc:    svc.userSvc.With(ctx),
	}
}

func (svc *record) FindByID(recordID uint64) (r *types.Record, err error) {
	err = svc.db.Transaction(func() (err error) {
		if r, err = svc.repository.FindByID(recordID); err != nil {
			return
		}

		if err = svc.preloadValues(r); err != nil {
			return
		}

		if err = svc.preloadUsers(r); err != nil {
			return
		}

		return
	})

	return r, errors.Wrap(err, "unable to find record")
}

func (svc *record) Report(moduleID uint64, metrics, dimensions, filter string) (out interface{}, err error) {
	var module *types.Module

	err = svc.db.Transaction(func() (err error) {
		if module, err = svc.moduleRepo.FindByID(moduleID); err != nil {
			return
		}

		out, err = svc.repository.Report(module, metrics, dimensions, filter)
		return
	})

	return out, errors.Wrap(err, "unable to build a report")
}

func (svc *record) Find(moduleID uint64, filter, sort string, page, perPage int) (rsp *repository.FindResponse, err error) {
	var module *types.Module

	err = svc.db.Transaction(func() (err error) {
		if module, err = svc.moduleRepo.FindByID(moduleID); err != nil {
			return
		}

		if rsp, err = svc.repository.Find(module, filter, sort, page, perPage); err != nil {
			return
		}

		if err = svc.preloadValues(rsp.Records...); err != nil {
			return
		}

		if err = svc.preloadUsers(rsp.Records...); err != nil {
			return
		}

		return
	})

	return rsp, errors.Wrap(err, "unable to find records")

}

func (svc *record) Create(new *types.Record) (record *types.Record, err error) {
	var module *types.Module

	err = svc.db.Transaction(func() (err error) {
		if module, err = svc.moduleRepo.FindByID(new.ModuleID); err != nil {
			return
		}

		if err = svc.sanitizeValues(module, new.Values); err != nil {
			return
		}

		if record, err = svc.repository.Create(new); err != nil {
			return
		}

		if err = svc.repository.UpdateValues(record.ID, new.Values); err != nil {
			return
		}

		if err = svc.preloadValues(record); err != nil {
			return
		}

		if err = svc.preloadUsers(record); err != nil {
			return
		}

		return
	})

	return record, errors.Wrap(err, "unable to create record")
}

func (svc *record) Update(updated *types.Record) (record *types.Record, err error) {
	var module *types.Module

	err = svc.db.Transaction(func() (err error) {
		if updated.ID == 0 {
			return errors.New("invalid record ID")
		}

		if record, err = svc.repository.FindByID(updated.ID); err != nil {
			return errors.Wrap(err, "nonexistent record")
		}

		updated.CreatedAt = record.CreatedAt
		updated.UserID = record.UserID

		if module, err = svc.moduleRepo.FindByID(updated.ModuleID); err != nil {
			return
		}

		if err = svc.sanitizeValues(module, updated.Values); err != nil {
			return
		}

		if record, err = svc.repository.Update(updated); err != nil {
			return
		}

		if err = svc.repository.UpdateValues(record.ID, updated.Values); err != nil {
			return
		}

		if err = svc.preloadUsers(record); err != nil {
			return
		}

		return
	})

	return record, errors.Wrap(err, "unable to update record")
}

// func (s *record) Fields(module *types.Module, record *types.Record) ([]*types.RecordValue, error) {
// 	return s.repository.Fields(module, record)
// }

func (svc *record) DeleteByID(id uint64) error {
	return svc.repository.DeleteByID(id)
}

// Validates and filters record values
func (svc *record) sanitizeValues(module *types.Module, values types.RecordValueSet) (err error) {
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

	var places = map[string]uint{}
	// var has bool

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

		value.Place = places[field.Name]
		places[field.Name]++

		return nil
	})
}

func (svc *record) preloadValues(rr ...*types.Record) error {
	if rvs, err := svc.repository.LoadValues(types.RecordSet(rr).IDs()...); err != nil {
		return err
	} else {
		return types.RecordSet(rr).Walk(func(r *types.Record) error {
			r.Values = rvs.FilterByRecordID(r.ID)
			return nil
		})
	}
}

func (svc *record) preloadUsers(rr ...*types.Record) error {
	if uu, err := svc.userSvc.FindByIDs(types.RecordSet(rr).UserIDs()...); err != nil {
		return err
	} else {
		return types.RecordSet(rr).Walk(func(r *types.Record) error {
			r.User = uu.FindByID(r.UserID)
			return nil
		})
	}
}
