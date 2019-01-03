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
		Find(module *types.Module, filter string, sort string, page int, perPage int) (*FindResponse, error)

		Create(mod *types.Record) (*types.Record, error)
		Update(mod *types.Record) (*types.Record, error)
		DeleteByID(id uint64) error

		Fields(module *types.Module, record *types.Record) ([]*types.RecordColumn, error)
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

func (r *record) Find(module *types.Module, filter string, sort string, page int, perPage int) (*FindResponse, error) {
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
		Where("(module_id = ? AND deleted_at IS NULL AND json IS NOT NULL)", module.ID)

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
	p = ql.NewParser()
	p.OnIdent = ql.MakeIdentOrderWrapHandler(jsonWrap, "id", "module_id", "user_id", "created_at", "updated_at")

	orderColumns, err := p.ParseColumns(sort)
	if err != nil {
		return nil, err
	}

	var argsOrder = make([]interface{}, 0)
	for _, column := range orderColumns {
		sql, args, err := column.ToSql()
		if err != nil {
			return nil, err
		}
		argsOrder = append(argsOrder, args...)
		query = query.OrderBy(sql)
	}

	// Create actual fetch SQL sentences.
	sqlSelect, argsSelect, err = query.ToSql()
	if err != nil {
		return nil, err
	}

	// Append order args to select args and execute actual query.
	argsSelect = append(argsSelect, argsOrder...)
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

func (r *record) Fields(module *types.Module, record *types.Record) ([]*types.RecordColumn, error) {
	result := make([]*types.RecordColumn, 0)

	if module.ID != record.ModuleID {
		return result, errors.New("Record does not belong to the module")
	}

	fieldNames := module.Fields.Names()

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
