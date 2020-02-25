package types

import (
	"reflect"
	"testing"
	"time"
)

func TestRecordValueSet_Set(t *testing.T) {
	tests := []struct {
		name string
		set  RecordValueSet
		new  RecordValue
		want RecordValueSet
	}{
		{
			name: "simple add on empty",
			set:  RecordValueSet{},
			new:  RecordValue{Name: "n", Value: "v"},
			want: RecordValueSet{{Name: "n", Value: "v"}},
		},
		{
			name: "update existing",
			set:  RecordValueSet{{Name: "a", Value: "b"}, {Name: "n", Value: "v"}, {Name: "x", Value: "y"}},
			new:  RecordValue{Name: "n", Value: "v2"},
			want: RecordValueSet{{Name: "a", Value: "b"}, {Name: "n", Value: "v2"}, {Name: "x", Value: "y"}},
		},
		{
			name: "multi-value",
			set:  RecordValueSet{{Name: "n", Value: "v"}},
			new:  RecordValue{Name: "n", Value: "v", Place: 1},
			want: RecordValueSet{{Name: "n", Value: "v", Place: 0}, {Name: "n", Value: "v", Place: 1}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.set.Set(&tt.new); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Set() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRecordValueSet_Merge(t *testing.T) {
	tests := []struct {
		name string
		set  RecordValueSet
		new  RecordValueSet
		want RecordValueSet
	}{
		{
			name: "simple update of an empty set",
			set:  RecordValueSet{},
			new:  RecordValueSet{{Name: "n", Value: "v"}},
			want: RecordValueSet{{Name: "n", Value: "v", updated: true}},
		},
		{
			name: "update with nil",
			set:  RecordValueSet{{Name: "n", Value: "v"}},
			new:  nil,
			want: RecordValueSet{{Name: "n", Value: "v", oldValue: "v", DeletedAt: &time.Time{}, updated: true}},
		},
		{
			name: "update with new value",
			set:  RecordValueSet{{Name: "n", Value: "1"}},
			new:  RecordValueSet{{Name: "n", Value: "2"}},
			want: RecordValueSet{{Name: "n", Value: "2", oldValue: "1", updated: true}},
		},
		{
			name: "update with less values",
			set:  RecordValueSet{{Name: "n", Value: "1"}, {Name: "deleted", Value: "d"}},
			new:  RecordValueSet{{Name: "n", Value: "2"}},
			want: RecordValueSet{{Name: "n", Value: "2", oldValue: "1", updated: true}, {Name: "deleted", Value: "d", oldValue: "d", updated: true, DeletedAt: &time.Time{}}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.set.Merge(tt.new); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Update() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
