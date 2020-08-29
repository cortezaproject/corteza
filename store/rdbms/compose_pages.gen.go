package rdbms

// This file is an auto-generated file
//
// Template:    pkg/codegen/assets/store_rdbms.gen.go.tpl
// Definitions: store/compose_pages.yaml
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
	TriggerBeforeComposePageCreate triggerKey = "composePageBeforeCreate"
	TriggerBeforeComposePageUpdate triggerKey = "composePageBeforeUpdate"
	TriggerBeforeComposePageUpsert triggerKey = "composePageBeforeUpsert"
	TriggerBeforeComposePageDelete triggerKey = "composePageBeforeDelete"
)

// SearchComposePages returns all matching rows
//
// This function calls convertComposePageFilter with the given
// types.PageFilter and expects to receive a working squirrel.SelectBuilder
func (s Store) SearchComposePages(ctx context.Context, f types.PageFilter) (types.PageSet, types.PageFilter, error) {
	var scap uint
	q, err := s.convertComposePageFilter(f)
	if err != nil {
		return nil, f, err
	}

	scap = f.Limit

	// Cleanup anything we've accidentally received...
	f.PrevPage, f.NextPage = nil, nil

	// When cursor for a previous page is used it's marked as reversed
	// This tells us to flip the descending flag on all used sort keys
	reverseCursor := f.PageCursor != nil && f.PageCursor.Reverse

	if err := f.Sort.Validate(s.sortableComposePageColumns()...); err != nil {
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
		set = make([]*types.Page, 0, scap)
		// fetches rows and scans them into types.Page resource this is then passed to Check function on filter
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
				res *types.Page

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
					res, err = s.internalComposePageRowScanner(rows)
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
					f.PrevPage = s.collectComposePageCursorValues(set[0], sort.Columns()...)
					f.PrevPage.Reverse = true
				}

				// Less items fetched then requested by page-limit
				// not very likely there's another page
				f.NextPage = s.collectComposePageCursorValues(set[len(set)-1], sort.Columns()...)
			}

			f.PageCursor = nil
			return nil
		}
	)

	return set, f, s.config.ErrorHandler(fetch())
}

// LookupComposePageByNamespaceIDHandle searches for page by handle (case-insensitive)
func (s Store) LookupComposePageByNamespaceIDHandle(ctx context.Context, namespace_id uint64, handle string) (*types.Page, error) {
	return s.execLookupComposePage(ctx, squirrel.Eq{
		s.preprocessColumn("cpg.rel_namespace", ""): s.preprocessValue(namespace_id, ""),
		s.preprocessColumn("cpg.handle", "lower"):   s.preprocessValue(handle, "lower"),
	})
}

// LookupComposePageByNamespaceIDModuleID searches for page by moduleID
func (s Store) LookupComposePageByNamespaceIDModuleID(ctx context.Context, namespace_id uint64, module_id uint64) (*types.Page, error) {
	return s.execLookupComposePage(ctx, squirrel.Eq{
		s.preprocessColumn("cpg.rel_namespace", ""): s.preprocessValue(namespace_id, ""),
		s.preprocessColumn("cpg.rel_module", ""):    s.preprocessValue(module_id, ""),
	})
}

// LookupComposePageByID searches for compose page by ID
//
// It returns compose page even if deleted
func (s Store) LookupComposePageByID(ctx context.Context, id uint64) (*types.Page, error) {
	return s.execLookupComposePage(ctx, squirrel.Eq{
		s.preprocessColumn("cpg.id", ""): s.preprocessValue(id, ""),
	})
}

// CreateComposePage creates one or more rows in compose_page table
func (s Store) CreateComposePage(ctx context.Context, rr ...*types.Page) (err error) {
	for _, res := range rr {
		err = s.checkComposePageConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execCreateComposePages(ctx, s.internalComposePageEncoder(res))
		if err != nil {
			return err
		}
	}

	return
}

// UpdateComposePage updates one or more existing rows in compose_page
func (s Store) UpdateComposePage(ctx context.Context, rr ...*types.Page) error {
	return s.config.ErrorHandler(s.PartialComposePageUpdate(ctx, nil, rr...))
}

// PartialComposePageUpdate updates one or more existing rows in compose_page
func (s Store) PartialComposePageUpdate(ctx context.Context, onlyColumns []string, rr ...*types.Page) (err error) {
	for _, res := range rr {
		err = s.checkComposePageConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpdateComposePages(
			ctx,
			squirrel.Eq{
				s.preprocessColumn("cpg.id", ""): s.preprocessValue(res.ID, ""),
			},
			s.internalComposePageEncoder(res).Skip("id").Only(onlyColumns...))
		if err != nil {
			return s.config.ErrorHandler(err)
		}
	}

	return
}

