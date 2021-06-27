package rdbms

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/ql"
	"github.com/cortezaproject/corteza-server/pkg/report"
	"github.com/cortezaproject/corteza-server/pkg/slice"
	"github.com/jmoiron/sqlx"
)

type (
	recordDatasource struct {
		name   string
		module *types.Module

		// @todo use these
		supportedAggregationFunctions map[string]bool
		supportedFilterFunctions      map[string]bool

		store *Store
		rows  *sql.Rows

		baseFilter *report.RowDefinition

		cols         report.FrameColumnSet
		qBuilder     squirrel.SelectBuilder
		nestLevel    int
		levelColumns map[string]bool
	}
)

func ComposeRecordDatasourceBuilder(s *Store, module *types.Module, ld *report.LoadStepDefinition) (report.Datasource, error) {
	var err error

	r := &recordDatasource{
		name:         ld.Name,
		module:       module,
		store:        s,
		cols:         ld.Columns,
		levelColumns: make(map[string]bool),

		supportedAggregationFunctions: slice.ToStringBoolMap([]string{
			"COUNT",
			"SUM",
			"MAX",
			"MIN",
			"AVG",
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

	r.qBuilder, err = r.baseQuery(ld.Rows)

	return r, err
}

func (r *recordDatasource) Name() string {
	if r.name != "" {
		return r.name
	}

	if r.module == nil {
		return ""
	}

	if r.module.Handle == "" {
		return r.module.Name
	}
	return r.module.Handle
}

// @todo add Transform
// @todo try to make Group and Transform use the base query

func (r *recordDatasource) Group(d report.GroupDefinition, name string) (bool, error) {
	defer func() {
		r.nestLevel++
		r.name = name
	}()

	cls := r.levelColumns
	r.levelColumns = make(map[string]bool)

	gCols := make(report.FrameColumnSet, 0, 10)

	q := squirrel.Select()

	// @todo allow some transformation functions within the agg. functions
	parser := ql.NewParser()
	parser.OnFunction = r.stdAggregationHandler
	parser.OnIdent = func(i ql.Ident) (ql.Ident, error) {
		if !cls[i.Value] {
			return i, fmt.Errorf("column %s does not exist on level %d", i.Value, r.nestLevel)
		}

		i.Value = fmt.Sprintf("l%d.%s", r.nestLevel, i.Value)
		return i, nil
	}

	for _, g := range d.Groups {
		e, err := parser.ParseExpression(g.Expr)
		if err != nil {
			return false, err
		}

		// @todo imply based on context
		c := report.MakeColumnOfKind(g.Kind)
		c.Name = g.Name
		gCols = append(gCols, c)

		r.levelColumns[g.Name] = true
		q = q.Column(fmt.Sprintf("(%s) as `%s`", e.String(), g.Name)).
			GroupBy(e.String())
	}

	var e ql.ASTNode
	var err error
	for _, c := range d.Columns {
		if c.Aggregate != "" {
			e, err = parser.ParseExpression(fmt.Sprintf("%s(%s)", c.Aggregate, c.Expr))
			if err != nil {
				return false, err
			}
		} else {
			e, err = parser.ParseExpression(c.Expr)
			if err != nil {
				return false, err
			}
		}

		var col *report.FrameColumn
		if c.Kind != "" {
			col = report.MakeColumnOfKind(c.Kind)
		} else {
			// @todo imply based on context
			col = report.MakeColumnOfKind("Number")
		}
		col.Name = c.Name
		gCols = append(gCols, col)
		r.levelColumns[c.Name] = true

		q = q.
			Column(fmt.Sprintf("%s as `%s`", e.String(), c.Name))
	}

	if d.Rows != nil {
		// @todo validate groupping functions
		hh, err := r.rowFilterToString("", gCols, d.Rows)
		if err != nil {
			return false, err
		}
		q = q.Having(hh)
	}

	r.cols = gCols
	r.qBuilder = q.FromSelect(r.qBuilder, fmt.Sprintf("l%d", r.nestLevel))
	return true, nil
}

func (r *recordDatasource) Load(ctx context.Context, dd ...*report.FrameDefinition) (l report.Loader, c report.Closer, err error) {
	def := dd[0]

	q, err := r.preloadQuery(def)
	if err != nil {
		return nil, nil, err
	}

	return r.load(ctx, def, q)
}

func (r *recordDatasource) Partition(ctx context.Context, partitionSize uint, partitionCol string, dd ...*report.FrameDefinition) (l report.Loader, c report.Closer, err error) {
	def := dd[0]

	q, err := r.preloadQuery(def)
	if err != nil {
		return nil, nil, err
	}

	// the partitioning wrap
	// @todo move this to the DB driver package?
	// @todo squash the query a bit? try to move most of this to the base query to remove
	//       one sub-select
	prt := squirrel.Select(fmt.Sprintf("*, row_number() over(partition by %s order by %s) as pp_rank", partitionCol, partitionCol)).
		FromSelect(q, "partition_base")

	// @odo make it better, please...
	ss, err := r.sortExpr(def)
	if err != nil {
		return nil, nil, err
	}

	q = squirrel.Select("*").
		FromSelect(prt, "partition_wrap").
		Where(fmt.Sprintf("pp_rank <= %d", partitionSize)).
		OrderBy(ss...)

	return r.load(ctx, def, q)
}

func (r *recordDatasource) preloadQuery(def *report.FrameDefinition) (squirrel.SelectBuilder, error) {
	q := r.qBuilder

	// assure columns
	// - undefined columns = all columns
	if len(def.Columns) == 0 {
		def.Columns = r.cols
	} else {
		// - make sure they exist
		cc := make(report.FrameColumnSet, len(def.Columns))
		for i, c := range def.Columns {
			ci := r.cols.Find(c.Name)
			if ci == -1 {
				return q, fmt.Errorf("column not found: %s", c.Name)
			}
			cc[i] = r.cols[ci]
		}
		def.Columns = cc
	}

	// when filtering/sorting, wrap the base query in a sub-select, so we don't need to
	// worry about exact column names.
	//
	// @todo flatten the query
	if def.Rows != nil || def.Sorting != nil {
		wrap := squirrel.Select("*").FromSelect(q, "w_base")

		// additional filtering
		if def.Rows != nil {
			f, err := r.rowFilterToString("", r.cols, def.Rows)
			if err != nil {
				return q, err
			}
			wrap = wrap.Where(f)
		}

		// additional sorting
		if len(def.Sorting) > 0 {
			ss, err := r.sortExpr(def)
			if err != nil {
				return q, err
			}

			wrap = wrap.OrderBy(ss...)
		}

		q = wrap
	}

	return q, nil
}

func (r *recordDatasource) sortExpr(def *report.FrameDefinition) ([]string, error) {
	ss := make([]string, len(def.Sorting))
	for i, c := range def.Sorting {
		ci := r.cols.Find(c.Column)
		if ci == -1 {
			return nil, fmt.Errorf("sort column not resolved: %s", c.Column)
		}

		_, _, typeCast, err := r.store.config.CastModuleFieldToColumnType(r.cols[ci], c.Column)
		if err != nil {
			return nil, err
		}

		ss[i] = r.store.config.SqlSortHandler(fmt.Sprintf(typeCast, c.Column), c.Descending)
	}

	return ss, nil
}

func (r *recordDatasource) load(ctx context.Context, def *report.FrameDefinition, q squirrel.SelectBuilder) (l report.Loader, c report.Closer, err error) {
	r.rows, err = r.store.Query(ctx, q)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot execute query: %w", err)
	}

	return func(cap int) ([]*report.Frame, error) {
			f := &report.Frame{
				Name:   def.Name,
				Source: def.Source,
				Ref:    def.Ref,
			}

			checkCap := cap > 0

			// fetch & convert the data
			i := 0
			// @todo make it in place
			f.Columns = def.Columns
			f.Rows = make(report.FrameRowSet, 0, cap)

			for r.rows.Next() {
				i++

				err = r.Cast(r.rows, f)
				if err != nil {
					return nil, err
				}

				if checkCap && i >= cap {
					out := []*report.Frame{f}
					f = &report.Frame{}
					i = 0
					return out, nil
				}
			}

			if i > 0 {
				return []*report.Frame{f}, nil
			}
			return nil, nil
		}, func() {
			if r.rows == nil {
				return
			}
			r.rows.Close()
		}, nil
}

// @todo handle those rv_ prefixes; for now omitted
func (r *recordDatasource) baseQuery(f *report.RowDefinition) (sqb squirrel.SelectBuilder, err error) {
	var (
		joinTpl = "compose_record_value AS %s ON (%s.record_id = crd.id AND %s.name = '%s' AND %s.deleted_at IS NULL)"

		report = r.store.composeRecordsSelectBuilder().
			Where("crd.deleted_at IS NULL").
			Where("crd.module_id = ?", r.module.ID)
	)

	// Prepare all of the mod columns
	// @todo make this as small as possible!
	for _, f := range r.module.Fields {
		report = report.LeftJoin(strings.ReplaceAll(joinTpl, "%s", f.Name)).
			Column(f.Name + ".value as " + f.Name)

		r.levelColumns[f.Name] = true
	}

	if f != nil {
		// @todo functions and function validation
		parser := ql.NewParser()
		parser.OnIdent = func(i ql.Ident) (ql.Ident, error) {
			if !r.levelColumns[i.Value] {
				return i, fmt.Errorf("column %s does not exist on level %d", i.Value, r.nestLevel)
			}

			return i, nil
		}

		fl, err := r.rowFilterToString("", r.cols, f)
		if err != nil {
			return sqb, err
		}
		astq, err := parser.ParseExpression(fl)
		if err != nil {
			return sqb, err
		}
		report = report.Where(astq.String())
	}

	return report, nil
}

func (b *recordDatasource) Cast(row sqlx.ColScanner, out *report.Frame) error {
	var err error
	aux := make(map[string]interface{})
	if err = sqlx.MapScan(row, aux); err != nil {
		return err
	}

	r := make(report.FrameRow, len(out.Columns))

	k := ""
	for i, c := range out.Columns {
		k = "" + c.Name
		v, ok := aux[k]
		if !ok {
			continue
		}

		// @todo improve json value casting; I couldn't figure it out at the time

		// @todo this doesn't work...
		// var aux interface{}
		// err = json.Unmarshal(v.([]byte), &aux)
		// if err != nil {
		// 	return err
		// }

		switch cv := v.(type) {
		case []byte:
			c, err := c.Caster(string(cv))
			if err != nil {
				return err
			}
			r[i] = c

		case uint64, int64:
			r[i], err = c.Caster(cv)
			if err != nil {
				return err
			}

		default:
			if isNil(cv) {
				continue
			}

			c, err := c.Caster(cv)
			if err != nil {
				return err
			}
			r[i] = c

		}

	}

	out.Rows = append(out.Rows, r)
	return nil
}

// Identifiers should be names of the fields (physical table columns OR json fields, defined in module)
func (b *recordDatasource) stdAggregationHandler(f ql.Function) (ql.ASTNode, error) {
	if !b.supportedAggregationFunctions[strings.ToUpper(f.Name)] {
		return f, fmt.Errorf("unsupported aggregate function %q", f.Name)
	}

	return b.store.SqlFunctionHandler(f)
}

func (ds *recordDatasource) rowFilterToString(conjunction string, cc report.FrameColumnSet, def ...*report.RowDefinition) (string, error) {
	// The fields on the root level of the definition take priority
	base := ""
	for _, f := range def {
		if f == nil {
			continue
		}

		if f.Cells != nil {
			for k, op := range f.Cells {
				ci := cc.Find(k)
				if ci == -1 {
					return "", fmt.Errorf("filtered column not found in the data frame: %s", k)
				}
				_, _, typeCast, err := ds.store.config.CastModuleFieldToColumnType(cc[ci], k)
				if err != nil {
					return "", err
				}

				col := fmt.Sprintf(typeCast, k)
				base += fmt.Sprintf("%s%s%s %s ", col, op.OpToCmp(), fmt.Sprintf(typeCast, op.Value), conjunction)
			}
		}
	}
	if strings.TrimSpace(base) != "" {
		base = strings.TrimSuffix(base, " "+conjunction+" ")
	}

	// Nested AND
	for _, f := range def {
		if f == nil {
			continue
		}

		if f.And != nil {
			na, err := ds.rowFilterToString("AND", cc, f.And...)
			if err != nil {
				return "", err
			}
			if strings.TrimSpace(base) != "" {
				base += fmt.Sprintf(" %s (%s)", conjunction, na)
			} else {
				base = na
			}
		}
	}

	// Nested OR
	for _, f := range def {
		if f == nil {
			continue
		}

		if f.Or != nil {
			na, err := ds.rowFilterToString("OR", cc, f.Or...)
			if err != nil {
				return "", err
			}
			if strings.TrimSpace(base) != "" {
				base += fmt.Sprintf(" %s (%s)", conjunction, na)
			} else {
				base = na
			}
		}
	}

	return strings.TrimSpace(base), nil
}

func isNil(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}
