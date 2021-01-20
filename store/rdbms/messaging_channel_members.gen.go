package rdbms

// This file is an auto-generated file
//
// Template:    pkg/codegen/assets/store_rdbms.gen.go.tpl
// Definitions: store/messaging_channel_members.yaml
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

// SearchMessagingChannelMembers returns all matching rows
//
// This function calls convertMessagingChannelMemberFilter with the given
// types.ChannelMemberFilter and expects to receive a working squirrel.SelectBuilder
func (s Store) SearchMessagingChannelMembers(ctx context.Context, f types.ChannelMemberFilter) (types.ChannelMemberSet, types.ChannelMemberFilter, error) {
	var (
		err error
		set []*types.ChannelMember
		q   squirrel.SelectBuilder
	)

	return set, f, func() error {
		q, err = s.convertMessagingChannelMemberFilter(f)
		if err != nil {
			return err
		}

		set, err = s.QueryMessagingChannelMembers(ctx, q, nil)
		return err
	}()
}

// QueryMessagingChannelMembers queries the database, converts and checks each row and
// returns collected set
//
// Fn also returns total number of fetched items and last fetched item so that the caller can construct cursor
// for next page of results
func (s Store) QueryMessagingChannelMembers(
	ctx context.Context,
	q squirrel.Sqlizer,
	check func(*types.ChannelMember) (bool, error),
) ([]*types.ChannelMember, error) {
	var (
		set = make([]*types.ChannelMember, 0, DefaultSliceCapacity)
		res *types.ChannelMember

		// Query rows with
		rows, err = s.Query(ctx, q)
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		if err = rows.Err(); err == nil {
			res, err = s.internalMessagingChannelMemberRowScanner(rows)
		}

		if err != nil {
			return nil, err
		}

		set = append(set, res)
	}

	return set, rows.Err()
}

// CreateMessagingChannelMember creates one or more rows in messaging_channel_member table
func (s Store) CreateMessagingChannelMember(ctx context.Context, rr ...*types.ChannelMember) (err error) {
	for _, res := range rr {
		err = s.checkMessagingChannelMemberConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execCreateMessagingChannelMembers(ctx, s.internalMessagingChannelMemberEncoder(res))
		if err != nil {
			return err
		}
	}

	return
}

// UpdateMessagingChannelMember updates one or more existing rows in messaging_channel_member
func (s Store) UpdateMessagingChannelMember(ctx context.Context, rr ...*types.ChannelMember) error {
	return s.partialMessagingChannelMemberUpdate(ctx, nil, rr...)
}

// partialMessagingChannelMemberUpdate updates one or more existing rows in messaging_channel_member
func (s Store) partialMessagingChannelMemberUpdate(ctx context.Context, onlyColumns []string, rr ...*types.ChannelMember) (err error) {
	for _, res := range rr {
		err = s.checkMessagingChannelMemberConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpdateMessagingChannelMembers(
			ctx,
			squirrel.Eq{
				s.preprocessColumn("mcm.rel_channel", ""): store.PreprocessValue(res.ChannelID, ""), s.preprocessColumn("mcm.rel_user", ""): store.PreprocessValue(res.UserID, ""),
			},
			s.internalMessagingChannelMemberEncoder(res).Skip("rel_channel", "rel_user").Only(onlyColumns...))
		if err != nil {
			return err
		}
	}

	return
}

