package rdbms

// This file is an auto-generated file
//
// Template:    pkg/codegen/assets/store_rdbms.gen.go.tpl
// Definitions: store/compose_charts.yaml
//
// Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated.

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/jmoiron/sqlx"
	"strings"
)

// SearchComposeCharts returns all matching rows
//
// This function calls convertComposeChartFilter with the given
// types.ChartFilter and expects to receive a working squirrel.SelectBuilder
func (s Store) SearchComposeCharts(ctx context.Context, f types.ChartFilter) (types.ChartSet, types.ChartFilter, error) {
	var scap uint
	q, err := s.convertComposeChartFilter(f)
	if err != nil {
		return nil, f, err
	}

	scap = f.Limit

	// Cleanup anything we've accidentally received...
	f.PrevPage, f.NextPage = nil, nil

	// When cursor for a previous page is used it's marked as reversed
	// This tells us to flip the descending flag on all used sort keys
	reverseCursor := f.PageCursor != nil && f.PageCursor.Reverse

	if err = f.Sort.Validate(s.sortableComposeChartColumns()...); err != nil {
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
		set = make([]*types.Chart, 0, scap)
		// fetches rows and scans them into Types.Chart resource this is then passed to Check function on filter
		// to help determine if fetched resource fits or not
		//
		// Note that limit is passed explicitly and is not necessarily equal to filter's limit. We want
		// to keep that value intact.
		//
		// The value for cursor is used and set directly from/to the filter!
		//
		// It returns total number of fetched pages and modifies PageCursor value for paging
		fetchPage = func(cursor *store.PagingCursor, limit uint) (fetched uint, err error) {
			var (
				res *types.Chart

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
				if res, err = s.internalComposeChartRowScanner(rows, rows.Err()); err != nil {
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
					f.PrevPage = s.collectComposeChartCursorValues(set[0], sort.Columns()...)
					f.PrevPage.Reverse = true
				}

				// Less items fetched then requested by page-limit
				// not very likely there's another page
				f.NextPage = s.collectComposeChartCursorValues(set[len(set)-1], sort.Columns()...)
			}

			f.PageCursor = nil
			return nil
		}
	)

	return set, f, fetch()
}

// LookupComposeChartByID searches for compose chart by ID
//
// It returns compose chart even if deleted
func (s Store) LookupComposeChartByID(ctx context.Context, id uint64) (*types.Chart, error) {
	return s.ComposeChartLookup(ctx, squirrel.Eq{
		"cch.id": id,
	})
}

// LookupComposeChartByHandle searches for compose chart by handle (case-insensitive)
func (s Store) LookupComposeChartByHandle(ctx context.Context, handle string) (*types.Chart, error) {
	return s.ComposeChartLookup(ctx, squirrel.Eq{
		"cch.handle": handle,
	})
}

