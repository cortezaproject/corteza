package dml

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cast"

	"github.com/crusttech/human/server/compose/dalutils"
	composeTypes "github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/pkg/id"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
)

// migratedMetaKey marks a record as one this importer wrote during a publish.
//
// It is what makes a retry converge. Migration runs outside the transaction
// that flips statuses and swaps namespaces -- it spans two connections and can
// take minutes, so it cannot sit inside a transaction -- which means a failure
// at the flip leaves a fully populated target namespace behind and the natural
// response, publishing again, used to copy every record a second time.
//
// Clearing the target module wholesale would be simpler, and a draft's modules
// are indeed populated only by migration today. But a draft namespace is a real
// namespace: nothing stops someone opening it and typing a record in before
// publishing, and that record is not ours to delete. The stamp keeps the
// distinction, and carries the source id so a migrated row can say where it
// came from.
const migratedMetaKey = "publishMigration"

type (
	// Migration is one publish's worth of record copying: several module
	// imports that have to agree with each other.
	//
	// They have to agree because of record links. A Record-kind value holds a
	// record id from the PREVIOUS namespace, which does not exist in the new
	// one, so writing it as-is has the whole row rejected -- every row of every
	// linked module, on every publish. Remapping needs the id the target got,
	// and the target may be imported after the module pointing at it (or by a
	// cycle, never "after" at all), so links are held back and written in a
	// second pass once every module has been through.
	Migration struct {
		im *Importer

		// idMap maps a source record id to the id the same record was given in
		// the target namespace.
		idMap map[uint64]uint64

		// modules in import order; only those a mapping actually imported into
		modules []*migratedModule
	}

	migratedModule struct {
		ns      *composeTypes.Namespace
		module  *composeTypes.Module
		records []*migratedRecord
	}

	// migratedRecord is the part of a source record that survives the create
	// but cannot be expressed in it: the links (deferred) and the provenance
	// (which procCreate overwrites by design).
	migratedRecord struct {
		srcID uint64
		newID uint64

		ownedBy   uint64
		createdAt time.Time
		createdBy uint64
		updatedAt *time.Time
		updatedBy uint64

		// Record-kind values, still holding source record ids
		links composeTypes.RecordValueSet
	}
)

// NewMigration starts a publish migration. Run it once per module mapping, then
// Finalize once.
func (im *Importer) NewMigration() *Migration {
	return &Migration{
		im:    im,
		idMap: make(map[uint64]uint64),
	}
}

// Run imports a single module mapping as part of this migration.
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
func (m *Migration) Run(ctx context.Context, mappingID uint64) error {
	mp, err := m.im.mapping.FindByID(ctx, mappingID)
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
	if err := store.CreateDmlImportRun(ctx, m.im.store, run); err != nil {
		return err
	}
	m.im.runImportWork(ctx, run, mp, false, m)
	if run.Status != "completed" {
		return fmt.Errorf("dml import failed: %s", run.Error)
	}

	// A record the importer could not write is data that would be missing from
	// the published revision, so it fails the publish rather than the run.
	// importTable writes with skipFailed, which is right for an external import
	// (one malformed row should not stop ten thousand good ones) and wrong
	// here: "completed" was reported even when EVERY row failed, and publish
	// went on to swap the namespaces and soft-delete the source.
	//
	// Failing here is safe: migration runs before the status flip, so the old
	// namespace is still live and still holds every record.
	if run.Failed > 0 {
		return fmt.Errorf(
			"%d of %d records could not be migrated for module %q; nothing was "+
				"published. Drop the offending field from the module's mapping to "+
				"migrate the rest",
			run.Failed, run.Failed+run.Processed, mp.ModuleHandle,
		)
	}

	return nil
}

// clear removes rows a previous attempt at this same migration left behind, so
// running the import again converges instead of doubling. See migratedMetaKey.
func (m *Migration) clear(ctx context.Context, ns *composeTypes.Namespace, module *composeTypes.Module) error {
	set, _, err := dalutils.ComposeRecordsList(ctx, m.im.dalSvc, module, composeTypes.RecordFilter{
		ModuleID:    module.ID,
		NamespaceID: ns.ID,
	})
	if err != nil {
		return fmt.Errorf("list previously migrated records of module %q: %w", module.Handle, err)
	}

	stale := make(composeTypes.RecordSet, 0, len(set))
	for _, r := range set {
		if isMigrated(r) {
			stale = append(stale, r)
		}
	}
	if len(stale) == 0 {
		return nil
	}

	if err = dalutils.ComposeRecordDelete(ctx, m.im.dalSvc, module, stale...); err != nil {
		return fmt.Errorf("clear previously migrated records of module %q: %w", module.Handle, err)
	}
	return nil
}

// track registers a module this migration is importing into and returns the
// bucket its records are collected in.
func (m *Migration) track(ns *composeTypes.Namespace, module *composeTypes.Module) *migratedModule {
	mm := &migratedModule{ns: ns, module: module}
	m.modules = append(m.modules, mm)
	return mm
}

