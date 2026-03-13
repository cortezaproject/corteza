package j7s

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetBySlicePath(t *testing.T) {
	root := mapNode{
		"a": map[string]any{
			"b": map[string]any{
				"c": "deep",
			},
			"num": 42,
		},
		"top": "value",
	}

	tt := []struct {
		name      string
		parts     []string
		wantVal   any
		wantFound bool
		wantErr   bool
	}{
		{
			name:    "empty parts",
			parts:   []string{},
			wantErr: true,
		},
		{
			name:    "empty segment",
			parts:   []string{"a", ""},
			wantErr: true,
		},
		{
			name:      "top-level key",
			parts:     []string{"top"},
			wantVal:   "value",
			wantFound: true,
		},
		{
			name:      "nested two levels",
			parts:     []string{"a", "num"},
			wantVal:   42,
			wantFound: true,
		},
		{
			name:      "nested three levels",
			parts:     []string{"a", "b", "c"},
			wantVal:   "deep",
			wantFound: true,
		},
		{
			name:      "missing key at top",
			parts:     []string{"missing"},
			wantVal:   nil,
			wantFound: false,
		},
		{
			name:      "missing key mid-path",
			parts:     []string{"a", "missing"},
			wantVal:   nil,
			wantFound: false,
		},
		{
			name:      "non-map mid-path",
			parts:     []string{"top", "child"},
			wantVal:   nil,
			wantFound: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			val, found, err := GetBySlicePath(root, tc.parts)

			if tc.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			require.Equal(t, tc.wantFound, found)
			require.Equal(t, tc.wantVal, val)
		})
	}
}
