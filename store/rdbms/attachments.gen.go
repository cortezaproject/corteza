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
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/jmoiron/sqlx"
)

// SearchAttachments returns all matching rows
//
// This function calls convertAttachmentFilter with the given
// types.AttachmentFilter and expects to receive a working squirrel.SelectBuilder
func (s Store) SearchAttachments(ctx context.Context, f types.AttachmentFilter) (types.AttachmentSet, types.AttachmentFilter, error) {
	q, err := s.convertAttachmentFilter(f)
	if err != nil {
		return nil, f, err
	}

	scap := f.PerPage
	if scap == 0 {
		scap = DefaultSliceCapacity
	}

	if f.Count, err = Count(ctx, s.db, q); err != nil || f.Count == 0 {
		return nil, f, err
	}

	var (
		set = make([]*types.Attachment, 0, scap)
		// @todo this offset needs to be removed and replaced with key-based-paging
		fetchPage = func(offset, limit uint) (fetched, skipped uint, err error) {
			var (
				res *types.Attachment
				chk bool
			)

			if limit > 0 {
				q = q.Limit(uint64(limit))
			}

			if offset > 0 {
				q = q.Offset(uint64(offset))
			}

			rows, err := s.Query(ctx, q)
			if err != nil {
				return
			}

			for rows.Next() {
				fetched++
				if res, err = s.internalAttachmentRowScanner(rows, rows.Err()); err != nil {
					if cerr := rows.Close(); cerr != nil {
						err = fmt.Errorf("could not close rows (%v) after scan error: %w", cerr, err)
					}

					return
				}

				// If check function is set, call it and act accordingly
				if f.Check != nil {
					if chk, err = f.Check(res); err != nil {
						if cerr := rows.Close(); cerr != nil {
							err = fmt.Errorf("could not close rows (%v) after check error: %w", cerr, err)
						}

						return
					} else if !chk {
						// did not pass the check
						// go with the next row
						skipped++
						continue
					}
				}

				set = append(set, res)

				// make sure we do not fetch more than requested!
				if f.Limit > 0 && uint(len(set)) >= f.Limit {
					break
				}
			}

			err = rows.Close()
			return
		}

		fetch = func() error {
			var (
				fetched uint

				// starting offset & limit are from filter arg
				// note that this will have to be improved with key-based pagination
				offset, limit = calculatePaging(f.PageFilter)
			)

			for refetch := 0; refetch < MaxRefetches; refetch++ {
				if fetched, _, err = fetchPage(offset, limit); err != nil {
					return err
				}

				// if limit is not set or we've already collected enough resources
				// we can break the loop right away
				if limit == 0 || fetched == 0 || uint(len(set)) >= f.Limit {
					break
				}

				// we've skipped fetched resources (due to check() fn)
				// and we still have less results (in set) than required by limit
				// inc offset by number of fetched items
				offset += fetched

				if limit < MinRefetchLimit {
					limit = MinRefetchLimit
				}

			}
			return nil
		}
	)

	return set, f, fetch()
}

// LookupAttachmentByID searches for attachment by its ID
//
// It returns attachment even if deleted
func (s Store) LookupAttachmentByID(ctx context.Context, id uint64) (*types.Attachment, error) {
	return s.AttachmentLookup(ctx, squirrel.Eq{
		"att.id": id,
	})
}

// CreateAttachment creates one or more rows in attachments table
func (s Store) CreateAttachment(ctx context.Context, rr ...*types.Attachment) error {
	if len(rr) == 0 {
		return nil
	}

	return Tx(ctx, s.db, s.config, nil, func(db *sqlx.Tx) (err error) {
		for _, res := range rr {
			err = ExecuteSqlizer(ctx, s.DB(), s.Insert(s.AttachmentTable()).SetMap(s.internalAttachmentEncoder(res)))
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// UpdateAttachment updates one or more existing rows in attachments
func (s Store) UpdateAttachment(ctx context.Context, rr ...*types.Attachment) error {
	return s.PartialUpdateAttachment(ctx, nil, rr...)
}

// PartialUpdateAttachment updates one or more existing rows in attachments
//
// It wraps the update into transaction and can perform partial update by providing list of updatable columns
func (s Store) PartialUpdateAttachment(ctx context.Context, onlyColumns []string, rr ...*types.Attachment) error {
	if len(rr) == 0 {
		return nil
	}

	return Tx(ctx, s.db, s.config, nil, func(db *sqlx.Tx) (err error) {
		for _, res := range rr {
			err = s.ExecUpdateAttachments(
				ctx,
				squirrel.Eq{s.preprocessColumn("att.id", ""): s.preprocessValue(res.ID, "")},
				s.internalAttachmentEncoder(res).Skip("id").Only(onlyColumns...))
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// RemoveAttachment removes one or more rows from attachments table
func (s Store) RemoveAttachment(ctx context.Context, rr ...*types.Attachment) error {
	if len(rr) == 0 {
		return nil
	}

	return Tx(ctx, s.db, s.config, nil, func(db *sqlx.Tx) (err error) {
		for _, res := range rr {
			err = ExecuteSqlizer(ctx, s.DB(), s.Delete(s.AttachmentTable("att")).Where(squirrel.Eq{s.preprocessColumn("att.id", ""): s.preprocessValue(res.ID, "")}))
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// RemoveAttachmentByID removes row from the attachments table
func (s Store) RemoveAttachmentByID(ctx context.Context, ID uint64) error {
	return ExecuteSqlizer(ctx, s.DB(), s.Delete(s.AttachmentTable("att")).Where(squirrel.Eq{s.preprocessColumn("att.id", ""): s.preprocessValue(ID, "")}))
}

// TruncateAttachments removes all rows from the attachments table
func (s Store) TruncateAttachments(ctx context.Context) error {
	return Truncate(ctx, s.DB(), s.AttachmentTable())
}

// ExecUpdateAttachments updates all matched (by cnd) rows in attachments with given data
func (s Store) ExecUpdateAttachments(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return ExecuteSqlizer(ctx, s.DB(), s.Update(s.AttachmentTable("att")).Where(cnd).SetMap(set))
}

// AttachmentLookup prepares Attachment query and executes it,
// returning types.Attachment (or error)
func (s Store) AttachmentLookup(ctx context.Context, cnd squirrel.Sqlizer) (*types.Attachment, error) {
	return s.internalAttachmentRowScanner(s.QueryRow(ctx, s.QueryAttachments().Where(cnd)))
}

func (s Store) internalAttachmentRowScanner(row rowScanner, err error) (*types.Attachment, error) {
	if err != nil {
		return nil, err
	}

	var res = &types.Attachment{}
	if _, has := s.config.RowScanners["attachment"]; has {
		scanner := s.config.RowScanners["attachment"].(func(rowScanner, *types.Attachment) error)
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
func (s Store) QueryAttachments() squirrel.SelectBuilder {
	return s.Select(s.AttachmentTable("att"), s.AttachmentColumns("att")...)
}

// AttachmentTable name of the db table
func (Store) AttachmentTable(aa ...string) string {
	var alias string
	if len(aa) > 0 {
		alias = " AS " + aa[0]
	}

	return "attachments" + alias
}

// AttachmentColumns returns all defined table columns
//
// With optional string arg, all columns are returned aliased
func (Store) AttachmentColumns(aa ...string) []string {
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
