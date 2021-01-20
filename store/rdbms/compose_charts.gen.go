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
	"github.com/cortezaproject/corteza-server/store/rdbms/builders"
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
		// Cleanup unwanted cursor values (only relevant is f.PageCursor, next&prev are reset and returned)
		f.PrevPage, f.NextPage = nil, nil

		if f.PageCursor != nil {
			// Page cursor exists so we need to validate it against used sort
			// To cover the case when paging cursor is set but sorting is empty, we collect the sorting instructions
			// from the cursor.
			// This (extracted sorting info) is then returned as part of response
			if f.Sort, err = f.PageCursor.Sort(f.Sort); err != nil {
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
		if q, err = setOrderBy(q, sort, s.sortableComposeChartColumns()); err != nil {
			return err
		}

		set, f.PrevPage, f.NextPage, err = s.fetchFullPageOfComposeCharts(
			ctx,
			q, f.Sort, f.PageCursor,
			f.Limit,
			f.Check,
			func(cur *filter.PagingCursor) squirrel.Sqlizer {
				return builders.CursorCondition(cur, nil)
			},
		)

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
	sort filter.SortExprSet,
	cursor *filter.PagingCursor,
	reqItems uint,
	check func(*types.Chart) (bool, error),
	cursorCond func(*filter.PagingCursor) squirrel.Sqlizer,
) (set []*types.Chart, prev, next *filter.PagingCursor, err error) {
	var (
		aux []*types.Chart

		// When cursor for a previous page is used it's marked as reversed
		// This tells us to flip the descending flag on all used sort keys
		reversedOrder = cursor != nil && cursor.ROrder

		// copy of the select builder
		tryQuery squirrel.SelectBuilder

		// Copy no. of required items to limit
		// Limit will change when doing subsequent queries to fill
		// the set with all required items
		limit = reqItems

		// cursor to prev. page is only calculated when cursor is used
		hasPrev = cursor != nil

		// next cursor is calculated when there are more pages to come
		hasNext bool
	)

	set = make([]*types.Chart, 0, DefaultSliceCapacity)

	for try := 0; try < MaxRefetches; try++ {
		if cursor != nil {
			tryQuery = q.Where(cursorCond(cursor))
		} else {
			tryQuery = q
		}

		if limit > 0 {
			// fetching + 1 so we know if there are more items
			// we can fetch (next-page cursor)
			tryQuery = tryQuery.Limit(uint64(limit + 1))
		}

		if aux, err = s.QueryComposeCharts(ctx, tryQuery, check); err != nil {
			return nil, nil, nil, err
		}

		if len(aux) == 0 {
			// nothing fetched
			break
		}

		// append fetched items
		set = append(set, aux...)

		if reqItems == 0 {
			// no max requested items specified, break out
			break
		}

		collected := uint(len(set))

		if reqItems > collected {
			// not enough items fetched, try again with adjusted limit
			limit = reqItems - collected

			if limit < MinEnsureFetchLimit {
				// In case limit is set very low and we've missed records in the first fetch,
				// make sure next fetch limit is a bit higher
				limit = MinEnsureFetchLimit
			}

			// Update cursor so that it points to the last item fetched
			cursor = s.collectComposeChartCursorValues(set[collected-1], sort...)

			// Copy reverse flag from sorting
			cursor.LThen = sort.Reversed()
			continue
		}

		if reqItems < collected {
			set = set[:reqItems]
			hasNext = true
		}

		break
	}

	collected := len(set)

	if collected == 0 {
		return nil, nil, nil, nil
	}

	if reversedOrder {
		// Fetched set needs to be reversed because we've forced a descending order to get the previous page
		for i, j := 0, collected-1; i < j; i, j = i+1, j-1 {
			set[i], set[j] = set[j], set[i]
		}

		// when in reverse-order rules on what cursor to return change
		hasPrev, hasNext = hasNext, hasPrev
	}

	if hasPrev {
		prev = s.collectComposeChartCursorValues(set[0], sort...)
		prev.ROrder = true
		prev.LThen = !sort.Reversed()
	}

	if hasNext {
		next = s.collectComposeChartCursorValues(set[collected-1], sort...)
		next.LThen = sort.Reversed()
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
) ([]*types.Chart, error) {
	var (
		set = make([]*types.Chart, 0, DefaultSliceCapacity)
		res *types.Chart

		// Query rows with
		rows, err = s.Query(ctx, q)
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		if err = rows.Err(); err == nil {
			res, err = s.internalComposeChartRowScanner(rows)
		}

		if err != nil {
			return nil, err
		}

		// check fn set, call it and see if it passed the test
		// if not, skip the item
		if check != nil {
			if chk, err := check(res); err != nil {
				return nil, err
			} else if !chk {
				continue
			}
		}

		set = append(set, res)
	}

	return set, rows.Err()
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
		s.preprocessColumn("id", ""),
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
		return nil, errors.Store("could not scan composeChart db row: %s", err).Wrap(err)
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
func (s Store) collectComposeChartCursorValues(res *types.Chart, cc ...*filter.SortExpr) *filter.PagingCursor {
	var (
		cursor = &filter.PagingCursor{LThen: filter.SortExprSet(cc).Reversed()}

		hasUnique bool

		// All known primary key columns

		pkId bool

		collect = func(cc ...*filter.SortExpr) {
			for _, c := range cc {
				switch c.Column {
				case "id":
					cursor.Set(c.Column, res.ID, c.Descending)

					pkId = true
				case "handle":
					cursor.Set(c.Column, res.Handle, c.Descending)
					hasUnique = true

				case "name":
					cursor.Set(c.Column, res.Name, c.Descending)

				case "created_at":
					cursor.Set(c.Column, res.CreatedAt, c.Descending)

				case "updated_at":
					cursor.Set(c.Column, res.UpdatedAt, c.Descending)

				case "deleted_at":
					cursor.Set(c.Column, res.DeletedAt, c.Descending)

				}
			}
		}
	)

	collect(cc...)
	if !hasUnique || !(pkId && true) {
		collect(&filter.SortExpr{Column: "id", Descending: false})
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
