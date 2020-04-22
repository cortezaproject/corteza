package service

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/cortezaproject/corteza-server/compose/service/values"

	"github.com/pkg/errors"
	"github.com/titpetric/factory"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/cortezaproject/corteza-server/compose/decoder"
	"github.com/cortezaproject/corteza-server/compose/repository"
	"github.com/cortezaproject/corteza-server/compose/service/event"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/cortezaproject/corteza-server/pkg/logger"
)

const (
	IMPORT_ON_ERROR_SKIP = "SKIP"
	IMPORT_ON_ERROR_FAIL = "FAIL"
)

type (
	record struct {
		db     *factory.DB
		ctx    context.Context
		logger *zap.Logger

		ac       recordAccessController
		eventbus eventDispatcher

		recordRepo repository.RecordRepository
		moduleRepo repository.ModuleRepository
		nsRepo     repository.NamespaceRepository

		formatter recordValuesFormatter
		sanitizer recordValuesSanitizer
		validator recordValuesValidator
	}

	recordValuesFormatter interface {
		Run(*types.Module, types.RecordValueSet) types.RecordValueSet
	}

	recordValuesSanitizer interface {
		Run(*types.Module, types.RecordValueSet) types.RecordValueSet
	}

	recordValuesValidator interface {
		Run(*types.Module, *types.Record) *types.RecordValueErrorSet
		UniqueChecker(fn values.UniqueChecker)
		RecordRefChecker(fn values.ReferenceChecker)
		UserRefChecker(fn values.ReferenceChecker)
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
		Import(*RecordImportSession, ImportSessionService) error

		Create(record *types.Record) (*types.Record, error)
		Update(record *types.Record) (*types.Record, error)

		DeleteByID(namespaceID, moduleID uint64, recordID ...uint64) error

		Organize(namespaceID, moduleID, recordID uint64, sortingField, sortingValue, sortingFilter, valueField, value string) error

		Iterator(f types.RecordFilter, fn eventbus.HandlerFn, action string) (err error)
	}

	Encoder interface {
		Record(*types.Record) error
	}

	Decoder interface {
		Header() []string
		EntryCount() (uint64, error)
		Records(fields map[string]string, Create decoder.RecordCreator) error
	}

	RecordImportSession struct {
		Decoder     Decoder              `json:"-"`
		CreatedAt   time.Time            `json:"createdAt"`
		UpdatedAt   time.Time            `json:"updatedAt"`
		OnError     string               `json:"onError"`
		SessionID   uint64               `json:"sessionID,string"`
		UserID      uint64               `json:"userID,string"`
		NamespaceID uint64               `json:"namespaceID,string"`
		ModuleID    uint64               `json:"moduleID,string"`
		Fields      map[string]string    `json:"fields"`
		Progress    RecordImportProgress `json:"progress"`
	}

	RecordImportProgress struct {
		StartedAt  *time.Time `json:"startedAt"`
		FinishedAt *time.Time `json:"finishedAt"`
		EntryCount uint64     `json:"entryCount"`
		Completed  uint64     `json:"completed"`
		Failed     uint64     `json:"failed"`
		FailReason string     `json:"failReason,omitempty"`
	}
)

func Record() RecordService {
	return (&record{
		logger:   DefaultLogger.Named("record"),
		ac:       DefaultAccessControl,
		eventbus: eventbus.Service(),
	}).With(context.Background())
}

func (svc record) With(ctx context.Context) RecordService {
	db := repository.DB(ctx)

	// Initialize validator and setup all checkers it needs
	validator := values.Validator()

	validator.UniqueChecker(func(v *types.RecordValue, f *types.ModuleField, m *types.Module) (uint64, error) {
		if v.Ref == 0 {
			return 0, nil
		}

		return repository.Record(ctx, db).RefValueLookup(m.ID, f.Name, v.Ref)
	})

	validator.RecordRefChecker(func(v *types.RecordValue, f *types.ModuleField, m *types.Module) (bool, error) {
		if v.Ref == 0 {
			return false, nil
		}

		r, err := repository.Record(ctx, db).FindByID(m.NamespaceID, v.Ref)
		return r != nil, err
	})

	validator.UserRefChecker(func(v *types.RecordValue, f *types.ModuleField, m *types.Module) (bool, error) {
		// @todo cross service check
		return true, nil
	})

	validator.FileRefChecker(func(v *types.RecordValue, f *types.ModuleField, m *types.Module) (bool, error) {
		if v.Ref == 0 {
			return false, nil
		}

		r, err := repository.Attachment(ctx, db).FindByID(m.NamespaceID, v.Ref)
		return r != nil, err
	})

	return &record{
		db:     db,
		ctx:    ctx,
		logger: svc.logger,

		ac:       svc.ac,
		eventbus: svc.eventbus,

		recordRepo: repository.Record(ctx, db),
		moduleRepo: repository.Module(ctx, db),
		nsRepo:     repository.Namespace(ctx, db),

		formatter: values.Formatter(),
		sanitizer: values.Sanitizer(),
		validator: validator,
	}
}

