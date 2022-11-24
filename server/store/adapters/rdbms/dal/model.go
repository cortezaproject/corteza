package dal

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/cortezaproject/corteza/server/pkg/errors"
	"github.com/modern-go/reflect2"

	"github.com/cortezaproject/corteza/server/pkg/dal"
	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/store/adapters/rdbms/drivers"
	"github.com/cortezaproject/corteza/server/store/adapters/rdbms/ql"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/jmoiron/sqlx"
)

type (
	queryRunner interface {
		sqlx.QueryerContext
		sqlx.ExecerContext
	}

	model struct {
		model *dal.Model
		conn  queryRunner

		dialect drivers.Dialect

		table drivers.TableCodec
	}

	parsedFilter interface {
		ExpressionParsed() *ql.ASTNode
	}

	queryParser interface {
		Convert(*ql.ASTNode) (out exp.Expression, err error)
		Parse(string) (out exp.Expression, err error)
	}
)

func validate(m *dal.Model) error {
	var (
		c2c = make(map[string]*dal.Attribute, len(m.Attributes))
	)

	for _, a := range m.Attributes {
		if a.StoreIdent() == "" {
			return fmt.Errorf("attribute %q has no ident", a.Ident)
		}

		if a.Store == nil {
			return fmt.Errorf("attribute %q has no store codec", a.Ident)
		}

		usedBy, has := c2c[a.StoreIdent()]
		if has {
			// column already in the map
			if a.Store.SingleValueOnly() {
				return fmt.Errorf("attribute %q has single value codec but column %q is already used by attribute %q", a.Ident, a.StoreIdent(), usedBy.Ident)
			}
		}

		c2c[a.StoreIdent()] = a
	}

	return nil
}

// Model returns fully initialized model store
//
// It abstracts database table and its columns and provides unified interface
// for fetching and storing records.
func Model(m *dal.Model, c queryRunner, d drivers.Dialect) *model {
	var (
		ms = &model{
			model:   m,
			conn:    c,
			dialect: d,
			table:   drivers.NewTableCodec(m, d),
		}
	)

	return ms
}

// parseQuery parses the query into the goqu expression
//
// The parse query initializes a fresh converter instance because QL converter
// uses some internal state to keep track of symbols and stuff.
//
// When doing parallel requests over the same model, unexpected... stuff happens.
//
// Alternative solution would introduce a mutes on the internal model but that
// is probably the same or worse as this.
//
func (d *model) parseQuery(q string) (out exp.Expression, err error) {
	return d.QueryParser().Parse(q)
}

func (d *model) convertQuery(n *ql.ASTNode) (out exp.Expression, err error) {
	return d.QueryParser().Convert(n)
}

// QueryParser returns ql struct that allows parsing query strings or converting AST into expression
// @todo benchmark to see if this re-init is a bad idea; I don't think it should be
//       since we're just initializing fairly light structs.
func (d *model) QueryParser() queryParser {
	return ql.Converter(
		ql.SymHandler(d.qlConverterGenericSymHandler()),
		ql.RefHandler(d.qlConverterGenericRefHandler()),
	)
}

func (d *model) qlConverterGenericSymHandler() func(node *ql.ASTNode) (exp.Expression, error) {
	return func(node *ql.ASTNode) (exp.Expression, error) {
		// @note normalize system idents on the RDBMS level for filters
		//       offloaded to the database.
		sym := dal.NormalizeAttrNames(node.Symbol)
		if d.model.ResourceType == "corteza::compose:module" {
			// temporary solution
			//
			// before DAL some fields were aliased (recordID => ID)
			switch sym {
			case "recordID":
				sym = "ID"
			}
		}

		if node.Meta == nil {
			node.Meta = make(map[string]any)
		}

		attr := d.model.Attributes.FindByIdent(sym)
		if attr == nil {
			return nil, fmt.Errorf("unknown attribute %q used in query expression", node.Symbol)
		}

		node.Meta["dal.Attribute"] = attr
		node.Meta["dal.Model"] = d.model

		if !attr.Filterable {
			return nil, fmt.Errorf("attribute %q can not be used in query expression", attr.Ident)
		}

		return d.table.AttributeExpression(sym)
	}
}

