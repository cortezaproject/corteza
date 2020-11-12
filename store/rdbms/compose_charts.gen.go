package rdbms

// This file is an auto-generated file
//
// Template:    pkg/codegen/assets/store_rdbms.gen.go.tpl
// Definitions: store/compose_charts.yaml
//
// Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated.

import (
	"context"
	"database/sql"
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/store"
)

var _ = errors.Is

// SearchComposeCharts returns all matching rows
//
// This function calls convertComposeChartFilter with the given
// types.ChartFilter and expects to receive a working squirrel.SelectBuilder
func (s Store) SearchComposeCharts(ctx context.Context, f types.ChartFilter) (types.ChartSet, types.ChartFilter, error) {
	var (
		err error
		set []*types.Chart
		q   squirrel.SelectBuilder
	)

	return set, f, func() error {
		q, err = s.convertComposeChartFilter(f)
		if err != nil {
			return err
		}

		// Paging enabled
		// {search: {enablePaging:true}}
		// Cleanup unwanted cursors (only relevant is f.PageCursor, next&prev are reset and returned)
		f.PrevPage, f.NextPage = nil, nil

		if f.PageCursor != nil {
			// Page cursor exists so we need to validate it against used sort
			if f.Sort, err = f.PageCursor.Sort(f.Sort); err != nil {
				return err
			}
		}

		if len(f.Sort) == 0 {
			f.Sort = filter.SortExprSet{}
		}

		// Make sure results are always sorted at least by primary keys
		if f.Sort.Get("id") == nil {
			f.Sort = append(f.Sort, &filter.SortExpr{Column: "id"})
		}

		sort := f.Sort.Clone()

		// When cursor for a previous page is used it's marked as reversed
		// This tells us to flip the descending flag on all used sort keys
		if f.PageCursor != nil && f.PageCursor.Reverse {
			sort.Reverse()
		}

		// Apply sorting expr from filter to query
		if q, err = setOrderBy(q, sort, s.sortableComposeChartColumns()); err != nil {
			return err
		}

		set, f.PrevPage, f.NextPage, err = s.fetchFullPageOfComposeCharts(ctx, q, sort.Columns(), sort.Reversed(), f.PageCursor, f.Limit, f.Check)
		if err != nil {
			return err
		}

		f.PageCursor = nil
		return nil
	}()
}

// fetchFullPageOfComposeCharts collects all requested results.
//
// Function applies:
//  - cursor conditions (where ...)
//  - limit
//
// Main responsibility of this function is to perform additional sequential queries in case when not enough results
// are collected due to failed check on a specific row (by check fn).
//
// Function then moves cursor to the last item fetched
func (s Store) fetchFullPageOfComposeCharts(
	ctx context.Context,
	q squirrel.SelectBuilder,
	sortColumns []string,
	sortDesc bool,
	cursor *filter.PagingCursor,
	limit uint,
	check func(*types.Chart) (bool, error),
) (set []*types.Chart, prev, next *filter.PagingCursor, err error) {
	var (
		aux []*types.Chart

		// When cursor for a previous page is used it's marked as reversed
		// This tells us to flip the descending flag on all used sort keys
		reversedOrder = cursor != nil && cursor.Reverse

		// copy of the select builder
		tryQuery squirrel.SelectBuilder

		fetched uint
	)

	set = make([]*types.Chart, 0, DefaultSliceCapacity)

	if cursor != nil {
		cursor.Reverse = sortDesc
	}

	for try := 0; try < MaxRefetches; try++ {
		tryQuery = setCursorCond(q, cursor)
		if limit > 0 {
			tryQuery = tryQuery.Limit(uint64(limit + 1))
		}

		if aux, fetched, _, err = s.QueryComposeCharts(ctx, tryQuery, check); err != nil {
			return nil, nil, nil, err
		}

		if cursor != nil && prev == nil && len(aux) > 0 {
			// Cursor for previous page is calculated only when cursor is used (so, not on first page)
			prev = s.collectComposeChartCursorValues(aux[0], sortColumns...)
		}

		// Point cursor to the last fetched element
		// if last != nil {
		if fetched >= limit && limit > 0 {
			next = s.collectComposeChartCursorValues(aux[limit-1], sortColumns...)
		}

		if limit > 0 && uint(len(aux)) >= limit {
			// we should use only as much as requested
			set = append(set, aux[:limit]...)
			break
		} else {
			set = append(set, aux...)
		}

		// if limit is not set or we've already collected enough items
		// we can break the loop right away
		if limit == 0 || fetched == 0 || fetched < limit {
			break
		}

		// In case limit is set very low and we've missed records in the first fetch,
		// make sure next fetch limit is a bit higher
		if limit < MinEnsureFetchLimit {
			limit = MinEnsureFetchLimit
		}

		// @todo improve strategy for collecting next page with lower limit
	}

	if reversedOrder {
		// Fetched set needs to be reversed because we've forced a descending order to get the previous page
		for i, j := 0, len(set)-1; i < j; i, j = i+1, j-1 {
			set[i], set[j] = set[j], set[i]
		}

		// and flip prev/next cursors too
		prev, next = next, prev
	}

	if prev != nil {
		prev.Reverse = true
	}

	return set, prev, next, nil
}

