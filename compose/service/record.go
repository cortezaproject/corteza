package service

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"time"

	"github.com/titpetric/factory"

	"github.com/cortezaproject/corteza-server/compose/decoder"
	"github.com/cortezaproject/corteza-server/compose/repository"
	"github.com/cortezaproject/corteza-server/compose/service/event"
	"github.com/cortezaproject/corteza-server/compose/service/values"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
)

const (
	IMPORT_ON_ERROR_SKIP         = "SKIP"
	IMPORT_ON_ERROR_FAIL         = "FAIL"
	IMPORT_ERROR_MAX_INDEX_COUNT = 500000
)

type (
	record struct {
		db  *factory.DB
		ctx context.Context

		actionlog actionlog.Recorder

		ac       recordAccessController
		eventbus eventDispatcher

		recordRepo repository.RecordRepository
		moduleRepo repository.ModuleRepository
		nsRepo     repository.NamespaceRepository

		formatter recordValuesFormatter
		sanitizer recordValuesSanitizer
		validator recordValuesValidator

		optEmitEvents bool
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
		Bulk(oo ...*types.RecordBulkOperation) (types.RecordSet, error)

		DeleteByID(namespaceID, moduleID uint64, recordID ...uint64) error

		Organize(namespaceID, moduleID, recordID uint64, sortingField, sortingValue, sortingFilter, valueField, value string) error

		Iterator(f types.RecordFilter, fn eventbus.HandlerFn, action string) (err error)

		EventEmitting(enable bool)
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
		FailLog    *FailLog   `json:"failLog,omitempty"`
	}

	FailLog struct {
		// Records holds an array of record indexes
		Records          RecordIndex `json:"records"`
		RecordsTruncated bool        `json:"recordsTruncated"`
		// Errors specifies a map of occurred errors & the number of
		Errors ErrorIndex `json:"errors"`
	}

	RecordIndex []int
	ErrorIndex  map[string]int
)

