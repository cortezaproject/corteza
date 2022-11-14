package auth

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestApplyRoleSecurity(t *testing.T) {
	tests := []struct {
		name       string
		permitted  []uint64
		prohibited []uint64
		forced     []uint64
		roles      []uint64
		wantOut    []uint64
	}{
		{
			"empty",
			[]uint64{},
			[]uint64{},
			[]uint64{},
			[]uint64{},
			[]uint64{},
		},
		{
			"nil",
			nil,
			nil,
			nil,
			nil,
			[]uint64{},
		},
		{
			"one",
			[]uint64{1},
			[]uint64{2},
			[]uint64{3},
			[]uint64{1, 2},
			[]uint64{1, 3},
		},
		{
			"forced only",
			[]uint64{},
			[]uint64{},
			[]uint64{3, 2, 1},
			[]uint64{2},
			[]uint64{1, 2, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.wantOut, ApplyRoleSecurity(tt.permitted, tt.prohibited, tt.forced, tt.roles...))
		})
	}
}
