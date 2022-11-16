package dal

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/pkg/ql"
	"github.com/spf13/cast"
)

type (
	// @todo consider reworking this
	SimpleAttr struct {
		Ident string
		Expr  string
		Src   string
		Props MapProperties
	}

	// Row is a generic implementation for ValueGetter and ValueSetter
	//
	// Primarily used within DAL pipeline execution steps, but may also be used
	// outside.
	Row struct {
		counters map[string]uint
		values   valueSet

		// Metadata to make it easier to work with
		// @todo add when needed
	}

	valueSet map[string][]any
)

const (
	attributeNestingSeparator     = "."
	attributeNestingGvalSeparator = "___DLTR___"
)

func (sa SimpleAttr) Identifier() string              { return sa.Ident }
func (sa SimpleAttr) Expression() (expression string) { return sa.Expr }
func (sa SimpleAttr) Source() (ident string)          { return sa.Src }
func (sa SimpleAttr) Properties() MapProperties       { return sa.Props }

// WithValue is a simple helper to construct rows with populated values
//
// @note The main use is for tests so restrain from using it in code.
func (r *Row) WithValue(name string, pos uint, v any) *Row {
	err := r.SetValue(name, pos, v)
	if err != nil {
		panic(err)
	}

	return r
}

func (r Row) SelectGVal(ctx context.Context, k string) (interface{}, error) {
	return r.GetValue(unwrapNestedGvalIdent(k), 0)
}

// Reset clears out the row so the same instance can be reused where possible
//
// Important: Reset only clears out the counters and does not re-init/clear out
// the underlaying values. Don't directly iterate over the values, but use the
// counters.
func (r *Row) Reset() {
	for k := range r.counters {
		r.counters[k] = 0
	}
}

func (r *Row) SetValue(name string, pos uint, v any) error {
	if r.values == nil {
		r.values = make(valueSet)
	}
	if r.counters == nil {
		r.counters = make(map[string]uint)
	}

	// Make sure there is space for it
	// @note benchmarking proves that the rest of the function introduces a lot of memory pressure.
	//       Investigate options on reworking this/reducing allocations.
	if int(pos)+1 > len(r.values[name]) {
		r.values[name] = append(r.values[name], make([]any, (int(pos)+1)-len(r.values[name]))...)
	}

	r.values[name][pos] = v
	if pos >= r.counters[name] {
		r.counters[name]++
	}

	return nil
}

func (r *Row) CountValues() map[string]uint {
	return r.counters
}

func (r *Row) GetValue(name string, pos uint) (any, error) {
	if r.values == nil {
		return nil, nil
	}
	if r.counters == nil {
		return nil, nil
	}
	if pos >= r.counters[name] {
		return nil, nil
	}

	return r.values[name][pos], nil
}

func (r *Row) String() string {
	out := make([]string, 0, 20)

	for k, cc := range r.counters {
		for i := uint(0); i < cc; i++ {
			v := r.values[k][i]
			out = append(out, fmt.Sprintf("%s [%d] %v", k, i, v))
		}
	}

	return strings.Join(out, " | ")
}

// compareGetters compares the two ValueGetters
// -1: a is less then b
// 0: a is equal to b
// 1: a is greater then b
//
// Multi value rules:
// - if a has less items then b, a is less then b (-1)
// - if a has more items then b, a is more then b (1)
// - if a and b have the same amount of items; if any of the corresponding values
//   are different, that outcome is used as the result
//
// This function is used to satisfy sort's less function requirement.
func compareGetters(a, b ValueGetter, ac, bc map[string]uint, attr string) int {
	// If a has less values then b, then a is less then b
	if ac[attr] < bc[attr] {
		return -1
	} else if ac[attr] > bc[attr] {
		return 1
	}

	// If a and b have the same number of values, then we need to compare them
	for i := uint(0); i < ac[attr]; i++ {
		va, err := a.GetValue(attr, i)
		if err != nil {
			return 1
		}

		vb, err := b.GetValue(attr, i)
		if err != nil {
			return 1
		}

		// Continue the cmp. until we find two values that are different
		cmp := compareValues(va, vb)
		if cmp != 0 {
			return cmp
		}
	}

	// If any value is different from the other, the loop above would end; so
	// here, we can safely say they are the same
	return 0
}

