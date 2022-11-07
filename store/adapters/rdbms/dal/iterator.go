package dal

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/cortezaproject/corteza-server/pkg/errors"

	"github.com/cortezaproject/corteza-server/pkg/dal"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
)

type (
	iterator struct {
		// source model; how data we are reading from is shaped
		src *model

		// destination model; how data we are reading into is shaped
		// this is used to create scan buffer
		// when not doing plain selection from one table final results might
		// require a different list of columns and scanning needs to be adjusted
		dst *model

		// buffer for scanned rows
		scanBuf []any

		// results from the last read
		rows *sql.Rows

		// last error
		err error

		query   *goqu.SelectDataset
		sorting filter.SortExprSet
		cursor  *filter.PagingCursor
		limit   uint

		// set this flag when iterating over aggregated results
		// it will affect how
		//  - cursor conditions are applied (WHERE )
		aggregated bool
	}
)

func (i *iterator) Next(ctx context.Context) bool {
	if i.err == nil && i.rows == nil {
		i.rows, i.err = i.fetch(ctx)
	}

	if i.err != nil {
		return false
	}

	return i.rows.Next()
}

// More fetches more records from the point of last record
func (i *iterator) More(max uint, last dal.ValueGetter) (err error) {
	if i.rows != nil {
		if err = i.rows.Close(); err != nil {
			return fmt.Errorf("could not close previous query: %w", err)
		}
		i.rows = nil
	}

	i.limit = max
	if last != nil {
		if i.cursor, err = i.collectCursorValues(last); err != nil {
			return fmt.Errorf("could not collect cursor values: %w", err)
		}
	}

	return nil
}

func (i *iterator) Preload(_ context.Context, max uint, cur *filter.PagingCursor) (err error) {
	if i.rows != nil {
		if err = i.rows.Close(); err != nil {
			return fmt.Errorf("could not close previous query: %w", err)
		}
		i.rows = nil
	}

	i.limit = max
	if cur != nil {
		i.cursor = cur
	}

	return i.err
}

func (i *iterator) Sorting() filter.SortExprSet {
	return i.sorting
}

func (i *iterator) fetch(ctx context.Context) (rows *sql.Rows, err error) {
	if i.err != nil {
		return nil, i.err
	}

	if i.query == nil {
		return nil, fmt.Errorf("can not fetch without query")
	}

	if i.scanBuf == nil {
		// we're going to init scan buffer only once
		// and rely on the query.Rows.Scan function to
		// fill it up with fresh values!
		i.scanBuf = i.dst.table.MakeScanBuffer()
	}

	var (
		cur = i.cursor

		tmp   exp.Expression
		query string
		args  []any

		// contains query with ORDER BY and WHERE clauses
		sqlExpr = i.query

		sort = i.sorting.Clone()
	)

	{
		// Apply limit from the filter
		sqlExpr = sqlExpr.Limit(i.limit)

		if cur != nil {
			// @todo this needs to work with embedded attributes (non physical columns) as well!
			tmp, err = rdbms.CursorExpression(
				cur,
				func(ident string) (exp.LiteralExpression, error) { return i.src.table.AttributeExpression(ident) },
				func(ident string, val any) (exp.LiteralExpression, error) {
					// @fixme vvv
					// attr := i.dst.model.Attributes.FindByStoreIdent(ident)

					attr := i.dst.model.Attributes.FindByIdent(ident)
					if attr == nil {
						panic("unknown attribute " + ident + " used in cursor expression cast callback")
					}

					enc, err := i.dst.dialect.TypeWrap(attr.Type).Encode(val)
					if err != nil {
						return nil, err
					}

					return exp.NewLiteralExpression("?", enc), nil
				},
			)

			if err != nil {
				return
			}

			if i.aggregated {
				sqlExpr = sqlExpr.Having(tmp)
			} else {
				sqlExpr = sqlExpr.Where(tmp)
			}

			if cur.IsROrder() {
				if i.limit > 0 {
					// When paging with the reverse cursor AND limit set
					// we need to do a do sub-query reverse sort to ensure
					// that we only select the rows that make sense

					innerSort := sort.Clone()
					innerSort.Reverse()

					// Wrap the fil & ordered sub-query with cursor-conditions
					sqlExpr = i.src.dialect.GOQU().From(sqlExpr.Order(i.orderByExp(innerSort)...).As(i.src.model.Ident))

					// make sure we don't reverse it again
				} else {
					// if limit is not set it does not make sense to wrap the select
					// and resort it, so let's reverse the main sorting
					sort.Reverse()
				}
			}
		}
	}

	{
		// Apply sort

		// @todo is this going to be a problem? do we need to properly address the columns
		//       from the sub-query?
		if len(sort) > 0 {
			sqlExpr = sqlExpr.Order(i.orderByExp(sort)...)
		}
	}

	if query, args, err = sqlExpr.ToSQL(); err != nil {
		return nil, err
	}

	rows, err = i.src.conn.QueryContext(ctx, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		// no rows, no error
		return nil, nil
	}

	return
}

