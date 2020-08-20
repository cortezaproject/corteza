package rdbms

// This file is an auto-generated file
//
// Template:    pkg/codegen/assets/store_rdbms.gen.go.tpl
// Definitions: store/compose_records.yaml
//
// Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/store"
	"strings"
)

var _ = errors.Is

const (
	TriggerBeforeComposeRecordCreate triggerKey = "composeRecordBeforeCreate"
	TriggerBeforeComposeRecordUpdate triggerKey = "composeRecordBeforeUpdate"
	TriggerBeforeComposeRecordUpsert triggerKey = "composeRecordBeforeUpsert"
	TriggerBeforeComposeRecordDelete triggerKey = "composeRecordBeforeDelete"
)

// searchComposeRecords returns all matching rows
//
// This function calls convertComposeRecordFilter with the given
// types.RecordFilter and expects to receive a working squirrel.SelectBuilder
func (s Store) searchComposeRecords(ctx context.Context, _mod *types.Module, f types.RecordFilter) (types.RecordSet, types.RecordFilter, error) {
	var scap uint
	q, err := s.convertComposeRecordFilter(_mod, f)
	if err != nil {
		return nil, f, err
	}

	scap = f.Limit

	// Cleanup anything we've accidentally received...
	f.PrevPage, f.NextPage = nil, nil

	// When cursor for a previous page is used it's marked as reversed
	// This tells us to flip the descending flag on all used sort keys
	reverseCursor := f.PageCursor != nil && f.PageCursor.Reverse

	if err := f.Sort.Validate(s.sortableComposeRecordColumns()...); err != nil {
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
		set = make([]*types.Record, 0, scap)
		// fetches rows and scans them into types.Record resource this is then passed to Check function on filter
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
				res *types.Record

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
					res, err = s.internalComposeRecordRowScanner(_mod, rows)
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
					f.PrevPage = s.collectComposeRecordCursorValues(set[0], sort.Columns()...)
					f.PrevPage.Reverse = true
				}

				// Less items fetched then requested by page-limit
				// not very likely there's another page
				f.NextPage = s.collectComposeRecordCursorValues(set[len(set)-1], sort.Columns()...)
			}

			f.PageCursor = nil
			return nil
		}
	)

	return set, f, s.config.ErrorHandler(fetch())
}

// lookupComposeRecordByID searches for compose record by ID
// It returns compose record even if deleted
func (s Store) lookupComposeRecordByID(ctx context.Context, _mod *types.Module, id uint64) (*types.Record, error) {
	return s.execLookupComposeRecord(ctx, _mod, squirrel.Eq{
		s.preprocessColumn("crd.id", ""): s.preprocessValue(id, ""),
	})
}

// createComposeRecord creates one or more rows in compose_record table
func (s Store) createComposeRecord(ctx context.Context, _mod *types.Module, rr ...*types.Record) (err error) {
	for _, res := range rr {
		err = s.checkComposeRecordConstraints(ctx, _mod, res)
		if err != nil {
			return err
		}

		// err = s.composeRecordHook(ctx, TriggerBeforeComposeRecordCreate, _mod, res)
		// if err != nil {
		// 	return err
		// }

		err = s.execCreateComposeRecords(ctx, s.internalComposeRecordEncoder(res))
		if err != nil {
			return err
		}
	}

	return
}

// updateComposeRecord updates one or more existing rows in compose_record
func (s Store) updateComposeRecord(ctx context.Context, _mod *types.Module, rr ...*types.Record) error {
	return s.config.ErrorHandler(s.partialComposeRecordUpdate(ctx, _mod, nil, rr...))
}

// partialComposeRecordUpdate updates one or more existing rows in compose_record
func (s Store) partialComposeRecordUpdate(ctx context.Context, _mod *types.Module, onlyColumns []string, rr ...*types.Record) (err error) {
	for _, res := range rr {
		err = s.checkComposeRecordConstraints(ctx, _mod, res)
		if err != nil {
			return err
		}

		// err = s.composeRecordHook(ctx, TriggerBeforeComposeRecordUpdate, _mod, res)
		// if err != nil {
		// 	return err
		// }

		err = s.execUpdateComposeRecords(
			ctx,
			squirrel.Eq{
				s.preprocessColumn("crd.id", ""): s.preprocessValue(res.ID, ""),
			},
			s.internalComposeRecordEncoder(res).Skip("id").Only(onlyColumns...))
		if err != nil {
			return s.config.ErrorHandler(err)
		}
	}

	return
}

