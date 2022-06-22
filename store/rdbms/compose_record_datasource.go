package rdbms

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/qlng"
	"github.com/cortezaproject/corteza-server/pkg/report"
	"github.com/cortezaproject/corteza-server/store/rdbms/builders"
	"github.com/jmoiron/sqlx"
)

type (
	recordDatasource struct {
		name   string
		module *types.Module

		store *Store
		rows  *sql.Rows

		partitioned   bool
		partitionSize uint
		partitionCol  string

		baseFilter *qlng.ASTNode

		cols         report.FrameColumnSet
		q            squirrel.SelectBuilder
		nestLevel    int
		nestLabel    string
		levelColumns map[string]string
	}
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

	r.q, err = r.baseQuery(ld.Filter)
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
			Ref:     r.Name(),
			Columns: r.cols,
		},
	}
}

func (r *recordDatasource) Load(ctx context.Context, dd ...*report.FrameDefinition) (l report.Loader, c report.Closer, err error) {
	if r.partitioned {
		return r.partition(ctx, r.partitionSize, r.partitionCol, dd...)
	}

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
		q   = squirrel.Select().PlaceholderFormat(r.store.config.PlaceholderFormat)
		err error
	)

	auxLevelColumns := r.levelColumns
	r.levelColumns = make(map[string]string)
	groupCols := make(report.FrameColumnSet, 0, 10)

	// firstly handle GROUP BY definitions
	for _, k := range d.Keys {
		// - validate columns & functions
		err = k.Def.Traverse(func(n *qlng.ASTNode) (bool, *qlng.ASTNode, error) {
			if n.Symbol != "" {
				if _, ok := auxLevelColumns[n.Symbol]; !ok {
					return false, nil, fmt.Errorf("column %s does not exist on level %d", n.Symbol, r.nestLevel)
				}
			}

			// @todo do we need additional function validation besides the built-in one?

			return true, n, nil
		})
		if err != nil {
			return false, err
		}

		// AST transformation tasks
		tr := r.store.ASTTransformer(k.Def.ASTNode)
		outType, err := tr.Analyze(auxLevelColumns)
		if err != nil {
			return false, err
		}
		// - prepare frame col. definition
		c := report.MakeColumnOfKind(outType)
		c.Name = k.Name
		c.Label = k.Label
		if c.Label == "" {
			c.Label = c.Name
		}
		groupCols = append(groupCols, c)
		r.levelColumns[k.Name] = outType

		// - SQL things
		tr.SetPlaceholder(false)
		sql, _, err := tr.ToSql()
		if err != nil {
			return false, err
		}

		q = q.Column(squirrel.Alias(tr, c.Name)).
			GroupBy(sql)
	}

	// secondly handle column definitions
	for _, k := range d.Columns {
		// - make sure the column results with an aggregated value
		if !r.isAggregated(k.Def.ASTNode) {
			return false, fmt.Errorf("group column %s does not aggregate data", k.Name)
		}

		// - validate columns & functions
		err = k.Def.Traverse(func(n *qlng.ASTNode) (bool, *qlng.ASTNode, error) {
			if n.Ref == "count" && len(n.Args) == 0 {
				pc := r.cols.FirstPrimary()
				if pc == nil {
					return false, nil, errors.New("cannot use count(): no primary key defined")
				}
				n.Args = append(n.Args, &qlng.ASTNode{
					Ref: "distinct",
					Args: qlng.ASTNodeSet{{
						Symbol: pc.Name,
					}},
				})
			}

			if n.Symbol != "" {
				if _, ok := auxLevelColumns[n.Symbol]; !ok {
					return false, nil, fmt.Errorf("column %s does not exist on level %d", n.Symbol, r.nestLevel)
				}
			}

			// @todo do we need additional function validation besides the built-in one?

			return true, n, nil
		})
		if err != nil {
			return false, err
		}

		// AST transformation tasks
		tr := r.store.ASTTransformer(k.Def.ASTNode)
		outType, err := tr.Analyze(auxLevelColumns)
		if err != nil {
			return false, err
		}
		// - prepare frame col. definition
		c := report.MakeColumnOfKind(outType)
		c.Name = k.Name
		c.Label = k.Label
		if c.Label == "" {
			c.Label = c.Name
		}
		groupCols = append(groupCols, c)
		r.levelColumns[k.Name] = outType

		// - SQL things
		q = q.Column(squirrel.Alias(tr, c.Name))
	}

	if d.Filter != nil && d.Filter.ASTNode != nil {
		q = q.Having(r.store.ASTTransformer(d.Filter.ASTNode))
	}

	r.cols = groupCols
	r.q = q.FromSelect(r.q, fmt.Sprintf("l%d", r.nestLevel))
	return true, nil
}

// @todo add Transform

func (r *recordDatasource) Partition(partitionSize uint, partitionCol string) (bool, error) {
	if r.partitioned {
		return true, nil
	}
	if partitionCol == "" {
		return false, errors.New("unable to partition: partition column not defined")
	}

	r.partitioned = true
	r.partitionCol = partitionCol
	r.partitionSize = partitionSize
	return true, nil
}

