package rdbms

// This file is an auto-generated file
//
// Template:    pkg/codegen/assets/store_rdbms.gen.go.tpl
// Definitions: store/actionlog.yaml
//
// Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated.

import (
	"context"
	"database/sql"
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/store"
)

var _ = errors.Is

// SearchActionlogs returns all matching rows
//
// This function calls convertActionlogFilter with the given
// actionlog.Filter and expects to receive a working squirrel.SelectBuilder
func (s Store) SearchActionlogs(ctx context.Context, f actionlog.Filter) (actionlog.ActionSet, actionlog.Filter, error) {
	var (
		err error
		set []*actionlog.Action
		q   squirrel.SelectBuilder
	)

	return set, f, func() error {
		q, err = s.convertActionlogFilter(f)
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
			f.Sort = append(f.Sort, &filter.SortExpr{Column: "id", Descending: true})
		}

		sort := f.Sort.Clone()

		// When cursor for a previous page is used it's marked as reversed
		// This tells us to flip the descending flag on all used sort keys
		if f.PageCursor != nil && f.PageCursor.Reverse {
			sort.Reverse()
		}

		// Apply sorting expr from filter to query
		if q, err = setOrderBy(q, sort, s.sortableActionlogColumns()); err != nil {
			return err
		}

		set, f.PrevPage, f.NextPage, err = s.fetchFullPageOfActionlogs(ctx, q, sort.Columns(), sort.Reversed(), f.PageCursor, f.Limit, nil)
		if err != nil {
			return err
		}

		f.PageCursor = nil
		return nil
	}()
}

