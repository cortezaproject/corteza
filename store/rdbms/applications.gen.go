package rdbms

// This file is an auto-generated file
//
// Template:    pkg/codegen/assets/store_rdbms.gen.go.tpl
// Definitions: store/applications.yaml
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
	TriggerBeforeApplicationCreate triggerKey = "applicationBeforeCreate"
	TriggerBeforeApplicationUpdate triggerKey = "applicationBeforeUpdate"
	TriggerBeforeApplicationUpsert triggerKey = "applicationBeforeUpsert"
	TriggerBeforeApplicationDelete triggerKey = "applicationBeforeDelete"
)

// SearchApplications returns all matching rows
//
// This function calls convertApplicationFilter with the given
// types.ApplicationFilter and expects to receive a working squirrel.SelectBuilder
func (s Store) SearchApplications(ctx context.Context, f types.ApplicationFilter) (types.ApplicationSet, types.ApplicationFilter, error) {
	var scap uint
	q, err := s.convertApplicationFilter(f)
	if err != nil {
		return nil, f, err
	}

	scap = f.Limit

	// Cleanup anything we've accidentally received...
	f.PrevPage, f.NextPage = nil, nil

	// When cursor for a previous page is used it's marked as reversed
	// This tells us to flip the descending flag on all used sort keys
	reverseCursor := f.PageCursor != nil && f.PageCursor.Reverse

	if err := f.Sort.Validate(s.sortableApplicationColumns()...); err != nil {
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
		set = make([]*types.Application, 0, scap)
		// fetches rows and scans them into types.Application resource this is then passed to Check function on filter
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
				res *types.Application

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
					res, err = s.internalApplicationRowScanner(rows)
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
					f.PrevPage = s.collectApplicationCursorValues(set[0], sort.Columns()...)
					f.PrevPage.Reverse = true
				}

				// Less items fetched then requested by page-limit
				// not very likely there's another page
				f.NextPage = s.collectApplicationCursorValues(set[len(set)-1], sort.Columns()...)
			}

			f.PageCursor = nil
			return nil
		}
	)

	return set, f, s.config.ErrorHandler(fetch())
}

// LookupApplicationByID searches for application by ID
//
// It returns application even if deleted
func (s Store) LookupApplicationByID(ctx context.Context, id uint64) (*types.Application, error) {
	return s.execLookupApplication(ctx, squirrel.Eq{
		s.preprocessColumn("app.id", ""): s.preprocessValue(id, ""),
	})
}

// CreateApplication creates one or more rows in applications table
func (s Store) CreateApplication(ctx context.Context, rr ...*types.Application) (err error) {
	for _, res := range rr {
		err = s.checkApplicationConstraints(ctx, res)
		if err != nil {
			return err
		}

		// err = s.applicationHook(ctx, TriggerBeforeApplicationCreate, res)
		// if err != nil {
		// 	return err
		// }

		err = s.execCreateApplications(ctx, s.internalApplicationEncoder(res))
		if err != nil {
			return err
		}
	}

	return
}

// UpdateApplication updates one or more existing rows in applications
func (s Store) UpdateApplication(ctx context.Context, rr ...*types.Application) error {
	return s.config.ErrorHandler(s.PartialApplicationUpdate(ctx, nil, rr...))
}

// PartialApplicationUpdate updates one or more existing rows in applications
func (s Store) PartialApplicationUpdate(ctx context.Context, onlyColumns []string, rr ...*types.Application) (err error) {
	for _, res := range rr {
		err = s.checkApplicationConstraints(ctx, res)
		if err != nil {
			return err
		}

		// err = s.applicationHook(ctx, TriggerBeforeApplicationUpdate, res)
		// if err != nil {
		// 	return err
		// }

		err = s.execUpdateApplications(
			ctx,
			squirrel.Eq{
				s.preprocessColumn("app.id", ""): s.preprocessValue(res.ID, ""),
			},
			s.internalApplicationEncoder(res).Skip("id").Only(onlyColumns...))
		if err != nil {
			return s.config.ErrorHandler(err)
		}
	}

	return
}

