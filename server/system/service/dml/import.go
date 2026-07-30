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
		go im.runImportWork(auth.SetIdentityToContext(context.Background(), identity), run, mp, true)
	}

	return run, nil
}

// RunImportForeground runs the import synchronously in the calling goroutine.
// Used by the project publish flow where migration must complete before state
// flips.
//
// It deliberately does NOT run the applier. Apply materialises a module FROM
// the mapping -- that is what an external import needs, because the target
// module does not exist yet or is only a projection of the source table. A
// publish is the opposite situation: the target module was authored by hand in
// the draft revision and is the authority; the mapping only says which column
// feeds which field.
//
// Letting Apply run here destroyed exactly that. It replaces Fields wholesale
// from the mapping's columns, and updateModuleFields matches old against new by
// field ID -- the mapping carries none -- so every field of every migrated
// module was soft-deleted and recreated as untyped text with its label reset to
// its name, and the module's Name was replaced by its handle. Verified on a
// live publish that used the server's own suggested mappings.
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
	im.runImportWork(ctx, run, mp, false)
	if run.Status != "completed" {
		return fmt.Errorf("dml import failed: %s", run.Error)
	}

	// A record the importer could not write is data that would be missing from
	// the published revision, so it fails the publish rather than the run.
	// importTable writes with skipFailed, which is right for an external import
	// (one malformed row should not stop ten thousand good ones) and wrong
	// here: "completed" was reported even when EVERY row failed, and publish
	// went on to swap the namespaces and soft-delete the source. The commonest
	// cause is a record-link value -- it holds a record id from the old
	// namespace, which does not exist in the new one, so the whole row is
	// rejected and a module with a link field silently arrives empty.
	//
	// Failing here is safe: migration runs before the status flip, so the old
	// namespace is still live and still holds every record.
	if run.Failed > 0 {
		return fmt.Errorf(
			"%d of %d records could not be migrated for module %q; "+
				"nothing was published. Record-link values are the usual cause: they still "+
				"name records in the previous revision. Drop that field from the module's "+
				"mapping to migrate the rest",
			run.Failed, run.Failed+run.Processed, mp.ModuleHandle,
		)
	}

	return nil
}

// applySchema decides whether the run may reshape the target module from the
// mapping before importing into it. True for imports of external data, false
// for a publish -- see RunImportForeground.
func (im *Importer) runImportWork(ctx context.Context, run *types.DmlImportRun, mp *types.DmlMapping, applySchema bool) {
	run.Status = "running"
	if err := store.UpdateDmlImportRun(ctx, im.store, run); err != nil {
		return
	}

	fail := func(format string, args ...any) {
		run.Status = "failed"
		run.Error = fmt.Sprintf(format, args...)
		_ = store.UpdateDmlImportRun(ctx, im.store, run)
	}

	if applySchema {
		if err := im.applier.Apply(ctx, mp.ID); err != nil {
			fail("apply: %s", err.Error())
			return
		}
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
			// A rejected record does not always come back as r.Error. Bulk
			// splits its failures: an operation-level error lands in Error,
			// but a VALUE-level rejection -- which is what an invalid
			// record-link produces, and record links are exactly what breaks
			// when records move to a new namespace -- lands in ValueError and
			// leaves Error nil (compose/service/record.go, the
			// IsRecordValueErrorSet branch). Counting only Error therefore
			// reported a clean run for records that were never written, which
			// is how a publish could report success and land an empty module.
			switch {
			case r.Error != nil, !r.ValueError.IsValid(), r.DuplicationError.HasStrictErrors():
				run.Failed++
			default:
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

		// One value per place, not just the first. importRow keeps every value
		// the source produced (the iterator sets them by place), but this loop
		// used to read position 0 alone -- so a multi-value field arrived with
		// its first entry and the rest were dropped without a word. On a
		// project publish that is silent data loss inside a record that
		// otherwise migrated fine.
		var values composeTypes.RecordValueSet
		counts := row.CountValues()
		for _, col := range mp.Columns {
			if col.Skip {
				continue
			}
			for place := uint(0); place < counts[col.SourceIdent]; place++ {
				val, err := row.GetValue(col.SourceIdent, place)
				if err != nil || val == nil {
					continue
				}
				values = append(values, &composeTypes.RecordValue{
					Name:  col.FieldName,
					Value: rawToString(val),
					Place: place,
				})
			}
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
