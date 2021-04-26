package rdbms

// This file is an auto-generated file
//
// Template:    pkg/codegen/assets/store_rdbms.gen.go.tpl
// Definitions: store/messagebus_queue_settings.yaml
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

// SearchMessagebusQueueSettings returns all matching rows
//
// This function calls convertMessagebusQueueSettingFilter with the given
// messagebus.QueueSettingsFilter and expects to receive a working squirrel.SelectBuilder
func (s Store) SearchMessagebusQueueSettings(ctx context.Context, f messagebus.QueueSettingsFilter) (messagebus.QueueSettingsSet, messagebus.QueueSettingsFilter, error) {
	var (
		err error
		set []*messagebus.QueueSettings
		q   squirrel.SelectBuilder
	)

	return set, f, func() error {
		q, err = s.convertMessagebusQueueSettingFilter(f)
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

		set, f.PrevPage, f.NextPage, err = s.fetchFullPageOfMessagebusQueueSettings(
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

// fetchFullPageOfMessagebusQueueSettings collects all requested results.
//
// Function applies:
//  - cursor conditions (where ...)
//  - limit
//
// Main responsibility of this function is to perform additional sequential queries in case when not enough results
// are collected due to failed check on a specific row (by check fn).
//
// Function then moves cursor to the last item fetched
func (s Store) fetchFullPageOfMessagebusQueueSettings(
	ctx context.Context,
	q squirrel.SelectBuilder,
	sort filter.SortExprSet,
	cursor *filter.PagingCursor,
	reqItems uint,
	check func(*messagebus.QueueSettings) (bool, error),
	cursorCond func(*filter.PagingCursor) squirrel.Sqlizer,
) (set []*messagebus.QueueSettings, prev, next *filter.PagingCursor, err error) {
	var (
		aux []*messagebus.QueueSettings

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

	set = make([]*messagebus.QueueSettings, 0, DefaultSliceCapacity)

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

		if aux, err = s.QueryMessagebusQueueSettings(ctx, tryQuery, check); err != nil {
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
			cursor = s.collectMessagebusQueueSettingCursorValues(set[collected-1], sort...)

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
		prev = s.collectMessagebusQueueSettingCursorValues(set[0], sort...)
		prev.ROrder = true
		prev.LThen = !sort.Reversed()
	}

	if hasNext {
		next = s.collectMessagebusQueueSettingCursorValues(set[collected-1], sort...)
		next.LThen = sort.Reversed()
	}

	return set, prev, next, nil
}

// QueryMessagebusQueueSettings queries the database, converts and checks each row and
// returns collected set
//
// Fn also returns total number of fetched items and last fetched item so that the caller can construct cursor
// for next page of results
func (s Store) QueryMessagebusQueueSettings(
	ctx context.Context,
	q squirrel.Sqlizer,
	check func(*messagebus.QueueSettings) (bool, error),
) ([]*messagebus.QueueSettings, error) {
	var (
		set = make([]*messagebus.QueueSettings, 0, DefaultSliceCapacity)
		res *messagebus.QueueSettings

		// Query rows with
		rows, err = s.Query(ctx, q)
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		if err = rows.Err(); err == nil {
			res, err = s.internalMessagebusQueueSettingRowScanner(rows)
		}

		if err != nil {
			return nil, err
		}

		set = append(set, res)
	}

	return set, rows.Err()
}

// LookupMessagebusQueueSettingByID searches for queue by ID
func (s Store) LookupMessagebusQueueSettingByID(ctx context.Context, id uint64) (*messagebus.QueueSettings, error) {
	return s.execLookupMessagebusQueueSetting(ctx, squirrel.Eq{
		s.preprocessColumn("mqs.id", ""): store.PreprocessValue(id, ""),
	})
}

// LookupMessagebusQueueSettingByQueue searches for queue by queue name
func (s Store) LookupMessagebusQueueSettingByQueue(ctx context.Context, queue string) (*messagebus.QueueSettings, error) {
	return s.execLookupMessagebusQueueSetting(ctx, squirrel.Eq{
		s.preprocessColumn("mqs.queue", ""): store.PreprocessValue(queue, ""),
	})
}

// CreateMessagebusQueueSetting creates one or more rows in queue_settings table
func (s Store) CreateMessagebusQueueSetting(ctx context.Context, rr ...*messagebus.QueueSettings) (err error) {
	for _, res := range rr {
		err = s.checkMessagebusQueueSettingConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execCreateMessagebusQueueSettings(ctx, s.internalMessagebusQueueSettingEncoder(res))
		if err != nil {
			return err
		}
	}

	return
}

// UpdateMessagebusQueueSetting updates one or more existing rows in queue_settings
func (s Store) UpdateMessagebusQueueSetting(ctx context.Context, rr ...*messagebus.QueueSettings) error {
	return s.partialMessagebusQueueSettingUpdate(ctx, nil, rr...)
}

// partialMessagebusQueueSettingUpdate updates one or more existing rows in queue_settings
func (s Store) partialMessagebusQueueSettingUpdate(ctx context.Context, onlyColumns []string, rr ...*messagebus.QueueSettings) (err error) {
	for _, res := range rr {
		err = s.checkMessagebusQueueSettingConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpdateMessagebusQueueSettings(
			ctx,
			squirrel.Eq{
				s.preprocessColumn("mqs.id", ""): store.PreprocessValue(res.ID, ""),
			},
			s.internalMessagebusQueueSettingEncoder(res).Skip("id").Only(onlyColumns...))
		if err != nil {
			return err
		}
	}

	return
}

// UpsertMessagebusQueueSetting updates one or more existing rows in queue_settings
func (s Store) UpsertMessagebusQueueSetting(ctx context.Context, rr ...*messagebus.QueueSettings) (err error) {
	for _, res := range rr {
		err = s.checkMessagebusQueueSettingConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpsertMessagebusQueueSettings(ctx, s.internalMessagebusQueueSettingEncoder(res))
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteMessagebusQueueSetting Deletes one or more rows from queue_settings table
func (s Store) DeleteMessagebusQueueSetting(ctx context.Context, rr ...*messagebus.QueueSettings) (err error) {
	for _, res := range rr {

		err = s.execDeleteMessagebusQueueSettings(ctx, squirrel.Eq{
			s.preprocessColumn("mqs.id", ""): store.PreprocessValue(res.ID, ""),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteMessagebusQueueSettingByID Deletes row from the queue_settings table
func (s Store) DeleteMessagebusQueueSettingByID(ctx context.Context, ID uint64) error {
	return s.execDeleteMessagebusQueueSettings(ctx, squirrel.Eq{
		s.preprocessColumn("mqs.id", ""): store.PreprocessValue(ID, ""),
	})
}

// TruncateMessagebusQueueSettings Deletes all rows from the queue_settings table
func (s Store) TruncateMessagebusQueueSettings(ctx context.Context) error {
	return s.Truncate(ctx, s.messagebusQueueSettingTable())
}

// execLookupMessagebusQueueSetting prepares MessagebusQueueSetting query and executes it,
// returning messagebus.QueueSettings (or error)
func (s Store) execLookupMessagebusQueueSetting(ctx context.Context, cnd squirrel.Sqlizer) (res *messagebus.QueueSettings, err error) {
	var (
		row rowScanner
	)

	row, err = s.QueryRow(ctx, s.messagebusQueueSettingsSelectBuilder().Where(cnd))
	if err != nil {
		return
	}

	res, err = s.internalMessagebusQueueSettingRowScanner(row)
	if err != nil {
		return
	}

	return res, nil
}

// execCreateMessagebusQueueSettings updates all matched (by cnd) rows in queue_settings with given data
func (s Store) execCreateMessagebusQueueSettings(ctx context.Context, payload store.Payload) error {
	return s.Exec(ctx, s.InsertBuilder(s.messagebusQueueSettingTable()).SetMap(payload))
}

// execUpdateMessagebusQueueSettings updates all matched (by cnd) rows in queue_settings with given data
func (s Store) execUpdateMessagebusQueueSettings(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return s.Exec(ctx, s.UpdateBuilder(s.messagebusQueueSettingTable("mqs")).Where(cnd).SetMap(set))
}

// execUpsertMessagebusQueueSettings inserts new or updates matching (by-primary-key) rows in queue_settings with given data
func (s Store) execUpsertMessagebusQueueSettings(ctx context.Context, set store.Payload) error {
	upsert, err := s.config.UpsertBuilder(
		s.config,
		s.messagebusQueueSettingTable(),
		set,
		s.preprocessColumn("id", ""),
	)

	if err != nil {
		return err
	}

	return s.Exec(ctx, upsert)
}

// execDeleteMessagebusQueueSettings Deletes all matched (by cnd) rows in queue_settings with given data
func (s Store) execDeleteMessagebusQueueSettings(ctx context.Context, cnd squirrel.Sqlizer) error {
	return s.Exec(ctx, s.DeleteBuilder(s.messagebusQueueSettingTable("mqs")).Where(cnd))
}

func (s Store) internalMessagebusQueueSettingRowScanner(row rowScanner) (res *messagebus.QueueSettings, err error) {
	res = &messagebus.QueueSettings{}

	if _, has := s.config.RowScanners["messagebusQueueSetting"]; has {
		scanner := s.config.RowScanners["messagebusQueueSetting"].(func(_ rowScanner, _ *messagebus.QueueSettings) error)
		err = scanner(row, res)
	} else {
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
	}

	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err != nil {
		return nil, errors.Store("could not scan messagebusQueueSetting db row: %s", err).Wrap(err)
	} else {
		return res, nil
	}
}

// QueryMessagebusQueueSettings returns squirrel.SelectBuilder with set table and all columns
func (s Store) messagebusQueueSettingsSelectBuilder() squirrel.SelectBuilder {
	return s.SelectBuilder(s.messagebusQueueSettingTable("mqs"), s.messagebusQueueSettingColumns("mqs")...)
}

// messagebusQueueSettingTable name of the db table
func (Store) messagebusQueueSettingTable(aa ...string) string {
	var alias string
	if len(aa) > 0 {
		alias = " AS " + aa[0]
	}

	return "queue_settings" + alias
}

// MessagebusQueueSettingColumns returns all defined table columns
//
// With optional string arg, all columns are returned aliased
func (Store) messagebusQueueSettingColumns(aa ...string) []string {
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

// internalMessagebusQueueSettingEncoder encodes fields from messagebus.QueueSettings to store.Payload (map)
//
// Encoding is done by using generic approach or by calling encodeMessagebusQueueSetting
// func when rdbms.customEncoder=true
func (s Store) internalMessagebusQueueSettingEncoder(res *messagebus.QueueSettings) store.Payload {
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

// collectMessagebusQueueSettingCursorValues collects values from the given resource that and sets them to the cursor
// to be used for pagination
//
// Values that are collected must come from sortable, unique or primary columns/fields
// At least one of the collected columns must be flagged as unique, otherwise fn appends primary keys at the end
//
// Known issue:
//   when collecting cursor values for query that sorts by unique column with partial index (ie: unique handle on
//   undeleted items)
func (s Store) collectMessagebusQueueSettingCursorValues(res *messagebus.QueueSettings, cc ...*filter.SortExpr) *filter.PagingCursor {
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

// checkMessagebusQueueSettingConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant but unfortunately we cannot rely
// on the full support (MySQL does not support conditional indexes)
func (s *Store) checkMessagebusQueueSettingConstraints(ctx context.Context, res *messagebus.QueueSettings) error {
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
