package rest

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"strconv"
	"strings"

	"github.com/titpetric/factory/resputil"

	"github.com/cortezaproject/corteza-server/compose/decoder"
	"github.com/cortezaproject/corteza-server/compose/encoder"
	"github.com/cortezaproject/corteza-server/compose/repository"
	"github.com/cortezaproject/corteza-server/compose/rest/request"
	"github.com/cortezaproject/corteza-server/compose/service"
	"github.com/cortezaproject/corteza-server/compose/service/event"
	"github.com/cortezaproject/corteza-server/compose/service/values"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/corredor"
	"github.com/cortezaproject/corteza-server/pkg/mime"
	"github.com/cortezaproject/corteza-server/pkg/payload"
	"github.com/cortezaproject/corteza-server/pkg/rh"
	systemService "github.com/cortezaproject/corteza-server/system/service"
	systemTypes "github.com/cortezaproject/corteza-server/system/types"
)

type (
	recordPayload struct {
		*types.Record

		Records types.RecordSet `json:"records,omitempty"`

		CanUpdateRecord bool `json:"canUpdateRecord"`
		CanDeleteRecord bool `json:"canDeleteRecord"`
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
		userFinder    systemService.UserService
	}

	recordAccessController interface {
		CanUpdateRecord(context.Context, *types.Module) bool
		CanDeleteRecord(context.Context, *types.Module) bool
	}
)

func (Record) New() *Record {
	return &Record{
		importSession: service.DefaultImportSession,
		record:        service.DefaultRecord,
		module:        service.DefaultModule,
		namespace:     service.DefaultNamespace,
		attachment:    service.DefaultAttachment,
		ac:            service.DefaultAccessControl,

		// See comment at DefaultSystemUser definition
		userFinder: service.DefaultSystemUser,
	}
}

func (ctrl *Record) Report(ctx context.Context, r *request.RecordReport) (interface{}, error) {
	return ctrl.record.With(ctx).Report(r.NamespaceID, r.ModuleID, r.Metrics, r.Dimensions, r.Filter)
}

func (ctrl *Record) List(ctx context.Context, r *request.RecordList) (interface{}, error) {
	var (
		m   *types.Module
		err error

		rf = types.RecordFilter{
			NamespaceID: r.NamespaceID,
			ModuleID:    r.ModuleID,
			Sort:        r.Sort,

			Deleted: rh.FilterState(r.Deleted),

			PageFilter: rh.Paging(r),
		}
	)

	if m, err = ctrl.module.With(ctx).FindByID(r.NamespaceID, r.ModuleID); err != nil {
		return nil, err
	}

	if r.Query != "" {
		// Query param takes preference
		rf.Query = r.Query
	} else if r.Filter != "" {
		// Backward compatibility
		// Filter param is deprecated
		rf.Query = r.Filter
	}

	rr, filter, err := ctrl.record.With(ctx).Find(rf)

	return ctrl.makeFilterPayload(ctx, m, rr, &filter, err)
}

func (ctrl *Record) Read(ctx context.Context, r *request.RecordRead) (interface{}, error) {
	var (
		m   *types.Module
		err error
	)

	if m, err = ctrl.module.With(ctx).FindByID(r.NamespaceID, r.ModuleID); err != nil {
		return nil, err
	}

	record, err := ctrl.record.With(ctx).FindByID(r.NamespaceID, r.RecordID)

	// Temp workaround until we do proper by-module filtering for record findByID
	if record != nil && record.ModuleID != r.ModuleID {
		return nil, repository.ErrRecordNotFound
	}

	return ctrl.makePayload(ctx, m, record, err)
}

