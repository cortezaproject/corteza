package rdbms

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
)

type (
	cursorCondition struct {
		cur          pagingCursor
		keyMapper    cursorKeyMapper
		sortableCols map[string]string
	}

	pagingCursor interface {
		Keys() []string
		KK() [][]string
		Modifiers() []string
		Values() []interface{}
		IsLThen() bool
		Desc() []bool
		IsROrder() bool
	}

	KeyMap struct {
		FieldCast    string
		TypeCast     string
		TypeCastPtrn string
	}

	// translates field to (store) column
	// this helps us with upper/lower case, underscore/camel-case and
	// when using fields on records
	//
	// @todo extend the return args to provide additional info (like is-nullable)
	//       to avoid IS NULL checks (see sql() fn)
	cursorKeyMapper func(string) (KeyMap, error)
)

func cursor(cursor *filter.PagingCursor) ([]goqu.Expression, error) {
	sql, args, err := CursorCondition(cursor, nil, nil).ToSQL()
	if err != nil {
		return nil, err
	}

	return []goqu.Expression{goqu.Literal(sql, args...)}, nil
}

func cursorWithSorting(cursor *filter.PagingCursor, sortableCols map[string]string) ([]goqu.Expression, error) {
	sql, args, err := CursorCondition(cursor, nil, sortableCols).ToSQL()
	if err != nil {
		return nil, err
	}

	return []goqu.Expression{goqu.Literal(sql, args...)}, nil
}

// CursorCondition builds a complex condition to filter rows before/after row that
// the paging cursor points to
func CursorCondition(pc pagingCursor, keyMapper cursorKeyMapper, sortableCols map[string]string) *cursorCondition {
	if keyMapper == nil {
		keyMapper = func(s string) (KeyMap, error) {
			return KeyMap{
				FieldCast:    s,
				TypeCast:     s,
				TypeCastPtrn: "%s",
			}, nil
		}
	}

	return &cursorCondition{cur: pc, keyMapper: keyMapper, sortableCols: sortableCols}
}

func (c *cursorCondition) ToSQL() (string, []interface{}, error) {
	sql, err := c.sql()
	if err != nil {
		return "", nil, err
	}

	return sql, c.values(), nil
}

// generates all values that we need for the generated condition SQL
func (c *cursorCondition) values() []interface{} {
	if len(c.cur.Keys()) == 0 {
		return nil
	}

	l := len(c.cur.Values())
	vv := make([]interface{}, 0, l*2-1)
	for i, v := range c.cur.Values() {
		vv = append(vv, v)
		if i < l-1 {
			// for all but 1st, use values 2 times.
			vv = append(vv, v)
		}
	}

	return vv
}

// builds cursor SQL expression
//
// this could be simple (f1, f2, ...) < (v1, v2, ...) but we need  to be a bit careful with NULL values
// So we need (f1 < v1 OR (f1 = v1 AND f2 < v2) pattern, extended to:
// ((f1 IS NULL AND v1 IS NOT NULL) OR f1 < v1 OR (((f1 IS NULL AND v1 IS NULL) OR f1 = v1) AND (f2...)
//
// Due to issues with param biding & types in Postgres (using ? IS NULL results in an error), we need do
// check (on app-side) if value is nil to replace "? IS (NOT) NULL" check with TRUE/FALSE constants.
func (c *cursorCondition) sql() (cnd string, err error) {
	const (
		// we start with this
		baseTpl = "((%s IS %s AND %s) OR (%s %s %s))"

		// and then wrap each iteration with base and this
		wrapTpl = "(%s OR (((%s IS NULL AND %s) OR %s = %s) AND %s))"
	)

	var (
		cc = c.cur.Keys()
		kk = c.cur.KK()
		mm = c.cur.Modifiers()
		vv = c.cur.Values()

		ltOp = map[bool]string{
			true:  "<",
			false: ">",
		}

		notOp = map[bool]string{
			true:  "NOT NULL",
			false: "NULL",
		}

		// Little utility to know for sure if some value is nil or not
		//
		// Interface variables can be a bit tricky here, so this is required.
		nilCheck = func(i interface{}) bool {
			if i == nil {
				return true
			}
			switch reflect.TypeOf(i).Kind() {
			case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
				return reflect.ValueOf(i).IsNil()
			}
			return false
		}

		// Modifying this function to use expressions instead of constant boolean
		// values because MSSQL doesn't have those.
		//
		// @todo rethink and redo the whole/all of the filtering logic surrounding paging
		// cursors to make them consistent/reusable
		isNull = func(i int, neg bool) string {
			if (nilCheck(vv[i]) && !neg) || (!nilCheck(vv[i]) && neg) {
				return "1=1"
			}

			return "1=0"
		}
	)

	if len(cc) == 0 {
		return
	}

	// going from the last key/column to the 1st one
	for i := len(cc) - 1; i >= 0; i-- {
		// Get the key context so we know how to format fields and format typecasts
		colName := cc[i]
		if c.sortableCols != nil {
			if v, ok := c.sortableCols[strings.ToLower(cc[i])]; ok {
				colName = v
			}

			// update column name as per the modifier if needed
			if len(kk[i]) > 0 && strings.ToLower(mm[i]) == filter.COALESCE {
				var tmp []string
				for _, k := range kk[i] {
					if v, ok := c.sortableCols[strings.ToLower(k)]; ok {
						tmp = append(tmp, v)
					}
				}
				colName = fmt.Sprintf("COALESCE(%s)", strings.Join(tmp, ", "))
			}
		}

		km, err := c.keyMapper(colName)
		if err != nil {
			return "", err
		}

		// We need to cut off the values that are before the cursor (when ascending)
		// and vice-versa for descending.
		lt := c.cur.Desc()[i]
		if c.cur.IsROrder() {
			lt = !lt
		}
		op := ltOp[lt]

		// Typecast the value so comparasement can work properly
		vc := fmt.Sprintf(km.TypeCastPtrn, "?")

		base := fmt.Sprintf(baseTpl, km.FieldCast, notOp[!lt], isNull(i, lt), km.TypeCast, op, vc)

		if cnd == "" {
			cnd = base
		} else {
			// wrap existing conditions (next key) and the generated base for the current key
			cnd = fmt.Sprintf(wrapTpl, base, km.FieldCast, isNull(i, false), km.TypeCast, vc, cnd)
		}
	}

	return
}

