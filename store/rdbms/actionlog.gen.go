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
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/store"
	"strings"
)

var _ = errors.Is

const (
	TriggerBeforeActionlogCreate triggerKey = "actionlogBeforeCreate"
)

// SearchActionlogs returns all matching rows
//
// This function calls convertActionlogFilter with the given
// actionlog.Filter and expects to receive a working squirrel.SelectBuilder
func (s Store) SearchActionlogs(ctx context.Context, f actionlog.Filter) (actionlog.ActionSet, actionlog.Filter, error) {
	var scap uint
	q, err := s.convertActionlogFilter(f)
	if err != nil {
		return nil, f, err
	}

	scap = f.Limit

	// Cleanup anything we've accidentally received...
	f.PrevPage, f.NextPage = nil, nil

	// When cursor for a previous page is used it's marked as reversed
	// This tells us to flip the descending flag on all used sort keys
	reverseCursor := f.PageCursor != nil && f.PageCursor.Reverse

	// Sorting is disabled in definition yaml file
	// {search: {enablePaging:false}}
	//
	// We still need to sort the results by primary key for paging purposes
	sort := filter.SortExprSet{
		&filter.SortExpr{Column: "id", Descending: true},
	}

	if scap == 0 {
		scap = DefaultSliceCapacity
	}

	var (
		set = make([]*actionlog.Action, 0, scap)
		// fetches rows and scans them into actionlog.Action resource this is then passed to Check function on filter
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
				res *actionlog.Action

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
					res, err = s.internalActionlogRowScanner(rows)
				}

				if err != nil {
					if cerr := rows.Close(); cerr != nil {
						err = fmt.Errorf("could not close rows (%v) after scan error: %w", cerr, err)
					}

					return
				}

				// If check function is set, call it and act accordingly
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
					f.PrevPage = s.collectActionlogCursorValues(set[0], sort.Columns()...)
					f.PrevPage.Reverse = true
				}

				// Less items fetched then requested by page-limit
				// not very likely there's another page
				f.NextPage = s.collectActionlogCursorValues(set[len(set)-1], sort.Columns()...)
			}

			f.PageCursor = nil
			return nil
		}
	)

	return set, f, s.config.ErrorHandler(fetch())
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
	return s.config.ErrorHandler(s.Truncate(ctx, s.actionlogTable()))
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
	return s.config.ErrorHandler(s.Exec(ctx, s.InsertBuilder(s.actionlogTable()).SetMap(payload)))
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
		return nil, store.ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("could not scan db row for Actionlog: %w", err)
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

// {true true true false false}

// internalActionlogEncoder encodes fields from actionlog.Action to store.Payload (map)
//
// Encoding is done by using generic approach or by calling encodeActionlog
// func when rdbms.customEncoder=true
func (s Store) internalActionlogEncoder(res *actionlog.Action) store.Payload {
	return s.encodeActionlog(res)
}

func (s Store) collectActionlogCursorValues(res *actionlog.Action, cc ...string) *filter.PagingCursor {
	var (
		cursor = &filter.PagingCursor{}

		hasUnique bool

		collect = func(cc ...string) {
			for _, c := range cc {
				switch c {
				case "id":
					cursor.Set(c, res.ID, false)

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

func (s *Store) checkActionlogConstraints(ctx context.Context, res *actionlog.Action) error {

	return nil
}
