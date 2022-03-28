package filter

import (
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
