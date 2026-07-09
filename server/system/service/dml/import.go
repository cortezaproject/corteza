package dml

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cast"

	composeTypes "github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/pkg/dal"
	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/pkg/id"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
)

const importBatchSize = 500

type (
	dalDataReader interface {
		SearchExternalModels(ctx context.Context, connectionID uint64) (dal.ModelSet, error)
		SearchExternalData(ctx context.Context, connectionID uint64, model *dal.Model, f filter.Filter) (dal.Iterator, error)
	}

	importAC interface {
		CanSearchDalConnections(ctx context.Context) bool
		CanCreateDalConnection(ctx context.Context) bool
	}

	ComposeRecordSvc interface {
		Bulk(ctx context.Context, skipFailed bool, oo ...*composeTypes.RecordBulkOperation) ([]composeTypes.RecordBulkOperationResult, error)
		Search(ctx context.Context, f composeTypes.RecordFilter) (composeTypes.RecordSet, composeTypes.RecordFilter, error)
	}

	Importer struct {
		mapping   mappingReader
		applier   *Applier
		dalSvc    dalDataReader
		recordSvc ComposeRecordSvc
		moduleSvc ComposeModuleSvc
		nsSvc     ComposeNamespaceSvc
		store     store.Storer
		ac        importAC
	}
)

func NewImporter(
	m mappingReader,
	a *Applier,
	mod ComposeModuleSvc,
	d dalDataReader,
	r ComposeRecordSvc,
	ns ComposeNamespaceSvc,
	s store.Storer,
	ac importAC,
) *Importer {
	return &Importer{
		mapping:   m,
		applier:   a,
		moduleSvc: mod,
		dalSvc:    d,
		recordSvc: r,
		nsSvc:     ns,
		store:     s,
		ac:        ac,
	}
}

func (im *Importer) GetRun(ctx context.Context, runID uint64) (*types.DmlImportRun, error) {
	if !im.ac.CanSearchDalConnections(ctx) {
		return nil, errors.Unauthorized("not allowed to access DML import runs")
	}
	run, err := store.LookupDmlImportRunByID(ctx, im.store, runID)
	if err != nil {
		return nil, err
	}
	if run == nil {
		return nil, fmt.Errorf("dml: import run %d not found", runID)
	}
	return run, nil
}

// RunImport starts an import for a single DmlMapping.
func (im *Importer) RunImport(ctx context.Context, mappingID uint64, method types.DmlImportMethod) (*types.DmlImportRun, error) {
	if !im.ac.CanCreateDalConnection(ctx) {
		return nil, errors.Unauthorized("not allowed to run DML imports")
	}
	switch method {
	case "":
		method = types.DmlImportMethodBackground
	case types.DmlImportMethodBackground:
	default:
		return nil, fmt.Errorf("unknown import method %q", method)
	}

	mp, err := im.mapping.FindByID(ctx, mappingID)
	if err != nil {
		return nil, err
	}

	run := &types.DmlImportRun{
		ID:           id.Next(),
		ConnectionID: mp.ConnectionID,
		MappingID:    mappingID,
		Method:       method,
		Status:       "pending",
	}
	if err := store.CreateDmlImportRun(ctx, im.store, run); err != nil {
		return nil, err
	}

	identity := auth.GetIdentityFromContext(ctx)
	switch method {
	case types.DmlImportMethodBackground:
		go im.runImportWork(auth.SetIdentityToContext(context.Background(), identity), run, mp)
	}

	return run, nil
}

// RunImportForeground runs the import synchronously in the calling goroutine.
// Used by the project publish flow where migration must complete before state flips.
func (im *Importer) RunImportForeground(ctx context.Context, mappingID uint64) error {
	mp, err := im.mapping.FindByID(ctx, mappingID)
	if err != nil {
		return err
	}
	run := &types.DmlImportRun{
		ID:           id.Next(),
		ConnectionID: mp.ConnectionID,
		MappingID:    mappingID,
		Method:       "foreground",
		Status:       "pending",
	}
	if err := store.CreateDmlImportRun(ctx, im.store, run); err != nil {
		return err
	}
	im.runImportWork(ctx, run, mp)
	if run.Status != "completed" {
		return fmt.Errorf("dml import failed: %s", run.Error)
	}
	return nil
}

