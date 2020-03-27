package rh

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLimit(t *testing.T) {
	var (
		r = require.New(t)
	)

	r.Equal(Limit(42).Limit, uint(42))
	r.Equal(Limit(0, 42).Offset, uint(42))
}

func Test_parsePagination(t *testing.T) {
	var (
		tests = []struct {
			name    string
			args    interface{}
			pf      PageFilter
			wantErr bool
		}{
			{
				"empty",
				nil,
				PageFilter{},
				false,
			},
			{
				"valid l/o",
				map[string]string{"limit": "42", "offset": "314"},
				PageFilter{Limit: 42, Offset: 314},
				false,
			},
			{
				"mixed",
				map[string]string{"page": "42", "limit": "314"},
				PageFilter{Limit: 314, Offset: 0},
				false,
			},
			{
				"invalid limit",
				map[string]string{"limit": "abc"},
				PageFilter{},
				true,
			},
			{
				"invalid page",
				map[string]string{"page": "abc"},
				PageFilter{},
				true,
			},
		}
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				pf = PageFilter{}
			)

			if err := parsePagination(&pf, tt.args); (err != nil) != tt.wantErr {
				t.Errorf("parsePagination() error = %v, wantErr %v", err, tt.wantErr)
			} else if !reflect.DeepEqual(pf, tt.pf) {
				t.Errorf("\n  actual: %v\nexpected: %v\n", pf, tt.pf)
			}
		})
	}
}