// compareValues compares the two values
// @todo identify what other types we should support
// -1: a is less then b
// 0: a is equal to b
// 1: a is greater then b
//
// @note I considered using GVal here but it introduces more overhead then
//       what I've conjured here.
// @todo look into using generics or some wrapping types here
func compareValues(va, vb any) int {
	// simple/edge cases
	if va == vb {
		return 0
	}
	if va == nil {
		return -1
	}
	if vb == nil {
		return 1
	}

	// Compare based on type
	switch ca := va.(type) {
	case string:
		cb, err := cast.ToStringE(vb)
		if err != nil {
			return -1
		}
		if ca < cb {
			return -1
		}
		if ca > cb {
			return 1
		}

	case int, int8, int16, int32, int64:
		// this one can't error since we know it's an ok value
		xa := cast.ToInt64(va)
		cb, err := cast.ToInt64E(vb)
		if err != nil {
			return -1
		}
		if xa < cb {
			return -1
		}
		if xa > cb {
			return 1
		}

	case uint, uint8, uint16, uint32, uint64:
		// this one can't error since we know it's an ok value
		xa := cast.ToUint64(va)
		cb, err := cast.ToUint64E(vb)
		if err != nil {
			return -1
		}
		if xa < cb {
			return -1
		}
		if xa > cb {
			return 1
		}

	case float32, float64:
		// this one can't error since we know it's an ok value
		xa := cast.ToFloat64(va)
		cb, err := cast.ToFloat64E(vb)
		if err != nil {
			return -1
		}
		if xa < cb {
			return -1
		}
		if xa > cb {
			return 1
		}

	case time.Time, *time.Time:
		// this one can't error since we know it's an ok value
		xa := cast.ToTime(va)
		cb, err := cast.ToTimeE(vb)
		if err != nil {
			return -1
		}
		if xa.Before(cb) {
			return -1
		}
		if xa.After(cb) {
			return 1
		}
	}

	panic(fmt.Sprintf("unsupported type for values %v, %v", va, vb))
}

// constraintsToExpression converts the given constraints map to a ql parsable expression
func constraintsToExpression(cc map[string][]any) string {
	out := make([]string, 0, 10)

	for k, vv := range cc {
		part := make([]string, len(vv))
		for i, v := range vv {
			if vs, ok := v.(string); ok {
				part[i] = fmt.Sprintf(`%s == '%s'`, k, vs)
			} else {
				part[i] = fmt.Sprintf("%s == %v", k, v)
			}
		}

		pt := strings.Join(part, " || ")
		if len(cc) > 1 {
			out = append(out, fmt.Sprintf("(%s)", pt))
		} else {
			out = append(out, pt)
		}

	}

	return strings.Join(out, " && ")
}

// stateConstraintsToExpression converts the given state expression to a ql parsable expression
func stateConstraintsToExpression(cc map[string]filter.State) string {
	out := make([]string, 0, 10)

	for k, s := range cc {
		// Inclusive one is omitted since the condition always evaluates
		// as true (field == null || field != null => true)
		switch s {
		case filter.StateExcluded:
			out = append(out, fmt.Sprintf("%s == null", k))

		// Only these ones
		case filter.StateExclusive:
			out = append(out, fmt.Sprintf("%s != null", k))
		}
	}

	return strings.Join(out, " && ")
}