// UpsertComposePage updates one or more existing rows in compose_page
func (s Store) UpsertComposePage(ctx context.Context, rr ...*types.Page) (err error) {
	for _, res := range rr {
		err = s.checkComposePageConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.config.ErrorHandler(s.execUpsertComposePages(ctx, s.internalComposePageEncoder(res)))
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteComposePage Deletes one or more rows from compose_page table
func (s Store) DeleteComposePage(ctx context.Context, rr ...*types.Page) (err error) {
	for _, res := range rr {

		err = s.execDeleteComposePages(ctx, squirrel.Eq{
			s.preprocessColumn("cpg.id", ""): s.preprocessValue(res.ID, ""),
		})
		if err != nil {
			return s.config.ErrorHandler(err)
		}
	}

	return nil
}

// DeleteComposePageByID Deletes row from the compose_page table
func (s Store) DeleteComposePageByID(ctx context.Context, ID uint64) error {
	return s.execDeleteComposePages(ctx, squirrel.Eq{
		s.preprocessColumn("cpg.id", ""): s.preprocessValue(ID, ""),
	})
}

// TruncateComposePages Deletes all rows from the compose_page table
func (s Store) TruncateComposePages(ctx context.Context) error {
	return s.config.ErrorHandler(s.Truncate(ctx, s.composePageTable()))
}

// execLookupComposePage prepares ComposePage query and executes it,
// returning types.Page (or error)
func (s Store) execLookupComposePage(ctx context.Context, cnd squirrel.Sqlizer) (res *types.Page, err error) {
	var (
		row rowScanner
	)

	row, err = s.QueryRow(ctx, s.composePagesSelectBuilder().Where(cnd))
	if err != nil {
		return
	}

	res, err = s.internalComposePageRowScanner(row)
	if err != nil {
		return
	}

	return res, nil
}

// execCreateComposePages updates all matched (by cnd) rows in compose_page with given data
func (s Store) execCreateComposePages(ctx context.Context, payload store.Payload) error {
	return s.config.ErrorHandler(s.Exec(ctx, s.InsertBuilder(s.composePageTable()).SetMap(payload)))
}

// execUpdateComposePages updates all matched (by cnd) rows in compose_page with given data
func (s Store) execUpdateComposePages(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return s.config.ErrorHandler(s.Exec(ctx, s.UpdateBuilder(s.composePageTable("cpg")).Where(cnd).SetMap(set)))
}

// execUpsertComposePages inserts new or updates matching (by-primary-key) rows in compose_page with given data
func (s Store) execUpsertComposePages(ctx context.Context, set store.Payload) error {
	upsert, err := s.config.UpsertBuilder(
		s.config,
		s.composePageTable(),
		set,
		"id",
	)

	if err != nil {
		return err
	}

	return s.config.ErrorHandler(s.Exec(ctx, upsert))
}

// execDeleteComposePages Deletes all matched (by cnd) rows in compose_page with given data
func (s Store) execDeleteComposePages(ctx context.Context, cnd squirrel.Sqlizer) error {
	return s.config.ErrorHandler(s.Exec(ctx, s.DeleteBuilder(s.composePageTable("cpg")).Where(cnd)))
}

func (s Store) internalComposePageRowScanner(row rowScanner) (res *types.Page, err error) {
	res = &types.Page{}

	if _, has := s.config.RowScanners["composePage"]; has {
		scanner := s.config.RowScanners["composePage"].(func(_ rowScanner, _ *types.Page) error)
		err = scanner(row, res)
	} else {
		err = row.Scan(
			&res.ID,
			&res.SelfID,
			&res.NamespaceID,
			&res.ModuleID,
			&res.Handle,
			&res.Title,
			&res.Description,
			&res.Blocks,
			&res.Visible,
			&res.Weight,
			&res.CreatedAt,
			&res.UpdatedAt,
			&res.DeletedAt,
		)
	}

	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("could not scan db row for ComposePage: %w", err)
	} else {
		return res, nil
	}
}

// QueryComposePages returns squirrel.SelectBuilder with set table and all columns
func (s Store) composePagesSelectBuilder() squirrel.SelectBuilder {
	return s.SelectBuilder(s.composePageTable("cpg"), s.composePageColumns("cpg")...)
}

// composePageTable name of the db table
func (Store) composePageTable(aa ...string) string {
	var alias string
	if len(aa) > 0 {
		alias = " AS " + aa[0]
	}

	return "compose_page" + alias
}

// ComposePageColumns returns all defined table columns
//
// With optional string arg, all columns are returned aliased
func (Store) composePageColumns(aa ...string) []string {
	var alias string
	if len(aa) > 0 {
		alias = aa[0] + "."
	}

	return []string{
		alias + "id",
		alias + "self_id",
		alias + "rel_namespace",
		alias + "rel_module",
		alias + "handle",
		alias + "title",
		alias + "description",
		alias + "blocks",
		alias + "visible",
		alias + "weight",
		alias + "created_at",
		alias + "updated_at",
		alias + "deleted_at",
	}
}

// {true true true true true}

// sortableComposePageColumns returns all ComposePage columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
func (Store) sortableComposePageColumns() []string {
	return []string{
		"id",
		"created_at",
		"updated_at",
		"deleted_at",
	}
}

// internalComposePageEncoder encodes fields from types.Page to store.Payload (map)
//
// Encoding is done by using generic approach or by calling encodeComposePage
// func when rdbms.customEncoder=true
func (s Store) internalComposePageEncoder(res *types.Page) store.Payload {
	return store.Payload{
		"id":            res.ID,
		"self_id":       res.SelfID,
		"rel_namespace": res.NamespaceID,
		"rel_module":    res.ModuleID,
		"handle":        res.Handle,
		"title":         res.Title,
		"description":   res.Description,
		"blocks":        res.Blocks,
		"visible":       res.Visible,
		"weight":        res.Weight,
		"created_at":    res.CreatedAt,
		"updated_at":    res.UpdatedAt,
		"deleted_at":    res.DeletedAt,
	}
}

func (s Store) collectComposePageCursorValues(res *types.Page, cc ...string) *filter.PagingCursor {
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

func (s *Store) checkComposePageConstraints(ctx context.Context, res *types.Page) error {

	return nil
}
