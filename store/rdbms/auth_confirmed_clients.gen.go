package rdbms

// This file is an auto-generated file
//
// Template:    pkg/codegen/assets/store_rdbms.gen.go.tpl
// Definitions: store/auth_confirmed_clients.yaml
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

// SearchAuthConfirmedClients returns all matching rows
//
// This function calls convertAuthConfirmedClientFilter with the given
// types.AuthConfirmedClientFilter and expects to receive a working squirrel.SelectBuilder
func (s Store) SearchAuthConfirmedClients(ctx context.Context, f types.AuthConfirmedClientFilter) (types.AuthConfirmedClientSet, types.AuthConfirmedClientFilter, error) {
	var (
		err error
		set []*types.AuthConfirmedClient
		q   squirrel.SelectBuilder
	)

	return set, f, func() error {
		q, err = s.convertAuthConfirmedClientFilter(f)
		if err != nil {
			return err
		}

		set, err = s.QueryAuthConfirmedClients(ctx, q, nil)
		return err
	}()
}

// QueryAuthConfirmedClients queries the database, converts and checks each row and
// returns collected set
//
// Fn also returns total number of fetched items and last fetched item so that the caller can construct cursor
// for next page of results
func (s Store) QueryAuthConfirmedClients(
	ctx context.Context,
	q squirrel.Sqlizer,
	check func(*types.AuthConfirmedClient) (bool, error),
) ([]*types.AuthConfirmedClient, error) {
	var (
		set = make([]*types.AuthConfirmedClient, 0, DefaultSliceCapacity)
		res *types.AuthConfirmedClient

		// Query rows with
		rows, err = s.Query(ctx, q)
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		if err = rows.Err(); err == nil {
			res, err = s.internalAuthConfirmedClientRowScanner(rows)
		}

		if err != nil {
			return nil, err
		}

		set = append(set, res)
	}

	return set, rows.Err()
}

// LookupAuthConfirmedClientByUserIDClientID
func (s Store) LookupAuthConfirmedClientByUserIDClientID(ctx context.Context, user_id uint64, client_id uint64) (*types.AuthConfirmedClient, error) {
	return s.execLookupAuthConfirmedClient(ctx, squirrel.Eq{
		s.preprocessColumn("acc.rel_user", ""):   store.PreprocessValue(user_id, ""),
		s.preprocessColumn("acc.rel_client", ""): store.PreprocessValue(client_id, ""),
	})
}

// CreateAuthConfirmedClient creates one or more rows in auth_confirmed_clients table
func (s Store) CreateAuthConfirmedClient(ctx context.Context, rr ...*types.AuthConfirmedClient) (err error) {
	for _, res := range rr {
		err = s.checkAuthConfirmedClientConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execCreateAuthConfirmedClients(ctx, s.internalAuthConfirmedClientEncoder(res))
		if err != nil {
			return err
		}
	}

	return
}

// UpdateAuthConfirmedClient updates one or more existing rows in auth_confirmed_clients
func (s Store) UpdateAuthConfirmedClient(ctx context.Context, rr ...*types.AuthConfirmedClient) error {
	return s.partialAuthConfirmedClientUpdate(ctx, nil, rr...)
}

// partialAuthConfirmedClientUpdate updates one or more existing rows in auth_confirmed_clients
func (s Store) partialAuthConfirmedClientUpdate(ctx context.Context, onlyColumns []string, rr ...*types.AuthConfirmedClient) (err error) {
	for _, res := range rr {
		err = s.checkAuthConfirmedClientConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpdateAuthConfirmedClients(
			ctx,
			squirrel.Eq{
				s.preprocessColumn("acc.rel_user", ""): store.PreprocessValue(res.UserID, ""), s.preprocessColumn("acc.rel_client", ""): store.PreprocessValue(res.ClientID, ""),
			},
			s.internalAuthConfirmedClientEncoder(res).Skip("rel_user", "rel_client").Only(onlyColumns...))
		if err != nil {
			return err
		}
	}

	return
}

