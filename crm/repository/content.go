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
)

type (
	ContentRepository interface {
		With(ctx context.Context, db *factory.DB) ContentRepository

		FindByID(id uint64) (*types.Content, error)

		Find(moduleID uint64, query string, page int, perPage int, sort string) (*FindResponse, error)

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
		Sort    string `json:"sort"`
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
	return mod, nil
}

func (r *content) Find(moduleID uint64, query string, page int, perPage int, sort string) (*FindResponse, error) {
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
			Sort:    sort,
		},
		Contents: make([]*types.Content, 0),
	}

	query = "%" + query + "%"

	sqlSelect := "SELECT * FROM crm_content"
	sqlCount := "SELECT count(*) FROM crm_content"
	sqlWhere := "WHERE module_id=? and deleted_at IS NULL"
	sqlOrder := "ORDER BY id DESC"
	sqlLimit := fmt.Sprintf("LIMIT %d, %d", page*perPage, perPage)

	chuncks := strings.Split(sort, ",")
	if len(chuncks) > 0 {

		// Ger module fields.
		modulRepo := Module(r.Context(), r.db())
		mod, err := modulRepo.FindByID(moduleID)
		if err != nil {
			return nil, err
		}
		modFields, err := modulRepo.FieldNames(mod)
		if err != nil {
			return nil, err
		}
		fieldMap := make(map[string]bool)
		for i := 0; i < len(modFields); i++ {
			fieldMap[modFields[i]] = true
		}

		orderFields := make([]string, 0)
		for _, c := range chuncks {
			args := strings.Split(c, " ")

			var field string
			if _, ok := fieldMap[args[0]]; ok {
				field = "JSON_UNQUOTE(JSON_EXTRACT(json, REPLACE(JSON_UNQUOTE(JSON_SEARCH(json, 'one', '" + args[0] + "')), '.name', '.value')))"
			} else {
				switch args[0] {
				case "moduleId":
					field = "module_id"
				case "userId":
					field = "user_id"
				case "createdAt":
					field = "created_at"
				case "updatedAt":
					field = "updated_at"
				case "deletedAt":
					field = "deleted_at"
				default:
					field = "id"
				}
			}

			// Check for second order parameter or use default value ASC.
			order := "ASC"
			if len(args) == 2 {
				order = strings.ToUpper(args[1])
				switch order {
				case "DESC":
					order = "DESC"
				default:
					order = "ASC"
				}
			}

			// We skip batch of parameters if there are more then 2 values.
			if len(args) > 2 {
				continue
			}

			// Add field and order to sort order fields.
			orderFields = append(orderFields, field+" "+order)
		}

		sqlOrder = "ORDER BY " + strings.Join(orderFields, ", ")
	}

	// One possibility to order by field value without JSON, is query written bellow with FIELD over column names and order by value:
	// SELECT * FROM crm_content
	// LEFT JOIN crm_content_column ON crm_content.id = crm_content_column.content_id"
	// WHERE column_name in ('name', 'email')
	// ORDER BY FIELD(column_name, 'email', 'name'), column_value;

	// Possibility to order with JSON:
	// SELECT *,
	// JSON_UNQUOTE(JSON_EXTRACT(json, REPLACE(JSON_UNQUOTE(JSON_SEARCH(json, 'all', 'email')), '.name', '.value'))) as emailField
	// FROM crm_content
	// ORDER by emailField asc;

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
		if err := r.db().Select(&response.Contents, sqlSelect+" "+sqlWhere+" "+sqlOrder+" "+sqlLimit, moduleID); err != nil {
			return nil, err
		}
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
	return mod, nil
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