// readRow lifts the identity and provenance the internal connection attached to
// a source row (see the Src* consts in internal_conn.go).
func readRow(row *importRow) *migratedRecord {
	get := func(k string) any {
		v, _ := row.GetValue(k, 0)
		return v
	}

	mr := &migratedRecord{
		srcID:     cast.ToUint64(get(SrcRecordID)),
		ownedBy:   cast.ToUint64(get(SrcOwnedBy)),
		createdBy: cast.ToUint64(get(SrcCreatedBy)),
		updatedBy: cast.ToUint64(get(SrcUpdatedBy)),
	}
	if t, ok := get(SrcCreatedAt).(time.Time); ok {
		mr.createdAt = t
	}
	// UpdatedAt is a *time.Time and is nil on a record nobody ever edited; the
	// type assertion has to be on the pointer or an untouched record would come
	// out looking edited at publish time.
	if t, ok := get(SrcUpdatedAt).(*time.Time); ok {
		mr.updatedAt = t
	}

	return mr
}

// meta is the stamp written on the created record. See migratedMetaKey.
func (mr *migratedRecord) meta() map[string]any {
	return map[string]any{
		migratedMetaKey: map[string]any{
			"sourceRecordID": strconv.FormatUint(mr.srcID, 10),
		},
	}
}

// isMigrated recognises the stamp meta() writes. The two must stay in step:
// clearing keys on this, so a stamp the reader does not recognise means a retry
// silently doubles every row.
func isMigrated(r *composeTypes.Record) bool {
	if r == nil || r.Meta == nil {
		return false
	}
	_, ok := r.Meta[migratedMetaKey]
	return ok
}

// Finalize runs the second pass: it rewrites the record links every module held
// back and restores the provenance the create overwrote.
//
// Both are done with one DAL-level write per record rather than through the
// record service. The service is the right door for a create -- it validates,
// defaults and sanitizes, and the target module may not have the shape the
// source had -- but it is the wrong door for this: procCreate resets id,
// revision, createdAt/By and ownership on purpose, and an update would stamp
// updatedAt/By just as firmly. Rather than teach ordinary record creation a
// "trust me" mode that every other caller then has to reason about, the fixup
// goes straight to the DAL, which is where the audit columns live.
//
// The record is re-read rather than reused from the create's result: what the
// service hands back has been through the value formatter, and writing those
// display-shaped values back would be a quiet corruption of, say, every
// DateTime field.
func (m *Migration) Finalize(ctx context.Context) error {
	var unresolved []string

	for _, mm := range m.modules {
		if len(mm.records) == 0 {
			continue
		}

		set, _, err := dalutils.ComposeRecordsList(ctx, m.im.dalSvc, mm.module, composeTypes.RecordFilter{
			ModuleID:    mm.module.ID,
			NamespaceID: mm.ns.ID,
		})
		if err != nil {
			return fmt.Errorf("re-read migrated records of module %q: %w", mm.module.Handle, err)
		}

		byID := make(map[uint64]*composeTypes.Record, len(set))
		for _, r := range set {
			byID[r.ID] = r
		}

		upd, uu := m.fixup(mm, byID)
		unresolved = append(unresolved, uu...)

		if len(upd) == 0 {
			continue
		}
		if err = dalutils.ComposeRecordUpdate(ctx, m.im.dalSvc, mm.module, upd...); err != nil {
			return fmt.Errorf("write migrated records of module %q: %w", mm.module.Handle, err)
		}
	}

	if len(unresolved) > 0 {
		return fmt.Errorf(
			"%d record link(s) could not be remapped; nothing was published. "+
				"They point at records whose module was not migrated: %v",
			len(unresolved), unresolved,
		)
	}

	return nil
}

// fixup patches the freshly re-read records of one module: links remapped
// through the migration's id map, provenance put back. Records the migration
// did not create are left alone -- someone may have typed a record into the
// draft namespace, and it is not ours to rewrite.
func (m *Migration) fixup(
	mm *migratedModule,
	byID map[uint64]*composeTypes.Record,
) (upd composeTypes.RecordSet, unresolved []string) {
	upd = make(composeTypes.RecordSet, 0, len(mm.records))

	for _, mr := range mm.records {
		rec, ok := byID[mr.newID]
		if !ok {
			continue
		}

		for _, lv := range mr.links {
			srcRef := cast.ToUint64(lv.Value)
			newRef, ok := m.idMap[srcRef]
			if !ok {
				// The row is written but its link would be empty, and a publish
				// that quietly drops a reference is worse than one that stops.
				// The draft is still a draft, so a retry re-clears and starts
				// over.
				unresolved = append(unresolved, fmt.Sprintf(
					"%s.%s -> record %d", mm.module.Handle, lv.Name, srcRef,
				))
				continue
			}
			rec.Values = rec.Values.Set(&composeTypes.RecordValue{
				Name:  lv.Name,
				Value: strconv.FormatUint(newRef, 10),
				Ref:   newRef,
				Place: lv.Place,
			})
		}

		rec.OwnedBy = mr.ownedBy
		rec.CreatedAt = mr.createdAt
		rec.CreatedBy = mr.createdBy
		rec.UpdatedAt = mr.updatedAt
		rec.UpdatedBy = mr.updatedBy
		rec.SetModule(mm.module)

		upd = append(upd, rec)
	}

	return
}
