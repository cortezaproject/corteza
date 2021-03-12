package rdbms

// This file is an auto-generated file
//
// Template:    pkg/codegen/assets/store_rdbms.gen.go.tpl
// Definitions: store/auth_oa2tokens.yaml
//
// Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated.

import (
	"context"
	"database/sql"
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
)

var _ = errors.Is

// SearchAuthOa2tokens returns all matching rows
//
// This function calls convertAuthOa2tokenFilter with the given
// types.AuthOa2tokenFilter and expects to receive a working squirrel.SelectBuilder
func (s Store) SearchAuthOa2tokens(ctx context.Context, f types.AuthOa2tokenFilter) (types.AuthOa2tokenSet, types.AuthOa2tokenFilter, error) {
	var (
		err error
		set []*types.AuthOa2token
		q   squirrel.SelectBuilder
	)

	return set, f, func() error {
		q, err = s.convertAuthOa2tokenFilter(f)
		if err != nil {
			return err
		}

		set, err = s.QueryAuthOa2tokens(ctx, q, nil)
		return err
	}()
}

// QueryAuthOa2tokens queries the database, converts and checks each row and
// returns collected set
//
// Fn also returns total number of fetched items and last fetched item so that the caller can construct cursor
// for next page of results
func (s Store) QueryAuthOa2tokens(
	ctx context.Context,
	q squirrel.Sqlizer,
	check func(*types.AuthOa2token) (bool, error),
) ([]*types.AuthOa2token, error) {
	var (
		set = make([]*types.AuthOa2token, 0, DefaultSliceCapacity)
		res *types.AuthOa2token

		// Query rows with
		rows, err = s.Query(ctx, q)
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		if err = rows.Err(); err == nil {
			res, err = s.internalAuthOa2tokenRowScanner(rows)
		}

		if err != nil {
			return nil, err
		}

		set = append(set, res)
	}

	return set, rows.Err()
}

// LookupAuthOa2tokenByCode
func (s Store) LookupAuthOa2tokenByCode(ctx context.Context, code string) (*types.AuthOa2token, error) {
	return s.execLookupAuthOa2token(ctx, squirrel.Eq{
		s.preprocessColumn("tkn.code", ""): store.PreprocessValue(code, ""),
	})
}

// LookupAuthOa2tokenByAccess
func (s Store) LookupAuthOa2tokenByAccess(ctx context.Context, access string) (*types.AuthOa2token, error) {
	return s.execLookupAuthOa2token(ctx, squirrel.Eq{
		s.preprocessColumn("tkn.access", ""): store.PreprocessValue(access, ""),
	})
}

// LookupAuthOa2tokenByRefresh
func (s Store) LookupAuthOa2tokenByRefresh(ctx context.Context, refresh string) (*types.AuthOa2token, error) {
	return s.execLookupAuthOa2token(ctx, squirrel.Eq{
		s.preprocessColumn("tkn.refresh", ""): store.PreprocessValue(refresh, ""),
	})
}

// CreateAuthOa2token creates one or more rows in auth_oa2tokens table
func (s Store) CreateAuthOa2token(ctx context.Context, rr ...*types.AuthOa2token) (err error) {
	for _, res := range rr {
		err = s.checkAuthOa2tokenConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execCreateAuthOa2tokens(ctx, s.internalAuthOa2tokenEncoder(res))
		if err != nil {
			return err
		}
	}

	return
}

