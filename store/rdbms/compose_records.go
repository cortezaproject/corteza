package rdbms

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/ql"
	"github.com/cortezaproject/corteza-server/pkg/rh"
	"github.com/cortezaproject/corteza-server/store"
	"strings"
)

// @todo support for partitioned records (records are partitioned into multiple record tables
//       in this case, values are no longer separated into record_value (key-value) table but encoded as JSON
// @todo support for partitioned record with (optional) physical columns for record values
//       physical columns are part of module-field configuration

// SearchComposeRecords returns all matching ComposeRecords from store
func (s Store) SearchComposeRecords(ctx context.Context, m *types.Module, f types.RecordFilter) (set types.RecordSet, filter types.RecordFilter, err error) {
	// In when module requires this,
	set, filter, err = s.searchComposeRecords(ctx, m, f)
	if err != nil {
		return
	}

	return
}

// LookupComposeRecordByID searches for compose record by ID
// It returns compose record even if deleted
func (s Store) LookupComposeRecordByID(ctx context.Context, _ *types.Module, id uint64) (res *types.Record, err error) {
	res, err = s.lookupComposeRecordByID(ctx, nil, id)
	if err != nil {
		return
	}

	res.Values, _, err = s.searchComposeRecordValues(ctx, nil, types.RecordValueFilter{RecordID: []uint64{id}})
	if err != nil {
		return nil, err
	}

	return
}

// CreateComposeRecord creates one or more ComposeRecords in store
func (s Store) CreateComposeRecord(ctx context.Context, m *types.Module, rr ...*types.Record) (err error) {
	if err = validateRecordModule(m, rr...); err != nil {
		return
	}

	for _, res := range rr {

		err = s.createComposeRecord(ctx, nil, res)
		if err != nil {
			return
		}

		err = s.createComposeRecordValue(ctx, nil, res.Values...)
		if err != nil {
			return
		}
	}

	return
}

// UpdateComposeRecord updates one or more (existing) ComposeRecords in store
func (s Store) UpdateComposeRecord(ctx context.Context, m *types.Module, rr ...*types.Record) (err error) {
	if err = validateRecordModule(m, rr...); err != nil {
		return
	}

	for _, res := range rr {
		err = s.updateComposeRecord(ctx, nil, res)
		if err != nil {
			return
		}

		cnd := squirrel.Eq{"record_id": res.ID}

		if res.DeletedAt != nil {
			// Record was deleted, set all values to deleted too
			err = s.execUpdateComposeRecordValues(ctx, cnd, store.Payload{"deleted_at": res.DeletedAt})
			if err != nil {
				return
			}
		} else {
			// we're following the old implementation here for now
			// all old values are Deleted and new inserted
			err = s.execDeleteComposeRecordValues(ctx, cnd)
			if err != nil {
				return
			}

			err = s.createComposeRecordValue(ctx, nil, res.Values...)
		}
	}

	return
}

// PartialComposeRecordUpdate updates one or more existing ComposeRecords in store
func (s Store) PartialComposeRecordUpdate(ctx context.Context, m *types.Module, onlyColumns []string, rr ...*types.Record) (err error) {
	if err = validateRecordModule(m, rr...); err != nil {
		return
	}

	return fmt.Errorf("partial compose record update is not supported")
}

// DeleteComposeRecord Deletes one or more ComposeRecords from store
func (s Store) DeleteComposeRecord(ctx context.Context, m *types.Module, rr ...*types.Record) (err error) {
	if err = validateRecordModule(m, rr...); err != nil {
		return
	}

	for _, res := range rr {
		err = s.DeleteComposeRecordByID(ctx, m, res.ID)
		if err != nil {
			return
		}
	}

	return
}

// DeleteComposeRecord Deletes one or more ComposeRecords from store
func (s Store) UpsertComposeRecord(ctx context.Context, m *types.Module, rr ...*types.Record) (err error) {
	if err = validateRecordModule(m, rr...); err != nil {
		return
	}

	for _, res := range rr {
		err = s.upsertComposeRecord(ctx, m, res)
		if err != nil {
			return
		}
	}

	return
}

// DeleteComposeRecordByID Deletes ComposeRecord from store
func (s Store) DeleteComposeRecordByID(ctx context.Context, _ *types.Module, ID uint64) (err error) {
	err = s.DeleteComposeRecordByID(ctx, nil, ID)
	if err != nil {
		return
	}

	err = s.Exec(ctx, s.DeleteBuilder(s.composeRecordValueTable("crv")).Where(squirrel.Eq{"crv.record_id": ID}))
	if err != nil {
		return
	}

	return
}

// TruncateComposeRecords Deletes all ComposeRecords from store
func (s Store) TruncateComposeRecords(ctx context.Context, _ *types.Module) (err error) {
	err = s.truncateComposeRecords(ctx, nil)
	if err != nil {
		return
	}

	err = s.truncateComposeRecordValues(ctx, nil)
	if err != nil {
		return
	}

	return
}