func (d *model) qlConverterGenericRefHandler() func(*ql.ASTNode, ...exp.Expression) (exp.Expression, error) {
	return func(node *ql.ASTNode, args ...exp.Expression) (exp.Expression, error) {
		return d.dialect.ExprHandler(node, args...)
	}
}

func (d *model) Truncate(ctx context.Context) error {
	sql, args, err := d.truncateSql().ToSQL()
	if err != nil {
		return err
	}

	_, err = d.conn.ExecContext(ctx, sql, args...)
	return err
}

func (d *model) Create(ctx context.Context, rr ...dal.ValueGetter) error {
	sql, args, err := d.insertSql(rr...).ToSQL()
	if err != nil {
		return err
	}

	_, err = d.conn.ExecContext(ctx, sql, args...)
	return err
}

func (d *model) Update(ctx context.Context, r dal.ValueGetter) error {
	sql, args, err := d.updateSql(r).ToSQL()
	if err != nil {
		return err
	}

	_, err = d.conn.ExecContext(ctx, sql, args...)
	return err
}

func (d *model) Delete(ctx context.Context, r dal.ValueGetter) error {
	sql, args, err := d.deleteSql(r).ToSQL()
	if err != nil {
		return err
	}

	_, err = d.conn.ExecContext(ctx, sql, args...)
	return err
}

func (d *model) Search(f filter.Filter) (i *iterator, err error) {
	var (
		orderBy = f.OrderBy()
	)

	if f.Cursor() != nil {
		// Page cursor exists; we need to validate it against used sort
		// To cover the case when paging cursor is set but sorting is empty, we collect the sorting instructions
		// from the cursor.
		// This (extracted sorting info) is then returned as part of response
		if orderBy, err = f.Cursor().Sort(orderBy); err != nil {
			return nil, err
		}
	}

	// validate order-by
	for _, s := range orderBy {
		if _, err = d.table.AttributeExpression(s.Column); err != nil {
			return nil, err
		}

		if att := d.model.Attributes.FindByIdent(s.Column); att != nil && att.MultiValue {
			return nil, fmt.Errorf("not allowed to sort by multi-value attribute: %s", s.Column)
		}

	}

	// sanitize filter a bit
	for _, c := range d.table.Columns() {
		if !c.IsPrimaryKey() {
			continue
		}

		attrIdent := c.Attribute().Ident
		if orderBy.Get(attrIdent) != nil {
			continue
		}

		if !c.Attribute().PrimaryKey && !c.Attribute().Sortable {
			return nil, fmt.Errorf("can not sort by %q; not sortable, not primary key", attrIdent)
		}

		// Make sure results are always sorted at least by primary key
		orderBy = append(orderBy, &filter.SortExpr{Column: attrIdent, Descending: orderBy.LastDescending()})
	}

	i = &iterator{
		// source and destination is the same
		src: d,
		dst: d,

		sorting: orderBy,
		cursor:  f.Cursor(),
		limit:   f.Limit(),
	}

	i.query = d.searchSql(f)
	if err = i.query.Error(); err != nil {
		return
	}

	return
}

// Aggregate constructs SELECT sql with group-by and an optional having CLAUSE
//
// All group-by attributes are prepended to aggregation
// expressions when constructing expressions & columns to select from.
//
// Passing in filter with cursor, empty groupBy or aggrExpr slice will result in an error
func (d *model) Aggregate(f filter.Filter, groupBy []dal.AggregateAttr, aggrExpr []dal.AggregateAttr, having *ql.ASTNode) (i *iterator, err error) {
	if len(groupBy) == 0 {
		return nil, fmt.Errorf("can not run aggregation without group-by")
	}

	i = &iterator{
		cursor:  f.Cursor(),
		sorting: f.OrderBy(),
		limit:   f.Limit(),
	}

	var (
		// source model; how data we are reading from is shaped
		srcModel = &dal.Model{}

		// destination model; how data we are reading into is shaped
		dstModel = &dal.Model{}

		srcAttr, dstAttr *dal.Attribute
	)

	// prepare a bit modified module that
	// describes aggregated columns (prepending attributes used for group-by)
	for _, c := range append(groupBy, aggrExpr...) {
		srcAttr = &dal.Attribute{
			Ident: c.Identifier,
			Type:  c.Type,
			Store: c.Store,
		}
		srcModel.Attributes = append(srcModel.Attributes, srcAttr)

		dstAttr = &dal.Attribute{
			Ident: c.Identifier,
			Type:  c.Type,
		}
		dstModel.Attributes = append(dstModel.Attributes, dstAttr)
	}

	i.src = Model(srcModel, d.conn, d.dialect)
	i.dst = Model(dstModel, d.conn, d.dialect)

	i.query = d.aggregateSql(f, groupBy, aggrExpr, having)

	if err = i.query.Error(); err != nil {
		return
	}

	return
}

