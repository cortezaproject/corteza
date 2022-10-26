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
				t.Errorf("got:\n%+v\n\nwant\n%+v", got, tt.want)
			}
		})
	}
}

func TestRecordValueSet_Merge(t *testing.T) {
	tests := []struct {
		name string
		set  RecordValueSet
		new  RecordValueSet
		mfs  ModuleFieldSet
		fn   func(f *ModuleField) bool
		want RecordValueSet
	}{
		{
			name: "update with partial access",
			set:  RecordValueSet{{Name: "n", Value: "1"}, {Name: "accessible", Value: "skip me"}, {Name: "inaccessible", Value: "don't skip me"}},
			new:  RecordValueSet{{Name: "n", Value: "2"}},
			mfs: ModuleFieldSet{
				&ModuleField{Name: "n"},
				&ModuleField{Name: "accessible"},
				&ModuleField{Name: "inaccessible"},
			},
			fn: func(f *ModuleField) bool {
				// is field accessible?
				return f.Name == "accessible"
			},
			want: RecordValueSet{
				{Name: "n", Value: "2", OldValue: "1", Updated: true},
				{Name: "accessible", Value: "skip me", OldValue: "skip me", Updated: true, DeletedAt: nil},
				{Name: "inaccessible", Value: "don't skip me", OldValue: "don't skip me"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.set.Merge(tt.mfs, tt.new, tt.fn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got:\n%+v\n\nwant\n%+v", got, tt.want)
			}
		})
	}
}

func TestRecordValueSet_merge(t *testing.T) {
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
			want: RecordValueSet{{Name: "n", Value: "v", Updated: true}},
		},
		{
			name: "update nil",
			set:  nil,
			new:  RecordValueSet{{Name: "n", Value: "v"}},
			want: RecordValueSet{{Name: "n", Value: "v", OldValue: "", Updated: true}},
		},
		{
			name: "update with nil",
			set:  RecordValueSet{{Name: "n", Value: "v"}},
			new:  nil,
			want: RecordValueSet{{Name: "n", Value: "v", OldValue: "v", DeletedAt: nil, Updated: true}},
		},
		{
			name: "update with new value",
			set:  RecordValueSet{{Name: "n", Value: "1"}},
			new:  RecordValueSet{{Name: "n", Value: "2"}},
			want: RecordValueSet{{Name: "n", Value: "2", OldValue: "1", Updated: true}},
		},
		{
			name: "update with less values",
			set:  RecordValueSet{{Name: "n", Value: "1"}, {Name: "deleted", Value: "d"}},
			new:  RecordValueSet{{Name: "n", Value: "2"}},
			want: RecordValueSet{{Name: "n", Value: "2", OldValue: "1", Updated: true}, {Name: "deleted", Value: "d", OldValue: "d", Updated: true, DeletedAt: nil}},
		},
		{
			name: "update multi value",
			set:  RecordValueSet{{Name: "c", Value: "1st", Place: 1}, {Name: "c", Value: "2nd", Place: 2}, {Name: "c", Value: "3rd", Place: 3}, {Name: "c", Value: "4th", Place: 4}},
			new:  RecordValueSet{{Name: "c", Value: "1st", Place: 1}, {Name: "c", Value: "2nd", Place: 2}, {Name: "c", Value: "4th", Place: 3}},
			want: RecordValueSet{
				{Name: "c", Value: "1st", Place: 1, OldValue: "1st"},
				{Name: "c", Value: "2nd", Place: 2, OldValue: "2nd"},
				{Name: "c", Value: "4th", Place: 3, OldValue: "3rd", Updated: true},
				{Name: "c", Value: "4th", Place: 4, OldValue: "4th", Updated: true, DeletedAt: nil},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.set.merge(tt.new); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got:\n%+v\n\nwant\n%+v", got, tt.want)
			}
		})
	}
}
