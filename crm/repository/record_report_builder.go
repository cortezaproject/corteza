package repository

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	"gopkg.in/Masterminds/squirrel.v1"

	"github.com/crusttech/crust/crm/types"
)

type (
	recordReportBuilder struct {
		jsonField string

		moduleID uint64
		params   *types.RecordReport
	}
)

var (
	recordReportExprMatch = regexp.MustCompile(`^\s*(\w+)\((.+)\)\s*$`)
)

func NewRecordReportBuilder(moduleID uint64, params *types.RecordReport) *recordReportBuilder {
	return &recordReportBuilder{
		moduleID:  moduleID,
		params:    params,
		jsonField: `JSON_UNQUOTE(JSON_EXTRACT(json, REPLACE(JSON_UNQUOTE(JSON_SEARCH(json, 'one', ?)), '.name', '.value')))`,
	}
}

func (b recordReportBuilder) field(name string) squirrel.Sqlizer {
	switch name {
	case "created_at", "updated_at":
		return squirrel.Expr(name)
	default:
		return squirrel.Expr(b.jsonField, name)
	}
}

func (b recordReportBuilder) alias(col squirrel.Sqlizer, alias, fallback string) (squirrel.Sqlizer, string) {
	if alias != "" {
		return squirrel.Alias(col, alias), alias
	}

	return squirrel.Alias(col, fallback), fallback
}

func (b recordReportBuilder) wrapInModifiers(col squirrel.Sqlizer, mm ...string) squirrel.Sqlizer {
	for _, m := range mm {
		switch strings.ToUpper(m) {
		case "WEEKDAY":
			col = SqlConcatExpr("DATE_FORMAT(", col, ", '%W')")
		case "DATE":
			col = SqlConcatExpr("DATE_FORMAT(", col, ", '%Y-%m-%d')")
		case "WEEK":
			col = SqlConcatExpr("DATE_FORMAT(", col, ", '%Y-%u')")
		case "MONTH":
			col = SqlConcatExpr("DATE_FORMAT(", col, ", '%Y-%m')")
		case "QUARTER":
			col = SqlConcatExpr("CONCAT(", "YEAR(", col, "), 'Q', ", "QUARTER(", col, ")", ")")
		case "YEAR":
			col = SqlConcatExpr("DATE_FORMAT(", col, ", '%Y')")
		}
	}

	return col
}

func (b recordReportBuilder) parseExpression(exp string) squirrel.Sqlizer {
	res := recordReportExprMatch.FindStringSubmatch(exp)
	if len(res) > 0 {
		aggrFuncName := strings.ToUpper(res[1])
		aggrFuncArgs := b.parseExpression(res[2])

		switch aggrFuncName {
		case "COUNTD":
			return SqlConcatExpr("COUNT(DISTINCT ", aggrFuncArgs, ")")

		case "SUM", "MAX", "MIN", "AVG", "STD":
			return SqlConcatExpr(aggrFuncName+"(CAST(", aggrFuncArgs, " AS DECIMAL(14,2)))")
		}
	} else {
		return b.field(exp)
	}

	return nil
}

func (b *recordReportBuilder) Build() (sql string, args []interface{}, err error) {
	report := squirrel.
		Select().
		Column(squirrel.Alias(squirrel.Expr("COUNT(*)"), "count")).
		From("crm_record").
		Where("module_id = ?", b.moduleID)

	if b.params == nil {
		return "", nil, errors.New("can not generate report without parameters")
	}

	for i, m := range b.params.Metrics {
		col := SqlConcatExpr("CAST(", b.parseExpression(m.Expression), " AS DECIMAL(14,2))")
		col, m.Alias = b.alias(col, m.Alias, fmt.Sprintf("metric_%d", i))

		report = report.Column(col)

		b.params.Metrics[i].Alias = m.Alias // copy generated alias back
	}

	for i, d := range b.params.Dimensions {
		col := b.field(d.Field)

		col = b.wrapInModifiers(col, d.Modifiers...)

		col, d.Alias = b.alias(col, d.Alias, fmt.Sprintf("dimension_%d", i))

		report = report.Column(col)
		report = report.GroupBy(d.Alias)
		report = report.OrderBy(d.Alias)

		b.params.Dimensions[i].Alias = d.Alias // copy generated alias back
	}

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
	for _, m := range b.params.Metrics {
		out[m.Alias], _ = strconv.ParseFloat(out[m.Alias].(string), 64)
	}

	return out
}
