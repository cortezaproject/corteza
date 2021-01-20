package rdbms

// This file is an auto-generated file
//
// Template:    pkg/codegen/assets/store_rdbms.gen.go.tpl
// Definitions: store/messaging_unread.yaml
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

// QueryMessagingUnreads queries the database, converts and checks each row and
// returns collected set
//
// Fn also returns total number of fetched items and last fetched item so that the caller can construct cursor
// for next page of results
func (s Store) QueryMessagingUnreads(
	ctx context.Context,
	q squirrel.Sqlizer,
	check func(*types.Unread) (bool, error),
) ([]*types.Unread, error) {
	var (
		set = make([]*types.Unread, 0, DefaultSliceCapacity)
		res *types.Unread

		// Query rows with
		rows, err = s.Query(ctx, q)
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		if err = rows.Err(); err == nil {
			res, err = s.internalMessagingUnreadRowScanner(rows)
		}

		if err != nil {
			return nil, err
		}

		set = append(set, res)
	}

	return set, rows.Err()
}

// CreateMessagingUnread creates one or more rows in messaging_unread table
func (s Store) CreateMessagingUnread(ctx context.Context, rr ...*types.Unread) (err error) {
	for _, res := range rr {
		err = s.checkMessagingUnreadConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execCreateMessagingUnreads(ctx, s.internalMessagingUnreadEncoder(res))
		if err != nil {
			return err
		}
	}

	return
}

// UpdateMessagingUnread updates one or more existing rows in messaging_unread
func (s Store) UpdateMessagingUnread(ctx context.Context, rr ...*types.Unread) error {
	return s.partialMessagingUnreadUpdate(ctx, nil, rr...)
}

// partialMessagingUnreadUpdate updates one or more existing rows in messaging_unread
func (s Store) partialMessagingUnreadUpdate(ctx context.Context, onlyColumns []string, rr ...*types.Unread) (err error) {
	for _, res := range rr {
		err = s.checkMessagingUnreadConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpdateMessagingUnreads(
			ctx,
			squirrel.Eq{
				s.preprocessColumn("mur.rel_channel", ""): store.PreprocessValue(res.ChannelID, ""), s.preprocessColumn("mur.rel_reply_to", ""): store.PreprocessValue(res.ReplyTo, ""), s.preprocessColumn("mur.rel_user", ""): store.PreprocessValue(res.UserID, ""),
			},
			s.internalMessagingUnreadEncoder(res).Skip("rel_channel", "rel_reply_to", "rel_user").Only(onlyColumns...))
		if err != nil {
			return err
		}
	}

	return
}

