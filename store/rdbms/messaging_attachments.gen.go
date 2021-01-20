package rdbms

// This file is an auto-generated file
//
// Template:    pkg/codegen/assets/store_rdbms.gen.go.tpl
// Definitions: store/messaging_attachments.yaml
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

// SearchMessagingAttachments returns all matching rows
//
// This function calls convertMessagingAttachmentFilter with the given
// types.AttachmentFilter and expects to receive a working squirrel.SelectBuilder
func (s Store) SearchMessagingAttachments(ctx context.Context, f types.AttachmentFilter) (types.AttachmentSet, types.AttachmentFilter, error) {
	var (
		err error
		set []*types.Attachment
		q   squirrel.SelectBuilder
	)

	return set, f, func() error {
		q = s.messagingAttachmentsSelectBuilder()

		set, err = s.QueryMessagingAttachments(ctx, q, nil)
		return err
	}()
}

// QueryMessagingAttachments queries the database, converts and checks each row and
// returns collected set
//
// Fn also returns total number of fetched items and last fetched item so that the caller can construct cursor
// for next page of results
func (s Store) QueryMessagingAttachments(
	ctx context.Context,
	q squirrel.Sqlizer,
	check func(*types.Attachment) (bool, error),
) ([]*types.Attachment, error) {
	var (
		set = make([]*types.Attachment, 0, DefaultSliceCapacity)
		res *types.Attachment

		// Query rows with
		rows, err = s.Query(ctx, q)
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		if err = rows.Err(); err == nil {
			res, err = s.internalMessagingAttachmentRowScanner(rows)
		}

		if err != nil {
			return nil, err
		}

		set = append(set, res)
	}

	return set, rows.Err()
}

// LookupMessagingAttachmentByID searches for attachment by its ID
//
// It returns attachment even if deleted
func (s Store) LookupMessagingAttachmentByID(ctx context.Context, id uint64) (*types.Attachment, error) {
	return s.execLookupMessagingAttachment(ctx, squirrel.Eq{
		s.preprocessColumn("att.id", ""): store.PreprocessValue(id, ""),
	})
}

// CreateMessagingAttachment creates one or more rows in messaging_attachment table
func (s Store) CreateMessagingAttachment(ctx context.Context, rr ...*types.Attachment) (err error) {
	for _, res := range rr {
		err = s.checkMessagingAttachmentConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execCreateMessagingAttachments(ctx, s.internalMessagingAttachmentEncoder(res))
		if err != nil {
			return err
		}
	}

	return
}

// UpdateMessagingAttachment updates one or more existing rows in messaging_attachment
func (s Store) UpdateMessagingAttachment(ctx context.Context, rr ...*types.Attachment) error {
	return s.partialMessagingAttachmentUpdate(ctx, nil, rr...)
}

// partialMessagingAttachmentUpdate updates one or more existing rows in messaging_attachment
func (s Store) partialMessagingAttachmentUpdate(ctx context.Context, onlyColumns []string, rr ...*types.Attachment) (err error) {
	for _, res := range rr {
		err = s.checkMessagingAttachmentConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpdateMessagingAttachments(
			ctx,
			squirrel.Eq{
				s.preprocessColumn("att.id", ""): store.PreprocessValue(res.ID, ""),
			},
			s.internalMessagingAttachmentEncoder(res).Skip("id").Only(onlyColumns...))
		if err != nil {
			return err
		}
	}

	return
}

