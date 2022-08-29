package dal

import (
	"fmt"
	"strings"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/ql"
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
)

func (sa SimpleAttr) Identifier() string              { return sa.Ident }
func (sa SimpleAttr) Expression() (expression string) { return sa.Expr }
func (sa SimpleAttr) Source() (ident string)          { return sa.Src }
func (sa SimpleAttr) Properties() MapProperties       { return sa.Props }

// compareGetters compares the two ValueGetters
// @todo multi-value support?
// -1: a is less then b
// 0: a is equal to b
// 1: a is greater then b
//
// This function is used to satisfy sort's less function requirement.
func compareGetters(a, b ValueGetter, ac, bc map[string]uint, attr string) int {
	va, err := a.GetValue(attr, 0)
	if err != nil {
		return 1
	}

	vb, err := b.GetValue(attr, 0)
	if err != nil {
		return 1
	}

	return compareValues(va, vb)
}

// compareValues compares the two values
// @todo support for other types and slices
// -1: a is less then b
// 0: a is equal to b
// 1: a is greater then b
//
// @note I considered using gval here but using gval proved to bring
// a bit too much overhead.
//
// @todo look into using generics or some wrapping types here
func compareValues(va, vb any) int {
	if va == vb {
		return 0
	}

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
		if xa.Equal(cb) {
			return 0
		}
		if xa.Before(cb) {
			return -1
		}
		if xa.After(cb) {
			return 1
		}
	}
	return -1
}

// // // // // // // // // // // // // // // // // // // // // // // // //

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
		parts  = make([]string, 0, 5)
		pcNode *ql.ASTNode
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

		// Convert the expression
		if expr := f.Expression(); len(expr) != 0 {
			parts = append(parts, expr)
		}

		// Convert the paging cursor
		if pc := f.Cursor(); pc != nil {
			pcNode, err = f.Cursor().ToAST(nil, nil)
			if err != nil {
				return
			}
		}
	}

	expr := strings.Join(parts, " && ")

	// Everything is empty, not doing anything
	if len(expr) == 0 && pcNode == nil {
		return nil, nil
	}

	// Parse the base expression and prepare the QL node
	if pcNode != nil {
		// Use just the paging cursor node
		if len(expr) == 0 {
			return newRunnerGvalParsed(pcNode)
		}

		// Use both the expression and the paging cursor node and-ed together
		expr, err := newConverterGval().Parse(expr)
		if err != nil {
			return nil, err
		}
		return newRunnerGvalParsed(&ql.ASTNode{
			Ref:  "and",
			Args: ql.ASTNodeSet{pcNode, expr},
		})
	}

	// Default, parse the expr from source
	return newRunnerGval(expr)
}

func valueGetterCounterComparator(ss filter.SortExprSet, a, b ValueGetter, ca, cb map[string]uint) bool {
	// @todo we can probably remove the branching here and write a bool alg. expr.
	for _, s := range ss {
		cmp := compareGetters(a, b, ca, cb, s.Column)

		if cmp != 0 {
			if s.Descending {
				return cmp > 0
			}
			return cmp < 0
		}

	}

	return false
}
