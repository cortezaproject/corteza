package rdbms

// This file is an auto-generated file
//
// Template:    pkg/codegen/assets/store_rdbms.gen.go.tpl
// Definitions: store/messaging_flags.yaml
//
// Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated.

import (
	"context"
	"database/sql"
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/store"
)

var _ = errors.Is

// SearchMessagingFlags returns all matching rows
//
// This function calls convertMessagingFlagFilter with the given
// types.MessageFlagFilter and expects to receive a working squirrel.SelectBuilder
func (s Store) SearchMessagingFlags(ctx context.Context, f types.MessageFlagFilter) (types.MessageFlagSet, types.MessageFlagFilter, error) {
	var (
		err error
		set []*types.MessageFlag
		q   squirrel.SelectBuilder
	)

	return set, f, func() error {
		q, err = s.convertMessagingFlagFilter(f)
		if err != nil {
			return err
		}

		set, err = s.QueryMessagingFlags(ctx, q, nil)
		return err
	}()
}

// QueryMessagingFlags queries the database, converts and checks each row and
// returns collected set
//
// Fn also returns total number of fetched items and last fetched item so that the caller can construct cursor
// for next page of results
func (s Store) QueryMessagingFlags(
	ctx context.Context,
	q squirrel.Sqlizer,
	check func(*types.MessageFlag) (bool, error),
) ([]*types.MessageFlag, error) {
	var (
		set = make([]*types.MessageFlag, 0, DefaultSliceCapacity)
		res *types.MessageFlag

		// Query rows with
		rows, err = s.Query(ctx, q)
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		if err = rows.Err(); err == nil {
			res, err = s.internalMessagingFlagRowScanner(rows)
		}

		if err != nil {
			return nil, err
		}

		set = append(set, res)
	}

	return set, rows.Err()
}

// LookupMessagingFlagByID searches for flags by ID
func (s Store) LookupMessagingFlagByID(ctx context.Context, id uint64) (*types.MessageFlag, error) {
	return s.execLookupMessagingFlag(ctx, squirrel.Eq{
		s.preprocessColumn("mmf.id", ""): store.PreprocessValue(id, ""),
	})
}

// CreateMessagingFlag creates one or more rows in messaging_message_flag table
func (s Store) CreateMessagingFlag(ctx context.Context, rr ...*types.MessageFlag) (err error) {
	for _, res := range rr {
		err = s.checkMessagingFlagConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execCreateMessagingFlags(ctx, s.internalMessagingFlagEncoder(res))
		if err != nil {
			return err
		}
	}

	return
}

// UpdateMessagingFlag updates one or more existing rows in messaging_message_flag
func (s Store) UpdateMessagingFlag(ctx context.Context, rr ...*types.MessageFlag) error {
	return s.partialMessagingFlagUpdate(ctx, nil, rr...)
}

// partialMessagingFlagUpdate updates one or more existing rows in messaging_message_flag
func (s Store) partialMessagingFlagUpdate(ctx context.Context, onlyColumns []string, rr ...*types.MessageFlag) (err error) {
	for _, res := range rr {
		err = s.checkMessagingFlagConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpdateMessagingFlags(
			ctx,
			squirrel.Eq{
				s.preprocessColumn("mmf.id", ""): store.PreprocessValue(res.ID, ""),
			},
			s.internalMessagingFlagEncoder(res).Skip("id").Only(onlyColumns...))
		if err != nil {
			return err
		}
	}

	return
}