// UpsertMessagingUnread updates one or more existing rows in messaging_unread
func (s Store) UpsertMessagingUnread(ctx context.Context, rr ...*types.Unread) (err error) {
	for _, res := range rr {
		err = s.checkMessagingUnreadConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpsertMessagingUnreads(ctx, s.internalMessagingUnreadEncoder(res))
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteMessagingUnread Deletes one or more rows from messaging_unread table
func (s Store) DeleteMessagingUnread(ctx context.Context, rr ...*types.Unread) (err error) {
	for _, res := range rr {

		err = s.execDeleteMessagingUnreads(ctx, squirrel.Eq{
			s.preprocessColumn("mur.rel_channel", ""): store.PreprocessValue(res.ChannelID, ""), s.preprocessColumn("mur.rel_reply_to", ""): store.PreprocessValue(res.ReplyTo, ""), s.preprocessColumn("mur.rel_user", ""): store.PreprocessValue(res.UserID, ""),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteMessagingUnreadByChannelIDReplyToUserID Deletes row from the messaging_unread table
func (s Store) DeleteMessagingUnreadByChannelIDReplyToUserID(ctx context.Context, channelID uint64, replyTo uint64, userID uint64) error {
	return s.execDeleteMessagingUnreads(ctx, squirrel.Eq{
		s.preprocessColumn("mur.rel_channel", ""):  store.PreprocessValue(channelID, ""),
		s.preprocessColumn("mur.rel_reply_to", ""): store.PreprocessValue(replyTo, ""),
		s.preprocessColumn("mur.rel_user", ""):     store.PreprocessValue(userID, ""),
	})
}

// TruncateMessagingUnreads Deletes all rows from the messaging_unread table
func (s Store) TruncateMessagingUnreads(ctx context.Context) error {
	return s.Truncate(ctx, s.messagingUnreadTable())
}

// execLookupMessagingUnread prepares MessagingUnread query and executes it,
// returning types.Unread (or error)
func (s Store) execLookupMessagingUnread(ctx context.Context, cnd squirrel.Sqlizer) (res *types.Unread, err error) {
	var (
		row rowScanner
	)

	row, err = s.QueryRow(ctx, s.messagingUnreadsSelectBuilder().Where(cnd))
	if err != nil {
		return
	}

	res, err = s.internalMessagingUnreadRowScanner(row)
	if err != nil {
		return
	}

	return res, nil
}

// execCreateMessagingUnreads updates all matched (by cnd) rows in messaging_unread with given data
func (s Store) execCreateMessagingUnreads(ctx context.Context, payload store.Payload) error {
	return s.Exec(ctx, s.InsertBuilder(s.messagingUnreadTable()).SetMap(payload))
}

// execUpdateMessagingUnreads updates all matched (by cnd) rows in messaging_unread with given data
func (s Store) execUpdateMessagingUnreads(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return s.Exec(ctx, s.UpdateBuilder(s.messagingUnreadTable("mur")).Where(cnd).SetMap(set))
}

// execUpsertMessagingUnreads inserts new or updates matching (by-primary-key) rows in messaging_unread with given data
func (s Store) execUpsertMessagingUnreads(ctx context.Context, set store.Payload) error {
	upsert, err := s.config.UpsertBuilder(
		s.config,
		s.messagingUnreadTable(),
		set,
		s.preprocessColumn("rel_channel", ""),
		s.preprocessColumn("rel_reply_to", ""),
		s.preprocessColumn("rel_user", ""),
	)

	if err != nil {
		return err
	}

	return s.Exec(ctx, upsert)
}

// execDeleteMessagingUnreads Deletes all matched (by cnd) rows in messaging_unread with given data
func (s Store) execDeleteMessagingUnreads(ctx context.Context, cnd squirrel.Sqlizer) error {
	return s.Exec(ctx, s.DeleteBuilder(s.messagingUnreadTable("mur")).Where(cnd))
}

func (s Store) internalMessagingUnreadRowScanner(row rowScanner) (res *types.Unread, err error) {
	res = &types.Unread{}

	if _, has := s.config.RowScanners["messagingUnread"]; has {
		scanner := s.config.RowScanners["messagingUnread"].(func(_ rowScanner, _ *types.Unread) error)
		err = scanner(row, res)
	} else {
		err = row.Scan(
			&res.ChannelID,
			&res.ReplyTo,
			&res.UserID,
			&res.LastMessageID,
			&res.Count,
		)
	}

	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err != nil {
		return nil, errors.Store("could not scan messagingUnread db row: %s", err).Wrap(err)
	} else {
		return res, nil
	}
}

// QueryMessagingUnreads returns squirrel.SelectBuilder with set table and all columns
func (s Store) messagingUnreadsSelectBuilder() squirrel.SelectBuilder {
	return s.SelectBuilder(s.messagingUnreadTable("mur"), s.messagingUnreadColumns("mur")...)
}

// messagingUnreadTable name of the db table
func (Store) messagingUnreadTable(aa ...string) string {
	var alias string
	if len(aa) > 0 {
		alias = " AS " + aa[0]
	}

	return "messaging_unread" + alias
}

// MessagingUnreadColumns returns all defined table columns
//
// With optional string arg, all columns are returned aliased
func (Store) messagingUnreadColumns(aa ...string) []string {
	var alias string
	if len(aa) > 0 {
		alias = aa[0] + "."
	}

	return []string{
		alias + "rel_channel",
		alias + "rel_reply_to",
		alias + "rel_user",
		alias + "rel_last_message",
		alias + "count",
	}
}

// {false true false false false false}

// internalMessagingUnreadEncoder encodes fields from types.Unread to store.Payload (map)
//
// Encoding is done by using generic approach or by calling encodeMessagingUnread
// func when rdbms.customEncoder=true
func (s Store) internalMessagingUnreadEncoder(res *types.Unread) store.Payload {
	return store.Payload{
		"rel_channel":      res.ChannelID,
		"rel_reply_to":     res.ReplyTo,
		"rel_user":         res.UserID,
		"rel_last_message": res.LastMessageID,
		"count":            res.Count,
	}
}

// checkMessagingUnreadConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant but unfortunately we can not rely
// on the full support (MySQL does not support conditional indexes)
func (s *Store) checkMessagingUnreadConstraints(ctx context.Context, res *types.Unread) error {
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
