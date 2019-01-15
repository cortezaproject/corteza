package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lann/builder"
	"github.com/pkg/errors"
	"github.com/titpetric/factory"
	sq "gopkg.in/Masterminds/squirrel.v1"

	"github.com/crusttech/crust/crm/repository/ql"

	"github.com/crusttech/crust/crm/types"
)

type (
	RecordRepository interface {
		With(ctx context.Context, db *factory.DB) RecordRepository

		FindByID(id uint64) (*types.Record, error)

		Report(module *types.Module, metrics, dimensions, filter string) (results interface{}, err error)
		Find(module *types.Module, filter string, sort string, page int, perPage int) (*FindResponse, error)

		Create(record *types.Record) (*types.Record, error)
		Update(record *types.Record) (*types.Record, error)
		DeleteByID(id uint64) error

		UpdateValues(recordID uint64, rvs types.RecordValueSet) (err error)
		LoadValues(IDs ...uint64) (rvs types.RecordValueSet, err error)
	}

	FindResponseMeta struct {
		Filter string `json:"filter,omitempty"`
		Sort   string `json:"sort,omitempty"`

		Page    int `json:"page"`
		PerPage int `json:"perPage"`
		Count   int `json:"count"`
	}

	FindResponse struct {
		Meta    FindResponseMeta `json:"meta"`
		Records types.RecordSet  `json:"records"`
	}

	record struct {
		*repository
	}
)

const (
	sortWrap = `sort`
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

func (r *record) Report(module *types.Module, metrics, dimensions, filter string) (results interface{}, err error) {
	crb := NewRecordReportBuilder(module)

	var result = make([]map[string]interface{}, 0)

	if query, args, err := crb.Build(metrics, dimensions, filter); err != nil {
		return nil, errors.Wrap(err, "can not generate report query")
	} else if rows, err := r.db().Query(query, args...); err != nil {
		return nil, errors.Wrapf(err, "can not execute report query (%s)", query)
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

	var query, err = r.buildQuery(module, filter, sort)
	if err != nil {
		return nil, err
	}

	// Assemble SQL for counting (includes only where)
	count := query.Column("COUNT(*)")
	count = builder.Delete(count, "OrderBys").(sq.SelectBuilder)
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

	// Assemble SQL for fetching record (where + sorting + paging)...
	query = query.
		Column("crm_record.*").
		Limit(uint64(perPage)).
		Offset(uint64(page))

	// Create actual fetch SQL sentences.
	sqlSelect, argsSelect, err = query.ToSql()
	if err != nil {
		return nil, err
	}

	// Append order args to select args and execute actual query.
	if err := r.db().Select(&response.Records, sqlSelect, argsSelect...); err != nil {
		return nil, err
	}

	return response, nil
}

func (r *record) buildQuery(module *types.Module, filter string, sort string) (query sq.SelectBuilder, err error) {
	// Create query for fetching and counting records.
	query = sq.Select().
		From("crm_record").
		Where(sq.Eq{"module_id": module.ID}).
		Where(sq.Eq{"deleted_at": nil})

	// Do not translate/wrap these
	var realColumns = []string{
		"id",
		"module_id",
		"user_id",
		"created_at",
		"updated_at",
	}

	const colWrap = `(SELECT value FROM crm_record_value WHERE name = ? AND record_id = crm_record.id AND deleted_at IS NULL)`

	// Parse filters.
	if filter != "" {
		var (
			// Filter parser
			fp = ql.NewParser()

			// Filter node
			fn ql.ASTNode
		)

		// Make a nice wrapper that will translate module fields to subqueries
		fp.OnIdent = func(i ql.Ident) (ql.Ident, error) {
			for _, s := range realColumns {
				if s == i.Value {
					return i, nil
				}
			}

			if !module.Fields.HasName(i.Value) {
				return i, errors.Errorf("unknown field %q", i.Value)
			}

			// @todo switch value for ref when doing Record/User lookup

			i.Args = []interface{}{i.Value}
			i.Value = colWrap

			return i, nil
		}

		if fn, err = fp.ParseExpression(filter); err != nil {
			return
		}

		query = query.Where(fn)
	}

	if sort != "" {
		var (
			// Sort parser
			sp = ql.NewParser()

			// Sort columns
			sc ql.Columns
		)

		sp.OnIdent = func(i ql.Ident) (ql.Ident, error) {
			for _, s := range realColumns {
				if s == i.Value {
					i.Value += " "
					return i, nil
				}
			}

			if !module.Fields.HasName(i.Value) {
				return i, errors.Errorf("unknown field %q", i.Value)
			}

			i.Value = strings.Replace(colWrap, "?", fmt.Sprintf("'%s'", i.Value), 1) + " "
			return i, nil
		}

		if sc, err = sp.ParseColumns(sort); err != nil {
			return
		}

		query = query.OrderBy(sc.Strings()...)
	}

	return
}

func (r *record) Create(record *types.Record) (*types.Record, error) {
	record.ID = factory.Sonyflake.NextID()
	record.CreatedAt = time.Now()
	record.UserID = Identity(r.Context())

	if err := r.db().Replace("crm_record", record); err != nil {
		return nil, errors.Wrap(err, "could not update record")
	}

	return record, nil
}

func (r *record) Update(record *types.Record) (*types.Record, error) {
	now := time.Now()
	record.UpdatedAt = &now

	if err := r.db().Replace("crm_record", record); err != nil {
		return nil, errors.Wrap(err, "could not update record")
	}

	return record, nil
}

func (r *record) DeleteByID(id uint64) error {
	_, err := r.db().Exec("update crm_record set deleted_at=? where id=?", time.Now(), id)
	return err
}

func (r *record) UpdateValues(recordID uint64, rvs types.RecordValueSet) (err error) {
	// Remove all records and prepare to be updated
	// @todo be more selective and delete only removed values
	if _, err = r.db().Exec("DELETE FROM crm_record_value WHERE record_id = ?", recordID); err != nil {
		return errors.Wrap(err, "could not remove record values")
	}

	err = rvs.Walk(func(value *types.RecordValue) error {
		value.RecordID = recordID
		return r.db().Replace("crm_record_value", value)
	})

	return errors.Wrap(err, "could not replace record values")

}

func (r *record) LoadValues(IDs ...uint64) (rvs types.RecordValueSet, err error) {
	if len(IDs) == 0 {
		return
	}

	var sql = "SELECT * FROM crm_record_value WHERE record_id IN (?) AND deleted_at IS NULL ORDER BY record_id, place"

	if sql, args, err := sqlx.In(sql, IDs); err != nil {
		return nil, err
	} else {
		return rvs, r.db().Select(&rvs, sql, args...)
	}
}
