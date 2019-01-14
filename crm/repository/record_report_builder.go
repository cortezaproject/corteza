package repository

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"gopkg.in/Masterminds/squirrel.v1"

	"github.com/crusttech/crust/crm/repository/ql"
	"github.com/crusttech/crust/crm/types"
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
		From("crm_record").
		Where("module_id = ?", module.ID)

	return &recordReportBuilder{
		parser: ql.NewParser(),
		module: module,
		report: report,
	}
}

func (b *recordReportBuilder) isRealCol(name string) bool {
	switch name {
	case "id",
		"module_id",
		"user_id",
		"created_at",
		"updated_at":
		return true
	}

	return false
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
		if b.isRealCol(i.Value) {
			return i, nil
		}

		if !b.module.Fields.HasName(i.Value) {
			return i, errors.Errorf("unknown field %q", i.Value)
		}

		if !alreadyJoined(i.Value) {
			b.report = b.report.LeftJoin(fmt.Sprintf(
				"crm_record_value AS rv_%s ON (rv_%s.record_id = crm_record.id AND rv_%s.name = ? AND rv_%s.deleted_at IS NULL)",
				i.Value, i.Value, i.Value, i.Value,
			), i.Value)
		}

		// @todo switch value for ref when doing Record/User lookup
		i.Value = fmt.Sprintf("rv_%s.value", i.Value)

		return i, nil
	}

	var columns ql.Columns
	b.parser.OnFunction = stdAggregationHandler
	if columns, err = b.parser.ParseColumns(metrics); err != nil {
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

	var filter ql.ASTNode
	if filter, err = b.parser.ParseExpression(filters); err != nil {
		return
	}

	b.report = b.report.Where(filter)

	return b.report.ToSql()
}

func (b recordReportBuilder) Cast(row sqlx.ColScanner) map[string]interface{} {
	out := map[string]interface{}{}
	sqlx.MapScan(row, out)
	for k, v := range out {
		switch cv := v.(type) {
		case []uint8:
			out[k] = string(cv)
		}
	}

	// Cast all metrics to float64
	for _, numeric := range b.numerics {
		out[numeric], _ = strconv.ParseFloat(out[numeric].(string), 64)
	}

	return out
}