// log() returns zap's logger with requestID from current context and fields.
func (svc record) log(ctx context.Context, fields ...zapcore.Field) *zap.Logger {
	return logger.AddRequestID(ctx, svc.logger).With(fields...)
}

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

	if !svc.ac.CanReadModule(svc.ctx, m) {
		return nil, ErrNoReadPermissions.withStack()
	}

	if m.Fields, err = svc.moduleRepo.FindFields(m.ID); err != nil {
		return
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

func (svc record) Import(ses *RecordImportSession, ssvc ImportSessionService) error {
	if ses.Decoder == nil {
		return nil
	}

	if ses.Progress.StartedAt != nil {
		return errors.New("Unable to start import: Import session already active")
	}

	sa := time.Now()
	ses.Progress.StartedAt = &sa
	ssvc.SetRecordByID(svc.ctx, ses.SessionID, 0, 0, nil, &ses.Progress, nil)

	return svc.db.Transaction(func() (err error) {
		err = ses.Decoder.Records(ses.Fields, func(mod *types.Record) error {
			mod.NamespaceID = ses.NamespaceID
			mod.ModuleID = ses.ModuleID
			mod.OwnedBy = ses.UserID

			_, err := svc.Create(mod)
			if err != nil {
				ses.Progress.Failed++
				ses.Progress.FailReason = err.Error()

				if ses.OnError == IMPORT_ON_ERROR_FAIL {
					fa := time.Now()
					ses.Progress.FinishedAt = &fa
					ssvc.SetRecordByID(svc.ctx, ses.SessionID, 0, 0, nil, &ses.Progress, nil)
					return err
				}
			} else {
				ses.Progress.Completed++
			}
			return nil
		})

		fa := time.Now()
		ses.Progress.FinishedAt = &fa
		ssvc.SetRecordByID(svc.ctx, ses.SessionID, 0, 0, nil, &ses.Progress, nil)
		return
	})
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

func (svc record) Create(new *types.Record) (rec *types.Record, err error) {
	var invokerID = auth.GetIdentityFromContext(svc.ctx).Identity()

	return rec, svc.db.Transaction(func() (err error) {
		var (
			ns *types.Namespace
			m  *types.Module
		)

		ns, m, _, err = svc.loadCombo(new.NamespaceID, new.ModuleID, 0)
		if err != nil {
			return
		}

		if !svc.ac.CanCreateRecord(svc.ctx, m) {
			return ErrNoCreatePermissions.withStack()
		}

		if err = svc.generalValueSetValidation(m, new.Values); err != nil {
			return
		}

		var (
			rve *types.RecordValueErrorSet
		)

		// Handle input payload
		if rve = svc.procCreate(invokerID, m, new); !rve.IsValid() {
			return rve
		}

		if err = svc.eventbus.WaitFor(svc.ctx, event.RecordBeforeCreate(new, nil, m, ns, rve)); err != nil {
			return
		} else if !rve.IsValid() {
			return rve
		}

		// Assign defaults (only on missing values)
		new.Values = svc.setDefaultValues(m, new.Values)

		// Handle payload from automation scripts
		if rve = svc.procCreate(invokerID, m, new); !rve.IsValid() {
			return rve
		}

		if new, err = svc.recordRepo.Create(new); err != nil {
			return
		}

		if err = svc.recordRepo.UpdateValues(new.ID, new.Values); err != nil {
			return
		}

		// At this point we can return the value
		rec = new

		defer svc.eventbus.Dispatch(svc.ctx, event.RecordAfterCreateImmutable(new, nil, m, ns, nil))
		return
	})
}

// Runs value sanitization, sets values that should be used
// and validates the final result
//
// This logic is kept in a utility function - it's used in the beginning
// of the creation procedure and after results are back from the automation scripts
//
// Both these points introduce external data that need to be checked fully in the same manner
func (svc record) procCreate(invokerID uint64, m *types.Module, new *types.Record) *types.RecordValueErrorSet {
	// Mark all values as updated (new)
	new.Values.SetUpdatedFlag(true)

	// Before values are processed further and
	// sent to automation scripts (if any)
	// we need to make sure it does not get un-sanitized data
	new.Values = svc.sanitizer.Run(m, new.Values)

	// Reset values to new record
	// to make sure nobody slips in something we do not want
	new.CreatedBy = invokerID
	new.CreatedAt = *nowPtr()
	new.UpdatedAt = nil
	new.UpdatedBy = 0
	new.DeletedAt = nil
	new.DeletedBy = 0

	if new.OwnedBy == 0 {
		// If od owner is not set, make current user
		// the owner of the record
		new.OwnedBy = invokerID
	}

	// Run validation of the updated records
	return svc.validator.Run(m, new)
}

func (svc record) Update(upd *types.Record) (rec *types.Record, err error) {
	var invokerID = auth.GetIdentityFromContext(svc.ctx).Identity()

	return rec, svc.db.Transaction(func() (err error) {
		if upd.ID == 0 {
			return ErrInvalidID.withStack()
		}

		ns, m, old, err := svc.loadCombo(upd.NamespaceID, upd.ModuleID, upd.ID)
		if err != nil {
			return
		}

		if err = svc.generalValueSetValidation(m, upd.Values); err != nil {
			return
		}

		// Test if stale (update has an older version of data)
		if isStale(upd.UpdatedAt, old.UpdatedAt, old.CreatedAt) {
			return ErrStaleData.withStack()
		}

		if !svc.ac.CanUpdateRecord(svc.ctx, m) {
			return ErrNoUpdatePermissions.withStack()
		}

		// Preload old record values so we can send it together with event
		if err = svc.preloadValues(m, old); err != nil {
			return
		}

		var (
			rve *types.RecordValueErrorSet
		)

		// Handle input payload
		if rve = svc.procUpdate(invokerID, m, upd, old); !rve.IsValid() {
			return rve
		}

		// Before we pass values to record-before-update handling events
		// values needs do be cleaned up
		//
		// Value merge inside procUpdate sets delete flag we need
		// when changes are applied but we do not want deleted values
		// to be sent to handler
		upd.Values = upd.Values.GetClean()

		// Scripts can (besides simple error value) return complex record value error set
		// that is passed back to the UI or any other API consumer
		//
		// rve (record-validation-errorset) struct is passed so it can be
		// used & filled by automation scripts
		if err = svc.eventbus.WaitFor(svc.ctx, event.RecordBeforeUpdate(upd, old, m, ns, rve)); err != nil {
			return
		} else if !rve.IsValid() {
			return rve
		}

		// Handle payload from automation scripts
		if rve = svc.procUpdate(invokerID, m, upd, old); !rve.IsValid() {
			return rve
		}

		if upd, err = svc.recordRepo.Update(upd); err != nil {
			return
		}

		if err = svc.recordRepo.UpdateValues(upd.ID, upd.Values); err != nil {
			return
		}

		// Final value cleanup
		// These (clean) values are returned (and sent to after-update handler)
		upd.Values = upd.Values.GetClean()

		// At this point we can return the value
		rec = upd

		defer svc.eventbus.Dispatch(svc.ctx, event.RecordAfterUpdateImmutable(upd, old, m, ns, nil))
		return
	})
}

// Runs value sanitization, copies values that should updated
// and validates the final result
//
// This logic is kept in a utility function - it's used in the beginning
// of the update procedure and after results are back from the automation scripts
//
// Both these points introduce external data that need to be checked fully in the same manner
func (svc record) procUpdate(invokerID uint64, m *types.Module, upd *types.Record, old *types.Record) *types.RecordValueErrorSet {
	// Mark all values as updated (new)
	upd.Values.SetUpdatedFlag(true)

	// First sanitization
	//
	// Before values are merged with existing data and
	// sent to automation scripts (if any)
	// we need to make sure it does not get sanitized data
	upd.Values = svc.sanitizer.Run(m, upd.Values)

	// Copy values to updated record
	// to make sure nobody slips in something we do not want
	upd.CreatedAt = old.CreatedAt
	upd.CreatedBy = old.CreatedBy
	upd.UpdatedAt = nowPtr()
	upd.UpdatedBy = invokerID
	upd.DeletedAt = old.DeletedAt
	upd.DeletedBy = old.DeletedBy

	// Merge new (updated) values with old ones
	// This way we get list of updated, stale and deleted values
	// that we can selectively update in the repository
	upd.Values = old.Values.Merge(upd.Values)

	if upd.OwnedBy == 0 {
		if old.OwnedBy > 0 {
			// Owner not set/send in the payload
			//
			// Fallback to old owner (if set)
			upd.OwnedBy = old.OwnedBy
		} else {
			// If od owner is not set, make current user
			// the owner of the record
			upd.OwnedBy = invokerID
		}
	}

	// Run validation of the updated records
	return svc.validator.Run(m, upd)
}

func (svc record) recordInfoUpdate(r *types.Record) {
	now := time.Now()
	r.UpdatedAt = &now
	r.UpdatedBy = auth.GetIdentityFromContext(svc.ctx).Identity()
}

// DeleteByID removes one or more records (all from the same module and namespace)
//
// Before and after each record is deleted beforeDelete and afterDelete events are emitted
// If beforeRecord aborts the action it does so for that specific record only

func (svc record) DeleteByID(namespaceID, moduleID uint64, recordIDs ...uint64) error {
	if namespaceID == 0 {
		return ErrInvalidID.withStack()
	}

	if moduleID == 0 {
		return ErrInvalidID.withStack()
	}

	var (
		isBulkDelete = len(recordIDs) > 0
		now          = *nowPtr()

		ns *types.Namespace
		m  *types.Module

		err error
	)

	ns, m, _, err = svc.loadCombo(namespaceID, moduleID, 0)
	if err != nil {
		return err
	}

	if !svc.ac.CanDeleteRecord(svc.ctx, m) {
		return ErrNoDeletePermissions.withStack()
	}

	for _, recordID := range recordIDs {
		if recordID == 0 {
			return ErrInvalidID.withStack()
		}

		err := svc.db.Transaction(func() (err error) {
			var (
				del *types.Record
			)

			del, err = svc.FindByID(namespaceID, recordID)
			if err != nil {
				return err
			}

			// Preload old record values so we can send it together with event
			if err = svc.preloadValues(m, del); err != nil {
				return err
			}

			// Calling before-record-delete scripts
			if err = svc.eventbus.WaitFor(svc.ctx, event.RecordBeforeDelete(nil, del, m, ns, nil)); err != nil {
				if isBulkDelete {
					// Not considered fatal,
					// continue with next record
					return nil
				} else {
					return err
				}
			}

			del.DeletedAt = &now
			del.DeletedBy = auth.GetIdentityFromContext(svc.ctx).Identity()

			if err = svc.recordRepo.Delete(del); err != nil {
				return err
			}

			if err = svc.recordRepo.DeleteValues(del); err != nil {
				return err
			}

			defer svc.eventbus.Dispatch(svc.ctx, event.RecordAfterDeleteImmutable(nil, del, m, ns, nil))

			return err
		})

		if err != nil {
			return errors.Wrap(err, "failed to delete record")
		}
	}

	return nil
}

// Organize - Record organizer
//
// Reorders records & sets field value
func (svc record) Organize(namespaceID, moduleID, recordID uint64, posField, position, filter, grpField, group string) error {
	var (
		_, module, record, err = svc.loadCombo(namespaceID, moduleID, recordID)

		recordValues = types.RecordValueSet{}

		reorderingRecords bool

		log = svc.log(svc.ctx,
			zap.String("position-field", posField),
			zap.String("position", position),
			zap.String("filter", filter),
			zap.String("group-field", grpField),
			zap.String("group", group),
		)
	)

	log.Debug("record organizer")

	if err != nil {
		return err
	}

	if !svc.ac.CanUpdateRecord(svc.ctx, module) {
		return ErrNoUpdatePermissions.withStack()
	}

	if posField != "" {
		reorderingRecords = true

		if !regexp.MustCompile(`^[0-9]+$`).MatchString(position) {
			return errors.Errorf("expecting number for sorting position %q", posField)
		}

		// Check field existence and permissions
		// check if numeric -- we can not reorder on any other field type

		sf := module.Fields.FindByName(posField)
		if sf == nil {
			return errors.Errorf("no such field %q", posField)
		}

		if !sf.IsNumeric() {
			return errors.Errorf("can not reorder on non numeric field %q", posField)
		}

		if sf.Multi {
			return errors.Errorf("can not reorder on multi-value field %q", posField)
		}

		if !svc.ac.CanUpdateRecordValue(svc.ctx, sf) {
			return ErrNoUpdatePermissions.withStack()
		}

		// Set new position
		recordValues = recordValues.Set(&types.RecordValue{
			RecordID: recordID,
			Name:     posField,
			Value:    position,
		})
	}

	if grpField != "" {
		// Check field existence and permissions

		vf := module.Fields.FindByName(grpField)
		if vf == nil {
			return errors.Errorf("no such field %q", grpField)
		}

		if vf.Multi {
			return errors.Errorf("can not update multi-value field %q", posField)
		}

		if !svc.ac.CanUpdateRecordValue(svc.ctx, vf) {
			return ErrNoUpdatePermissions.withStack()
		}

		// Set new value
		recordValues = recordValues.Set(&types.RecordValue{
			RecordID: recordID,
			Name:     grpField,
			Value:    group,
		})
	}

	return svc.db.Transaction(func() (err error) {
		if len(recordValues) > 0 {
			svc.recordInfoUpdate(record)
			if _, err = svc.recordRepo.Update(record); err != nil {
				return
			}

			if err = svc.recordRepo.PartialUpdateValues(recordValues...); err != nil {
				return
			}

			log.Info("record moved")
		}

		if reorderingRecords {
			var (
				set              types.RecordSet
				recordOrderPlace uint64
			)

			// If we already have filter, wrap it in parenthesis
			if filter != "" {
				filter = fmt.Sprintf("(%s) AND ", filter)
			}

			if recordOrderPlace, err = strconv.ParseUint(position, 0, 64); err != nil {
				return
			}

			// Assemble record filter:
			// We are interested only in records that have value of a sorting field greater than
			// the place we're moving our record to.
			// and sort the set with sorting field
			set, _, err = svc.recordRepo.Find(module, types.RecordFilter{
				Filter: fmt.Sprintf("%s(%s >= %d)", filter, posField, recordOrderPlace),
				Sort:   posField,
			})

			log.Info("reordering other records", zap.Int("count", len(set)))

			if err != nil {
				return
			}

			// Update value on each record
			return set.Walk(func(r *types.Record) error {
				recordOrderPlace++

				// Update each and every set
				return svc.recordRepo.PartialUpdateValues(&types.RecordValue{
					RecordID: r.ID,
					Name:     posField,
					Value:    strconv.FormatUint(recordOrderPlace, 10),
				})
			})
		}

		return
	})
}

// Iterator loads and iterates through list of records
//
// For each record, RecordOnIteration is generated and passed to fn()
// to be then passed to automation script that invoked the iteration
//
// No other triggers (before/after update/delete/create) are fired when (if)
// records are changed
//
// action arg enables one of the following scenarios:
//   - clone:   make new record (unless aborted)
//   - update:  update records (unless aborted)
//   - delete:  delete records (unless aborted)
//   - default: only iterates over records, records are not changed, return value is ignored
//
//
// Iterator can be invoked only when defined in corredor script:
//
// return default {
//   iterator (each) {
//     return each({
//       resourceType: 'compose:record',
//       // action: 'update',
//       filter: {
//         namespace: '122709101053521922',
//         module: '122709116471783426',
//         query: 'Status = "foo"',
//         sort: 'Status DESC',
//         limit: 3,
//       },
//     })
//   },
//
//   // this is required in case of a deferred iterator
//   // security: { runAs: .... } }
//
//   // exec gets called for every record found by iterator
//   exec () { ... }
// }
func (svc record) Iterator(f types.RecordFilter, fn eventbus.HandlerFn, action string) (err error) {
	var (
		invokerID = auth.GetIdentityFromContext(svc.ctx).Identity()

		ns  *types.Namespace
		m   *types.Module
		set types.RecordSet
	)

	return svc.db.Transaction(func() (err error) {
		if ns, m, _, err = svc.loadCombo(f.NamespaceID, f.ModuleID, 0); err != nil {
			return
		}

		if !svc.ac.CanUpdateRecord(svc.ctx, m) {
			return ErrNoUpdatePermissions.withStack()
		}

		// @todo might be good to split set into smaller chunks
		set, f, err = svc.recordRepo.Find(m, f)
		if err != nil {
			return
		}

		if err = svc.preloadValues(m, set...); err != nil {
			return
		}

		for _, rec := range set {
			if err = fn(svc.ctx, event.RecordOnIteration(rec, nil, m, ns, nil)); err != nil {
				if err.Error() != "Aborted" {
					// When script was  softly aborted (return false),
					// proceed with iteration but do not clone, update or delete
					// current record!
					return
				}
			}

			switch action {
			case "clone":
				var cln *types.Record

				// Assign defaults (only on missing values)
				rec.Values = svc.setDefaultValues(m, rec.Values)

				// Handle payload from automation scripts
				if rve := svc.procCreate(invokerID, m, rec); !rve.IsValid() {
					return rve
				}

				if cln, err = svc.recordRepo.Create(rec); err != nil {
					return
				} else if err = svc.recordRepo.UpdateValues(cln.ID, cln.Values); err != nil {
					return
				}
			case "update":
				// Handle input payload
				if rve := svc.procUpdate(invokerID, m, rec, rec); !rve.IsValid() {
					return rve
				}

				if rec, err = svc.recordRepo.Update(rec); err != nil {
					return
				} else if err = svc.recordRepo.UpdateValues(rec.ID, rec.Values); err != nil {
					return
				}
			case "delete":
				if err = svc.recordRepo.Delete(rec); err != nil {
					return err
				} else if err = svc.recordRepo.DeleteValues(rec); err != nil {
					return err
				}
			}
		}

		return
	})
}

// loadCombo Loads everything we need for record manipulation
//
// Loads namespace, module, record and set of triggers.
func (svc record) loadCombo(namespaceID, moduleID, recordID uint64) (ns *types.Namespace, m *types.Module, r *types.Record, err error) {
	if namespaceID == 0 {
		err = ErrNamespaceRequired.withStack()
		return
	}
	if ns, err = svc.loadNamespace(namespaceID); err != nil {
		return
	}

	if recordID > 0 {
		if r, err = svc.recordRepo.FindByID(namespaceID, recordID); err != nil {
			return
		}

		if r.ModuleID != moduleID && moduleID > 0 {
			return nil, nil, nil, ErrInvalidModuleID.withStack()
		}
	}

	if moduleID > 0 {
		if m, err = svc.loadModule(ns.ID, moduleID); err != nil {
			return
		}
	}

	return
}

func (svc record) setDefaultValues(m *types.Module, vv types.RecordValueSet) (out types.RecordValueSet) {
	out = vv

	for _, f := range m.Fields {
		if f.DefaultValue == nil {
			continue
		}

		for i, dv := range f.DefaultValue {
			// Default values on field are (might be) without field name and place
			if !out.Has(f.Name, uint(i)) {
				out = append(out, &types.RecordValue{
					Name:  f.Name,
					Value: dv.Value,
					Place: uint(i),
				})
			}
		}
	}

	return
}

// Does basic field and format validation
//
// Received values must fit the data model: on unknown fields
// or multi/single value mismatch we return an error
//
func (svc record) generalValueSetValidation(m *types.Module, vv types.RecordValueSet) (err error) {
	var (
		numeric = regexp.MustCompile(`^[1-9](\d+)$`)
	)

	err = vv.Walk(func(v *types.RecordValue) error {
		var field = m.Fields.FindByName(v.Name)
		if field == nil {
			return errors.Errorf("no such field %q", v.Name)
		}

		if field.IsRef() {
			if v.Value == "" {
				return nil
			}

			if !numeric.MatchString(v.Value) {
				return errors.Errorf("invalid reference format")
			}
		}

		return nil
	})

	if err != nil {
		return
	}

	// Make sure there are no multi values in a non-multi value fields
	err = m.Fields.Walk(func(field *types.ModuleField) error {
		if !field.Multi && len(vv.FilterByName(field.Name)) > 1 {
			return errors.Errorf("more than one value for a single-value field %q", field.Name)
		}

		return nil
	})

	if err != nil {
		return
	}

	return
}

func (svc record) preloadValues(m *types.Module, rr ...*types.Record) error {
	if rvs, err := svc.recordRepo.LoadValues(svc.readableFields(m), types.RecordSet(rr).IDs()); err != nil {
		return err
	} else {
		return types.RecordSet(rr).Walk(func(r *types.Record) error {
			r.Values = svc.formatter.Run(m, rvs.FilterByRecordID(r.ID))
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
