package rdbms

// This file is an auto-generated file
//
// Template:    pkg/codegen/assets/store_rdbms.gen.go.tpl
// Definitions: store/attachments.yaml
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

// SearchAttachments returns all matching rows
//
// This function calls convertAttachmentFilter with the given
// types.AttachmentFilter and expects to receive a working squirrel.SelectBuilder
func (s Store) SearchAttachments(ctx context.Context, f types.AttachmentFilter) (types.AttachmentSet, types.AttachmentFilter, error) {
	var (
		err error
		set []*types.Attachment
		q   squirrel.SelectBuilder
	)

	return set, f, func() error {
		q, err = s.convertAttachmentFilter(f)
		if err != nil {
			return err
		}

		set, err = s.QueryAttachments(ctx, q, f.Check)
		return err
	}()
}

// QueryAttachments queries the database, converts and checks each row and
// returns collected set
//
// Fn also returns total number of fetched items and last fetched item so that the caller can construct cursor
// for next page of results
func (s Store) QueryAttachments(
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
			res, err = s.internalAttachmentRowScanner(rows)
		}

		if err != nil {
			return nil, err
		}

		// check fn set, call it and see if it passed the test
		// if not, skip the item
		if check != nil {
			if chk, err := check(res); err != nil {
				return nil, err
			} else if !chk {
				continue
			}
		}

		set = append(set, res)
	}

	return set, rows.Err()
}

// LookupAttachmentByID searches for attachment by its ID
//
// It returns attachment even if deleted
func (s Store) LookupAttachmentByID(ctx context.Context, id uint64) (*types.Attachment, error) {
	return s.execLookupAttachment(ctx, squirrel.Eq{
		s.preprocessColumn("att.id", ""): store.PreprocessValue(id, ""),
	})
}

// CreateAttachment creates one or more rows in attachments table
func (s Store) CreateAttachment(ctx context.Context, rr ...*types.Attachment) (err error) {
	for _, res := range rr {
		err = s.checkAttachmentConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execCreateAttachments(ctx, s.internalAttachmentEncoder(res))
		if err != nil {
			return err
		}
	}

	return
}

// UpdateAttachment updates one or more existing rows in attachments
func (s Store) UpdateAttachment(ctx context.Context, rr ...*types.Attachment) error {
	return s.partialAttachmentUpdate(ctx, nil, rr...)
}

// partialAttachmentUpdate updates one or more existing rows in attachments
func (s Store) partialAttachmentUpdate(ctx context.Context, onlyColumns []string, rr ...*types.Attachment) (err error) {
	for _, res := range rr {
		err = s.checkAttachmentConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpdateAttachments(
			ctx,
			squirrel.Eq{
				s.preprocessColumn("att.id", ""): store.PreprocessValue(res.ID, ""),
			},
			s.internalAttachmentEncoder(res).Skip("id").Only(onlyColumns...))
		if err != nil {
			return err
		}
	}

	return
}