// @todo see if the rest of the "conversion" functions should return a QL node
//       like the cursor one does.
func prepareGenericRowTester(f internalFilter) (_ tester, err error) {
	var (
		parts    = make([]string, 0, 5)
		pcNode   *ql.ASTNode
		exprNode *ql.ASTNode
	)

	{
		// Convert the regular constraints
		if cc := f.Constraints(); len(cc) != 0 {
			parts = append(parts, constraintsToExpression(cc))
		}

		// Convert state constraints
		// @todo check if the attributes in the state constraints are nullable.
		if sc := f.StateConstraints(); len(sc) != 0 {
			parts = append(parts, stateConstraintsToExpression(sc))
		}

		exprNode = f.ExpressionParsed()
		if exprNode == nil {
			expr := f.Expression()
			if expr != "" {
				parts = append(parts, f.Expression())
			}
		}

		// Convert the paging cursor
		if pc := f.Cursor(); pc != nil {
			pcNode, err = f.Cursor().ToAST(nil, nil)
			if err != nil {
				return
			}
			pcNode.Traverse(func(a *ql.ASTNode) (bool, *ql.ASTNode, error) {
				if a.Symbol != "" {
					a.Symbol = wrapNestedGvalIdent(a.Symbol)
				}
				return true, a, nil
			})
		}
	}

	expr := strings.Join(parts, " && ")

	// Everything is empty, not doing anything
	if len(expr) == 0 && exprNode == nil && pcNode == nil {
		return nil, nil
	}

	args := make([]*ql.ASTNode, 0, 5)

	// Paging cursors
	if pcNode != nil {
		args = append(args, pcNode)
	}

	// Parsed filter expression
	if exprNode != nil {
		args = append(args, exprNode)
	}

	// Rest of the generated expression string
	if len(expr) > 0 {
		expr, err := newConverterGval().Parse(expr)
		if err != nil {
			return nil, err
		}
		args = append(args, expr)
	}

	return newRunnerGvalParsed(
		&ql.ASTNode{
			Ref:  "and",
			Args: args,
		},
	)
}

// makeRowComparator returns a ValueGetter comparator for the given sort expr
func makeRowComparator(ss ...*filter.SortExpr) func(a, b ValueGetter) bool {
	return func(a, b ValueGetter) bool {
		for _, s := range ss {
			cmp := compareGetters(a, b, a.CountValues(), b.CountValues(), s.Column)

			less, skip := evalCmpResult(cmp, s)
			if !skip {
				return less
			}
		}

		return false
	}
}

func evalCmpResult(cmp int, s *filter.SortExpr) (less, skip bool) {
	if cmp != 0 {
		if s.Descending {
			return cmp > 0, false
		}
		return cmp < 0, false
	}

	return false, true
}

func indexAttrs(aa ...AttributeMapping) (out map[string]bool) {
	out = make(map[string]bool, len(aa))
	indexAttrsInto(out, aa...)
	return
}

func indexAttrsInto(dst map[string]bool, aa ...AttributeMapping) {
	for _, a := range aa {
		dst[a.Identifier()] = true
	}
}

// keysFromExpr returns all of the identifiers used in agg. group expressions
//
// The hasConstants return argument is true if any of the expressions returns a
// constant value, such as year(now()) or 42
func keysFromExpr(nn ...*ql.ASTNode) (out []string, hasConstants bool) {
	out = make([]string, 0, (len(nn)+1)*2)
	auxOut := make(map[string]bool, (len(nn)+1)*2)

	for _, n := range nn {
		symbols := n.CollectSymbols()
		if len(symbols) == 0 {
			hasConstants = true
		}

		for _, s := range symbols {
			auxOut[s] = true
		}
	}

	for k := range auxOut {
		out = append(out, k)
	}

	return
}

// Assure sort validates that the filter's definition includes all of the primary
// keys and that the paging cursor's sort is compatible
func assureSort(f internalFilter, primaries []string) (out internalFilter, err error) {
	out = f

	// make sure all primary keys are in there
	for _, p := range primaries {
		if out.orderBy.Get(p) == nil {
			out.orderBy = append(out.orderBy, &filter.SortExpr{
				Column:     p,
				Descending: out.orderBy.LastDescending(),
			})
		}
	}

	// No cursor, no problem
	if f.cursor == nil {
		return
	}

	// Make sure the cursor can handle this sort def
	out.orderBy, err = out.cursor.Sort(out.orderBy)
	if err != nil {
		return
	}

	return
}

func wrapNestedGvalIdent(ident string) string {
	return strings.ReplaceAll(ident, attributeNestingSeparator, attributeNestingGvalSeparator)
}

func unwrapNestedGvalIdent(ident string) string {
	return strings.ReplaceAll(ident, attributeNestingGvalSeparator, attributeNestingSeparator)
}
