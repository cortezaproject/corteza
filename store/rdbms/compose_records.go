package rdbms

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/ql"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/store/rdbms/builders"
)

const (
	composeRecordValueAliasPfx = "rv_"
	composeRecordValueJoinTpl  = "compose_record_value AS {alias} ON ({alias}.record_id = crd.id AND {alias}.name = '{field}' AND {alias}.deleted_at IS NULL)"
)

type (
	mftd struct {
		boolean  bool
		numeric  bool
		dateTime bool
		ref      bool
	}
)

func buildComposeRecordsCursor(cfg *Config, m *types.Module) func(cur *filter.PagingCursor) squirrel.Sqlizer {
	return func(cur *filter.PagingCursor) squirrel.Sqlizer {
		return builders.CursorCondition(cur, func(key string) (builders.KeyMap, error) {
			if col, fd, is := isRealRecordCol(key); is {
				_, _, tcp, _ := cfg.CastModuleFieldToColumnType(fd, key)
				// These values here won't be casted
				return builders.KeyMap{
					FieldCast:    col,
					TypeCast:     col,
					TypeCastPtrn: tcp,
				}, nil
			}

			f := m.Fields.FindByName(key)
			if f == nil {
				return builders.KeyMap{}, fmt.Errorf("unknown module field %q used in a cursor", key)
			}

			_, fcp, tcp, err := cfg.CastModuleFieldToColumnType(f, f.Name)
			if err != nil {
				return builders.KeyMap{}, err
			}

			fc := fmt.Sprintf(fcp, key)
			tc := fmt.Sprintf(tcp, fc)
			rr := builders.KeyMap{
				FieldCast:    fc,
				TypeCast:     tc,
				TypeCastPtrn: tcp,
			}

			return rr, nil
		})
	}
}

// @todo support for partitioned records (records are partitioned into multiple record tables
//       in this case, values are no longer separated into record_value (key-value) table but encoded as JSON
// @todo support for partitioned record with (optional) physical columns for record values
//       physical columns are part of module-field configuration

// SearchComposeRecords returns all matching ComposeRecords from store
func (s Store) SearchComposeRecords(ctx context.Context, m *types.Module, f types.RecordFilter) (types.RecordSet, types.RecordFilter, error) {
	var (
		set []*types.Record
		q   squirrel.SelectBuilder
	)

	return set, f, func() (err error) {
		q, err = s.convertComposeRecordFilter(m, f)
		if err != nil {
			return
		}

		f.PrevPage, f.NextPage = nil, nil

		if f.PageCursor != nil {
			if f.IncPageNavigation || f.IncTotal {
				return fmt.Errorf("not allowed to fetch page navigation or total item count with page cursor")
			}

			// Page cursor exists so we need to validate it against used sort
			// To cover the case when paging cursor is set but sorting is empty, we collect the sorting instructions
			// from the cursor.
			// This (extracted sorting info) is then returned as part of response
			if f.Sort, err = f.PageCursor.Sort(nil); err != nil {
				return err
			}
		}

		// Make sure results are always sorted at least by primary keys
		if f.Sort.Get("id") == nil {
			f.Sort = append(f.Sort, &filter.SortExpr{
				Column:     "id",
				Descending: f.Sort.LastDescending(),
			})
		}

		// Prevent sorting over multi-value fields
		//
		// Due to how values are currently storred, this causes duplication.
		// For now, we'll prevent this and address in future releases.
		for _, s := range f.Sort {
			f := m.Fields.FindByName(s.Column)
			if f != nil && f.Multi {
				return fmt.Errorf("not allowed to sort by multi-value fields: %s", s.Column)
			}
		}

		// Cloned sorting instructions for the actual sorting
		// Original are passed to the fetchFullPageOfUsers fn used for cursor creation so it MUST keep the initial
		// direction information
		sort := f.Sort.Clone()

		// When cursor for a previous page is used it's marked as reversed
		// This tells us to flip the descending flag on all used sort keys
		if f.PageCursor != nil && f.PageCursor.ROrder {
			sort.Reverse()
		}

		// Apply sorting expr from filter to query
		if q, err = s.composeRecordsSorter(m, q, sort); err != nil {
			return
		}

		set, f.PrevPage, f.NextPage, err = s.fetchFullPageOfComposeRecords(
			ctx, m, q,
			f.Sort, f.PageCursor, f.Limit,
			f.Check,
			buildComposeRecordsCursor(s.config, m),
		)

		if err != nil {
			return
		}

		if f.IncPageNavigation || f.IncTotal {
			if f.IncTotal {
				// Calc total from the number of items fetched
				// even if we do build the page navigation
				f.Total = uint(len(set))
			}
			if f.Limit > 0 && uint(len(set)) == f.Limit {
				// Build page navigation ONLY when limit is set and
				// there are less items fetched then requested limit
				if nav, err := s.composeRecordsPageNavigation(ctx, m, q, f, f.Sort, buildComposeRecordsCursor(s.config, m)); err != nil {
					return err
				} else {
					f.Total = nav.Total
					f.PageNavigation = nav.PageNavigation
				}

				if !f.IncPageNavigation {
					// remove page navigation if not requested
					f.PageNavigation = nil
				}

			}
		}

		f.PageCursor = nil
		return nil
	}()
}

