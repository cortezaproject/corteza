package payload

import (
	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParseFilterState(t *testing.T) {
	var (
		req = require.New(t)
	)
	tests := []struct {
		name     string
		input    string
		expected filter.State
	}{
		{
			"invalid string should be Excluded",
			"zero",
			filter.StateExcluded,
		},
		{
			"empty string should be Excluded",
			"0",
			filter.StateExcluded,
		},
		{
			"Excluded",
			"0",
			filter.StateExcluded,
		},
		{
			"Excluded",
			"1",
			filter.StateInclusive,
		},
		{
			"Excluded",
			"2",
			filter.StateExclusive,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req.Equal(tt.expected, ParseFilterState(tt.input))
		})
	}
}
