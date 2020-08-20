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
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
)

var _ = errors.Is

const (
	TriggerBeforeAttachmentCreate triggerKey = "attachmentBeforeCreate"
	TriggerBeforeAttachmentUpdate triggerKey = "attachmentBeforeUpdate"
	TriggerBeforeAttachmentUpsert triggerKey = "attachmentBeforeUpsert"
	TriggerBeforeAttachmentDelete triggerKey = "attachmentBeforeDelete"
)

// SearchAttachments returns all matching rows
//
// This function calls convertAttachmentFilter with the given
// types.AttachmentFilter and expects to receive a working squirrel.SelectBuilder
func (s Store) SearchAttachments(ctx context.Context, f types.AttachmentFilter) (types.AttachmentSet, types.AttachmentFilter, error) {
	var scap uint
	q, err := s.convertAttachmentFilter(f)
	if err != nil {
		return nil, f, err
	}

	if scap == 0 {
		scap = DefaultSliceCapacity
	}

	var (
		set = make([]*types.Attachment, 0, scap)
		// Paging is disabled in definition yaml file
		// {search: {enablePaging:false}} and this allows
		// a much simpler row fetching logic
		fetch = func() error {
			var (
				res       *types.Attachment
				rows, err = s.Query(ctx, q)
			)

			if err != nil {
				return err
			}

			for rows.Next() {
				if rows.Err() == nil {
					res, err = s.internalAttachmentRowScanner(rows)
				}

				if err != nil {
					if cerr := rows.Close(); cerr != nil {
						err = fmt.Errorf("could not close rows (%v) after scan error: %w", cerr, err)
					}

					return err
				}

				// If check function is set, call it and act accordingly

				if f.Check != nil {
					if chk, err := f.Check(res); err != nil {
						if cerr := rows.Close(); cerr != nil {
							err = fmt.Errorf("could not close rows (%v) after check error: %w", cerr, err)
						}

						return err
					} else if !chk {
						// did not pass the check
						// go with the next row
						continue
					}
				}
				set = append(set, res)
			}

			return rows.Close()
		}
	)

	return set, f, s.config.ErrorHandler(fetch())
}

// LookupAttachmentByID searches for attachment by its ID
//
// It returns attachment even if deleted
func (s Store) LookupAttachmentByID(ctx context.Context, id uint64) (*types.Attachment, error) {
	return s.execLookupAttachment(ctx, squirrel.Eq{
		s.preprocessColumn("att.id", ""): s.preprocessValue(id, ""),
	})
}

// CreateAttachment creates one or more rows in attachments table
func (s Store) CreateAttachment(ctx context.Context, rr ...*types.Attachment) (err error) {
	for _, res := range rr {
		err = s.checkAttachmentConstraints(ctx, res)
		if err != nil {
			return err
		}

		// err = s.attachmentHook(ctx, TriggerBeforeAttachmentCreate, res)
		// if err != nil {
		// 	return err
		// }

		err = s.execCreateAttachments(ctx, s.internalAttachmentEncoder(res))
		if err != nil {
			return err
		}
	}

	return
}

// UpdateAttachment updates one or more existing rows in attachments
func (s Store) UpdateAttachment(ctx context.Context, rr ...*types.Attachment) error {
	return s.config.ErrorHandler(s.PartialAttachmentUpdate(ctx, nil, rr...))
}

// PartialAttachmentUpdate updates one or more existing rows in attachments
func (s Store) PartialAttachmentUpdate(ctx context.Context, onlyColumns []string, rr ...*types.Attachment) (err error) {
	for _, res := range rr {
		err = s.checkAttachmentConstraints(ctx, res)
		if err != nil {
			return err
		}

		// err = s.attachmentHook(ctx, TriggerBeforeAttachmentUpdate, res)
		// if err != nil {
		// 	return err
		// }

		err = s.execUpdateAttachments(
			ctx,
			squirrel.Eq{
				s.preprocessColumn("att.id", ""): s.preprocessValue(res.ID, ""),
			},
			s.internalAttachmentEncoder(res).Skip("id").Only(onlyColumns...))
		if err != nil {
			return s.config.ErrorHandler(err)
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

		// err = s.attachmentHook(ctx, TriggerBeforeAttachmentUpsert, res)
		// if err != nil {
		// 	return err
		// }

		err = s.config.ErrorHandler(s.execUpsertAttachments(ctx, s.internalAttachmentEncoder(res)))
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteAttachment Deletes one or more rows from attachments table
func (s Store) DeleteAttachment(ctx context.Context, rr ...*types.Attachment) (err error) {
	for _, res := range rr {
		// err = s.attachmentHook(ctx, TriggerBeforeAttachmentDelete, res)
		// if err != nil {
		// 	return err
		// }

		err = s.execDeleteAttachments(ctx, squirrel.Eq{
			s.preprocessColumn("att.id", ""): s.preprocessValue(res.ID, ""),
		})
		if err != nil {
			return s.config.ErrorHandler(err)
		}
	}

	return nil
}

// DeleteAttachmentByID Deletes row from the attachments table
func (s Store) DeleteAttachmentByID(ctx context.Context, ID uint64) error {
	return s.execDeleteAttachments(ctx, squirrel.Eq{
		s.preprocessColumn("att.id", ""): s.preprocessValue(ID, ""),
	})
}

// TruncateAttachments Deletes all rows from the attachments table
func (s Store) TruncateAttachments(ctx context.Context) error {
	return s.config.ErrorHandler(s.Truncate(ctx, s.attachmentTable()))
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
	return s.config.ErrorHandler(s.Exec(ctx, s.InsertBuilder(s.attachmentTable()).SetMap(payload)))
}

// execUpdateAttachments updates all matched (by cnd) rows in attachments with given data
func (s Store) execUpdateAttachments(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return s.config.ErrorHandler(s.Exec(ctx, s.UpdateBuilder(s.attachmentTable("att")).Where(cnd).SetMap(set)))
}

// execUpsertAttachments inserts new or updates matching (by-primary-key) rows in attachments with given data
func (s Store) execUpsertAttachments(ctx context.Context, set store.Payload) error {
	upsert, err := s.config.UpsertBuilder(
		s.config,
		s.attachmentTable(),
		set,
		"id",
	)

	if err != nil {
		return err
	}

	return s.config.ErrorHandler(s.Exec(ctx, upsert))
}

// execDeleteAttachments Deletes all matched (by cnd) rows in attachments with given data
func (s Store) execDeleteAttachments(ctx context.Context, cnd squirrel.Sqlizer) error {
	return s.config.ErrorHandler(s.Exec(ctx, s.DeleteBuilder(s.attachmentTable("att")).Where(cnd)))
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
		return nil, store.ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("could not scan db row for Attachment: %w", err)
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

// {true true false false true}

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

func (s *Store) checkAttachmentConstraints(ctx context.Context, res *types.Attachment) error {

	return nil
}

// func (s *Store) attachmentHook(ctx context.Context, key triggerKey, res *types.Attachment) error {
// 	if fn, has := s.config.TriggerHandlers[key]; has {
// 		return fn.(func (ctx context.Context, s *Store, res *types.Attachment) error)(ctx, s, res)
// 	}
//
// 	return nil
// }