// upsertComposeRecord updates one or more existing rows in compose_record
func (s Store) upsertComposeRecord(ctx context.Context, _mod *types.Module, rr ...*types.Record) (err error) {
	for _, res := range rr {
		err = s.checkComposeRecordConstraints(ctx, _mod, res)
		if err != nil {
			return err
		}

		// err = s.composeRecordHook(ctx, TriggerBeforeComposeRecordUpsert, _mod, res)
		// if err != nil {
		// 	return err
		// }

		err = s.config.ErrorHandler(s.execUpsertComposeRecords(ctx, s.internalComposeRecordEncoder(res)))
		if err != nil {
			return err
		}
	}

	return nil
}

// deleteComposeRecord Deletes one or more rows from compose_record table
func (s Store) deleteComposeRecord(ctx context.Context, _mod *types.Module, rr ...*types.Record) (err error) {
	for _, res := range rr {
		// err = s.composeRecordHook(ctx, TriggerBeforeComposeRecordDelete, _mod, res)
		// if err != nil {
		// 	return err
		// }

		err = s.execDeleteComposeRecords(ctx, squirrel.Eq{
			s.preprocessColumn("crd.id", ""): s.preprocessValue(res.ID, ""),
		})
		if err != nil {
			return s.config.ErrorHandler(err)
		}
	}

	return nil
}

// deleteComposeRecordByID Deletes row from the compose_record table
func (s Store) deleteComposeRecordByID(ctx context.Context, _mod *types.Module, ID uint64) error {
	return s.execDeleteComposeRecords(ctx, squirrel.Eq{
		s.preprocessColumn("crd.id", ""): s.preprocessValue(ID, ""),
	})
}

// truncateComposeRecords Deletes all rows from the compose_record table
func (s Store) truncateComposeRecords(ctx context.Context, _mod *types.Module) error {
	return s.config.ErrorHandler(s.Truncate(ctx, s.composeRecordTable()))
}

// execLookupComposeRecord prepares ComposeRecord query and executes it,
// returning types.Record (or error)
func (s Store) execLookupComposeRecord(ctx context.Context, _mod *types.Module, cnd squirrel.Sqlizer) (res *types.Record, err error) {
	var (
		row rowScanner
	)

	row, err = s.QueryRow(ctx, s.composeRecordsSelectBuilder().Where(cnd))
	if err != nil {
		return
	}

	res, err = s.internalComposeRecordRowScanner(_mod, row)
	if err != nil {
		return
	}

	return res, nil
}

// execCreateComposeRecords updates all matched (by cnd) rows in compose_record with given data
func (s Store) execCreateComposeRecords(ctx context.Context, payload store.Payload) error {
	return s.config.ErrorHandler(s.Exec(ctx, s.InsertBuilder(s.composeRecordTable()).SetMap(payload)))
}

// execUpdateComposeRecords updates all matched (by cnd) rows in compose_record with given data
func (s Store) execUpdateComposeRecords(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return s.config.ErrorHandler(s.Exec(ctx, s.UpdateBuilder(s.composeRecordTable("crd")).Where(cnd).SetMap(set)))
}

// execUpsertComposeRecords inserts new or updates matching (by-primary-key) rows in compose_record with given data
func (s Store) execUpsertComposeRecords(ctx context.Context, set store.Payload) error {
	upsert, err := s.config.UpsertBuilder(
		s.config,
		s.composeRecordTable(),
		set,
		"id",
	)

	if err != nil {
		return err
	}

	return s.config.ErrorHandler(s.Exec(ctx, upsert))
}

