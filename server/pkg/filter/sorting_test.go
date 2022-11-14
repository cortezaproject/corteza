package filter

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func Test_parseSort(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		wantSet SortExprSet
		wantErr bool
	}{
		{
			"one simple column",
			"name",
			SortExprSet{&SortExpr{Column: "name"}},
			false,
		},
		{
			"one simple column, descending",
			"name desc",
			SortExprSet{&SortExpr{Column: "name", Descending: true}},
			false,
		},
		{
			"combo",
			"name desc, email asc, age desc",
			SortExprSet{
				&SortExpr{Column: "name", Descending: true},
				&SortExpr{Column: "email", Descending: false},
				&SortExpr{Column: "age", Descending: true},
			},
			false,
		},
		{
			"empty",
			"",
			SortExprSet{},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSet, err := parseSort(tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseSort() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotSet, tt.wantSet) {
				t.Errorf("parseSort() gotSet = %v, want %v", gotSet, tt.wantSet)
			}
		})
	}
}

func TestSortUmarshaling(t *testing.T) {
	type tmp struct {
		Sorting
		Other bool
	}

	tests := []struct {
		name string
		in   string
		out  *tmp
	}{
		{
			"one simple column",
			`{"sort": "name DESC", "other": true}`,
			&tmp{Sorting: Sorting{Sort: SortExprSet{&SortExpr{Column: "name", Descending: true}}}, Other: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var unm = &tmp{}

			req := require.New(t)
			req.NoError(json.Unmarshal([]byte(tt.in), unm))
			req.Equal(tt.out, unm)
		})
	}
}