func (d *model) Lookup(ctx context.Context, pkv dal.ValueGetter, r dal.ValueSetter) (err error) {
	query, args, err := d.lookupSql(pkv).ToSQL()
	if err != nil {
		return
	}

	// using sql.Rows instead of a row
	// this gives us more control over closing (the rows resource)
	// and ability to use sql.RawBytes
	var rows *sql.Rows
	rows, err = d.conn.QueryContext(ctx, query, args...)
	if err != nil {
		return
	}

	if err = rows.Err(); err != nil {
		return
	}

	defer rows.Close()
	if !rows.Next() {
		return errors.NotFound("not found")
	}

	scanBuf := d.table.MakeScanBuffer()
	if err = rows.Scan(scanBuf...); err != nil {
		return
	}

	if err = d.table.Decode(scanBuf, r); err != nil {
		return
	}

	return rows.Close()
}

// constructs SQL for selecting records from a table,
// converting parts of record filter into conditions
//
// Does not add any limits, sorting or any cursor conditions!
func (d *model) searchSql(f filter.Filter) *goqu.SelectDataset {
	var (
		err  error
		base = d.selectSql()
		tmp  exp.Expression
		cnd  []exp.Expression
	)

	{
		// Add model & namespace constraints when model expects (has configured attributes) for them
		//
		// This covers both scenarios:
		// 1) Model is configured to store records in a dedicated table
		//    without model and/or namespace attributes
		//
		// 2) Model has model and/or namespace attribute and saves records
		//    from different modules in the same table

		// if d.sysExprNamespaceID != nil {
		//	cnd = append(cnd, d.sysExprNamespaceID.Eq(f.NamespaceID))
		// } else {
		//	// @todo check if f.NamespaceID is compatible
		// }
		//
		// if d.sysExprModuleID != nil {
		//	cnd = append(cnd, d.sysExprModuleID.Eq(f.ModuleID))
		// } else {
		//	// @todo check if f.ModuleID is compatible
		// }
	}

	cc := f.Constraints()
	if d.model.Constraints != nil {
		if cc == nil {
			cc = d.model.Constraints
		} else {
			for k, c := range d.model.Constraints {
				// Overwrite user-provided constraints as the system ones are more important
				cc[k] = c
			}
		}
	}

	for ident, vv := range cc {
		attr := d.model.Attributes.FindByIdent(ident)
		if attr == nil {
			return base.SetError(fmt.Errorf("unknown attribute %q used for constrant", ident))
		}

		// @note why?
		// if !attr.PrimaryKey {
		// 	continue
		// }

		var attrExpr exp.Expression
		attrExpr, err = d.table.AttributeExpression(attr.Ident)
		if err != nil {
			return base.SetError(err)
		}

		if len(vv) > 0 {
			cnd = append(cnd, exp.NewBooleanExpression(exp.InOp, attrExpr, vv))
		}
	}

	for ident, state := range f.StateConstraints() {
		attr := d.model.Attributes.FindByIdent(ident)
		if attr == nil {
			if ident == "deletedAt" {
				// workaround for situation where filter constraints
				// contain deletedAt but attribute does not exist
				continue
			}

			return base.SetError(fmt.Errorf("unknown attribute %q used for state constrant", ident))
		}

		if !attr.Type.IsNullable() {
			// @todo this must be checked much earlier
			return base.SetError(fmt.Errorf("can not use non-nullable attribute %q soft-deleting", attr.Ident))
		}

		if state == filter.StateInclusive {
			// not used so we don't rea
			continue
		}

		var attrExpr exp.Expression
		attrExpr, err = d.table.AttributeExpression(attr.Ident)
		if err != nil {
			return base.SetError(err)
		}

		switch state {
		case filter.StateExclusive:
			// only not-null values
			cnd = append(cnd, exp.NewLiteralExpression("? IS NULL", attrExpr))

		case filter.StateExcluded:
			// exclude all non-null values
			cnd = append(cnd, exp.NewLiteralExpression("? IS NULL", attrExpr))
		}
	}

	if mc := f.MetaConstraints(); len(mc) > 0 {
		attr := d.model.Attributes.FindByIdent("meta")
		if attr == nil {
			return base.SetError(fmt.Errorf("can not filter records in this module by meta: no meta attribute defined"))
		}

		var (
			metaKeyExpr   exp.Expression
			metaAttrIdent = exp.NewIdentifierExpression("", d.model.Ident, attr.Ident)
		)

		for mKey, mVal := range f.MetaConstraints() {
			metaKeyExpr, err = d.dialect.JsonExtractUnquote(metaAttrIdent, mKey)
			if err != nil {
				return base.SetError(err)
			}

			if reflect2.IsNil(mVal) {
				cnd = append(cnd, exp.NewBooleanExpression(exp.IsNotOp, metaKeyExpr, nil))
			} else {
				cnd = append(cnd, exp.NewBooleanExpression(exp.EqOp, metaKeyExpr, mVal))
			}
		}
	}

	if pf, ok := f.(parsedFilter); ok {
		tmp := pf.ExpressionParsed()
		if tmp != nil {
			n, err := d.convertQuery(tmp)
			if err != nil {
				return base.SetError(err)
			}
			cnd = append(cnd, n)
		}
	}

	if q := strings.TrimSpace(f.Expression()); len(q) > 0 {
		if tmp, err = d.parseQuery(q); err != nil {
			return base.SetError(err)
		}

		cnd = append(cnd, tmp)
	}

	return base.Where(cnd...)
}