// execDeleteComposeRecords Deletes all matched (by cnd) rows in compose_record with given data
func (s Store) execDeleteComposeRecords(ctx context.Context, cnd squirrel.Sqlizer) error {
	return s.config.ErrorHandler(s.Exec(ctx, s.DeleteBuilder(s.composeRecordTable("crd")).Where(cnd)))
}

func (s Store) internalComposeRecordRowScanner(_mod *types.Module, row rowScanner) (res *types.Record, err error) {
	res = &types.Record{}

	if _, has := s.config.RowScanners["composeRecord"]; has {
		scanner := s.config.RowScanners["composeRecord"].(func(_mod *types.Module, _ rowScanner, _ *types.Record) error)
		err = scanner(_mod, row, res)
	} else {
		err = row.Scan(
			&res.ID,
			&res.ModuleID,
			&res.NamespaceID,
			&res.OwnedBy,
			&res.CreatedBy,
			&res.UpdatedBy,
			&res.DeletedBy,
			&res.CreatedAt,
			&res.UpdatedAt,
			&res.DeletedAt,
		)
	}

	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("could not scan db row for ComposeRecord: %w", err)
	} else {
		return res, nil
	}
}

// QueryComposeRecords returns squirrel.SelectBuilder with set table and all columns
func (s Store) composeRecordsSelectBuilder() squirrel.SelectBuilder {
	return s.SelectBuilder(s.composeRecordTable("crd"), s.composeRecordColumns("crd")...)
}

// composeRecordTable name of the db table
func (Store) composeRecordTable(aa ...string) string {
	var alias string
	if len(aa) > 0 {
		alias = " AS " + aa[0]
	}

	return "compose_record" + alias
}

// ComposeRecordColumns returns all defined table columns
//
// With optional string arg, all columns are returned aliased
func (Store) composeRecordColumns(aa ...string) []string {
	var alias string
	if len(aa) > 0 {
		alias = aa[0] + "."
	}

	return []string{
		alias + "id",
		alias + "module_id",
		alias + "rel_namespace",
		alias + "owned_by",
		alias + "created_by",
		alias + "updated_by",
		alias + "deleted_by",
		alias + "created_at",
		alias + "updated_at",
		alias + "deleted_at",
	}
}

// {true false true true true}

// sortableComposeRecordColumns returns all ComposeRecord columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
func (Store) sortableComposeRecordColumns() []string {
	return []string{
		"id",
		"created_at",
		"updated_at",
		"deleted_at",
	}
}

// internalComposeRecordEncoder encodes fields from types.Record to store.Payload (map)
//
// Encoding is done by using generic approach or by calling encodeComposeRecord
// func when rdbms.customEncoder=true
func (s Store) internalComposeRecordEncoder(res *types.Record) store.Payload {
	return store.Payload{
		"id":            res.ID,
		"module_id":     res.ModuleID,
		"rel_namespace": res.NamespaceID,
		"owned_by":      res.OwnedBy,
		"created_by":    res.CreatedBy,
		"updated_by":    res.UpdatedBy,
		"deleted_by":    res.DeletedBy,
		"created_at":    res.CreatedAt,
		"updated_at":    res.UpdatedAt,
		"deleted_at":    res.DeletedAt,
	}
}

func (s Store) collectComposeRecordCursorValues(res *types.Record, cc ...string) *filter.PagingCursor {
	var (
		cursor = &filter.PagingCursor{}

		hasUnique bool

		collect = func(cc ...string) {
			for _, c := range cc {
				switch c {
				case "id":
					cursor.Set(c, res.ID, false)
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

func (s *Store) checkComposeRecordConstraints(ctx context.Context, _mod *types.Module, res *types.Record) error {

	return nil
}

// func (s *Store) composeRecordHook(ctx context.Context, key triggerKey, _mod  *types.Module, res *types.Record) error {
// 	if fn, has := s.config.TriggerHandlers[key]; has {
// 		return fn.(func (ctx context.Context, s *Store, _mod  *types.Module, res *types.Record) error)(ctx, s, _mod, res)
// 	}
//
// 	return nil
// }