// QueryComposeCharts queries the database, converts and checks each row and
// returns collected set
//
// Fn also returns total number of fetched items and last fetched item so that the caller can construct cursor
// for next page of results
func (s Store) QueryComposeCharts(
	ctx context.Context,
	q squirrel.Sqlizer,
	check func(*types.Chart) (bool, error),
) ([]*types.Chart, uint, *types.Chart, error) {
	var (
		set = make([]*types.Chart, 0, DefaultSliceCapacity)
		res *types.Chart

		// Query rows with
		rows, err = s.Query(ctx, q)

		fetched uint
	)

	if err != nil {
		return nil, 0, nil, err
	}

	defer rows.Close()
	for rows.Next() {
		fetched++
		if err = rows.Err(); err == nil {
			res, err = s.internalComposeChartRowScanner(rows)
		}

		if err != nil {
			return nil, 0, nil, err
		}

		// check fn set, call it and see if it passed the test
		// if not, skip the item
		if check != nil {
			if chk, err := check(res); err != nil {
				return nil, 0, nil, err
			} else if !chk {
				continue
			}
		}

		set = append(set, res)
	}

	return set, fetched, res, rows.Err()
}

// LookupComposeChartByID searches for compose chart by ID
//
// It returns compose chart even if deleted
func (s Store) LookupComposeChartByID(ctx context.Context, id uint64) (*types.Chart, error) {
	return s.execLookupComposeChart(ctx, squirrel.Eq{
		s.preprocessColumn("cch.id", ""): store.PreprocessValue(id, ""),
	})
}

// LookupComposeChartByNamespaceIDHandle searches for compose chart by handle (case-insensitive)
func (s Store) LookupComposeChartByNamespaceIDHandle(ctx context.Context, namespace_id uint64, handle string) (*types.Chart, error) {
	return s.execLookupComposeChart(ctx, squirrel.Eq{
		s.preprocessColumn("cch.rel_namespace", ""): store.PreprocessValue(namespace_id, ""),
		s.preprocessColumn("cch.handle", "lower"):   store.PreprocessValue(handle, "lower"),

		"cch.deleted_at": nil,
	})
}

// CreateComposeChart creates one or more rows in compose_chart table
func (s Store) CreateComposeChart(ctx context.Context, rr ...*types.Chart) (err error) {
	for _, res := range rr {
		err = s.checkComposeChartConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execCreateComposeCharts(ctx, s.internalComposeChartEncoder(res))
		if err != nil {
			return err
		}
	}

	return
}

// UpdateComposeChart updates one or more existing rows in compose_chart
func (s Store) UpdateComposeChart(ctx context.Context, rr ...*types.Chart) error {
	return s.partialComposeChartUpdate(ctx, nil, rr...)
}