// DeleteAuthOa2token Deletes one or more rows from auth_oa2tokens table
func (s Store) DeleteAuthOa2token(ctx context.Context, rr ...*types.AuthOa2token) (err error) {
	for _, res := range rr {

		err = s.execDeleteAuthOa2tokens(ctx, squirrel.Eq{
			s.preprocessColumn("tkn.id", ""): store.PreprocessValue(res.ID, ""),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteAuthOa2tokenByID Deletes row from the auth_oa2tokens table
func (s Store) DeleteAuthOa2tokenByID(ctx context.Context, ID uint64) error {
	return s.execDeleteAuthOa2tokens(ctx, squirrel.Eq{
		s.preprocessColumn("tkn.id", ""): store.PreprocessValue(ID, ""),
	})
}

// TruncateAuthOa2tokens Deletes all rows from the auth_oa2tokens table
func (s Store) TruncateAuthOa2tokens(ctx context.Context) error {
	return s.Truncate(ctx, s.authOa2tokenTable())
}

// execLookupAuthOa2token prepares AuthOa2token query and executes it,
// returning types.AuthOa2token (or error)
func (s Store) execLookupAuthOa2token(ctx context.Context, cnd squirrel.Sqlizer) (res *types.AuthOa2token, err error) {
	var (
		row rowScanner
	)

	row, err = s.QueryRow(ctx, s.authOa2tokensSelectBuilder().Where(cnd))
	if err != nil {
		return
	}

	res, err = s.internalAuthOa2tokenRowScanner(row)
	if err != nil {
		return
	}

	return res, nil
}

// execCreateAuthOa2tokens updates all matched (by cnd) rows in auth_oa2tokens with given data
func (s Store) execCreateAuthOa2tokens(ctx context.Context, payload store.Payload) error {
	return s.Exec(ctx, s.InsertBuilder(s.authOa2tokenTable()).SetMap(payload))
}

// execDeleteAuthOa2tokens Deletes all matched (by cnd) rows in auth_oa2tokens with given data
func (s Store) execDeleteAuthOa2tokens(ctx context.Context, cnd squirrel.Sqlizer) error {
	return s.Exec(ctx, s.DeleteBuilder(s.authOa2tokenTable("tkn")).Where(cnd))
}

func (s Store) internalAuthOa2tokenRowScanner(row rowScanner) (res *types.AuthOa2token, err error) {
	res = &types.AuthOa2token{}

	if _, has := s.config.RowScanners["authOa2token"]; has {
		scanner := s.config.RowScanners["authOa2token"].(func(_ rowScanner, _ *types.AuthOa2token) error)
		err = scanner(row, res)
	} else {
		err = row.Scan(
			&res.ID,
			&res.Code,
			&res.Access,
			&res.Refresh,
			&res.ExpiresAt,
			&res.CreatedAt,
			&res.Data,
			&res.ClientID,
			&res.UserID,
			&res.RemoteAddr,
			&res.UserAgent,
		)
	}

	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err != nil {
		return nil, errors.Store("could not scan authOa2token db row: %s", err).Wrap(err)
	} else {
		return res, nil
	}
}

// QueryAuthOa2tokens returns squirrel.SelectBuilder with set table and all columns
func (s Store) authOa2tokensSelectBuilder() squirrel.SelectBuilder {
	return s.SelectBuilder(s.authOa2tokenTable("tkn"), s.authOa2tokenColumns("tkn")...)
}

// authOa2tokenTable name of the db table
func (Store) authOa2tokenTable(aa ...string) string {
	var alias string
	if len(aa) > 0 {
		alias = " AS " + aa[0]
	}

	return "auth_oa2tokens" + alias
}

// AuthOa2tokenColumns returns all defined table columns
//
// With optional string arg, all columns are returned aliased
func (Store) authOa2tokenColumns(aa ...string) []string {
	var alias string
	if len(aa) > 0 {
		alias = aa[0] + "."
	}

	return []string{
		alias + "id",
		alias + "code",
		alias + "access",
		alias + "refresh",
		alias + "expires_at",
		alias + "created_at",
		alias + "data",
		alias + "rel_client",
		alias + "rel_user",
		alias + "remote_addr",
		alias + "user_agent",
	}
}

// {true true false false false false}

// internalAuthOa2tokenEncoder encodes fields from types.AuthOa2token to store.Payload (map)
//
// Encoding is done by using generic approach or by calling encodeAuthOa2token
// func when rdbms.customEncoder=true
func (s Store) internalAuthOa2tokenEncoder(res *types.AuthOa2token) store.Payload {
	return store.Payload{
		"id":          res.ID,
		"code":        res.Code,
		"access":      res.Access,
		"refresh":     res.Refresh,
		"expires_at":  res.ExpiresAt,
		"created_at":  res.CreatedAt,
		"data":        res.Data,
		"rel_client":  res.ClientID,
		"rel_user":    res.UserID,
		"remote_addr": res.RemoteAddr,
		"user_agent":  res.UserAgent,
	}
}

// checkAuthOa2tokenConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant but unfortunately we can not rely
// on the full support (MySQL does not support conditional indexes)
func (s *Store) checkAuthOa2tokenConstraints(ctx context.Context, res *types.AuthOa2token) error {
	// Consider resource valid when all fields in unique constraint check lookups
	// have valid (non-empty) value
	//
	// Only string and uint64 are supported for now
	// feel free to add additional types if needed
	var valid = true

	valid = valid && len(res.Code) > 0

	valid = valid && len(res.Access) > 0

	valid = valid && len(res.Refresh) > 0

	if !valid {
		return nil
	}

	{
		ex, err := s.LookupAuthOa2tokenByCode(ctx, res.Code)
		if err == nil && ex != nil && ex.ID != res.ID {
			return store.ErrNotUnique.Stack(1)
		} else if !errors.IsNotFound(err) {
			return err
		}
	}

	{
		ex, err := s.LookupAuthOa2tokenByAccess(ctx, res.Access)
		if err == nil && ex != nil && ex.ID != res.ID {
			return store.ErrNotUnique.Stack(1)
		} else if !errors.IsNotFound(err) {
			return err
		}
	}

	{
		ex, err := s.LookupAuthOa2tokenByRefresh(ctx, res.Refresh)
		if err == nil && ex != nil && ex.ID != res.ID {
			return store.ErrNotUnique.Stack(1)
		} else if !errors.IsNotFound(err) {
			return err
		}
	}

	return nil
}