func (s Store) convertComposeRecordFilter(m *types.Module, f types.RecordFilter) (query squirrel.SelectBuilder, err error) {
	if m == nil {
		err = fmt.Errorf("module not provided")
		return
	}

	if f.ModuleID == 0 {
		// In case Module ID on filter is not set,
		// values from provided module can be used
		f.ModuleID = m.ID
	} else if m.ID != f.ModuleID {
		err = fmt.Errorf("provided module does not match filter module ID")
	}

	if f.NamespaceID == 0 {
		// In case Namespace ID on filter is not set,
		// values from provided module can be used
		f.NamespaceID = m.NamespaceID
	} else if m.NamespaceID != f.NamespaceID {
		err = fmt.Errorf("provided module namespace ID does not match filter namespace ID")
	}

	var (
		joinedFields  = []string{}
		alreadyJoined = func(f string) bool {
			for _, a := range joinedFields {
				if a == f {
					return true
				}
			}

			joinedFields = append(joinedFields, f)
			return false
		}

		identResolver = func(i ql.Ident) (ql.Ident, error) {
			var is bool
			if i.Value, is = isRealRecordCol(i.Value); is {
				i.Value += " "
				return i, nil
			}

			if !m.Fields.HasName(i.Value) {
				return i, fmt.Errorf("unknown field %q", i.Value)
			}

			if !alreadyJoined(i.Value) {
				query = query.LeftJoin(fmt.Sprintf(
					"compose_record_value AS rv_%s ON (rv_%s.record_id = crd.id AND rv_%s.name = ? AND rv_%s.deleted_at IS NULL)",
					i.Value, i.Value, i.Value, i.Value,
				), i.Value)
			}

			field := m.Fields.FindByName(i.Value)

			switch true {
			case field.IsBoolean():
				i.Value = fmt.Sprintf("(rv_%s.value NOT IN ('', '0', 'false', 'f',  'FALSE', 'F', false))", i.Value)
			case field.IsNumeric():
				i.Value = fmt.Sprintf("CAST(rv_%s.value AS SIGNED)", i.Value)
			case field.IsDateTime():
				i.Value = fmt.Sprintf("CAST(rv_%s.value AS DATETIME)", i.Value)
			case field.IsRef():
				i.Value = fmt.Sprintf("rv_%s.ref ", i.Value)
			default:
				i.Value = fmt.Sprintf("rv_%s.value ", i.Value)
			}

			return i, nil
		}
	)

	// Create query for fetching and counting records.
	query = s.composeRecordsSelectBuilder().
		Where("crd.module_id = ?", m.ID).
		Where("crd.rel_namespace = ?", m.NamespaceID)

	// Inc/exclude deleted records according to filter settings
	query = rh.FilterNullByState(query, "crd.deleted_at", f.Deleted)

	// Parse filters.
	if f.Query != "" {
		var (
			// Filter parser
			fp = ql.NewParser()

			// Filter node
			fn ql.ASTNode
		)

		// Resolve all identifiers found in the query
		// into their table/column counterparts
		fp.OnIdent = identResolver

		if fn, err = fp.ParseExpression(f.Query); err != nil {
			return
		} else if filterSql, filterArgs, err := fn.ToSql(); err != nil {
			return query, err
		} else {
			query = query.Where("("+filterSql+")", filterArgs...)
		}
	}

	// @todo refactor
	//if f.Sort != "" {
	//	var (
	//		// Sort parser
	//		sp = ql.NewParser()
	//
	//		// Sort columns
	//		sc ql.Columns
	//	)
	//
	//	// Resolve all identifiers found in sort
	//	// into their table/column counterparts
	//	sp.OnIdent = identResolver
	//
	//	if sc, err = sp.ParseColumns(f.Sort); err != nil {
	//		return
	//	}
	//
	//	query = query.OrderBy(sc.Strings()...)
	//}

	return
}

//// Checks if field name is "real column", reformats it and returns
func isRealRecordCol(name string) (string, bool) {
	switch name {
	case
		"recordID",
		"id":
		return "crd.id", true
	case
		"module_id",
		"owned_by",
		"created_by",
		"created_at",
		"updated_by",
		"updated_at",
		"deleted_by",
		"deleted_at":
		return "crd." + name, true

	case
		"moduleID",
		"ownedBy",
		"createdBy",
		"createdAt",
		"updatedBy",
		"updatedAt",
		"deletedBy",
		"deletedAt":
		return "crd." + name[0:len(name)-2] + "_" + strings.ToLower(name[len(name)-2:]), true
	}

	return name, false
}

// Verifies if module and namespace ID on record match IDs on module
func validateRecordModule(m *types.Module, rr ...*types.Record) error {
	for _, r := range rr {
		if m.ID != r.ModuleID {
			return fmt.Errorf("provided module does not match module ID on record")
		}

		if m.NamespaceID != r.NamespaceID {
			return fmt.Errorf("provided module namespace ID does not match namespace ID on record")
		}
	}

	return nil
}
