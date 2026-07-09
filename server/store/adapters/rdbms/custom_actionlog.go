package rdbms

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/crusttech/human/server/pkg/actionlog"
	"github.com/crusttech/human/server/store/adapters/rdbms/ql"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
)

// Dimension/metric keys are resolved through the pkg/actionlog registries;
// this adapter only maps registry kinds to portable goqu expressions.
func (s *Store) ActionlogReport(ctx context.Context, rr actionlog.ReportRequest) (_ actionlog.ReportRowSet, err error) {
	var (
		dims = make([]actionlog.ReportDimension, 0, len(rr.Dimensions))
		mets = make([]actionlog.ReportMetric, 0, len(rr.Metrics))

		selects []any
		groups  []any
		orders  []exp.OrderedExpression
	)

	for _, key := range rr.Dimensions {
		d, ok := actionlog.ReportDimensionByKey(key)
		if !ok {
			return nil, fmt.Errorf("unknown report dimension %q", key)
		}

		dims = append(dims, d)

		e, err := s.actionlogDimensionExpr(d)
		if err != nil {
			return nil, err
		}

		selects = append(selects, exp.NewAliasExpression(e, d.Key))
		groups = append(groups, e)
		orders = append(orders, exp.NewOrderedExpression(e, exp.AscDir, exp.NoNullsSortType))
	}

	for _, key := range rr.Metrics {
		m, ok := actionlog.ReportMetricByKey(key)
		if !ok {
			return nil, fmt.Errorf("unknown report metric %q", key)
		}

		mets = append(mets, m)

		e, err := actionlogMetricExpr(m)
		if err != nil {
			return nil, err
		}

		selects = append(selects, exp.NewAliasExpression(e, m.Key))
	}

	where, f, err := s.Filters.Actionlog(s, rr.Filter)
	if err != nil {
		return nil, err
	}

	query := s.Dialect.GOQU().
		From(actionlogTable).
		Select(selects...).
		Where(where...).
		GroupBy(groups...).
		Order(orders...).
		Limit(uint(f.Limit))

	rows, err := s.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var (
		set = make(actionlog.ReportRowSet, 0, 64)
		raw = make([]any, len(dims)+len(mets))
		ptr = make([]any, len(raw))
	)

	for i := range raw {
		ptr[i] = &raw[i]
	}

	for rows.Next() {
		if err = rows.Scan(ptr...); err != nil {
			return nil, err
		}

		row := &actionlog.ReportRow{
			Dimensions: make(map[string]any, len(dims)),
			Metrics:    make(map[string]float64, len(mets)),
		}

		for i, d := range dims {
			row.Dimensions[d.Key] = normalizeReportDimensionValue(d, raw[i])
		}

		for i, m := range mets {
			if row.Metrics[m.Key], err = normalizeReportMetricValue(raw[len(dims)+i]); err != nil {
				return nil, fmt.Errorf("metric %q: %w", m.Key, err)
			}
		}

		set = append(set, row)
	}

	return set, rows.Err()
}

func (s *Store) actionlogDimensionExpr(d actionlog.ReportDimension) (exp.Expression, error) {
	switch d.Kind {
	case actionlog.ReportDimensionColumn, actionlog.ReportDimensionRef:
		return goqu.C(d.Column), nil

	case actionlog.ReportDimensionDate:
		// ql "date" ref resolves to the dialect's native date truncation
		// (pg ::DATE, mysql/sqlite DATE(), mssql override)
		return s.Dialect.ExprHandler(&ql.ASTNode{Ref: "date"}, goqu.C(d.Column))
	}

	return nil, fmt.Errorf("unsupported report dimension kind %q", d.Kind)
}

func actionlogMetricExpr(m actionlog.ReportMetric) (exp.Expression, error) {
	switch m.Kind {
	case actionlog.ReportMetricCount:
		return goqu.COUNT(goqu.Star()), nil

	case actionlog.ReportMetricCountDistinct:
		return goqu.COUNT(goqu.DISTINCT(m.Column)), nil

	case actionlog.ReportMetricCountNotEmpty:
		return goqu.SUM(goqu.Case().
			When(goqu.C(m.Column).Neq(""), goqu.L("1")).
			Else(goqu.L("0")),
		), nil
	}

	return nil, fmt.Errorf("unsupported report metric kind %q", m.Kind)
}

// Drivers return mixed types for computed columns (int64, []byte, time.Time);
// normalize so the payload is stable across all supported databases.
func normalizeReportDimensionValue(d actionlog.ReportDimension, v any) any {
	if b, is := v.([]byte); is {
		v = string(b)
	}

	switch d.Kind {
	case actionlog.ReportDimensionRef:
		switch id := v.(type) {
		case int64:
			return strconv.FormatUint(uint64(id), 10)
		case uint64:
			return strconv.FormatUint(id, 10)
		}

	case actionlog.ReportDimensionDate:
		const format = "2006-01-02"
		switch ts := v.(type) {
		case time.Time:
			return ts.Format(format)
		case string:
			if len(ts) > len(format) {
				return ts[:len(format)]
			}
		}
	}

	return v
}

func normalizeReportMetricValue(v any) (float64, error) {
	switch n := v.(type) {
	case nil:
		return 0, nil
	case int64:
		return float64(n), nil
	case uint64:
		return float64(n), nil
	case float64:
		return n, nil
	case []byte:
		return strconv.ParseFloat(string(n), 64)
	case string:
		return strconv.ParseFloat(n, 64)
	}

	return 0, fmt.Errorf("unexpected value type %T", v)
}