// UpsertAttachment updates one or more existing rows in attachments
func (s Store) UpsertAttachment(ctx context.Context, rr ...*types.Attachment) (err error) {
	for _, res := range rr {
		err = s.checkAttachmentConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpsertAttachments(ctx, s.internalAttachmentEncoder(res))
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteAttachment Deletes one or more rows from attachments table
func (s Store) DeleteAttachment(ctx context.Context, rr ...*types.Attachment) (err error) {
	for _, res := range rr {

		err = s.execDeleteAttachments(ctx, squirrel.Eq{
			s.preprocessColumn("att.id", ""): store.PreprocessValue(res.ID, ""),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteAttachmentByID Deletes row from the attachments table
func (s Store) DeleteAttachmentByID(ctx context.Context, ID uint64) error {
	return s.execDeleteAttachments(ctx, squirrel.Eq{
		s.preprocessColumn("att.id", ""): store.PreprocessValue(ID, ""),
	})
}

// TruncateAttachments Deletes all rows from the attachments table
func (s Store) TruncateAttachments(ctx context.Context) error {
	return s.Truncate(ctx, s.attachmentTable())
}

// execLookupAttachment prepares Attachment query and executes it,
// returning types.Attachment (or error)
func (s Store) execLookupAttachment(ctx context.Context, cnd squirrel.Sqlizer) (res *types.Attachment, err error) {
	var (
		row rowScanner
	)

	row, err = s.QueryRow(ctx, s.attachmentsSelectBuilder().Where(cnd))
	if err != nil {
		return
	}

	res, err = s.internalAttachmentRowScanner(row)
	if err != nil {
		return
	}

	return res, nil
}

// execCreateAttachments updates all matched (by cnd) rows in attachments with given data
func (s Store) execCreateAttachments(ctx context.Context, payload store.Payload) error {
	return s.Exec(ctx, s.InsertBuilder(s.attachmentTable()).SetMap(payload))
}

// execUpdateAttachments updates all matched (by cnd) rows in attachments with given data
func (s Store) execUpdateAttachments(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return s.Exec(ctx, s.UpdateBuilder(s.attachmentTable("att")).Where(cnd).SetMap(set))
}

// execUpsertAttachments inserts new or updates matching (by-primary-key) rows in attachments with given data
func (s Store) execUpsertAttachments(ctx context.Context, set store.Payload) error {
	upsert, err := s.config.UpsertBuilder(
		s.config,
		s.attachmentTable(),
		set,
		s.preprocessColumn("id", ""),
	)

	if err != nil {
		return err
	}

	return s.Exec(ctx, upsert)
}

// execDeleteAttachments Deletes all matched (by cnd) rows in attachments with given data
func (s Store) execDeleteAttachments(ctx context.Context, cnd squirrel.Sqlizer) error {
	return s.Exec(ctx, s.DeleteBuilder(s.attachmentTable("att")).Where(cnd))
}

func (s Store) internalAttachmentRowScanner(row rowScanner) (res *types.Attachment, err error) {
	res = &types.Attachment{}

	if _, has := s.config.RowScanners["attachment"]; has {
		scanner := s.config.RowScanners["attachment"].(func(_ rowScanner, _ *types.Attachment) error)
		err = scanner(row, res)
	} else {
		err = row.Scan(
			&res.ID,
			&res.OwnerID,
			&res.Kind,
			&res.Url,
			&res.PreviewUrl,
			&res.Name,
			&res.Meta,
			&res.CreatedAt,
			&res.UpdatedAt,
			&res.DeletedAt,
		)
	}

	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err != nil {
		return nil, errors.Store("could not scan attachment db row: %s", err).Wrap(err)
	} else {
		return res, nil
	}
}

// QueryAttachments returns squirrel.SelectBuilder with set table and all columns
func (s Store) attachmentsSelectBuilder() squirrel.SelectBuilder {
	return s.SelectBuilder(s.attachmentTable("att"), s.attachmentColumns("att")...)
}

// attachmentTable name of the db table
func (Store) attachmentTable(aa ...string) string {
	var alias string
	if len(aa) > 0 {
		alias = " AS " + aa[0]
	}

	return "attachments" + alias
}

// AttachmentColumns returns all defined table columns
//
// With optional string arg, all columns are returned aliased
func (Store) attachmentColumns(aa ...string) []string {
	var alias string
	if len(aa) > 0 {
		alias = aa[0] + "."
	}

	return []string{
		alias + "id",
		alias + "rel_owner",
		alias + "kind",
		alias + "url",
		alias + "preview_url",
		alias + "name",
		alias + "meta",
		alias + "created_at",
		alias + "updated_at",
		alias + "deleted_at",
	}
}

// {true true false false false true}

// internalAttachmentEncoder encodes fields from types.Attachment to store.Payload (map)
//
// Encoding is done by using generic approach or by calling encodeAttachment
// func when rdbms.customEncoder=true
func (s Store) internalAttachmentEncoder(res *types.Attachment) store.Payload {
	return store.Payload{
		"id":          res.ID,
		"rel_owner":   res.OwnerID,
		"kind":        res.Kind,
		"url":         res.Url,
		"preview_url": res.PreviewUrl,
		"name":        res.Name,
		"meta":        res.Meta,
		"created_at":  res.CreatedAt,
		"updated_at":  res.UpdatedAt,
		"deleted_at":  res.DeletedAt,
	}
}

// checkAttachmentConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant but unfortunately we can not rely
// on the full support (MySQL does not support conditional indexes)
func (s *Store) checkAttachmentConstraints(ctx context.Context, res *types.Attachment) error {
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
