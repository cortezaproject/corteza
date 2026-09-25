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
			SortExprSet{
				&SortExpr{
					columns: []string{"name"},
					Column:  "name",
				},
			},
			false,
		},
		{
			"one simple column, descending",
			"name desc",
			SortExprSet{
				&SortExpr{
					columns:    []string{"name"},
					Column:     "name",
					Descending: true},
			},
			false,
		},
		{
			"combo",
			"name desc, email asc, age desc",
			SortExprSet{
				&SortExpr{
					columns:    []string{"name"},
					Column:     "name",
					Descending: true,
				},
				&SortExpr{
					columns:    []string{"email"},
					Column:     "email",
					Descending: false,
				},
				&SortExpr{
					columns:    []string{"age"},
					Column:     "age",
					Descending: true,
				},
			},
			false,
		},
		{
			"empty",
			"",
			SortExprSet{},
			false,
		},
		{
			"COALESCE with multiple column",
			"COALESCE(updated_at,created_at)",
			SortExprSet{
				&SortExpr{
					modifier:   "COALESCE",
					columns:    []string{"updated_at", "created_at"},
					Descending: false,
				},
			},
			false,
		},
		{
			"COALESCE with multiple column, descending",
			"COALESCE(updated_at,created_at) desc",
			SortExprSet{
				&SortExpr{
					modifier:   "COALESCE",
					columns:    []string{"updated_at", "created_at"},
					Descending: true,
				},
			},
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
			&tmp{Sorting: Sorting{Sort: SortExprSet{&SortExpr{Column: "name", columns: []string{"name"}, Descending: true}}}, Other: true},
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

// TestParseSortRefusesAnUnknownModifier pins the case the old guard could not
// reach: it tested s.modifier, which is only ever assigned in the COALESCE
// branch, so for a fresh expression the check was dead and "nonsense(createdAt)"
// parsed to a bare "createdAt" — a typo silently changing the ordering instead
// of failing.
func TestParseSortRefusesAnUnknownModifier(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		wantErr bool
	}{
		{name: "no modifier", in: "createdAt DESC"},
		{name: "the one known modifier", in: "coalesce(deletedAt, updatedAt, createdAt) DESC"},
		{name: "modifier casing is not significant", in: "COALESCE(deletedAt, createdAt)"},
		{name: "isnull ahead of the last-change column", in: "isnull(deletedAt), coalesce(deletedAt, updatedAt, createdAt) DESC"},
		{name: "a typo is refused, not dropped", in: "coalese(createdAt)", wantErr: true},
		{name: "an invented modifier is refused", in: "nonsense(createdAt)", wantErr: true},
		{name: "a bad modifier on a later term still fails", in: "createdAt, nonsense(updatedAt)", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set, err := parseSort(tt.in)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("parseSort(%q) = %v, want an error", tt.in, set.String())
				}
				return
			}

			if err != nil {
				t.Fatalf("parseSort(%q) returned %v", tt.in, err)
			}
		})
	}
}