// fetches a large number of records and generates page navigation and/or counts all records
//
// Func tries to minimize amount of data fetched & processed by loading 1000x set of records per page
// We're storing only records (and fetch their values) used for generating values
//
// Caveats:
//  - this will only work if Check function DOES NOT RELY on record values
func (s Store) composeRecordsPageNavigation(
	ctx context.Context,
	m *types.Module,
	q squirrel.SelectBuilder,
	f types.RecordFilter,
	sort filter.SortExprSet,
	cursorCond func(*filter.PagingCursor) squirrel.Sqlizer,
) (types.RecordFilter, error) {
	return f, func() (err error) {
		var (
			// holds set of last records from each page
			set []*types.Record

			supPageLimitFactor = 1000
			supPageLimit       = f.Limit * uint(supPageLimitFactor)
			// super-page query; fetches large number of records to optimize page navigation and counting
			supPageQuery = q.Limit(uint64(supPageLimit))

			rows *sql.Rows
			res  *types.Record

			cursor *filter.PagingCursor
		)

		// reset output on filter
		f.Total = 0
		f.PageNavigation = make([]*filter.Page, 0, supPageLimitFactor)

		for {
			rowsValid, rowsProcessed := 0, 0

			// reuse q for query with/without cursor
			if cursor == nil {
				// first page
				q = supPageQuery
			} else {
				//cursor.LThen = sort.Reversed()
				q = supPageQuery.Where(cursorCond(cursor))
			}

			rows, err = s.Query(ctx, q)
			if err != nil {
				return err
			}

			defer rows.Close()
			for rows.Next() {
				rowsProcessed++
				if err = rows.Err(); err == nil {
					res, err = s.internalComposeRecordRowScanner(m, rows)
				}

				if err != nil {
					return err
				}

				// check fn set, call it and see if it passed the test
				// if not, skip the item
				if f.Check != nil {
					if chk, err := f.Check(res); err != nil {
						return err
					} else if !chk {
						continue
					}
				}

				rowsValid++
				f.Total++

				if f.Total%f.Limit == 0 {
					set = append(set, res)
				}
			}

			if f.Total%f.Limit != 0 {
				// append the rest
				set = append(set, res)
			}

			// fetch values for all last-records fetched in this iteration
			if err = s.composeRecordPostLoadProcessor(ctx, m, set[len(f.PageNavigation):]...); err != nil {
				return err
			}

			// generate cursor values for all records fetched in this iteration
			for _, r := range set[len(f.PageNavigation):] {
				f.PageNavigation = append(f.PageNavigation, &filter.Page{
					Page:   uint(len(f.PageNavigation) + 1),
					Count:  uint(rowsValid),
					Cursor: s.collectComposeRecordCursorValues(m, r, sort...),
				})
			}

			if uint(rowsProcessed) < supPageLimit {
				// fetched less than planned, assume entire matching record set is processed
				break
			}

			// get ready for next iteration and generate cursor
			cursor = f.PageNavigation[len(f.PageNavigation)-1].Cursor
		}

		// And prefix on at the beginning to point to 1st page (for consistency)
		if len(f.PageNavigation) > 0 {
			// reorder cursors and move them forward for 1 page,
			// going from last to 1st
			for p := len(f.PageNavigation) - 1; p > 0; p-- {
				f.PageNavigation[p].Cursor = f.PageNavigation[p-1].Cursor
			}

			// no need for cursor on the 1st page.
			f.PageNavigation[0].Cursor = nil
		}

		return nil
	}()
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

			err = s.createComposeRecordValue(ctx, nil, res.Values.GetClean()...)
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

		identResolver = func(sortBy bool) func(i ql.Ident) (ql.Ident, error) {
			return func(i ql.Ident) (ql.Ident, error) {
				var is bool
				if i.Value, _, is = isRealRecordCol(i.Value); is {
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

					if sortBy {
						var sortCol string
						sortCol, _, _, err = s.config.CastModuleFieldToColumnType(m.Fields.FindByName(i.Value), i.Value)
						if err != nil {
							return i, err
						}

						query = query.Column(squirrel.Alias(squirrel.Expr(sortCol), composeRecordValueAliasPfx+i.Value))
					}

				}

				return s.FieldToColumnTypeCaster(m.Fields.FindByName(i.Value), i)
			}
		}
	)

	// Create query for fetching and counting records.
	query = s.composeRecordsSelectBuilder().
		Prefix("SELECT "+strings.Join(s.composeRecordColumns("sub"), ", ")+" FROM (").
		Suffix(") AS sub").
		Distinct().
		Where("crd.module_id = ?", m.ID).
		Where("crd.rel_namespace = ?", m.NamespaceID)

	// Inc/exclude deleted records according to filter settings
	query = filter.StateCondition(query, "crd.deleted_at", f.Deleted)

	if len(f.LabeledIDs) > 0 {
		query = query.Where(squirrel.Eq{"crd.id": f.LabeledIDs})
	}

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
		fp.OnIdent = identResolver(false)

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
		sp.OnIdent = identResolver(true)

		if _, err = sp.ParseColumns(f.Sort.String()); err != nil {
			return
		}
	}

	if f.PageCursor != nil {
		var (
			// Sort parser
			sp = ql.NewParser()
		)

		// Resolve all identifiers found in sort
		// into their table/column counterparts
		sp.OnIdent = identResolver(true)

		if _, err = sp.ParseColumns(strings.Join(f.PageCursor.Keys(), ", ")); err != nil {
			return
		}
	}

	return
}

