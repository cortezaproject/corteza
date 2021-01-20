package rdbms

// This file is an auto-generated file
//
// Template:    pkg/codegen/assets/store_rdbms.gen.go.tpl
// Definitions: store/messaging_message_attachments.yaml
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

// SearchMessagingMessageAttachments returns all matching rows
//
// This function calls convertMessagingMessageAttachmentFilter with the given
// types.MessageAttachmentFilter and expects to receive a working squirrel.SelectBuilder
func (s Store) SearchMessagingMessageAttachments(ctx context.Context, f types.MessageAttachmentFilter) (types.MessageAttachmentSet, types.MessageAttachmentFilter, error) {
	var (
		err error
		set []*types.MessageAttachment
		q   squirrel.SelectBuilder
	)

	return set, f, func() error {
		q = s.messagingMessageAttachmentsSelectBuilder()

		set, err = s.QueryMessagingMessageAttachments(ctx, q, nil)
		return err
	}()
}

// QueryMessagingMessageAttachments queries the database, converts and checks each row and
// returns collected set
//
// Fn also returns total number of fetched items and last fetched item so that the caller can construct cursor
// for next page of results
func (s Store) QueryMessagingMessageAttachments(
	ctx context.Context,
	q squirrel.Sqlizer,
	check func(*types.MessageAttachment) (bool, error),
) ([]*types.MessageAttachment, error) {
	var (
		set = make([]*types.MessageAttachment, 0, DefaultSliceCapacity)
		res *types.MessageAttachment

		// Query rows with
		rows, err = s.Query(ctx, q)
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		if err = rows.Err(); err == nil {
			res, err = s.internalMessagingMessageAttachmentRowScanner(rows)
		}

		if err != nil {
			return nil, err
		}

		set = append(set, res)
	}

	return set, rows.Err()
}

// LookupMessagingMessageAttachmentByMessageID searches for message attachment by message ID
func (s Store) LookupMessagingMessageAttachmentByMessageID(ctx context.Context, message_id uint64) (*types.MessageAttachment, error) {
	return s.execLookupMessagingMessageAttachment(ctx, squirrel.Eq{
		s.preprocessColumn("mma.rel_message", ""): store.PreprocessValue(message_id, ""),
	})
}

// CreateMessagingMessageAttachment creates one or more rows in messaging_message_attachment table
func (s Store) CreateMessagingMessageAttachment(ctx context.Context, rr ...*types.MessageAttachment) (err error) {
	for _, res := range rr {
		err = s.checkMessagingMessageAttachmentConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execCreateMessagingMessageAttachments(ctx, s.internalMessagingMessageAttachmentEncoder(res))
		if err != nil {
			return err
		}
	}

	return
}

// UpdateMessagingMessageAttachment updates one or more existing rows in messaging_message_attachment
func (s Store) UpdateMessagingMessageAttachment(ctx context.Context, rr ...*types.MessageAttachment) error {
	return s.partialMessagingMessageAttachmentUpdate(ctx, nil, rr...)
}

// partialMessagingMessageAttachmentUpdate updates one or more existing rows in messaging_message_attachment
func (s Store) partialMessagingMessageAttachmentUpdate(ctx context.Context, onlyColumns []string, rr ...*types.MessageAttachment) (err error) {
	for _, res := range rr {
		err = s.checkMessagingMessageAttachmentConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpdateMessagingMessageAttachments(
			ctx,
			squirrel.Eq{
				s.preprocessColumn("mma.rel_message", ""): store.PreprocessValue(res.MessageID, ""),
			},
			s.internalMessagingMessageAttachmentEncoder(res).Skip("rel_message").Only(onlyColumns...))
		if err != nil {
			return err
		}
	}

	return
}