// UpsertMessagingFlag updates one or more existing rows in messaging_message_flag
func (s Store) UpsertMessagingFlag(ctx context.Context, rr ...*types.MessageFlag) (err error) {
	for _, res := range rr {
		err = s.checkMessagingFlagConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpsertMessagingFlags(ctx, s.internalMessagingFlagEncoder(res))
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteMessagingFlag Deletes one or more rows from messaging_message_flag table
func (s Store) DeleteMessagingFlag(ctx context.Context, rr ...*types.MessageFlag) (err error) {
	for _, res := range rr {

		err = s.execDeleteMessagingFlags(ctx, squirrel.Eq{
			s.preprocessColumn("mmf.id", ""): store.PreprocessValue(res.ID, ""),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteMessagingFlagByID Deletes row from the messaging_message_flag table
func (s Store) DeleteMessagingFlagByID(ctx context.Context, ID uint64) error {
	return s.execDeleteMessagingFlags(ctx, squirrel.Eq{
		s.preprocessColumn("mmf.id", ""): store.PreprocessValue(ID, ""),
	})
}

// TruncateMessagingFlags Deletes all rows from the messaging_message_flag table
func (s Store) TruncateMessagingFlags(ctx context.Context) error {
	return s.Truncate(ctx, s.messagingFlagTable())
}

// execLookupMessagingFlag prepares MessagingFlag query and executes it,
// returning types.MessageFlag (or error)
func (s Store) execLookupMessagingFlag(ctx context.Context, cnd squirrel.Sqlizer) (res *types.MessageFlag, err error) {
	var (
		row rowScanner
	)

	row, err = s.QueryRow(ctx, s.messagingFlagsSelectBuilder().Where(cnd))
	if err != nil {
		return
	}

	res, err = s.internalMessagingFlagRowScanner(row)
	if err != nil {
		return
	}

	return res, nil
}

// execCreateMessagingFlags updates all matched (by cnd) rows in messaging_message_flag with given data
func (s Store) execCreateMessagingFlags(ctx context.Context, payload store.Payload) error {
	return s.Exec(ctx, s.InsertBuilder(s.messagingFlagTable()).SetMap(payload))
}

// execUpdateMessagingFlags updates all matched (by cnd) rows in messaging_message_flag with given data
func (s Store) execUpdateMessagingFlags(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return s.Exec(ctx, s.UpdateBuilder(s.messagingFlagTable("mmf")).Where(cnd).SetMap(set))
}

// execUpsertMessagingFlags inserts new or updates matching (by-primary-key) rows in messaging_message_flag with given data
func (s Store) execUpsertMessagingFlags(ctx context.Context, set store.Payload) error {
	upsert, err := s.config.UpsertBuilder(
		s.config,
		s.messagingFlagTable(),
		set,
		s.preprocessColumn("id", ""),
	)

	if err != nil {
		return err
	}

	return s.Exec(ctx, upsert)
}

// execDeleteMessagingFlags Deletes all matched (by cnd) rows in messaging_message_flag with given data
func (s Store) execDeleteMessagingFlags(ctx context.Context, cnd squirrel.Sqlizer) error {
	return s.Exec(ctx, s.DeleteBuilder(s.messagingFlagTable("mmf")).Where(cnd))
}

func (s Store) internalMessagingFlagRowScanner(row rowScanner) (res *types.MessageFlag, err error) {
	res = &types.MessageFlag{}

	if _, has := s.config.RowScanners["messagingFlag"]; has {
		scanner := s.config.RowScanners["messagingFlag"].(func(_ rowScanner, _ *types.MessageFlag) error)
		err = scanner(row, res)
	} else {
		err = row.Scan(
			&res.ID,
			&res.UserID,
			&res.MessageID,
			&res.ChannelID,
			&res.Flag,
			&res.CreatedAt,
		)
	}

	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err != nil {
		return nil, errors.Store("could not scan messagingFlag db row: %s", err).Wrap(err)
	} else {
		return res, nil
	}
}

// QueryMessagingFlags returns squirrel.SelectBuilder with set table and all columns
func (s Store) messagingFlagsSelectBuilder() squirrel.SelectBuilder {
	return s.SelectBuilder(s.messagingFlagTable("mmf"), s.messagingFlagColumns("mmf")...)
}

// messagingFlagTable name of the db table
func (Store) messagingFlagTable(aa ...string) string {
	var alias string
	if len(aa) > 0 {
		alias = " AS " + aa[0]
	}

	return "messaging_message_flag" + alias
}

// MessagingFlagColumns returns all defined table columns
//
// With optional string arg, all columns are returned aliased
func (Store) messagingFlagColumns(aa ...string) []string {
	var alias string
	if len(aa) > 0 {
		alias = aa[0] + "."
	}

	return []string{
		alias + "id",
		alias + "rel_user",
		alias + "rel_message",
		alias + "rel_channel",
		alias + "flag",
		alias + "created_at",
	}
}

// {true true false false false false}

// internalMessagingFlagEncoder encodes fields from types.MessageFlag to store.Payload (map)
//
// Encoding is done by using generic approach or by calling encodeMessagingFlag
// func when rdbms.customEncoder=true
func (s Store) internalMessagingFlagEncoder(res *types.MessageFlag) store.Payload {
	return store.Payload{
		"id":          res.ID,
		"rel_user":    res.UserID,
		"rel_message": res.MessageID,
		"rel_channel": res.ChannelID,
		"flag":        res.Flag,
		"created_at":  res.CreatedAt,
	}
}

// checkMessagingFlagConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant but unfortunately we can not rely
// on the full support (MySQL does not support conditional indexes)
func (s *Store) checkMessagingFlagConstraints(ctx context.Context, res *types.MessageFlag) error {
	// Consider resource valid when all fields in unique constraint check lookups
	// have valid (non-empty) value
	//
	// Only string and uint64 are supported for now
	// feel free to add additional types if needed
	var valid = true

	if !valid {
		return nil
	}

	return nil
}
