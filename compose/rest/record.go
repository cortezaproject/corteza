package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/cortezaproject/corteza-server/compose/rest/request"
	"github.com/cortezaproject/corteza-server/compose/service"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/api"
	"github.com/cortezaproject/corteza-server/pkg/corredor"
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/csv"
	ejson "github.com/cortezaproject/corteza-server/pkg/envoy/json"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	estore "github.com/cortezaproject/corteza-server/pkg/envoy/store"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/payload"
	"github.com/cortezaproject/corteza-server/store"
)

type (
	recordPayload struct {
		*types.Record

		Records types.RecordSet `json:"records,omitempty"`

		CanUpdateRecord bool `json:"canUpdateRecord"`
		CanReadRecord   bool `json:"canReadRecord"`
		CanDeleteRecord bool `json:"canDeleteRecord"`
		CanGrant        bool `json:"canGrant"`
	}

	recordSetPayload struct {
		Filter *types.RecordFilter `json:"filter,omitempty"`
		Set    []*recordPayload    `json:"set"`
	}

	Record struct {
		importSession service.ImportSessionService
		record        service.RecordService
		module        service.ModuleService
		namespace     service.NamespaceService
		attachment    service.AttachmentService
		ac            recordAccessController
	}

	recordAccessController interface {
		CanGrant(context.Context) bool

		CanUpdateRecord(context.Context, *types.Record) bool
		CanReadRecord(context.Context, *types.Record) bool
		CanDeleteRecord(context.Context, *types.Record) bool
	}
)

const (
	defaultRecordSearchSize uint = 500
	maxRecordSearchSize          = 1000
)

func (Record) New() *Record {
	return &Record{
		importSession: service.DefaultImportSession,
		record:        service.DefaultRecord,
		module:        service.DefaultModule,
		namespace:     service.DefaultNamespace,
		attachment:    service.DefaultAttachment,
		ac:            service.DefaultAccessControl,
	}
}

func (ctrl *Record) Report(ctx context.Context, r *request.RecordReport) (interface{}, error) {
	return ctrl.record.Report(ctx, r.NamespaceID, r.ModuleID, r.Metrics, r.Dimensions, r.Filter)
}

func (ctrl *Record) List(ctx context.Context, r *request.RecordList) (interface{}, error) {
	var (
		m   *types.Module
		err error

		f = types.RecordFilter{
			NamespaceID: r.NamespaceID,
			ModuleID:    r.ModuleID,
			Labels:      r.Labels,
			Deleted:     filter.State(r.Deleted),
		}
	)

	if err = f.Sort.Set(r.Sort); err != nil {
		return nil, err
	}

	if m, err = ctrl.module.FindByID(ctx, r.NamespaceID, r.ModuleID); err != nil {
		return nil, err
	}

	if r.Query != "" {
		// Query param takes preference
		f.Query = r.Query
	}

	if r.Limit == 0 {
		r.Limit = defaultRecordSearchSize
	}

	r.Limit = uint(math.Min(float64(r.Limit), float64(maxRecordSearchSize)))

	if f.Paging, err = filter.NewPaging(r.Limit, r.PageCursor); err != nil {
		return nil, err
	}

	f.IncTotal = r.IncTotal
	f.IncPageNavigation = r.IncPageNavigation

	if f.Sorting, err = filter.NewSorting(r.Sort); err != nil {
		return nil, err
	}

	rr, filter, err := ctrl.record.Find(ctx, f)

	return ctrl.makeFilterPayload(ctx, m, rr, &filter, err)
}

func (ctrl *Record) Read(ctx context.Context, r *request.RecordRead) (interface{}, error) {
	var (
		m   *types.Module
		err error
	)

	if m, err = ctrl.module.FindByID(ctx, r.NamespaceID, r.ModuleID); err != nil {
		return nil, err
	}

	record, err := ctrl.record.FindByID(ctx, r.NamespaceID, r.ModuleID, r.RecordID)

	// Temp workaround until we do proper by-module filtering for record findByID
	if record != nil && record.ModuleID != r.ModuleID {
		return nil, store.ErrNotFound
	}

	return ctrl.makePayload(ctx, m, record, err)
}

