package rdbms

// This file is an auto-generated file
//
// Template:    pkg/codegen/assets/store_rdbms.gen.go.tpl
// Definitions: store/messagebus_queue_message.yaml
//
// Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated.

import (
	"context"
	"database/sql"
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/messagebus"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/store/rdbms/builders"
)

var _ = errors.Is

// SearchMessagebusQueueMessages returns all matching rows
//
// This function calls convertMessagebusQueueMessageFilter with the given
// messagebus.QueueMessageFilter and expects to receive a working squirrel.SelectBuilder
func (s Store) SearchMessagebusQueueMessages(ctx context.Context, f messagebus.QueueMessageFilter) (messagebus.QueueMessageSet, messagebus.QueueMessageFilter, error) {
	var (
		err error
		set []*messagebus.QueueMessage
		q   squirrel.SelectBuilder
	)

	return set, f, func() error {
		q, err = s.convertMessagebusQueueMessageFilter(f)
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
		if q, err = setOrderBy(q, sort, s.sortableMessagebusQueueMessageColumns()); err != nil {
			return err
		}

		set, f.PrevPage, f.NextPage, err = s.fetchFullPageOfMessagebusQueueMessages(
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

// fetchFullPageOfMessagebusQueueMessages collects all requested results.
//
// Function applies:
//  - cursor conditions (where ...)
//  - limit
//
// Main responsibility of this function is to perform additional sequential queries in case when not enough results
// are collected due to failed check on a specific row (by check fn).
//
// Function then moves cursor to the last item fetched
func (s Store) fetchFullPageOfMessagebusQueueMessages(
	ctx context.Context,
	q squirrel.SelectBuilder,
	sort filter.SortExprSet,
	cursor *filter.PagingCursor,
	reqItems uint,
	check func(*messagebus.QueueMessage) (bool, error),
	cursorCond func(*filter.PagingCursor) squirrel.Sqlizer,
) (set []*messagebus.QueueMessage, prev, next *filter.PagingCursor, err error) {
	var (
		aux []*messagebus.QueueMessage

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

	set = make([]*messagebus.QueueMessage, 0, DefaultSliceCapacity)

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

		if aux, err = s.QueryMessagebusQueueMessages(ctx, tryQuery, check); err != nil {
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
			cursor = s.collectMessagebusQueueMessageCursorValues(set[collected-1], sort...)

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
		prev = s.collectMessagebusQueueMessageCursorValues(set[0], sort...)
		prev.ROrder = true
		prev.LThen = !sort.Reversed()
	}

	if hasNext {
		next = s.collectMessagebusQueueMessageCursorValues(set[collected-1], sort...)
		next.LThen = sort.Reversed()
	}

	return set, prev, next, nil
}

// QueryMessagebusQueueMessages queries the database, converts and checks each row and
// returns collected set
//
// Fn also returns total number of fetched items and last fetched item so that the caller can construct cursor
// for next page of results
func (s Store) QueryMessagebusQueueMessages(
	ctx context.Context,
	q squirrel.Sqlizer,
	check func(*messagebus.QueueMessage) (bool, error),
) ([]*messagebus.QueueMessage, error) {
	var (
		tmp = make([]*messagebus.QueueMessage, 0, DefaultSliceCapacity)
		set = make([]*messagebus.QueueMessage, 0, DefaultSliceCapacity)
		res *messagebus.QueueMessage

		// Query rows with
		rows, err = s.Query(ctx, q)
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		if err = rows.Err(); err == nil {
			res, err = s.internalMessagebusQueueMessageRowScanner(rows)
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

// CreateMessagebusQueueMessage creates one or more rows in queue_messages table
func (s Store) CreateMessagebusQueueMessage(ctx context.Context, rr ...*messagebus.QueueMessage) (err error) {
	for _, res := range rr {
		err = s.checkMessagebusQueueMessageConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execCreateMessagebusQueueMessages(ctx, s.internalMessagebusQueueMessageEncoder(res))
		if err != nil {
			return err
		}
	}

	return
}

// UpdateMessagebusQueueMessage updates one or more existing rows in queue_messages
func (s Store) UpdateMessagebusQueueMessage(ctx context.Context, rr ...*messagebus.QueueMessage) error {
	return s.partialMessagebusQueueMessageUpdate(ctx, nil, rr...)
}

// partialMessagebusQueueMessageUpdate updates one or more existing rows in queue_messages
func (s Store) partialMessagebusQueueMessageUpdate(ctx context.Context, onlyColumns []string, rr ...*messagebus.QueueMessage) (err error) {
	for _, res := range rr {
		err = s.checkMessagebusQueueMessageConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpdateMessagebusQueueMessages(
			ctx,
			squirrel.Eq{
				s.preprocessColumn("mqm.id", ""): store.PreprocessValue(res.ID, ""),
			},
			s.internalMessagebusQueueMessageEncoder(res).Skip("id").Only(onlyColumns...))
		if err != nil {
			return err
		}
	}

	return
}

// DeleteMessagebusQueueMessage Deletes one or more rows from queue_messages table
func (s Store) DeleteMessagebusQueueMessage(ctx context.Context, rr ...*messagebus.QueueMessage) (err error) {
	for _, res := range rr {

		err = s.execDeleteMessagebusQueueMessages(ctx, squirrel.Eq{
			s.preprocessColumn("mqm.id", ""): store.PreprocessValue(res.ID, ""),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteMessagebusQueueMessageByID Deletes row from the queue_messages table
func (s Store) DeleteMessagebusQueueMessageByID(ctx context.Context, ID uint64) error {
	return s.execDeleteMessagebusQueueMessages(ctx, squirrel.Eq{
		s.preprocessColumn("mqm.id", ""): store.PreprocessValue(ID, ""),
	})
}

// TruncateMessagebusQueueMessages Deletes all rows from the queue_messages table
func (s Store) TruncateMessagebusQueueMessages(ctx context.Context) error {
	return s.Truncate(ctx, s.messagebusQueueMessageTable())
}

// execLookupMessagebusQueueMessage prepares MessagebusQueueMessage query and executes it,
// returning messagebus.QueueMessage (or error)
func (s Store) execLookupMessagebusQueueMessage(ctx context.Context, cnd squirrel.Sqlizer) (res *messagebus.QueueMessage, err error) {
	var (
		row rowScanner
	)

	row, err = s.QueryRow(ctx, s.messagebusQueueMessagesSelectBuilder().Where(cnd))
	if err != nil {
		return
	}

	res, err = s.internalMessagebusQueueMessageRowScanner(row)
	if err != nil {
		return
	}

	return res, nil
}

// execCreateMessagebusQueueMessages updates all matched (by cnd) rows in queue_messages with given data
func (s Store) execCreateMessagebusQueueMessages(ctx context.Context, payload store.Payload) error {
	return s.Exec(ctx, s.InsertBuilder(s.messagebusQueueMessageTable()).SetMap(payload))
}

// execUpdateMessagebusQueueMessages updates all matched (by cnd) rows in queue_messages with given data
func (s Store) execUpdateMessagebusQueueMessages(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return s.Exec(ctx, s.UpdateBuilder(s.messagebusQueueMessageTable("mqm")).Where(cnd).SetMap(set))
}

// execDeleteMessagebusQueueMessages Deletes all matched (by cnd) rows in queue_messages with given data
func (s Store) execDeleteMessagebusQueueMessages(ctx context.Context, cnd squirrel.Sqlizer) error {
	return s.Exec(ctx, s.DeleteBuilder(s.messagebusQueueMessageTable("mqm")).Where(cnd))
}

func (s Store) internalMessagebusQueueMessageRowScanner(row rowScanner) (res *messagebus.QueueMessage, err error) {
	res = &messagebus.QueueMessage{}

	if _, has := s.config.RowScanners["messagebusQueueMessage"]; has {
		scanner := s.config.RowScanners["messagebusQueueMessage"].(func(_ rowScanner, _ *messagebus.QueueMessage) error)
		err = scanner(row, res)
	} else {
		err = row.Scan(
			&res.ID,
			&res.Queue,
			&res.Payload,
			&res.Processed,
			&res.Created,
		)
	}

	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err != nil {
		return nil, errors.Store("could not scan messagebusQueueMessage db row: %s", err).Wrap(err)
	} else {
		return res, nil
	}
}

// QueryMessagebusQueueMessages returns squirrel.SelectBuilder with set table and all columns
func (s Store) messagebusQueueMessagesSelectBuilder() squirrel.SelectBuilder {
	return s.SelectBuilder(s.messagebusQueueMessageTable("mqm"), s.messagebusQueueMessageColumns("mqm")...)
}

// messagebusQueueMessageTable name of the db table
func (Store) messagebusQueueMessageTable(aa ...string) string {
	var alias string
	if len(aa) > 0 {
		alias = " AS " + aa[0]
	}

	return "queue_messages" + alias
}

// MessagebusQueueMessageColumns returns all defined table columns
//
// With optional string arg, all columns are returned aliased
func (Store) messagebusQueueMessageColumns(aa ...string) []string {
	var alias string
	if len(aa) > 0 {
		alias = aa[0] + "."
	}

	return []string{
		alias + "id",
		alias + "queue",
		alias + "payload",
		alias + "processed",
		alias + "created",
	}
}

// {true true false true true false}

// sortableMessagebusQueueMessageColumns returns all MessagebusQueueMessage columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
func (Store) sortableMessagebusQueueMessageColumns() map[string]string {
	return map[string]string{
		"id": "id",
	}
}

// internalMessagebusQueueMessageEncoder encodes fields from messagebus.QueueMessage to store.Payload (map)
//
// Encoding is done by using generic approach or by calling encodeMessagebusQueueMessage
// func when rdbms.customEncoder=true
func (s Store) internalMessagebusQueueMessageEncoder(res *messagebus.QueueMessage) store.Payload {
	return store.Payload{
		"id":        res.ID,
		"queue":     res.Queue,
		"payload":   res.Payload,
		"processed": res.Processed,
		"created":   res.Created,
	}
}

// collectMessagebusQueueMessageCursorValues collects values from the given resource that and sets them to the cursor
// to be used for pagination
//
// Values that are collected must come from sortable, unique or primary columns/fields
// At least one of the collected columns must be flagged as unique, otherwise fn appends primary keys at the end
//
// Known issue:
//   when collecting cursor values for query that sorts by unique column with partial index (ie: unique handle on
//   undeleted items)
func (s Store) collectMessagebusQueueMessageCursorValues(res *messagebus.QueueMessage, cc ...*filter.SortExpr) *filter.PagingCursor {
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

// checkMessagebusQueueMessageConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant but unfortunately we cannot rely
// on the full support (MySQL does not support conditional indexes)
func (s *Store) checkMessagebusQueueMessageConstraints(ctx context.Context, res *messagebus.QueueMessage) error {
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
