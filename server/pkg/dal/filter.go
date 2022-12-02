package dal

import (
	"fmt"

	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/pkg/ql"
)

type (
	// internalFilter is a wrap to simplify interactions for cases where we might
	// want to interact with the received filters; i.e. when optimizing pipeline steps
	internalFilter struct {
		constraints      map[string][]any
		stateConstraints map[string]filter.State
		metaConstraints  map[string]any
		expression       string
		expParsed        *ql.ASTNode
		orderBy          filter.SortExprSet
		limit            uint
		cursor           *filter.PagingCursor
	}

	parsedFilter interface {
		ExpressionParsed() *ql.ASTNode
	}
)

func (f internalFilter) Constraints() map[string][]any             { return f.constraints }
func (f internalFilter) StateConstraints() map[string]filter.State { return f.stateConstraints }
func (f internalFilter) MetaConstraints() map[string]any           { return f.metaConstraints }
func (f internalFilter) Expression() string                        { return f.expression }
func (f internalFilter) ExpressionParsed() *ql.ASTNode             { return f.expParsed }
func (f internalFilter) OrderBy() filter.SortExprSet               { return f.orderBy }
func (f internalFilter) Limit() uint                               { return f.limit }
func (f internalFilter) Cursor() *filter.PagingCursor              { return f.cursor }

// toInternalFilter converts filter.Filter to internalFilter for easier manipulation
//
// nil filters are returned as zeroed-out internalFilters.
func toInternalFilter(f filter.Filter) (out internalFilter, err error) {
	if f == nil {
		return
	}

	// Base
	out = internalFilter{
		constraints:      f.Constraints(),
		stateConstraints: f.StateConstraints(),
		expression:       f.Expression(),
		orderBy:          f.OrderBy(),
		limit:            f.Limit(),
		cursor:           f.Cursor(),
	}

	pf, ok := f.(parsedFilter)
	if ok {
		out.expParsed = pf.ExpressionParsed()

		// In case the filter was already provided, we need to make sure the idents
		// are wrapped
		out.expParsed.Traverse(func(a *ql.ASTNode) (bool, *ql.ASTNode, error) {
			a.Symbol = wrapNestedGvalIdent(a.Symbol)
			return true, a, nil
		})
	}
	if out.expParsed != nil {
		// We can't mix the two so just clear it out to avoid confusion
		out.expression = ""
	}

	// Parse expression for later use
	if out.expression != "" {
		out.expParsed, err = newQlParser().Parse(out.expression)
		if err != nil {
			return
		}
	}

	return
}

func FilterFromExpr(n *ql.ASTNode) internalFilter {
	// @todo consider adding string expr for consistency
	return internalFilter{
		expParsed: n,
	}
}

// MergeFilters returns a new filter based on a overwritten by b
func (a internalFilter) MergeFilters(b filter.Filter) (c internalFilter, err error) {
	// In case we got a generic b filter, convert it to the internal one for easier handling
	aux, ok := b.(internalFilter)
	if !ok {
		var err error
		aux, err = toInternalFilter(b)
		if err != nil {
			return c, err
		}
	}

	return a.mergeFilters(aux), nil
}

func (a internalFilter) mergeFilters(b internalFilter) (c internalFilter) {

	c = a

	// expression
	if b.expParsed != nil {
		if c.expParsed == nil {
			c.expParsed = b.expParsed
			c.expression = b.expression
		} else {
			parsedA := a.expParsed
			rawA := a.expression
			if parsedA.Ref != "group" {
				parsedA = &ql.ASTNode{
					Ref: "group",
					Args: []*ql.ASTNode{
						parsedA,
					},
				}
				rawA = fmt.Sprintf("(%s)", rawA)
			}

			parsedB := b.expParsed
			rawB := b.expression
			if parsedB.Ref != "group" {
				parsedB = &ql.ASTNode{
					Ref: "group",
					Args: []*ql.ASTNode{
						parsedB,
					},
				}
				rawB = fmt.Sprintf("(%s)", rawB)
			}

			c.expParsed = &ql.ASTNode{
				Ref:  "and",
				Args: ql.ASTNodeSet{parsedA, parsedB},
			}
			c.expression = fmt.Sprintf("%s && %s", rawA, rawB)
		}

	}

	// constraints
	if b.constraints != nil {
		if c.constraints == nil {
			c.constraints = make(map[string][]any)
		}

		for k, v := range b.constraints {
			c.constraints[k] = append(c.constraints[k], v...)
		}
	}

	// state constraints
	if b.stateConstraints != nil {
		if c.stateConstraints == nil {
			c.stateConstraints = make(map[string]filter.State)
		}
		for k, v := range b.stateConstraints {
			c.stateConstraints[k] = v
		}
	}

	// order by
	if b.orderBy != nil {
		c.orderBy = append(c.orderBy, b.orderBy...)
	}

	// cursor
	// always use the latest paging cursor
	if b.cursor != nil {
		c.cursor = b.cursor
	}

	return
}

// empty reports if the filter is completely empty or not
func (f internalFilter) empty() bool {
	return f.constraints == nil &&
		f.stateConstraints == nil &&
		f.expression == "" &&
		f.orderBy == nil &&
		f.limit == 0 &&
		f.cursor == nil
}
