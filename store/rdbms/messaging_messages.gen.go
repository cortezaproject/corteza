package rdbms

// This file is an auto-generated file
//
// Template:    pkg/codegen/assets/store_rdbms.gen.go.tpl
// Definitions: store/messaging_messages.yaml
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

// SearchMessagingMessages returns all matching rows
//
// This function calls convertMessagingMessageFilter with the given
// types.MessageFilter and expects to receive a working squirrel.SelectBuilder
func (s Store) SearchMessagingMessages(ctx context.Context, f types.MessageFilter) (types.MessageSet, types.MessageFilter, error) {
	var (
		err error
		set []*types.Message
		q   squirrel.SelectBuilder
	)

	return set, f, func() error {
		q, err = s.convertMessagingMessageFilter(f)
		if err != nil {
			return err
		}

		set, err = s.QueryMessagingMessages(ctx, q, nil)
		return err
	}()
}

// QueryMessagingMessages queries the database, converts and checks each row and
// returns collected set
//
// Fn also returns total number of fetched items and last fetched item so that the caller can construct cursor
// for next page of results
func (s Store) QueryMessagingMessages(
	ctx context.Context,
	q squirrel.Sqlizer,
	check func(*types.Message) (bool, error),
) ([]*types.Message, error) {
	var (
		set = make([]*types.Message, 0, DefaultSliceCapacity)
		res *types.Message

		// Query rows with
		rows, err = s.Query(ctx, q)
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		if err = rows.Err(); err == nil {
			res, err = s.internalMessagingMessageRowScanner(rows)
		}

		if err != nil {
			return nil, err
		}

		set = append(set, res)
	}

	return set, rows.Err()
}

// LookupMessagingMessageByID searches for message by its ID
//
// It returns message even if deleted
func (s Store) LookupMessagingMessageByID(ctx context.Context, id uint64) (*types.Message, error) {
	return s.execLookupMessagingMessage(ctx, squirrel.Eq{
		s.preprocessColumn("msg.id", ""): store.PreprocessValue(id, ""),
	})
}

// CreateMessagingMessage creates one or more rows in messaging_message table
func (s Store) CreateMessagingMessage(ctx context.Context, rr ...*types.Message) (err error) {
	for _, res := range rr {
		err = s.checkMessagingMessageConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execCreateMessagingMessages(ctx, s.internalMessagingMessageEncoder(res))
		if err != nil {
			return err
		}
	}

	return
}

// UpdateMessagingMessage updates one or more existing rows in messaging_message
func (s Store) UpdateMessagingMessage(ctx context.Context, rr ...*types.Message) error {
	return s.partialMessagingMessageUpdate(ctx, nil, rr...)
}

// partialMessagingMessageUpdate updates one or more existing rows in messaging_message
func (s Store) partialMessagingMessageUpdate(ctx context.Context, onlyColumns []string, rr ...*types.Message) (err error) {
	for _, res := range rr {
		err = s.checkMessagingMessageConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpdateMessagingMessages(
			ctx,
			squirrel.Eq{
				s.preprocessColumn("msg.id", ""): store.PreprocessValue(res.ID, ""),
			},
			s.internalMessagingMessageEncoder(res).Skip("id").Only(onlyColumns...))
		if err != nil {
			return err
		}
	}

	return
}

