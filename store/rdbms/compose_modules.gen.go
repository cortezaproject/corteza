package rdbms

// This file is an auto-generated file
//
// Template:    pkg/codegen/assets/store_rdbms.gen.go.tpl
// Definitions: store/compose_modules.yaml
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

// SearchComposeModules returns all matching rows
//
// This function calls convertComposeModuleFilter with the given
// types.ModuleFilter and expects to receive a working squirrel.SelectBuilder
func (s Store) SearchComposeModules(ctx context.Context, f types.ModuleFilter) (types.ModuleSet, types.ModuleFilter, error) {
	q, err := s.convertComposeModuleFilter(f)
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
		set = make([]*types.Module, 0, scap)
		// @todo this offset needs to be removed and replaced with key-based-paging
		fetchPage = func(offset, limit uint) (fetched, skipped uint, err error) {
			var (
				res *types.Module
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
				if res, err = s.internalComposeModuleRowScanner(rows, rows.Err()); err != nil {
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

// LookupComposeModuleByHandle searches for compose module by handle (case-insensitive)
func (s Store) LookupComposeModuleByHandle(ctx context.Context, handle string) (*types.Module, error) {
	return s.ComposeModuleLookup(ctx, squirrel.Eq{
		"cmd.handle": handle,
	})
}

// LookupComposeModuleByID searches for compose module by ID
//
// It returns compose module even if deleted
func (s Store) LookupComposeModuleByID(ctx context.Context, id uint64) (*types.Module, error) {
	return s.ComposeModuleLookup(ctx, squirrel.Eq{
		"cmd.id": id,
	})
}

// CreateComposeModule creates one or more rows in compose_module table
func (s Store) CreateComposeModule(ctx context.Context, rr ...*types.Module) error {
	if len(rr) == 0 {
		return nil
	}

	return Tx(ctx, s.db, s.config, nil, func(db *sqlx.Tx) (err error) {
		for _, res := range rr {
			err = ExecuteSqlizer(ctx, s.DB(), s.Insert(s.ComposeModuleTable()).SetMap(s.internalComposeModuleEncoder(res)))
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// UpdateComposeModule updates one or more existing rows in compose_module
func (s Store) UpdateComposeModule(ctx context.Context, rr ...*types.Module) error {
	return s.PartialUpdateComposeModule(ctx, nil, rr...)
}

// PartialUpdateComposeModule updates one or more existing rows in compose_module
//
// It wraps the update into transaction and can perform partial update by providing list of updatable columns
func (s Store) PartialUpdateComposeModule(ctx context.Context, onlyColumns []string, rr ...*types.Module) error {
	if len(rr) == 0 {
		return nil
	}

	return Tx(ctx, s.db, s.config, nil, func(db *sqlx.Tx) (err error) {
		for _, res := range rr {
			err = s.ExecUpdateComposeModules(
				ctx,
				squirrel.Eq{s.preprocessColumn("cmd.id", ""): s.preprocessValue(res.ID, "")},
				s.internalComposeModuleEncoder(res).Skip("id").Only(onlyColumns...))
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// RemoveComposeModule removes one or more rows from compose_module table
func (s Store) RemoveComposeModule(ctx context.Context, rr ...*types.Module) error {
	if len(rr) == 0 {
		return nil
	}

	return Tx(ctx, s.db, s.config, nil, func(db *sqlx.Tx) (err error) {
		for _, res := range rr {
			err = ExecuteSqlizer(ctx, s.DB(), s.Delete(s.ComposeModuleTable("cmd")).Where(squirrel.Eq{s.preprocessColumn("cmd.id", ""): s.preprocessValue(res.ID, "")}))
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// RemoveComposeModuleByID removes row from the compose_module table
func (s Store) RemoveComposeModuleByID(ctx context.Context, ID uint64) error {
	return ExecuteSqlizer(ctx, s.DB(), s.Delete(s.ComposeModuleTable("cmd")).Where(squirrel.Eq{s.preprocessColumn("cmd.id", ""): s.preprocessValue(ID, "")}))
}

// TruncateComposeModules removes all rows from the compose_module table
func (s Store) TruncateComposeModules(ctx context.Context) error {
	return Truncate(ctx, s.DB(), s.ComposeModuleTable())
}

// ExecUpdateComposeModules updates all matched (by cnd) rows in compose_module with given data
func (s Store) ExecUpdateComposeModules(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return ExecuteSqlizer(ctx, s.DB(), s.Update(s.ComposeModuleTable("cmd")).Where(cnd).SetMap(set))
}

// ComposeModuleLookup prepares ComposeModule query and executes it,
// returning types.Module (or error)
func (s Store) ComposeModuleLookup(ctx context.Context, cnd squirrel.Sqlizer) (*types.Module, error) {
	return s.internalComposeModuleRowScanner(s.QueryRow(ctx, s.QueryComposeModules().Where(cnd)))
}

func (s Store) internalComposeModuleRowScanner(row rowScanner, err error) (*types.Module, error) {
	if err != nil {
		return nil, err
	}

	var res = &types.Module{}
	if _, has := s.config.RowScanners["composeModule"]; has {
		scanner := s.config.RowScanners["composeModule"].(func(rowScanner, *types.Module) error)
		err = scanner(row, res)
	} else {
		err = row.Scan(
			&res.ID,
			&res.Handle,
			&res.Name,
			&res.Meta,
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
		return nil, fmt.Errorf("could not scan db row for ComposeModule: %w", err)
	} else {
		return res, nil
	}
}

// QueryComposeModules returns squirrel.SelectBuilder with set table and all columns
func (s Store) QueryComposeModules() squirrel.SelectBuilder {
	return s.Select(s.ComposeModuleTable("cmd"), s.ComposeModuleColumns("cmd")...)
}

// ComposeModuleTable name of the db table
func (Store) ComposeModuleTable(aa ...string) string {
	var alias string
	if len(aa) > 0 {
		alias = " AS " + aa[0]
	}

	return "compose_module" + alias
}

// ComposeModuleColumns returns all defined table columns
//
// With optional string arg, all columns are returned aliased
func (Store) ComposeModuleColumns(aa ...string) []string {
	var alias string
	if len(aa) > 0 {
		alias = aa[0] + "."
	}

	return []string{
		alias + "id",
		alias + "handle",
		alias + "name",
		alias + "meta",
		alias + "rel_namespace",
		alias + "created_at",
		alias + "updated_at",
		alias + "deleted_at",
	}
}

// internalComposeModuleEncoder encodes fields from types.Module to store.Payload (map)
//
// Encoding is done by using generic approach or by calling encodeComposeModule
// func when rdbms.customEncoder=true
func (s Store) internalComposeModuleEncoder(res *types.Module) store.Payload {
	return store.Payload{
		"id":            res.ID,
		"handle":        res.Handle,
		"name":          res.Name,
		"meta":          res.Meta,
		"rel_namespace": res.NamespaceID,
		"created_at":    res.CreatedAt,
		"updated_at":    res.UpdatedAt,
		"deleted_at":    res.DeletedAt,
	}
}
