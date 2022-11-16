package rdbms

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/ql"
	"github.com/cortezaproject/corteza/server/pkg/slice"
	"github.com/jmoiron/sqlx"
)

type (
	ComposeRecordsReport []map[string]interface{}

	recordReportBuilder struct {
		// This is set by metric/column building to assist Cast()
		numerics []string

		parser *ql.Parser

		module     *types.Module
		metrics    string
		dimensions string
		filters    string

		store recordReportBuilderStoreQuerier

		supportedAggregationFunctions map[string]bool
		supportedFilterFunctions      map[string]bool
	}

	recordReportBuilderStoreQuerier interface {
		SelectBuilder(string, ...string) squirrel.SelectBuilder
		Query(context.Context, squirrel.Sqlizer) (*sql.Rows, error)
		SqlFunctionHandler(f ql.Function) (ql.ASTNode, error)
		FieldToColumnTypeCaster(f ModuleFieldTypeDetector, i ql.Ident) (ql.Ident, error)
	}
)

func ComposeRecordReportBuilder(s *Store, module *types.Module, metrics, dimensions, filters string) *recordReportBuilder {
	return &recordReportBuilder{
		parser:     ql.NewParser(),
		store:      s,
		module:     module,
		metrics:    metrics,
		dimensions: dimensions,
		filters:    filters,

		supportedAggregationFunctions: slice.ToStringBoolMap([]string{
			"COUNT",
			"SUM",
			"MAX",
			"MIN",
			"AVG",
			"STD",
		}),

		supportedFilterFunctions: slice.ToStringBoolMap([]string{
			"CONCAT",
			"QUARTER",
			"YEAR",
			"DATE",
			"NOW",
			"DATE_ADD",
			"DATE_SUB",
			"DATE_FORMAT",
		}),
	}
}

func (b *recordReportBuilder) Run(ctx context.Context) (ComposeRecordsReport, error) {
	var (
		result = make(ComposeRecordsReport, 0)
	)

	if sb, err := b.Build(); err != nil {
		return nil, fmt.Errorf("cannot generate report query: %w", err)
	} else if rows, err := b.store.Query(ctx, sb); err != nil {
		return nil, fmt.Errorf("cannot execute report query: %w", err)
	} else {
		err = func() error {
			defer rows.Close()
			for rows.Next() {
				r, err := b.Cast(rows)
				if err != nil {
					return err
				}
				result = append(result, r)
			}

			return nil
		}()

		if err = rows.Err(); err != nil {
			return nil, err
		}

		return result, nil
	}
}

func (b recordReportBuilder) ToSql() (string, []interface{}, error) {
	if sb, err := b.Build(); err != nil {
		return "", nil, err
	} else {
		return sb.ToSql()
	}
}

// Identifiers should be names of the fields (physical table columns OR json fields, defined in module)
func (b *recordReportBuilder) stdAggregationHandler(f ql.Function) (ql.ASTNode, error) {
	if !b.supportedAggregationFunctions[strings.ToUpper(f.Name)] {
		return f, fmt.Errorf("unsupported aggregate function %q", f.Name)
	}

	return b.store.SqlFunctionHandler(f)
}

// Identifiers should be names of the fields (physical table columns OR json fields, defined in module)
func (b *recordReportBuilder) stdFilterFuncHandler(f ql.Function) (ql.ASTNode, error) {
	if !b.supportedFilterFunctions[strings.ToUpper(f.Name)] {
		return f, fmt.Errorf("unsupported filter function %q", f.Name)
	}

	return b.store.SqlFunctionHandler(f)
}

