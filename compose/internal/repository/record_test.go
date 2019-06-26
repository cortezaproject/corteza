package repository

import (
	"fmt"
	"strings"
	"testing"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/internal/test"
)

func TestRecordFinder(t *testing.T) {
	r := record{}
	m := &types.Module{
		ID:          123,
		NamespaceID: 456,
		Fields: types.ModuleFieldSet{
			&types.ModuleField{Name: "foo"},
			&types.ModuleField{Name: "bar"},
		},
	}

	ttc := []struct {
		f     types.RecordFilter
		match []string
		args  []interface{}
		err   error
	}{
		{
			match: []string{"SELECT * FROM compose_record AS r WHERE r.deleted_at IS NULL AND r.module_id = ?"},
		},
		{
			f: types.RecordFilter{Filter: "id = 5 AND foo = 7"},
			match: []string{
				"id = 5",
				"rv_foo.value = 7"},
			args: []interface{}{"foo"},
		},
		{
			f: types.RecordFilter{Sort: "id ASC, bar DESC"},
			match: []string{
				" id ASC",
				" rv_bar.value DESC",
			},
			args: []interface{}{"bar"},
		},
	}

	for _, tc := range ttc {
		sb, err := r.buildQuery(m, tc.f)

		if tc.err != nil {
			test.Assert(t, tc.err.Error() == fmt.Sprintf("%v", err), "buildQuery(%+v) did not return an expected error %q but %q", tc.f, tc.err, err)
		} else {
			test.Assert(t, err == nil, "buildQuery(%+v) returned an unexpected error: %v", tc.f, err)
		}

		sb = sb.Column("*")
		sql, args, err := sb.ToSql()

		for _, m := range tc.match {
			test.Assert(t, strings.Contains(sql, m),
				"assertion failed; query %q \n  "+
					"             did not contain  %q", sql, m)
		}

		tc.args = append(tc.args, m.ID, m.NamespaceID)
		test.Assert(t, fmt.Sprintf("%+v", args) == fmt.Sprintf("%+v", tc.args),
			"assertion failed; args %+v \n  "+
				"     do not match expected %+v", args, tc.args)
	}
}
