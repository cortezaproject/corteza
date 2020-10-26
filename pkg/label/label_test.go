package label

import "testing"

func TestChanged(t *testing.T) {
	tests := []struct {
		name string
		old  map[string]string
		new  map[string]string
		want bool
	}{
		// TODO: Add test cases.
		{
			"2x nil",
			nil,
			nil,
			false,
		},
		{
			"2x empty",
			map[string]string{},
			map[string]string{},
			false,
		},
		{
			"nil & empty",
			nil,
			map[string]string{},
			false,
		},
		{
			"same",
			map[string]string{"a": "a"},
			map[string]string{"a": "a"},
			false,
		},
		{
			"diff1",
			map[string]string{"a": "a"},
			map[string]string{"a": "b"},
			true,
		},
		{
			"diff2",
			map[string]string{"a": "b"},
			map[string]string{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Changed(tt.old, tt.new); got != tt.want {
				t.Errorf("Changed() = %v, want %v", got, tt.want)
			}
		})
	}
}
