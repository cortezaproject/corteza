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
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
	"strings"
)

var _ = errors.Is

const (
	TriggerBeforeReminderCreate triggerKey = "reminderBeforeCreate"
	TriggerBeforeReminderUpdate triggerKey = "reminderBeforeUpdate"
	TriggerBeforeReminderUpsert triggerKey = "reminderBeforeUpsert"
	TriggerBeforeReminderDelete triggerKey = "reminderBeforeDelete"
)

// SearchReminders returns all matching rows
//
// This function calls convertReminderFilter with the given
// types.ReminderFilter and expects to receive a working squirrel.SelectBuilder
func (s Store) SearchReminders(ctx context.Context, f types.ReminderFilter) (types.ReminderSet, types.ReminderFilter, error) {
	var scap uint
	q, err := s.convertReminderFilter(f)
	if err != nil {
		return nil, f, err
	}

	scap = f.Limit

	// Cleanup anything we've accidentally received...
	f.PrevPage, f.NextPage = nil, nil

	// When cursor for a previous page is used it's marked as reversed
	// This tells us to flip the descending flag on all used sort keys
	reverseCursor := f.PageCursor != nil && f.PageCursor.Reverse

	if err := f.Sort.Validate(s.sortableReminderColumns()...); err != nil {
		return nil, f, fmt.Errorf("could not validate sort: %v", err)
	}

	// If paging with reverse cursor, change the sorting
	// direction for all columns we're sorting by
	sort := f.Sort.Clone()
	if reverseCursor {
		sort.Reverse()
	}

	// Apply sorting expr from filter to query
	if len(sort) > 0 {
		sqlSort := make([]string, len(sort))
		for i := range sort {
			sqlSort[i] = sort[i].Column
			if sort[i].Descending {
				sqlSort[i] += " DESC"
			}
		}

		q = q.OrderBy(sqlSort...)
	}

	if scap == 0 {
		scap = DefaultSliceCapacity
	}

	var (
		set = make([]*types.Reminder, 0, scap)
		// fetches rows and scans them into types.Reminder resource this is then passed to Check function on filter
		// to help determine if fetched resource fits or not
		//
		// Note that limit is passed explicitly and is not necessarily equal to filter's limit. We want
		// to keep that value intact.
		//
		// The value for cursor is used and set directly from/to the filter!
		//
		// It returns total number of fetched pages and modifies PageCursor value for paging
		fetchPage = func(cursor *filter.PagingCursor, limit uint) (fetched uint, err error) {
			var (
				res *types.Reminder

				// Make a copy of the select query builder so that we don't change
				// the original query
				slct = q.Options()
			)

			if limit > 0 {
				slct = slct.Limit(uint64(limit))

				if cursor != nil && len(cursor.Keys()) > 0 {
					const cursorTpl = `(%s) %s (?%s)`
					op := ">"
					if cursor.Reverse {
						op = "<"
					}

					pred := fmt.Sprintf(cursorTpl, strings.Join(cursor.Keys(), ", "), op, strings.Repeat(", ?", len(cursor.Keys())-1))
					slct = slct.Where(pred, cursor.Values()...)
				}
			}

			rows, err := s.Query(ctx, slct)
			if err != nil {
				return
			}

			for rows.Next() {
				fetched++

				if rows.Err() == nil {
					res, err = s.internalReminderRowScanner(rows)
				}

				if err != nil {
					if cerr := rows.Close(); cerr != nil {
						err = fmt.Errorf("could not close rows (%v) after scan error: %w", cerr, err)
					}

					return
				}

				// If check function is set, call it and act accordingly

				if f.Check != nil {
					var chk bool
					if chk, err = f.Check(res); err != nil {
						if cerr := rows.Close(); cerr != nil {
							err = fmt.Errorf("could not close rows (%v) after check error: %w", cerr, err)
						}

						return
					} else if !chk {
						// did not pass the check
						// go with the next row
						continue
					}
				}
				set = append(set, res)

				if f.Limit > 0 {
					if uint(len(set)) >= f.Limit {
						// make sure we do not fetch more than requested!
						break
					}
				}
			}

			err = rows.Close()
			return
		}

		fetch = func() error {
			var (
				// how many items were actually fetched
				fetched uint

				// starting offset & limit are from filter arg
				// note that this will have to be improved with key-based pagination
				limit = f.Limit

				// Copy cursor value
				//
				// This is where we'll start fetching and this value will be overwritten when
				// results come back
				cursor = f.PageCursor

				lastSetFull bool
			)

			for refetch := 0; refetch < MaxRefetches; refetch++ {
				if fetched, err = fetchPage(cursor, limit); err != nil {
					return err
				}

				// if limit is not set or we've already collected enough items
				// we can break the loop right away
				if limit == 0 || fetched == 0 || fetched < limit {
					break
				}

				if uint(len(set)) >= f.Limit {
					// we should return as much as requested
					set = set[0:f.Limit]
					lastSetFull = true
					break
				}

				// In case limit is set very low and we've missed records in the first fetch,
				// make sure next fetch limit is a bit higher
				if limit < MinRefetchLimit {
					limit = MinRefetchLimit
				}

				// @todo it might be good to implement different kind of strategies
				//       (beyond min-refetch-limit above) that can adjust limit on
				//       retry to more optimal number
			}

			if reverseCursor {
				// Cursor for previous page was used
				// Fetched set needs to be reverseCursor because we've forced a descending order to
				// get the previus page
				for i, j := 0, len(set)-1; i < j; i, j = i+1, j-1 {
					set[i], set[j] = set[j], set[i]
				}
			}

			if f.Limit > 0 && len(set) > 0 {
				if f.PageCursor != nil && (!f.PageCursor.Reverse || lastSetFull) {
					f.PrevPage = s.collectReminderCursorValues(set[0], sort.Columns()...)
					f.PrevPage.Reverse = true
				}

				// Less items fetched then requested by page-limit
				// not very likely there's another page
				f.NextPage = s.collectReminderCursorValues(set[len(set)-1], sort.Columns()...)
			}

			f.PageCursor = nil
			return nil
		}
	)

	return set, f, s.config.ErrorHandler(fetch())
}

