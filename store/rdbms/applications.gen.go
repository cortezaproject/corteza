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
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/jmoiron/sqlx"
)

// SearchApplications returns all matching rows
//
// This function calls convertApplicationFilter with the given
// types.ApplicationFilter and expects to receive a working squirrel.SelectBuilder
func (s Store) SearchApplications(ctx context.Context, f types.ApplicationFilter) (types.ApplicationSet, types.ApplicationFilter, error) {
	q, err := s.convertApplicationFilter(f)
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
		set = make([]*types.Application, 0, scap)
		// @todo this offset needs to be removed and replaced with key-based-paging
		fetchPage = func(offset, limit uint) (fetched, skipped uint, err error) {
			var (
				res *types.Application
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
				if res, err = s.internalApplicationRowScanner(rows, rows.Err()); err != nil {
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

// LookupApplicationByID searches for application by ID
//
// It returns application even if deleted
func (s Store) LookupApplicationByID(ctx context.Context, id uint64) (*types.Application, error) {
	return s.ApplicationLookup(ctx, squirrel.Eq{
		"app.id": id,
	})
}

// CreateApplication creates one or more rows in applications table
func (s Store) CreateApplication(ctx context.Context, rr ...*types.Application) error {
	if len(rr) == 0 {
		return nil
	}

	return Tx(ctx, s.db, s.config, nil, func(db *sqlx.Tx) (err error) {
		for _, res := range rr {
			err = ExecuteSqlizer(ctx, s.DB(), s.Insert(s.ApplicationTable()).SetMap(s.internalApplicationEncoder(res)))
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// UpdateApplication updates one or more existing rows in applications
func (s Store) UpdateApplication(ctx context.Context, rr ...*types.Application) error {
	return s.PartialUpdateApplication(ctx, nil, rr...)
}

// PartialUpdateApplication updates one or more existing rows in applications
//
// It wraps the update into transaction and can perform partial update by providing list of updatable columns
func (s Store) PartialUpdateApplication(ctx context.Context, onlyColumns []string, rr ...*types.Application) error {
	if len(rr) == 0 {
		return nil
	}

	return Tx(ctx, s.db, s.config, nil, func(db *sqlx.Tx) (err error) {
		for _, res := range rr {
			err = s.ExecUpdateApplications(
				ctx,
				squirrel.Eq{s.preprocessColumn("app.id", ""): s.preprocessValue(res.ID, "")},
				s.internalApplicationEncoder(res).Skip("id").Only(onlyColumns...))
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// RemoveApplication removes one or more rows from applications table
func (s Store) RemoveApplication(ctx context.Context, rr ...*types.Application) error {
	if len(rr) == 0 {
		return nil
	}

	return Tx(ctx, s.db, s.config, nil, func(db *sqlx.Tx) (err error) {
		for _, res := range rr {
			err = ExecuteSqlizer(ctx, s.DB(), s.Delete(s.ApplicationTable("app")).Where(squirrel.Eq{s.preprocessColumn("app.id", ""): s.preprocessValue(res.ID, "")}))
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// RemoveApplicationByID removes row from the applications table
func (s Store) RemoveApplicationByID(ctx context.Context, ID uint64) error {
	return ExecuteSqlizer(ctx, s.DB(), s.Delete(s.ApplicationTable("app")).Where(squirrel.Eq{s.preprocessColumn("app.id", ""): s.preprocessValue(ID, "")}))
}

// TruncateApplications removes all rows from the applications table
func (s Store) TruncateApplications(ctx context.Context) error {
	return Truncate(ctx, s.DB(), s.ApplicationTable())
}

// ExecUpdateApplications updates all matched (by cnd) rows in applications with given data
func (s Store) ExecUpdateApplications(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return ExecuteSqlizer(ctx, s.DB(), s.Update(s.ApplicationTable("app")).Where(cnd).SetMap(set))
}

// ApplicationLookup prepares Application query and executes it,
// returning types.Application (or error)
func (s Store) ApplicationLookup(ctx context.Context, cnd squirrel.Sqlizer) (*types.Application, error) {
	return s.internalApplicationRowScanner(s.QueryRow(ctx, s.QueryApplications().Where(cnd)))
}

func (s Store) internalApplicationRowScanner(row rowScanner, err error) (*types.Application, error) {
	if err != nil {
		return nil, err
	}

	var res = &types.Application{}
	if _, has := s.config.RowScanners["application"]; has {
		scanner := s.config.RowScanners["application"].(func(rowScanner, *types.Application) error)
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
func (s Store) QueryApplications() squirrel.SelectBuilder {
	return s.Select(s.ApplicationTable("app"), s.ApplicationColumns("app")...)
}

// ApplicationTable name of the db table
func (Store) ApplicationTable(aa ...string) string {
	var alias string
	if len(aa) > 0 {
		alias = " AS " + aa[0]
	}

	return "applications" + alias
}

// ApplicationColumns returns all defined table columns
//
// With optional string arg, all columns are returned aliased
func (Store) ApplicationColumns(aa ...string) []string {
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
