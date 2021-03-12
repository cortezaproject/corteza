package filter

import (
	"github.com/Masterminds/squirrel"
	"strconv"
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

func (s State) String() string {
	return strconv.Itoa(int(s))
}

// squirrel.SelectBuilder
func StateCondition(q squirrel.SelectBuilder, field string, fs State) squirrel.SelectBuilder {
	switch fs {
	case StateExclusive:
		// only null values
		return q.Where(squirrel.NotEq{field: nil})

	case StateInclusive:
		// no filter
		return q

	default:
		// exclude all non-null values
		return q.Where(squirrel.Eq{field: nil})
	}
}

// squirrel.SelectBuilder
func StateConditionNegBool(q squirrel.SelectBuilder, field string, fs State) squirrel.SelectBuilder {
	switch fs {
	case StateExcluded:
		// only true
		return q.Where(squirrel.Eq{field: true})

	case StateExclusive:
		// only false
		return q.Where(squirrel.Eq{field: false})

	default:
		return q
	}
}