// generates slice of ordered-expressions
func (i *iterator) orderByExp(sort filter.SortExprSet) (oe []exp.OrderedExpression) {
	for _, s := range sort {
		// assuming all columns were pre-validated!
		tmp, _ := i.src.table.AttributeExpression(s.Column)

		if s.Descending {
			oe = append(oe, i.src.dialect.OrderedExpression(tmp, exp.DescSortDir, exp.NullsLastSortType))
		} else {
			oe = append(oe, i.src.dialect.OrderedExpression(tmp, exp.AscDir, exp.NullsFirstSortType))
		}
	}

	return
}

func (i *iterator) Scan(r dal.ValueSetter) (err error) {
	if i.err != nil {
		return i.err
	}

	if err = i.rows.Scan(i.scanBuf...); err != nil {
		return err
	}

	if err = i.dst.table.Decode(i.scanBuf, r); err != nil {
		return
	}

	return nil
}

func (i *iterator) Err() error {
	return i.err
}

// Close iterator and cleanup
func (i *iterator) Close() error {
	if i.rows == nil {
		return nil
	}
	return i.rows.Close()
}

func (i *iterator) BackCursor(r dal.ValueGetter) (cur *filter.PagingCursor, err error) {
	cur, err = i.collectCursorValues(r)
	if err != nil {
		return
	}

	// if this cursor is used, we need to reverse the sorting order
	cur.ROrder = true

	// if this cursor is used, we need to flip the conditional operator
	// from less-then to greater-then
	cur.LThen = i.sorting.Reversed()
	return
}

func (i *iterator) ForwardCursor(r dal.ValueGetter) (*filter.PagingCursor, error) {
	return i.collectCursorValues(r)
}

func (i *iterator) collectCursorValues(r dal.ValueGetter) (_ *filter.PagingCursor, err error) {
	var (
		cur = &filter.PagingCursor{LThen: i.sorting.Reversed()}

		// @todo this will not work when using multiple primary keys!
		pkUsed bool
		value  any

		pKeys = make(map[string]bool)
	)

	for _, c := range i.dst.table.Columns() {
		if c.IsPrimaryKey() {
			attrIdent := c.Attribute().Ident
			pKeys[attrIdent] = true
		}
	}

	if len(pKeys) == 0 {
		//return nil, fmt.Errorf("can not construct cursor without primary key attributes")
	}

	for _, c := range i.sorting {
		if pKeys[c.Column] {
			pkUsed = true
		}

		if value, err = r.GetValue(c.Column, 0); err != nil {
			return
		}

		cur.Set(c.Column, value, c.Descending)
	}

	if !pkUsed {
		for key := range pKeys {
			value, err = r.GetValue(key, 0)
			if err != nil {
				return
			}

			cur.Set(key, value, false)
		}
	}

	return cur, nil
}