func (d *model) aggregateSql(f filter.Filter, groupBy []dal.AggregateAttr, out []dal.AggregateAttr, having *ql.ASTNode) (q *goqu.SelectDataset) {
	// get SELECT query based on
	// the given filter
	q = d.searchSql(f)

	var (
		err  error
		expr exp.Expression

		alias string

		// store map alias-expression pairs to power-up
		// HAVING clause query parsing (but only for the
		// aggregation expression!)
		a2expr = make(map[string]exp.Expression)

		selected []any

		field = func(c dal.AggregateAttr) (expr exp.Expression, err error) {
			switch {
			case len(c.RawExpr) > 0:
				// @todo could probably be removed since RawExpr is only a temporary solution?
				return d.parseQuery(c.RawExpr)
			case c.Expression != nil:
				return d.convertQuery(c.Expression)
			}

			return d.table.AttributeExpression(c.Identifier)
		}
	)

	for i, c := range groupBy {
		if expr, err = field(c); err != nil {
			return q.SetError(err)
		}

		alias = c.Identifier
		if alias == "" {
			alias = fmt.Sprintf("group_by_%d", i)
		}

		a2expr[alias] = expr

		// grouping by selected
		q = q.GroupByAppend(alias)

		expr = exp.NewAliasExpression(expr, alias)

		// Add all group-by columns at the start of selection fields
		selected = append(selected, expr)
	}

	q = q.Select(selected...)

	for i, c := range out {
		if expr, err = field(c); err != nil {
			return q.SetError(err)
		}

		alias = c.Identifier
		if alias == "" {
			alias = fmt.Sprintf("aggr_%d", i)
		}

		a2expr[alias] = expr

		expr = exp.NewAliasExpression(expr, alias)
		q = q.SelectAppend(expr)
	}

	if having != nil {
		var (
			symHandler = d.qlConverterGenericSymHandler()

			converter = ql.Converter(
				// we need a more specialized symbol handler for the HAVING clause
				// since it can contain aggregation expressions
				ql.SymHandler(func(node *ql.ASTNode) (exp.Expression, error) {
					sym := dal.NormalizeAttrNames(node.Symbol)

					if a2expr[sym] != nil {
						if d.dialect.Nuances().HavingClauseMustUseAlias {
							// is aliased expression?
							return a2expr[sym], nil
						} else {
							return exp.NewIdentifierExpression("", "", sym), nil
						}
					}

					// if not, use the default handler
					return symHandler(node)
				}),
				ql.RefHandler(d.qlConverterGenericRefHandler()),
			)
		)

		// using special symbol handler that when converting HAVING clause expression
		// this handler looks at the
		if expr, err = converter.Convert(having); err != nil {
			q.SetError(err)
			return
		}

		q = q.Having(expr)
	}

	return
}

