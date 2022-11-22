package rdbms

import (
	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/doug-martin/goqu/v9"
	"reflect"
	"testing"
)

func Test_generateSorting(t *testing.T) {
	tests := []struct {
		name      string
		sort      filter.SortExpr
		columns   []string
		modifier  string
		sortables map[string]string
		wantOut   goqu.Expression
		wantErr   bool
	}{
		{
			name: "single column sorting",
			sort: filter.SortExpr{Column: "createdat"},
			sortables: map[string]string{
				"createdat": "created_at",
			},
			wantOut: goqu.I("created_at"),
			wantErr: false,
		},
		{
			name: "path sorting",
			sort: filter.SortExpr{Column: "user.createdat"},
			sortables: map[string]string{
				"user.createdat": "user.created_at",
			},
			wantOut: goqu.I("user.created_at"),
			wantErr: false,
		},
		{
			name:     "coalesce single column sorting",
			sort:     filter.SortExpr{},
			columns:  []string{"createdat"},
			modifier: "coalesce",
			sortables: map[string]string{
				"createdat": "created_at",
				"updatedat": "updated_at",
			},
			wantOut: goqu.COALESCE(goqu.I("created_at")),
			wantErr: false,
		},
		{
			name:     "coalesce multiple column sorting",
			sort:     filter.SortExpr{},
			columns:  []string{"createdat", "updatedat"},
			modifier: "coalesce",
			sortables: map[string]string{
				"createdat": "created_at",
				"updatedat": "updated_at",
			},
			wantOut: goqu.COALESCE(goqu.I("created_at"), goqu.I("updated_at")),
			wantErr: false,
		},
		{
			name:     "invalid coalesce multiple column sorting",
			sort:     filter.SortExpr{},
			columns:  []string{"createdat", "deletedat"},
			modifier: "coalesce",
			sortables: map[string]string{
				"createdat": "created_at",
				"updatedat": "updated_at",
			},
			wantOut: nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.sort.SetColumns(tt.columns...)
			_ = tt.sort.SetModifier(tt.modifier)
			gotOut, err := generateSorting(tt.sortables, &tt.sort)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateSorting() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotOut, tt.wantOut) {
				t.Errorf("generateSorting() gotOut = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}
