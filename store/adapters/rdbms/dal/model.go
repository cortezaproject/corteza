package dal

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/dal"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/drivers"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/ql"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/jmoiron/sqlx"
)

type (
	queryRunner interface {
		sqlx.QueryerContext
		sqlx.ExecerContext
	}

	queryParser interface {
		Parse(string) (exp.Expression, error)
	}

	model struct {
		model *dal.Model
		conn  queryRunner

		queryParser queryParser
		dialect     drivers.Dialect

		table drivers.TableCodec
	}
)

// Model returns fully initialized model store
//
// It abstracts database table and its columns and provides unified interface
// for fetching and storing records.
func Model(m *dal.Model, c queryRunner, d drivers.Dialect) *model {
	var (
		ms = &model{
			model:       m,
			conn:        c,
			dialect:     d,
			queryParser: ql.Converter(),
			table:       drivers.NewTableCodec(m, d),
		}
	)

	ms.queryParser = ql.Converter(
		ql.SymHandler(func(node *ql.ASTNode) (exp.Expression, error) {
			attr := ms.model.Attributes.FindByIdent(node.Symbol)
			if attr == nil {
				return nil, fmt.Errorf("unknown attribute %q used in query expression", node.Symbol)
			}

			if !attr.Filterable {
				return nil, fmt.Errorf("attribute %q can not be used in query expression", attr.Ident)
			}

			return ms.table.AttributeExpression(node.Symbol)
		}),
	)

	return ms
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

	for _, s := range orderBy {
		if _, err = d.table.AttributeExpression(s.Column); err != nil {
			return nil, err
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

	var (
		q *goqu.SelectDataset
	)

	q = d.searchSql(f)
	if err = q.Error(); err != nil {
		return
	}

	return &iterator{
		ms:      d,
		query:   q,
		sorting: orderBy,
		cursor:  f.Cursor(),
		limit:   f.Limit(),
	}, nil
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
		return sql.ErrNoRows
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

		//if d.sysExprNamespaceID != nil {
		//	cnd = append(cnd, d.sysExprNamespaceID.Eq(f.NamespaceID))
		//} else {
		//	// @todo check if f.NamespaceID is compatible
		//}
		//
		//if d.sysExprModuleID != nil {
		//	cnd = append(cnd, d.sysExprModuleID.Eq(f.ModuleID))
		//} else {
		//	// @todo check if f.ModuleID is compatible
		//}
	}

	for ident, vv := range f.Constraints() {
		attr := d.model.Attributes.FindByIdent(ident)
		if attr == nil {
			return base.SetError(fmt.Errorf("unknown attribute %q used for constrant", ident))
		}

		// @note why?
		// if !attr.PrimaryKey {
		// 	continue
		// }

		var attrExpr exp.LiteralExpression
		attrExpr, err = d.table.AttributeExpression(attr.Ident)
		if err != nil {
			return base.SetError(err)
		}

		if len(vv) > 0 {
			cnd = append(cnd, attrExpr.In(vv...))
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

		var attrExpr exp.LiteralExpression
		attrExpr, err = d.table.AttributeExpression(attr.Ident)
		if err != nil {
			return base.SetError(err)
		}

		switch state {
		case filter.StateExclusive:
			// only not-null values
			cnd = append(cnd, attrExpr.IsNotNull())

		case filter.StateExcluded:
			// exclude all non-null values
			cnd = append(cnd, attrExpr.IsNull())
		}
	}

	if q := strings.TrimSpace(f.Expression()); len(q) > 0 {
		if tmp, err = d.queryParser.Parse(q); err != nil {
			return base.SetError(err)
		}

		cnd = append(cnd, tmp)
	}

	return base.Where(cnd...)
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
			From(d.table.Ident()).
			Select(d.table.Ident().Col(cols[0].Name()))
	)

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