// UpsertMessagingChannelMember updates one or more existing rows in messaging_channel_member
func (s Store) UpsertMessagingChannelMember(ctx context.Context, rr ...*types.ChannelMember) (err error) {
	for _, res := range rr {
		err = s.checkMessagingChannelMemberConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpsertMessagingChannelMembers(ctx, s.internalMessagingChannelMemberEncoder(res))
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteMessagingChannelMember Deletes one or more rows from messaging_channel_member table
func (s Store) DeleteMessagingChannelMember(ctx context.Context, rr ...*types.ChannelMember) (err error) {
	for _, res := range rr {

		err = s.execDeleteMessagingChannelMembers(ctx, squirrel.Eq{
			s.preprocessColumn("mcm.rel_channel", ""): store.PreprocessValue(res.ChannelID, ""), s.preprocessColumn("mcm.rel_user", ""): store.PreprocessValue(res.UserID, ""),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteMessagingChannelMemberByChannelIDUserID Deletes row from the messaging_channel_member table
func (s Store) DeleteMessagingChannelMemberByChannelIDUserID(ctx context.Context, channelID uint64, userID uint64) error {
	return s.execDeleteMessagingChannelMembers(ctx, squirrel.Eq{
		s.preprocessColumn("mcm.rel_channel", ""): store.PreprocessValue(channelID, ""),
		s.preprocessColumn("mcm.rel_user", ""):    store.PreprocessValue(userID, ""),
	})
}

// TruncateMessagingChannelMembers Deletes all rows from the messaging_channel_member table
func (s Store) TruncateMessagingChannelMembers(ctx context.Context) error {
	return s.Truncate(ctx, s.messagingChannelMemberTable())
}

// execLookupMessagingChannelMember prepares MessagingChannelMember query and executes it,
// returning types.ChannelMember (or error)
func (s Store) execLookupMessagingChannelMember(ctx context.Context, cnd squirrel.Sqlizer) (res *types.ChannelMember, err error) {
	var (
		row rowScanner
	)

	row, err = s.QueryRow(ctx, s.messagingChannelMembersSelectBuilder().Where(cnd))
	if err != nil {
		return
	}

	res, err = s.internalMessagingChannelMemberRowScanner(row)
	if err != nil {
		return
	}

	return res, nil
}

// execCreateMessagingChannelMembers updates all matched (by cnd) rows in messaging_channel_member with given data
func (s Store) execCreateMessagingChannelMembers(ctx context.Context, payload store.Payload) error {
	return s.Exec(ctx, s.InsertBuilder(s.messagingChannelMemberTable()).SetMap(payload))
}

// execUpdateMessagingChannelMembers updates all matched (by cnd) rows in messaging_channel_member with given data
func (s Store) execUpdateMessagingChannelMembers(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return s.Exec(ctx, s.UpdateBuilder(s.messagingChannelMemberTable("mcm")).Where(cnd).SetMap(set))
}

// execUpsertMessagingChannelMembers inserts new or updates matching (by-primary-key) rows in messaging_channel_member with given data
func (s Store) execUpsertMessagingChannelMembers(ctx context.Context, set store.Payload) error {
	upsert, err := s.config.UpsertBuilder(
		s.config,
		s.messagingChannelMemberTable(),
		set,
		s.preprocessColumn("rel_channel", ""),
		s.preprocessColumn("rel_user", ""),
	)

	if err != nil {
		return err
	}

	return s.Exec(ctx, upsert)
}

// execDeleteMessagingChannelMembers Deletes all matched (by cnd) rows in messaging_channel_member with given data
func (s Store) execDeleteMessagingChannelMembers(ctx context.Context, cnd squirrel.Sqlizer) error {
	return s.Exec(ctx, s.DeleteBuilder(s.messagingChannelMemberTable("mcm")).Where(cnd))
}

func (s Store) internalMessagingChannelMemberRowScanner(row rowScanner) (res *types.ChannelMember, err error) {
	res = &types.ChannelMember{}

	if _, has := s.config.RowScanners["messagingChannelMember"]; has {
		scanner := s.config.RowScanners["messagingChannelMember"].(func(_ rowScanner, _ *types.ChannelMember) error)
		err = scanner(row, res)
	} else {
		err = row.Scan(
			&res.ChannelID,
			&res.UserID,
			&res.Type,
			&res.Flag,
			&res.CreatedAt,
			&res.UpdatedAt,
		)
	}

	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err != nil {
		return nil, errors.Store("could not scan messagingChannelMember db row: %s", err).Wrap(err)
	} else {
		return res, nil
	}
}

// QueryMessagingChannelMembers returns squirrel.SelectBuilder with set table and all columns
func (s Store) messagingChannelMembersSelectBuilder() squirrel.SelectBuilder {
	return s.SelectBuilder(s.messagingChannelMemberTable("mcm"), s.messagingChannelMemberColumns("mcm")...)
}

// messagingChannelMemberTable name of the db table
func (Store) messagingChannelMemberTable(aa ...string) string {
	var alias string
	if len(aa) > 0 {
		alias = " AS " + aa[0]
	}

	return "messaging_channel_member" + alias
}

// MessagingChannelMemberColumns returns all defined table columns
//
// With optional string arg, all columns are returned aliased
func (Store) messagingChannelMemberColumns(aa ...string) []string {
	var alias string
	if len(aa) > 0 {
		alias = aa[0] + "."
	}

	return []string{
		alias + "rel_channel",
		alias + "rel_user",
		alias + "type",
		alias + "flag",
		alias + "created_at",
		alias + "updated_at",
	}
}

// {true true false false false false}

// internalMessagingChannelMemberEncoder encodes fields from types.ChannelMember to store.Payload (map)
//
// Encoding is done by using generic approach or by calling encodeMessagingChannelMember
// func when rdbms.customEncoder=true
func (s Store) internalMessagingChannelMemberEncoder(res *types.ChannelMember) store.Payload {
	return store.Payload{
		"rel_channel": res.ChannelID,
		"rel_user":    res.UserID,
		"type":        res.Type,
		"flag":        res.Flag,
		"created_at":  res.CreatedAt,
		"updated_at":  res.UpdatedAt,
	}
}

// checkMessagingChannelMemberConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant but unfortunately we can not rely
// on the full support (MySQL does not support conditional indexes)
func (s *Store) checkMessagingChannelMemberConstraints(ctx context.Context, res *types.ChannelMember) error {
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
