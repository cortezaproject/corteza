package rest

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/titpetric/factory/resputil"

	"github.com/pkg/errors"

	"github.com/cortezaproject/corteza-server/compose/decoder"
	"github.com/cortezaproject/corteza-server/compose/encoder"
	"github.com/cortezaproject/corteza-server/compose/repository"
	"github.com/cortezaproject/corteza-server/compose/rest/request"
	"github.com/cortezaproject/corteza-server/compose/service"
	"github.com/cortezaproject/corteza-server/compose/service/event"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/corredor"
	"github.com/cortezaproject/corteza-server/pkg/mime"
	"github.com/cortezaproject/corteza-server/pkg/rh"
)

var _ = errors.Wrap

type (
	recordPayload struct {
		*types.Record

		CanUpdateRecord bool `json:"canUpdateRecord"`
		CanDeleteRecord bool `json:"canDeleteRecord"`
	}

	recordSetPayload struct {
		Filter types.RecordFilter `json:"filter"`
		Set    []*recordPayload   `json:"set"`
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
	}
}

func (ctrl *Record) Report(ctx context.Context, r *request.RecordReport) (interface{}, error) {
	return ctrl.record.With(ctx).Report(r.NamespaceID, r.ModuleID, r.Metrics, r.Dimensions, r.Filter)
}

func (ctrl *Record) List(ctx context.Context, r *request.RecordList) (interface{}, error) {
	var (
		m   *types.Module
		err error
	)

	if m, err = ctrl.module.With(ctx).FindByID(r.NamespaceID, r.ModuleID); err != nil {
		return nil, err
	}

	rr, filter, err := ctrl.record.With(ctx).Find(types.RecordFilter{
		NamespaceID: r.NamespaceID,
		ModuleID:    r.ModuleID,
		Filter:      r.Filter,
		Sort:        r.Sort,

		PageFilter: rh.Paging(r.Page, r.PerPage),
	})

	return ctrl.makeFilterPayload(ctx, m, rr, filter, err)
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

	record, err := ctrl.record.With(ctx).Create(&types.Record{
		NamespaceID: r.NamespaceID,
		ModuleID:    r.ModuleID,
		Values:      r.Values,
	})

	return ctrl.makePayload(ctx, m, record, err)
}

func (ctrl *Record) Update(ctx context.Context, r *request.RecordUpdate) (interface{}, error) {
	var (
		m   *types.Module
		err error
	)

	if m, err = ctrl.module.With(ctx).FindByID(r.NamespaceID, r.ModuleID); err != nil {
		return nil, err
	}

	record, err := ctrl.record.With(ctx).Update(&types.Record{
		ID:          r.RecordID,
		NamespaceID: r.NamespaceID,
		ModuleID:    r.ModuleID,
		Values:      r.Values,
	})

	return ctrl.makePayload(ctx, m, record, err)
}

func (ctrl *Record) Delete(ctx context.Context, r *request.RecordDelete) (interface{}, error) {
	return resputil.OK(), ctrl.record.With(ctx).DeleteByID(r.NamespaceID, r.RecordID)
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
		}
	}

	// determine decoder
	switch strings.ToLower(ext) {
	case "json", "jsonl", "ldjson", "ndjson":
		recordDecoder = decoder.NewStructuredDecoder(json.NewDecoder(f), f)

	case "csv":
		recordDecoder = decoder.NewFlatReader(csv.NewReader(f), f)

	default:
		return nil, service.ErrRecordImportFormatNotSupported

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

	return ctrl.importSession.SetRecordByID(
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
	ses, err := ctrl.importSession.FindRecordByID(ctx, r.SessionID)
	if err != nil {
		return nil, err
	}

	if ses.Progress.StartedAt != nil {
		return nil, service.ErrRecordImportSessionAlreadyStarted
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
	ses, err := ctrl.importSession.FindRecordByID(ctx, r.SessionID)
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
			Filter:      r.Filter,
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

		switch strings.ToLower(r.Ext) {
		case "json", "jsonl", "ldjson", "ndjson":
			contentType = "application/jsonl"
			recordEncoder = encoder.NewStructuredEncoder(json.NewEncoder(w), ff...)

		case "csv":
			contentType = "text/csv"
			recordEncoder = encoder.NewFlatWriter(csv.NewWriter(w), true, ff...)

		case "xlsx":
			contentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
			recordEncoder = encoder.NewExcelizeEncoder(w, true, ff...)

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
		return nil, errors.New("unknown procedure")
	}

	return nil, nil
}

func (ctrl *Record) Trigger(ctx context.Context, r *request.RecordTrigger) (rsp interface{}, err error) {
	var (
		record    *types.Record
		module    *types.Module
		namespace *types.Namespace
	)

	if record, err = ctrl.record.With(ctx).FindByID(r.NamespaceID, r.RecordID); err != nil {
		return
	}

	if module, err = ctrl.module.With(ctx).FindByID(r.NamespaceID, r.ModuleID); err != nil {
		return
	}

	if namespace, err = ctrl.namespace.With(ctx).FindByID(r.NamespaceID); err != nil {
		return
	}

	return resputil.OK(), corredor.Service().ExecOnManual(ctx, r.Script, event.RecordOnManual(record, nil, module, namespace))
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

func (ctrl Record) makeFilterPayload(ctx context.Context, m *types.Module, rr types.RecordSet, f types.RecordFilter, err error) (*recordSetPayload, error) {
	if err != nil {
		return nil, err
	}

	modp := &recordSetPayload{Filter: f, Set: make([]*recordPayload, len(rr))}

	for i := range rr {
		modp.Set[i], _ = ctrl.makePayload(ctx, m, rr[i], nil)
	}

	return modp, nil
}
