package repository

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/titpetric/factory"
	"gopkg.in/Masterminds/squirrel.v1"

	"github.com/crusttech/crust/crm/repository/ql"

	"github.com/crusttech/crust/crm/types"
)

type (
	RecordRepository interface {
		With(ctx context.Context, db *factory.DB) RecordRepository

		FindByID(id uint64) (*types.Record, error)

		Report(moduleID uint64, metrics, dimensions, filter string) (results interface{}, err error)
		Find(moduleID uint64, filter string, sort string, page int, perPage int) (*FindResponse, error)

		Create(mod *types.Record) (*types.Record, error)
		Update(mod *types.Record) (*types.Record, error)
		DeleteByID(id uint64) error

		Fields(mod *types.Record) ([]*types.RecordColumn, error)
	}

	FindResponseMeta struct {
		Filter  string `json:"filter,omitempty"`
		Page    int    `json:"page"`
		PerPage int    `json:"perPage"`
		Count   int    `json:"count"`
		Sort    string `json:"sort"`
	}

	FindResponse struct {
		Meta    FindResponseMeta `json:"meta"`
		Records []*types.Record  `json:"records"`
	}

	record struct {
		*repository
	}
)

const (
	jsonWrap = `JSON_UNQUOTE(JSON_EXTRACT(json, REPLACE(JSON_UNQUOTE(JSON_SEARCH(json, 'one', ?)), '.name', '.value')))`
)

func Record(ctx context.Context, db *factory.DB) RecordRepository {
	return (&record{}).With(ctx, db)
}

func (r *record) With(ctx context.Context, db *factory.DB) RecordRepository {
	return &record{
		repository: r.repository.With(ctx, db),
	}
}

// @todo: update to accepted DeletedAt column semantics from SAM

func (r *record) FindByID(id uint64) (*types.Record, error) {
	mod := &types.Record{}
	if err := r.db().Get(mod, "SELECT * FROM crm_record WHERE id=? and deleted_at IS NULL", id); err != nil {
		return nil, err
	}
	return mod, nil
}

func (r *record) Report(moduleID uint64, metrics, dimensions, filter string) (results interface{}, err error) {
	crb := NewRecordReportBuilder(moduleID)

	if err = crb.SetMetrics(metrics); err != nil {
		return
	}

	if err = crb.SetDimensions(dimensions); err != nil {
		return
	}

	if err = crb.SetFilter(filter); err != nil {
		return
	}

	var result = make([]map[string]interface{}, 0)

	if query, args, err := crb.Build(); err != nil {
		return nil, errors.Wrap(err, "Can not generate report query")
	} else if rows, err := r.db().Query(query, args...); err != nil {
		return nil, errors.Wrapf(err, "Can not execute report query (%s)", query)
	} else {
		for rows.Next() {
			result = append(result, crb.Cast(rows))
		}

		return result, nil
	}
}

