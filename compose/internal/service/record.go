package service

import (
	"context"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/titpetric/factory"

	"github.com/crusttech/crust/compose/internal/repository"
	"github.com/crusttech/crust/compose/types"
	"github.com/crusttech/crust/internal/auth"

	systemService "github.com/crusttech/crust/system/service"
)

type (
	record struct {
		db  *factory.DB
		ctx context.Context

		prmSvc  PermissionsService
		userSvc systemService.UserService

		recordRepo repository.RecordRepository
		moduleRepo repository.ModuleRepository
	}

	RecordService interface {
		With(ctx context.Context) RecordService

		FindByID(namespaceID, recordID uint64) (*types.Record, error)

		Report(namespaceID, moduleID uint64, metrics, dimensions, filter string) (interface{}, error)
		Find(filter types.RecordFilter) (set types.RecordSet, f types.RecordFilter, err error)

		Create(record *types.Record) (*types.Record, error)
		Update(record *types.Record) (*types.Record, error)

		DeleteByID(namespaceID, recordID uint64) error

		// Fields(module *types.Module, record *types.Record) ([]*types.RecordValue, error)
	}
)

func Record() RecordService {
	return (&record{
		prmSvc:  DefaultPermissions,
		userSvc: systemService.DefaultUser,
	}).With(context.Background())
}

func (svc *record) With(ctx context.Context) RecordService {
	db := repository.DB(ctx)
	return &record{
		db:  db,
		ctx: ctx,

		prmSvc:  svc.prmSvc.With(ctx),
		userSvc: systemService.User(ctx),

		recordRepo: repository.Record(ctx, db),
		moduleRepo: repository.Module(ctx, db),
	}
}

func (svc *record) FindByID(namespaceID, recordID uint64) (r *types.Record, err error) {
	if namespaceID == 0 {
		return nil, ErrNamespaceRequired
	}

	if r, err = svc.recordRepo.FindByID(namespaceID, recordID); err != nil {
		return
	}

	if !svc.prmSvc.CanReadRecord(r) {
		return nil, ErrNoReadPermissions.withStack()
	}

	if err = svc.preloadValues(r); err != nil {
		return
	}

	return
}

func (svc record) loadModule(namespaceID, moduleID uint64) (m *types.Module, err error) {
	if m, err = svc.moduleRepo.FindByID(namespaceID, moduleID); err != nil {
		return
	}

	if m.Fields, err = svc.moduleRepo.FindFields(m.ID); err != nil {
		return
	}

	if !svc.prmSvc.CanReadRecord(m) {
		return nil, ErrNoReadPermissions.withStack()
	}

	return
}

func (svc *record) Report(namespaceID, moduleID uint64, metrics, dimensions, filter string) (out interface{}, err error) {
	var m *types.Module
	if m, err = svc.loadModule(namespaceID, moduleID); err != nil {
		return
	}

	if !svc.prmSvc.CanReadRecord(m) {
		return nil, ErrNoReadPermissions.withStack()
	}

	return svc.recordRepo.Report(m, metrics, dimensions, filter)
}

func (svc *record) Find(filter types.RecordFilter) (set types.RecordSet, f types.RecordFilter, err error) {
	var m *types.Module
	if m, err = svc.loadModule(filter.NamespaceID, filter.ModuleID); err != nil {
		return
	}

	if !svc.prmSvc.CanReadRecord(m) {
		return nil, filter, ErrNoReadPermissions.withStack()
	}

	set, f, err = svc.recordRepo.Find(m, filter)
	if err != nil {
		return
	}

	if err = svc.preloadValues(set...); err != nil {
		return
	}

	return
}

func (svc *record) Create(mod *types.Record) (r *types.Record, err error) {
	if mod.NamespaceID == 0 {
		return nil, ErrNamespaceRequired
	}

	var m *types.Module
	if m, err = svc.loadModule(mod.NamespaceID, mod.ModuleID); err != nil {
		return
	}

	if !svc.prmSvc.CanCreateRecord(m) {
		return nil, ErrNoCreatePermissions.withStack()
	}

	if err = svc.sanitizeValues(m, mod.Values); err != nil {
		return
	}

	mod.OwnedBy = auth.GetIdentityFromContext(svc.ctx).Identity()
	mod.CreatedBy = mod.OwnedBy
	mod.CreatedAt = time.Now()

	return r, svc.db.Transaction(func() (err error) {
		if r, err = svc.recordRepo.Create(mod); err != nil {
			return
		}

		if err = svc.recordRepo.UpdateValues(r.ID, mod.Values); err != nil {
			return
		}

		if err = svc.preloadValues(r); err != nil {
			return
		}

		return
	})
}

func (svc *record) Update(mod *types.Record) (r *types.Record, err error) {
	if mod.ID == 0 {
		return nil, ErrInvalidID.withStack()
	}

	if mod.NamespaceID == 0 {
		return nil, ErrNamespaceRequired
	}

	var m *types.Module
	if m, err = svc.loadModule(mod.NamespaceID, mod.ModuleID); err != nil {
		return
	}

	if r, err = svc.recordRepo.FindByID(mod.NamespaceID, mod.ID); err != nil {
		return
	}

	if !svc.prmSvc.CanUpdateRecord(r) {
		return nil, ErrNoUpdatePermissions.withStack()
	}

	if isStale(mod.UpdatedAt, r.UpdatedAt, r.CreatedAt) {
		return nil, ErrStaleData.withStack()
	}

	if err = svc.sanitizeValues(m, mod.Values); err != nil {
		return
	}

	now := time.Now()
	r.UpdatedAt = &now
	r.UpdatedBy = auth.GetIdentityFromContext(svc.ctx).Identity()

	return r, svc.db.Transaction(func() (err error) {
		if r, err = svc.recordRepo.Update(r); err != nil {
			return
		}

		if err = svc.recordRepo.UpdateValues(r.ID, mod.Values); err != nil {
			return
		}

		return
	})
}

func (svc *record) DeleteByID(namespaceID, recordID uint64) (err error) {
	err = svc.db.Transaction(func() (err error) {
		var record *types.Record

		if record, err = svc.recordRepo.FindByID(namespaceID, recordID); err != nil {
			return errors.Wrap(err, "nonexistent record")
		}

		now := time.Now()
		record.DeletedAt = &now
		record.DeletedBy = auth.GetIdentityFromContext(svc.ctx).Identity()

		if err = svc.recordRepo.Delete(record); err != nil {
			return
		}

		if err = svc.recordRepo.DeleteValues(record); err != nil {
			return
		}

		return
	})

	return errors.Wrap(err, "unable to delete record")
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
	if rvs, err := svc.recordRepo.LoadValues(types.RecordSet(rr).IDs()...); err != nil {
		return err
	} else {
		return types.RecordSet(rr).Walk(func(r *types.Record) error {
			r.Values = rvs.FilterByRecordID(r.ID)
			return nil
		})
	}
}
