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
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/jmoiron/sqlx"
)

// SearchComposePages returns all matching rows
//
// This function calls convertComposePageFilter with the given
// types.PageFilter and expects to receive a working squirrel.SelectBuilder
func (s Store) SearchComposePages(ctx context.Context, f types.PageFilter) (types.PageSet, types.PageFilter, error) {
	q, err := s.convertComposePageFilter(f)
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
		set = make([]*types.Page, 0, scap)
		// @todo this offset needs to be removed and replaced with key-based-paging
		fetchPage = func(offset, limit uint) (fetched, skipped uint, err error) {
			var (
				res *types.Page
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
				if res, err = s.internalComposePageRowScanner(rows, rows.Err()); err != nil {
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

// LookupComposePageByHandle searches for page chart by handle (case-insensitive)
func (s Store) LookupComposePageByHandle(ctx context.Context, handle string) (*types.Page, error) {
	return s.ComposePageLookup(ctx, squirrel.Eq{
		"cpg.handle": handle,
	})
}

// LookupComposePageByID searches for compose page by ID
//
// It returns compose page even if deleted
func (s Store) LookupComposePageByID(ctx context.Context, id uint64) (*types.Page, error) {
	return s.ComposePageLookup(ctx, squirrel.Eq{
		"cpg.id": id,
	})
}

// CreateComposePage creates one or more rows in compose_page table
func (s Store) CreateComposePage(ctx context.Context, rr ...*types.Page) error {
	if len(rr) == 0 {
		return nil
	}

	return Tx(ctx, s.db, s.config, nil, func(db *sqlx.Tx) (err error) {
		for _, res := range rr {
			err = ExecuteSqlizer(ctx, s.DB(), s.Insert(s.ComposePageTable()).SetMap(s.internalComposePageEncoder(res)))
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// UpdateComposePage updates one or more existing rows in compose_page
func (s Store) UpdateComposePage(ctx context.Context, rr ...*types.Page) error {
	return s.PartialUpdateComposePage(ctx, nil, rr...)
}

// PartialUpdateComposePage updates one or more existing rows in compose_page
//
// It wraps the update into transaction and can perform partial update by providing list of updatable columns
func (s Store) PartialUpdateComposePage(ctx context.Context, onlyColumns []string, rr ...*types.Page) error {
	if len(rr) == 0 {
		return nil
	}

	return Tx(ctx, s.db, s.config, nil, func(db *sqlx.Tx) (err error) {
		for _, res := range rr {
			err = s.ExecUpdateComposePages(
				ctx,
				squirrel.Eq{s.preprocessColumn("cpg.id", ""): s.preprocessValue(res.ID, "")},
				s.internalComposePageEncoder(res).Skip("id").Only(onlyColumns...))
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// RemoveComposePage removes one or more rows from compose_page table
func (s Store) RemoveComposePage(ctx context.Context, rr ...*types.Page) error {
	if len(rr) == 0 {
		return nil
	}

	return Tx(ctx, s.db, s.config, nil, func(db *sqlx.Tx) (err error) {
		for _, res := range rr {
			err = ExecuteSqlizer(ctx, s.DB(), s.Delete(s.ComposePageTable("cpg")).Where(squirrel.Eq{s.preprocessColumn("cpg.id", ""): s.preprocessValue(res.ID, "")}))
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// RemoveComposePageByID removes row from the compose_page table
func (s Store) RemoveComposePageByID(ctx context.Context, ID uint64) error {
	return ExecuteSqlizer(ctx, s.DB(), s.Delete(s.ComposePageTable("cpg")).Where(squirrel.Eq{s.preprocessColumn("cpg.id", ""): s.preprocessValue(ID, "")}))
}

// TruncateComposePages removes all rows from the compose_page table
func (s Store) TruncateComposePages(ctx context.Context) error {
	return Truncate(ctx, s.DB(), s.ComposePageTable())
}

// ExecUpdateComposePages updates all matched (by cnd) rows in compose_page with given data
func (s Store) ExecUpdateComposePages(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return ExecuteSqlizer(ctx, s.DB(), s.Update(s.ComposePageTable("cpg")).Where(cnd).SetMap(set))
}

// ComposePageLookup prepares ComposePage query and executes it,
// returning types.Page (or error)
func (s Store) ComposePageLookup(ctx context.Context, cnd squirrel.Sqlizer) (*types.Page, error) {
	return s.internalComposePageRowScanner(s.QueryRow(ctx, s.QueryComposePages().Where(cnd)))
}

func (s Store) internalComposePageRowScanner(row rowScanner, err error) (*types.Page, error) {
	if err != nil {
		return nil, err
	}

	var res = &types.Page{}
	if _, has := s.config.RowScanners["composePage"]; has {
		scanner := s.config.RowScanners["composePage"].(func(rowScanner, *types.Page) error)
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
func (s Store) QueryComposePages() squirrel.SelectBuilder {
	return s.Select(s.ComposePageTable("cpg"), s.ComposePageColumns("cpg")...)
}

// ComposePageTable name of the db table
func (Store) ComposePageTable(aa ...string) string {
	var alias string
	if len(aa) > 0 {
		alias = " AS " + aa[0]
	}

	return "compose_page" + alias
}

// ComposePageColumns returns all defined table columns
//
// With optional string arg, all columns are returned aliased
func (Store) ComposePageColumns(aa ...string) []string {
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
