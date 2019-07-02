package repository

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"gopkg.in/Masterminds/squirrel.v1"

	"github.com/cortezaproject/corteza-server/compose/internal/repository/ql"
	"github.com/cortezaproject/corteza-server/compose/types"
)

type (
	recordReportBuilder struct {
		module *types.Module

		// This is set by metric/column building to assist Cast()
		numerics []string

		report squirrel.SelectBuilder
		parser *ql.Parser
	}
)

// Identifiers should be names of the fields (physical table columns OR json fields, defined in module)
func stdAggregationHandler(f ql.Function) (ql.Function, error) {
	switch strings.ToUpper(f.Name) {
	// case "COUNTD":
	// 	return SqlConcatExpr("COUNT(DISTINCT ", aggrFuncArgs, ")")
	//
	case "COUNT", "SUM", "MAX", "MIN", "AVG", "STD":
		return f, nil
	default:
		return f, fmt.Errorf("unsupported aggregate function %q", f.Name)
	}
}

// Identifiers should be names of the fields (physical table columns OR json fields, defined in module)
func stdFilterFuncHandler(f ql.Function) (ql.Function, error) {
	switch strings.ToUpper(f.Name) {
	case "CONCAT", "QUARTER", "YEAR", "DATE", "NOW", "DATE_ADD", "DATE_SUB", "DATE_FORMAT":
		return f, nil

	default:
		return f, fmt.Errorf("unsupported group-by function %q", f.Name)
	}
}

func NewRecordReportBuilder(module *types.Module) *recordReportBuilder {
	var report = squirrel.
		Select().
		Column(squirrel.Alias(squirrel.Expr("COUNT(*)"), "count")).
		From("compose_record AS r").
		Where("r.deleted_at IS NULL").
		Where("r.module_id = ?", module.ID)

	return &recordReportBuilder{
		parser: ql.NewParser(),
		module: module,
		report: report,
	}
}

func (b *recordReportBuilder) Build(metrics, dimensions, filters string) (sql string, args []interface{}, err error) {
	var joinedFields = []string{}
	var alreadyJoined = func(f string) bool {
		for _, a := range joinedFields {
			if a == f {
				return true
			}
		}

		joinedFields = append(joinedFields, f)
		return false
	}

	b.parser.OnIdent = func(i ql.Ident) (ql.Ident, error) {
		var is bool
		if i.Value, is = isRealRecordCol(i.Value); is {
			return i, nil
		}

		if !b.module.Fields.HasName(i.Value) {
			return i, errors.Errorf("unknown field %q", i.Value)
		}

		if !alreadyJoined(i.Value) {
			b.report = b.report.LeftJoin(fmt.Sprintf(
				"compose_record_value AS rv_%s ON (rv_%s.record_id = r.id AND rv_%s.name = ? AND rv_%s.deleted_at IS NULL)",
				i.Value, i.Value, i.Value, i.Value,
			), i.Value)
		}

		// @todo switch value for ref when doing Record/Owner lookup
		i.Value = fmt.Sprintf("rv_%s.value", i.Value)

		return i, nil
	}

	var columns ql.Columns
	b.parser.OnFunction = stdAggregationHandler
	if columns, err = b.parser.ParseColumns(metrics); err != nil {
		err = errors.Wrapf(err, "could not parse metrics %q", metrics)
		return
	}

	// Add all metrics to columns
	for i, m := range columns {
		if m.Alias == "" {
			// Generate alias
			m.Alias = fmt.Sprintf("metric_%d", i)
		}

		// Wrap to cast func to ensure numeric output
		col := squirrel.Alias(SqlConcatExpr("CAST(", m.Expr, " AS DECIMAL(14,2))"), m.Alias)
		b.report = b.report.Column(col)

		b.numerics = append(b.numerics, m.Alias)
	}

	b.parser.OnFunction = stdFilterFuncHandler
	if columns, err = b.parser.ParseColumns(dimensions); err != nil {
		err = errors.Wrapf(err, "could not parse dimensions %q", dimensions)
		return
	}

	// Add dimensions
	for i, d := range columns {
		if d.Alias == "" {
			d.Alias = fmt.Sprintf("dimension_%d", i)
		}

		b.report = b.report.
			Column(d).
			GroupBy(d.Alias).
			OrderBy(d.Alias)
	}

	// Use a different handler for filter functions for this
	b.parser.OnFunction = stdFilterFuncHandler

	if len(filters) > 0 {
		var filter ql.ASTNode
		if filter, err = b.parser.ParseExpression(filters); err != nil {
			err = errors.Wrapf(err, "could not parse filters %q", filters)
			return
		}

		b.report = b.report.Where(filter)
	}

	return b.report.ToSql()
}

func (b recordReportBuilder) Cast(row sqlx.ColScanner) map[string]interface{} {
	out := map[string]interface{}{}
	sqlx.MapScan(row, out)
	for k, v := range out {
		switch cv := v.(type) {
		case []uint8:
			out[k] = string(cv)
		default:
		}
	}

	// Cast all metrics to float64
	for _, fname := range b.numerics {
		switch num := out[fname].(type) {
		case string:
			out[fname], _ = strconv.ParseFloat(num, 64)
		}
	}

	return out
}
