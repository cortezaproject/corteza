package rdbms

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/ql"
	"github.com/cortezaproject/corteza-server/pkg/report"
	"github.com/cortezaproject/corteza-server/pkg/slice"
	"github.com/cortezaproject/corteza-server/store/rdbms/builders"
	"github.com/jmoiron/sqlx"
)

type (
	recordDatasource struct {
		name   string
		module *types.Module

		// @todo use these
		supportedAggregationFunctions map[string]string
		supportedFilterFunctions      map[string]bool

		store *Store
		rows  *sql.Rows

		baseFilter *report.RowDefinition

		cols         report.FrameColumnSet
		q            squirrel.SelectBuilder
		nestLevel    int
		nestLabel    string
		levelColumns map[string]string
	}
)

var (
	supportedAggregationFunctions = slice.ToStringBoolMap([]string{
		"COUNT",
		"SUM",
		"MAX",
		"MIN",
		"AVG",
	})

	// supportedGroupingFunctions = ...
)

// ComposeRecordDatasourceBuilder initializes and returns a datasource builder for compose record resource
//
// @todo try to make the resulting query as flat as possible
func ComposeRecordDatasourceBuilder(s *Store, module *types.Module, ld *report.LoadStepDefinition) (report.Datasource, error) {
	var err error

	r := &recordDatasource{
		name:   ld.Name,
		module: module,
		store:  s,
		cols:   ld.Columns,

		// levelColumns help us keep track of what columns are currently available
		levelColumns: make(map[string]string),
	}

	r.q, err = r.baseQuery(ld.Rows)
	return r, err
}

// Name returns the name we should use when referencing this datasource
//
// The name is determined from the user-specified name, or implied from the context.
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

func (r *recordDatasource) Describe() report.FrameDescriptionSet {
	return report.FrameDescriptionSet{
		&report.FrameDescription{
			Source:  r.Name(),
			Columns: r.cols,
		},
	}
}

func (r *recordDatasource) Load(ctx context.Context, dd ...*report.FrameDefinition) (l report.Loader, c report.Closer, err error) {
	def := dd[0]

	q, err := r.preloadQuery(def)
	if err != nil {
		return nil, nil, err
	}

	return r.load(ctx, def, q)
}

// Group instructs the datasource to provide grouped and aggregated output
func (r *recordDatasource) Group(d report.GroupDefinition, name string) (bool, error) {
	defer func() {
		r.nestLevel++
		r.nestLabel = "group"
		r.name = name
	}()

	var (
		q       = squirrel.Select()
		auxKind = ""
		ok      = false
	)

	auxLevelColumns := r.levelColumns
	r.levelColumns = make(map[string]string)
	groupCols := make(report.FrameColumnSet, 0, 10)

	for _, g := range d.Keys {
		auxKind, ok = auxLevelColumns[g.Column]
		if !ok {
			return false, fmt.Errorf("column %s does not exist on level %d", g.Column, r.nestLevel)
		}

		// @todo...
		// if g.Group != "" {...}

		c := report.MakeColumnOfKind(auxKind)
		c.Name = g.Name
		c.Label = g.Label
		if c.Label == "" {
			c.Label = c.Name
		}
		groupCols = append(groupCols, c)

		r.levelColumns[g.Name] = auxKind
		q = q.Column(fmt.Sprintf("%s as `%s`", g.Column, g.Name)).
			GroupBy(g.Column)
	}

	var aggregate string
	for _, c := range d.Columns {
		aggregate = strings.ToUpper(c.Aggregate)

		if c.Column == "" {
			if c.Aggregate == "" {
				return false, fmt.Errorf("the aggregation function is required when the column is omitted")
			}
		} else {
			auxKind, ok = auxLevelColumns[c.Column]
			if !ok {
				return false, fmt.Errorf("column %s does not exist on level %d", c.Column, r.nestLevel)
			}
		}

		qParam := c.Column
		if c.Aggregate != "" {
			if !supportedAggregationFunctions[aggregate] {
				return false, fmt.Errorf("aggregation function not supported: %s", c.Aggregate)
			}

			// when an aggregation function is defined, the output is always numeric
			auxKind = "Number"

			qParam = fmt.Sprintf("%s(%s)", aggregate, c.Column)
		} else if qParam == "" {
			qParam = "*"
		}

		col := report.MakeColumnOfKind(auxKind)
		col.Name = c.Name
		col.Label = c.Label
		if col.Label == "" {
			col.Label = col.Name
		}
		groupCols = append(groupCols, col)
		r.levelColumns[c.Name] = auxKind

		q = q.
			Column(squirrel.Alias(squirrel.Expr(qParam), c.Name))
	}

	if d.Rows != nil {
		// @todo validate groupping functions
		hh, err := r.rowFilterToString("", groupCols, d.Rows)
		if err != nil {
			return false, err
		}
		q = q.Having(hh)
	}

	r.cols = groupCols
	r.q = q.FromSelect(r.q, fmt.Sprintf("l%d", r.nestLevel))
	return true, nil
}