func (r *record) Find(moduleID uint64, filter string, sort string, page int, perPage int) (*FindResponse, error) {
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
			Filter:  filter,
			Page:    page,
			PerPage: perPage,
			Sort:    sort,
		},
		Records: make([]*types.Record, 0),
	}

	// Create query for fetching and counting records.
	query := squirrel.
		Select().
		From("crm_record").
		Where("(module_id = ? AND deleted_at IS NULL AND json IS NOT NULL)", moduleID)

	// Parse filters.
	p := ql.NewParser()
	p.OnIdent = ql.MakeIdentWrapHandler(jsonWrap, "created_at", "updated_at", "id", "user_id")

	where, err := p.ParseExpression(filter)
	if err != nil {
		return nil, err
	}

	// Append filtering to query.
	query = query.Where(squirrel.And{where})

	// Create count SQL sentences.
	count := query.Column(squirrel.Alias(squirrel.Expr("COUNT(*)"), "count"))
	sqlSelect, argsSelect, err := count.ToSql()
	if err != nil {
		return nil, err
	}

	// Execute count query.
	if err := r.db().Get(&response.Meta.Count, sqlSelect, argsSelect...); err != nil {
		return nil, err
	}

	// Return empty response if count of records is zero.
	if response.Meta.Count == 0 {
		return response, nil
	}

	// Create query for fetching records.
	query = query.
		Column("*").
		Limit(uint64(perPage)).
		Offset(uint64(page))

	// Append Sorting.
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

		query = query.OrderBy(orderFields...)
	}

	/*
		p = ql.NewParser()
		p.OnIdent = ql.MakeIdentWrapHandler(jsonWrap, "id", "module_id", "user_id", "created_at", "updated_at")

		order, err := p.ParseExpression(sort)
		if err != nil {
			return nil, err
		}

		sqlOrder, argsOrder, err := order.ToSql()
		if err != nil {
			return nil, err
		}
		query = query.OrderBy(sqlOrder)
	*/

	// Create actual fetch SQL sentences.
	sqlSelect, argsSelect, err = query.ToSql()
	if err != nil {
		return nil, err
	}

	// Append order args to select args and execute actual query.
	// argsSelect = append(argsSelect, argsOrder...)
	if err := r.db().Select(&response.Records, sqlSelect, argsSelect...); err != nil {
		return nil, err
	}

	return response, nil
}

func (r *record) Create(mod *types.Record) (*types.Record, error) {
	mod.ID = factory.Sonyflake.NextID()
	mod.CreatedAt = time.Now()
	mod.UserID = Identity(r.Context())

	fields := make([]types.RecordColumn, 0)
	if err := json.Unmarshal(mod.Fields, &fields); err != nil {
		return nil, errors.Wrap(err, "No content")
	}

	r.db().Exec("delete from crm_record_links where record_id=?", mod.ID)
	for _, v := range fields {
		v.RecordID = mod.ID
		if err := r.db().Replace("crm_record_column", v); err != nil {
			return nil, errors.Wrap(err, "Error adding columns")
		}
		for _, related := range v.Related {
			row := types.Related{
				RecordID:        v.RecordID,
				Name:            v.Name,
				RelatedRecordID: related,
			}
			if err := r.db().Replace("crm_record_links", row); err != nil {
				return nil, errors.Wrap(err, "Error adding column links")
			}
		}
	}

	if err := r.db().Insert("crm_record", mod); err != nil {
		return nil, err
	}
	return mod, nil
}

func (r *record) Update(mod *types.Record) (*types.Record, error) {
	now := time.Now()
	mod.UpdatedAt = &now

	fields := make([]types.RecordColumn, 0)
	if err := json.Unmarshal(mod.Fields, &fields); err != nil {
		return nil, errors.Wrap(err, "Error when saving record, no content")
	}

	r.db().Exec("delete from crm_record_links where record_id=?", mod.ID)
	for _, v := range fields {
		v.RecordID = mod.ID
		if err := r.db().Replace("crm_record_column", v); err != nil {
			return nil, errors.Wrap(err, "Error adding columns to database")
		}
		for _, related := range v.Related {
			row := types.Related{
				RecordID:        v.RecordID,
				Name:            v.Name,
				RelatedRecordID: related,
			}
			if err := r.db().Replace("crm_record_links", row); err != nil {
				return nil, errors.Wrap(err, "Error adding column links")
			}
		}
	}

	return mod, r.db().Replace("crm_record", mod)
}

func (r *record) DeleteByID(id uint64) error {
	_, err := r.db().Exec("update crm_record set deleted_at=? where id=?", time.Now(), id)
	return err
}

func (r *record) Fields(record *types.Record) ([]*types.RecordColumn, error) {
	result := make([]*types.RecordColumn, 0)
	module := Module(r.ctx, r.db())

	mod, err := module.FindByID(record.ModuleID)
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
		record.ID,
	}
	for _, v := range fieldNames {
		args = append(args, v)
	}
	return result, r.db().Select(&result, "select * FROM crm_record_column where record_id=? order by "+order, args...)
}
