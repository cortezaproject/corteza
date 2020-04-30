package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/titpetric/factory"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/ql"
	"github.com/cortezaproject/corteza-server/pkg/rh"
)

type (
	RecordRepository interface {
		With(ctx context.Context, db *factory.DB) RecordRepository

		FindByID(namespaceID, recordID uint64) (*types.Record, error)

		Report(module *types.Module, metrics, dimensions, filter string) (results interface{}, err error)
		Find(module *types.Module, filter types.RecordFilter) (set types.RecordSet, f types.RecordFilter, err error)
		Export(module *types.Module, filter types.RecordFilter) (set types.RecordSet, err error)

		Create(record *types.Record) (*types.Record, error)
		Update(record *types.Record) (*types.Record, error)
		Delete(record *types.Record) error

		RefValueLookup(moduleID uint64, field string, ref uint64) (recordID uint64, err error)
		LoadValues(fieldNames []string, IDs []uint64) (rvs types.RecordValueSet, err error)
		DeleteValues(record *types.Record) error
		UpdateValues(recordID uint64, rvs types.RecordValueSet) (err error)
		PartialUpdateValues(rvs ...*types.RecordValue) (err error)
	}

	record struct {
		*repository
	}
)

const (
	ErrRecordNotFound = repositoryError("RecordNotFound")
)

func Record(ctx context.Context, db *factory.DB) RecordRepository {
	return (&record{}).With(ctx, db)
}

func (r record) With(ctx context.Context, db *factory.DB) RecordRepository {
	return &record{
		repository: r.repository.With(ctx, db),
	}
}

func (r record) table() string {
	return "compose_record"
}

func (r record) columns() []string {
	return []string{
		"r.id",
		"r.module_id",
		"r.rel_namespace",
		"r.owned_by",
		"r.created_at",
		"r.created_by",
		"r.updated_at",
		"r.updated_by",
		"r.deleted_at",
		"r.deleted_by",
	}
}

func (r record) query() squirrel.SelectBuilder {
	return squirrel.
		Select(r.columns()...).
		From(r.table() + " AS r")
}

// @todo: update to accepted DeletedAt column semantics from Messaging

func (r record) FindByID(namespaceID, recordID uint64) (*types.Record, error) {
	return r.findOneBy(namespaceID, "id", recordID)
}

func (r record) findOneBy(namespaceID uint64, field string, value interface{}) (*types.Record, error) {
	var (
		rec = &types.Record{}

		q = r.query().
			Where("r.deleted_at IS NULL").
			Where(squirrel.Eq{field: value, "rel_namespace": namespaceID})

		err = rh.FetchOne(r.db(), q, rec)
	)

	if err != nil {
		return nil, err
	} else if rec.ID == 0 {
		return nil, ErrRecordNotFound
	}

	return rec, nil
}

