package filter

import (
	"github.com/Masterminds/squirrel"
)

type (
	// State for filtering by state,
	// for example: include, exclude or return only deleted values
	State uint
)

// State* constants aid with Filter*
const (
	// StateExcluded do not include entries
	StateExcluded State = 0

	// StateInclusive include entries
	StateInclusive State = 1

	// StateExclusive only entries that have this state
	StateExclusive State = 2
)

// squirrel.SelectBuilder
func StateCondition(q squirrel.SelectBuilder, field string, fs State) squirrel.SelectBuilder {
	switch fs {
	case StateExclusive:
		// only null values
		return q.Where(squirrel.NotEq{field: nil})

	case StateInclusive:
		// mo filter
		return q

	default:
		// exclude all non-null values
		return q.Where(squirrel.Eq{field: nil})
	}
}
