package service

import (
	"context"
	"regexp"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/titpetric/factory"
	"go.uber.org/zap"

	"github.com/cortezaproject/corteza-server/compose/internal/repository"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/internal/auth"
)

type (
	record struct {
		db     *factory.DB
		ctx    context.Context
		logger *zap.Logger

		ac recordAccessController
		sr *scriptRunner

		recordRepo repository.RecordRepository
		moduleRepo repository.ModuleRepository
		nsRepo     repository.NamespaceRepository
		tRepo      repository.TriggerRepository
	}

	recordAccessController interface {
		CanCreateRecord(context.Context, *types.Module) bool
		CanReadNamespace(context.Context, *types.Namespace) bool
		CanReadModule(context.Context, *types.Module) bool
		CanReadRecord(context.Context, *types.Module) bool
		CanUpdateRecord(context.Context, *types.Module) bool
		CanDeleteRecord(context.Context, *types.Module) bool
		CanReadRecordValue(context.Context, *types.ModuleField) bool
		CanUpdateRecordValue(context.Context, *types.ModuleField) bool
	}

	RecordService interface {
		With(ctx context.Context) RecordService

		FindByID(namespaceID, recordID uint64) (*types.Record, error)

		Report(namespaceID, moduleID uint64, metrics, dimensions, filter string) (interface{}, error)
		Find(filter types.RecordFilter) (set types.RecordSet, f types.RecordFilter, err error)
		Export(types.RecordFilter, Encoder) error

		Create(record *types.Record) (*types.Record, error)
		Update(record *types.Record) (*types.Record, error)

		DeleteByID(namespaceID, recordID uint64) error

		// Fields(module *types.Module, record *types.Record) ([]*types.RecordValue, error)
	}

	Encoder interface {
		Record(*types.Record) error
	}
)

func Record() RecordService {
	return (&record{
		logger: DefaultLogger.Named("record"),
		ac:     DefaultAccessControl,
		sr:     DefaultScriptRunner,
	}).With(context.Background())
}

func (svc record) With(ctx context.Context) RecordService {
	db := repository.DB(ctx)

	return &record{
		db:     db,
		ctx:    ctx,
		logger: svc.logger,

		ac: svc.ac,
		sr: svc.sr,

		recordRepo: repository.Record(ctx, db),
		moduleRepo: repository.Module(ctx, db),
		nsRepo:     repository.Namespace(ctx, db),
		tRepo:      repository.Trigger(ctx, db),
	}
}

// log() returns zap's logger with requestID from current context and fields.
// func (svc record) log(fields ...zapcore.Field) *zap.Logger {
// 	return logger.AddRequestID(svc.ctx, svc.logger).With(fields...)
// }

