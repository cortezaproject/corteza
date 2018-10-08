package repository

import (
	"context"
	"strings"
	"time"

	"encoding/json"

	"github.com/pkg/errors"
	"github.com/titpetric/factory"

	"github.com/crusttech/crust/crm/types"
)

type (
	ContentRepository interface {
		With(ctx context.Context, db *factory.DB) ContentRepository

		FindByID(id uint64) (*types.Content, error)
		Find() ([]*types.Content, error)
		Create(mod *types.Content) (*types.Content, error)
		Update(mod *types.Content) (*types.Content, error)
		DeleteByID(id uint64) error

		Fields(mod *types.Content) ([]*types.ContentColumn, error)
	}

	content struct {
		*repository
	}
)

func Content(ctx context.Context, db *factory.DB) ContentRepository {
	return (&content{}).With(ctx, db)
}

func (r *content) With(ctx context.Context, db *factory.DB) ContentRepository {
	return &content{
		repository: r.repository.With(ctx, db),
	}
}

// @todo: update to accepted DeletedAt column semantics from SAM

func (r *content) FindByID(id uint64) (*types.Content, error) {
	mod := &types.Content{}
	return mod, r.db().Get(mod, "SELECT * FROM crm_content WHERE id=? and deleted_at IS NULL", id)
}

func (r *content) Find() ([]*types.Content, error) {
	mod := make([]*types.Content, 0)
	return mod, r.db().Select(&mod, "SELECT * FROM crm_content WHERE deleted_at IS NULL ORDER BY id DESC")
}

func (r *content) Create(mod *types.Content) (*types.Content, error) {
	mod.ID = factory.Sonyflake.NextID()
	mod.CreatedAt = time.Now()

	fields := make([]types.ContentColumn, 0)
	if err := json.Unmarshal(mod.Fields, &fields); err != nil {
		return nil, errors.Wrap(err, "No content")
	}

	for _, v := range fields {
		v.ContentID = mod.ID
		if err := r.db().Replace("crm_content_column", v); err != nil {
			return nil, errors.Wrap(err, "Error adding columns")
		}
	}

	return mod, r.db().Insert("crm_content", mod)
}

func (r *content) Update(mod *types.Content) (*types.Content, error) {
	now := time.Now()
	mod.UpdatedAt = &now

	fields := make([]types.ContentColumn, 0)
	if err := json.Unmarshal(mod.Fields, &fields); err != nil {
		return nil, errors.Wrap(err, "Error when saving content, no content")
	}

	for _, v := range fields {
		v.ContentID = mod.ID
		if err := r.db().Replace("crm_content_column", v); err != nil {
			return nil, errors.Wrap(err, "Error adding columns to database")
		}
	}

	return mod, r.db().Replace("crm_content", mod)
}

func (r *content) DeleteByID(id uint64) error {
	_, err := r.db().Exec("update crm_content set deleted_at=? where id=?", time.Now(), id)
	return err
}

func (r *content) Fields(content *types.Content) ([]*types.ContentColumn, error) {
	result := make([]*types.ContentColumn, 0)
	module := Module(r.ctx, r.db())

	mod, err := module.FindByID(content.ModuleID)
	if err != nil {
		return result, err
	}

	fieldNames, err := module.FieldNames(mod)
	if err != nil {
		return result, err
	}
	if len(fieldNames) == 0 {
		return result, errors.New("Module has no fields")
	}

	order := "FIELD(column_name" + strings.Repeat(",?", len(fieldNames)) + ")"
	args := []interface{}{
		content.ID,
	}
	for _, v := range fieldNames {
		args = append(args, v)
	}
	return result, r.db().Select(&result, "select * from crm_content_column where content_id=? order by "+order, args...)
}