// UpsertMessagingMessage updates one or more existing rows in messaging_message
func (s Store) UpsertMessagingMessage(ctx context.Context, rr ...*types.Message) (err error) {
	for _, res := range rr {
		err = s.checkMessagingMessageConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpsertMessagingMessages(ctx, s.internalMessagingMessageEncoder(res))
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteMessagingMessage Deletes one or more rows from messaging_message table
func (s Store) DeleteMessagingMessage(ctx context.Context, rr ...*types.Message) (err error) {
	for _, res := range rr {

		err = s.execDeleteMessagingMessages(ctx, squirrel.Eq{
			s.preprocessColumn("msg.id", ""): store.PreprocessValue(res.ID, ""),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteMessagingMessageByID Deletes row from the messaging_message table
func (s Store) DeleteMessagingMessageByID(ctx context.Context, ID uint64) error {
	return s.execDeleteMessagingMessages(ctx, squirrel.Eq{
		s.preprocessColumn("msg.id", ""): store.PreprocessValue(ID, ""),
	})
}

// TruncateMessagingMessages Deletes all rows from the messaging_message table
func (s Store) TruncateMessagingMessages(ctx context.Context) error {
	return s.Truncate(ctx, s.messagingMessageTable())
}

// execLookupMessagingMessage prepares MessagingMessage query and executes it,
// returning types.Message (or error)
func (s Store) execLookupMessagingMessage(ctx context.Context, cnd squirrel.Sqlizer) (res *types.Message, err error) {
	var (
		row rowScanner
	)

	row, err = s.QueryRow(ctx, s.messagingMessagesSelectBuilder().Where(cnd))
	if err != nil {
		return
	}

	res, err = s.internalMessagingMessageRowScanner(row)
	if err != nil {
		return
	}

	return res, nil
}

// execCreateMessagingMessages updates all matched (by cnd) rows in messaging_message with given data
func (s Store) execCreateMessagingMessages(ctx context.Context, payload store.Payload) error {
	return s.Exec(ctx, s.InsertBuilder(s.messagingMessageTable()).SetMap(payload))
}

// execUpdateMessagingMessages updates all matched (by cnd) rows in messaging_message with given data
func (s Store) execUpdateMessagingMessages(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return s.Exec(ctx, s.UpdateBuilder(s.messagingMessageTable("msg")).Where(cnd).SetMap(set))
}

// execUpsertMessagingMessages inserts new or updates matching (by-primary-key) rows in messaging_message with given data
func (s Store) execUpsertMessagingMessages(ctx context.Context, set store.Payload) error {
	upsert, err := s.config.UpsertBuilder(
		s.config,
		s.messagingMessageTable(),
		set,
		s.preprocessColumn("id", ""),
	)

	if err != nil {
		return err
	}

	return s.Exec(ctx, upsert)
}

// execDeleteMessagingMessages Deletes all matched (by cnd) rows in messaging_message with given data
func (s Store) execDeleteMessagingMessages(ctx context.Context, cnd squirrel.Sqlizer) error {
	return s.Exec(ctx, s.DeleteBuilder(s.messagingMessageTable("msg")).Where(cnd))
}

func (s Store) internalMessagingMessageRowScanner(row rowScanner) (res *types.Message, err error) {
	res = &types.Message{}

	if _, has := s.config.RowScanners["messagingMessage"]; has {
		scanner := s.config.RowScanners["messagingMessage"].(func(_ rowScanner, _ *types.Message) error)
		err = scanner(row, res)
	} else {
		err = row.Scan(
			&res.ID,
			&res.Type,
			&res.Message,
			&res.Meta,
			&res.UserID,
			&res.ChannelID,
			&res.ReplyTo,
			&res.Replies,
			&res.CreatedAt,
			&res.UpdatedAt,
			&res.DeletedAt,
		)
	}

	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err != nil {
		return nil, errors.Store("could not scan messagingMessage db row: %s", err).Wrap(err)
	} else {
		return res, nil
	}
}

// QueryMessagingMessages returns squirrel.SelectBuilder with set table and all columns
func (s Store) messagingMessagesSelectBuilder() squirrel.SelectBuilder {
	return s.SelectBuilder(s.messagingMessageTable("msg"), s.messagingMessageColumns("msg")...)
}

// messagingMessageTable name of the db table
func (Store) messagingMessageTable(aa ...string) string {
	var alias string
	if len(aa) > 0 {
		alias = " AS " + aa[0]
	}

	return "messaging_message" + alias
}

// MessagingMessageColumns returns all defined table columns
//
// With optional string arg, all columns are returned aliased
func (Store) messagingMessageColumns(aa ...string) []string {
	var alias string
	if len(aa) > 0 {
		alias = aa[0] + "."
	}

	return []string{
		alias + "id",
		alias + "type",
		alias + "message",
		alias + "meta",
		alias + "rel_user",
		alias + "rel_channel",
		alias + "reply_to",
		alias + "replies",
		alias + "created_at",
		alias + "updated_at",
		alias + "deleted_at",
	}
}

// {true true false false false false}

// internalMessagingMessageEncoder encodes fields from types.Message to store.Payload (map)
//
// Encoding is done by using generic approach or by calling encodeMessagingMessage
// func when rdbms.customEncoder=true
func (s Store) internalMessagingMessageEncoder(res *types.Message) store.Payload {
	return s.encodeMessagingMessage(res)
}

// checkMessagingMessageConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant but unfortunately we can not rely
// on the full support (MySQL does not support conditional indexes)
func (s *Store) checkMessagingMessageConstraints(ctx context.Context, res *types.Message) error {
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
