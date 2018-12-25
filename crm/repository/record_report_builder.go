package repository

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	"gopkg.in/Masterminds/squirrel.v1"

	"github.com/crusttech/crust/crm/repository/ql"
)

type (
	recordReportBuilder struct {
		moduleID uint64

		metrics    ql.Columns
		dimensions ql.Columns
		filter     ql.ASTNode

		// This is set by metric/column building to assist Cast()
		numerics []string
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
func stdGroupByFuncHandler(f ql.Function) (ql.Function, error) {
	switch strings.ToUpper(f.Name) {
	case "DATE_FORMAT":
		if len(f.Arguments) == 2 {
			return f, nil
		} else {
			return f, fmt.Errorf("incorrect parameter count for group-by function '%s'", f.Name)
		}
	case "CONCAT", "QUARTER", "YEAR", "DATE", "NOW":
		return f, nil

	default:
		return f, fmt.Errorf("unsupported group-by function %q", f.Name)
	}

}

func NewRecordReportBuilder(moduleID uint64) *recordReportBuilder {
	return &recordReportBuilder{moduleID: moduleID}
}

func (b *recordReportBuilder) SetMetrics(metrics string) (err error) {
	p := ql.NewParser()

	p.OnIdent = ql.MakeIdentWrapHandler(jsonWrap, "created_at", "updated_at")
	p.OnFunction = stdAggregationHandler

	b.metrics, err = p.ParseColumns(metrics)
	return
}

func (b *recordReportBuilder) SetDimensions(dimensions string) (err error) {
	p := ql.NewParser()

	p.OnIdent = ql.MakeIdentWrapHandler(jsonWrap, "created_at", "updated_at")
	p.OnFunction = stdGroupByFuncHandler

	b.dimensions, err = p.ParseColumns(dimensions)
	return
}

func (b *recordReportBuilder) SetFilter(filters string) (err error) {
	p := ql.NewParser()

	p.OnIdent = ql.MakeIdentWrapHandler(jsonWrap, "created_at", "updated_at", "id", "user_id")
	p.OnFunction = stdGroupByFuncHandler

	b.filter, err = p.ParseExpression(filters)
	return
}

func (b *recordReportBuilder) Build() (sql string, args []interface{}, err error) {
	report := squirrel.
		Select().
		Column(squirrel.Alias(squirrel.Expr("COUNT(*)"), "count")).
		From("crm_record").
		Where("module_id = ?", b.moduleID)

	// Add all metrics to columns
	for i, m := range b.metrics {
		if m.Alias == "" {
			// Generate alias
			m.Alias = fmt.Sprintf("metric_%d", i)
		}

		// Wrap to cast func to ensure numeric output
		col := squirrel.Alias(SqlConcatExpr("CAST(", m.Expr, " AS DECIMAL(14,2))"), m.Alias)
		report = report.Column(col)

		b.numerics = append(b.numerics, m.Alias)
	}

	// Add all dimensions to columns
	for i, d := range b.dimensions {
		if d.Alias == "" {
			d.Alias = fmt.Sprintf("dimension_%d", i)
		}

		report = report.Column(d)
		report = report.GroupBy(d.Alias)
		report = report.OrderBy(d.Alias)
	}

	report = report.Where(b.filter)

	return report.ToSql()
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