// CursorExpression builds cursor SQL expression using goqu/exp package
//
// this could be simple (f1, f2, ...) < (v1, v2, ...) but we need  to be a bit careful with NULL values
// So we need (f1 < v1 OR (f1 = v1 AND f2 < v2) pattern, extended to:
// ((f1 IS NULL AND v1 IS NOT NULL) OR f1 < v1 OR (((f1 IS NULL AND v1 IS NULL) OR f1 = v1) AND (f2...)
//
// Due to issues with param biding & types in Postgres (using ? IS NULL results in an error), we need do
// check (on app-side) if value is nil to replace "? IS (NOT) NULL" check with TRUE/FALSE constants.
func CursorExpression(
	cur *filter.PagingCursor,
	identLookup func(i string) (exp.Expression, error),
	castFn func(i string, val any) (exp.Expression, error),
) (e exp.Expression, err error) {
	var (
		cc = cur.Keys()
		vv = cur.Values()

		value any

		ident exp.Expression

		ltOp = map[bool]exp.BooleanOperation{
			true:  exp.LtOp,
			false: exp.GtOp,
		}

		notOp = map[bool]exp.LiteralExpression{
			true:  exp.NewLiteralExpression("NOT NULL"),
			false: exp.NewLiteralExpression("NULL"),
		}

		// Little utility to know for sure if some value is nil or not
		//
		// Interface variables can be a bit tricky here, so this is required.
		nilCheck = func(i interface{}) bool {
			if i == nil {
				return true
			}
			switch reflect.TypeOf(i).Kind() {
			case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
				return reflect.ValueOf(i).IsNil()
			}
			return false
		}

		// Modifying this function to use expressions instead of constant boolean
		// values because MSSQL doesn't have those.
		//
		// @todo rethink and redo the whole/all of the filtering logic surrounding paging
		// cursors to make them consistent/reusable
		isValueNull = func(i int, neg bool) exp.Expression {
			if (nilCheck(vv[i]) && !neg) || (!nilCheck(vv[i]) && neg) {
				return exp.NewLiteralExpression("1=1")
			}

			return exp.NewLiteralExpression("1=0")
		}

		curCond exp.ExpressionList
	)

	if len(cc) == 0 {
		return
	}

	// ((($1 IS NOT NULL AND FALSE) OR ($2 > $3)) OR (((? IS NULL AND FALSE) OR ? = ?) AND (("test_tbl"."id" IS NOT NULL AND FALSE) OR ("test_tbl"."id" > ?))))

	// going from the last key/column to the 1st one
	for i := len(cc) - 1; i >= 0; i-- {
		if identLookup != nil {
			// Get the key context so we know how to format fields and format typecasts
			ident, err = identLookup(cc[i])
			if err != nil {
				return
			}
		} else {
			ident = exp.NewLiteralExpression("?", exp.NewIdentifierExpression("", "", cc[i]))
		}

		if castFn == nil {
			value = vv[i]
		} else {
			value, err = castFn(cc[i], vv[i])
			if err != nil {
				return
			}
		}

		// We need to cut off the values that are before the cursor (when ascending)
		// and vice-versa for descending.
		lt := cur.Desc()[i]
		if cur.IsROrder() {
			lt = !lt
		}

		op := ltOp[lt]

		// // Typecast the value so comparison can work properly

		// Either BOTH (field and value) are NULL or field is grater-then value
		base := exp.NewExpressionList(
			exp.OrType,
			// both NULL
			exp.NewExpressionList(
				exp.AndType,
				exp.NewLiteralExpression(`(? IS ?)`, ident, notOp[!lt]),
				isValueNull(i, lt),
			),
			// or GT/LT value
			exp.NewBooleanExpression(op, ident, value),
		)

		if curCond == nil {
			curCond = base
		} else {
			curCond = exp.NewExpressionList(
				exp.OrType,
				base,
				exp.NewExpressionList(
					exp.AndType,
					exp.NewExpressionList(
						exp.OrType,
						// both NULL
						exp.NewExpressionList(
							exp.AndType,
							exp.NewLiteralExpression(`(? IS NULL)`, ident),
							isValueNull(i, false),
						),
						exp.NewBooleanExpression(exp.EqOp, ident, value),
					),
					curCond,
				),
			)
		}
	}

	return curCond.Expression(), nil
}