// UpsertApplication updates one or more existing rows in applications
func (s Store) UpsertApplication(ctx context.Context, rr ...*types.Application) (err error) {
	for _, res := range rr {
		err = s.checkApplicationConstraints(ctx, res)
		if err != nil {
			return err
		}

		// err = s.applicationHook(ctx, TriggerBeforeApplicationUpsert, res)
		// if err != nil {
		// 	return err
		// }

		err = s.config.ErrorHandler(s.execUpsertApplications(ctx, s.internalApplicationEncoder(res)))
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteApplication Deletes one or more rows from applications table
func (s Store) DeleteApplication(ctx context.Context, rr ...*types.Application) (err error) {
	for _, res := range rr {
		// err = s.applicationHook(ctx, TriggerBeforeApplicationDelete, res)
		// if err != nil {
		// 	return err
		// }

		err = s.execDeleteApplications(ctx, squirrel.Eq{
			s.preprocessColumn("app.id", ""): s.preprocessValue(res.ID, ""),
		})
		if err != nil {
			return s.config.ErrorHandler(err)
		}
	}

	return nil
}

// DeleteApplicationByID Deletes row from the applications table
func (s Store) DeleteApplicationByID(ctx context.Context, ID uint64) error {
	return s.execDeleteApplications(ctx, squirrel.Eq{
		s.preprocessColumn("app.id", ""): s.preprocessValue(ID, ""),
	})
}

// TruncateApplications Deletes all rows from the applications table
func (s Store) TruncateApplications(ctx context.Context) error {
	return s.config.ErrorHandler(s.Truncate(ctx, s.applicationTable()))
}

// execLookupApplication prepares Application query and executes it,
// returning types.Application (or error)
func (s Store) execLookupApplication(ctx context.Context, cnd squirrel.Sqlizer) (res *types.Application, err error) {
	var (
		row rowScanner
	)

	row, err = s.QueryRow(ctx, s.applicationsSelectBuilder().Where(cnd))
	if err != nil {
		return
	}

	res, err = s.internalApplicationRowScanner(row)
	if err != nil {
		return
	}

	return res, nil
}

// execCreateApplications updates all matched (by cnd) rows in applications with given data
func (s Store) execCreateApplications(ctx context.Context, payload store.Payload) error {
	return s.config.ErrorHandler(s.Exec(ctx, s.InsertBuilder(s.applicationTable()).SetMap(payload)))
}

// execUpdateApplications updates all matched (by cnd) rows in applications with given data
func (s Store) execUpdateApplications(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return s.config.ErrorHandler(s.Exec(ctx, s.UpdateBuilder(s.applicationTable("app")).Where(cnd).SetMap(set)))
}

// execUpsertApplications inserts new or updates matching (by-primary-key) rows in applications with given data
func (s Store) execUpsertApplications(ctx context.Context, set store.Payload) error {
	upsert, err := s.config.UpsertBuilder(
		s.config,
		s.applicationTable(),
		set,
		"id",
	)

	if err != nil {
		return err
	}

	return s.config.ErrorHandler(s.Exec(ctx, upsert))
}

// execDeleteApplications Deletes all matched (by cnd) rows in applications with given data
func (s Store) execDeleteApplications(ctx context.Context, cnd squirrel.Sqlizer) error {
	return s.config.ErrorHandler(s.Exec(ctx, s.DeleteBuilder(s.applicationTable("app")).Where(cnd)))
}

func (s Store) internalApplicationRowScanner(row rowScanner) (res *types.Application, err error) {
	res = &types.Application{}

	if _, has := s.config.RowScanners["application"]; has {
		scanner := s.config.RowScanners["application"].(func(_ rowScanner, _ *types.Application) error)
		err = scanner(row, res)
	} else {
		err = row.Scan(
			&res.ID,
			&res.Name,
			&res.OwnerID,
			&res.Enabled,
			&res.Unify,
			&res.CreatedAt,
			&res.UpdatedAt,
			&res.DeletedAt,
		)
	}

	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("could not scan db row for Application: %w", err)
	} else {
		return res, nil
	}
}

// QueryApplications returns squirrel.SelectBuilder with set table and all columns
func (s Store) applicationsSelectBuilder() squirrel.SelectBuilder {
	return s.SelectBuilder(s.applicationTable("app"), s.applicationColumns("app")...)
}

// applicationTable name of the db table
func (Store) applicationTable(aa ...string) string {
	var alias string
	if len(aa) > 0 {
		alias = " AS " + aa[0]
	}

	return "applications" + alias
}

// ApplicationColumns returns all defined table columns
//
// With optional string arg, all columns are returned aliased
func (Store) applicationColumns(aa ...string) []string {
	var alias string
	if len(aa) > 0 {
		alias = aa[0] + "."
	}

	return []string{
		alias + "id",
		alias + "name",
		alias + "rel_owner",
		alias + "enabled",
		alias + "unify",
		alias + "created_at",
		alias + "updated_at",
		alias + "deleted_at",
	}
}

// {true true true true true}

// sortableApplicationColumns returns all Application columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
func (Store) sortableApplicationColumns() []string {
	return []string{
		"id",
		"name",
		"created_at",
		"updated_at",
		"deleted_at",
	}
}

// internalApplicationEncoder encodes fields from types.Application to store.Payload (map)
//
// Encoding is done by using generic approach or by calling encodeApplication
// func when rdbms.customEncoder=true
func (s Store) internalApplicationEncoder(res *types.Application) store.Payload {
	return store.Payload{
		"id":         res.ID,
		"name":       res.Name,
		"rel_owner":  res.OwnerID,
		"enabled":    res.Enabled,
		"unify":      res.Unify,
		"created_at": res.CreatedAt,
		"updated_at": res.UpdatedAt,
		"deleted_at": res.DeletedAt,
	}
}

func (s Store) collectApplicationCursorValues(res *types.Application, cc ...string) *filter.PagingCursor {
	var (
		cursor = &filter.PagingCursor{}

		hasUnique bool

		collect = func(cc ...string) {
			for _, c := range cc {
				switch c {
				case "id":
					cursor.Set(c, res.ID, false)
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
	if !hasUnique {
		collect(
			"id",
		)
	}

	return cursor
}

func (s *Store) checkApplicationConstraints(ctx context.Context, res *types.Application) error {

	return nil
}

// func (s *Store) applicationHook(ctx context.Context, key triggerKey, res *types.Application) error {
// 	if fn, has := s.config.TriggerHandlers[key]; has {
// 		return fn.(func (ctx context.Context, s *Store, res *types.Application) error)(ctx, s, res)
// 	}
//
// 	return nil
// }
