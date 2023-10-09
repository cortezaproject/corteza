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

	composeEnvoy "github.com/cortezaproject/corteza/server/compose/envoy"
	"github.com/cortezaproject/corteza/server/compose/rest/request"
	"github.com/cortezaproject/corteza/server/compose/service"
	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/api"
	"github.com/cortezaproject/corteza/server/pkg/corredor"
	"github.com/cortezaproject/corteza/server/pkg/dal"
	"github.com/cortezaproject/corteza/server/pkg/envoy"
	"github.com/cortezaproject/corteza/server/pkg/envoy/csv"
	envoyJson "github.com/cortezaproject/corteza/server/pkg/envoy/json"
	estore "github.com/cortezaproject/corteza/server/pkg/envoy/store"
	"github.com/cortezaproject/corteza/server/pkg/envoyx"
	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/pkg/revisions"
	"github.com/cortezaproject/corteza/server/store"
)

type (
	recordBulkPatchRecord struct {
		Record      *types.Record              `json:"record"`
		Error       error                      `json:"error,omitempty"`
		ValueErrors *types.RecordValueErrorSet `json:"valueErrors,omitempty"`
	}
	recordBulkPatchPayload struct {
		Records []recordBulkPatchRecord `json:"records"`
	}

	recordPayload struct {
		*types.Record

		Records           types.RecordSet            `json:"records,omitempty"`
		RecordValueErrors *types.RecordValueErrorSet `json:"valueErrors"`

		CanManageOwnerOnRecord bool `json:"canManageOwnerOnRecord"`
		CanUpdateRecord        bool `json:"canUpdateRecord"`
		CanReadRecord          bool `json:"canReadRecord"`
		CanDeleteRecord        bool `json:"canDeleteRecord"`
		CanUndeleteRecord      bool `json:"canUndeleteRecord"`
		CanSearchRevisions     bool `json:"canSearchRevisions"`

		CanGrant bool `json:"canGrant"`
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
		CanUndeleteRecord(context.Context, *types.Record) bool
		CanManageOwnerOnRecord(context.Context, *types.Record) bool
		CanSearchRevisionsOnRecord(context.Context, *types.Record) bool
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
			Meta:        r.Meta,
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

	record, dd, err := ctrl.record.FindByID(ctx, r.NamespaceID, r.ModuleID, r.RecordID)

	// Temp workaround until we do proper by-module filtering for record findByID
	if record != nil && record.ModuleID != r.ModuleID {
		return nil, store.ErrNotFound
	}

	return ctrl.makePayload(ctx, m, record, dd, err)
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
			Meta:        r.Meta,
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

	results, err := ctrl.record.Bulk(ctx, false, oo...)
	if rve := types.IsRecordValueErrorSet(err); rve != nil {
		return ctrl.handleValidationError(rve), nil
	}

	var (
		rr types.RecordSet
		dd = &types.RecordValueErrorSet{}
	)

	for _, r := range results {
		rr = append(rr, r.Record)
		dd.Merge(r.DuplicationError)
	}

	return ctrl.makeBulkPayload(ctx, m, dd, err, rr...)
}

func (ctrl *Record) Patch(ctx context.Context, req *request.RecordPatch) (interface{}, error) {
	var (
		f = types.RecordFilter{
			Query:       req.Query,
			NamespaceID: req.NamespaceID,
			ModuleID:    req.ModuleID,
			Deleted:     filter.State(0),
		}

		err error
	)

	counters := make(map[string]uint)
	for _, v := range req.Values {
		v.Place = counters[v.Name]
		counters[v.Name]++
	}

	err = ctrl.record.BulkModifyByFilter(ctx, f, req.Values, types.OperationTypePatch)

	if rve := types.IsRecordValueErrorSet(err); rve != nil {
		return ctrl.handleValidationError(rve), nil
	}

	return api.OK(), err
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
			Meta:        r.Meta,
			OwnedBy:     r.OwnedBy,
			UpdatedAt:   r.UpdatedAt,
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

	results, err := ctrl.record.Bulk(ctx, false, oo...)
	if rve := types.IsRecordValueErrorSet(err); rve != nil {
		return ctrl.handleValidationError(rve), nil
	}

	var rr types.RecordSet
	dd := &types.RecordValueErrorSet{}

	for _, r := range results {
		rr = append(rr, r.Record)
		dd.Merge(r.DuplicationError)
	}

	return ctrl.makeBulkPayload(ctx, m, dd, err, rr...)
}

func (ctrl *Record) Delete(ctx context.Context, r *request.RecordDelete) (interface{}, error) {
	return api.OK(), ctrl.record.DeleteByID(ctx, r.NamespaceID, r.ModuleID, r.RecordID)
}

func (ctrl *Record) BulkDelete(ctx context.Context, r *request.RecordBulkDelete) (interface{}, error) {
	var (
		f = types.RecordFilter{
			Query:       r.Query,
			NamespaceID: r.NamespaceID,
			ModuleID:    r.ModuleID,
			Deleted:     filter.State(0),
		}
	)

	if r.Truncate {
		return nil, fmt.Errorf("pending implementation")
	}

	return api.OK(), ctrl.record.BulkModifyByFilter(ctx, f, nil, types.OperationTypeDelete)
}

func (ctrl *Record) Undelete(ctx context.Context, r *request.RecordUndelete) (interface{}, error) {
	return api.OK(), ctrl.record.UndeleteByID(ctx, r.NamespaceID, r.ModuleID, r.RecordID)
}

func (ctrl *Record) BulkUndelete(ctx context.Context, r *request.RecordBulkUndelete) (interface{}, error) {
	var (
		f = types.RecordFilter{
			Query:       r.Query,
			NamespaceID: r.NamespaceID,
			ModuleID:    r.ModuleID,
			Deleted:     filter.State(1),
		}
	)

	return api.OK(), ctrl.record.BulkModifyByFilter(ctx, f, nil, types.OperationTypeUndelete)
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

// @todo :')
func (ctrl *Record) ImportRun(ctx context.Context, r *request.RecordImportRun) (_ interface{}, err error) {
	var (
		ns  *types.Namespace
		mod *types.Module
	)
	if mod, err = ctrl.module.FindByID(ctx, r.NamespaceID, r.ModuleID); err != nil {
		return nil, err
	}
	if ns, err = ctrl.namespace.FindByID(ctx, r.NamespaceID); err != nil {
		return nil, err
	}

	var importSession *service.RecordImportSession
	err = func() (err error) {
		// Check if session ok
		{
			importSession, err = ctrl.importSession.FindByID(ctx, r.SessionID)
			if err != nil {
				return
			}

			if importSession.Progress.StartedAt != nil {
				return fmt.Errorf("unable to start import: import session already active")
			}
		}

		sa := time.Now()
		importSession.Progress.StartedAt = &sa

		// Some prereq
		{
			importSession.Fields = make(map[string]string)
			err = json.Unmarshal(r.Fields, &importSession.Fields)
			if err != nil {
				return err
			}

			importSession.OnError = r.OnError
		}

		// Prep envoy bits
		var (
			envoySvc     *envoyx.Service
			encodeParams envoyx.EncodeParams

			nodeScope    envoyx.Scope
			nodes        envoyx.NodeSet
			node         *envoyx.Node
			storeEncoder = composeEnvoy.StoreEncoder{}
		)
		{
			envoySvc = envoyx.New()
			// @todo add when/if needed
			envoySvc.AddEncoder(envoyx.EncodeTypeStore,
				storeEncoder,
			)

			encodeParams = envoyx.EncodeParams{
				Type: envoyx.EncodeTypeStore,
				Params: map[string]any{
					"storer": service.DefaultStore,
					"dal":    dal.Service(),
				},
				DeferOk: func() {
					importSession.Progress.Completed++
				},
				DeferNok: func(err error) error {
					importSession.Progress.Failed++

					if importSession.Progress.FailLog == nil {
						importSession.Progress.FailLog = &service.FailLog{
							Errors: make(service.ErrorIndex),
						}
					}

					if rve, is := err.(*types.RecordValueErrorSet); is {
						for _, ve := range rve.Set {
							for k, v := range ve.Meta {
								importSession.Progress.FailLog.Errors.Add(fmt.Sprintf("%s %s %v", ve.Kind, k, v))
							}
						}
					} else {
						importSession.Progress.FailLog.Errors.Add(err.Error())
					}

					if len(importSession.Progress.FailLog.Records) < service.IMPORT_ERROR_MAX_INDEX_COUNT {
						// +1 because we indexed them with 1 before
						importSession.Progress.FailLog.Records = append(importSession.Progress.FailLog.Records, int(importSession.Progress.Completed)+1)
					} else {
						importSession.Progress.FailLog.RecordsTruncated = true
					}

					if importSession.OnError == service.IMPORT_ON_ERROR_SKIP {
						return nil
					}
					return err
				},
			}

			nodeScope = envoyx.Scope{
				ResourceType: types.NamespaceResourceType,
				Identifiers:  envoyx.MakeIdentifiers(importSession.NamespaceID),
			}

			fieldMapping := map[string]envoyx.MapEntry{}
			for c, f := range importSession.Fields {
				fieldMapping[c] = envoyx.MapEntry{
					Column: c,
					Field:  f,
				}
			}

			nsNode := &envoyx.Node{
				Resource:     ns,
				ResourceType: types.NamespaceResourceType,
				Identifiers:  envoyx.MakeIdentifiers(ns.Slug, ns.ID),
				Scope:        nodeScope,
				Placeholder:  true,
			}

			modNode := &envoyx.Node{
				Resource:     mod,
				ResourceType: types.ModuleResourceType,
				Identifiers:  envoyx.MakeIdentifiers(mod.Handle, mod.ID),
				Scope:        nodeScope,
				References: map[string]envoyx.Ref{
					"NamespaceID": {
						ResourceType: types.NamespaceResourceType,
						Identifiers:  nsNode.Identifiers,
						Scope:        nsNode.Scope,
					},
				},
				Placeholder: true,
			}

			node = &envoyx.Node{
				Datasource: &composeEnvoy.RecordDatasource{
					Mapping: envoyx.DatasourceMapping{
						SourceIdent: importSession.Name,
						References:  map[string]string{},
						Scope:       map[string]string{},
						Defaultable: false,
						Mapping: envoyx.FieldMapping{
							Map: fieldMapping,
						},
					},
				},
				ResourceType: composeEnvoy.ComposeRecordDatasourceAuxType,
				Identifiers:  envoyx.MakeIdentifiers(importSession.Name),

				References: map[string]envoyx.Ref{
					"ModuleID": {
						ResourceType: types.ModuleResourceType,
						Identifiers:  envoyx.MakeIdentifiers(importSession.ModuleID),
						Scope:        nodeScope,
					},
					"NamespaceID": {
						ResourceType: types.NamespaceResourceType,
						Identifiers:  envoyx.MakeIdentifiers(importSession.NamespaceID),
						Scope:        nodeScope,
					},
				},
				Scope: nodeScope,
			}

			nodes = envoyx.NodeSet{nsNode, modNode, node}
		}

		// encoding stuff
		var (
			depGraph *envoyx.DepGraph
		)
		{
			depGraph, err = envoySvc.Bake(ctx, encodeParams, importSession.Providers, nodes...)
			if err != nil {
				return
			}

			{
				// @todo this is temporary because the service's logic is a bit flawed for this case
				err = storeEncoder.Prepare(ctx, encodeParams, composeEnvoy.ComposeRecordDatasourceAuxType, envoyx.NodeSet{node})
				if err != nil {
					return
				}

				err = storeEncoder.Encode(ctx, encodeParams, composeEnvoy.ComposeRecordDatasourceAuxType, envoyx.NodeSet{node}, depGraph)
				// @note err is handled lower down; bare with
			}

			// err = envoySvc.Encode(ctx, encodeParams, depGraph)
			now := time.Now()
			importSession.Progress.FinishedAt = &now
			if err != nil {
				importSession.Progress.FailReason = err.Error()
				return
			}
			return
		}
	}()
	return importSession, ctrl.record.RecordImport(ctx, err)
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
		nn, err := sd.Decode(ctx, service.DefaultStore, dal.Service(), f)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to fetch records: %s", err.Error()), http.StatusBadRequest)
		}

		var encoder envoy.PrepareEncodeStreamer

		switch strings.ToLower(r.Ext) {
		case "json", "jsonl", "ldjson", "ndjson":
			contentType = "application/jsonl"
			encoder = envoyJson.NewBulkRecordEncoder(&envoyJson.EncoderConfig{
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
	return ctrl.makePayload(ctx, module, record, nil, err)
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

func (ctrl *Record) Revisions(ctx context.Context, r *request.RecordRevisions) (interface{}, error) {
	var (
		makeRev = func() dal.ValueSetter { return &revisions.Revision{} }
	)

	iter, err := ctrl.record.SearchRevisions(ctx, r.NamespaceID, r.ModuleID, r.RecordID)
	if err != nil {
		return nil, err
	}

	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if _, err = w.Write([]byte(`{"response":{"set":[`)); err != nil {
			return
		}

		err = dal.IteratorEncodeJSON(ctx, w, iter, makeRev)
		if err != nil {
			return
		}

		if _, err = w.Write([]byte(`]}}`)); err != nil {
			return
		}

		return
	}, err
}

func (ctrl Record) makeBulkPayload(ctx context.Context, m *types.Module, dd *types.RecordValueErrorSet, err error, rr ...*types.Record) (*recordPayload, error) {
	if err != nil || rr == nil {
		return nil, err
	}

	return &recordPayload{
		Record:            rr[0],
		Records:           rr[1:],
		RecordValueErrors: dd,

		CanManageOwnerOnRecord: ctrl.ac.CanManageOwnerOnRecord(ctx, rr[0]),
		CanUpdateRecord:        ctrl.ac.CanUpdateRecord(ctx, rr[0]),
		CanReadRecord:          ctrl.ac.CanReadRecord(ctx, rr[0]),
		CanDeleteRecord:        ctrl.ac.CanDeleteRecord(ctx, rr[0]),
		CanUndeleteRecord:      ctrl.ac.CanUndeleteRecord(ctx, rr[0]),
		CanSearchRevisions:     ctrl.ac.CanSearchRevisionsOnRecord(ctx, rr[0]),
	}, nil
}

func (ctrl Record) makeRecordBulkPatchPayload(ctx context.Context, rr []types.RecordBulkOperationResult, err error) (*recordBulkPatchPayload, error) {
	if err != nil {
		return nil, err
	}

	out := &recordBulkPatchPayload{
		Records: make([]recordBulkPatchRecord, 0, len(rr)),
	}

	for _, r := range rr {
		vr := r.ValueError
		vr.Merge(r.DuplicationError)
		out.Records = append(out.Records, recordBulkPatchRecord{
			Record:      r.Record,
			ValueErrors: vr,
			Error:       r.Error,
		})
	}

	return out, nil
}

func (ctrl Record) makePayload(ctx context.Context, m *types.Module, r *types.Record, dd *types.RecordValueErrorSet, err error) (*recordPayload, error) {
	if err != nil || r == nil {
		return nil, err
	}

	return &recordPayload{
		Record:            r,
		RecordValueErrors: dd,

		CanGrant: ctrl.ac.CanGrant(ctx),

		CanManageOwnerOnRecord: ctrl.ac.CanManageOwnerOnRecord(ctx, r),
		CanUpdateRecord:        ctrl.ac.CanUpdateRecord(ctx, r),
		CanReadRecord:          ctrl.ac.CanReadRecord(ctx, r),
		CanDeleteRecord:        ctrl.ac.CanDeleteRecord(ctx, r),
		CanUndeleteRecord:      ctrl.ac.CanUndeleteRecord(ctx, r),
		CanSearchRevisions:     ctrl.ac.CanSearchRevisionsOnRecord(ctx, r),
	}, nil
}

func (ctrl Record) makeFilterPayload(ctx context.Context, m *types.Module, rr types.RecordSet, f *types.RecordFilter, err error) (*recordSetPayload, error) {
	if err != nil {
		return nil, err
	}

	modp := &recordSetPayload{Filter: f, Set: make([]*recordPayload, len(rr))}

	for i := range rr {
		modp.Set[i], _ = ctrl.makePayload(ctx, m, rr[i], nil, nil)
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
