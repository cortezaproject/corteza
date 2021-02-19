package builders

import (
	"fmt"
	"reflect"

	"github.com/davecgh/go-spew/spew"
)

type (
	cursorCondition struct {
		cur       pagingCursor
		keyMapper cursorKeyMapper
	}

	pagingCursor interface {
		Keys() []string
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

// Builds a complex condition to filter rows before/after row that
// the paging cursor points to
func CursorCondition(pc pagingCursor, keyMapper cursorKeyMapper) *cursorCondition {
	if keyMapper == nil {
		keyMapper = func(s string) (KeyMap, error) {
			return KeyMap{
				FieldCast:    s,
				TypeCast:     s,
				TypeCastPtrn: "%s",
			}, nil
		}
	}

	return &cursorCondition{cur: pc, keyMapper: keyMapper}
}

func (c *cursorCondition) ToSql() (string, []interface{}, error) {
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

// builds cursor SQL condition
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

		isNull = func(i int, neg bool) string {
			if (nilCheck(vv[i]) && !neg) || (!nilCheck(vv[i]) && neg) {
				return "TRUE"
			}

			return "FALSE"
		}
	)

	if len(cc) == 0 {
		return
	}

	// going from the last key/column to the 1st one
	for i := len(cc) - 1; i >= 0; i-- {
		// Get the key context so we know how to format fields and format typecasts
		km, err := c.keyMapper(cc[i])
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

		spew.Dump("lt", lt)
		base := fmt.Sprintf(baseTpl, km.FieldCast, notOp[!lt], isNull(i, lt), km.TypeCast, op, vc)

		if cnd == "" {
			cnd = base
		} else {
			// wrap existing conditions (next key) and the generated base for the current key
			cnd = fmt.Sprintf(wrapTpl, base, km.FieldCast, isNull(i, false), km.TypeCast, vc, cnd)
		}
	}

	spew.Dump(cnd)
	return
}
