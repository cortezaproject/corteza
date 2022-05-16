package filter

type (
	Filter interface {
		// Constraints returns map of attribute idents and values
		// used for structured filtering ({a1: [v1], a2: [v2, v3]} => "a1 = v1 AND a2 = (v2,v4)")
		Constraints() map[string][]any

		// StateConstraints returns map of attribute idents and states
		// used for structured filtering ({a1: s1, a2: s2} => "a1 = s1 AND a2 = s2")
		StateConstraints() map[string]State

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