func (svc record) FindByID(namespaceID, recordID uint64) (r *types.Record, err error) {
	if namespaceID == 0 {
		return nil, ErrNamespaceRequired
	}

	if r, err = svc.recordRepo.FindByID(namespaceID, recordID); err != nil {
		return
	}

	var m *types.Module
	if m, err = svc.loadModule(namespaceID, r.ModuleID); err != nil {
		return
	}

	if !svc.ac.CanReadRecord(svc.ctx, m) {
		return nil, ErrNoReadPermissions.withStack()
	}

	if err = svc.preloadValues(m, r); err != nil {
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

	if !svc.ac.CanReadModule(svc.ctx, m) {
		return nil, ErrNoReadPermissions.withStack()
	}

	return
}

func (svc record) loadNamespace(namespaceID uint64) (ns *types.Namespace, err error) {
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

func (svc record) Report(namespaceID, moduleID uint64, metrics, dimensions, filter string) (out interface{}, err error) {
	var m *types.Module
	if m, err = svc.loadModule(namespaceID, moduleID); err != nil {
		return
	}

	return svc.recordRepo.
		Report(m, metrics, dimensions, filter)
}

func (svc record) Find(filter types.RecordFilter) (set types.RecordSet, f types.RecordFilter, err error) {
	var m *types.Module
	if m, err = svc.loadModule(filter.NamespaceID, filter.ModuleID); err != nil {
		return
	}

	set, f, err = svc.recordRepo.Find(m, filter)
	if err != nil {
		return
	}

	if err = svc.preloadValues(m, set...); err != nil {
		return
	}

	return
}

// Export returns all records
//
// @todo better value handling
func (svc record) Export(filter types.RecordFilter, enc Encoder) error {
	m, err := svc.loadModule(filter.NamespaceID, filter.ModuleID)
	if err != nil {
		return err
	}

	set, err := svc.recordRepo.Export(m, filter)
	if err != nil {
		return err
	}

	if err = svc.preloadValues(m, set...); err != nil {
		return err
	}

	return set.Walk(enc.Record)
}

func (svc record) Create(mod *types.Record) (r *types.Record, err error) {
	ns, m, r, tt, err := svc.loadCombo(mod.NamespaceID, mod.ModuleID, 0)
	if err != nil {
		return
	}

	if !svc.ac.CanCreateRecord(svc.ctx, m) {
		return nil, ErrNoCreatePermissions.withStack()
	}

	creatorID := auth.GetIdentityFromContext(svc.ctx).Identity()
	r = &types.Record{
		ModuleID:    mod.ModuleID,
		NamespaceID: mod.NamespaceID,

		CreatedBy: creatorID,
		OwnedBy:   creatorID,

		CreatedAt: time.Now(),
	}

	if err = svc.copyChanges(m, mod, r); err != nil {
		return
	}

	mod = nil // make sure we do not use it anymore

	if err = tt.WalkByAction("beforeCreate", svc.runTrigger(svc.ctx, ns, m, r)); err != nil {
		return
	}

	defer func() {
		_ = tt.WalkByAction("afterCreate", svc.runTrigger(svc.ctx, ns, m, r))
	}()

	return r, svc.db.Transaction(func() (err error) {
		if r, err = svc.recordRepo.Create(r); err != nil {
			return
		}

		if err = svc.recordRepo.UpdateValues(r.ID, r.Values); err != nil {
			return
		}

		return
	})
}

func (svc record) Update(mod *types.Record) (r *types.Record, err error) {
	if mod.ID == 0 {
		return nil, ErrInvalidID.withStack()
	}

	ns, m, r, tt, err := svc.loadCombo(mod.NamespaceID, mod.ModuleID, mod.ID)
	if err != nil {
		return
	}

	if !svc.ac.CanUpdateRecord(svc.ctx, m) {
		return nil, ErrNoUpdatePermissions.withStack()
	}

	// Test if stale (update has an older copy)
	if isStale(mod.UpdatedAt, r.UpdatedAt, r.CreatedAt) {
		return nil, ErrStaleData.withStack()
	}

	now := time.Now()
	r.UpdatedAt = &now
	r.UpdatedBy = auth.GetIdentityFromContext(svc.ctx).Identity()

	if err = svc.copyChanges(m, mod, r); err != nil {
		return
	}

	mod = nil // make sure we do not use it anymore

	if err = tt.WalkByAction("beforeUpdate", svc.runTrigger(svc.ctx, ns, m, r)); err != nil {
		return
	}

	defer func() {
		_ = tt.WalkByAction("afterUpdate", svc.runTrigger(svc.ctx, ns, m, r))
	}()

	return r, svc.db.Transaction(func() (err error) {
		if r, err = svc.recordRepo.Update(r); err != nil {
			return
		}

		if err = svc.recordRepo.UpdateValues(r.ID, r.Values); err != nil {
			return
		}

		return
	})
}

func (svc record) DeleteByID(namespaceID, recordID uint64) (err error) {
	if recordID == 0 {
		return ErrInvalidID.withStack()
	}

	ns, m, r, tt, err := svc.loadCombo(namespaceID, 0, recordID)
	if err != nil {
		return
	}

	if err = tt.WalkByAction("beforeDelete", svc.runTrigger(svc.ctx, ns, m, r)); err != nil {
		return
	}

	defer func() {
		_ = tt.WalkByAction("afterDelete", svc.runTrigger(svc.ctx, ns, m, r))
	}()

	err = svc.db.Transaction(func() (err error) {
		now := time.Now()
		r.DeletedAt = &now
		r.DeletedBy = auth.GetIdentityFromContext(svc.ctx).Identity()

		if err = svc.recordRepo.Delete(r); err != nil {
			return
		}

		if err = svc.recordRepo.DeleteValues(r); err != nil {
			return
		}

		return
	})

	return errors.Wrap(err, "unable to delete record")
}

func (svc record) loadCombo(namespaceID, moduleID, recordID uint64) (ns *types.Namespace, m *types.Module, r *types.Record, tt types.TriggerSet, err error) {
	if namespaceID == 0 {
		err = ErrNamespaceRequired
		return
	}
	if ns, err = svc.loadNamespace(namespaceID); err != nil {
		return
	}

	if recordID > 0 {
		if r, err = svc.recordRepo.FindByID(namespaceID, recordID); err != nil {
			return
		}

		moduleID = r.ModuleID
	}

	if m, err = svc.loadModule(ns.ID, moduleID); err != nil {
		return
	}

	tt, _, err = svc.tRepo.Find(types.TriggerFilter{
		NamespaceID: ns.ID,
		ModuleID:    m.ID,
	})

	return
}

func (svc record) runTrigger(ctx context.Context, ns *types.Namespace, m *types.Module, r *types.Record) func(t *types.Trigger) error {
	svc.logger.Debug("initializing trigger runner")
	return func(t *types.Trigger) error {
		svc.logger.Debug("running trigger", zap.Uint64("triggerID", t.ID))
		if svc.sr == nil {
			// No script runner set
			svc.logger.Debug("script runner not set")
			return nil
		}

		// pr == processed record
		pr, err := svc.sr.Record(svc.ctx, t, ns, m, r)
		if err != nil {
			svc.logger.Debug("failed to run record script", zap.Error(err))
			return err
		}

		if pr == nil {
			// Did not get any processed record,
			// consider canceled
			return errors.New("aborted by automation")
		}

		return svc.copyChanges(m, pr, r)
	}
}

// Copies changes from mod to r(ecord)
func (svc record) copyChanges(m *types.Module, mod, r *types.Record) (err error) {
	// Automation scripts are allowed to modify record owner & values.
	if mod.OwnedBy > 0 {
		r.OwnedBy = mod.OwnedBy
	}

	r.Values, err = svc.sanitizeValues(m, mod.Values)
	return err
}

// Validates and filters record values
func (svc record) sanitizeValues(module *types.Module, values types.RecordValueSet) (out types.RecordValueSet, err error) {
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

	// Remove all values on un-updatable fields
	out, err = values.Filter(func(v *types.RecordValue) (bool, error) {
		var field = module.Fields.FindByName(v.Name)
		if field == nil {
			return false, errors.Errorf("no such field %q", v.Name)
		}

		if field.IsRef() && v.Value == "" {
			// Skip empty values on ref fields
			return false, nil
		}

		return svc.ac.CanUpdateRecordValue(svc.ctx, field), nil
	})

	if err != nil {
		return
	}

	var (
		places  = map[string]uint{}
		numeric = regexp.MustCompile(`^(\d+)$`)
	)

	return out, out.Walk(func(value *types.RecordValue) (err error) {
		var field = module.Fields.FindByName(value.Name)
		if field == nil {
			return errors.Errorf("no such field %q", value.Name)
		}

		if field.IsRef() {
			if !numeric.MatchString(value.Value) {
				return errors.Errorf("invalid reference format")
			}

			if value.Ref, err = strconv.ParseUint(value.Value, 10, 64); err != nil {
				return err
			}
		}

		value.Place = places[field.Name]
		places[field.Name]++

		return nil
	})
}

func (svc record) preloadValues(m *types.Module, rr ...*types.Record) error {
	if rvs, err := svc.recordRepo.LoadValues(svc.readableFields(m), types.RecordSet(rr).IDs()); err != nil {
		return err
	} else {
		return types.RecordSet(rr).Walk(func(r *types.Record) error {
			r.Values = rvs.FilterByRecordID(r.ID)
			return nil
		})
	}
}

// readableFields creates a slice of module fields that current user has permission to read
func (svc record) readableFields(m *types.Module) []string {
	ff := make([]string, 0)

	_ = m.Fields.Walk(func(f *types.ModuleField) error {
		if svc.ac.CanReadRecordValue(svc.ctx, f) {
			ff = append(ff, f.Name)
		}

		return nil
	})

	return ff
}