func (b *recordReportBuilder) Build() (sb squirrel.SelectBuilder, err error) {
	var (
		joinTpl = "compose_record_value AS rv_%s ON (rv_%s.record_id = crd.id AND rv_%s.name = '%s' AND rv_%s.deleted_at IS NULL)"

		report = b.store.SelectBuilder("compose_record AS crd").
			// record rows can duplicate due to multi value fields; we need to remove duplicates when counting
			Column(squirrel.Alias(squirrel.Expr("COUNT(DISTINCT(crd.id))"), "count")).
			Where("crd.deleted_at IS NULL").
			Where("crd.module_id = ?", b.module.ID)

		joinedFields = []string{}

		alreadyJoined = func(f string) bool {
			for _, a := range joinedFields {
				if a == f {
					return true
				}
			}

			joinedFields = append(joinedFields, f)
			return false
		}
	)

	b.parser.OnIdent = func(i ql.Ident) (ql.Ident, error) {
		var is bool
		if i.Value, _, is = isRealRecordCol(i.Value); is {
			return i, nil
		}

		if !b.module.Fields.HasName(i.Value) {
			return i, fmt.Errorf("unknown field %q", i.Value)
		}

		if !alreadyJoined(i.Value) {
			report = report.LeftJoin(strings.ReplaceAll(joinTpl, "%s", i.Value))
		}

		return b.store.FieldToColumnTypeCaster(b.module.Fields.FindByName(i.Value), i)
	}

	var columns ql.Columns
	b.parser.OnFunction = b.stdAggregationHandler
	if columns, err = b.parser.ParseColumns(b.metrics); err != nil {
		err = fmt.Errorf("could not parse metrics %q: %w", b.metrics, err)
		return
	}

	// Add all metrics to columns
	for i, m := range columns {
		if m.Alias == "" {
			// Generate alias
			m.Alias = fmt.Sprintf("metric_%d", i)
		}

		// Wrap to cast func to ensure numeric output
		col := squirrel.Alias(SquirrelConcatExpr("CAST(", m.Expr, " AS DECIMAL(14,2))"), m.Alias)
		report = report.Column(col)

		b.numerics = append(b.numerics, m.Alias)
	}

	b.parser.OnFunction = b.stdFilterFuncHandler
	if columns, err = b.parser.ParseColumns(b.dimensions); err != nil {
		err = fmt.Errorf("could not parse dimensions %q: %w", b.dimensions, err)
		return
	}

	// Add dimensions
	for i, d := range columns {
		if d.Alias == "" {
			d.Alias = fmt.Sprintf("dimension_%d", i)
		}

		report = report.
			Column(d).
			GroupBy(d.Alias).
			OrderBy(d.Alias)
	}

	// Use a different handler for filter functions for this
	b.parser.OnFunction = b.stdFilterFuncHandler

	if len(b.filters) > 0 {
		var filter ql.ASTNode
		if filter, err = b.parser.ParseExpression(b.filters); err != nil {
			err = fmt.Errorf("could not parse filters %q: %w", b.filters, err)
			return
		}

		// We need to wrap this one level deeper, since additional filters should
		// be evaluated as a whole.
		// For example A AND B OR C =should be> (A AND B OR C)
		// so the output becomes BASE AND (ADDITIONAL)
		report = report.Where(ql.ASTNodes{filter})
	}

	return report, nil
}

func (b recordReportBuilder) Cast(row sqlx.ColScanner) (out map[string]interface{}, err error) {
	out = map[string]interface{}{}
	if err = sqlx.MapScan(row, out); err != nil {
		return nil, err
	}

	for k, v := range out {
		switch cv := v.(type) {
		case []uint8:
			out[k] = string(cv)
		case uint64, int64:
			// Just to make sure we don't break anything old
			if strings.Contains(k, "dimension") {
				out[k] = fmt.Sprintf("%d", cv)
			} else {
				out[k] = cv
			}
		default:
		}
	}

	// Cast all metrics to float64
	for _, fname := range b.numerics {
		switch num := out[fname].(type) {
		case string:
			out[fname], err = strconv.ParseFloat(num, 64)
			if err != nil {
				return nil, err
			}
		}
	}

	return
}
