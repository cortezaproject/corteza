package rdbms

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/ql"
	"github.com/cortezaproject/corteza-server/pkg/slice"
	"github.com/cortezaproject/corteza-server/store"
	"strings"
)

const (
	composeRecordValueAliasPfx = "rv_"
	composeRecordValueJoinTpl  = "compose_record_value AS {alias} ON ({alias}.record_id = crd.id AND {alias}.name = '{field}' AND {alias}.deleted_at IS NULL)"
)

// @todo support for partitioned records (records are partitioned into multiple record tables
//       in this case, values are no longer separated into record_value (key-value) table but encoded as JSON
// @todo support for partitioned record with (optional) physical columns for record values
//       physical columns are part of module-field configuration

// SearchComposeRecords returns all matching ComposeRecords from store
func (s Store) SearchComposeRecords(ctx context.Context, m *types.Module, filter types.RecordFilter) (set types.RecordSet, f types.RecordFilter, err error) {
	// In when module requires this,
	set, f, err = s.searchComposeRecords(ctx, m, filter)
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

		// Make sure all record-values are linked to the record
		res.Values.SetRecordID(res.ID)

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

			// Make sure all record-values are linked to the record
			res.Values.SetRecordID(res.ID)

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
	err = s.deleteComposeRecordByID(ctx, nil, ID)
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

func (s Store) ComposeRecordReport(ctx context.Context, m *types.Module, metrics, dimensions, filter string) ([]map[string]interface{}, error) {
	return ComposeRecordReportBuilder(&s, m, metrics, dimensions, filter).Run(ctx)
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
				join := composeRecordValueJoinTpl
				join = strings.ReplaceAll(join, "{alias}", composeRecordValueAliasPfx+i.Value)
				join = strings.ReplaceAll(join, "{field}", i.Value)
				query = query.LeftJoin(join)
			}

			return s.FieldToColumnTypeCaster(m.Fields.FindByName(i.Value), i)
		}
	)

	// Create query for fetching and counting records.
	query = s.composeRecordsSelectBuilder().
		Where("crd.module_id = ?", m.ID).
		Where("crd.rel_namespace = ?", m.NamespaceID)

	// Inc/exclude deleted records according to filter settings
	query = filter.StateCondition(query, "crd.deleted_at", f.Deleted)

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

	if len(f.Sort) > 0 {
		var (
			// Sort parser
			sp = ql.NewParser()
		)

		// Resolve all identifiers found in sort
		// into their table/column counterparts
		sp.OnIdent = identResolver

		if _, err = sp.ParseColumns(f.Sort.String()); err != nil {
			return
		}
	}

	return
}

func (s Store) composeRecordPostLoadProcessor(ctx context.Context, m *types.Module, set ...*types.Record) (err error) {
	if len(set) > 0 {
		// Load all related record values and append them to each record
		var (
			rvs types.RecordValueSet
		)
		rvs, _, err = s.searchComposeRecordValues(ctx, nil, types.RecordValueFilter{
			RecordID: types.RecordSet(set).IDs(),
			Deleted:  filter.StateInclusive,
		})
		if err != nil {
			return
		}

		for r := range set {
			set[r].Values = rvs.FilterByRecordID(set[r].ID)
		}
	}

	return nil
}

func (s Store) composeRecordsSorter(m *types.Module, q squirrel.SelectBuilder, sort filter.SortExprSet) (squirrel.SelectBuilder, error) {
	var (
		sortable = slice.ToStringBoolMap(s.sortableComposeRecordColumns())
		sqlSort  = make([]string, len(sort))
	)

	for i, c := range sort {
		var err error
		if sortable[c.Column] {
			sqlSort[i] = c.Column
		} else if f := m.Fields.FindByName(c.Column); f != nil {
			sqlSort[i], err = s.config.CastModuleFieldToColumnType(f, c.Column)
		} else {
			err = fmt.Errorf("could not sort by unknown column: %s", c.Column)
		}

		if err != nil {
			return q, err
		}

		// Apply proper sorting param
		sqlSort[i] = s.config.SqlSortHandler(sqlSort[i], sort[i].Descending)
	}

	return q.OrderBy(sqlSort...), nil
}

// Custom implementation for collecting cursor values from compose records AND it's values!
func (s Store) collectComposeRecordCursorValues(res *types.Record, cc ...string) *filter.PagingCursor {
	var (
		cursor = &filter.PagingCursor{}

		hasUnique bool
		pkID      bool

		collect = func(cc ...string) {
			for _, c := range cc {
				switch c {
				case "crd.id":
					cursor.Set(c, res.ID, false)
					pkID = true
				case "crd.created_at":
					cursor.Set(c, res.CreatedAt, false)
				case "crd.updated_at":
					cursor.Set(c, res.UpdatedAt, false)
				case "crd.deleted_at":
					cursor.Set(c, res.DeletedAt, false)
				default:
					if rv := res.Values.Get(c, 0); rv != nil {
						cursor.Set(fmt.Sprintf("%s%s.value", composeRecordValueAliasPfx, c), rv.Value, false)
					} else {
						cursor.Set(fmt.Sprintf("%s%s.value", composeRecordValueAliasPfx, c), nil, false)
					}
				}
			}
		}
	)

	collect(cc...)
	if !hasUnique || !pkID {
		collect(
			"crd.id",
		)
	}

	return cursor
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