func (ctrl *Record) Create(ctx context.Context, r *request.RecordCreate) (interface{}, error) {
	var (
		m   *types.Module
		err error
	)

	if m, err = ctrl.module.With(ctx).FindByID(r.NamespaceID, r.ModuleID); err != nil {
		return nil, err
	}

	oo := make([]*types.RecordBulkOperation, 0)

	// If defined, initialize parent record
	if r.Values != nil {
		rr := &types.Record{
			NamespaceID: r.NamespaceID,
			ModuleID:    r.ModuleID,
			Values:      r.Values,
		}
		oo = append(oo, &types.RecordBulkOperation{
			Record:    rr,
			Operation: types.OperationTypeCreate,
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

	rr, err := ctrl.record.With(ctx).Bulk(oo...)
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

	if m, err = ctrl.module.With(ctx).FindByID(r.NamespaceID, r.ModuleID); err != nil {
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

	rr, err := ctrl.record.With(ctx).Bulk(oo...)

	if rve := types.IsRecordValueErrorSet(err); rve != nil {
		return ctrl.handleValidationError(rve), nil
	}

	return ctrl.makeBulkPayload(ctx, m, err, rr...)
}

func (ctrl *Record) Delete(ctx context.Context, r *request.RecordDelete) (interface{}, error) {
	return resputil.OK(), ctrl.record.With(ctx).DeleteByID(r.NamespaceID, r.ModuleID, r.RecordID)
}

func (ctrl *Record) BulkDelete(ctx context.Context, r *request.RecordBulkDelete) (interface{}, error) {
	if r.Truncate {
		return nil, fmt.Errorf("pending implementation")
	}

	return resputil.OK(), ctrl.record.With(ctx).DeleteByID(
		r.NamespaceID,
		r.ModuleID,
		payload.ParseUInt64s(r.RecordIDs)...,
	)
}

func (ctrl *Record) Upload(ctx context.Context, r *request.RecordUpload) (interface{}, error) {
	file, err := r.Upload.Open()
	if err != nil {
		return nil, err
	}

	defer file.Close()

	a, err := ctrl.attachment.With(ctx).CreateRecordAttachment(
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
	var (
		err           error
		recordDecoder service.Decoder
		entryCount    uint64
	)

	// Access control.
	if _, err = ctrl.module.With(ctx).FindByID(r.NamespaceID, r.ModuleID); err != nil {
		return nil, err
	}

	f, err := r.Upload.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()

	_, ext, err := mime.Type(f)
	if err != nil {
		return nil, err
	}

	if ext == "txt" {
		if is, err := mime.JsonL(f); err != nil {
			return nil, err
		} else if is {
			ext = "jsonl"
		} else {
			// As last resort, use extension of the upload filename
			ext = strings.TrimLeft(path.Ext(r.Upload.Filename), ".")
		}
	}

	// determine decoder
	switch strings.ToLower(ext) {
	case "json", "jsonl", "ldjson", "ndjson":
		recordDecoder = decoder.NewStructuredDecoder(json.NewDecoder(f), f)

	case "csv":
		recordDecoder = decoder.NewFlatReader(csv.NewReader(f), f)

	default:
		// copied here from service/errors.go for backward compatibility
		// @todo move this logic to service and use action/error pattern
		return nil, fmt.Errorf("compose.service.RecordImportFormatNotSupported")
	}

	entryCount, err = recordDecoder.EntryCount()
	if err != nil {
		return nil, err
	}

	header := recordDecoder.Header()
	hh := make(map[string]string)
	for _, h := range header {
		hh[h] = ""
	}

	return ctrl.importSession.SetByID(
		ctx,
		0,
		r.NamespaceID,
		r.ModuleID,
		hh,
		&service.RecordImportProgress{EntryCount: entryCount},
		recordDecoder)
}

func (ctrl *Record) ImportRun(ctx context.Context, r *request.RecordImportRun) (interface{}, error) {
	var (
		err error
	)

	// Access control.
	if _, err = ctrl.module.With(ctx).FindByID(r.NamespaceID, r.ModuleID); err != nil {
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

	// @todo routine
	ctrl.record.With(ctx).Import(ses, ctrl.importSession)

	return ses, nil
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
	type (
		// ad-hoc interface for our encoder
		Encoder interface {
			service.Encoder
			Flush()
		}
	)

	var (
		err error

		// Record encoder
		recordEncoder Encoder

		filename = fmt.Sprintf("; filename=%s.%s", r.Filename, r.Ext)

		f = types.RecordFilter{
			NamespaceID: r.NamespaceID,
			ModuleID:    r.ModuleID,
			Query:       r.Filter,
		}

		contentType string
	)
	// Access control.
	if _, err = ctrl.module.With(ctx).FindByID(r.NamespaceID, r.ModuleID); err != nil {
		return nil, err
	}

	if len(r.Fields) == 1 {
		r.Fields = strings.Split(r.Fields[0], ",")
	}

	return func(w http.ResponseWriter, req *http.Request) {
		ff := encoder.MakeFields(r.Fields...)

		if len(ff) == 0 {
			http.Error(w, "no record value fields provided", http.StatusBadRequest)
		}

		// Custom user getter function for the underlying encoders.
		//
		// not the most optimal solution; we have no other means to do a proper preload of users
		// @todo preload users
		users := map[uint64]*systemTypes.User{}

		uf := func(ID uint64) (*systemTypes.User, error) {
			var err error

			if _, exists := users[ID]; exists {
				// nonexistent users are also cached!
				return users[ID], nil
			}

			// @todo this "communication" between system and compose
			//       services is ad-hoc solution
			users[ID], err = ctrl.userFinder.With(ctx).FindByID(ID)
			if err != nil {
				return nil, err
			}
			return users[ID], nil
		}

		switch strings.ToLower(r.Ext) {
		case "json", "jsonl", "ldjson", "ndjson":
			contentType = "application/jsonl"
			recordEncoder = encoder.NewStructuredEncoder(json.NewEncoder(w), uf, r.Timezone, ff...)

		case "csv":
			contentType = "text/csv"
			recordEncoder = encoder.NewFlatWriter(csv.NewWriter(w), true, uf, r.Timezone, ff...)

		case "xlsx":
			contentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
			recordEncoder = encoder.NewExcelizeEncoder(w, true, uf, r.Timezone, ff...)

		default:
			http.Error(w, "unsupported format ("+r.Ext+")", http.StatusBadRequest)
			return
		}

		w.Header().Add("Content-Type", contentType)
		w.Header().Add("Content-Disposition", "attachment"+filename)

		if err = ctrl.record.With(ctx).Export(f, recordEncoder); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		recordEncoder.Flush()
	}, nil
}

func (ctrl Record) Exec(ctx context.Context, r *request.RecordExec) (interface{}, error) {
	aa := request.ProcedureArgs(r.Args)

	switch r.Procedure {
	case "organize":
		return resputil.OK(), ctrl.record.With(ctx).Organize(
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

	return nil, nil
}

func (ctrl *Record) TriggerScript(ctx context.Context, r *request.RecordTriggerScript) (rsp interface{}, err error) {
	var (
		record    *types.Record
		oldRecord *types.Record
		module    *types.Module
		namespace *types.Namespace
	)

	if oldRecord, err = ctrl.record.With(ctx).FindByID(r.NamespaceID, r.RecordID); err != nil {
		return
	}

	if module, err = ctrl.module.With(ctx).FindByID(r.NamespaceID, r.ModuleID); err != nil {
		return
	}

	if namespace, err = ctrl.namespace.With(ctx).FindByID(r.NamespaceID); err != nil {
		return
	}

	record = oldRecord
	record.Values = values.Sanitizer().Run(module, r.Values)
	validated := values.Validator().Run(module, record)

	err = corredor.Service().Exec(
		ctx,
		r.Script,
		event.RecordOnManual(record, oldRecord, module, namespace, validated),
	)

	// Script can return modified record and we'll pass it on to the caller
	return ctrl.makePayload(ctx, module, record, err)
}

func (ctrl *Record) TriggerScriptOnList(ctx context.Context, r *request.RecordTriggerScriptOnList) (rsp interface{}, err error) {
	//var (
	//	module    *types.Module
	//	namespace *types.Namespace
	//)
	//
	//if module, err = ctrl.module.With(ctx).FindByID(r.NamespaceID, r.ModuleID); err != nil {
	//	return
	//}
	//
	//if namespace, err = ctrl.namespace.With(ctx).FindByID(r.NamespaceID); err != nil {
	//	return
	//}

	// @todo this does not need to be under /record ... where then?!?!
	err = corredor.Service().ExecIterator(ctx, r.Script)

	// Script can return modified record and we'll pass it on to the caller
	return resputil.OK(), err
}

func (ctrl Record) makeBulkPayload(ctx context.Context, m *types.Module, err error, rr ...*types.Record) (*recordPayload, error) {
	if err != nil || rr == nil {
		return nil, err
	}

	return &recordPayload{
		Record:  rr[0],
		Records: rr[1:],

		CanUpdateRecord: ctrl.ac.CanUpdateRecord(ctx, m),
		CanDeleteRecord: ctrl.ac.CanDeleteRecord(ctx, m),
	}, nil
}

func (ctrl Record) makePayload(ctx context.Context, m *types.Module, r *types.Record, err error) (*recordPayload, error) {
	if err != nil || r == nil {
		return nil, err
	}

	return &recordPayload{
		Record: r,

		CanUpdateRecord: ctrl.ac.CanUpdateRecord(ctx, m),
		CanDeleteRecord: ctrl.ac.CanDeleteRecord(ctx, m),
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
