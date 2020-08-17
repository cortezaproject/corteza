package rdbms

// This file is an auto-generated file
//
// Template:    pkg/codegen/assets/store_rdbms.gen.go.tpl
// Definitions: store/credentials.yaml
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

// SearchCredentials returns all matching rows
//
// This function calls convertCredentialsFilter with the given
// types.CredentialsFilter and expects to receive a working squirrel.SelectBuilder
func (s Store) SearchCredentials(ctx context.Context, f types.CredentialsFilter) (types.CredentialsSet, types.CredentialsFilter, error) {
	q, err := s.convertCredentialsFilter(f)
	if err != nil {
		return nil, f, err
	}

	scap := DefaultSliceCapacity

	var (
		set   = make([]*types.Credentials, 0, scap)
		fetch = func() error {
			var (
				res       *types.Credentials
				rows, err = s.Query(ctx, q)
			)

			if err != nil {
				return err
			}

			for rows.Next() {
				if res, err = s.internalCredentialsRowScanner(rows, rows.Err()); err != nil {
					if cerr := rows.Close(); cerr != nil {
						return fmt.Errorf("could not close rows (%v) after scan error: %w", cerr, err)
					}

					return err
				}

				set = append(set, res)
			}

			return rows.Close()
		}
	)

	return set, f, fetch()
}

// LookupCredentialsByID searches for credentials by ID
//
// It returns credentials even if deleted
func (s Store) LookupCredentialsByID(ctx context.Context, id uint64) (*types.Credentials, error) {
	return s.CredentialsLookup(ctx, squirrel.Eq{
		"crd.id": id,
	})
}

// CreateCredentials creates one or more rows in credentials table
func (s Store) CreateCredentials(ctx context.Context, rr ...*types.Credentials) error {
	if len(rr) == 0 {
		return nil
	}

	return Tx(ctx, s.db, s.config, nil, func(db *sqlx.Tx) (err error) {
		for _, res := range rr {
			err = ExecuteSqlizer(ctx, s.DB(), s.Insert(s.CredentialsTable()).SetMap(s.internalCredentialsEncoder(res)))
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// UpdateCredentials updates one or more existing rows in credentials
func (s Store) UpdateCredentials(ctx context.Context, rr ...*types.Credentials) error {
	return s.PartialUpdateCredentials(ctx, nil, rr...)
}

// PartialUpdateCredentials updates one or more existing rows in credentials
//
// It wraps the update into transaction and can perform partial update by providing list of updatable columns
func (s Store) PartialUpdateCredentials(ctx context.Context, onlyColumns []string, rr ...*types.Credentials) error {
	if len(rr) == 0 {
		return nil
	}

	return Tx(ctx, s.db, s.config, nil, func(db *sqlx.Tx) (err error) {
		for _, res := range rr {
			err = s.ExecUpdateCredentials(
				ctx,
				squirrel.Eq{s.preprocessColumn("crd.id", ""): s.preprocessValue(res.ID, "")},
				s.internalCredentialsEncoder(res).Skip("id").Only(onlyColumns...))
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// RemoveCredentials removes one or more rows from credentials table
func (s Store) RemoveCredentials(ctx context.Context, rr ...*types.Credentials) error {
	if len(rr) == 0 {
		return nil
	}

	return Tx(ctx, s.db, s.config, nil, func(db *sqlx.Tx) (err error) {
		for _, res := range rr {
			err = ExecuteSqlizer(ctx, s.DB(), s.Delete(s.CredentialsTable("crd")).Where(squirrel.Eq{s.preprocessColumn("crd.id", ""): s.preprocessValue(res.ID, "")}))
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// RemoveCredentialsByID removes row from the credentials table
func (s Store) RemoveCredentialsByID(ctx context.Context, ID uint64) error {
	return ExecuteSqlizer(ctx, s.DB(), s.Delete(s.CredentialsTable("crd")).Where(squirrel.Eq{s.preprocessColumn("crd.id", ""): s.preprocessValue(ID, "")}))
}

// TruncateCredentials removes all rows from the credentials table
func (s Store) TruncateCredentials(ctx context.Context) error {
	return Truncate(ctx, s.DB(), s.CredentialsTable())
}

// ExecUpdateCredentials updates all matched (by cnd) rows in credentials with given data
func (s Store) ExecUpdateCredentials(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return ExecuteSqlizer(ctx, s.DB(), s.Update(s.CredentialsTable("crd")).Where(cnd).SetMap(set))
}

// CredentialsLookup prepares Credentials query and executes it,
// returning types.Credentials (or error)
func (s Store) CredentialsLookup(ctx context.Context, cnd squirrel.Sqlizer) (*types.Credentials, error) {
	return s.internalCredentialsRowScanner(s.QueryRow(ctx, s.QueryCredentials().Where(cnd)))
}

func (s Store) internalCredentialsRowScanner(row rowScanner, err error) (*types.Credentials, error) {
	if err != nil {
		return nil, err
	}

	var res = &types.Credentials{}
	if _, has := s.config.RowScanners["credentials"]; has {
		scanner := s.config.RowScanners["credentials"].(func(rowScanner, *types.Credentials) error)
		err = scanner(row, res)
	} else {
		err = row.Scan(
			&res.ID,
			&res.OwnerID,
			&res.Kind,
			&res.Label,
			&res.Credentials,
			&res.Meta,
			&res.LastUsedAt,
			&res.ExpiresAt,
			&res.CreatedAt,
			&res.UpdatedAt,
			&res.DeletedAt,
		)
	}

	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("could not scan db row for Credentials: %w", err)
	} else {
		return res, nil
	}
}

// QueryCredentials returns squirrel.SelectBuilder with set table and all columns
func (s Store) QueryCredentials() squirrel.SelectBuilder {
	return s.Select(s.CredentialsTable("crd"), s.CredentialsColumns("crd")...)
}

// CredentialsTable name of the db table
func (Store) CredentialsTable(aa ...string) string {
	var alias string
	if len(aa) > 0 {
		alias = " AS " + aa[0]
	}

	return "credentials" + alias
}

// CredentialsColumns returns all defined table columns
//
// With optional string arg, all columns are returned aliased
func (Store) CredentialsColumns(aa ...string) []string {
	var alias string
	if len(aa) > 0 {
		alias = aa[0] + "."
	}

	return []string{
		alias + "id",
		alias + "rel_owner",
		alias + "kind",
		alias + "label",
		alias + "credentials",
		alias + "meta",
		alias + "last_used_at",
		alias + "expires_at",
		alias + "created_at",
		alias + "updated_at",
		alias + "deleted_at",
	}
}

// internalCredentialsEncoder encodes fields from types.Credentials to store.Payload (map)
//
// Encoding is done by using generic approach or by calling encodeCredentials
// func when rdbms.customEncoder=true
func (s Store) internalCredentialsEncoder(res *types.Credentials) store.Payload {
	return store.Payload{
		"id":           res.ID,
		"rel_owner":    res.OwnerID,
		"kind":         res.Kind,
		"label":        res.Label,
		"credentials":  res.Credentials,
		"meta":         res.Meta,
		"last_used_at": res.LastUsedAt,
		"expires_at":   res.ExpiresAt,
		"created_at":   res.CreatedAt,
		"updated_at":   res.UpdatedAt,
		"deleted_at":   res.DeletedAt,
	}
}