// fetchFullPageOfActionlogs collects all requested results.
//
// Function applies:
//  - cursor conditions (where ...)
//  - limit
//
// Main responsibility of this function is to perform additional sequential queries in case when not enough results
// are collected due to failed check on a specific row (by check fn).
//
// Function then moves cursor to the last item fetched
func (s Store) fetchFullPageOfActionlogs(
	ctx context.Context,
	q squirrel.SelectBuilder,
	sortColumns []string,
	sortDesc bool,
	cursor *filter.PagingCursor,
	limit uint,
	check func(*actionlog.Action) (bool, error),
) (set []*actionlog.Action, prev, next *filter.PagingCursor, err error) {
	var (
		aux []*actionlog.Action

		// When cursor for a previous page is used it's marked as reversed
		// This tells us to flip the descending flag on all used sort keys
		reversedOrder = cursor != nil && cursor.Reverse

		// copy of the select builder
		tryQuery squirrel.SelectBuilder

		fetched uint
	)

	set = make([]*actionlog.Action, 0, DefaultSliceCapacity)

	if cursor != nil {
		cursor.Reverse = sortDesc
	}

	for try := 0; try < MaxRefetches; try++ {
		tryQuery = setCursorCond(q, cursor)
		if limit > 0 {
			tryQuery = tryQuery.Limit(uint64(limit + 1))
		}

		if aux, err = s.QueryActionlogs(ctx, tryQuery, check); err != nil {
			return nil, nil, nil, err
		}

		fetched = uint(len(aux))
		if cursor != nil && prev == nil && fetched > 0 {
			// Cursor for previous page is calculated only when cursor is used (so, not on first page)
			prev = s.collectActionlogCursorValues(aux[0], sortColumns...)
		}

		// Point cursor to the last fetched element
		if fetched > limit && limit > 0 {
			next = s.collectActionlogCursorValues(aux[limit-1], sortColumns...)

			// we should use only as much as requested
			set = append(set, aux[:limit]...)
			break
		} else {
			set = append(set, aux...)
		}

		// if limit is not set or we've already collected enough items
		// we can break the loop right away
		if limit == 0 || fetched == 0 || fetched <= limit {
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

// QueryActionlogs queries the database, converts and checks each row and
// returns collected set
//
// Fn also returns total number of fetched items and last fetched item so that the caller can construct cursor
// for next page of results
func (s Store) QueryActionlogs(
	ctx context.Context,
	q squirrel.Sqlizer,
	check func(*actionlog.Action) (bool, error),
) ([]*actionlog.Action, error) {
	var (
		set = make([]*actionlog.Action, 0, DefaultSliceCapacity)
		res *actionlog.Action

		// Query rows with
		rows, err = s.Query(ctx, q)
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		if err = rows.Err(); err == nil {
			res, err = s.internalActionlogRowScanner(rows)
		}

		if err != nil {
			return nil, err
		}

		set = append(set, res)
	}

	return set, rows.Err()
}

// CreateActionlog creates one or more rows in actionlog table
func (s Store) CreateActionlog(ctx context.Context, rr ...*actionlog.Action) (err error) {
	for _, res := range rr {
		err = s.checkActionlogConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execCreateActionlogs(ctx, s.internalActionlogEncoder(res))
		if err != nil {
			return err
		}
	}

	return
}

// TruncateActionlogs Deletes all rows from the actionlog table
func (s Store) TruncateActionlogs(ctx context.Context) error {
	return s.Truncate(ctx, s.actionlogTable())
}

// execLookupActionlog prepares Actionlog query and executes it,
// returning actionlog.Action (or error)
func (s Store) execLookupActionlog(ctx context.Context, cnd squirrel.Sqlizer) (res *actionlog.Action, err error) {
	var (
		row rowScanner
	)

	row, err = s.QueryRow(ctx, s.actionlogsSelectBuilder().Where(cnd))
	if err != nil {
		return
	}

	res, err = s.internalActionlogRowScanner(row)
	if err != nil {
		return
	}

	return res, nil
}

// execCreateActionlogs updates all matched (by cnd) rows in actionlog with given data
func (s Store) execCreateActionlogs(ctx context.Context, payload store.Payload) error {
	return s.Exec(ctx, s.InsertBuilder(s.actionlogTable()).SetMap(payload))
}

func (s Store) internalActionlogRowScanner(row rowScanner) (res *actionlog.Action, err error) {
	res = &actionlog.Action{}

	if _, has := s.config.RowScanners["actionlog"]; has {
		scanner := s.config.RowScanners["actionlog"].(func(_ rowScanner, _ *actionlog.Action) error)
		err = scanner(row, res)
	} else {
		err = s.scanActionlogRow(row, res)
	}

	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err != nil {
		return nil, errors.Store("could not scan actionlog db row").Wrap(err)
	} else {
		return res, nil
	}
}

// QueryActionlogs returns squirrel.SelectBuilder with set table and all columns
func (s Store) actionlogsSelectBuilder() squirrel.SelectBuilder {
	return s.SelectBuilder(s.actionlogTable("alg"), s.actionlogColumns("alg")...)
}

// actionlogTable name of the db table
func (Store) actionlogTable(aa ...string) string {
	var alias string
	if len(aa) > 0 {
		alias = " AS " + aa[0]
	}

	return "actionlog" + alias
}

// ActionlogColumns returns all defined table columns
//
// With optional string arg, all columns are returned aliased
func (Store) actionlogColumns(aa ...string) []string {
	var alias string
	if len(aa) > 0 {
		alias = aa[0] + "."
	}

	return []string{
		alias + "id",
		alias + "ts",
		alias + "request_origin",
		alias + "request_id",
		alias + "actor_ip_addr",
		alias + "actor_id",
		alias + "resource",
		alias + "action",
		alias + "error",
		alias + "severity",
		alias + "description",
		alias + "meta",
	}
}

// {true true false true true false}

// sortableActionlogColumns returns all Actionlog columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
func (Store) sortableActionlogColumns() map[string]string {
	return map[string]string{
		"id": "id",
	}
}

// internalActionlogEncoder encodes fields from actionlog.Action to store.Payload (map)
//
// Encoding is done by using generic approach or by calling encodeActionlog
// func when rdbms.customEncoder=true
func (s Store) internalActionlogEncoder(res *actionlog.Action) store.Payload {
	return s.encodeActionlog(res)
}

// collectActionlogCursorValues collects values from the given resource that and sets them to the cursor
// to be used for pagination
//
// Values that are collected must come from sortable, unique or primary columns/fields
// At least one of the collected columns must be flagged as unique, otherwise fn appends primary keys at the end
//
// Known issue:
//   when collecting cursor values for query that sorts by unique column with partial index (ie: unique handle on
//   undeleted items)
func (s Store) collectActionlogCursorValues(res *actionlog.Action, cc ...string) *filter.PagingCursor {
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

// checkActionlogConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant but unfortunately we can not rely
// on the full support (MySQL does not support conditional indexes)
func (s *Store) checkActionlogConstraints(ctx context.Context, res *actionlog.Action) error {
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
