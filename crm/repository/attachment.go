package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/titpetric/factory"
	sq "gopkg.in/Masterminds/squirrel.v1"

	"github.com/crusttech/crust/crm/types"
)

type (
	AttachmentRepository interface {
		With(ctx context.Context, db *factory.DB) AttachmentRepository

		Find(filter types.AttachmentFilter) (types.AttachmentSet, types.AttachmentFilter, error)
		FindByID(id uint64) (*types.Attachment, error)
		FindByIDs(IDs ...uint64) (types.AttachmentSet, error)
		Create(mod *types.Attachment) (*types.Attachment, error)
		DeleteByID(id uint64) error
	}

	attachment struct {
		*repository
	}
)

const (
	sqlAttachmentColumns = `
		a.id, a.rel_owner, a.kind,
		a.url, a.preview_url,
		a.name,
		a.meta,
		a.created_at, a.updated_at,	a.deleted_at
	`
	sqlAttachmentScope = "deleted_at IS NULL"

	sqlAttachmentByID = `SELECT ` + sqlAttachmentColumns +
		` FROM crm_attachment AS a WHERE id = ? AND ` + sqlAttachmentScope

	sqlAttachmentsByIDs = `SELECT ` + sqlAttachmentColumns +
		` FROM crm_attachment AS a WHERE id IN (?) AND ` + sqlAttachmentScope
)

func Attachment(ctx context.Context, db *factory.DB) AttachmentRepository {
	return (&attachment{}).With(ctx, db)
}

func (r *attachment) With(ctx context.Context, db *factory.DB) AttachmentRepository {
	return &attachment{
		repository: r.repository.With(ctx, db),
	}
}

func (r *attachment) FindByID(id uint64) (*types.Attachment, error) {
	mod := &types.Attachment{}

	return mod, r.db().Get(mod, sqlAttachmentByID, id)
}

func (r *attachment) FindByIDs(IDs ...uint64) (rval types.AttachmentSet, err error) {
	rval = make([]*types.Attachment, 0)

	if len(IDs) == 0 {
		return
	}

	if sql, args, err := sqlx.In(sqlAttachmentsByIDs, IDs); err != nil {
		return nil, err
	} else {
		return rval, r.db().Select(&rval, sql, args...)
	}
}

func (r *attachment) Find(filter types.AttachmentFilter) (set types.AttachmentSet, f types.AttachmentFilter, err error) {
	f = filter
	if f.PerPage > 100 {
		f.PerPage = 100
	} else if f.PerPage == 0 {
		f.PerPage = 50
	}

	set = types.AttachmentSet{}

	query := sq.Select().From("crm_attachment AS a").Where(sq.Eq{"a.kind": f.Kind})

	switch f.Kind {
	case types.PageAttachment:
		// @todo implement filtering by page
		if f.PageID > 0 {
			err = errors.New("filtering by pageID not implemented")
			return
		}
	case types.RecordAttachment:
		query = query.
			Join("crm_record_value AS v ON (v.ref = a.id)")

		if f.ModuleID > 0 {
			query = query.
				Join("crm_record AS r ON (r.id = v.record_id)").
				Where(sq.Eq{"r.module_id": f.ModuleID})
		}

		if f.RecordID > 0 {
			query = query.Where(sq.Eq{"v.record_id": f.RecordID})
		}

		if f.FieldName != "" {
			query = query.Where(sq.Eq{"v.name": f.FieldName})
		}
	default:
		err = errors.New("unsupported kind value")
	}

	if f.Filter != "" {
		err = errors.New("filtering by filter not implemented")
		return
	}

	// Assemble SQL for counting (includes only where)
	count := query.Column("COUNT(*)")
	if sqlSelect, argsSelect, err := count.ToSql(); err != nil {
		return set, f, err
	} else {
		// Execute count query.
		if err := r.db().Get(&f.Count, sqlSelect, argsSelect...); err != nil {
			return set, f, err
		}

		// Return empty response if count of records is zero.
		if f.Count == 0 {
			return set, f, nil
		}
	}

	// Assemble SQL for fetching attachments (where + sorting + paging)...
	query = query.
		Column(sqlAttachmentColumns).
		Limit(uint64(f.PerPage)).
		Offset(uint64(f.Page * f.PerPage))

	if sqlSelect, argsSelect, err := query.ToSql(); err != nil {
		return set, f, err
	} else {
		return set, f, r.db().Select(&set, sqlSelect, argsSelect...)
	}
}

func (r *attachment) Create(mod *types.Attachment) (*types.Attachment, error) {
	if mod.ID == 0 {
		mod.ID = factory.Sonyflake.NextID()
	}

	mod.CreatedAt = time.Now()

	return mod, r.db().Insert("crm_attachment", mod)
}

func (r *attachment) DeleteByID(id uint64) error {
	_, err := r.db().Exec("UPDATE crm_attachment SET deleted_at = NOW() WHERE id = ?", id)
	return err
}