// UpsertMessagingAttachment updates one or more existing rows in messaging_attachment
func (s Store) UpsertMessagingAttachment(ctx context.Context, rr ...*types.Attachment) (err error) {
	for _, res := range rr {
		err = s.checkMessagingAttachmentConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpsertMessagingAttachments(ctx, s.internalMessagingAttachmentEncoder(res))
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteMessagingAttachment Deletes one or more rows from messaging_attachment table
func (s Store) DeleteMessagingAttachment(ctx context.Context, rr ...*types.Attachment) (err error) {
	for _, res := range rr {

		err = s.execDeleteMessagingAttachments(ctx, squirrel.Eq{
			s.preprocessColumn("att.id", ""): store.PreprocessValue(res.ID, ""),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteMessagingAttachmentByID Deletes row from the messaging_attachment table
func (s Store) DeleteMessagingAttachmentByID(ctx context.Context, ID uint64) error {
	return s.execDeleteMessagingAttachments(ctx, squirrel.Eq{
		s.preprocessColumn("att.id", ""): store.PreprocessValue(ID, ""),
	})
}

// TruncateMessagingAttachments Deletes all rows from the messaging_attachment table
func (s Store) TruncateMessagingAttachments(ctx context.Context) error {
	return s.Truncate(ctx, s.messagingAttachmentTable())
}

// execLookupMessagingAttachment prepares MessagingAttachment query and executes it,
// returning types.Attachment (or error)
func (s Store) execLookupMessagingAttachment(ctx context.Context, cnd squirrel.Sqlizer) (res *types.Attachment, err error) {
	var (
		row rowScanner
	)

	row, err = s.QueryRow(ctx, s.messagingAttachmentsSelectBuilder().Where(cnd))
	if err != nil {
		return
	}

	res, err = s.internalMessagingAttachmentRowScanner(row)
	if err != nil {
		return
	}

	return res, nil
}

// execCreateMessagingAttachments updates all matched (by cnd) rows in messaging_attachment with given data
func (s Store) execCreateMessagingAttachments(ctx context.Context, payload store.Payload) error {
	return s.Exec(ctx, s.InsertBuilder(s.messagingAttachmentTable()).SetMap(payload))
}

// execUpdateMessagingAttachments updates all matched (by cnd) rows in messaging_attachment with given data
func (s Store) execUpdateMessagingAttachments(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return s.Exec(ctx, s.UpdateBuilder(s.messagingAttachmentTable("att")).Where(cnd).SetMap(set))
}

// execUpsertMessagingAttachments inserts new or updates matching (by-primary-key) rows in messaging_attachment with given data
func (s Store) execUpsertMessagingAttachments(ctx context.Context, set store.Payload) error {
	upsert, err := s.config.UpsertBuilder(
		s.config,
		s.messagingAttachmentTable(),
		set,
		s.preprocessColumn("id", ""),
	)

	if err != nil {
		return err
	}

	return s.Exec(ctx, upsert)
}

// execDeleteMessagingAttachments Deletes all matched (by cnd) rows in messaging_attachment with given data
func (s Store) execDeleteMessagingAttachments(ctx context.Context, cnd squirrel.Sqlizer) error {
	return s.Exec(ctx, s.DeleteBuilder(s.messagingAttachmentTable("att")).Where(cnd))
}

func (s Store) internalMessagingAttachmentRowScanner(row rowScanner) (res *types.Attachment, err error) {
	res = &types.Attachment{}

	if _, has := s.config.RowScanners["messagingAttachment"]; has {
		scanner := s.config.RowScanners["messagingAttachment"].(func(_ rowScanner, _ *types.Attachment) error)
		err = scanner(row, res)
	} else {
		err = row.Scan(
			&res.ID,
			&res.Url,
			&res.PreviewUrl,
			&res.Name,
			&res.Meta,
			&res.OwnerID,
			&res.CreatedAt,
			&res.UpdatedAt,
			&res.DeletedAt,
		)
	}

	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err != nil {
		return nil, errors.Store("could not scan messagingAttachment db row: %s", err).Wrap(err)
	} else {
		return res, nil
	}
}

// QueryMessagingAttachments returns squirrel.SelectBuilder with set table and all columns
func (s Store) messagingAttachmentsSelectBuilder() squirrel.SelectBuilder {
	return s.SelectBuilder(s.messagingAttachmentTable("att"), s.messagingAttachmentColumns("att")...)
}

// messagingAttachmentTable name of the db table
func (Store) messagingAttachmentTable(aa ...string) string {
	var alias string
	if len(aa) > 0 {
		alias = " AS " + aa[0]
	}

	return "messaging_attachment" + alias
}

// MessagingAttachmentColumns returns all defined table columns
//
// With optional string arg, all columns are returned aliased
func (Store) messagingAttachmentColumns(aa ...string) []string {
	var alias string
	if len(aa) > 0 {
		alias = aa[0] + "."
	}

	return []string{
		alias + "id",
		alias + "url",
		alias + "preview_url",
		alias + "name",
		alias + "meta",
		alias + "rel_owner",
		alias + "created_at",
		alias + "updated_at",
		alias + "deleted_at",
	}
}

// {true true false false false false}

// internalMessagingAttachmentEncoder encodes fields from types.Attachment to store.Payload (map)
//
// Encoding is done by using generic approach or by calling encodeMessagingAttachment
// func when rdbms.customEncoder=true
func (s Store) internalMessagingAttachmentEncoder(res *types.Attachment) store.Payload {
	return store.Payload{
		"id":          res.ID,
		"url":         res.Url,
		"preview_url": res.PreviewUrl,
		"name":        res.Name,
		"meta":        res.Meta,
		"rel_owner":   res.OwnerID,
		"created_at":  res.CreatedAt,
		"updated_at":  res.UpdatedAt,
		"deleted_at":  res.DeletedAt,
	}
}

// checkMessagingAttachmentConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant but unfortunately we can not rely
// on the full support (MySQL does not support conditional indexes)
func (s *Store) checkMessagingAttachmentConstraints(ctx context.Context, res *types.Attachment) error {
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
