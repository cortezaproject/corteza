package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"encoding/json"

	"github.com/pkg/errors"
	"github.com/titpetric/factory"

	"github.com/crusttech/crust/crm/types"
	systemRepository "github.com/crusttech/crust/system/repository"
)

type (
	ContentRepository interface {
		With(ctx context.Context, db *factory.DB) ContentRepository

		FindByID(id uint64) (*types.Content, error)

		Find(moduleID uint64, query string, page int, perPage int) (*FindResponse, error)

		Create(mod *types.Content) (*types.Content, error)
		Update(mod *types.Content) (*types.Content, error)
		DeleteByID(id uint64) error

		Fields(mod *types.Content) ([]*types.ContentColumn, error)
	}

	FindResponseMeta struct {
		Query   string `json:"query,omitempty"`
		Page    int    `json:"page"`
		PerPage int    `json:"perPage"`
		Count   int    `json:"count"`
	}

	FindResponse struct {
		Meta     FindResponseMeta `json:"meta"`
		Contents []*types.Content `json:"contents"`
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
	if err := r.db().Get(mod, "SELECT * FROM crm_content WHERE id=? and deleted_at IS NULL", id); err != nil {
		return nil, err
	}
	return mod, r.prepare(mod, "page", "user", "fields")
}

func (r *content) Find(moduleID uint64, query string, page int, perPage int) (*FindResponse, error) {
	if page < 0 {
		page = 0
	}
	if perPage <= 0 {
		perPage = 50
	}
	if perPage > 100 {
		perPage = 100
	}
	if perPage < 10 {
		perPage = 10
	}
	response := &FindResponse{
		Meta: FindResponseMeta{
			Page:    page,
			PerPage: perPage,
			Query:   query,
		},
		Contents: make([]*types.Content, 0),
	}

	query = "%" + query + "%"

	sqlSelect := "SELECT * FROM crm_content"
	sqlCount := "SELECT count(*) FROM crm_content"
	sqlWhere := "WHERE module_id=? and deleted_at IS NULL"

	sqlOrder := "ORDER BY id DESC"
	sqlLimit := fmt.Sprintf("LIMIT %d, %d", page, perPage)

	switch true {
	case query != "":
		sqlWhere = sqlWhere + " AND id in (select distinct content_id from crm_content_column where column_value like ?)"
		if err := r.db().Get(&response.Meta.Count, sqlCount+" "+sqlWhere, moduleID, query); err != nil {
			return nil, err
		}
		if err := r.db().Select(&response.Contents, sqlSelect+" "+sqlWhere+" "+sqlOrder+" "+sqlLimit, moduleID, query); err != nil {
			return nil, err
		}
	default:
		if err := r.db().Get(&response.Meta.Count, sqlCount+" "+sqlWhere, moduleID); err != nil {
			return nil, err
		}
		if err := r.db().Select(&response.Contents, fmt.Sprintf("SELECT * FROM crm_content WHERE module_id=? and deleted_at IS NULL ORDER BY id DESC LIMIT %d, %d", page, perPage), moduleID); err != nil {
			return nil, err
		}
	}

	if err := r.prepareAll(response.Contents, "user", "fields"); err != nil {
		return nil, err
	}

	return response, nil
}

func (r *content) Create(mod *types.Content) (*types.Content, error) {
	mod.ID = factory.Sonyflake.NextID()
	mod.CreatedAt = time.Now()
	mod.UserID = Identity(r.Context())

	fields := make([]types.ContentColumn, 0)
	if err := json.Unmarshal(mod.Fields, &fields); err != nil {
		return nil, errors.Wrap(err, "No content")
	}

	r.db().Exec("delete from crm_content_links where content_id=?", mod.ID)
	for _, v := range fields {
		v.ContentID = mod.ID
		if err := r.db().Replace("crm_content_column", v); err != nil {
			return nil, errors.Wrap(err, "Error adding columns")
		}
		for _, related := range v.Related {
			row := types.Related{
				ContentID:        v.ContentID,
				Name:             v.Name,
				RelatedContentID: related,
			}
			if err := r.db().Replace("crm_content_links", row); err != nil {
				return nil, errors.Wrap(err, "Error adding column links")
			}
		}
	}

	if err := r.db().Insert("crm_content", mod); err != nil {
		return nil, err
	}

	return mod, r.prepare(mod, "user", "fields")
}

func (r *content) Update(mod *types.Content) (*types.Content, error) {
	now := time.Now()
	mod.UpdatedAt = &now

	fields := make([]types.ContentColumn, 0)
	if err := json.Unmarshal(mod.Fields, &fields); err != nil {
		return nil, errors.Wrap(err, "Error when saving content, no content")
	}

	r.db().Exec("delete from crm_content_links where content_id=?", mod.ID)
	for _, v := range fields {
		v.ContentID = mod.ID
		if err := r.db().Replace("crm_content_column", v); err != nil {
			return nil, errors.Wrap(err, "Error adding columns to database")
		}
		for _, related := range v.Related {
			row := types.Related{
				ContentID:        v.ContentID,
				Name:             v.Name,
				RelatedContentID: related,
			}
			if err := r.db().Replace("crm_content_links", row); err != nil {
				return nil, errors.Wrap(err, "Error adding column links")
			}
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
	return result, r.db().Select(&result, "select * FROM crm_content_column where content_id=? order by "+order, args...)
}

func (r *content) prepareAll(contents []*types.Content, fields ...string) error {
	for _, content := range contents {
		if err := r.prepare(content, fields...); err != nil {
			return err
		}
	}
	return nil
}

func (r *content) prepare(content *types.Content, fields ...string) (err error) {
	api := Page(r.Context(), r.db())
	usersAPI := systemRepository.User(r.Context(), r.db())
	for _, field := range fields {
		switch field {
		case "fields":
			fields, err := r.Fields(content)
			if err != nil {
				return err
			}
			json, err := json.Marshal(fields)
			if err != nil {
				return err
			}
			if err := (&content.Fields).Scan(json); err != nil {
				return err
			}
		case "page":
			if content.Page, err = api.FindByModuleID(content.ModuleID); err != nil {
				return
			}
		case "user":
			if content.UserID > 0 {
				if content.User, err = usersAPI.FindByID(content.UserID); err != nil {
					return
				}
			}
		}
	}
	return
}