func (d *model) lookupSql(pkv dal.ValueGetter) *goqu.SelectDataset {
	var (
		sel       = d.selectSql().Limit(1)
		cond, err = d.pkLookupCondition(pkv)
	)

	if err != nil {
		sel = sel.SetError(err)
	}

	return sel.Where(cond)
}

func (d *model) selectSql() *goqu.SelectDataset {
	var (
		cols = d.table.Columns()

		// working around a bug inside goqu lib that adds
		// * to the list of columns to be selected
		// even if we clear the columns first
		q = d.dialect.GOQU().
			From(d.table.Ident())
	)

	if len(cols) == 0 {
		return q.SetError(fmt.Errorf("can not create SELECT without columns"))
	}

	q = q.Select(d.table.Ident().Col(cols[0].Name()))
	for _, col := range cols[1:] {
		q = q.SelectAppend(d.table.Ident().Col(col.Name()))
	}

	return q
}

func (d *model) truncateSql() (_ *goqu.TruncateDataset) {
	return d.dialect.GOQU().Truncate(d.table.Ident())
}

func (d *model) insertSql(rr ...dal.ValueGetter) (_ *goqu.InsertDataset) {
	var (
		ins = d.dialect.GOQU().Insert(d.table.Ident())
		cc  = d.table.Columns()

		rows = make([][]any, len(rr))
		cols = make([]any, len(cc))

		err error
	)

	for c := range cc {
		cols[c] = cc[c].Name()
	}

	for i, r := range rr {
		rows[i], err = d.table.Encode(r)
		if err != nil {
			return ins.SetError(err)
		}
	}

	return ins.Cols(cols...).Vals(rows...)
}

// updateSql generates SQL command for updating record
func (d *model) updateSql(r dal.ValueGetter) *goqu.UpdateDataset {
	var (
		upd = d.dialect.GOQU().Update(d.table.Ident())

		values    = exp.Record{}
		condition = exp.Ex{}

		encoded, err = d.table.Encode(r)
	)

	if err != nil {
		return upd.SetError(err)
	}

	for i, c := range d.table.Columns() {
		if c.IsPrimaryKey() {
			// values[]
			condition[c.Name()] = encoded[i]
		} else {
			values[c.Name()] = encoded[i]
		}
	}

	return upd.Where(condition).Set(values)
}

func (d *model) deleteSql(pkv dal.ValueGetter) *goqu.DeleteDataset {
	var (
		del       = d.dialect.GOQU().Delete(d.table.Ident())
		cond, err = d.pkLookupCondition(pkv)
	)

	if err != nil {
		del.SetError(err)
	}

	return del.Where(cond)
}

// Constructs primary-key-lookup expression from pk values
func (d *model) pkLookupCondition(pkv dal.ValueGetter) (_ exp.Expression, err error) {
	var (
		cnd = exp.NewExpressionList(exp.AndType)
		val any
	)
	for _, c := range d.table.Columns() {
		if !c.IsPrimaryKey() {
			continue
		}

		val, err = pkv.GetValue(c.Name(), 0)
		if err != nil {
			return nil, fmt.Errorf("could not get value for primary key %q: %w", c.Name(), err)
		}

		cnd = cnd.Append(exp.NewBooleanExpression(exp.EqOp, d.table.Ident().Col(c.Name()), val))
	}

	return cnd, nil
}
