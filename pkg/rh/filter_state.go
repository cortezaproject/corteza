package rh

import (
	"github.com/Masterminds/squirrel"
)

type (
	// FilterState for filtering by state,
	// for example: include, exclude or return only deleted values
	FilterState uint
)

// FilterState* constants aid with Filter*
const (
	// FilterStateExcluded do not include entries
	FilterStateExcluded FilterState = 0

	// FilterStateInclusive include entries
	FilterStateInclusive FilterState = 1

	// FilterStateExclusive only entries that have this state
	FilterStateExclusive FilterState = 2
)

// squirrel.SelectBuilder
func FilterNullByState(q squirrel.SelectBuilder, field string, fs FilterState) squirrel.SelectBuilder {
	switch fs {
	case FilterStateExclusive:
		// only null values
		return q.Where(squirrel.NotEq{field: nil})

	case FilterStateInclusive:
		// mo filter
		return q

	default:
		// exclude all non-null values
		return q.Where(squirrel.Eq{field: nil})
	}
}