// @todo add Transform

func (r *recordDatasource) Partition(ctx context.Context, partitionSize uint, partitionCol string, dd ...*report.FrameDefinition) (l report.Loader, c report.Closer, err error) {
	def := dd[0]

	q, err := r.preloadQuery(def)
	if err != nil {
		return nil, nil, err
	}

	var ss []string
	if len(def.Sort) > 0 {
		ss, err = r.sortExpr(def.Sort)
		if err != nil {
			return nil, nil, err
		}
	}

	// the partitioning wrap
	// @todo move this to the DB driver package?
	// @todo squash the query a bit? try to move most of this to the base query to remove
	//       one sub-select
	prt := squirrel.Select(fmt.Sprintf("*, row_number() over(partition by %s order by %s) as pp_rank", partitionCol, strings.Join(ss, ","))).
		FromSelect(q, "partition_base")

	// the sort is already defined when partitioning so it's unneeded here
	q = squirrel.Select("*").
		FromSelect(prt, "partition_wrap").
		Where(fmt.Sprintf("pp_rank <= %d", partitionSize))

	return r.load(ctx, def, q)
}

func (r *recordDatasource) preloadQuery(def *report.FrameDefinition) (squirrel.SelectBuilder, error) {
	q := r.q

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
	if def.Rows != nil || def.Sort != nil {
		q = squirrel.Select("*").FromSelect(q, "w_base")
	}

	// - filtering
	if def.Rows != nil {
		f, err := r.rowFilterToString("", r.cols, def.Rows)
		if err != nil {
			return q, err
		}
		q = q.Where(f)
	}

	return q, nil
}

func (r *recordDatasource) load(ctx context.Context, def *report.FrameDefinition, q squirrel.SelectBuilder) (l report.Loader, c report.Closer, err error) {
	sort := def.Sort

	// - paging related stuff
	if def.Paging.PageCursor != nil {
		// Page cursor exists so we need to validate it against used sort
		// To cover the case when paging cursor is set but sorting is empty, we collect the sorting instructions
		// from the cursor.
		// This (extracted sorting info) is then returned as part of response
		if def.Sort, err = def.Paging.PageCursor.Sort(def.Sort); err != nil {
			return nil, nil, err
		}
	}

	// Cloned sorting instructions for the actual sorting
	// Original must be kept for cursor creation
	sort = def.Sort.Clone()

	// When cursor for a previous page is used it's marked as reversed
	// This tells us to flip the descending flag on all used sort keys
	if def.Paging.PageCursor != nil && def.Paging.PageCursor.ROrder {
		sort.Reverse()
	}

	if def.Paging.PageCursor != nil {
		q = q.Where(builders.CursorCondition(def.Paging.PageCursor, nil))
	}

	if len(sort) > 0 {
		ss, err := r.sortExpr(sort)
		if err != nil {
			return nil, nil, err
		}

		q = q.OrderBy(ss...)
	}

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

			// Fetch & convert the data.
			// Go 1 over the requested cap to be able to determine if there are
			// any additional pages
			i := 0
			f.Columns = def.Columns
			f.Rows = make(report.FrameRowSet, 0, cap+1)

			for r.rows.Next() {
				i++

				err = r.cast(r.rows, f)
				if err != nil {
					return nil, err
				}

				// If the count goes over the capacity, then we have a next page
				if checkCap && i > cap {
					out := []*report.Frame{f}
					f = &report.Frame{}
					i = 0
					return r.calculatePaging(out, def.Sort, uint(cap), def.Paging.PageCursor), nil
				}
			}

			if i > 0 {
				return r.calculatePaging([]*report.Frame{f}, def.Sort, uint(cap), def.Paging.PageCursor), nil
			}
			return nil, nil
		}, func() {
			if r.rows == nil {
				return
			}
			r.rows.Close()
		}, nil
}