func (ctrl *Record) Create(ctx context.Context, r *request.RecordCreate) (interface{}, error) {
	var (
		m   *types.Module
		err error
	)

	if m, err = ctrl.module.FindByID(ctx, r.NamespaceID, r.ModuleID); err != nil {
		return nil, err
	}

	oo := make([]*types.RecordBulkOperation, 0)

	// If defined, initialize parent record
	if r.Values != nil {
		rr := &types.Record{
			NamespaceID: r.NamespaceID,
			ModuleID:    r.ModuleID,
			Values:      r.Values,
			Labels:      r.Labels,
			OwnedBy:     r.OwnedBy,
		}
		oo = append(oo, &types.RecordBulkOperation{
			Record:    rr,
			Operation: types.OperationTypeCreate,
			ID:        "parent:0",
		})
	}

	// If defined, initialize sub records for creation
	oob, err := r.Records.ToBulkOperations(r.ModuleID, r.NamespaceID)
	if err != nil {
		return nil, err
	}

	// Validate returned bulk operations
	for _, o := range oob {
		if o.LinkBy != "" && len(oo) == 0 {
			return nil, fmt.Errorf("missing parent record definition")
		}
	}
	oo = append(oo, oob...)

	rr, err := ctrl.record.Bulk(ctx, oo...)
	if rve := types.IsRecordValueErrorSet(err); rve != nil {
		return ctrl.handleValidationError(rve), nil
	}

	return ctrl.makeBulkPayload(ctx, m, err, rr...)
}

func (ctrl *Record) Update(ctx context.Context, r *request.RecordUpdate) (interface{}, error) {
	var (
		m   *types.Module
		err error
	)

	if m, err = ctrl.module.FindByID(ctx, r.NamespaceID, r.ModuleID); err != nil {
		return nil, err
	}

	oo := make([]*types.RecordBulkOperation, 0)

	// If defined, initialize parent record for creation
	if r.Values != nil {
		rr := &types.Record{
			ID:          r.RecordID,
			NamespaceID: r.NamespaceID,
			ModuleID:    r.ModuleID,
			Values:      r.Values,
			Labels:      r.Labels,
			OwnedBy:     r.OwnedBy,
		}
		oo = append(oo, &types.RecordBulkOperation{
			Record:    rr,
			Operation: types.OperationTypeUpdate,
			ID:        strconv.FormatUint(rr.ID, 10),
		})
	}

	// If defined, initialize sub records for creation
	oob, err := r.Records.ToBulkOperations(r.ModuleID, r.NamespaceID)
	if err != nil {
		return nil, err
	}

	// Validate returned bulk operations
	for _, o := range oob {
		if o.LinkBy != "" && len(oo) == 0 {
			return nil, fmt.Errorf("missing parent record definition")
		}
	}
	oo = append(oo, oob...)

	rr, err := ctrl.record.Bulk(ctx, oo...)

	if rve := types.IsRecordValueErrorSet(err); rve != nil {
		return ctrl.handleValidationError(rve), nil
	}

	return ctrl.makeBulkPayload(ctx, m, err, rr...)
}

func (ctrl *Record) Delete(ctx context.Context, r *request.RecordDelete) (interface{}, error) {
	return api.OK(), ctrl.record.DeleteByID(ctx, r.NamespaceID, r.ModuleID, r.RecordID)
}

func (ctrl *Record) BulkDelete(ctx context.Context, r *request.RecordBulkDelete) (interface{}, error) {
	if r.Truncate {
		return nil, fmt.Errorf("pending implementation")
	}

	return api.OK(), ctrl.record.DeleteByID(ctx,
		r.NamespaceID,
		r.ModuleID,
		payload.ParseUint64s(r.RecordIDs)...,
	)
}

func (ctrl *Record) Upload(ctx context.Context, r *request.RecordUpload) (interface{}, error) {
	file, err := r.Upload.Open()
	if err != nil {
		return nil, err
	}

	defer file.Close()

	a, err := ctrl.attachment.CreateRecordAttachment(
		ctx,
		r.NamespaceID,
		r.Upload.Filename,
		r.Upload.Size,
		file,
		r.ModuleID,
		r.RecordID,
		r.FieldName,
	)

	return makeAttachmentPayload(ctx, a, err)
}

func (ctrl *Record) ImportInit(ctx context.Context, r *request.RecordImportInit) (interface{}, error) {
	if _, err := ctrl.module.FindByID(ctx, r.NamespaceID, r.ModuleID); err != nil {
		return nil, err
	}

	f, err := r.Upload.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Mime type detection library fails for some .csv files, so let's help them out a bit.
	// The detection can now fallback to the user-provided content-type.
	ct := r.Upload.Header.Get("Content-Type")
	return ctrl.importSession.Create(ctx, f, r.Upload.Filename, ct, r.NamespaceID, r.ModuleID)
}