// UpsertMessagingMessageAttachment updates one or more existing rows in messaging_message_attachment
func (s Store) UpsertMessagingMessageAttachment(ctx context.Context, rr ...*types.MessageAttachment) (err error) {
	for _, res := range rr {
		err = s.checkMessagingMessageAttachmentConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpsertMessagingMessageAttachments(ctx, s.internalMessagingMessageAttachmentEncoder(res))
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteMessagingMessageAttachment Deletes one or more rows from messaging_message_attachment table
func (s Store) DeleteMessagingMessageAttachment(ctx context.Context, rr ...*types.MessageAttachment) (err error) {
	for _, res := range rr {

		err = s.execDeleteMessagingMessageAttachments(ctx, squirrel.Eq{
			s.preprocessColumn("mma.rel_message", ""): store.PreprocessValue(res.MessageID, ""),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteMessagingMessageAttachmentByMessageID Deletes row from the messaging_message_attachment table
func (s Store) DeleteMessagingMessageAttachmentByMessageID(ctx context.Context, messageID uint64) error {
	return s.execDeleteMessagingMessageAttachments(ctx, squirrel.Eq{
		s.preprocessColumn("mma.rel_message", ""): store.PreprocessValue(messageID, ""),
	})
}

// TruncateMessagingMessageAttachments Deletes all rows from the messaging_message_attachment table
func (s Store) TruncateMessagingMessageAttachments(ctx context.Context) error {
	return s.Truncate(ctx, s.messagingMessageAttachmentTable())
}

// execLookupMessagingMessageAttachment prepares MessagingMessageAttachment query and executes it,
// returning types.MessageAttachment (or error)
func (s Store) execLookupMessagingMessageAttachment(ctx context.Context, cnd squirrel.Sqlizer) (res *types.MessageAttachment, err error) {
	var (
		row rowScanner
	)

	row, err = s.QueryRow(ctx, s.messagingMessageAttachmentsSelectBuilder().Where(cnd))
	if err != nil {
		return
	}

	res, err = s.internalMessagingMessageAttachmentRowScanner(row)
	if err != nil {
		return
	}

	return res, nil
}

// execCreateMessagingMessageAttachments updates all matched (by cnd) rows in messaging_message_attachment with given data
func (s Store) execCreateMessagingMessageAttachments(ctx context.Context, payload store.Payload) error {
	return s.Exec(ctx, s.InsertBuilder(s.messagingMessageAttachmentTable()).SetMap(payload))
}

// execUpdateMessagingMessageAttachments updates all matched (by cnd) rows in messaging_message_attachment with given data
func (s Store) execUpdateMessagingMessageAttachments(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return s.Exec(ctx, s.UpdateBuilder(s.messagingMessageAttachmentTable("mma")).Where(cnd).SetMap(set))
}

// execUpsertMessagingMessageAttachments inserts new or updates matching (by-primary-key) rows in messaging_message_attachment with given data
func (s Store) execUpsertMessagingMessageAttachments(ctx context.Context, set store.Payload) error {
	upsert, err := s.config.UpsertBuilder(
		s.config,
		s.messagingMessageAttachmentTable(),
		set,
		s.preprocessColumn("rel_message", ""),
	)

	if err != nil {
		return err
	}

	return s.Exec(ctx, upsert)
}

// execDeleteMessagingMessageAttachments Deletes all matched (by cnd) rows in messaging_message_attachment with given data
func (s Store) execDeleteMessagingMessageAttachments(ctx context.Context, cnd squirrel.Sqlizer) error {
	return s.Exec(ctx, s.DeleteBuilder(s.messagingMessageAttachmentTable("mma")).Where(cnd))
}

func (s Store) internalMessagingMessageAttachmentRowScanner(row rowScanner) (res *types.MessageAttachment, err error) {
	res = &types.MessageAttachment{}

	if _, has := s.config.RowScanners["messagingMessageAttachment"]; has {
		scanner := s.config.RowScanners["messagingMessageAttachment"].(func(_ rowScanner, _ *types.MessageAttachment) error)
		err = scanner(row, res)
	} else {
		err = row.Scan(
			&res.MessageID,
			&res.AttachmentID,
		)
	}

	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err != nil {
		return nil, errors.Store("could not scan messagingMessageAttachment db row: %s", err).Wrap(err)
	} else {
		return res, nil
	}
}

// QueryMessagingMessageAttachments returns squirrel.SelectBuilder with set table and all columns
func (s Store) messagingMessageAttachmentsSelectBuilder() squirrel.SelectBuilder {
	return s.SelectBuilder(s.messagingMessageAttachmentTable("mma"), s.messagingMessageAttachmentColumns("mma")...)
}

// messagingMessageAttachmentTable name of the db table
func (Store) messagingMessageAttachmentTable(aa ...string) string {
	var alias string
	if len(aa) > 0 {
		alias = " AS " + aa[0]
	}

	return "messaging_message_attachment" + alias
}

// MessagingMessageAttachmentColumns returns all defined table columns
//
// With optional string arg, all columns are returned aliased
func (Store) messagingMessageAttachmentColumns(aa ...string) []string {
	var alias string
	if len(aa) > 0 {
		alias = aa[0] + "."
	}

	return []string{
		alias + "rel_message",
		alias + "rel_attachment",
	}
}

// {true true false false false false}

// internalMessagingMessageAttachmentEncoder encodes fields from types.MessageAttachment to store.Payload (map)
//
// Encoding is done by using generic approach or by calling encodeMessagingMessageAttachment
// func when rdbms.customEncoder=true
func (s Store) internalMessagingMessageAttachmentEncoder(res *types.MessageAttachment) store.Payload {
	return store.Payload{
		"rel_message":    res.MessageID,
		"rel_attachment": res.AttachmentID,
	}
}

// checkMessagingMessageAttachmentConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant but unfortunately we can not rely
// on the full support (MySQL does not support conditional indexes)
func (s *Store) checkMessagingMessageAttachmentConstraints(ctx context.Context, res *types.MessageAttachment) error {
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
