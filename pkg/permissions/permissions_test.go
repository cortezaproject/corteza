package permissions

import (
	"reflect"
	"testing"
)

func TestDynamicRoles(t *testing.T) {
	tests := []struct {
		name string
		u    uint64
		cc   []uint64
		exp  []uint64
	}{
		{
			"empty",
			42,
			nil,
			[]uint64{},
		},
		{
			"only one",
			42,
			[]uint64{42, 2},
			[]uint64{2},
		},
		{
			"none",
			42,
			[]uint64{1, 2},
			[]uint64{},
		},
		{
			"few",
			42,
			[]uint64{42, 2, 43, 3},
			[]uint64{2},
		},
		{
			"all",
			42,
			[]uint64{42, 1, 42, 2},
			[]uint64{1, 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRr := DynamicRoles(tt.u, tt.cc...); !reflect.DeepEqual(gotRr, tt.exp) {
				t.Errorf("DynamicRoles() = %v, want %v", gotRr, tt.exp)
			}
		})
	}
}