func (ctrl *Record) ImportRun(ctx context.Context, r *request.RecordImportRun) (interface{}, error) {
	var (
		err error
	)

	// Access control.
	if _, err = ctrl.module.FindByID(ctx, r.NamespaceID, r.ModuleID); err != nil {
		return nil, err
	}

	// Check if session ok
	ses, err := ctrl.importSession.FindByID(ctx, r.SessionID)
	if err != nil {
		return nil, err
	}

	ses.Fields = make(map[string]string)
	err = json.Unmarshal(r.Fields, &ses.Fields)
	if err != nil {
		return nil, err
	}

	ses.OnError = r.OnError

	// Errors are presented in the session
	err = func() (err error) {
		if ses.Progress.StartedAt != nil {
			return fmt.Errorf("unable to start import: import session already active")
		}

		sa := time.Now()
		ses.Progress.StartedAt = &sa

		// Prepare additional metadata
		tpl := resource.NewComposeRecordTemplate(
			strconv.FormatUint(ses.ModuleID, 10),
			strconv.FormatUint(ses.NamespaceID, 10),
			ses.Name,
			false,
			resource.MapToMappingTplSet(ses.Fields),
			ses.Key,
		)

		// Shape the data
		ses.Resources = append(ses.Resources, tpl)
		rt := resource.ComposeRecordShaper()
		ses.Resources, err = resource.Shape(ses.Resources, rt)

		// Build
		cfg := &estore.EncoderConfig{
			// For now the identifier is ignored, so this will never occur
			OnExisting: resource.Skip,
			DeferOk: func() {
				ses.Progress.Completed++
			},
		}
		cfg.DeferNok = func(err error) error {
			ses.Progress.Failed++

			if ses.Progress.FailLog == nil {
				ses.Progress.FailLog = &service.FailLog{
					Errors: make(service.ErrorIndex),
				}
			}

			if rve, is := err.(*types.RecordValueErrorSet); is {
				for _, ve := range rve.Set {
					for k, v := range ve.Meta {
						ses.Progress.FailLog.Errors.Add(fmt.Sprintf("%s %s %v", ve.Kind, k, v))
					}
				}
			} else {
				ses.Progress.FailLog.Errors.Add(err.Error())
			}

			if len(ses.Progress.FailLog.Records) < service.IMPORT_ERROR_MAX_INDEX_COUNT {
				// +1 because we indexed them with 1 before
				ses.Progress.FailLog.Records = append(ses.Progress.FailLog.Records, int(ses.Progress.Completed)+1)
			} else {
				ses.Progress.FailLog.RecordsTruncated = true
			}

			if ses.OnError == service.IMPORT_ON_ERROR_SKIP {
				return nil
			}
			return err
		}
		se := estore.NewStoreEncoder(service.DefaultStore, cfg)
		bld := envoy.NewBuilder(se)
		g, err := bld.Build(ctx, ses.Resources...)
		if err != nil {
			return err
		}

		// Encode
		err = envoy.Encode(ctx, g, se)
		now := time.Now()
		ses.Progress.FinishedAt = &now
		if err != nil {
			ses.Progress.FailReason = err.Error()
			return err
		}

		return
	}()

	return ses, ctrl.record.RecordImport(ctx, err)
}

func (ctrl *Record) ImportProgress(ctx context.Context, r *request.RecordImportProgress) (interface{}, error) {
	// Get session
	ses, err := ctrl.importSession.FindByID(ctx, r.SessionID)
	if err != nil {
		return nil, err
	}

	return ses, nil
}

func (ctrl *Record) Export(ctx context.Context, r *request.RecordExport) (interface{}, error) {
	var (
		err error

		filename = fmt.Sprintf("; filename=%s.%s", r.Filename, r.Ext)

		rf = &types.RecordFilter{
			Query:       r.Filter,
			NamespaceID: r.NamespaceID,
			ModuleID:    r.ModuleID,
		}
		f = estore.NewDecodeFilter().
			ComposeRecord(rf)

		contentType string
	)

	// Access control
	if _, err = ctrl.module.FindByID(ctx, r.NamespaceID, r.ModuleID); err != nil {
		return nil, err
	}

	if len(r.Fields) == 1 {
		r.Fields = strings.Split(r.Fields[0], ",")
	}

	return func(w http.ResponseWriter, req *http.Request) {
		if len(r.Fields) == 0 {
			http.Error(w, "no record value fields provided", http.StatusBadRequest)
		}

		fx := make(map[string]bool)
		for _, f := range r.Fields {
			fx[f] = true
		}

		sd := estore.Decoder()
		nn, err := sd.Decode(ctx, service.DefaultStore, f)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to fetch records: %s", err.Error()), http.StatusBadRequest)
		}

		var encoder envoy.PrepareEncodeStreamer

		switch strings.ToLower(r.Ext) {
		case "json", "jsonl", "ldjson", "ndjson":
			contentType = "application/jsonl"
			encoder = ejson.NewBulkRecordEncoder(&ejson.EncoderConfig{
				Fields:   fx,
				Timezone: r.Timezone,
			})

		case "csv":
			contentType = "text/csv"
			encoder = csv.NewBulkRecordEncoder(&csv.EncoderConfig{
				Fields:   fx,
				Timezone: r.Timezone,
			})

		default:
			http.Error(w, "unsupported format ("+r.Ext+")", http.StatusBadRequest)
			return
		}

		w.Header().Add("Content-Type", contentType)
		w.Header().Add("Content-Disposition", "attachment"+filename)

		bld := envoy.NewBuilder(encoder)
		g, err := bld.Build(ctx, nn...)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		err = envoy.Encode(ctx, g, encoder)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		ss := encoder.Stream()

		// Find only the stream we are interested in
		for _, s := range ss {
			if s.Resource == types.RecordResourceType {
				io.Copy(w, s.Source)
			}
		}

		err = ctrl.record.RecordExport(ctx, *rf)

	}, err
}