// partialComposeChartUpdate updates one or more existing rows in compose_chart
func (s Store) partialComposeChartUpdate(ctx context.Context, onlyColumns []string, rr ...*types.Chart) (err error) {
	for _, res := range rr {
		err = s.checkComposeChartConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpdateComposeCharts(
			ctx,
			squirrel.Eq{
				s.preprocessColumn("cch.id", ""): store.PreprocessValue(res.ID, ""),
			},
			s.internalComposeChartEncoder(res).Skip("id").Only(onlyColumns...))
		if err != nil {
			return err
		}
	}

	return
}

// UpsertComposeChart updates one or more existing rows in compose_chart
func (s Store) UpsertComposeChart(ctx context.Context, rr ...*types.Chart) (err error) {
	for _, res := range rr {
		err = s.checkComposeChartConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpsertComposeCharts(ctx, s.internalComposeChartEncoder(res))
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteComposeChart Deletes one or more rows from compose_chart table
func (s Store) DeleteComposeChart(ctx context.Context, rr ...*types.Chart) (err error) {
	for _, res := range rr {

		err = s.execDeleteComposeCharts(ctx, squirrel.Eq{
			s.preprocessColumn("cch.id", ""): store.PreprocessValue(res.ID, ""),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteComposeChartByID Deletes row from the compose_chart table
func (s Store) DeleteComposeChartByID(ctx context.Context, ID uint64) error {
	return s.execDeleteComposeCharts(ctx, squirrel.Eq{
		s.preprocessColumn("cch.id", ""): store.PreprocessValue(ID, ""),
	})
}

// TruncateComposeCharts Deletes all rows from the compose_chart table
func (s Store) TruncateComposeCharts(ctx context.Context) error {
	return s.Truncate(ctx, s.composeChartTable())
}

// execLookupComposeChart prepares ComposeChart query and executes it,
// returning types.Chart (or error)
func (s Store) execLookupComposeChart(ctx context.Context, cnd squirrel.Sqlizer) (res *types.Chart, err error) {
	var (
		row rowScanner
	)

	row, err = s.QueryRow(ctx, s.composeChartsSelectBuilder().Where(cnd))
	if err != nil {
		return
	}

	res, err = s.internalComposeChartRowScanner(row)
	if err != nil {
		return
	}

	return res, nil
}

// execCreateComposeCharts updates all matched (by cnd) rows in compose_chart with given data
func (s Store) execCreateComposeCharts(ctx context.Context, payload store.Payload) error {
	return s.Exec(ctx, s.InsertBuilder(s.composeChartTable()).SetMap(payload))
}

// execUpdateComposeCharts updates all matched (by cnd) rows in compose_chart with given data
func (s Store) execUpdateComposeCharts(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return s.Exec(ctx, s.UpdateBuilder(s.composeChartTable("cch")).Where(cnd).SetMap(set))
}

// execUpsertComposeCharts inserts new or updates matching (by-primary-key) rows in compose_chart with given data
func (s Store) execUpsertComposeCharts(ctx context.Context, set store.Payload) error {
	upsert, err := s.config.UpsertBuilder(
		s.config,
		s.composeChartTable(),
		set,
		"id",
	)

	if err != nil {
		return err
	}

	return s.Exec(ctx, upsert)
}

// execDeleteComposeCharts Deletes all matched (by cnd) rows in compose_chart with given data
func (s Store) execDeleteComposeCharts(ctx context.Context, cnd squirrel.Sqlizer) error {
	return s.Exec(ctx, s.DeleteBuilder(s.composeChartTable("cch")).Where(cnd))
}

func (s Store) internalComposeChartRowScanner(row rowScanner) (res *types.Chart, err error) {
	res = &types.Chart{}

	if _, has := s.config.RowScanners["composeChart"]; has {
		scanner := s.config.RowScanners["composeChart"].(func(_ rowScanner, _ *types.Chart) error)
		err = scanner(row, res)
	} else {
		err = row.Scan(
			&res.ID,
			&res.Handle,
			&res.Name,
			&res.Config,
			&res.NamespaceID,
			&res.CreatedAt,
			&res.UpdatedAt,
			&res.DeletedAt,
		)
	}

	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err != nil {
		return nil, errors.Store("could not scan composeChart db row").Wrap(err)
	} else {
		return res, nil
	}
}

// QueryComposeCharts returns squirrel.SelectBuilder with set table and all columns
func (s Store) composeChartsSelectBuilder() squirrel.SelectBuilder {
	return s.SelectBuilder(s.composeChartTable("cch"), s.composeChartColumns("cch")...)
}

// composeChartTable name of the db table
func (Store) composeChartTable(aa ...string) string {
	var alias string
	if len(aa) > 0 {
		alias = " AS " + aa[0]
	}

	return "compose_chart" + alias
}

// ComposeChartColumns returns all defined table columns
//
// With optional string arg, all columns are returned aliased
func (Store) composeChartColumns(aa ...string) []string {
	var alias string
	if len(aa) > 0 {
		alias = aa[0] + "."
	}

	return []string{
		alias + "id",
		alias + "handle",
		alias + "name",
		alias + "config",
		alias + "rel_namespace",
		alias + "created_at",
		alias + "updated_at",
		alias + "deleted_at",
	}
}

// {true true false true true true}

// sortableComposeChartColumns returns all ComposeChart columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
func (Store) sortableComposeChartColumns() map[string]string {
	return map[string]string{
		"id": "id", "handle": "handle", "name": "name", "created_at": "created_at",
		"createdat":  "created_at",
		"updated_at": "updated_at",
		"updatedat":  "updated_at",
		"deleted_at": "deleted_at",
		"deletedat":  "deleted_at",
	}
}

// internalComposeChartEncoder encodes fields from types.Chart to store.Payload (map)
//
// Encoding is done by using generic approach or by calling encodeComposeChart
// func when rdbms.customEncoder=true
func (s Store) internalComposeChartEncoder(res *types.Chart) store.Payload {
	return store.Payload{
		"id":            res.ID,
		"handle":        res.Handle,
		"name":          res.Name,
		"config":        res.Config,
		"rel_namespace": res.NamespaceID,
		"created_at":    res.CreatedAt,
		"updated_at":    res.UpdatedAt,
		"deleted_at":    res.DeletedAt,
	}
}

// collectComposeChartCursorValues collects values from the given resource that and sets them to the cursor
// to be used for pagination
//
// Values that are collected must come from sortable, unique or primary columns/fields
// At least one of the collected columns must be flagged as unique, otherwise fn appends primary keys at the end
//
// Known issue:
//   when collecting cursor values for query that sorts by unique column with partial index (ie: unique handle on
//   undeleted items)
func (s Store) collectComposeChartCursorValues(res *types.Chart, cc ...string) *filter.PagingCursor {
	var (
		cursor = &filter.PagingCursor{}

		hasUnique bool

		// All known primary key columns

		pkId bool

		collect = func(cc ...string) {
			for _, c := range cc {
				switch c {
				case "id":
					cursor.Set(c, res.ID, false)

					pkId = true
				case "handle":
					cursor.Set(c, res.Handle, false)
					hasUnique = true

				case "name":
					cursor.Set(c, res.Name, false)

				case "created_at":
					cursor.Set(c, res.CreatedAt, false)

				case "updated_at":
					cursor.Set(c, res.UpdatedAt, false)

				case "deleted_at":
					cursor.Set(c, res.DeletedAt, false)

				}
			}
		}
	)

	collect(cc...)
	if !hasUnique || !(pkId && true) {
		collect("id")
	}

	return cursor
}

// checkComposeChartConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant but unfortunately we can not rely
// on the full support (MySQL does not support conditional indexes)
func (s *Store) checkComposeChartConstraints(ctx context.Context, res *types.Chart) error {
	// Consider resource valid when all fields in unique constraint check lookups
	// have valid (non-empty) value
	//
	// Only string and uint64 are supported for now
	// feel free to add additional types if needed
	var valid = true

	if !valid {
		return nil
	}

	return nil
}