func (im *Importer) runImportWork(ctx context.Context, run *types.DmlImportRun, mp *types.DmlMapping) {
	run.Status = "running"
	if err := store.UpdateDmlImportRun(ctx, im.store, run); err != nil {
		return
	}

	fail := func(format string, args ...any) {
		run.Status = "failed"
		run.Error = fmt.Sprintf(format, args...)
		_ = store.UpdateDmlImportRun(ctx, im.store, run)
	}

	if err := im.applier.Apply(ctx, mp.ID); err != nil {
		fail("apply: %s", err.Error())
		return
	}

	nsHandle := mp.NamespaceHandle
	if nsHandle == "" {
		nsHandle = fmt.Sprintf("dml_%d", mp.ConnectionID)
	}
	ns, err := im.nsSvc.FindByHandle(ctx, nsHandle)
	if err != nil {
		fail("namespace %q not found — run Apply first", nsHandle)
		return
	}

	externalModels, err := im.dalSvc.SearchExternalModels(ctx, mp.ConnectionID)
	if err != nil {
		fail("%s", err.Error())
		return
	}
	var srcModel *dal.Model
	for _, m := range externalModels {
		if m.Ident == mp.SourceIdent {
			srcModel = m
			break
		}
	}
	if srcModel == nil {
		fail("source table %q not found on connection", mp.SourceIdent)
		return
	}

	module, err := im.moduleSvc.FindByHandle(ctx, ns.ID, mp.ModuleHandle)
	if err != nil {
		fail("target module %q not found in namespace %q — run Apply first", mp.ModuleHandle, nsHandle)
		return
	}

	if err := im.importTable(ctx, run, ns, module, mp, srcModel); err != nil {
		fail("%s", err.Error())
		return
	}

	run.Status = "completed"
	_ = store.UpdateDmlImportRun(ctx, im.store, run)
}

func (im *Importer) importTable(
	ctx context.Context,
	run *types.DmlImportRun,
	ns *composeTypes.Namespace,
	module *composeTypes.Module,
	mp *types.DmlMapping,
	srcModel *dal.Model,
) error {
	iter, err := im.dalSvc.SearchExternalData(ctx, mp.ConnectionID, srcModel, filter.Generic())
	if err != nil {
		return fmt.Errorf("search external data %q: %w", mp.SourceIdent, err)
	}
	defer iter.Close()

	batch := make([]*composeTypes.RecordBulkOperation, 0, importBatchSize)

	flush := func() error {
		if len(batch) == 0 {
			return nil
		}
		results, err := im.recordSvc.Bulk(ctx, true, batch...)
		if err != nil {
			return err
		}
		for _, r := range results {
			if r.Error != nil {
				run.Failed++
			} else {
				run.Processed++
			}
		}
		batch = batch[:0]
		return nil
	}

	row := &importRow{}
	for iter.Next(ctx) {
		row.reset()
		if err := iter.Scan(row); err != nil {
			run.Failed++
			continue
		}

		var values composeTypes.RecordValueSet
		for _, col := range mp.Columns {
			if col.Skip {
				continue
			}
			val, err := row.get(col.SourceIdent)
			if err != nil || val == nil {
				continue
			}
			values = append(values, &composeTypes.RecordValue{
				Name:  col.FieldName,
				Value: rawToString(val),
			})
		}

		batch = append(batch, &composeTypes.RecordBulkOperation{
			Record: &composeTypes.Record{
				ModuleID:    module.ID,
				NamespaceID: ns.ID,
				Values:      values,
			},
			Operation: composeTypes.OperationTypeCreate,
		})

		if len(batch) >= importBatchSize {
			if err := flush(); err != nil {
				return err
			}
		}
	}

	if err := iter.Err(); err != nil {
		return err
	}

	return flush()
}

func rawToString(v any) string {
	switch t := v.(type) {
	case nil:
		return ""
	case string:
		return t
	case time.Time:
		return t.Format(time.RFC3339)
	case *time.Time:
		if t == nil {
			return ""
		}
		return t.Format(time.RFC3339)
	default:
		return cast.ToString(v)
	}
}

type importRow struct {
	vals map[string][]any
}

func (r *importRow) reset() {
	r.vals = make(map[string][]any)
}

func (r *importRow) CountValues() map[string]uint {
	out := make(map[string]uint, len(r.vals))
	for k, vv := range r.vals {
		out[k] = uint(len(vv))
	}
	return out
}

func (r *importRow) GetValue(name string, pos uint) (any, error) {
	vv, ok := r.vals[name]
	if !ok || int(pos) >= len(vv) {
		return nil, nil
	}
	return vv[pos], nil
}

func (r *importRow) SetValue(name string, pos uint, val any) error {
	for int(pos) >= len(r.vals[name]) {
		r.vals[name] = append(r.vals[name], nil)
	}
	r.vals[name][pos] = val
	return nil
}

func (r *importRow) get(name string) (any, error) {
	return r.GetValue(name, 0)
}