func Record() RecordService {
	return (&record{
		ac:            DefaultAccessControl,
		eventbus:      eventbus.Service(),
		optEmitEvents: true,
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
		db:  db,
		ctx: ctx,

		actionlog: DefaultActionlog,

		ac:       svc.ac,
		eventbus: svc.eventbus,

		recordRepo: repository.Record(ctx, db),
		moduleRepo: repository.Module(ctx, db),
		nsRepo:     repository.Namespace(ctx, db),

		formatter: values.Formatter(),
		sanitizer: values.Sanitizer(),
		validator: validator,

		optEmitEvents: svc.optEmitEvents,
	}
}

func (svc *record) EventEmitting(enable bool) {
	svc.optEmitEvents = enable
}

// lookup fn() orchestrates record lookup, namespace preload and check
func (svc record) lookup(namespaceID uint64, lookup func(*recordActionProps) (*types.Record, error)) (r *types.Record, err error) {
	var (
		ns     *types.Namespace
		m      *types.Module
		aProps = &recordActionProps{record: &types.Record{NamespaceID: namespaceID}}
	)

	err = func() error {
		if ns, err = svc.loadNamespace(namespaceID); err != nil {
			return err
		}

		aProps.setNamespace(ns)

		if r, err = lookup(aProps); err != nil {
			if repository.ErrRecordNotFound.Eq(err) {
				return RecordErrNotFound()
			}

			return err
		}

		aProps.setRecord(r)

		if m, err = svc.loadModule(namespaceID, r.ModuleID); err != nil {
			return err
		}

		aProps.setModule(m)

		if !svc.ac.CanReadRecord(svc.ctx, m) {
			return RecordErrNotAllowedToRead()
		}

		if err = svc.preloadValues(m, r); err != nil {
			return err
		}

		return nil
	}()

	return r, svc.recordAction(svc.ctx, aProps, RecordActionLookup, err)
}

func (svc record) FindByID(namespaceID, recordID uint64) (r *types.Record, err error) {
	return svc.lookup(namespaceID, func(props *recordActionProps) (*types.Record, error) {
		props.record.ID = recordID
		return svc.recordRepo.FindByID(namespaceID, recordID)
	})
}

func (svc record) loadModule(namespaceID, moduleID uint64) (m *types.Module, err error) {
	return m, func() error {
		if namespaceID == 0 {
			return RecordErrInvalidNamespaceID()
		}

		if moduleID == 0 {
			return RecordErrInvalidModuleID()
		}

		if m, err = svc.moduleRepo.FindByID(namespaceID, moduleID); err != nil {
			if repository.ErrModuleNotFound.Eq(err) {
				return RecordErrModuleNotFoundModule()
			}

			return err
		}

		if !svc.ac.CanReadModule(svc.ctx, m) {
			return RecordErrNotAllowedToReadModule()
		}

		if m.Fields, err = svc.moduleRepo.FindFields(m.ID); err != nil {
			return err
		}

		return nil
	}()
}

func (svc record) loadNamespace(namespaceID uint64) (ns *types.Namespace, err error) {
	return ns, func() error {
		if namespaceID == 0 {
			return RecordErrInvalidNamespaceID()
		}

		if ns, err = svc.nsRepo.FindByID(namespaceID); err != nil {
			if repository.ErrNamespaceNotFound.Eq(err) {
				return RecordErrNamespaceNotFound()
			}

			return err
		}

		if !svc.ac.CanReadNamespace(svc.ctx, ns) {
			return RecordErrNotAllowedToReadNamespace()
		}

		return err
	}()
}

// Report generates report for a given module using metrics, dimensions and filter
func (svc record) Report(namespaceID, moduleID uint64, metrics, dimensions, filter string) (out interface{}, err error) {
	var (
		ns     *types.Namespace
		m      *types.Module
		aProps = &recordActionProps{record: &types.Record{NamespaceID: namespaceID}}
	)

	err = func() error {
		if ns, err = svc.loadNamespace(namespaceID); err != nil {
			return err
		}

		aProps.setNamespace(ns)

		if m, err = svc.loadModule(namespaceID, moduleID); err != nil {
			return err
		}

		aProps.setModule(m)

		out, err = svc.recordRepo.Report(m, metrics, dimensions, filter)
		return err
	}()

	return out, svc.recordAction(svc.ctx, aProps, RecordActionReport, err)
}

func (svc record) Find(filter types.RecordFilter) (set types.RecordSet, f types.RecordFilter, err error) {
	var (
		m      *types.Module
		aProps = &recordActionProps{filter: &filter}
	)

	err = func() error {
		if m, err = svc.loadModule(filter.NamespaceID, filter.ModuleID); err != nil {
			return err
		}

		set, f, err = svc.recordRepo.Find(m, filter)
		if err != nil {
			return err
		}

		if err = svc.preloadValues(m, set...); err != nil {
			return err
		}

		return nil
	}()

	return set, f, svc.recordAction(svc.ctx, aProps, RecordActionSearch, err)
}

func (svc record) Import(ses *RecordImportSession, ssvc ImportSessionService) (err error) {
	var (
		aProps = &recordActionProps{}
	)

	if ses.Decoder == nil {
		return nil
	}

	err = func() (err error) {

		if ses.Progress.StartedAt != nil {
			return fmt.Errorf("Unable to start import: Import session already active")
		}

		sa := time.Now()
		ses.Progress.StartedAt = &sa
		ssvc.SetByID(svc.ctx, ses.SessionID, 0, 0, nil, &ses.Progress, nil)

		index := 0
		err = ses.Decoder.Records(ses.Fields, func(rec *types.Record) error {
			index++

			rec.NamespaceID = ses.NamespaceID
			rec.ModuleID = ses.ModuleID
			rec.OwnedBy = ses.UserID

			_, err := svc.Create(rec)
			if err != nil {
				recErr, isRecErr := err.(*recordError)

				ses.Progress.Failed++
				ses.Progress.FailReason = err.Error()

				if ses.Progress.FailLog == nil {
					ses.Progress.FailLog = &FailLog{
						Errors: make(ErrorIndex),
					}
				}

				if isRecErr {
					if evErr, ok := recErr.wrap.(*types.RecordValueErrorSet); ok {
						for _, ve := range evErr.Set {
							for k, v := range ve.Meta {
								ses.Progress.FailLog.Errors.Add(fmt.Sprintf("%s %s %v", ve.Kind, k, v))
							}
						}
					} else {
						ses.Progress.FailLog.Errors.Add(err.Error())
					}
				} else {
					ses.Progress.FailLog.Errors.Add(err.Error())
				}

				if len(ses.Progress.FailLog.Records) < IMPORT_ERROR_MAX_INDEX_COUNT {
					ses.Progress.FailLog.Records = append(ses.Progress.FailLog.Records, index)
				} else {
					ses.Progress.FailLog.RecordsTruncated = true
				}

				if ses.OnError == IMPORT_ON_ERROR_FAIL {
					fa := time.Now()
					ses.Progress.FinishedAt = &fa
					ssvc.SetByID(svc.ctx, ses.SessionID, 0, 0, nil, &ses.Progress, nil)
					return err
				}
			} else {
				ses.Progress.Completed++
			}
			return nil
		})

		fa := time.Now()
		ses.Progress.FinishedAt = &fa
		ssvc.SetByID(svc.ctx, ses.SessionID, 0, 0, nil, &ses.Progress, nil)
		return
	}()

	return svc.recordAction(svc.ctx, aProps, RecordActionImport, err)
}

// Export returns all records
//
// @todo better value handling
func (svc record) Export(filter types.RecordFilter, enc Encoder) (err error) {
	var (
		aProps = &recordActionProps{filter: &filter}
	)

	err = func() error {
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
	}()

	return svc.recordAction(svc.ctx, aProps, RecordActionExport, err)
}

// Bulk handles provided set of bulk record operations.
// It's able to create, update or delete records in a single transaction.
func (svc record) Bulk(oo ...*types.RecordBulkOperation) (rr types.RecordSet, err error) {
	var pr *types.Record

	err = func() error {
		// pre-verify all
		for _, p := range oo {
			switch p.Operation {
			case types.OperationTypeCreate, types.OperationTypeUpdate, types.OperationTypeDelete:
				// ok
			default:
				return RecordErrUnknownBulkOperation(&recordActionProps{bulkOperation: string(p.Operation)})
			}
		}

		var (
			// in case we get record value errors from create or update operations
			// we ll merge the errors into one slice and return it all together
			//
			// this is done under assumption that potential before-record-update/create automation
			// scripts are playing by the rules and do not do any changes before any potential
			// record value errors are returned
			//
			// @todo all records/values could and should be pre-validated
			//       before we start storing any changes
			rves = &types.RecordValueErrorSet{}

			action func(props ...*recordActionProps) *recordAction
			r      *types.Record

			aProp = &recordActionProps{}
		)

		for _, p := range oo {
			r = p.Record

			aProp.setChanged(r)

			// Handle any pre processing, such as defining parent recordID.
			if p.LinkBy != "" {
				// As is, we can use the first record as the master record.
				// This is valid, since we do not allow this, if the master record is not defined
				rv := &types.RecordValue{
					Name: p.LinkBy,
				}
				if pr != nil {
					rv.Value = strconv.FormatUint(rr[0].ID, 10)
					rv.Ref = rr[0].ID
				}
				r.Values = r.Values.Set(rv)
			}

			switch p.Operation {
			case types.OperationTypeCreate:
				action = RecordActionCreate
				r, err = svc.create(r)

			case types.OperationTypeUpdate:
				action = RecordActionUpdate
				r, err = svc.update(r)

			case types.OperationTypeDelete:
				action = RecordActionDelete
				r, err = svc.delete(r.NamespaceID, r.ModuleID, r.ID)
			}

			if rve := types.IsRecordValueErrorSet(err); rve != nil {
				// Attach additional meta to each value error for FE identification
				for _, re := range rve.Set {
					re.Meta["id"] = p.ID

					rves.Push(re)
				}

				// log record value error for this record
				_ = svc.recordAction(svc.ctx, aProp, action, err)

				// do not return errors just yet, values on other records from the payload (if any)
				// might have errors too
				continue
			}

			if err != nil {
				return svc.recordAction(svc.ctx, aProp, action, err)
			}

			rr = append(rr, r)
			if pr == nil {
				pr = r
			}
		}

		if !rves.IsValid() {
			// Any errors gathered?
			return RecordErrValueInput().Wrap(rves)
		}

		return nil
	}()

	if len(oo) == 1 {
		// was not really a bulk operation and we already recorded the action
		// inside transaction loop
		return rr, err
	} else {
		// when doing bulk op (updating and/or creating more than one record at once),
		// we already log action for each operation
		//
		// to log the fact that the bulk op was done, we do one additional recording
		// without any props
		return rr, svc.recordAction(svc.ctx, &recordActionProps{}, RecordActionBulk, err)
	}
}

// Raw create function that is responsible for value validation, event dispatching
// and creation.
func (svc record) create(new *types.Record) (rec *types.Record, err error) {
	var (
		aProps    = &recordActionProps{changed: new}
		invokerID = auth.GetIdentityFromContext(svc.ctx).Identity()

		ns *types.Namespace
		m  *types.Module
	)

	ns, m, _, err = svc.loadCombo(new.NamespaceID, new.ModuleID, 0)
	if err != nil {
		return
	}

	aProps.setNamespace(ns)
	aProps.setModule(m)

	if !svc.ac.CanCreateRecord(svc.ctx, m) {
		return nil, RecordErrNotAllowedToCreate()
	}

	if err = svc.generalValueSetValidation(m, new.Values); err != nil {
		return
	}

	var (
		rve *types.RecordValueErrorSet
	)

	if svc.optEmitEvents {
		// Handle input payload
		if rve = svc.procCreate(invokerID, m, new); !rve.IsValid() {
			return nil, RecordErrValueInput().Wrap(rve)
		}

		new.Values = svc.formatter.Run(m, new.Values)
		if err = svc.eventbus.WaitFor(svc.ctx, event.RecordBeforeCreate(new, nil, m, ns, rve)); err != nil {
			return
		} else if !rve.IsValid() {
			return nil, RecordErrValueInput().Wrap(rve)
		}
	}

	// Assign defaults (only on missing values)
	new.Values = svc.setDefaultValues(m, new.Values)

	// Handle payload from automation scripts
	if rve = svc.procCreate(invokerID, m, new); !rve.IsValid() {
		return nil, RecordErrValueInput().Wrap(rve)
	}

	err = svc.db.Transaction(func() error {
		if new, err = svc.recordRepo.Create(new); err != nil {
			return err
		}

		return svc.recordRepo.UpdateValues(new.ID, new.Values)
	})

	if err != nil {
		return nil, err
	}

	// At this point we can return the value
	rec = new

	if svc.optEmitEvents {
		new.Values = svc.formatter.Run(m, new.Values)
		_ = svc.eventbus.WaitFor(svc.ctx, event.RecordAfterCreateImmutable(new, nil, m, ns, nil))
	}

	return
}

// Raw update function that is responsible for value validation, event dispatching
// and update.
func (svc record) update(upd *types.Record) (rec *types.Record, err error) {
	var (
		aProps    = &recordActionProps{changed: upd}
		invokerID = auth.GetIdentityFromContext(svc.ctx).Identity()

		ns  *types.Namespace
		m   *types.Module
		old *types.Record
	)

	if upd.ID == 0 {
		return nil, RecordErrInvalidID()
	}

	ns, m, old, err = svc.loadCombo(upd.NamespaceID, upd.ModuleID, upd.ID)
	if err != nil {
		return
	}

	aProps.setNamespace(ns)
	aProps.setModule(m)
	aProps.setRecord(old)

	if !svc.ac.CanUpdateRecord(svc.ctx, m) {
		return nil, RecordErrNotAllowedToUpdate()
	}

	// Test if stale (update has an older version of data)
	if isStale(upd.UpdatedAt, old.UpdatedAt, old.CreatedAt) {
		return nil, RecordErrStaleData()
	}

	if err = svc.generalValueSetValidation(m, upd.Values); err != nil {
		return
	}

	// Preload old record values so we can send it together with event
	if err = svc.preloadValues(m, old); err != nil {
		return
	}

	var (
		rve *types.RecordValueErrorSet
	)

	if svc.optEmitEvents {
		// Handle input payload
		if rve = svc.procUpdate(invokerID, m, upd, old); !rve.IsValid() {
			return nil, RecordErrValueInput().Wrap(rve)
		}

		// Before we pass values to record-before-update handling events
		// values needs do be cleaned up
		//
		// Value merge inside procUpdate sets delete flag we need
		// when changes are applied but we do not want deleted values
		// to be sent to handler
		upd.Values = upd.Values.GetClean()

		// Before we pass values to automation scripts, they should be formatted
		upd.Values = svc.formatter.Run(m, upd.Values)

		// Scripts can (besides simple error value) return complex record value error set
		// that is passed back to the UI or any other API consumer
		//
		// rve (record-validation-errorset) struct is passed so it can be
		// used & filled by automation scripts
		if err = svc.eventbus.WaitFor(svc.ctx, event.RecordBeforeUpdate(upd, old, m, ns, rve)); err != nil {
			return
		} else if !rve.IsValid() {
			return nil, RecordErrValueInput().Wrap(rve)
		}
	}

	// Handle payload from automation scripts
	if rve = svc.procUpdate(invokerID, m, upd, old); !rve.IsValid() {
		return nil, RecordErrValueInput().Wrap(rve)
	}

	err = svc.db.Transaction(func() error {
		if upd, err = svc.recordRepo.Update(upd); err != nil {
			return nil
		}

		return svc.recordRepo.UpdateValues(upd.ID, upd.Values)

	})

	if err != nil {
		return nil, err
	}

	// Final value cleanup
	// These (clean) values are returned (and sent to after-update handler)
	upd.Values = upd.Values.GetClean()

	// At this point we can return the value
	rec = upd

	if svc.optEmitEvents {
		// Before we pass values to automation scripts, they should be formatted
		upd.Values = svc.formatter.Run(m, upd.Values)
		_ = svc.eventbus.WaitFor(svc.ctx, event.RecordAfterUpdateImmutable(upd, old, m, ns, nil))
	}
	return
}

func (svc record) Create(new *types.Record) (rec *types.Record, err error) {
	var (
		aProps = &recordActionProps{changed: new}
	)

	err = func() error {
		rec, err = svc.create(new)
		aProps.setRecord(rec)
		return err
	}()

	return rec, svc.recordAction(svc.ctx, aProps, RecordActionCreate, err)
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

	rve := &types.RecordValueErrorSet{}
	_ = new.Values.Walk(func(v *types.RecordValue) error {
		if v.IsUpdated() && !svc.ac.CanUpdateRecordValue(svc.ctx, m.Fields.FindByName(v.Name)) {
			rve.Push(types.RecordValueError{Kind: "updateDenied", Meta: map[string]interface{}{"field": v.Name, "value": v.Value}})
		}

		return nil
	})

	if !rve.IsValid() {
		return rve
	}

	// Run validation of the updated records
	return svc.validator.Run(m, new)
}

func (svc record) Update(upd *types.Record) (rec *types.Record, err error) {
	var (
		aProps = &recordActionProps{changed: upd}
	)

	err = func() error {
		rec, err = svc.update(upd)
		aProps.setRecord(rec)
		return err
	}()

	return rec, svc.recordAction(svc.ctx, aProps, RecordActionUpdate, err)
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

	// Value merge process does not know anything about permissions so
	// in case when new values are missing but do exist in the old set and their update/read is denied
	// we need to copy them to ensure value merge process them correctly
	for _, f := range m.Fields {
		if len(upd.Values.FilterByName(f.Name)) == 0 && !svc.ac.CanUpdateRecordValue(svc.ctx, m.Fields.FindByName(f.Name)) {
			// copy all fields from old to new
			upd.Values = append(upd.Values, old.Values.FilterByName(f.Name).GetClean()...)
		}
	}

	// Merge new (updated) values with old ones
	// This way we get list of updated, stale and deleted values
	// that we can selectively update in the repository
	upd.Values = old.Values.Merge(upd.Values)

	rve := &types.RecordValueErrorSet{}
	_ = upd.Values.Walk(func(v *types.RecordValue) error {
		if v.IsUpdated() && !svc.ac.CanUpdateRecordValue(svc.ctx, m.Fields.FindByName(v.Name)) {
			rve.Push(types.RecordValueError{Kind: "updateDenied", Meta: map[string]interface{}{"field": v.Name, "value": v.Value}})
		}

		return nil
	})

	if !rve.IsValid() {
		return rve
	}

	// Run validation of the updated records
	return svc.validator.Run(m, upd)
}

func (svc record) recordInfoUpdate(r *types.Record) {
	now := time.Now()
	r.UpdatedAt = &now
	r.UpdatedBy = auth.GetIdentityFromContext(svc.ctx).Identity()
}

func (svc record) delete(namespaceID, moduleID, recordID uint64) (del *types.Record, err error) {
	var (
		ns *types.Namespace
		m  *types.Module

		invokerID = auth.GetIdentityFromContext(svc.ctx).Identity()
	)

	if namespaceID == 0 {
		return nil, RecordErrInvalidNamespaceID()
	}
	if moduleID == 0 {
		return nil, RecordErrInvalidModuleID()
	}
	if recordID == 0 {
		return nil, RecordErrInvalidID()
	}

	ns, m, del, err = svc.loadCombo(namespaceID, moduleID, recordID)
	if err != nil {
		return nil, err
	}

	if !svc.ac.CanDeleteRecord(svc.ctx, m) {
		return nil, RecordErrNotAllowedToDelete()
	}

	if svc.optEmitEvents {
		// Preload old record values so we can send it together with event
		if err = svc.preloadValues(m, del); err != nil {
			return nil, err
		}

		// Calling before-record-delete scripts
		if err = svc.eventbus.WaitFor(svc.ctx, event.RecordBeforeDelete(nil, del, m, ns, nil)); err != nil {
			return nil, err
		}
	}

	del.DeletedAt = nowPtr()
	del.DeletedBy = invokerID

	err = svc.db.Transaction(func() error {
		if err = svc.recordRepo.Delete(del); err != nil {
			return err
		}

		return svc.recordRepo.DeleteValues(del)
	})

	if err != nil {
		return nil, err
	}

	if svc.optEmitEvents {
		_ = svc.eventbus.WaitFor(svc.ctx, event.RecordAfterDeleteImmutable(nil, del, m, ns, nil))
	}

	return del, nil
}

// DeleteByID removes one or more records (all from the same module and namespace)
//
// Before and after each record is deleted beforeDelete and afterDelete events are emitted
// If beforeRecord aborts the action it does so for that specific record only
func (svc record) DeleteByID(namespaceID, moduleID uint64, recordIDs ...uint64) (err error) {
	var (
		aProps = &recordActionProps{
			namespace: &types.Namespace{ID: namespaceID},
			module:    &types.Module{ID: moduleID},
		}

		isBulkDelete = len(recordIDs) > 1

		ns *types.Namespace
		m  *types.Module
		r  *types.Record
	)

	err = func() error {
		if namespaceID == 0 {
			return RecordErrInvalidNamespaceID()
		}
		if moduleID == 0 {
			return RecordErrInvalidModuleID()
		}

		ns, m, _, err = svc.loadCombo(namespaceID, moduleID, 0)
		if err != nil {
			return err
		}

		aProps.setNamespace(ns)
		aProps.setModule(m)

		if !svc.ac.CanDeleteRecord(svc.ctx, m) {
			return RecordErrNotAllowedToDelete()
		}

		return nil
	}()

	if err != nil {
		return svc.recordAction(svc.ctx, aProps, RecordActionDelete, err)
	}

	for _, recordID := range recordIDs {
		err := func() (err error) {
			r, err = svc.delete(namespaceID, moduleID, recordID)
			aProps.setRecord(r)

			// Record each record deletion action
			return svc.recordAction(svc.ctx, aProps, RecordActionDelete, err)
		}()

		// We'll not break for failed delete,
		// if we are deleting records in bulk.
		if err != nil && !isBulkDelete {
			return err
		}

	}

	// all errors (if any) were recorded
	// and in case of error for a non-bulk record deletion
	// error is already returned
	return nil
}

func (svc record) Organize(namespaceID, moduleID, recordID uint64, posField, position, filter, grpField, group string) (err error) {
	var (
		ns *types.Namespace
		m  *types.Module
		r  *types.Record

		recordValues = types.RecordValueSet{}

		aProps = &recordActionProps{record: &types.Record{NamespaceID: namespaceID, ModuleID: moduleID, ID: recordID}}

		reorderingRecords bool
	)

	err = func() error {
		ns, m, r, err = svc.loadCombo(namespaceID, moduleID, recordID)
		if err != nil {
			return err
		}

		aProps.setNamespace(ns)
		aProps.setModule(m)
		aProps.setRecord(r)

		if !svc.ac.CanUpdateRecord(svc.ctx, m) {
			return RecordErrNotAllowedToUpdate()
		}

		if posField != "" {
			reorderingRecords = true

			if !regexp.MustCompile(`^[0-9]+$`).MatchString(position) {
				return fmt.Errorf("expecting number for sorting position %q", posField)
			}

			// Check field existence and permissions
			// check if numeric -- we can not reorder on any other field type

			sf := m.Fields.FindByName(posField)
			if sf == nil {
				return fmt.Errorf("no such field %q", posField)
			}

			if !sf.IsNumeric() {
				return fmt.Errorf("can not reorder on non numeric field %q", posField)
			}

			if sf.Multi {
				return fmt.Errorf("can not reorder on multi-value field %q", posField)
			}

			if !svc.ac.CanUpdateRecordValue(svc.ctx, sf) {
				return RecordErrNotAllowedToUpdate()
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

			vf := m.Fields.FindByName(grpField)
			if vf == nil {
				return fmt.Errorf("no such field %q", grpField)
			}

			if vf.Multi {
				return fmt.Errorf("can not update multi-value field %q", posField)
			}

			if !svc.ac.CanUpdateRecordValue(svc.ctx, vf) {
				return RecordErrNotAllowedToUpdate()
			}

			// Set new value
			recordValues = recordValues.Set(&types.RecordValue{
				RecordID: recordID,
				Name:     grpField,
				Value:    group,
			})
		}

		return svc.db.Transaction(func() error {
			if len(recordValues) > 0 {
				svc.recordInfoUpdate(r)
				if _, err = svc.recordRepo.Update(r); err != nil {
					return err
				}

				if err = svc.recordRepo.PartialUpdateValues(recordValues...); err != nil {
					return err
				}
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
					return err
				}

				// Assemble record filter:
				// We are interested only in records that have value of a sorting field greater than
				// the place we're moving our record to.
				// and sort the set with sorting field
				set, _, err = svc.recordRepo.Find(m, types.RecordFilter{
					Query: fmt.Sprintf("%s(%s >= %d)", filter, posField, recordOrderPlace),
					Sort:  posField,
				})

				if err != nil {
					return err
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

			return nil
		})
	}()

	return svc.recordAction(svc.ctx, aProps, RecordActionOrganize, err)

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

		aProps = &recordActionProps{}
	)

	err = func() error {
		if ns, m, _, err = svc.loadCombo(f.NamespaceID, f.ModuleID, 0); err != nil {
			return err
		}

		if !svc.ac.CanReadRecord(svc.ctx, m) {
			return RecordErrNotAllowedToRead()
		}

		switch action {
		case "clone":
			if !svc.ac.CanCreateRecord(svc.ctx, m) {
				return RecordErrNotAllowedToCreate()
			}

		case "update":
			if !svc.ac.CanUpdateRecord(svc.ctx, m) {
				return RecordErrNotAllowedToUpdate()
			}

		case "delete":
			if !svc.ac.CanDeleteRecord(svc.ctx, m) {
				return RecordErrNotAllowedToDelete()
			}
		}

		// @todo might be good to split set into smaller chunks
		set, f, err = svc.recordRepo.Find(m, f)
		if err != nil {
			return err
		}

		if err = svc.preloadValues(m, set...); err != nil {
			return err
		}

		for _, rec := range set {
			recordableAction := RecordActionIteratorIteration

			err = func() error {
				if err = fn(svc.ctx, event.RecordOnIteration(rec, nil, m, ns, nil)); err != nil {
					if err.Error() != "Aborted" {
						// When script was softly aborted (return false),
						// proceed with iteration but do not clone, update or delete
						// current record!
						return nil
					}
				}

				switch action {
				case "clone":
					recordableAction = RecordActionIteratorClone

					var cln *types.Record

					// Assign defaults (only on missing values)
					rec.Values = svc.setDefaultValues(m, rec.Values)

					// Handle payload from automation scripts
					if rve := svc.procCreate(invokerID, m, rec); !rve.IsValid() {
						return RecordErrValueInput().Wrap(rve)
					}

					return svc.db.Transaction(func() error {
						if cln, err = svc.recordRepo.Create(rec); err != nil {
							return err
						} else if err = svc.recordRepo.UpdateValues(cln.ID, cln.Values); err != nil {
							return err
						}

						return nil
					})
				case "update":
					recordableAction = RecordActionIteratorUpdate

					// Handle input payload
					if rve := svc.procUpdate(invokerID, m, rec, rec); !rve.IsValid() {
						return RecordErrValueInput().Wrap(rve)
					}

					return svc.db.Transaction(func() error {
						if rec, err = svc.recordRepo.Update(rec); err != nil {
							return err
						} else if err = svc.recordRepo.UpdateValues(rec.ID, rec.Values); err != nil {
							return err
						}

						return nil
					})
				case "delete":
					recordableAction = RecordActionIteratorDelete

					return svc.db.Transaction(func() error {
						if err = svc.recordRepo.Delete(rec); err != nil {
							return err
						} else if err = svc.recordRepo.DeleteValues(rec); err != nil {
							return err
						}

						return nil
					})
				}

				return nil
			}()

			// record iteration action and
			// break the loop in case of an error
			_ = svc.recordAction(svc.ctx, aProps, recordableAction, err)
			if err != nil {
				return err
			}
		}

		return nil
	}()

	return svc.recordAction(svc.ctx, aProps, RecordActionIteratorInvoked, err)

}

// loadCombo Loads everything we need for record manipulation
//
// Loads namespace, module, record and set of triggers.
func (svc record) loadCombo(namespaceID, moduleID, recordID uint64) (ns *types.Namespace, m *types.Module, r *types.Record, err error) {
	if namespaceID == 0 {
		return nil, nil, nil, RecordErrInvalidNamespaceID()
	}
	if ns, err = svc.loadNamespace(namespaceID); err != nil {
		return
	}

	if recordID > 0 {
		if r, err = svc.recordRepo.FindByID(namespaceID, recordID); err != nil {
			return
		}

		if r.ModuleID != moduleID && moduleID > 0 {
			return nil, nil, nil, RecordErrInvalidModuleID()
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
// Record value errors is intentionally NOT used here; if input fails here
// we can assume that form builder (or whatever it was that assembled the record values)
// was misconfigured and will most likely failed to properly parse the
// record value errors payload too
func (svc record) generalValueSetValidation(m *types.Module, vv types.RecordValueSet) (err error) {
	var (
		aProps  = &recordActionProps{}
		numeric = regexp.MustCompile(`^[1-9](\d+)$`)
	)

	err = vv.Walk(func(v *types.RecordValue) error {
		var field = m.Fields.FindByName(v.Name)
		if field == nil {
			return RecordErrFieldNotFound(aProps.setField(v.Name))
		}

		if field.IsRef() {
			if v.Value == "" {
				return nil
			}

			if !numeric.MatchString(v.Value) {
				return RecordErrInvalidReferenceFormat(aProps.setField(v.Name).setValue(v.Value))
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
			return RecordErrInvalidValueStructure(aProps.setField(field.Name))
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

func (ei ErrorIndex) Add(err string) {
	if _, has := ei[err]; has {
		ei[err]++
	} else {
		ei[err] = 1
	}
}

func (ri RecordIndex) MarshalJSON() ([]byte, error) {
	sort.Ints(ri)

	rr := make([][]int, 0, len(ri))
	start := -1
	crt := -1

	for i := 0; i < len(ri); i++ {
		if start == -1 {
			start = ri[i]
			crt = ri[i]
			continue
		}

		// If the index increases for more then 1, the set is complete
		if ri[i]-crt > 1 {
			rr = append(rr, []int{start, crt})
			start = ri[i]
		}

		crt = ri[i]
	}

	rr = append(rr, []int{start, crt})
	return json.Marshal(rr)
}