//func (s Store) convertComposeRecordCursor(m *types.Module, from *filter.PagingCursor) (to *filter.PagingCursor) {
//	if from != nil {
//		to = &filter.PagingCursor{ROrder: from.ROrder, LThen: from.LThen}
//		// convert cursor keys field names (if used)
//		from.Walk(func(key string, val interface{}, desc bool) {
//			if col, has := s.sortableComposeRecordColumns()[strings.ToLower(key)]; has {
//				key = col
//			} else if f := m.Fields.FindByName(key); f != nil {
//				key, _ = s.config.CastModuleFieldToColumnType(f, key)
//			} else {
//				return
//			}
//
//			to.Set(key, val, desc)
//		})
//
//	}
//
//	return to
//}

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
		sortable = s.sortableComposeRecordColumns()
		sqlSort  = make([]string, len(sort))
	)

	for i, c := range sort {
		var err error
		if col, has := sortable[strings.ToLower(c.Column)]; has {
			sqlSort[i] = col
		} else if f := m.Fields.FindByName(c.Column); f != nil {

			sqlSort[i], _, _, err = s.config.CastModuleFieldToColumnType(f, c.Column)
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
func (s Store) collectComposeRecordCursorValues(m *types.Module, res *types.Record, sort ...*filter.SortExpr) *filter.PagingCursor {
	var (
		cursor = &filter.PagingCursor{LThen: filter.SortExprSet(sort).Reversed()}

		hasUnique bool
		pkID      bool

		conv = s.sortableComposeRecordColumns()

		collect = func(cc ...*filter.SortExpr) {
			for _, c := range cc {
				switch conv[strings.ToLower(c.Column)] {
				case "id":
					cursor.Set(c.Column, res.ID, c.Descending)
					pkID = true
				case "created_at":
					cursor.Set(c.Column, res.CreatedAt, c.Descending)
				case "created_by":
					cursor.Set(c.Column, res.CreatedBy, c.Descending)
				case "updated_at":
					cursor.Set(c.Column, res.UpdatedAt, c.Descending)
				case "updated_by":
					cursor.Set(c.Column, res.UpdatedBy, c.Descending)
				case "deleted_at":
					cursor.Set(c.Column, res.DeletedAt, c.Descending)
				case "deleted_by":
					cursor.Set(c.Column, res.DeletedBy, c.Descending)
				case "owned_by":
					cursor.Set(c.Column, res.OwnedBy, c.Descending)
				default:
					if rv := res.Values.Get(c.Column, 0); rv != nil {
						cursor.Set(c.Column, rv.Value, c.Descending)
					} else {
						cursor.Set(c.Column, nil, c.Descending)
					}
				}
			}
		}
	)

	collect(sort...)
	if !hasUnique || !pkID {
		collect(&filter.SortExpr{Column: "id"})
	}

	return cursor
}

//// Checks if field name is "real column", reformats it and returns
func isRealRecordCol(name string) (string, mftd, bool) {
	switch name {
	case
		"recordID",
		"id":
		return "crd.id", mftd{ref: true}, true

	case
		"module_id",
		"owned_by",
		"created_by",
		"updated_by",
		"deleted_by":
		return "crd." + name, mftd{ref: true}, true

	case
		"created_at",
		"updated_at",
		"deleted_at":
		return "crd." + name, mftd{dateTime: true}, true

	case
		"moduleID",
		"ownedBy",
		"createdBy",
		"updatedBy",
		"deletedBy":
		return "crd." + name[0:len(name)-2] + "_" + strings.ToLower(name[len(name)-2:]), mftd{ref: true}, true

	case
		"createdAt",
		"updatedAt",
		"deletedAt":
		return "crd." + name[0:len(name)-2] + "_" + strings.ToLower(name[len(name)-2:]), mftd{dateTime: true}, true
	}

	return name, mftd{}, false
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

func (t mftd) IsBoolean() bool {
	return t.boolean
}
func (t mftd) IsNumeric() bool {
	return t.numeric
}
func (t mftd) IsDateTime() bool {
	return t.dateTime
}
func (t mftd) IsRef() bool {
	return t.ref
}