func (ctrl Record) Exec(ctx context.Context, r *request.RecordExec) (interface{}, error) {
	aa := request.ProcedureArgs(r.Args)

	switch r.Procedure {
	case "organize":
		return api.OK(), ctrl.record.Organize(ctx,
			r.NamespaceID,
			r.ModuleID,
			aa.GetUint64("recordID"),
			aa.Get("positionField"),
			aa.Get("position"),
			aa.Get("filter"),
			aa.Get("groupField"),
			aa.Get("group"),
		)
	default:
		return nil, fmt.Errorf("unknown procedure")
	}
}

func (ctrl *Record) TriggerScript(ctx context.Context, r *request.RecordTriggerScript) (interface{}, error) {
	module, record, err := ctrl.record.TriggerScript(ctx, r.NamespaceID, r.ModuleID, r.RecordID, r.Values, r.Script)

	// Script can return modified record and we'll pass it on to the caller
	return ctrl.makePayload(ctx, module, record, err)
}

func (ctrl *Record) TriggerScriptOnList(ctx context.Context, r *request.RecordTriggerScriptOnList) (rsp interface{}, err error) {
	//var (
	//	module    *types.Module
	//	namespace *types.Namespace
	//)
	//
	//if module, err = ctrl.module.FindByID(ctx, r.NamespaceID, r.ModuleID); err != nil {
	//	return
	//}
	//
	//if namespace, err = ctrl.namespace.With(ctx).FindByID(r.NamespaceID); err != nil {
	//	return
	//}

	// @todo this does not need to be under /record ... where then?!?!
	err = corredor.Service().ExecIterator(ctx, r.Script)

	// Script can return modified record and we'll pass it on to the caller
	return api.OK(), err
}

func (ctrl Record) makeBulkPayload(ctx context.Context, m *types.Module, err error, rr ...*types.Record) (*recordPayload, error) {
	if err != nil || rr == nil {
		return nil, err
	}

	return &recordPayload{
		Record:  rr[0],
		Records: rr[1:],

		CanUpdateRecord: ctrl.ac.CanUpdateRecord(ctx, rr[0]),
		CanReadRecord:   ctrl.ac.CanReadRecord(ctx, rr[0]),
		CanDeleteRecord: ctrl.ac.CanDeleteRecord(ctx, rr[0]),
	}, nil
}

func (ctrl Record) makePayload(ctx context.Context, m *types.Module, r *types.Record, err error) (*recordPayload, error) {
	if err != nil || r == nil {
		return nil, err
	}

	return &recordPayload{
		Record: r,

		CanGrant: ctrl.ac.CanGrant(ctx),

		CanUpdateRecord: ctrl.ac.CanUpdateRecord(ctx, r),
		CanReadRecord:   ctrl.ac.CanReadRecord(ctx, r),
		CanDeleteRecord: ctrl.ac.CanDeleteRecord(ctx, r),
	}, nil
}

func (ctrl Record) makeFilterPayload(ctx context.Context, m *types.Module, rr types.RecordSet, f *types.RecordFilter, err error) (*recordSetPayload, error) {
	if err != nil {
		return nil, err
	}

	modp := &recordSetPayload{Filter: f, Set: make([]*recordPayload, len(rr))}

	for i := range rr {
		modp.Set[i], _ = ctrl.makePayload(ctx, m, rr[i], nil)
	}

	return modp, nil
}

// Special care for record validation errors
//
// We need to return a bit different format of response
// with all details that were collected through validation
func (ctrl Record) handleValidationError(rve *types.RecordValueErrorSet) interface{} {
	return func(w http.ResponseWriter, _ *http.Request) {
		rval := struct {
			Error struct {
				Message string                   `json:"message"`
				Details []types.RecordValueError `json:"details,omitempty"`
			} `json:"error"`
		}{}

		rval.Error.Message = rve.Error()
		rval.Error.Details = rve.Set

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(rval)
	}
}