// baseQuery prepares the initial SQL that will be used for data access
//
// The query includes all of the requested columns in the required types to avid the need to type cast.
func (r *recordDatasource) baseQuery(f *report.RowDefinition) (sqb squirrel.SelectBuilder, err error) {
	var (
		joinTpl = "compose_record_value AS %s ON (%s.record_id = crd.id AND %s.name = '%s' AND %s.deleted_at IS NULL)"
	)

	// - the initial set of available columns
	//
	// @todo at what level should the requested columns be validated?
	r.nestLevel = 0
	r.nestLabel = "base"
	for _, c := range r.cols {
		r.levelColumns[c.Name] = c.Kind
	}

	// - base query
	sqb = r.store.SelectBuilder("compose_record AS crd").
		Where("crd.deleted_at IS NULL").
		Where("crd.module_id = ?", r.module.ID).
		Where("crd.rel_namespace = ?", r.module.NamespaceID)

	// - based on the definition, preload the columns
	var (
		col      string
		is       bool
		isJoined = make(map[string]bool)
	)
	for _, c := range r.cols {
		if isJoined[c.Name] {
			continue
		}
		isJoined[c.Name] = true

		// native record columns don't need any extra handling
		if col, _, is = isRealRecordCol(c.Name); is {
			sqb = sqb.Column(squirrel.Alias(squirrel.Expr(col), c.Name))
			continue
		}

		// non-native record columns need to have their type casted before use
		_, _, tcp, _ := r.store.config.CastModuleFieldToColumnType(c, c.Name)
		sqb = sqb.LeftJoin(strings.ReplaceAll(joinTpl, "%s", c.Name)).
			Column(squirrel.Alias(squirrel.Expr(fmt.Sprintf(tcp, c.Name+".value")), c.Name))
	}

	// - any initial filtering we may need to do
	//
	// @todo better support functions and their validation.
	if f != nil {
		parser := ql.NewParser()
		parser.OnIdent = func(i ql.Ident) (ql.Ident, error) {
			if _, ok := r.levelColumns[i.Value]; !ok {
				return i, fmt.Errorf("column %s does not exist on level %d (%s)", i.Value, r.nestLevel, r.nestLabel)
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

		sqb = sqb.Where(astq.String())
	}

	return sqb, nil
}

func (b *recordDatasource) calculatePaging(out []*report.Frame, sorting filter.SortExprSet, limit uint, cursor *filter.PagingCursor) []*report.Frame {
	for _, o := range out {
		var (
			hasPrev       = cursor != nil
			hasNext       bool
			ignoreLimit   = limit == 0
			reversedOrder = cursor != nil && cursor.ROrder
		)

		hasNext = uint(len(o.Rows)) > limit
		if !ignoreLimit && uint(len(o.Rows)) > limit {
			o.Rows = o.Rows[:limit]
		}

		if reversedOrder {
			// Fetched set needs to be reversed because we've forced a descending order to get the previous page
			for i, j := 0, len(o.Rows)-1; i < j; i, j = i+1, j-1 {
				o.Rows[i], o.Rows[j] = o.Rows[j], o.Rows[i]
			}

			// when in reverse-order rules on what cursor to return change
			hasPrev, hasNext = hasNext, hasPrev
		}

		if ignoreLimit {
			return out
		}

		if hasPrev {
			o.Paging = &filter.Paging{}
			o.Paging.PrevPage = o.CollectCursorValues(o.FirstRow(), sorting...)
			o.Paging.PrevPage.ROrder = true
			o.Paging.PrevPage.LThen = !sorting.Reversed()
		}

		if hasNext {
			if o.Paging == nil {
				o.Paging = &filter.Paging{}
			}
			o.Paging.NextPage = o.CollectCursorValues(o.LastRow(), sorting...)
			o.Paging.NextPage.LThen = sorting.Reversed()
		}
	}

	return out
}

func (b *recordDatasource) cast(row sqlx.ColScanner, out *report.Frame) error {
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
	if !supportedAggregationFunctions[strings.ToUpper(f.Name)] {
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

func (r *recordDatasource) sortExpr(sorting filter.SortExprSet) ([]string, error) {
	ss := make([]string, len(sorting))
	for i, c := range sorting {
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
