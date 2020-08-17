package rdbms

// This file is an auto-generated file
//
// Template:    pkg/codegen/assets/store_rdbms.gen.go.tpl
// Definitions: store/compose_namespaces.yaml
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

// SearchComposeNamespaces returns all matching rows
//
// This function calls convertComposeNamespaceFilter with the given
// types.NamespaceFilter and expects to receive a working squirrel.SelectBuilder
func (s Store) SearchComposeNamespaces(ctx context.Context, f types.NamespaceFilter) (types.NamespaceSet, types.NamespaceFilter, error) {
	q, err := s.convertComposeNamespaceFilter(f)
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
		set = make([]*types.Namespace, 0, scap)
		// @todo this offset needs to be removed and replaced with key-based-paging
		fetchPage = func(offset, limit uint) (fetched, skipped uint, err error) {
			var (
				res *types.Namespace
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
				if res, err = s.internalComposeNamespaceRowScanner(rows, rows.Err()); err != nil {
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

// LookupComposeNamespaceBySlug searches for namespace by slug (case-insensitive)
func (s Store) LookupComposeNamespaceBySlug(ctx context.Context, slug string) (*types.Namespace, error) {
	return s.ComposeNamespaceLookup(ctx, squirrel.Eq{
		"cns.slug": slug,
	})
}

// LookupComposeNamespaceByID searches for compose namespace by ID
//
// It returns compose namespace even if deleted
func (s Store) LookupComposeNamespaceByID(ctx context.Context, id uint64) (*types.Namespace, error) {
	return s.ComposeNamespaceLookup(ctx, squirrel.Eq{
		"cns.id": id,
	})
}

// CreateComposeNamespace creates one or more rows in compose_namespace table
func (s Store) CreateComposeNamespace(ctx context.Context, rr ...*types.Namespace) error {
	if len(rr) == 0 {
		return nil
	}

	return Tx(ctx, s.db, s.config, nil, func(db *sqlx.Tx) (err error) {
		for _, res := range rr {
			err = ExecuteSqlizer(ctx, s.DB(), s.Insert(s.ComposeNamespaceTable()).SetMap(s.internalComposeNamespaceEncoder(res)))
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// UpdateComposeNamespace updates one or more existing rows in compose_namespace
func (s Store) UpdateComposeNamespace(ctx context.Context, rr ...*types.Namespace) error {
	return s.PartialUpdateComposeNamespace(ctx, nil, rr...)
}

// PartialUpdateComposeNamespace updates one or more existing rows in compose_namespace
//
// It wraps the update into transaction and can perform partial update by providing list of updatable columns
func (s Store) PartialUpdateComposeNamespace(ctx context.Context, onlyColumns []string, rr ...*types.Namespace) error {
	if len(rr) == 0 {
		return nil
	}

	return Tx(ctx, s.db, s.config, nil, func(db *sqlx.Tx) (err error) {
		for _, res := range rr {
			err = s.ExecUpdateComposeNamespaces(
				ctx,
				squirrel.Eq{s.preprocessColumn("cns.id", ""): s.preprocessValue(res.ID, "")},
				s.internalComposeNamespaceEncoder(res).Skip("id").Only(onlyColumns...))
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// RemoveComposeNamespace removes one or more rows from compose_namespace table
func (s Store) RemoveComposeNamespace(ctx context.Context, rr ...*types.Namespace) error {
	if len(rr) == 0 {
		return nil
	}

	return Tx(ctx, s.db, s.config, nil, func(db *sqlx.Tx) (err error) {
		for _, res := range rr {
			err = ExecuteSqlizer(ctx, s.DB(), s.Delete(s.ComposeNamespaceTable("cns")).Where(squirrel.Eq{s.preprocessColumn("cns.id", ""): s.preprocessValue(res.ID, "")}))
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// RemoveComposeNamespaceByID removes row from the compose_namespace table
func (s Store) RemoveComposeNamespaceByID(ctx context.Context, ID uint64) error {
	return ExecuteSqlizer(ctx, s.DB(), s.Delete(s.ComposeNamespaceTable("cns")).Where(squirrel.Eq{s.preprocessColumn("cns.id", ""): s.preprocessValue(ID, "")}))
}

// TruncateComposeNamespaces removes all rows from the compose_namespace table
func (s Store) TruncateComposeNamespaces(ctx context.Context) error {
	return Truncate(ctx, s.DB(), s.ComposeNamespaceTable())
}

// ExecUpdateComposeNamespaces updates all matched (by cnd) rows in compose_namespace with given data
func (s Store) ExecUpdateComposeNamespaces(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return ExecuteSqlizer(ctx, s.DB(), s.Update(s.ComposeNamespaceTable("cns")).Where(cnd).SetMap(set))
}

// ComposeNamespaceLookup prepares ComposeNamespace query and executes it,
// returning types.Namespace (or error)
func (s Store) ComposeNamespaceLookup(ctx context.Context, cnd squirrel.Sqlizer) (*types.Namespace, error) {
	return s.internalComposeNamespaceRowScanner(s.QueryRow(ctx, s.QueryComposeNamespaces().Where(cnd)))
}

func (s Store) internalComposeNamespaceRowScanner(row rowScanner, err error) (*types.Namespace, error) {
	if err != nil {
		return nil, err
	}

	var res = &types.Namespace{}
	if _, has := s.config.RowScanners["composeNamespace"]; has {
		scanner := s.config.RowScanners["composeNamespace"].(func(rowScanner, *types.Namespace) error)
		err = scanner(row, res)
	} else {
		err = row.Scan(
			&res.ID,
			&res.Name,
			&res.Slug,
			&res.Enabled,
			&res.Meta,
			&res.CreatedAt,
			&res.UpdatedAt,
			&res.DeletedAt,
		)
	}

	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("could not scan db row for ComposeNamespace: %w", err)
	} else {
		return res, nil
	}
}

// QueryComposeNamespaces returns squirrel.SelectBuilder with set table and all columns
func (s Store) QueryComposeNamespaces() squirrel.SelectBuilder {
	return s.Select(s.ComposeNamespaceTable("cns"), s.ComposeNamespaceColumns("cns")...)
}

// ComposeNamespaceTable name of the db table
func (Store) ComposeNamespaceTable(aa ...string) string {
	var alias string
	if len(aa) > 0 {
		alias = " AS " + aa[0]
	}

	return "compose_namespace" + alias
}

// ComposeNamespaceColumns returns all defined table columns
//
// With optional string arg, all columns are returned aliased
func (Store) ComposeNamespaceColumns(aa ...string) []string {
	var alias string
	if len(aa) > 0 {
		alias = aa[0] + "."
	}

	return []string{
		alias + "id",
		alias + "name",
		alias + "slug",
		alias + "enabled",
		alias + "meta",
		alias + "created_at",
		alias + "updated_at",
		alias + "deleted_at",
	}
}

// internalComposeNamespaceEncoder encodes fields from types.Namespace to store.Payload (map)
//
// Encoding is done by using generic approach or by calling encodeComposeNamespace
// func when rdbms.customEncoder=true
func (s Store) internalComposeNamespaceEncoder(res *types.Namespace) store.Payload {
	return store.Payload{
		"id":         res.ID,
		"name":       res.Name,
		"slug":       res.Slug,
		"enabled":    res.Enabled,
		"meta":       res.Meta,
		"created_at": res.CreatedAt,
		"updated_at": res.UpdatedAt,
		"deleted_at": res.DeletedAt,
	}
}