func (r *recordDatasource) partition(ctx context.Context, partitionSize uint, partitionCol string, dd ...*report.FrameDefinition) (l report.Loader, c report.Closer, err error) {
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
	var prt squirrel.SelectBuilder
	if len(ss) > 0 {
		prt = squirrel.Select(fmt.Sprintf("*, row_number() over(partition by %s order by %s) as pp_rank", partitionCol, strings.Join(ss, ","))).
			PlaceholderFormat(r.store.config.PlaceholderFormat).
			FromSelect(q, "partition_base")
	} else {
		prt = squirrel.Select(fmt.Sprintf("*, row_number() over(partition by %s) as pp_rank", partitionCol)).
			PlaceholderFormat(r.store.config.PlaceholderFormat).
			FromSelect(q, "partition_base")
	}

	// the sort is already defined when partitioning so it's unneeded here
	q = squirrel.Select("*").
		PlaceholderFormat(r.store.config.PlaceholderFormat).
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
	if def.Filter != nil && def.Filter.ASTNode != nil || def.Sort != nil {
		q = squirrel.Select("*").PlaceholderFormat(r.store.config.PlaceholderFormat).FromSelect(q, "w_base")
	}

	// - filtering
	if def.Filter != nil && def.Filter.ASTNode != nil {
		q = q.Where(r.store.ASTTransformer(def.Filter.ASTNode))
	}

	return q, nil
}

func (r *recordDatasource) load(ctx context.Context, def *report.FrameDefinition, q squirrel.SelectBuilder) (l report.Loader, c report.Closer, err error) {
	var (
		sort filter.SortExprSet
	)

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

	// Make sure results are always sorted at least by primary keys
	var canPage bool
	if canPage, err = r.validateSort(def); err != nil {
		return
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

	return func(cap int, processed bool) ([]*report.Frame, error) {
			f := &report.Frame{
				Name:   def.Name,
				Source: def.Source,
				Ref:    r.Name(),
				Sort:   def.Sort,
				Filter: def.Filter,
			}

			checkCap := cap > 0

			// Fetch & convert the data.
			// Go 1 over the requested cap to be able to determine if there are
			// any additional pages
			i := 0
			f.Columns = def.Columns
			f.Rows = make(report.FrameRowSet, 0, cap)

			for r.rows.Next() {
				i++

				err = r.cast(r.rows, f)
				if err != nil {
					return nil, err
				}

				// If the count goes over the capacity, then we have a next page
				if checkCap && (processed && i > cap || !processed && i >= cap) {
					out := []*report.Frame{f}
					f = &report.Frame{
						Name:   def.Name,
						Source: def.Source,
						Ref:    r.Name(),
						Sort:   def.Sort,
						Filter: def.Filter,
					}
					i = 0
					if processed {
						return r.calculatePaging(out, def.Sort, uint(cap), def.Paging.PageCursor, canPage), nil
					}
					return out, nil
				}
			}

			if i > 0 {
				if processed {
					return r.calculatePaging([]*report.Frame{f}, def.Sort, uint(cap), def.Paging.PageCursor, canPage), nil
				} else {
					return []*report.Frame{f}, nil
				}
			}

			if processed {
				return []*report.Frame{f}, nil
			}

			// This indicates that the DS is empty
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
func (r *recordDatasource) baseQuery(f *report.Filter) (sqb squirrel.SelectBuilder, err error) {
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

	sqb = squirrel.Select("*").
		PlaceholderFormat(r.store.config.PlaceholderFormat).
		FromSelect(sqb, "q_base")

	// - any initial filtering we may need to do
	//
	// @todo better support functions and their validation.
	if f != nil {
		err = f.ASTNode.Traverse(func(n *qlng.ASTNode) (bool, *qlng.ASTNode, error) {
			// for now, a symbol indicates a column
			if n.Symbol != "" {
				if _, ok := r.levelColumns[n.Symbol]; !ok {
					return false, n, fmt.Errorf("column %s does not exist on level %d (%s)", n.Symbol, r.nestLevel, r.nestLabel)
				}
			}

			return true, n, nil
		})
		if err != nil {
			return
		}

		if f != nil && f.ASTNode != nil {
			sqb = sqb.Where(r.store.ASTTransformer(f.ASTNode))
		}
	}

	return sqb, nil
}

func (b *recordDatasource) calculatePaging(out []*report.Frame, sorting filter.SortExprSet, limit uint, cursor *filter.PagingCursor, canPage bool) []*report.Frame {
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

		if ignoreLimit || !canPage {
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

	// Handling multi value fields; rows with duplicated PKs get their multi value fields
	// merged together.
	next := false
	hasPrimary := false
	if len(out.Rows) > 0 {
		for i, c := range out.Columns {
			hasPrimary = hasPrimary || c.Primary
			if c.Primary {
				a := out.Rows[len(out.Rows)-1][i].(expr.Comparable)
				if o, err := a.Compare(r[i]); err != nil || o != 0 {
					next = true
					break
				}
			}
		}
	} else {
		next = true
	}

	next = next || !hasPrimary
	if next {
		out.Rows = append(out.Rows, r)
	} else {
		lastR := out.Rows[len(out.Rows)-1]
		for i, c := range out.Columns {
			if c.Multivalue {
				lastR.AppendCell(i, r[i])
			}
		}
	}
	return nil
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

func (r *recordDatasource) validateSort(def *report.FrameDefinition) (canPage bool, err error) {
	unique := ""
	def.Sort = func() filter.SortExprSet {
		for _, c := range r.cols {
			if c.Primary || c.Unique {
				if unique == "" {
					unique = c.Name
				}
				if def.Sort.Get(c.Name) != nil {
					unique = c.Name
					return def.Sort
				}
			}
		}
		if unique != "" {
			return append(def.Sort, &filter.SortExpr{Column: unique, Descending: def.Sort.LastDescending()})
		}
		return def.Sort
	}()
	return unique != "", nil
}

func (r *recordDatasource) isAggregated(n *qlng.ASTNode) bool {
	switch strings.ToLower(n.Ref) {
	case "count",
		"sum",
		"max",
		"min",
		"avg":
		return true
	}
	return false
}
