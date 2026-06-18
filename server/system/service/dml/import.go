package dml

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cast"

	composeTypes "github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/pkg/dal"
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

	// ComposeRecordSvc is the subset of compose record service that Importer needs.
	ComposeRecordSvc interface {
		Bulk(ctx context.Context, skipFailed bool, oo ...*composeTypes.RecordBulkOperation) ([]composeTypes.RecordBulkOperationResult, error)
		Search(ctx context.Context, f composeTypes.RecordFilter) (composeTypes.RecordSet, composeTypes.RecordFilter, error)
	}

	// Importer copies source rows from an external connection into compose records.
	Importer struct {
		mapping   *Mapping
		applier   *Applier
		dalSvc    dalDataReader
		recordSvc ComposeRecordSvc
		moduleSvc ComposeModuleSvc
		nsSvc     ComposeNamespaceSvc
		store     store.Storer
	}
)

func NewImporter(
	m *Mapping,
	a *Applier,
	mod ComposeModuleSvc,
	d dalDataReader,
	r ComposeRecordSvc,
	ns ComposeNamespaceSvc,
) *Importer {
	return &Importer{
		mapping:   m,
		applier:   a,
		moduleSvc: mod,
		dalSvc:    d,
		recordSvc: r,
		nsSvc:     ns,
		store:     m.store,
	}
}

// GetRun returns a stored import run by ID.
func (im *Importer) GetRun(ctx context.Context, runID uint64) (*types.DmlImportRun, error) {
	return loadRun(ctx, im.store, runID)
}

// RunImport starts a synchronous import for the given mapping.
// Returns an import run summary. Call GetRun(id) to poll status.
func (im *Importer) RunImport(ctx context.Context, mappingID uint64) (*types.DmlImportRun, error) {
	mp, err := im.mapping.FindByID(ctx, mappingID)
	if err != nil {
		return nil, err
	}

	run := &types.DmlImportRun{
		ID:           id.Next(),
		ConnectionID: mp.ConnectionID,
		MappingID:    mappingID,
		Status:       "running",
	}
	if err := saveRun(ctx, im.store, run); err != nil {
		return nil, err
	}

	finish := func() (*types.DmlImportRun, error) {
		return run, saveRun(ctx, im.store, run)
	}
	fail := func(format string, args ...any) (*types.DmlImportRun, error) {
		run.Status = "failed"
		run.Error = fmt.Sprintf(format, args...)
		return finish()
	}

	// resolve namespace (must already exist — Apply creates it)
	nsHandle := mp.NamespaceHandle
	if nsHandle == "" {
		nsHandle = fmt.Sprintf("dml_%d", mp.ConnectionID)
	}
	ns, err := im.nsSvc.FindByHandle(ctx, nsHandle)
	if err != nil {
		return fail("namespace %q not found — run Apply first", nsHandle)
	}

	// fetch all external models for the connection
	externalModels, err := im.dalSvc.SearchExternalModels(ctx, mp.ConnectionID)
	if err != nil {
		return fail("%s", err.Error())
	}
	modelByIdent := make(map[string]*dal.Model, len(externalModels))
	for _, m := range externalModels {
		modelByIdent[m.Ident] = m
	}

	for _, tbl := range mp.Tables {
		if tbl.Skip {
			continue
		}

		srcModel, ok := modelByIdent[tbl.SourceIdent]
		if !ok {
			// Mapping references a source table that no longer exists. Record
			// it as a run-level warning rather than inflating the per-row
			// Failed counter, then move on.
			if run.Error != "" {
				run.Error += "; "
			}
			run.Error += fmt.Sprintf("source table %q not found, skipped", tbl.SourceIdent)
			continue
		}

		// resolve the target module so records carry a valid ModuleID and the
		// compose create path can sanitize/format values per field kind.
		module, err := im.moduleSvc.FindByHandle(ctx, ns.ID, tbl.ModuleHandle)
		if err != nil {
			return fail("target module %q not found in namespace %q — run Apply first", tbl.ModuleHandle, nsHandle)
		}

		if err := im.importTable(ctx, run, ns, module, tbl, srcModel, mp.ConnectionID); err != nil {
			return fail("%s", err.Error())
		}

		// persist progress after each table
		if err := saveRun(ctx, im.store, run); err != nil {
			return nil, err
		}
	}

	run.Status = "completed"
	return finish()
}

func (im *Importer) importTable(
	ctx context.Context,
	run *types.DmlImportRun,
	ns *composeTypes.Namespace,
	module *composeTypes.Module,
	tbl *types.DmlTableMap,
	srcModel *dal.Model,
	connectionID uint64,
) error {
	iter, err := im.dalSvc.SearchExternalData(ctx, connectionID, srcModel, filter.Generic())
	if err != nil {
		return fmt.Errorf("search external data %q: %w", tbl.SourceIdent, err)
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

		// Build the record against the target module. ModuleID + NamespaceID
		// let the compose create path resolve the module and run the
		// sanitizer/formatter, which canonicalizes each value per field kind
		// and sets refs — so we only need to hand it the raw string form.
		var values composeTypes.RecordValueSet
		for _, col := range tbl.Columns {
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

// rawToString converts a scanned DAL value into the canonical string form
// compose records expect. The rdbms DAL driver already decodes column values
// to canonical strings, so the time.Time branches are defensive (other future
// drivers); the compose sanitizer/formatter, keyed on the record's ModuleID,
// re-encodes every value per the target field kind on create.
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

// importRow is a simple ValueGetter+ValueSetter for DAL iterator scanning.
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
