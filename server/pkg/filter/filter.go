package filter

import "github.com/cortezaproject/corteza-server/pkg/ql"

type (
	filterOpt func(*filter)

	filter struct {
		constaints       map[string][]any
		stateConditions  map[string]State
		metaConditions   map[string]any
		expression       string
		expressionParsed *ql.ASTNode
		orderBy          SortExprSet
		limit            uint
		cursor           *PagingCursor
	}

	Filter interface {
		// Constraints returns map of attribute idents and values
		// used for structured filtering ({a1: [v1], a2: [v2, v3]} => "a1 = v1 AND a2 = (v2,v4)")
		Constraints() map[string][]any

		// StateConstraints returns map of attribute idents and states
		// used for structured filtering ({a1: s1, a2: s2} => "a1 = s1 AND a2 = s2")
		StateConstraints() map[string]State

		// MetaConstraints returns meta constraints
		MetaConstraints() map[string]any

		// Expression returns string, parseable by ql package
		Expression() string

		// OrderBy one or more fields
		OrderBy() SortExprSet

		// Limit amount of returned results
		Limit() uint

		// Cursor from the last fetch
		Cursor() *PagingCursor
	}
)

func Generic(oo ...filterOpt) *filter {
	f := &filter{}

	for _, o := range oo {
		o(f)
	}

	return f
}

// WithConstraint sets single constraint to filter
func WithConstraint(k string, v any) filterOpt {
	return func(f *filter) {
		if f.constaints == nil {
			f.constaints = make(map[string][]any)
		}

		f.constaints[k] = append(f.constaints[k], v)
	}
}

// WithConstraints set multiple constraints to filter
func WithConstraints(c map[string][]any) filterOpt {
	return func(f *filter) {
		f.constaints = c
	}
}

// WithStateConstraint sets single state constraint to filter
func WithStateConstraint(k string, s State) filterOpt {
	return func(f *filter) {
		if f.stateConditions == nil {
			f.stateConditions = make(map[string]State)
		}

		f.stateConditions[k] = s
	}
}

// WithStateConstraints sets multiple state constraints to filter
func WithStateConstraints(sc map[string]State) filterOpt {
	return func(f *filter) {
		f.stateConditions = sc
	}
}

// WithMetaConstraints sets multiple state constraints to filter
func WithMetaConstraints(mc map[string]any) filterOpt {
	return func(f *filter) {
		f.metaConditions = mc
	}
}

// WithExpression sets expression to filter
func WithExpression(e string) filterOpt {
	return func(f *filter) {
		f.expression = e
	}
}

// WithExpressionParsed sets parsed expression to filter
func WithExpressionParsed(e *ql.ASTNode) filterOpt {
	return func(f *filter) {
		f.expressionParsed = e
	}
}

// WithOrderBy sets order by expression
func WithOrderBy(o SortExprSet) filterOpt {
	return func(f *filter) {
		f.orderBy = o
	}
}

// WithLimit sets limit to filter
func WithLimit(l uint) filterOpt {
	return func(f *filter) {
		f.limit = l
	}
}

// WithCursor sets cursor to filter
func WithCursor(p *PagingCursor) filterOpt {
	return func(f *filter) {
		f.cursor = p
	}
}

func (f *filter) With(oo ...filterOpt) *filter {
	for _, o := range oo {
		o(f)
	}

	return f
}

func (f *filter) Constraints() map[string][]any      { return f.constaints }
func (f *filter) StateConstraints() map[string]State { return f.stateConditions }
func (f *filter) MetaConstraints() map[string]any    { return f.metaConditions }
func (f *filter) Expression() string                 { return f.expression }
func (f *filter) ExpressionParsed() *ql.ASTNode      { return f.expressionParsed }
func (f *filter) OrderBy() SortExprSet               { return f.orderBy }
func (f *filter) Limit() uint                        { return f.limit }
func (f *filter) Cursor() *PagingCursor              { return f.cursor }
