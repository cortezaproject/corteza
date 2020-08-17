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
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/jmoiron/sqlx"
)

// SearchReminders returns all matching rows
//
// This function calls convertReminderFilter with the given
// types.ReminderFilter and expects to receive a working squirrel.SelectBuilder
func (s Store) SearchReminders(ctx context.Context, f types.ReminderFilter) (types.ReminderSet, types.ReminderFilter, error) {
	q, err := s.convertReminderFilter(f)
	if err != nil {
		return nil, f, err
	}

	scap := f.PerPage
	if scap == 0 {
		scap = DefaultSliceCapacity
	}

	if f.Count, err = Count(ctx, s.db, q); err != nil || f.Count == 0 {
		return nil, f, err
	}

	var (
		set = make([]*types.Reminder, 0, scap)
		// @todo this offset needs to be removed and replaced with key-based-paging
		fetchPage = func(offset, limit uint) (fetched, skipped uint, err error) {
			var (
				res *types.Reminder
				chk bool
			)

			if limit > 0 {
				q = q.Limit(uint64(limit))
			}

			if offset > 0 {
				q = q.Offset(uint64(offset))
			}

			rows, err := s.Query(ctx, q)
			if err != nil {
				return
			}

			for rows.Next() {
				fetched++
				if res, err = s.internalReminderRowScanner(rows, rows.Err()); err != nil {
					if cerr := rows.Close(); cerr != nil {
						err = fmt.Errorf("could not close rows (%v) after scan error: %w", cerr, err)
					}

					return
				}

				// If check function is set, call it and act accordingly
				if f.Check != nil {
					if chk, err = f.Check(res); err != nil {
						if cerr := rows.Close(); cerr != nil {
							err = fmt.Errorf("could not close rows (%v) after check error: %w", cerr, err)
						}

						return
					} else if !chk {
						// did not pass the check
						// go with the next row
						skipped++
						continue
					}
				}

				set = append(set, res)

				// make sure we do not fetch more than requested!
				if f.Limit > 0 && uint(len(set)) >= f.Limit {
					break
				}
			}

			err = rows.Close()
			return
		}

		fetch = func() error {
			var (
				fetched uint

				// starting offset & limit are from filter arg
				// note that this will have to be improved with key-based pagination
				offset, limit = calculatePaging(f.PageFilter)
			)

			for refetch := 0; refetch < MaxRefetches; refetch++ {
				if fetched, _, err = fetchPage(offset, limit); err != nil {
					return err
				}

				// if limit is not set or we've already collected enough resources
				// we can break the loop right away
				if limit == 0 || fetched == 0 || uint(len(set)) >= f.Limit {
					break
				}

				// we've skipped fetched resources (due to check() fn)
				// and we still have less results (in set) than required by limit
				// inc offset by number of fetched items
				offset += fetched

				if limit < MinRefetchLimit {
					limit = MinRefetchLimit
				}

			}
			return nil
		}
	)

	return set, f, fetch()
}

// LookupReminderByID searches for reminder by its ID
//
// It returns reminder even if deleted or suspended
func (s Store) LookupReminderByID(ctx context.Context, id uint64) (*types.Reminder, error) {
	return s.ReminderLookup(ctx, squirrel.Eq{
		"rmd.id": id,
	})
}

// CreateReminder creates one or more rows in reminders table
func (s Store) CreateReminder(ctx context.Context, rr ...*types.Reminder) error {
	if len(rr) == 0 {
		return nil
	}

	return Tx(ctx, s.db, s.config, nil, func(db *sqlx.Tx) (err error) {
		for _, res := range rr {
			err = ExecuteSqlizer(ctx, s.DB(), s.Insert(s.ReminderTable()).SetMap(s.internalReminderEncoder(res)))
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// UpdateReminder updates one or more existing rows in reminders
func (s Store) UpdateReminder(ctx context.Context, rr ...*types.Reminder) error {
	return s.PartialUpdateReminder(ctx, nil, rr...)
}

// PartialUpdateReminder updates one or more existing rows in reminders
//
// It wraps the update into transaction and can perform partial update by providing list of updatable columns
func (s Store) PartialUpdateReminder(ctx context.Context, onlyColumns []string, rr ...*types.Reminder) error {
	if len(rr) == 0 {
		return nil
	}

	return Tx(ctx, s.db, s.config, nil, func(db *sqlx.Tx) (err error) {
		for _, res := range rr {
			err = s.ExecUpdateReminders(
				ctx,
				squirrel.Eq{s.preprocessColumn("rmd.id", ""): s.preprocessValue(res.ID, "")},
				s.internalReminderEncoder(res).Skip("id").Only(onlyColumns...))
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// RemoveReminder removes one or more rows from reminders table
func (s Store) RemoveReminder(ctx context.Context, rr ...*types.Reminder) error {
	if len(rr) == 0 {
		return nil
	}

	return Tx(ctx, s.db, s.config, nil, func(db *sqlx.Tx) (err error) {
		for _, res := range rr {
			err = ExecuteSqlizer(ctx, s.DB(), s.Delete(s.ReminderTable("rmd")).Where(squirrel.Eq{s.preprocessColumn("rmd.id", ""): s.preprocessValue(res.ID, "")}))
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// RemoveReminderByID removes row from the reminders table
func (s Store) RemoveReminderByID(ctx context.Context, ID uint64) error {
	return ExecuteSqlizer(ctx, s.DB(), s.Delete(s.ReminderTable("rmd")).Where(squirrel.Eq{s.preprocessColumn("rmd.id", ""): s.preprocessValue(ID, "")}))
}

// TruncateReminders removes all rows from the reminders table
func (s Store) TruncateReminders(ctx context.Context) error {
	return Truncate(ctx, s.DB(), s.ReminderTable())
}

// ExecUpdateReminders updates all matched (by cnd) rows in reminders with given data
func (s Store) ExecUpdateReminders(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return ExecuteSqlizer(ctx, s.DB(), s.Update(s.ReminderTable("rmd")).Where(cnd).SetMap(set))
}

// ReminderLookup prepares Reminder query and executes it,
// returning types.Reminder (or error)
func (s Store) ReminderLookup(ctx context.Context, cnd squirrel.Sqlizer) (*types.Reminder, error) {
	return s.internalReminderRowScanner(s.QueryRow(ctx, s.QueryReminders().Where(cnd)))
}

func (s Store) internalReminderRowScanner(row rowScanner, err error) (*types.Reminder, error) {
	if err != nil {
		return nil, err
	}

	var res = &types.Reminder{}
	if _, has := s.config.RowScanners["reminder"]; has {
		scanner := s.config.RowScanners["reminder"].(func(rowScanner, *types.Reminder) error)
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
func (s Store) QueryReminders() squirrel.SelectBuilder {
	return s.Select(s.ReminderTable("rmd"), s.ReminderColumns("rmd")...)
}

// ReminderTable name of the db table
func (Store) ReminderTable(aa ...string) string {
	var alias string
	if len(aa) > 0 {
		alias = " AS " + aa[0]
	}

	return "reminders" + alias
}

// ReminderColumns returns all defined table columns
//
// With optional string arg, all columns are returned aliased
func (Store) ReminderColumns(aa ...string) []string {
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
