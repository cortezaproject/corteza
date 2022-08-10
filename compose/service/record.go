package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/revisions"
	"regexp"
	"sort"
	"strconv"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/dal"
	"github.com/cortezaproject/corteza-server/pkg/locale"

	"github.com/cortezaproject/corteza-server/compose/dalutils"
	"github.com/cortezaproject/corteza-server/compose/service/event"
	"github.com/cortezaproject/corteza-server/compose/service/values"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/corredor"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/cortezaproject/corteza-server/pkg/report"
	"github.com/cortezaproject/corteza-server/store"
)

const (
	IMPORT_ON_ERROR_SKIP         = "SKIP"
	IMPORT_ON_ERROR_FAIL         = "FAIL"
	IMPORT_ERROR_MAX_INDEX_COUNT = 500000
)

type (
	record struct {
		dal dalDater

		actionlog actionlog.Recorder

		ac       recordAccessController
		eventbus eventDispatcher

		store store.Storer

		namespace namespaceFinder
		module    moduleFinder

		revisions *recordRevisions

		formatter recordValuesFormatter
		sanitizer recordValuesSanitizer
		validator recordValuesValidator
	}

	recordValuesFormatter interface {
		Run(*types.Module, types.RecordValueSet) types.RecordValueSet
	}

	recordValuesSanitizer interface {
		Run(*types.Module, types.RecordValueSet) types.RecordValueSet
		RunXSS(*types.Module, types.RecordValueSet) types.RecordValueSet
	}

	recordValuesValidator interface {
		Run(context.Context, store.Storer, *types.Module, *types.Record) *types.RecordValueErrorSet
		UniqueChecker(fn values.UniqueChecker)
		RecordRefChecker(fn values.ReferenceChecker)
		UserRefChecker(fn values.ReferenceChecker)
	}

	recordValueAccessController interface {
		CanReadRecordValueOnModuleField(context.Context, *types.ModuleField) bool
		CanUpdateRecordValueOnModuleField(context.Context, *types.ModuleField) bool
	}

	recordManageOwnerAccessController interface {
		CanManageOwnerOnRecord(context.Context, *types.Record) bool
		CanCreateOwnedRecordOnModule(context.Context, *types.Module) bool
	}

	recordAccessController interface {
		CanCreateRecordOnModule(context.Context, *types.Module) bool
		CanSearchRecordsOnModule(context.Context, *types.Module) bool
		CanReadNamespace(context.Context, *types.Namespace) bool
		CanReadModule(context.Context, *types.Module) bool
		CanReadRecord(context.Context, *types.Record) bool
		CanUpdateRecord(context.Context, *types.Record) bool
		CanDeleteRecord(context.Context, *types.Record) bool
		CanSearchRevisionsOnRecord(context.Context, *types.Record) bool

		recordManageOwnerAccessController
		recordValueAccessController
	}

	moduleFinder interface {
		Find(ctx context.Context, filter types.ModuleFilter) (set types.ModuleSet, f types.ModuleFilter, err error)
	}

	namespaceFinder interface {
		Find(context.Context, types.NamespaceFilter) (types.NamespaceSet, types.NamespaceFilter, error)
	}

	RecordService interface {
		FindByID(ctx context.Context, namespaceID, moduleID, recordID uint64) (*types.Record, error)

		Report(ctx context.Context, namespaceID, moduleID uint64, metrics, dimensions, filter string) (interface{}, error)
		Find(ctx context.Context, filter types.RecordFilter) (set types.RecordSet, f types.RecordFilter, err error)
		FindSensitive(ctx context.Context) (set []types.SensitiveRecordSet, err error)
		SearchRevisions(ctx context.Context, namespaceID, moduleID, recordID uint64) (dal.Iterator, error)
		RecordExport(context.Context, types.RecordFilter) error
		RecordImport(context.Context, error) error

		Datasource(context.Context, *report.LoadStepDefinition) (report.Datasource, error)

		Create(ctx context.Context, record *types.Record) (*types.Record, error)
		Update(ctx context.Context, record *types.Record) (*types.Record, error)
		Bulk(ctx context.Context, oo ...*types.RecordBulkOperation) (types.RecordSet, error)

		Validate(ctx context.Context, rec *types.Record) error

		DeleteByID(ctx context.Context, namespaceID, moduleID uint64, recordID ...uint64) error

		Organize(ctx context.Context, namespaceID, moduleID, recordID uint64, sortingField, sortingValue, sortingFilter, valueField, value string) error

		Iterator(ctx context.Context, f types.RecordFilter, fn eventbus.HandlerFn, action string) (err error)

		TriggerScript(ctx context.Context, namespaceID, moduleID, recordID uint64, rvs types.RecordValueSet, script string) (*types.Module, *types.Record, error)
	}

	recordImportSession struct {
		Name        string `json:"-"`
		SessionID   uint64 `json:"sessionID,string"`
		UserID      uint64 `json:"userID,string"`
		NamespaceID uint64 `json:"namespaceID,string"`
		ModuleID    uint64 `json:"moduleID,string"`

		OnError  string                `json:"onError"`
		Fields   map[string]string     `json:"fields"`
		Key      string                `json:"key"`
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

		FailLog *FailLog `json:"failLog,omitempty"`
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

func Record() *record {
	svc := &record{
		actionlog: DefaultActionlog,
		ac:        DefaultAccessControl,
		eventbus:  eventbus.Service(),
		store:     DefaultStore,
		dal:       dal.Service(),

		namespace: DefaultNamespace,
		module:    DefaultModule,

		revisions: &recordRevisions{revisions.Service(dal.Service())},

		formatter: values.Formatter(),
		sanitizer: values.Sanitizer(),
	}

	svc.validator = defaultValidator(svc)

	return svc
}

func defaultValidator(svc RecordService) recordValuesValidator {
	// Initialize validator and setup all checkers it needs
	validator := values.Validator()

	validator.UniqueChecker(func(ctx context.Context, s store.Storer, v *types.RecordValue, f *types.ModuleField, m *types.Module) (uint64, error) {
		if v.Ref == 0 {
			return 0, nil
		}

		// @todo re-implement record-value ref lookup through DAL
		panic("implement me")
	})

	validator.RecordRefChecker(func(ctx context.Context, s store.Storer, v *types.RecordValue, f *types.ModuleField, m *types.Module) (bool, error) {
		if svc == nil && v.Ref == 0 {
			return false, nil
		}

		r, err := svc.FindByID(ctx, f.NamespaceID, f.ModuleID, v.Ref)
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

	return validator
}

// lookup fn() orchestrates record lookup, namespace preload and check
func (svc record) lookup(ctx context.Context, namespaceID, moduleID uint64, lookup func(*types.Module, *recordActionProps) (*types.Record, error)) (r *types.Record, err error) {
	var (
		ns     *types.Namespace
		m      *types.Module
		aProps = &recordActionProps{record: &types.Record{NamespaceID: namespaceID}}
	)

	err = func() error {
		if ns, m, err = loadModuleCombo(ctx, svc.store, namespaceID, moduleID); err != nil {
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

		if !svc.ac.CanReadRecord(ctx, r) {
			return RecordErrNotAllowedToRead()
		}

		ComposeRecordFilterAC(ctx, svc.ac, m, r)

		r.SetModule(m)
		r.Values = svc.sanitizer.RunXSS(m, r.Values)

		return nil
	}()

	return r, svc.recordAction(ctx, aProps, RecordActionLookup, err)
}

func (svc record) FindByID(ctx context.Context, namespaceID, moduleID, recordID uint64) (r *types.Record, err error) {
	return svc.lookup(ctx, namespaceID, moduleID, func(m *types.Module, props *recordActionProps) (*types.Record, error) {
		props.record.ID = recordID

		return dalutils.ComposeRecordsFind(ctx, svc.dal, m, recordID)
	})
}

// Report generates report for a given module using metrics, dimensions and filter
//
// @todo remove; will be replaced with reporter endpoints
func (svc record) Report(ctx context.Context, namespaceID, moduleID uint64, metrics, dimensions, filter string) (out interface{}, err error) {
	err = fmt.Errorf("record report endpoints pending removal")
	return
}

func (svc record) Find(ctx context.Context, filter types.RecordFilter) (set types.RecordSet, f types.RecordFilter, err error) {
	var (
		m      *types.Module
		aProps = &recordActionProps{filter: &filter}
	)

	err = func() error {
		if m, err = loadModule(ctx, svc.store, filter.NamespaceID, filter.ModuleID); err != nil {
			return err
		}

		if !svc.ac.CanSearchRecordsOnModule(ctx, m) {
			return RecordErrNotAllowedToSearch()
		}

		filter.Check = ComposeRecordFilterChecker(ctx, svc.ac, m)

		if set, f, err = dalutils.ComposeRecordsList(ctx, svc.dal, m, filter); err != nil {
			return err
		}

		_ = set.Walk(func(r *types.Record) error {
			r.SetModule(m)
			r.Values = svc.sanitizer.RunXSS(m, r.Values)
			return nil
		})

		ComposeRecordFilterAC(ctx, svc.ac, m, set...)

		return nil
	}()

	return set, f, svc.recordAction(ctx, aProps, RecordActionSearch, err)
}

// FindSensitive returns stripped down records for all namespaces/modules where fields define a sensitivity level
func (svc record) FindSensitive(ctx context.Context) (set []types.SensitiveRecordSet, err error) {
	var (
		namespaces types.NamespaceSet
		modules    types.ModuleSet

		userID = auth.GetIdentityFromContext(ctx).Identity()
	)

	err = func() error {
		// Get namespaces
		namespaces, _, err = svc.namespace.Find(ctx, types.NamespaceFilter{})
		if err != nil {
			return err
		}

		for _, namespace := range namespaces {
			// Get corresponding modules
			modules, _, err = svc.module.Find(ctx, types.ModuleFilter{NamespaceID: namespace.ID})
			if err != nil {
				return err
			}

			for _, module := range modules {
				// Get sensitive record data
				aux, err := svc.findSensitive(ctx, userID, namespace, module, types.RecordFilter{
					ModuleID:    module.ID,
					NamespaceID: namespace.ID,
				})
				if err != nil {
					return err
				}
				if len(aux.Records) == 0 {
					continue
				}

				set = append(set, aux)
			}
		}

		return nil
	}()

	return set, err
}

func (svc record) findSensitive(ctx context.Context, userID uint64, namespace *types.Namespace, module *types.Module, filter types.RecordFilter) (out types.SensitiveRecordSet, err error) {
	out = types.SensitiveRecordSet{
		Namespace: namespace,
		Module:    module,
	}

	// Force the query to only show owned records
	// @todo allow additional querying
	filter.Query = fmt.Sprintf("ownedBy='%d'", userID)

	rr, _, err := svc.Find(ctx, filter)
	if err != nil {
		return
	}

	for _, r := range rr {
		vv := make([]map[string]any, 0, len(r.Values))

		for _, f := range module.Fields {
			// Skip the ones with no privacy
			// @todo allow the request to specify what level we wish to see
			if !f.IsSensitive() {
				continue
			}

			values := make([]any, 0, 2)
			for _, v := range r.Values.FilterByName(f.Name) {
				values = append(values, v.Value)
			}

			// Make value
			vv = append(vv, map[string]any{
				"name":    f.Name,
				"kind":    f.Kind,
				"isMulti": f.Multi,
				"value":   values,
			})
		}

		out.Records = append(out.Records, types.SensitiveRecord{
			RecordID: r.ID,
			Values:   vv,
		})
	}

	return
}

// SearchRevisions returns iterator for revisions of a record
//
func (svc record) SearchRevisions(ctx context.Context, namespaceID, moduleID, recordID uint64) (dal.Iterator, error) {
	var (
		aProps = &recordActionProps{record: &types.Record{NamespaceID: namespaceID, ModuleID: moduleID, ID: recordID}}

		rec  *types.Record
		iter dal.Iterator
		mod  *types.Module
		ns   *types.Namespace
	)

	err := func() (err error) {
		ns, mod, rec, err = loadRecordCombo(ctx, svc.store, svc.dal, namespaceID, moduleID, recordID)
		if err != nil {
			return
		}

		aProps.setModule(mod)
		aProps.setNamespace(ns)
		aProps.setRecord(rec)

		if !mod.Config.RecordRevisions.Enabled {
			return RecordErrRevisionsDisabledOnModule()
		}

		if !svc.ac.CanSearchRevisionsOnRecord(ctx, rec) {
			return RecordErrNotAllowedToSearchRevisions()
		}

		iter, err = svc.revisions.search(ctx, rec)
		return
	}()

	return iter, svc.recordAction(ctx, aProps, RecordActionSearchRevisions, err)
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

			aProp.setRecord(r)

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

			aProp.setChanged(r)

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

			_ = svc.recordAction(ctx, aProp, action, err)
			if err != nil {
				return err
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
		// was not really a bulk operation, and we already recorded the action
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
		aProps    = &recordActionProps{record: new}
		invokerID = auth.GetIdentityFromContext(ctx).Identity()

		ns *types.Namespace
		m  *types.Module
	)

	ns, m, err = loadModuleCombo(ctx, svc.store, new.NamespaceID, new.ModuleID)
	if err != nil {
		return
	}

	aProps.setNamespace(ns)
	aProps.setModule(m)

	if !svc.ac.CanCreateRecordOnModule(ctx, m) {
		return nil, RecordErrNotAllowedToCreate()
	}

	if err = RecordValueSanitization(m, new.Values); err != nil {
		return
	}

	var (
		rve *types.RecordValueErrorSet
	)

	// ensure module ref is set before running through records workflows and scripts
	new.SetModule(m)

	{
		if rve = svc.procCreate(ctx, invokerID, m, new); !rve.IsValid() {
			return nil, RecordErrValueInput().Wrap(rve)
		}

		if err = svc.eventbus.WaitFor(ctx, event.RecordBeforeCreate(new, nil, m, ns, rve, nil)); err != nil {
			return
		} else if !rve.IsValid() {
			return nil, RecordErrValueInput().Wrap(rve)
		}
	}

	new.Values = RecordValueDefaults(m, new.Values)

	// Handle payload from automation scripts
	if rve = svc.procCreate(ctx, invokerID, m, new); !rve.IsValid() {
		return nil, RecordErrValueInput().Wrap(rve)
	}

	aProps.setChanged(new)

	if err = dalutils.ComposeRecordCreate(ctx, svc.dal, m, new); err != nil {
		return
	}

	// store revision
	if m.Config.RecordRevisions.Enabled {
		if err = svc.revisions.created(ctx, new); err != nil {
			return
		}
	}

	// ensure module ref is set before running through records workflows and scripts
	new.SetModule(m)

	// At this point we can return the value
	rec = new

	{
		new.Values = svc.formatter.Run(m, new.Values)
		_ = svc.eventbus.WaitFor(ctx, event.RecordAfterCreateImmutable(new, nil, m, ns, nil, nil))
	}

	return
}

// RecordValueSanitization does basic field and format validation
//
// Received values must fit the data model: on unknown fields
// or multi/single value mismatch we return an error
//
// Record value errors is intentionally NOT used here; if input fails here
// we can assume that form builder (or whatever it was that assembled the record values)
// was misconfigured and will most likely failed to properly parse the
// record value errors payload too
func RecordValueSanitization(m *types.Module, vv types.RecordValueSet) (err error) {
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
			if v.Value == "" || v.Value == "0" {
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

func SetRecordOwner(ctx context.Context, ac recordManageOwnerAccessController, s store.Storer, old, upd *types.Record, invoker uint64) *types.RecordValueErrorSet {
	if upd == nil {
		// no-op
		return nil
	}

	var (
		curOwner uint64
		updOwner = upd.OwnedBy
	)

	if old != nil {
		curOwner = old.OwnedBy
	}

	updOwner = CalcRecordOwner(curOwner, updOwner, invoker)

	var (
		mkError = func(kind, tkey string) *types.RecordValueErrorSet {
			return &types.RecordValueErrorSet{Set: []types.RecordValueError{{
				Kind:    kind,
				Meta:    map[string]interface{}{"field": "", "value": updOwner},
				Message: locale.Global().T(ctx, "compose", tkey),
			}}}
		}
	)

	if old == nil && updOwner != invoker {
		// Can be set on new records
		if !ac.CanCreateOwnedRecordOnModule(ctx, upd.GetModule()) {
			return mkError("accessDenied", "record.errors.ownershipChangeDenied")
		}
	}

	if old != nil && curOwner != updOwner {
		// Can ownership be changed on existing record?
		if !ac.CanManageOwnerOnRecord(ctx, upd) {
			return mkError("accessDenied", "record.errors.ownershipChangeDenied")
		}
	}

	if _, err := store.LookupUserByID(ctx, s, updOwner); err != nil {
		if errors.IsNotFound(err) {
			return mkError("invalidValue", "record.errors.invalidOwner")
		} else {
			return mkError("internal", "record.errors.store")
		}
	}

	upd.OwnedBy = updOwner
	return nil
}

func CalcRecordOwner(current, new, invoker uint64) uint64 {
	if invoker == 0 {
		// invoker is, for some reason 0,
		// use current owner as invoker
		invoker = current
	}

	if new == 0 {
		// if new owner is not set, use invoker
		return invoker
	}

	// keep owner unchanged
	return new
}

func RecordValueUpdateOpCheck(ctx context.Context, ac recordValueAccessController, m *types.Module, vv types.RecordValueSet) *types.RecordValueErrorSet {
	rve := &types.RecordValueErrorSet{}
	if ac == nil {
		return rve
	}

	_ = vv.Walk(func(v *types.RecordValue) error {
		f := m.Fields.FindByName(v.Name)

		// when f is nil, the module field was deleted so we shouldn't do any AC
		if f == nil {
			return nil
		}

		if v.IsUpdated() && !ac.CanUpdateRecordValueOnModuleField(ctx, f) {
			rve.Push(types.RecordValueError{Kind: "updateDenied", Meta: map[string]interface{}{"field": v.Name, "value": v.Value}})
		}

		return nil
	})

	return rve
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
		aProps    = &recordActionProps{record: upd}
		invokerID = auth.GetIdentityFromContext(ctx).Identity()

		ns  *types.Namespace
		m   *types.Module
		old *types.Record
	)

	if upd.ID == 0 {
		return nil, RecordErrInvalidID()
	}

	ns, m, old, err = loadRecordCombo(ctx, svc.store, svc.dal, upd.NamespaceID, upd.ModuleID, upd.ID)
	if err != nil {
		return
	}

	aProps.setNamespace(ns)
	aProps.setModule(m)
	aProps.setRecord(old)

	if !svc.ac.CanUpdateRecord(ctx, old) {
		return nil, RecordErrNotAllowedToUpdate()
	}

	// Test if stale (update has an older version of data)
	if isStale(upd.UpdatedAt, old.UpdatedAt, old.CreatedAt) {
		return nil, RecordErrStaleData()
	}

	if err = RecordValueSanitization(m, upd.Values); err != nil {
		return
	}

	var (
		rve *types.RecordValueErrorSet
	)

	// ensure module ref is set before running through records workflows and scripts
	upd.SetModule(m)
	old.SetModule(m)

	{
		// Handle input payload
		if rve = svc.procUpdate(ctx, invokerID, m, upd, old); !rve.IsValid() {
			return nil, RecordErrValueInput().Wrap(rve)
		}

		// Scripts can (besides simple error value) return complex record value error set
		// that is passed back to the UI or any other API consumer
		//
		// rve (record-validation-errorset) struct is passed so it can be
		// used & filled by automation scripts
		if err = svc.eventbus.WaitFor(ctx, event.RecordBeforeUpdate(upd, old, m, ns, rve, nil)); err != nil {
			return
		} else if !rve.IsValid() {
			return nil, RecordErrValueInput().Wrap(rve)
		}
	}

	// Handle payload from automation scripts
	if rve = svc.procUpdate(ctx, invokerID, m, upd, old); !rve.IsValid() {
		return nil, RecordErrValueInput().Wrap(rve)
	}

	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) error {
		aProps.setChanged(upd)

		if err = dalutils.ComposeRecordUpdate(ctx, svc.dal, m, upd); err != nil {
			return err
		}

		if m.Config.RecordRevisions.Enabled {
			// Prepare record revision for update
			if err = svc.revisions.updated(ctx, upd, old); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// ensure module ref is set before running through records workflows and scripts
	upd.SetModule(m)
	old.SetModule(m)

	// Final value cleanup
	// These (clean) values are returned (and sent to after-update handler)
	upd.Values = upd.Values.GetClean()

	// At this point we can return the value
	rec = upd

	{
		// Before we pass values to automation scripts, they should be formatted
		upd.Values = svc.formatter.Run(m, upd.Values)
		_ = svc.eventbus.WaitFor(ctx, event.RecordAfterUpdateImmutable(upd, old, m, ns, nil, nil))
	}
	return
}

func (svc record) Create(ctx context.Context, new *types.Record) (rec *types.Record, err error) {
	var (
		aProps = &recordActionProps{record: new}
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
func (svc record) procCreate(ctx context.Context, invokerID uint64, m *types.Module, new *types.Record) (rve *types.RecordValueErrorSet) {
	new.Values.SetUpdatedFlag(true)

	new.Values.Walk(func(v *types.RecordValue) error {
		f := m.Fields.FindByName(v.Name)
		if f == nil {
			return nil
		}

		d := f.DefaultValue.Get("", v.Place)
		if d == nil {
			// just so that we do not miss any defaults that MIGHT have
			// field name set to it
			// this is highly unlikely but it does not hurt to try
			d = f.DefaultValue.Get(v.Name, v.Place)
		}

		// Mark as updated ONLY if value set is different from the default one
		v.Updated = d == nil || d.Value != v.Value

		return nil
	})

	// Reset values to new record
	// to make sure nobody slips in something we do not want
	new.ID = nextID()
	new.Revision = 1
	new.CreatedBy = invokerID
	new.CreatedAt = *nowUTC()
	new.UpdatedAt = nil
	new.UpdatedBy = 0
	new.DeletedAt = nil
	new.DeletedBy = 0

	if err := SetRecordOwner(ctx, svc.ac, svc.store, nil, new, invokerID); err != nil {
		return err
	}

	if rve = RecordValueUpdateOpCheck(ctx, svc.ac, m, new.Values); !rve.IsValid() {
		return
	}

	rve = RecordPreparer(ctx, svc.store, svc.sanitizer, svc.validator, svc.formatter, m, new)
	return rve
}

func (svc record) Update(ctx context.Context, upd *types.Record) (rec *types.Record, err error) {
	var (
		aProps = &recordActionProps{record: upd}
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
func (svc record) procUpdate(ctx context.Context, invokerID uint64, m *types.Module, upd *types.Record, old *types.Record) (rve *types.RecordValueErrorSet) {
	upd.Revision = old.Revision + 1

	// Mark all values as updated (new)
	upd.Values.SetUpdatedFlag(true)

	// First sanitization
	//
	// Before values are merged with existing data and
	// sent to automation scripts (if any)
	// we need to make sure it does not get sanitized data
	upd.Values = svc.sanitizer.Run(m, upd.Values)

	if upd.Meta == nil {
		// meta set to nil means we need to keep the old values!
		upd.Meta = old.Meta
	}

	// Copy values to updated record
	// to make sure nobody slips in something we do not want
	upd.CreatedAt = old.CreatedAt
	upd.CreatedBy = old.CreatedBy
	upd.UpdatedAt = nowUTC()
	upd.UpdatedBy = invokerID
	upd.DeletedAt = old.DeletedAt
	upd.DeletedBy = old.DeletedBy

	if rve = SetRecordOwner(ctx, svc.ac, svc.store, old, upd, invokerID); !rve.IsValid() {
		return
	}

	upd.Values = old.Values.Merge(m.Fields, upd.Values, func(f *types.ModuleField) bool {
		return svc.ac.CanUpdateRecordValueOnModuleField(ctx, m.Fields.FindByName(f.Name))
	})

	if rve = RecordValueUpdateOpCheck(ctx, svc.ac, m, upd.Values); !rve.IsValid() {
		return
	}

	rve = RecordPreparer(ctx, svc.store, svc.sanitizer, svc.validator, svc.formatter, m, upd)
	return
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

	ns, m, del, err = loadRecordCombo(ctx, svc.store, svc.dal, namespaceID, moduleID, recordID)
	if err != nil {
		return nil, err
	}

	if !svc.ac.CanDeleteRecord(ctx, del) {
		return nil, RecordErrNotAllowedToDelete()
	}

	del.DeletedAt = nowUTC()
	del.DeletedBy = invokerID

	// ensure module ref is set before running through records workflows and scripts
	del.SetModule(m)

	// deleted, revision need to be set when RecordBeforeDelete is triggered
	del.DeletedAt = nowUTC()
	del.DeletedBy = invokerID
	del.Revision = del.Revision + 1

	{
		// Calling before-record-delete scripts
		if err = svc.eventbus.WaitFor(ctx, event.RecordBeforeDelete(nil, del, m, ns, nil, nil)); err != nil {
			return nil, err
		}
	}

	if m.Config.RecordRevisions.Enabled {
		// Prepare record revision for update
		if err = svc.revisions.softDeleted(ctx, del); err != nil {
			return
		}
	}

	if err = dalutils.ComposeRecordSoftDelete(ctx, svc.dal, m, del); err != nil {
		return nil, err
	}

	// ensure module ref is set before running through records workflows and scripts
	del.SetModule(m)

	{
		_ = svc.eventbus.WaitFor(ctx, event.RecordAfterDeleteImmutable(nil, del, m, ns, nil, nil))
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

		ns, m, err = loadModuleCombo(ctx, svc.store, namespaceID, moduleID)
		if err != nil {
			return err
		}

		aProps.setNamespace(ns)
		aProps.setModule(m)

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

		aProps = &recordActionProps{record: &types.Record{NamespaceID: namespaceID, ModuleID: moduleID, ID: recordID}}

		reorderingRecords bool
	)

	err = func() error {
		ns, m, r, err = loadRecordCombo(ctx, svc.store, svc.dal, namespaceID, moduleID, recordID)
		if err != nil {
			return err
		}

		aProps.setNamespace(ns)
		aProps.setModule(m)
		aProps.setRecord(r)

		if !svc.ac.CanUpdateRecord(ctx, r) {
			return RecordErrNotAllowedToUpdate()
		}

		if posField != "" {
			reorderingRecords = true

			// Position field checks
			{
				// Check field existence and permissions
				// check if numeric -- we cannot reorder on any other field type
				sf := m.Fields.FindByName(posField)
				if sf == nil {
					return RecordErrMissingPositionField()
				}

				aProps.setPositionField(sf)

				if !sf.IsNumeric() {
					return RecordErrInvalidPositionFieldKind()
				}

				if sf.Multi {
					return RecordErrInvalidPositionFieldConfigMultiValue()
				}

				if !svc.ac.CanUpdateRecordValueOnModuleField(ctx, sf) {
					return RecordErrNotAllowedToUpdate()
				}
			}

			// Value checks
			{
				if !regexp.MustCompile(`^[0-9]+$`).MatchString(position) {
					return RecordErrInvalidPositionValueType()
				}
			}

			// Set new position
			if err = r.SetValue(posField, 0, position); err != nil {
				return err
			}
		}

		if grpField != "" {
			// Group field checks
			{
				vf := m.Fields.FindByName(grpField)
				if vf == nil {
					return RecordErrMissingGroupField()
				}

				aProps.setGroupField(vf)

				if vf.Multi {
					return RecordErrInvalidGroupFieldConfigMultiValue()
				}

				if !svc.ac.CanUpdateRecordValueOnModuleField(ctx, vf) {
					return RecordErrNotAllowedToUpdate()
				}
			}

			// Set new value
			if err = r.SetValue(grpField, 0, group); err != nil {
				return err
			}
		}

		return store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) error {
			if err = dalutils.ComposeRecordUpdate(ctx, svc.dal, m, r); err != nil {
				return err
			}

			if reorderingRecords {
				var (
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
				reorderFilter := types.RecordFilter{
					ModuleID:    moduleID,
					NamespaceID: namespaceID,
				}
				reorderFilter.Query = fmt.Sprintf("%s(%s >= %d)", filter, posField, recordOrderPlace)
				if err = reorderFilter.Sort.Set(posField); err != nil {
					return err
				}

				var iter dal.Iterator
				if iter, _, err = dalutils.ComposeRecordsIterator(ctx, svc.dal, m, reorderFilter); err != nil {
					return err
				}

				const updateChunkSize = 300
				updateChunk := make(types.RecordSet, 0, updateChunkSize)
				err = dalutils.WalkIterator(ctx, iter, m, func(r *types.Record) (err error) {
					recordOrderPlace++
					if err = r.SetValue(posField, 0, strconv.FormatUint(recordOrderPlace, 10)); err != nil {
						return
					}

					// Update in chunks so we don't update one at a time and we don't keep
					// too much data in memory.
					updateChunk = append(updateChunk, r)
					if len(updateChunk) >= updateChunkSize {
						if err = dalutils.ComposeRecordUpdate(ctx, svc.dal, m, updateChunk...); err != nil {
							return
						}
						updateChunk = make(types.RecordSet, 0, updateChunkSize)
					}

					return
				})
				if err != nil {
					return err
				}

				// Assure the last chunk is handled
				if len(updateChunk) > 0 {
					if err = dalutils.ComposeRecordUpdate(ctx, svc.dal, m, updateChunk...); err != nil {
						return err
					}
				}
			}

			return nil
		})
	}()

	return svc.recordAction(ctx, aProps, RecordActionOrganize, err)
}

func (svc record) Validate(ctx context.Context, rec *types.Record) error {
	if m, err := loadModule(ctx, svc.store, rec.NamespaceID, rec.ModuleID); err != nil {
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
		ns, m, r, err = loadRecordCombo(ctx, svc.store, svc.dal, namespaceID, moduleID, recordID)
	)

	if err != nil {
		return nil, nil, err
	}

	original := r.Clone()
	r.Values = values.Sanitizer().Run(m, rvs)
	validated := values.Validator().Run(ctx, svc.store, m, r)

	err = corredor.Service().Exec(ctx, script, event.RecordOnManual(r, original, m, ns, validated, nil))
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

		ns   *types.Namespace
		m    *types.Module
		iter dal.Iterator

		aProps = &recordActionProps{}
	)

	err = func() error {
		ns, m, err = loadModuleCombo(ctx, svc.store, f.NamespaceID, f.ModuleID)
		if err != nil {
			return err
		}

		iter, _, err = dalutils.ComposeRecordsIterator(ctx, svc.dal, m, f)
		if err != nil {
			return err
		}

		err = dalutils.WalkIterator(ctx, iter, m, func(rec *types.Record) (err error) {
			switch action {
			case "clone":
				if !svc.ac.CanCreateRecordOnModule(ctx, m) {
					return RecordErrNotAllowedToCreate()
				}

			case "update":
				if !svc.ac.CanUpdateRecord(ctx, rec) {
					return RecordErrNotAllowedToUpdate()
				}

			case "delete":
				if !svc.ac.CanDeleteRecord(ctx, rec) {
					return RecordErrNotAllowedToDelete()
				}
			}
			recordableAction := RecordActionIteratorIteration

			if !svc.ac.CanReadRecord(ctx, rec) {
				return RecordErrNotAllowedToRead()
			}

			err = func() error {
				if err = fn(ctx, event.RecordOnIteration(rec, nil, m, ns, nil, nil)); err != nil {
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
					if rve := svc.procCreate(ctx, invokerID, m, rec); !rve.IsValid() {
						return RecordErrValueInput().Wrap(rve)
					}

					return store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) error {
						return dalutils.ComposeRecordCreate(ctx, svc.dal, m, rec)
					})
				case "update":
					recordableAction = RecordActionIteratorUpdate

					// Handle input payload
					if rve := svc.procUpdate(ctx, invokerID, m, rec, rec); !rve.IsValid() {
						return RecordErrValueInput().Wrap(rve)
					}

					return store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) error {
						return dalutils.ComposeRecordUpdate(ctx, svc.dal, m, rec)
					})
				case "delete":
					recordableAction = RecordActionIteratorDelete

					return store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) error {
						rec.DeletedAt = nowUTC()
						rec.DeletedBy = invokerID
						return dalutils.ComposeRecordSoftDelete(ctx, svc.dal, m, rec)
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

			return
		})
		if err != nil {
			return err
		}

		return nil
	}()

	return svc.recordAction(ctx, aProps, RecordActionIteratorInvoked, err)
}

func ComposeRecordFilterChecker(ctx context.Context, ac recordAccessController, m *types.Module) func(*types.Record) (bool, error) {
	return func(rec *types.Record) (bool, error) {
		// Setting module right before we do access control
		//
		// Why?
		//  - Access control can use one of the contextual roles
		//  - Contextual role can use expression that accesses values
		//  - Record's values are only exported into expression's scope when
		//    module is set on record at the time when Dict() fn is called.
		rec.SetModule(m)

		if !ac.CanReadRecord(ctx, rec) {
			return false, nil
		}

		return true, nil
	}
}

// checks record-value-read access permissions for all module fields and removes unreadable fields from all records
func ComposeRecordFilterAC(ctx context.Context, ac recordValueAccessController, m *types.Module, rr ...*types.Record) {
	var (
		readableFields = map[string]bool{}
	)

	for _, f := range m.Fields {
		readableFields[f.Name] = ac.CanReadRecordValueOnModuleField(ctx, f)
	}

	for _, r := range rr {
		r.Values, _ = r.Values.Filter(func(v *types.RecordValue) (bool, error) {
			return readableFields[v.Name], nil
		})
	}
}

// loadRecordCombo Loads namespace, module and record
func loadRecordCombo(ctx context.Context, s store.Storer, dal dalDater, namespaceID, moduleID, recordID uint64) (ns *types.Namespace, m *types.Module, r *types.Record, err error) {
	if recordID == 0 {
		return nil, nil, nil, RecordErrInvalidID()
	}

	if ns, m, err = loadModuleCombo(ctx, s, namespaceID, moduleID); err != nil {
		return
	}

	if r, err = dalutils.ComposeRecordsFind(ctx, dal, m, recordID); err != nil {
		return
	}

	if r.ModuleID != moduleID {
		return nil, nil, nil, RecordErrInvalidModuleID()
	}

	return
}

// loadRecord loads record
//
// function uses global DAL service to load records
// this is because we need to be able to call it from AccessControl service
// that does not have DAL
func loadRecord(ctx context.Context, s store.Storer, namespaceID, moduleID, recordID uint64) (res *types.Record, err error) {
	_, _, res, err = loadRecordCombo(ctx, s, dal.Service(), namespaceID, moduleID, recordID)
	return
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