func (r record) Report(module *types.Module, metrics, dimensions, filter string) (results interface{}, err error) {
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

func (r record) Find(module *types.Module, filter types.RecordFilter) (set types.RecordSet, f types.RecordFilter, err error) {
	var query squirrel.SelectBuilder
	f = filter

	query, err = r.buildQuery(module, filter)
	if err != nil {
		return
	}

	if f.Count, err = rh.Count(r.db(), query); err != nil || f.Count == 0 {
		return
	}

	return set, f, rh.FetchPaged(r.db(), query, f.PageFilter, &set)
}

// Export ignores paging and does not return filter
//
// @todo optimize and include value loading
func (r record) Export(module *types.Module, filter types.RecordFilter) (set types.RecordSet, err error) {
	filter.PageFilter = rh.PageFilter{}

	query, err := r.buildQuery(module, filter)
	if err != nil {
		return
	}

	return set, rh.FetchAll(r.db(), query, &set)
}

func (r record) buildQuery(module *types.Module, f types.RecordFilter) (query squirrel.SelectBuilder, err error) {
	var (
		joinedFields  = []string{}
		alreadyJoined = func(f string) bool {
			for _, a := range joinedFields {
				if a == f {
					return true
				}
			}

			joinedFields = append(joinedFields, f)
			return false
		}

		identResolver = func(i ql.Ident) (ql.Ident, error) {
			var is bool
			if i.Value, is = isRealRecordCol(i.Value); is {
				i.Value += " "
				return i, nil
			}

			if !module.Fields.HasName(i.Value) {
				return i, errors.Errorf("unknown field %q", i.Value)
			}

			if !alreadyJoined(i.Value) {
				query = query.LeftJoin(fmt.Sprintf(
					"compose_record_value AS rv_%s ON (rv_%s.record_id = r.id AND rv_%s.name = ? AND rv_%s.deleted_at IS NULL)",
					i.Value, i.Value, i.Value, i.Value,
				), i.Value)
			}

			field := module.Fields.FindByName(i.Value)

			switch true {
			case field.IsBoolean():
				i.Value = fmt.Sprintf("(rv_%s.value NOT IN ('', '0', 'false', 'f',  'FALSE', 'F', false))", i.Value)
			case field.IsNumeric():
				i.Value = fmt.Sprintf("CAST(rv_%s.value AS SIGNED)", i.Value)
			case field.IsDateTime():
				i.Value = fmt.Sprintf("CAST(rv_%s.value AS DATETIME)", i.Value)
			case field.IsRef():
				i.Value = fmt.Sprintf("rv_%s.ref ", i.Value)
			default:
				i.Value = fmt.Sprintf("rv_%s.value ", i.Value)
			}

			return i, nil
		}
	)

	// Create query for fetching and counting records.
	query = r.query().
		Where("r.module_id = ?", module.ID).
		Where("r.rel_namespace = ?", module.NamespaceID)

	// Inc/exclude deleted records according to filter settings
	query = rh.FilterNullByState(query, "r.deleted_at", f.Deleted)

	// Parse filters.
	if f.Query != "" {
		var (
			// Filter parser
			fp = ql.NewParser()

			// Filter node
			fn ql.ASTNode
		)

		// Resolve all identifiers found in the query
		// into their table/column counterparts
		fp.OnIdent = identResolver

		if fn, err = fp.ParseExpression(f.Query); err != nil {
			return
		} else if filterSql, filterArgs, err := fn.ToSql(); err != nil {
			return query, err
		} else {
			query = query.Where("("+filterSql+")", filterArgs...)
		}
	}

	if f.Sort != "" {
		var (
			// Sort parser
			sp = ql.NewParser()

			// Sort columns
			sc ql.Columns
		)

		// Resolve all identifiers found in sort
		// into their table/column counterparts
		sp.OnIdent = identResolver

		if sc, err = sp.ParseColumns(f.Sort); err != nil {
			return
		}

		query = query.OrderBy(sc.Strings()...)
	}

	return
}

func (r record) Create(record *types.Record) (*types.Record, error) {
	record.ID = factory.Sonyflake.NextID()

	if err := r.db().Insert("compose_record", record); err != nil {
		return nil, errors.Wrap(err, "could not update record")
	}

	return record, nil
}

func (r record) Update(record *types.Record) (*types.Record, error) {
	if err := r.db().Update("compose_record", record, "id"); err != nil {
		return nil, errors.Wrap(err, "could not update record")
	}

	return record, nil
}

func (r record) Delete(record *types.Record) error {
	_, err := r.db().Exec(
		"UPDATE compose_record SET deleted_at = ?, deleted_by = ? WHERE rel_namespace = ? AND id = ?",
		record.DeletedAt,
		record.DeletedBy,
		record.NamespaceID,
		record.ID,
	)

	return err
}

func (r record) DeleteValues(record *types.Record) error {
	_, err := r.db().Exec(
		"UPDATE compose_record_value SET deleted_at = ? WHERE record_id = ?",
		record.DeletedAt,
		record.ID)

	return err
}

func (r record) UpdateValues(recordID uint64, rvs types.RecordValueSet) (err error) {
	// Remove all records and prepare to be updated
	// @todo be more selective and delete only removed and update/insert changed/new values
	if _, err = r.db().Exec("DELETE FROM compose_record_value WHERE record_id = ?", recordID); err != nil {
		return errors.Wrap(err, "could not remove record values")
	}

	err = rvs.Walk(func(value *types.RecordValue) error {
		if value.DeletedAt != nil {
			return nil
		}

		value.RecordID = recordID
		return r.db().Insert("compose_record_value", value)
	})

	return errors.Wrap(err, "could not insert record values")

}

func (r record) PartialUpdateValues(rvs ...*types.RecordValue) (err error) {
	err = types.RecordValueSet(rvs).Walk(func(value *types.RecordValue) error {
		return r.db().Replace("compose_record_value", value)
	})

	return errors.Wrap(err, "could not replace record values")
}

func (r record) RefValueLookup(moduleID uint64, field string, ref uint64) (recordID uint64, err error) {
	var sql = "SELECT record_id" +
		"  FROM compose_record AS r INNER JOIN compose_record_value AS v " +
		" WHERE r.module_id = ? " +
		"   AND v.name = ? " +
		"   AND v.ref = ? " +
		"   AND r.deleted_at IS NULL " +
		"   AND v.deleted_at IS NULL " +
		"       LIMIT 1"

	return recordID, r.db().Get(&recordID, sql, moduleID, field, ref)
}

func (r record) LoadValues(fieldNames []string, IDs []uint64) (rvs types.RecordValueSet, err error) {
	if len(fieldNames) == 0 || len(IDs) == 0 {
		return
	}

	var sql = "SELECT record_id, name, value, ref, place, deleted_at " +
		"  FROM compose_record_value " +
		" WHERE record_id IN (?) " +
		"   AND name IN (?) " +
		"   AND deleted_at IS NULL " +
		" ORDER BY record_id, place"

	if sql, args, err := sqlx.In(sql, IDs, fieldNames); err != nil {
		return nil, err
	} else {
		return rvs, r.db().Select(&rvs, sql, args...)
	}
}

// Checks if field name is "real column", reformats it and returns
func isRealRecordCol(name string) (string, bool) {
	switch name {
	case
		"recordID",
		"id":
		return "r.id", true
	case
		"module_id",
		"owned_by",
		"created_by",
		"created_at",
		"updated_by",
		"updated_at",
		"deleted_by",
		"deleted_at":
		return "r." + name, true

	case
		"moduleID",
		"ownedBy",
		"createdBy",
		"createdAt",
		"updatedBy",
		"updatedAt",
		"deletedBy",
		"deletedAt":
		return "r." + name[0:len(name)-2] + "_" + strings.ToLower(name[len(name)-2:]), true
	}

	return name, false
}
