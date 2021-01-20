package rdbms

// This file is an auto-generated file
//
// Template:    pkg/codegen/assets/store_rdbms.gen.go.tpl
// Definitions: store/auth_sessions.yaml
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

// SearchAuthSessions returns all matching rows
//
// This function calls convertAuthSessionFilter with the given
// types.AuthSessionFilter and expects to receive a working squirrel.SelectBuilder
func (s Store) SearchAuthSessions(ctx context.Context, f types.AuthSessionFilter) (types.AuthSessionSet, types.AuthSessionFilter, error) {
	var (
		err error
		set []*types.AuthSession
		q   squirrel.SelectBuilder
	)

	return set, f, func() error {
		q, err = s.convertAuthSessionFilter(f)
		if err != nil {
			return err
		}

		set, err = s.QueryAuthSessions(ctx, q, nil)
		return err
	}()
}

// QueryAuthSessions queries the database, converts and checks each row and
// returns collected set
//
// Fn also returns total number of fetched items and last fetched item so that the caller can construct cursor
// for next page of results
func (s Store) QueryAuthSessions(
	ctx context.Context,
	q squirrel.Sqlizer,
	check func(*types.AuthSession) (bool, error),
) ([]*types.AuthSession, error) {
	var (
		set = make([]*types.AuthSession, 0, DefaultSliceCapacity)
		res *types.AuthSession

		// Query rows with
		rows, err = s.Query(ctx, q)
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		if err = rows.Err(); err == nil {
			res, err = s.internalAuthSessionRowScanner(rows)
		}

		if err != nil {
			return nil, err
		}

		set = append(set, res)
	}

	return set, rows.Err()
}

// LookupAuthSessionByID
func (s Store) LookupAuthSessionByID(ctx context.Context, id string) (*types.AuthSession, error) {
	return s.execLookupAuthSession(ctx, squirrel.Eq{
		s.preprocessColumn("ses.id", ""): store.PreprocessValue(id, ""),
	})
}

// CreateAuthSession creates one or more rows in auth_sessions table
func (s Store) CreateAuthSession(ctx context.Context, rr ...*types.AuthSession) (err error) {
	for _, res := range rr {
		err = s.checkAuthSessionConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execCreateAuthSessions(ctx, s.internalAuthSessionEncoder(res))
		if err != nil {
			return err
		}
	}

	return
}

// UpdateAuthSession updates one or more existing rows in auth_sessions
func (s Store) UpdateAuthSession(ctx context.Context, rr ...*types.AuthSession) error {
	return s.partialAuthSessionUpdate(ctx, nil, rr...)
}

// partialAuthSessionUpdate updates one or more existing rows in auth_sessions
func (s Store) partialAuthSessionUpdate(ctx context.Context, onlyColumns []string, rr ...*types.AuthSession) (err error) {
	for _, res := range rr {
		err = s.checkAuthSessionConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpdateAuthSessions(
			ctx,
			squirrel.Eq{
				s.preprocessColumn("ses.id", ""): store.PreprocessValue(res.ID, ""),
			},
			s.internalAuthSessionEncoder(res).Skip("id").Only(onlyColumns...))
		if err != nil {
			return err
		}
	}

	return
}

