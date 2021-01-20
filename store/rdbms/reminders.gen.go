package rdbms

// This file is an auto-generated file
//
// Template:    pkg/codegen/assets/store_rdbms.gen.go.tpl
// Definitions: store/reminders.yaml
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

// SearchReminders returns all matching rows
//
// This function calls convertReminderFilter with the given
// types.ReminderFilter and expects to receive a working squirrel.SelectBuilder
func (s Store) SearchReminders(ctx context.Context, f types.ReminderFilter) (types.ReminderSet, types.ReminderFilter, error) {
	var (
		err error
		set []*types.Reminder
		q   squirrel.SelectBuilder
	)

	return set, f, func() error {
		q, err = s.convertReminderFilter(f)
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
		if q, err = setOrderBy(q, sort, s.sortableReminderColumns()); err != nil {
			return err
		}

		set, f.PrevPage, f.NextPage, err = s.fetchFullPageOfReminders(
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

// fetchFullPageOfReminders collects all requested results.
//
// Function applies:
//  - cursor conditions (where ...)
//  - limit
//
// Main responsibility of this function is to perform additional sequential queries in case when not enough results
// are collected due to failed check on a specific row (by check fn).
//
// Function then moves cursor to the last item fetched
func (s Store) fetchFullPageOfReminders(
	ctx context.Context,
	q squirrel.SelectBuilder,
	sort filter.SortExprSet,
	cursor *filter.PagingCursor,
	reqItems uint,
	check func(*types.Reminder) (bool, error),
	cursorCond func(*filter.PagingCursor) squirrel.Sqlizer,
) (set []*types.Reminder, prev, next *filter.PagingCursor, err error) {
	var (
		aux []*types.Reminder

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

	set = make([]*types.Reminder, 0, DefaultSliceCapacity)

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

		if aux, err = s.QueryReminders(ctx, tryQuery, check); err != nil {
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
			cursor = s.collectReminderCursorValues(set[collected-1], sort...)

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
		prev = s.collectReminderCursorValues(set[0], sort...)
		prev.ROrder = true
		prev.LThen = !sort.Reversed()
	}

	if hasNext {
		next = s.collectReminderCursorValues(set[collected-1], sort...)
		next.LThen = sort.Reversed()
	}

	return set, prev, next, nil
}

// QueryReminders queries the database, converts and checks each row and
// returns collected set
//
// Fn also returns total number of fetched items and last fetched item so that the caller can construct cursor
// for next page of results
func (s Store) QueryReminders(
	ctx context.Context,
	q squirrel.Sqlizer,
	check func(*types.Reminder) (bool, error),
) ([]*types.Reminder, error) {
	var (
		set = make([]*types.Reminder, 0, DefaultSliceCapacity)
		res *types.Reminder

		// Query rows with
		rows, err = s.Query(ctx, q)
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		if err = rows.Err(); err == nil {
			res, err = s.internalReminderRowScanner(rows)
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

// LookupReminderByID searches for reminder by its ID
//
// It returns reminder even if deleted or suspended
func (s Store) LookupReminderByID(ctx context.Context, id uint64) (*types.Reminder, error) {
	return s.execLookupReminder(ctx, squirrel.Eq{
		s.preprocessColumn("rmd.id", ""): store.PreprocessValue(id, ""),
	})
}

// CreateReminder creates one or more rows in reminders table
func (s Store) CreateReminder(ctx context.Context, rr ...*types.Reminder) (err error) {
	for _, res := range rr {
		err = s.checkReminderConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execCreateReminders(ctx, s.internalReminderEncoder(res))
		if err != nil {
			return err
		}
	}

	return
}

// UpdateReminder updates one or more existing rows in reminders
func (s Store) UpdateReminder(ctx context.Context, rr ...*types.Reminder) error {
	return s.partialReminderUpdate(ctx, nil, rr...)
}

// partialReminderUpdate updates one or more existing rows in reminders
func (s Store) partialReminderUpdate(ctx context.Context, onlyColumns []string, rr ...*types.Reminder) (err error) {
	for _, res := range rr {
		err = s.checkReminderConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpdateReminders(
			ctx,
			squirrel.Eq{
				s.preprocessColumn("rmd.id", ""): store.PreprocessValue(res.ID, ""),
			},
			s.internalReminderEncoder(res).Skip("id").Only(onlyColumns...))
		if err != nil {
			return err
		}
	}

	return
}

// UpsertReminder updates one or more existing rows in reminders
func (s Store) UpsertReminder(ctx context.Context, rr ...*types.Reminder) (err error) {
	for _, res := range rr {
		err = s.checkReminderConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpsertReminders(ctx, s.internalReminderEncoder(res))
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteReminder Deletes one or more rows from reminders table
func (s Store) DeleteReminder(ctx context.Context, rr ...*types.Reminder) (err error) {
	for _, res := range rr {

		err = s.execDeleteReminders(ctx, squirrel.Eq{
			s.preprocessColumn("rmd.id", ""): store.PreprocessValue(res.ID, ""),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteReminderByID Deletes row from the reminders table
func (s Store) DeleteReminderByID(ctx context.Context, ID uint64) error {
	return s.execDeleteReminders(ctx, squirrel.Eq{
		s.preprocessColumn("rmd.id", ""): store.PreprocessValue(ID, ""),
	})
}

// TruncateReminders Deletes all rows from the reminders table
func (s Store) TruncateReminders(ctx context.Context) error {
	return s.Truncate(ctx, s.reminderTable())
}

// execLookupReminder prepares Reminder query and executes it,
// returning types.Reminder (or error)
func (s Store) execLookupReminder(ctx context.Context, cnd squirrel.Sqlizer) (res *types.Reminder, err error) {
	var (
		row rowScanner
	)

	row, err = s.QueryRow(ctx, s.remindersSelectBuilder().Where(cnd))
	if err != nil {
		return
	}

	res, err = s.internalReminderRowScanner(row)
	if err != nil {
		return
	}

	return res, nil
}

// execCreateReminders updates all matched (by cnd) rows in reminders with given data
func (s Store) execCreateReminders(ctx context.Context, payload store.Payload) error {
	return s.Exec(ctx, s.InsertBuilder(s.reminderTable()).SetMap(payload))
}

// execUpdateReminders updates all matched (by cnd) rows in reminders with given data
func (s Store) execUpdateReminders(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return s.Exec(ctx, s.UpdateBuilder(s.reminderTable("rmd")).Where(cnd).SetMap(set))
}

// execUpsertReminders inserts new or updates matching (by-primary-key) rows in reminders with given data
func (s Store) execUpsertReminders(ctx context.Context, set store.Payload) error {
	upsert, err := s.config.UpsertBuilder(
		s.config,
		s.reminderTable(),
		set,
		s.preprocessColumn("id", ""),
	)

	if err != nil {
		return err
	}

	return s.Exec(ctx, upsert)
}

// execDeleteReminders Deletes all matched (by cnd) rows in reminders with given data
func (s Store) execDeleteReminders(ctx context.Context, cnd squirrel.Sqlizer) error {
	return s.Exec(ctx, s.DeleteBuilder(s.reminderTable("rmd")).Where(cnd))
}

func (s Store) internalReminderRowScanner(row rowScanner) (res *types.Reminder, err error) {
	res = &types.Reminder{}

	if _, has := s.config.RowScanners["reminder"]; has {
		scanner := s.config.RowScanners["reminder"].(func(_ rowScanner, _ *types.Reminder) error)
		err = scanner(row, res)
	} else {
		err = row.Scan(
			&res.ID,
			&res.Resource,
			&res.Payload,
			&res.SnoozeCount,
			&res.AssignedTo,
			&res.AssignedBy,
			&res.AssignedAt,
			&res.DismissedBy,
			&res.DismissedAt,
			&res.RemindAt,
			&res.CreatedAt,
			&res.UpdatedAt,
			&res.DeletedAt,
		)
	}

	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err != nil {
		return nil, errors.Store("could not scan reminder db row: %s", err).Wrap(err)
	} else {
		return res, nil
	}
}

// QueryReminders returns squirrel.SelectBuilder with set table and all columns
func (s Store) remindersSelectBuilder() squirrel.SelectBuilder {
	return s.SelectBuilder(s.reminderTable("rmd"), s.reminderColumns("rmd")...)
}

// reminderTable name of the db table
func (Store) reminderTable(aa ...string) string {
	var alias string
	if len(aa) > 0 {
		alias = " AS " + aa[0]
	}

	return "reminders" + alias
}

// ReminderColumns returns all defined table columns
//
// With optional string arg, all columns are returned aliased
func (Store) reminderColumns(aa ...string) []string {
	var alias string
	if len(aa) > 0 {
		alias = aa[0] + "."
	}

	return []string{
		alias + "id",
		alias + "resource",
		alias + "payload",
		alias + "snooze_count",
		alias + "assigned_to",
		alias + "assigned_by",
		alias + "assigned_at",
		alias + "dismissed_by",
		alias + "dismissed_at",
		alias + "remind_at",
		alias + "created_at",
		alias + "updated_at",
		alias + "deleted_at",
	}
}

// {true true false true true true}

// sortableReminderColumns returns all Reminder columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
func (Store) sortableReminderColumns() map[string]string {
	return map[string]string{
		"id": "id", "remind_at": "remind_at",
		"remindat":   "remind_at",
		"created_at": "created_at",
		"createdat":  "created_at",
		"updated_at": "updated_at",
		"updatedat":  "updated_at",
		"deleted_at": "deleted_at",
		"deletedat":  "deleted_at",
	}
}

// internalReminderEncoder encodes fields from types.Reminder to store.Payload (map)
//
// Encoding is done by using generic approach or by calling encodeReminder
// func when rdbms.customEncoder=true
func (s Store) internalReminderEncoder(res *types.Reminder) store.Payload {
	return store.Payload{
		"id":           res.ID,
		"resource":     res.Resource,
		"payload":      res.Payload,
		"snooze_count": res.SnoozeCount,
		"assigned_to":  res.AssignedTo,
		"assigned_by":  res.AssignedBy,
		"assigned_at":  res.AssignedAt,
		"dismissed_by": res.DismissedBy,
		"dismissed_at": res.DismissedAt,
		"remind_at":    res.RemindAt,
		"created_at":   res.CreatedAt,
		"updated_at":   res.UpdatedAt,
		"deleted_at":   res.DeletedAt,
	}
}

// collectReminderCursorValues collects values from the given resource that and sets them to the cursor
// to be used for pagination
//
// Values that are collected must come from sortable, unique or primary columns/fields
// At least one of the collected columns must be flagged as unique, otherwise fn appends primary keys at the end
//
// Known issue:
//   when collecting cursor values for query that sorts by unique column with partial index (ie: unique handle on
//   undeleted items)
func (s Store) collectReminderCursorValues(res *types.Reminder, cc ...*filter.SortExpr) *filter.PagingCursor {
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
				case "remind_at":
					cursor.Set(c.Column, res.RemindAt, c.Descending)

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

// checkReminderConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant but unfortunately we can not rely
// on the full support (MySQL does not support conditional indexes)
func (s *Store) checkReminderConstraints(ctx context.Context, res *types.Reminder) error {
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
