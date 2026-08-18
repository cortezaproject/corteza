package types

import (
	"testing"
)

func TestModuleFieldOptions_Int64Def(t *testing.T) {
	tests := []struct {
		name string
		opt  ModuleFieldOptions
		key  string
		def  int64
		want int64
	}{
		{"unexisting", ModuleFieldOptions{}, "k", 42, 42},
		{"nil", ModuleFieldOptions{"k": nil}, "k", 42, 42},
		{"bool", ModuleFieldOptions{"k": true}, "k", 42, 42},
		{"int", ModuleFieldOptions{"k": 1}, "k", 42, 1},
		{"float", ModuleFieldOptions{"k": 1.00000000001}, "k", 42, 1},
		{"stringed-int", ModuleFieldOptions{"k": "1"}, "k", 42, 1},
		{"stringed-float-1", ModuleFieldOptions{"k": "1.0"}, "k", 42, 1},
		{"stringed-float-2", ModuleFieldOptions{"k": "1.01"}, "k", 42, 1},
		{"stringed-float-3", ModuleFieldOptions{"k": "1.00000000001"}, "k", 42, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.opt.Int64Def(tt.key, tt.def); got != tt.want {
				t.Errorf("Int64Def() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModuleFieldOptions_Precision(t *testing.T) {
	tt := []struct {
		name string
		opt  ModuleFieldOptions
		want uint
	}{
		{"unset falls back to the webapp default", ModuleFieldOptions{}, 3},
		{"an explicit zero stays zero", ModuleFieldOptions{"precision": 0}, 0},
		{"an explicit value wins", ModuleFieldOptions{"precision": 2}, 2},
		{"above the maximum clamps", ModuleFieldOptions{"precision": 9}, 6},
	}

	for _, c := range tt {
		if got := c.opt.Precision(); got != c.want {
			t.Errorf("%s: got %d, want %d", c.name, got, c.want)
		}
	}
}