// UpsertAuthSession updates one or more existing rows in auth_sessions
func (s Store) UpsertAuthSession(ctx context.Context, rr ...*types.AuthSession) (err error) {
	for _, res := range rr {
		err = s.checkAuthSessionConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpsertAuthSessions(ctx, s.internalAuthSessionEncoder(res))
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteAuthSession Deletes one or more rows from auth_sessions table
func (s Store) DeleteAuthSession(ctx context.Context, rr ...*types.AuthSession) (err error) {
	for _, res := range rr {

		err = s.execDeleteAuthSessions(ctx, squirrel.Eq{
			s.preprocessColumn("ses.id", ""): store.PreprocessValue(res.ID, ""),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteAuthSessionByID Deletes row from the auth_sessions table
func (s Store) DeleteAuthSessionByID(ctx context.Context, ID string) error {
	return s.execDeleteAuthSessions(ctx, squirrel.Eq{
		s.preprocessColumn("ses.id", ""): store.PreprocessValue(ID, ""),
	})
}

// TruncateAuthSessions Deletes all rows from the auth_sessions table
func (s Store) TruncateAuthSessions(ctx context.Context) error {
	return s.Truncate(ctx, s.authSessionTable())
}

// execLookupAuthSession prepares AuthSession query and executes it,
// returning types.AuthSession (or error)
func (s Store) execLookupAuthSession(ctx context.Context, cnd squirrel.Sqlizer) (res *types.AuthSession, err error) {
	var (
		row rowScanner
	)

	row, err = s.QueryRow(ctx, s.authSessionsSelectBuilder().Where(cnd))
	if err != nil {
		return
	}

	res, err = s.internalAuthSessionRowScanner(row)
	if err != nil {
		return
	}

	return res, nil
}

// execCreateAuthSessions updates all matched (by cnd) rows in auth_sessions with given data
func (s Store) execCreateAuthSessions(ctx context.Context, payload store.Payload) error {
	return s.Exec(ctx, s.InsertBuilder(s.authSessionTable()).SetMap(payload))
}

// execUpdateAuthSessions updates all matched (by cnd) rows in auth_sessions with given data
func (s Store) execUpdateAuthSessions(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return s.Exec(ctx, s.UpdateBuilder(s.authSessionTable("ses")).Where(cnd).SetMap(set))
}

// execUpsertAuthSessions inserts new or updates matching (by-primary-key) rows in auth_sessions with given data
func (s Store) execUpsertAuthSessions(ctx context.Context, set store.Payload) error {
	upsert, err := s.config.UpsertBuilder(
		s.config,
		s.authSessionTable(),
		set,
		s.preprocessColumn("id", ""),
	)

	if err != nil {
		return err
	}

	return s.Exec(ctx, upsert)
}

// execDeleteAuthSessions Deletes all matched (by cnd) rows in auth_sessions with given data
func (s Store) execDeleteAuthSessions(ctx context.Context, cnd squirrel.Sqlizer) error {
	return s.Exec(ctx, s.DeleteBuilder(s.authSessionTable("ses")).Where(cnd))
}

func (s Store) internalAuthSessionRowScanner(row rowScanner) (res *types.AuthSession, err error) {
	res = &types.AuthSession{}

	if _, has := s.config.RowScanners["authSession"]; has {
		scanner := s.config.RowScanners["authSession"].(func(_ rowScanner, _ *types.AuthSession) error)
		err = scanner(row, res)
	} else {
		err = row.Scan(
			&res.ID,
			&res.Data,
			&res.UserID,
			&res.RemoteAddr,
			&res.UserAgent,
			&res.CreatedAt,
			&res.ExpiresAt,
		)
	}

	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err != nil {
		return nil, errors.Store("could not scan authSession db row: %s", err).Wrap(err)
	} else {
		return res, nil
	}
}

// QueryAuthSessions returns squirrel.SelectBuilder with set table and all columns
func (s Store) authSessionsSelectBuilder() squirrel.SelectBuilder {
	return s.SelectBuilder(s.authSessionTable("ses"), s.authSessionColumns("ses")...)
}

// authSessionTable name of the db table
func (Store) authSessionTable(aa ...string) string {
	var alias string
	if len(aa) > 0 {
		alias = " AS " + aa[0]
	}

	return "auth_sessions" + alias
}

// AuthSessionColumns returns all defined table columns
//
// With optional string arg, all columns are returned aliased
func (Store) authSessionColumns(aa ...string) []string {
	var alias string
	if len(aa) > 0 {
		alias = aa[0] + "."
	}

	return []string{
		alias + "id",
		alias + "data",
		alias + "rel_user",
		alias + "remote_addr",
		alias + "user_agent",
		alias + "created_at",
		alias + "expires_at",
	}
}

// {true true false false false false}

// internalAuthSessionEncoder encodes fields from types.AuthSession to store.Payload (map)
//
// Encoding is done by using generic approach or by calling encodeAuthSession
// func when rdbms.customEncoder=true
func (s Store) internalAuthSessionEncoder(res *types.AuthSession) store.Payload {
	return store.Payload{
		"id":          res.ID,
		"data":        res.Data,
		"rel_user":    res.UserID,
		"remote_addr": res.RemoteAddr,
		"user_agent":  res.UserAgent,
		"created_at":  res.CreatedAt,
		"expires_at":  res.ExpiresAt,
	}
}

// checkAuthSessionConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant but unfortunately we can not rely
// on the full support (MySQL does not support conditional indexes)
func (s *Store) checkAuthSessionConstraints(ctx context.Context, res *types.AuthSession) error {
	// Consider resource valid when all fields in unique constraint check lookups
	// have valid (non-empty) value
	//
	// Only string and uint64 are supported for now
	// feel free to add additional types if needed
	var valid = true

	valid = valid && len(res.ID) > 0

	if !valid {
		return nil
	}

	{
		ex, err := s.LookupAuthSessionByID(ctx, res.ID)
		if err == nil && ex != nil && ex.ID != res.ID {
			return store.ErrNotUnique.Stack(1)
		} else if !errors.IsNotFound(err) {
			return err
		}
	}

	return nil
}
