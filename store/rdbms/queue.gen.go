package rdbms

// This file is an auto-generated file
//
// Template:    pkg/codegen/assets/store_rdbms.gen.go.tpl
// Definitions: store/queue.yaml
//
// Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated.

import (
	"context"
	"database/sql"
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/store/rdbms/builders"
	"github.com/cortezaproject/corteza-server/system/types"
)

var _ = errors.Is

// SearchQueues returns all matching rows
//
// This function calls convertQueueFilter with the given
// types.QueueFilter and expects to receive a working squirrel.SelectBuilder
func (s Store) SearchQueues(ctx context.Context, f types.QueueFilter) (types.QueueSet, types.QueueFilter, error) {
	var (
		err error
		set []*types.Queue
		q   squirrel.SelectBuilder
	)

	return set, f, func() error {
		q, err = s.convertQueueFilter(f)
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

		set, f.PrevPage, f.NextPage, err = s.fetchFullPageOfQueues(
			ctx,
			q, f.Sort, f.PageCursor,
			f.Limit,
			nil,
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

// fetchFullPageOfQueues collects all requested results.
//
// Function applies:
//  - cursor conditions (where ...)
//  - limit
//
// Main responsibility of this function is to perform additional sequential queries in case when not enough results
// are collected due to failed check on a specific row (by check fn).
//
// Function then moves cursor to the last item fetched
func (s Store) fetchFullPageOfQueues(
	ctx context.Context,
	q squirrel.SelectBuilder,
	sort filter.SortExprSet,
	cursor *filter.PagingCursor,
	reqItems uint,
	check func(*types.Queue) (bool, error),
	cursorCond func(*filter.PagingCursor) squirrel.Sqlizer,
) (set []*types.Queue, prev, next *filter.PagingCursor, err error) {
	var (
		aux []*types.Queue

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

	set = make([]*types.Queue, 0, DefaultSliceCapacity)

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

		if aux, err = s.QueryQueues(ctx, tryQuery, check); err != nil {
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
			cursor = s.collectQueueCursorValues(set[collected-1], sort...)

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
		prev = s.collectQueueCursorValues(set[0], sort...)
		prev.ROrder = true
		prev.LThen = !sort.Reversed()
	}

	if hasNext {
		next = s.collectQueueCursorValues(set[collected-1], sort...)
		next.LThen = sort.Reversed()
	}

	return set, prev, next, nil
}

// QueryQueues queries the database, converts and checks each row and
// returns collected set
//
// Fn also returns total number of fetched items and last fetched item so that the caller can construct cursor
// for next page of results
func (s Store) QueryQueues(
	ctx context.Context,
	q squirrel.Sqlizer,
	check func(*types.Queue) (bool, error),
) ([]*types.Queue, error) {
	var (
		tmp = make([]*types.Queue, 0, DefaultSliceCapacity)
		set = make([]*types.Queue, 0, DefaultSliceCapacity)
		res *types.Queue

		// Query rows with
		rows, err = s.Query(ctx, q)
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		if err = rows.Err(); err == nil {
			res, err = s.internalQueueRowScanner(rows)
		}

		if err != nil {
			return nil, err
		}

		tmp = append(tmp, res)
	}

	for _, res = range tmp {

		set = append(set, res)
	}

	return set, nil
}

// LookupQueueByID searches for queue by ID
func (s Store) LookupQueueByID(ctx context.Context, id uint64) (*types.Queue, error) {
	return s.execLookupQueue(ctx, squirrel.Eq{
		s.preprocessColumn("mqs.id", ""): store.PreprocessValue(id, ""),
	})
}

// LookupQueueByQueue searches for queue by queue name
func (s Store) LookupQueueByQueue(ctx context.Context, queue string) (*types.Queue, error) {
	return s.execLookupQueue(ctx, squirrel.Eq{
		s.preprocessColumn("mqs.queue", ""): store.PreprocessValue(queue, ""),
	})
}

// CreateQueue creates one or more rows in queue_settings table
func (s Store) CreateQueue(ctx context.Context, rr ...*types.Queue) (err error) {
	for _, res := range rr {
		err = s.checkQueueConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execCreateQueues(ctx, s.internalQueueEncoder(res))
		if err != nil {
			return err
		}
	}

	return
}

// UpdateQueue updates one or more existing rows in queue_settings
func (s Store) UpdateQueue(ctx context.Context, rr ...*types.Queue) error {
	return s.partialQueueUpdate(ctx, nil, rr...)
}

// partialQueueUpdate updates one or more existing rows in queue_settings
func (s Store) partialQueueUpdate(ctx context.Context, onlyColumns []string, rr ...*types.Queue) (err error) {
	for _, res := range rr {
		err = s.checkQueueConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpdateQueues(
			ctx,
			squirrel.Eq{
				s.preprocessColumn("mqs.id", ""): store.PreprocessValue(res.ID, ""),
			},
			s.internalQueueEncoder(res).Skip("id").Only(onlyColumns...))
		if err != nil {
			return err
		}
	}

	return
}

// UpsertQueue updates one or more existing rows in queue_settings
func (s Store) UpsertQueue(ctx context.Context, rr ...*types.Queue) (err error) {
	for _, res := range rr {
		err = s.checkQueueConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpsertQueues(ctx, s.internalQueueEncoder(res))
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteQueue Deletes one or more rows from queue_settings table
func (s Store) DeleteQueue(ctx context.Context, rr ...*types.Queue) (err error) {
	for _, res := range rr {

		err = s.execDeleteQueues(ctx, squirrel.Eq{
			s.preprocessColumn("mqs.id", ""): store.PreprocessValue(res.ID, ""),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteQueueByID Deletes row from the queue_settings table
func (s Store) DeleteQueueByID(ctx context.Context, ID uint64) error {
	return s.execDeleteQueues(ctx, squirrel.Eq{
		s.preprocessColumn("mqs.id", ""): store.PreprocessValue(ID, ""),
	})
}

// TruncateQueues Deletes all rows from the queue_settings table
func (s Store) TruncateQueues(ctx context.Context) error {
	return s.Truncate(ctx, s.queueTable())
}

// execLookupQueue prepares Queue query and executes it,
// returning types.Queue (or error)
func (s Store) execLookupQueue(ctx context.Context, cnd squirrel.Sqlizer) (res *types.Queue, err error) {
	var (
		row rowScanner
	)

	row, err = s.QueryRow(ctx, s.queuesSelectBuilder().Where(cnd))
	if err != nil {
		return
	}

	res, err = s.internalQueueRowScanner(row)
	if err != nil {
		return
	}

	return res, nil
}

// execCreateQueues updates all matched (by cnd) rows in queue_settings with given data
func (s Store) execCreateQueues(ctx context.Context, payload store.Payload) error {
	return s.Exec(ctx, s.InsertBuilder(s.queueTable()).SetMap(payload))
}

// execUpdateQueues updates all matched (by cnd) rows in queue_settings with given data
func (s Store) execUpdateQueues(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return s.Exec(ctx, s.UpdateBuilder(s.queueTable("mqs")).Where(cnd).SetMap(set))
}

// execUpsertQueues inserts new or updates matching (by-primary-key) rows in queue_settings with given data
func (s Store) execUpsertQueues(ctx context.Context, set store.Payload) error {
	upsert, err := s.config.UpsertBuilder(
		s.config,
		s.queueTable(),
		set,
		s.preprocessColumn("id", ""),
	)

	if err != nil {
		return err
	}

	return s.Exec(ctx, upsert)
}

// execDeleteQueues Deletes all matched (by cnd) rows in queue_settings with given data
func (s Store) execDeleteQueues(ctx context.Context, cnd squirrel.Sqlizer) error {
	return s.Exec(ctx, s.DeleteBuilder(s.queueTable("mqs")).Where(cnd))
}

func (s Store) internalQueueRowScanner(row rowScanner) (res *types.Queue, err error) {
	res = &types.Queue{}

	err = row.Scan(
		&res.ID,
		&res.Consumer,
		&res.Queue,
		&res.Meta,
		&res.CreatedBy,
		&res.UpdatedBy,
		&res.DeletedBy,
		&res.CreatedAt,
		&res.UpdatedAt,
		&res.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err != nil {
		return nil, errors.Store("could not scan queue db row: %s", err).Wrap(err)
	} else {
		return res, nil
	}
}

// QueryQueues returns squirrel.SelectBuilder with set table and all columns
func (s Store) queuesSelectBuilder() squirrel.SelectBuilder {
	return s.SelectBuilder(s.queueTable("mqs"), s.queueColumns("mqs")...)
}

// queueTable name of the db table
func (Store) queueTable(aa ...string) string {
	var alias string
	if len(aa) > 0 {
		alias = " AS " + aa[0]
	}

	return "queue_settings" + alias
}

// QueueColumns returns all defined table columns
//
// With optional string arg, all columns are returned aliased
func (Store) queueColumns(aa ...string) []string {
	var alias string
	if len(aa) > 0 {
		alias = aa[0] + "."
	}

	return []string{
		alias + "id",
		alias + "consumer",
		alias + "queue",
		alias + "meta",
		alias + "created_by",
		alias + "updated_by",
		alias + "deleted_by",
		alias + "created_at",
		alias + "updated_at",
		alias + "deleted_at",
	}
}

// {true true false true false false}

// internalQueueEncoder encodes fields from types.Queue to store.Payload (map)
//
// Encoding is done by using generic approach or by calling encodeQueue
// func when rdbms.customEncoder=true
func (s Store) internalQueueEncoder(res *types.Queue) store.Payload {
	return store.Payload{
		"id":         res.ID,
		"consumer":   res.Consumer,
		"queue":      res.Queue,
		"meta":       res.Meta,
		"created_by": res.CreatedBy,
		"updated_by": res.UpdatedBy,
		"deleted_by": res.DeletedBy,
		"created_at": res.CreatedAt,
		"updated_at": res.UpdatedAt,
		"deleted_at": res.DeletedAt,
	}
}

// collectQueueCursorValues collects values from the given resource that and sets them to the cursor
// to be used for pagination
//
// Values that are collected must come from sortable, unique or primary columns/fields
// At least one of the collected columns must be flagged as unique, otherwise fn appends primary keys at the end
//
// Known issue:
//   when collecting cursor values for query that sorts by unique column with partial index (ie: unique handle on
//   undeleted items)
func (s Store) collectQueueCursorValues(res *types.Queue, cc ...*filter.SortExpr) *filter.PagingCursor {
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
				case "queue":
					cursor.Set(c.Column, res.Queue, c.Descending)

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

// checkQueueConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant but unfortunately we cannot rely
// on the full support (MySQL does not support conditional indexes)
func (s *Store) checkQueueConstraints(ctx context.Context, res *types.Queue) error {
	// Consider resource valid when all fields in unique constraint check lookups
	// have valid (non-empty) value
	//
	// Only string and uint64 are supported for now
	// feel free to add additional types if needed
	var valid = true

	if !valid {
		return nil
	}

	var checks = make([]func() error, 0)

	for _, check := range checks {
		if err := check(); err != nil {
			return err
		}
	}

	return nil
}
