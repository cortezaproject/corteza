package service

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/cortezaproject/corteza-server/compose/service/event"
	"github.com/cortezaproject/corteza-server/compose/service/values"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/corredor"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/cortezaproject/corteza-server/pkg/label"
	"github.com/cortezaproject/corteza-server/store"
)

const (
	IMPORT_ON_ERROR_SKIP = "SKIP"
	IMPORT_ON_ERROR_FAIL = "FAIL"
)

type (
	record struct {
		actionlog actionlog.Recorder

		ac       recordAccessController
		eventbus eventDispatcher

		store store.Storer

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
		Run(context.Context, store.Storer, *types.Module, *types.Record) *types.RecordValueErrorSet
		UniqueChecker(fn values.UniqueChecker)
		RecordRefChecker(fn values.ReferenceChecker)
		UserRefChecker(fn values.ReferenceChecker)
	}

	recordValueAccessController interface {
		CanReadRecordValue(context.Context, *types.ModuleField) bool
		CanUpdateRecordValue(context.Context, *types.ModuleField) bool
	}

	recordAccessController interface {
		CanCreateRecord(context.Context, *types.Module) bool
		CanReadNamespace(context.Context, *types.Namespace) bool
		CanReadModule(context.Context, *types.Module) bool
		CanReadRecord(context.Context, *types.Module) bool
		CanUpdateRecord(context.Context, *types.Module) bool
		CanDeleteRecord(context.Context, *types.Module) bool

		recordValueAccessController
	}

	RecordService interface {
		FindByID(ctx context.Context, namespaceID, moduleID, recordID uint64) (*types.Record, error)

		Report(ctx context.Context, namespaceID, moduleID uint64, metrics, dimensions, filter string) (interface{}, error)
		Find(ctx context.Context, filter types.RecordFilter) (set types.RecordSet, f types.RecordFilter, err error)
		RecordExport(context.Context, types.RecordFilter) error
		RecordImport(context.Context, error) error

		Create(ctx context.Context, record *types.Record) (*types.Record, error)
		Update(ctx context.Context, record *types.Record) (*types.Record, error)
		Bulk(ctx context.Context, oo ...*types.RecordBulkOperation) (types.RecordSet, error)

		Validate(ctx context.Context, rec *types.Record) error

		DeleteByID(ctx context.Context, namespaceID, moduleID uint64, recordID ...uint64) error

		Organize(ctx context.Context, namespaceID, moduleID, recordID uint64, sortingField, sortingValue, sortingFilter, valueField, value string) error

		Iterator(ctx context.Context, f types.RecordFilter, fn eventbus.HandlerFn, action string) (err error)

		TriggerScript(ctx context.Context, namespaceID, moduleID, recordID uint64, rvs types.RecordValueSet, script string) (*types.Module, *types.Record, error)

		EventEmitting(enable bool)
	}

	recordImportSession struct {
		Name        string `json:"-"`
		SessionID   uint64 `json:"sessionID,string"`
		UserID      uint64 `json:"userID,string"`
		NamespaceID uint64 `json:"namespaceID,string"`
		ModuleID    uint64 `json:"moduleID,string"`

		OnError  string                `json:"onError"`
		Fields   map[string]string     `json:"fields"`
		Progress *RecordImportProgress `json:"progress"`

		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`

		Resources []resource.Interface `json:"-"`
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
	// Initialize validator and setup all checkers it needs
	validator := values.Validator()

	validator.UniqueChecker(func(ctx context.Context, s store.Storer, v *types.RecordValue, f *types.ModuleField, m *types.Module) (uint64, error) {
		if v.Ref == 0 {
			return 0, nil
		}

		return store.ComposeRecordValueRefLookup(ctx, s, m, f.Name, v.Ref)
	})

	validator.RecordRefChecker(func(ctx context.Context, s store.Storer, v *types.RecordValue, f *types.ModuleField, m *types.Module) (bool, error) {
		if v.Ref == 0 {
			return false, nil
		}

		r, err := store.LookupComposeRecordByID(ctx, s, m, v.Ref)
		return r != nil, err
	})

	validator.UserRefChecker(func(ctx context.Context, s store.Storer, v *types.RecordValue, f *types.ModuleField, m *types.Module) (bool, error) {
		r, err := store.LookupUserByID(ctx, s, v.Ref)
		return r != nil, err
	})

	validator.FileRefChecker(func(ctx context.Context, s store.Storer, v *types.RecordValue, f *types.ModuleField, m *types.Module) (bool, error) {
		if v.Ref == 0 {
			return false, nil
		}

		r, err := store.LookupComposeAttachmentByID(ctx, s, v.Ref)
		return r != nil, err
	})

	return &record{
		actionlog:     DefaultActionlog,
		ac:            DefaultAccessControl,
		eventbus:      eventbus.Service(),
		optEmitEvents: true,
		store:         DefaultStore,

		formatter: values.Formatter(),
		sanitizer: values.Sanitizer(),
		validator: validator,
	}
}

func (svc *record) EventEmitting(enable bool) {
	svc.optEmitEvents = enable
}

// lookup fn() orchestrates record lookup, namespace preload and check
func (svc record) lookup(ctx context.Context, namespaceID, moduleID uint64, lookup func(*types.Module, *recordActionProps) (*types.Record, error)) (r *types.Record, err error) {
	var (
		ns     *types.Namespace
		m      *types.Module
		aProps = &recordActionProps{record: &types.Record{NamespaceID: namespaceID}}
	)

	err = func() error {
		if ns, m, err = loadModuleWithNamespace(ctx, svc.store, namespaceID, moduleID); err != nil {
			return err
		}

		aProps.setNamespace(ns)
		aProps.setModule(m)

		if r, err = lookup(m, aProps); errors.IsNotFound(err) {
			return RecordErrNotFound()
		} else if err != nil {
			return err
		}

		aProps.setRecord(r)

		if !svc.ac.CanReadRecord(ctx, m) {
			return RecordErrNotAllowedToRead()
		}

		trimUnreadableRecordFields(ctx, svc.ac, m, r)

		if err = label.Load(ctx, svc.store, r); err != nil {
			return err
		}

		r.SetModule(m)

		return nil
	}()

	return r, svc.recordAction(ctx, aProps, RecordActionLookup, err)
}

func (svc record) FindByID(ctx context.Context, namespaceID, moduleID, recordID uint64) (r *types.Record, err error) {
	return svc.lookup(ctx, namespaceID, moduleID, func(m *types.Module, props *recordActionProps) (*types.Record, error) {
		props.record.ID = recordID
		return store.LookupComposeRecordByID(ctx, svc.store, m, recordID)
	})
}

// Report generates report for a given module using metrics, dimensions and filter
func (svc record) Report(ctx context.Context, namespaceID, moduleID uint64, metrics, dimensions, filter string) (out interface{}, err error) {
	var (
		ns     *types.Namespace
		m      *types.Module
		aProps = &recordActionProps{record: &types.Record{NamespaceID: namespaceID}}
	)

	err = func() error {
		if ns, m, err = loadModuleWithNamespace(ctx, svc.store, namespaceID, moduleID); err != nil {
			return err
		}

		aProps.setNamespace(ns)
		aProps.setModule(m)

		out, err = store.ComposeRecordReport(ctx, svc.store, m, metrics, dimensions, filter)
		return err
	}()

	return out, svc.recordAction(ctx, aProps, RecordActionReport, err)
}

func (svc record) Find(ctx context.Context, filter types.RecordFilter) (set types.RecordSet, f types.RecordFilter, err error) {
	var (
		m      *types.Module
		aProps = &recordActionProps{filter: &filter}
	)

	err = func() error {
		if m, err = loadModule(ctx, svc.store, filter.ModuleID); err != nil {
			return err
		}

		filter.Check = func(res *types.Record) (bool, error) {
			if !svc.ac.CanReadRecord(ctx, m) {
				return false, nil
			}

			return true, nil
		}

		if len(filter.Labels) > 0 {
			filter.LabeledIDs, err = label.Search(
				ctx,
				svc.store,
				types.Record{}.LabelResourceKind(),
				filter.Labels,
			)

			if err != nil {
				return err
			}

			// labels specified but no labeled resources found
			if len(filter.LabeledIDs) == 0 {
				return nil
			}
		}

		set, f, err = store.SearchComposeRecords(ctx, svc.store, m, filter)
		if err != nil {
			return err
		}

		if err = label.Load(ctx, svc.store, toLabeledRecords(set)...); err != nil {
			return err
		}

		_ = set.Walk(func(r *types.Record) error {
			r.SetModule(m)
			return nil
		})

		trimUnreadableRecordFields(ctx, svc.ac, m, set...)

		return nil
	}()

	return set, f, svc.recordAction(ctx, aProps, RecordActionSearch, err)
}

func (svc record) RecordImport(ctx context.Context, err error) error {
	return svc.recordAction(ctx, &recordActionProps{}, RecordActionImport, err)
}

// RecordExport records that the export has occurred
func (svc record) RecordExport(ctx context.Context, f types.RecordFilter) (err error) {
	return svc.recordAction(ctx, &recordActionProps{filter: &f}, RecordActionExport, err)
}

// Bulk handles provided set of bulk record operations.
// It's able to create, update or delete records in a single transaction.
func (svc record) Bulk(ctx context.Context, oo ...*types.RecordBulkOperation) (rr types.RecordSet, err error) {
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
				r, err = svc.create(ctx, r)

			case types.OperationTypeUpdate:
				action = RecordActionUpdate
				r, err = svc.update(ctx, r)

			case types.OperationTypeDelete:
				action = RecordActionDelete
				r, err = svc.delete(ctx, r.NamespaceID, r.ModuleID, r.ID)
			}

			if rve := types.IsRecordValueErrorSet(err); rve != nil {
				// Attach additional meta to each value error for FE identification
				for _, re := range rve.Set {
					re.Meta["id"] = p.ID

					rves.Push(re)
				}

				// log record value error for this record
				_ = svc.recordAction(ctx, aProp, action, err)

				// do not return errors just yet, values on other records from the payload (if any)
				// might have errors too
				continue
			}

			if err != nil {
				return svc.recordAction(ctx, aProp, action, err)
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
		return rr, svc.recordAction(ctx, &recordActionProps{}, RecordActionBulk, err)
	}
}

// Raw create function that is responsible for value validation, event dispatching
// and creation.
func (svc record) create(ctx context.Context, new *types.Record) (rec *types.Record, err error) {
	var (
		aProps    = &recordActionProps{changed: new}
		invokerID = auth.GetIdentityFromContext(ctx).Identity()

		ns *types.Namespace
		m  *types.Module
	)

	ns, m, err = loadModuleWithNamespace(ctx, svc.store, new.NamespaceID, new.ModuleID)
	if err != nil {
		return
	}

	aProps.setNamespace(ns)
	aProps.setModule(m)

	if !svc.ac.CanCreateRecord(ctx, m) {
		return nil, RecordErrNotAllowedToCreate()
	}

	if err = RecordValueSanitazion(m, new.Values); err != nil {
		return
	}

	var (
		rve *types.RecordValueErrorSet
	)

	if svc.optEmitEvents {
		if rve = svc.procCreate(ctx, svc.store, invokerID, m, new); !rve.IsValid() {
			return nil, RecordErrValueInput().Wrap(rve)
		}

		if err = svc.eventbus.WaitFor(ctx, event.RecordBeforeCreate(new, nil, m, ns, rve)); err != nil {
			return
		} else if !rve.IsValid() {
			return nil, RecordErrValueInput().Wrap(rve)
		}
	}

	new.Values = RecordValueDefaults(m, new.Values)

	// Handle payload from automation scripts
	if rve = svc.procCreate(ctx, svc.store, invokerID, m, new); !rve.IsValid() {
		return nil, RecordErrValueInput().Wrap(rve)
	}

	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) error {
		return store.CreateComposeRecord(ctx, s, m, new)
	})

	if err != nil {
		return nil, err
	}

	if err = label.Create(ctx, svc.store, new); err != nil {
		return
	}

	// At this point we can return the value
	rec = new

	if svc.optEmitEvents {
		new.Values = svc.formatter.Run(m, new.Values)
		_ = svc.eventbus.WaitFor(ctx, event.RecordAfterCreateImmutable(new, nil, m, ns, nil))
	}

	return
}

// RecordValueSanitazion does basic field and format validation
//
// Received values must fit the data model: on unknown fields
// or multi/single value mismatch we return an error
//
// Record value errors is intentionally NOT used here; if input fails here
// we can assume that form builder (or whatever it was that assembled the record values)
// was misconfigured and will most likely failed to properly parse the
// record value errors payload too
func RecordValueSanitazion(m *types.Module, vv types.RecordValueSet) (err error) {
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

func RecordUpdateOwner(invokerID uint64, r, old *types.Record) *types.Record {
	if old == nil {
		if r.OwnedBy == 0 {
			// If od owner is not set, make current user
			// the owner of the record
			r.OwnedBy = invokerID
		}
	} else {
		if r.OwnedBy == 0 {
			if old.OwnedBy > 0 {
				// Owner not set/send in the payload
				//
				// Fallback to old owner (if set)
				r.OwnedBy = old.OwnedBy
			} else {
				// If od owner is not set, make current user
				// the owner of the record
				r.OwnedBy = invokerID
			}
		}
	}

	return r
}

func RecordValueMerger(ctx context.Context, ac recordValueAccessController, m *types.Module, vv, old types.RecordValueSet) (types.RecordValueSet, *types.RecordValueErrorSet) {
	if old != nil {
		// Value merge process does not know anything about permissions so
		// in case when new values are missing but do exist in the old set and their update/read is denied
		// we need to copy them to ensure value merge process them correctly
		for _, f := range m.Fields {
			if len(vv.FilterByName(f.Name)) == 0 && !ac.CanUpdateRecordValue(ctx, m.Fields.FindByName(f.Name)) {
				// copy all fields from old to new
				vv = append(vv, old.FilterByName(f.Name).GetClean()...)
			}
		}

		// Merge new (updated) values with old ones
		// This way we get list of updated, stale and deleted values
		// that we can selectively update in the repository
		vv = old.Merge(vv)
	}

	rve := &types.RecordValueErrorSet{}
	_ = vv.Walk(func(v *types.RecordValue) error {
		if v.IsUpdated() && !ac.CanUpdateRecordValue(ctx, m.Fields.FindByName(v.Name)) {
			rve.Push(types.RecordValueError{Kind: "updateDenied", Meta: map[string]interface{}{"field": v.Name, "value": v.Value}})
		}

		return nil
	})

	return vv, rve
}

func RecordPreparer(ctx context.Context, s store.Storer, ss recordValuesSanitizer, vv recordValuesValidator, ff recordValuesFormatter, m *types.Module, new *types.Record) *types.RecordValueErrorSet {
	// Before values are processed further and
	// sent to automation scripts (if any)
	// we need to make sure it does not get un-sanitized data
	new.Values = ss.Run(m, new.Values)

	rve := &types.RecordValueErrorSet{}
	values.Expression(ctx, m, new, nil, rve)

	if !rve.IsValid() {
		return rve
	}

	// Run validation of the updated records
	rve = vv.Run(ctx, s, m, new)
	if !rve.IsValid() {
		return rve
	}

	// Cleanup the values
	new.Values = new.Values.GetClean()

	// Formatting
	new.Values = ff.Run(m, new.Values)

	return nil
}

func RecordValueDefaults(m *types.Module, vv types.RecordValueSet) (out types.RecordValueSet) {
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

// Raw update function that is responsible for value validation, event dispatching
// and update.
func (svc record) update(ctx context.Context, upd *types.Record) (rec *types.Record, err error) {
	var (
		aProps    = &recordActionProps{changed: upd}
		invokerID = auth.GetIdentityFromContext(ctx).Identity()

		ns  *types.Namespace
		m   *types.Module
		old *types.Record
	)

	if upd.ID == 0 {
		return nil, RecordErrInvalidID()
	}

	ns, m, old, err = loadRecordCombo(ctx, svc.store, upd.NamespaceID, upd.ModuleID, upd.ID)
	if err != nil {
		return
	}

	aProps.setNamespace(ns)
	aProps.setModule(m)
	aProps.setRecord(old)

	if !svc.ac.CanUpdateRecord(ctx, m) {
		return nil, RecordErrNotAllowedToUpdate()
	}

	// Test if stale (update has an older version of data)
	if isStale(upd.UpdatedAt, old.UpdatedAt, old.CreatedAt) {
		return nil, RecordErrStaleData()
	}

	if err = RecordValueSanitazion(m, upd.Values); err != nil {
		return
	}

	var (
		rve *types.RecordValueErrorSet
	)

	if svc.optEmitEvents {
		// Handle input payload
		if rve = svc.procUpdate(ctx, svc.store, invokerID, m, upd, old); !rve.IsValid() {
			return nil, RecordErrValueInput().Wrap(rve)
		}

		// Scripts can (besides simple error value) return complex record value error set
		// that is passed back to the UI or any other API consumer
		//
		// rve (record-validation-errorset) struct is passed so it can be
		// used & filled by automation scripts
		if err = svc.eventbus.WaitFor(ctx, event.RecordBeforeUpdate(upd, old, m, ns, rve)); err != nil {
			return
		} else if !rve.IsValid() {
			return nil, RecordErrValueInput().Wrap(rve)
		}
	}

	// Handle payload from automation scripts
	if rve = svc.procUpdate(ctx, svc.store, invokerID, m, upd, old); !rve.IsValid() {
		return nil, RecordErrValueInput().Wrap(rve)
	}

	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) error {
		if label.Changed(old.Labels, upd.Labels) {
			if err = label.Update(ctx, s, upd); err != nil {
				return err
			}
		}

		return store.UpdateComposeRecord(ctx, s, m, upd)
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
		_ = svc.eventbus.WaitFor(ctx, event.RecordAfterUpdateImmutable(upd, old, m, ns, nil))
	}
	return
}

func (svc record) Create(ctx context.Context, new *types.Record) (rec *types.Record, err error) {
	var (
		aProps = &recordActionProps{changed: new}
	)

	err = func() error {
		rec, err = svc.create(ctx, new)
		aProps.setRecord(rec)
		return err
	}()

	return rec, svc.recordAction(ctx, aProps, RecordActionCreate, err)
}

// Runs value sanitization, sets values that should be used
// and validates the final result
//
// This logic is kept in a utility function - it's used in the beginning
// of the creation procedure and after results are back from the automation scripts
//
// Both these points introduce external data that need to be checked fully in the same manner
func (svc record) procCreate(ctx context.Context, s store.Storer, invokerID uint64, m *types.Module, new *types.Record) (rve *types.RecordValueErrorSet) {
	new.Values.SetUpdatedFlag(true)

	// Reset values to new record
	// to make sure nobody slips in something we do not want
	new.ID = nextID()
	new.CreatedBy = invokerID
	new.CreatedAt = *now()
	new.UpdatedAt = nil
	new.UpdatedBy = 0
	new.DeletedAt = nil
	new.DeletedBy = 0

	new = RecordUpdateOwner(invokerID, new, nil)
	new.Values, rve = RecordValueMerger(ctx, svc.ac, m, new.Values, nil)
	if !rve.IsValid() {
		return rve
	}
	rve = RecordPreparer(ctx, svc.store, svc.sanitizer, svc.validator, svc.formatter, m, new)
	return rve
}

func (svc record) Update(ctx context.Context, upd *types.Record) (rec *types.Record, err error) {
	var (
		aProps = &recordActionProps{changed: upd}
	)

	err = func() error {
		rec, err = svc.update(ctx, upd)
		aProps.setRecord(rec)
		return err
	}()

	return rec, svc.recordAction(ctx, aProps, RecordActionUpdate, err)
}

// Runs value sanitization, copies values that should updated
// and validates the final result
//
// This logic is kept in a utility function - it's used in the beginning
// of the update procedure and after results are back from the automation scripts
//
// Both these points introduce external data that need to be checked fully in the same manner
func (svc record) procUpdate(ctx context.Context, s store.Storer, invokerID uint64, m *types.Module, upd *types.Record, old *types.Record) (rve *types.RecordValueErrorSet) {
	// Mark all values as updated (new)
	upd.Values.SetUpdatedFlag(true)

	// Copy values to updated record
	// to make sure nobody slips in something we do not want
	upd.CreatedAt = old.CreatedAt
	upd.CreatedBy = old.CreatedBy
	upd.UpdatedAt = now()
	upd.UpdatedBy = invokerID
	upd.DeletedAt = old.DeletedAt
	upd.DeletedBy = old.DeletedBy

	upd = RecordUpdateOwner(invokerID, upd, old)
	upd.Values, rve = RecordValueMerger(ctx, svc.ac, m, upd.Values, old.Values)
	if !rve.IsValid() {
		return rve
	}
	rve = RecordPreparer(ctx, svc.store, svc.sanitizer, svc.validator, svc.formatter, m, upd)
	return rve
}

func (svc record) recordInfoUpdate(ctx context.Context, r *types.Record) {
	r.UpdatedAt = now()
	r.UpdatedBy = auth.GetIdentityFromContext(ctx).Identity()
}

func (svc record) delete(ctx context.Context, namespaceID, moduleID, recordID uint64) (del *types.Record, err error) {
	var (
		ns *types.Namespace
		m  *types.Module

		invokerID = auth.GetIdentityFromContext(ctx).Identity()
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

	ns, m, del, err = loadRecordCombo(ctx, svc.store, namespaceID, moduleID, recordID)
	if err != nil {
		return nil, err
	}

	if !svc.ac.CanDeleteRecord(ctx, m) {
		return nil, RecordErrNotAllowedToDelete()
	}

	if svc.optEmitEvents {
		// Calling before-record-delete scripts
		if err = svc.eventbus.WaitFor(ctx, event.RecordBeforeDelete(nil, del, m, ns, nil)); err != nil {
			return nil, err
		}
	}

	del.DeletedAt = now()
	del.DeletedBy = invokerID

	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) error {
		return store.UpdateComposeRecord(ctx, s, m, del)
	})

	if err != nil {
		return nil, err
	}

	if svc.optEmitEvents {
		_ = svc.eventbus.WaitFor(ctx, event.RecordAfterDeleteImmutable(nil, del, m, ns, nil))
	}

	return del, nil
}

// DeleteByID removes one or more records (all from the same module and namespace)
//
// Before and after each record is deleted beforeDelete and afterDelete events are emitted
// If beforeRecord aborts the action it does so for that specific record only
func (svc record) DeleteByID(ctx context.Context, namespaceID, moduleID uint64, recordIDs ...uint64) (err error) {
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

		ns, m, err = loadModuleWithNamespace(ctx, svc.store, namespaceID, moduleID)
		if err != nil {
			return err
		}

		aProps.setNamespace(ns)
		aProps.setModule(m)

		if !svc.ac.CanDeleteRecord(ctx, m) {
			return RecordErrNotAllowedToDelete()
		}

		return nil
	}()

	if err != nil {
		return svc.recordAction(ctx, aProps, RecordActionDelete, err)
	}

	for _, recordID := range recordIDs {
		err := func() (err error) {
			r, err = svc.delete(ctx, namespaceID, moduleID, recordID)
			aProps.setRecord(r)

			// Record each record deletion action
			return svc.recordAction(ctx, aProps, RecordActionDelete, err)
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

func (svc record) Organize(ctx context.Context, namespaceID, moduleID, recordID uint64, posField, position, filter, grpField, group string) (err error) {
	var (
		ns *types.Namespace
		m  *types.Module
		r  *types.Record

		recordValues = types.RecordValueSet{}

		aProps = &recordActionProps{record: &types.Record{NamespaceID: namespaceID, ModuleID: moduleID, ID: recordID}}

		reorderingRecords bool
	)

	err = func() error {
		ns, m, r, err = loadRecordCombo(ctx, svc.store, namespaceID, moduleID, recordID)
		if err != nil {
			return err
		}

		aProps.setNamespace(ns)
		aProps.setModule(m)
		aProps.setRecord(r)

		if !svc.ac.CanUpdateRecord(ctx, m) {
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

			if !svc.ac.CanUpdateRecordValue(ctx, sf) {
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

			if !svc.ac.CanUpdateRecordValue(ctx, vf) {
				return RecordErrNotAllowedToUpdate()
			}

			// Set new value
			recordValues = recordValues.Set(&types.RecordValue{
				RecordID: recordID,
				Name:     grpField,
				Value:    group,
			})
		}

		return store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) error {
			if len(recordValues) > 0 {
				//svc.recordInfoUpdate(r)
				//if err = store.UpdateComposeRecord(ctx, s, m, r); err != nil {
				//	return err
				//}
			}

			if err = store.PartialComposeRecordValueUpdate(ctx, s, m, recordValues...); err != nil {
				return err
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
				reorderFilter := types.RecordFilter{}
				reorderFilter.Query = fmt.Sprintf("%s(%s >= %d)", filter, posField, recordOrderPlace)
				if err = reorderFilter.Sort.Set(posField); err != nil {
					return err
				}

				set, _, err = store.SearchComposeRecords(ctx, s, m, reorderFilter)
				if err != nil {
					return err
				}

				// Update value on each record
				var vv = make([]*types.RecordValue, 0, len(set))
				_ = set.Walk(func(r *types.Record) error {
					recordOrderPlace++
					vv = append(vv, &types.RecordValue{
						RecordID: r.ID,
						Name:     posField,
						Value:    strconv.FormatUint(recordOrderPlace, 10),
					})

					return nil
				})

				if err = store.PartialComposeRecordValueUpdate(ctx, s, m, vv...); err != nil {
					return err
				}
			}

			return nil
		})
	}()

	return svc.recordAction(ctx, aProps, RecordActionOrganize, err)
}

func (svc record) Validate(ctx context.Context, rec *types.Record) error {
	if m, err := loadModule(ctx, svc.store, rec.ModuleID); err != nil {
		return err
	} else {
		rec.Values = values.Sanitizer().Run(m, rec.Values)
		return values.Validator().Run(ctx, svc.store, m, rec)
	}
}

// TriggerScript loads requested record sanitizes and validates values and passes all to the automation script
//
// For backward compatibility (of controllers), it returns module+record
func (svc record) TriggerScript(ctx context.Context, namespaceID, moduleID, recordID uint64, rvs types.RecordValueSet, script string) (*types.Module, *types.Record, error) {
	var (
		ns, m, r, err = loadRecordCombo(ctx, svc.store, namespaceID, moduleID, recordID)
	)

	if err != nil {
		return nil, nil, err
	}

	original := r.Clone()
	r.Values = values.Sanitizer().Run(m, rvs)
	validated := values.Validator().Run(ctx, svc.store, m, r)

	err = corredor.Service().Exec(ctx, script, event.RecordOnManual(r, original, m, ns, validated))
	if err != nil {
		return nil, nil, err
	}

	return m, r, nil
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
func (svc record) Iterator(ctx context.Context, f types.RecordFilter, fn eventbus.HandlerFn, action string) (err error) {
	var (
		invokerID = auth.GetIdentityFromContext(ctx).Identity()

		ns  *types.Namespace
		m   *types.Module
		set types.RecordSet

		aProps = &recordActionProps{}
	)

	err = func() error {
		ns, m, err = loadModuleWithNamespace(ctx, svc.store, f.NamespaceID, f.ModuleID)
		if err != nil {
			return err
		}

		if !svc.ac.CanReadRecord(ctx, m) {
			return RecordErrNotAllowedToRead()
		}

		switch action {
		case "clone":
			if !svc.ac.CanCreateRecord(ctx, m) {
				return RecordErrNotAllowedToCreate()
			}

		case "update":
			if !svc.ac.CanUpdateRecord(ctx, m) {
				return RecordErrNotAllowedToUpdate()
			}

		case "delete":
			if !svc.ac.CanDeleteRecord(ctx, m) {
				return RecordErrNotAllowedToDelete()
			}
		}

		// @todo might be good to split set into smaller chunks
		set, f, err = store.SearchComposeRecords(ctx, svc.store, m, f)
		if err != nil {
			return err
		}

		for _, rec := range set {
			recordableAction := RecordActionIteratorIteration

			err = func() error {
				if err = fn(ctx, event.RecordOnIteration(rec, nil, m, ns, nil)); err != nil {
					if errors.Is(err, corredor.ScriptExecAborted) {
						// When script was softly aborted (return false),
						// proceed with iteration but do not clone, update or delete
						// current record!
						return nil
					}
				}

				switch action {
				case "clone":
					recordableAction = RecordActionIteratorClone

					// Assign defaults (only on missing values)
					rec.Values = RecordValueDefaults(m, rec.Values)

					// Handle payload from automation scripts
					if rve := svc.procCreate(ctx, svc.store, invokerID, m, rec); !rve.IsValid() {
						return RecordErrValueInput().Wrap(rve)
					}

					return store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) error {
						return store.CreateComposeRecord(ctx, s, m, rec)
					})
				case "update":
					recordableAction = RecordActionIteratorUpdate

					// Handle input payload
					if rve := svc.procUpdate(ctx, svc.store, invokerID, m, rec, rec); !rve.IsValid() {
						return RecordErrValueInput().Wrap(rve)
					}

					return store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) error {
						return store.UpdateComposeRecord(ctx, s, m, rec)
					})
				case "delete":
					recordableAction = RecordActionIteratorDelete

					return store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) error {
						rec.DeletedAt = now()
						rec.DeletedBy = invokerID
						return store.UpdateComposeRecord(ctx, s, m, rec)
					})
				}

				return nil
			}()

			// record iteration action and
			// break the loop in case of an error
			_ = svc.recordAction(ctx, aProps, recordableAction, err)
			if err != nil {
				return err
			}
		}

		return nil
	}()

	return svc.recordAction(ctx, aProps, RecordActionIteratorInvoked, err)

}

// checks record-value-read access permissions for all module fields and removes unreadable fields from all records
func trimUnreadableRecordFields(ctx context.Context, ac recordValueAccessController, m *types.Module, rr ...*types.Record) {
	var (
		readableFields = map[string]bool{}
	)

	for _, f := range m.Fields {
		readableFields[f.Name] = ac.CanReadRecordValue(ctx, f)
	}

	for _, r := range rr {
		r.Values, _ = r.Values.Filter(func(v *types.RecordValue) (bool, error) {
			return readableFields[v.Name], nil
		})
	}
}

// loadRecordCombo Loads namespace, module and record
func loadRecordCombo(ctx context.Context, s store.Storer, namespaceID, moduleID, recordID uint64) (ns *types.Namespace, m *types.Module, r *types.Record, err error) {
	if ns, m, err = loadModuleWithNamespace(ctx, s, namespaceID, moduleID); err != nil {
		return
	}

	if r, err = store.LookupComposeRecordByID(ctx, s, m, recordID); err != nil {
		return
	}

	if r.ModuleID != moduleID {
		return nil, nil, nil, RecordErrInvalidModuleID()
	}

	return
}

// toLabeledRecords converts to []label.LabeledResource
//
// This function is auto-generated.
func toLabeledRecords(set []*types.Record) []label.LabeledResource {
	if len(set) == 0 {
		return nil
	}

	ll := make([]label.LabeledResource, len(set))
	for i := range set {
		ll[i] = set[i]
	}

	return ll
}