// UpsertAuthConfirmedClient updates one or more existing rows in auth_confirmed_clients
func (s Store) UpsertAuthConfirmedClient(ctx context.Context, rr ...*types.AuthConfirmedClient) (err error) {
	for _, res := range rr {
		err = s.checkAuthConfirmedClientConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpsertAuthConfirmedClients(ctx, s.internalAuthConfirmedClientEncoder(res))
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteAuthConfirmedClient Deletes one or more rows from auth_confirmed_clients table
func (s Store) DeleteAuthConfirmedClient(ctx context.Context, rr ...*types.AuthConfirmedClient) (err error) {
	for _, res := range rr {

		err = s.execDeleteAuthConfirmedClients(ctx, squirrel.Eq{
			s.preprocessColumn("acc.rel_user", ""): store.PreprocessValue(res.UserID, ""), s.preprocessColumn("acc.rel_client", ""): store.PreprocessValue(res.ClientID, ""),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteAuthConfirmedClientByUserIDClientID Deletes row from the auth_confirmed_clients table
func (s Store) DeleteAuthConfirmedClientByUserIDClientID(ctx context.Context, userID uint64, clientID uint64) error {
	return s.execDeleteAuthConfirmedClients(ctx, squirrel.Eq{
		s.preprocessColumn("acc.rel_user", ""):   store.PreprocessValue(userID, ""),
		s.preprocessColumn("acc.rel_client", ""): store.PreprocessValue(clientID, ""),
	})
}

// TruncateAuthConfirmedClients Deletes all rows from the auth_confirmed_clients table
func (s Store) TruncateAuthConfirmedClients(ctx context.Context) error {
	return s.Truncate(ctx, s.authConfirmedClientTable())
}

// execLookupAuthConfirmedClient prepares AuthConfirmedClient query and executes it,
// returning types.AuthConfirmedClient (or error)
func (s Store) execLookupAuthConfirmedClient(ctx context.Context, cnd squirrel.Sqlizer) (res *types.AuthConfirmedClient, err error) {
	var (
		row rowScanner
	)

	row, err = s.QueryRow(ctx, s.authConfirmedClientsSelectBuilder().Where(cnd))
	if err != nil {
		return
	}

	res, err = s.internalAuthConfirmedClientRowScanner(row)
	if err != nil {
		return
	}

	return res, nil
}

// execCreateAuthConfirmedClients updates all matched (by cnd) rows in auth_confirmed_clients with given data
func (s Store) execCreateAuthConfirmedClients(ctx context.Context, payload store.Payload) error {
	return s.Exec(ctx, s.InsertBuilder(s.authConfirmedClientTable()).SetMap(payload))
}

// execUpdateAuthConfirmedClients updates all matched (by cnd) rows in auth_confirmed_clients with given data
func (s Store) execUpdateAuthConfirmedClients(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return s.Exec(ctx, s.UpdateBuilder(s.authConfirmedClientTable("acc")).Where(cnd).SetMap(set))
}

// execUpsertAuthConfirmedClients inserts new or updates matching (by-primary-key) rows in auth_confirmed_clients with given data
func (s Store) execUpsertAuthConfirmedClients(ctx context.Context, set store.Payload) error {
	upsert, err := s.config.UpsertBuilder(
		s.config,
		s.authConfirmedClientTable(),
		set,
		s.preprocessColumn("rel_user", ""),
		s.preprocessColumn("rel_client", ""),
	)

	if err != nil {
		return err
	}

	return s.Exec(ctx, upsert)
}

// execDeleteAuthConfirmedClients Deletes all matched (by cnd) rows in auth_confirmed_clients with given data
func (s Store) execDeleteAuthConfirmedClients(ctx context.Context, cnd squirrel.Sqlizer) error {
	return s.Exec(ctx, s.DeleteBuilder(s.authConfirmedClientTable("acc")).Where(cnd))
}

func (s Store) internalAuthConfirmedClientRowScanner(row rowScanner) (res *types.AuthConfirmedClient, err error) {
	res = &types.AuthConfirmedClient{}

	if _, has := s.config.RowScanners["authConfirmedClient"]; has {
		scanner := s.config.RowScanners["authConfirmedClient"].(func(_ rowScanner, _ *types.AuthConfirmedClient) error)
		err = scanner(row, res)
	} else {
		err = row.Scan(
			&res.UserID,
			&res.ClientID,
			&res.ConfirmedAt,
		)
	}

	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err != nil {
		return nil, errors.Store("could not scan authConfirmedClient db row: %s", err).Wrap(err)
	} else {
		return res, nil
	}
}

// QueryAuthConfirmedClients returns squirrel.SelectBuilder with set table and all columns
func (s Store) authConfirmedClientsSelectBuilder() squirrel.SelectBuilder {
	return s.SelectBuilder(s.authConfirmedClientTable("acc"), s.authConfirmedClientColumns("acc")...)
}

// authConfirmedClientTable name of the db table
func (Store) authConfirmedClientTable(aa ...string) string {
	var alias string
	if len(aa) > 0 {
		alias = " AS " + aa[0]
	}

	return "auth_confirmed_clients" + alias
}

// AuthConfirmedClientColumns returns all defined table columns
//
// With optional string arg, all columns are returned aliased
func (Store) authConfirmedClientColumns(aa ...string) []string {
	var alias string
	if len(aa) > 0 {
		alias = aa[0] + "."
	}

	return []string{
		alias + "rel_user",
		alias + "rel_client",
		alias + "confirmed_at",
	}
}

// {true true false false false false}

// internalAuthConfirmedClientEncoder encodes fields from types.AuthConfirmedClient to store.Payload (map)
//
// Encoding is done by using generic approach or by calling encodeAuthConfirmedClient
// func when rdbms.customEncoder=true
func (s Store) internalAuthConfirmedClientEncoder(res *types.AuthConfirmedClient) store.Payload {
	return store.Payload{
		"rel_user":     res.UserID,
		"rel_client":   res.ClientID,
		"confirmed_at": res.ConfirmedAt,
	}
}

// checkAuthConfirmedClientConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant but unfortunately we can not rely
// on the full support (MySQL does not support conditional indexes)
func (s *Store) checkAuthConfirmedClientConstraints(ctx context.Context, res *types.AuthConfirmedClient) error {
	// Consider resource valid when all fields in unique constraint check lookups
	// have valid (non-empty) value
	//
	// Only string and uint64 are supported for now
	// feel free to add additional types if needed
	var valid = true

	valid = valid && res.UserID > 0

	valid = valid && res.ClientID > 0

	if !valid {
		return nil
	}

	{
		ex, err := s.LookupAuthConfirmedClientByUserIDClientID(ctx, res.UserID, res.ClientID)
		if err == nil && ex != nil && ex.UserID != res.UserID && ex.ClientID != res.ClientID {
			return store.ErrNotUnique.Stack(1)
		} else if !errors.IsNotFound(err) {
			return err
		}
	}

	return nil
}
