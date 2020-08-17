package rdbms

// This file is an auto-generated file
//
// Template:    pkg/codegen/assets/store_rdbms.gen.go.tpl
// Definitions: store/users.yaml
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

// SearchUsers returns all matching rows
//
// This function calls convertUserFilter with the given
// types.UserFilter and expects to receive a working squirrel.SelectBuilder
func (s Store) SearchUsers(ctx context.Context, f types.UserFilter) (types.UserSet, types.UserFilter, error) {
	q, err := s.convertUserFilter(f)
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
		set = make([]*types.User, 0, scap)
		// @todo this offset needs to be removed and replaced with key-based-paging
		fetchPage = func(offset, limit uint) (fetched, skipped uint, err error) {
			var (
				res *types.User
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
				if res, err = s.internalUserRowScanner(rows, rows.Err()); err != nil {
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

// LookupUserByID searches for user by ID
//
// It returns user even if deleted or suspended
func (s Store) LookupUserByID(ctx context.Context, id uint64) (*types.User, error) {
	return s.UserLookup(ctx, squirrel.Eq{
		"usr.id": id,
	})
}

// LookupUserByEmail searches for user by their email
//
// It returns only valid users (not deleted, not suspended)
func (s Store) LookupUserByEmail(ctx context.Context, email string) (*types.User, error) {
	return s.UserLookup(ctx, squirrel.Eq{
		"usr.email":        email,
		"usr.deleted_at":   nil,
		"usr.suspended_at": nil,
	})
}

// LookupUserByHandle searches for user by their email
//
// It returns only valid users (not deleted, not suspended)
func (s Store) LookupUserByHandle(ctx context.Context, handle string) (*types.User, error) {
	return s.UserLookup(ctx, squirrel.Eq{
		"usr.handle":       handle,
		"usr.deleted_at":   nil,
		"usr.suspended_at": nil,
	})
}

// LookupUserByUsername searches for user by their username
//
// It returns only valid users (not deleted, not suspended)
func (s Store) LookupUserByUsername(ctx context.Context, username string) (*types.User, error) {
	return s.UserLookup(ctx, squirrel.Eq{
		"usr.username":     username,
		"usr.deleted_at":   nil,
		"usr.suspended_at": nil,
	})
}

// CreateUser creates one or more rows in users table
func (s Store) CreateUser(ctx context.Context, rr ...*types.User) error {
	if len(rr) == 0 {
		return nil
	}

	return Tx(ctx, s.db, s.config, nil, func(db *sqlx.Tx) (err error) {
		for _, res := range rr {
			err = ExecuteSqlizer(ctx, s.DB(), s.Insert(s.UserTable()).SetMap(s.internalUserEncoder(res)))
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// UpdateUser updates one or more existing rows in users
func (s Store) UpdateUser(ctx context.Context, rr ...*types.User) error {
	return s.PartialUpdateUser(ctx, nil, rr...)
}

// PartialUpdateUser updates one or more existing rows in users
//
// It wraps the update into transaction and can perform partial update by providing list of updatable columns
func (s Store) PartialUpdateUser(ctx context.Context, onlyColumns []string, rr ...*types.User) error {
	if len(rr) == 0 {
		return nil
	}

	return Tx(ctx, s.db, s.config, nil, func(db *sqlx.Tx) (err error) {
		for _, res := range rr {
			err = s.ExecUpdateUsers(
				ctx,
				squirrel.Eq{s.preprocessColumn("usr.id", ""): s.preprocessValue(res.ID, "")},
				s.internalUserEncoder(res).Skip("id").Only(onlyColumns...))
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// RemoveUser removes one or more rows from users table
func (s Store) RemoveUser(ctx context.Context, rr ...*types.User) error {
	if len(rr) == 0 {
		return nil
	}

	return Tx(ctx, s.db, s.config, nil, func(db *sqlx.Tx) (err error) {
		for _, res := range rr {
			err = ExecuteSqlizer(ctx, s.DB(), s.Delete(s.UserTable("usr")).Where(squirrel.Eq{s.preprocessColumn("usr.id", ""): s.preprocessValue(res.ID, "")}))
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// RemoveUserByID removes row from the users table
func (s Store) RemoveUserByID(ctx context.Context, ID uint64) error {
	return ExecuteSqlizer(ctx, s.DB(), s.Delete(s.UserTable("usr")).Where(squirrel.Eq{s.preprocessColumn("usr.id", ""): s.preprocessValue(ID, "")}))
}

// TruncateUsers removes all rows from the users table
func (s Store) TruncateUsers(ctx context.Context) error {
	return Truncate(ctx, s.DB(), s.UserTable())
}

// ExecUpdateUsers updates all matched (by cnd) rows in users with given data
func (s Store) ExecUpdateUsers(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return ExecuteSqlizer(ctx, s.DB(), s.Update(s.UserTable("usr")).Where(cnd).SetMap(set))
}

// UserLookup prepares User query and executes it,
// returning types.User (or error)
func (s Store) UserLookup(ctx context.Context, cnd squirrel.Sqlizer) (*types.User, error) {
	return s.internalUserRowScanner(s.QueryRow(ctx, s.QueryUsers().Where(cnd)))
}

func (s Store) internalUserRowScanner(row rowScanner, err error) (*types.User, error) {
	if err != nil {
		return nil, err
	}

	var res = &types.User{}
	if _, has := s.config.RowScanners["user"]; has {
		scanner := s.config.RowScanners["user"].(func(rowScanner, *types.User) error)
		err = scanner(row, res)
	} else {
		err = row.Scan(
			&res.ID,
			&res.Email,
			&res.EmailConfirmed,
			&res.Username,
			&res.Name,
			&res.Handle,
			&res.Meta,
			&res.Kind,
			&res.CreatedAt,
			&res.UpdatedAt,
			&res.SuspendedAt,
			&res.DeletedAt,
		)
	}

	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("could not scan db row for User: %w", err)
	} else {
		return res, nil
	}
}

// QueryUsers returns squirrel.SelectBuilder with set table and all columns
func (s Store) QueryUsers() squirrel.SelectBuilder {
	return s.Select(s.UserTable("usr"), s.UserColumns("usr")...)
}

// UserTable name of the db table
func (Store) UserTable(aa ...string) string {
	var alias string
	if len(aa) > 0 {
		alias = " AS " + aa[0]
	}

	return "users" + alias
}

// UserColumns returns all defined table columns
//
// With optional string arg, all columns are returned aliased
func (Store) UserColumns(aa ...string) []string {
	var alias string
	if len(aa) > 0 {
		alias = aa[0] + "."
	}

	return []string{
		alias + "id",
		alias + "email",
		alias + "email_confirmed",
		alias + "username",
		alias + "name",
		alias + "handle",
		alias + "meta",
		alias + "kind",
		alias + "created_at",
		alias + "updated_at",
		alias + "suspended_at",
		alias + "deleted_at",
	}
}

// internalUserEncoder encodes fields from types.User to store.Payload (map)
//
// Encoding is done by using generic approach or by calling encodeUser
// func when rdbms.customEncoder=true
func (s Store) internalUserEncoder(res *types.User) store.Payload {
	return store.Payload{
		"id":              res.ID,
		"email":           res.Email,
		"email_confirmed": res.EmailConfirmed,
		"username":        res.Username,
		"name":            res.Name,
		"handle":          res.Handle,
		"meta":            res.Meta,
		"kind":            res.Kind,
		"created_at":      res.CreatedAt,
		"updated_at":      res.UpdatedAt,
		"suspended_at":    res.SuspendedAt,
		"deleted_at":      res.DeletedAt,
	}
}