// LookupReminderByID searches for reminder by its ID
//
// It returns reminder even if deleted or suspended
func (s Store) LookupReminderByID(ctx context.Context, id uint64) (*types.Reminder, error) {
	return s.execLookupReminder(ctx, squirrel.Eq{
		s.preprocessColumn("rmd.id", ""): s.preprocessValue(id, ""),
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
	return s.config.ErrorHandler(s.PartialReminderUpdate(ctx, nil, rr...))
}

// PartialReminderUpdate updates one or more existing rows in reminders
func (s Store) PartialReminderUpdate(ctx context.Context, onlyColumns []string, rr ...*types.Reminder) (err error) {
	for _, res := range rr {
		err = s.checkReminderConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpdateReminders(
			ctx,
			squirrel.Eq{
				s.preprocessColumn("rmd.id", ""): s.preprocessValue(res.ID, ""),
			},
			s.internalReminderEncoder(res).Skip("id").Only(onlyColumns...))
		if err != nil {
			return s.config.ErrorHandler(err)
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

		err = s.config.ErrorHandler(s.execUpsertReminders(ctx, s.internalReminderEncoder(res)))
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
			s.preprocessColumn("rmd.id", ""): s.preprocessValue(res.ID, ""),
		})
		if err != nil {
			return s.config.ErrorHandler(err)
		}
	}

	return nil
}

// DeleteReminderByID Deletes row from the reminders table
func (s Store) DeleteReminderByID(ctx context.Context, ID uint64) error {
	return s.execDeleteReminders(ctx, squirrel.Eq{
		s.preprocessColumn("rmd.id", ""): s.preprocessValue(ID, ""),
	})
}

// TruncateReminders Deletes all rows from the reminders table
func (s Store) TruncateReminders(ctx context.Context) error {
	return s.config.ErrorHandler(s.Truncate(ctx, s.reminderTable()))
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
	return s.config.ErrorHandler(s.Exec(ctx, s.InsertBuilder(s.reminderTable()).SetMap(payload)))
}

// execUpdateReminders updates all matched (by cnd) rows in reminders with given data
func (s Store) execUpdateReminders(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return s.config.ErrorHandler(s.Exec(ctx, s.UpdateBuilder(s.reminderTable("rmd")).Where(cnd).SetMap(set)))
}

// execUpsertReminders inserts new or updates matching (by-primary-key) rows in reminders with given data
func (s Store) execUpsertReminders(ctx context.Context, set store.Payload) error {
	upsert, err := s.config.UpsertBuilder(
		s.config,
		s.reminderTable(),
		set,
		"id",
	)

	if err != nil {
		return err
	}

	return s.config.ErrorHandler(s.Exec(ctx, upsert))
}

// execDeleteReminders Deletes all matched (by cnd) rows in reminders with given data
func (s Store) execDeleteReminders(ctx context.Context, cnd squirrel.Sqlizer) error {
	return s.config.ErrorHandler(s.Exec(ctx, s.DeleteBuilder(s.reminderTable("rmd")).Where(cnd)))
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
		return nil, store.ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("could not scan db row for Reminder: %w", err)
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

// {true true true true true}

// sortableReminderColumns returns all Reminder columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
func (Store) sortableReminderColumns() []string {
	return []string{
		"id",
		"remind_at",
		"created_at",
		"updated_at",
		"deleted_at",
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

func (s Store) collectReminderCursorValues(res *types.Reminder, cc ...string) *filter.PagingCursor {
	var (
		cursor = &filter.PagingCursor{}

		hasUnique bool

		collect = func(cc ...string) {
			for _, c := range cc {
				switch c {
				case "id":
					cursor.Set(c, res.ID, false)
				case "remind_at":
					cursor.Set(c, res.RemindAt, false)
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
	if !hasUnique {
		collect(
			"id",
		)
	}

	return cursor
}

func (s *Store) checkReminderConstraints(ctx context.Context, res *types.Reminder) error {

	return nil
}