// CreateComposeChart creates one or more rows in compose_chart table
func (s Store) CreateComposeChart(ctx context.Context, rr ...*types.Chart) error {
	if len(rr) == 0 {
		return nil
	}

	return Tx(ctx, s.db, s.config, nil, func(db *sqlx.Tx) (err error) {
		for _, res := range rr {
			err = ExecuteSqlizer(ctx, s.DB(), s.Insert(s.ComposeChartTable()).SetMap(s.internalComposeChartEncoder(res)))
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// UpdateComposeChart updates one or more existing rows in compose_chart
func (s Store) UpdateComposeChart(ctx context.Context, rr ...*types.Chart) error {
	return s.PartialUpdateComposeChart(ctx, nil, rr...)
}

// PartialUpdateComposeChart updates one or more existing rows in compose_chart
//
// It wraps the update into transaction and can perform partial update by providing list of updatable columns
func (s Store) PartialUpdateComposeChart(ctx context.Context, onlyColumns []string, rr ...*types.Chart) error {
	if len(rr) == 0 {
		return nil
	}

	return Tx(ctx, s.db, s.config, nil, func(db *sqlx.Tx) (err error) {
		for _, res := range rr {
			err = s.ExecUpdateComposeCharts(
				ctx,
				squirrel.Eq{s.preprocessColumn("cch.id", ""): s.preprocessValue(res.ID, "")},
				s.internalComposeChartEncoder(res).Skip("id").Only(onlyColumns...))
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// RemoveComposeChart removes one or more rows from compose_chart table
func (s Store) RemoveComposeChart(ctx context.Context, rr ...*types.Chart) error {
	if len(rr) == 0 {
		return nil
	}

	return Tx(ctx, s.db, s.config, nil, func(db *sqlx.Tx) (err error) {
		for _, res := range rr {
			err = ExecuteSqlizer(ctx, s.DB(), s.Delete(s.ComposeChartTable("cch")).Where(squirrel.Eq{s.preprocessColumn("cch.id", ""): s.preprocessValue(res.ID, "")}))
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// RemoveComposeChartByID removes row from the compose_chart table
func (s Store) RemoveComposeChartByID(ctx context.Context, ID uint64) error {
	return ExecuteSqlizer(ctx, s.DB(), s.Delete(s.ComposeChartTable("cch")).Where(squirrel.Eq{s.preprocessColumn("cch.id", ""): s.preprocessValue(ID, "")}))
}

// TruncateComposeCharts removes all rows from the compose_chart table
func (s Store) TruncateComposeCharts(ctx context.Context) error {
	return Truncate(ctx, s.DB(), s.ComposeChartTable())
}

// ExecUpdateComposeCharts updates all matched (by cnd) rows in compose_chart with given data
func (s Store) ExecUpdateComposeCharts(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return ExecuteSqlizer(ctx, s.DB(), s.Update(s.ComposeChartTable("cch")).Where(cnd).SetMap(set))
}

// ComposeChartLookup prepares ComposeChart query and executes it,
// returning types.Chart (or error)
func (s Store) ComposeChartLookup(ctx context.Context, cnd squirrel.Sqlizer) (*types.Chart, error) {
	return s.internalComposeChartRowScanner(s.QueryRow(ctx, s.QueryComposeCharts().Where(cnd)))
}

func (s Store) internalComposeChartRowScanner(row rowScanner, err error) (*types.Chart, error) {
	if err != nil {
		return nil, err
	}

	var res = &types.Chart{}
	if _, has := s.config.RowScanners["composeChart"]; has {
		scanner := s.config.RowScanners["composeChart"].(func(rowScanner, *types.Chart) error)
		err = scanner(row, res)
	} else {
		err = row.Scan(
			&res.ID,
			&res.Handle,
			&res.Name,
			&res.Config,
			&res.NamespaceID,
			&res.CreatedAt,
			&res.UpdatedAt,
			&res.DeletedAt,
		)
	}

	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("could not scan db row for ComposeChart: %w", err)
	} else {
		return res, nil
	}
}

// QueryComposeCharts returns squirrel.SelectBuilder with set table and all columns
func (s Store) QueryComposeCharts() squirrel.SelectBuilder {
	return s.Select(s.ComposeChartTable("cch"), s.ComposeChartColumns("cch")...)
}

// ComposeChartTable name of the db table
func (Store) ComposeChartTable(aa ...string) string {
	var alias string
	if len(aa) > 0 {
		alias = " AS " + aa[0]
	}

	return "compose_chart" + alias
}

// ComposeChartColumns returns all defined table columns
//
// With optional string arg, all columns are returned aliased
func (Store) ComposeChartColumns(aa ...string) []string {
	var alias string
	if len(aa) > 0 {
		alias = aa[0] + "."
	}

	return []string{
		alias + "id",
		alias + "handle",
		alias + "name",
		alias + "config",
		alias + "rel_namespace",
		alias + "created_at",
		alias + "updated_at",
		alias + "deleted_at",
	}
}

// {false false false false}

// sortableComposeChartColumns returns all ComposeChart columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
func (Store) sortableComposeChartColumns() []string {
	return []string{
		"id",
		"handle",
		"name",
		"created_at",
		"updated_at",
		"deleted_at",
	}
}

// internalComposeChartEncoder encodes fields from types.Chart to store.Payload (map)
//
// Encoding is done by using generic approach or by calling encodeComposeChart
// func when rdbms.customEncoder=true
func (s Store) internalComposeChartEncoder(res *types.Chart) store.Payload {
	return store.Payload{
		"id":            res.ID,
		"handle":        res.Handle,
		"name":          res.Name,
		"config":        res.Config,
		"rel_namespace": res.NamespaceID,
		"created_at":    res.CreatedAt,
		"updated_at":    res.UpdatedAt,
		"deleted_at":    res.DeletedAt,
	}
}

func (s Store) collectComposeChartCursorValues(res *types.Chart, cc ...string) *store.PagingCursor {
	var (
		cursor = &store.PagingCursor{}

		hasUnique bool

		collect = func(cc ...string) {
			for _, c := range cc {
				switch c {
				case "id":
					cursor.Set(c, res.ID, false)
				case "handle":
					cursor.Set(c, res.Handle, false)
					hasUnique = true
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
