package repository

import (
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/rh"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cortezaproject/corteza-server/compose/types"
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
		f       types.RecordFilter
		match   []string
		noMatch []string
		args    []interface{}
		err     error
	}{
		{
			match: []string{
				"SELECT r.id, r.module_id, r.rel_namespace, r.owned_by, r.created_at, " +
					"r.created_by, r.updated_at, r.updated_by, r.deleted_at, r.deleted_by " +
					"FROM compose_record AS r " +
					"WHERE r.module_id = ? AND r.rel_namespace = ? AND r.deleted_at IS NULL",
			},
		},
		{
			f: types.RecordFilter{Query: "id = 5 AND foo = 7"},
			match: []string{
				"r.id = 5",
				"rv_foo.value = 7"},
			args: []interface{}{"foo"},
		},
		{
			f: types.RecordFilter{Sort: "id ASC, bar DESC"},
			match: []string{
				" r.id ASC",
				" rv_bar.value DESC",
			},
			args: []interface{}{"bar"},
		},
		{
			f:     types.RecordFilter{Deleted: rh.FilterStateExcluded},
			match: []string{" r.deleted_at IS "},
		},
		{
			f:       types.RecordFilter{Deleted: rh.FilterStateInclusive},
			noMatch: []string{" r.deleted_at IS "},
		},
		{
			f:     types.RecordFilter{Deleted: rh.FilterStateExclusive},
			match: []string{" r.deleted_at IS NOT NULL"},
		},
	}

	for _, tc := range ttc {
		sb, err := r.buildQuery(m, tc.f)

		if tc.err != nil {
			require.True(t, tc.err.Error() == fmt.Sprintf("%v", err), "buildQuery(%+v) did not return an expected error %q but %q", tc.f, tc.err, err)
		} else {
			require.True(t, err == nil, "buildQuery(%+v) returned an unexpected error: %v", tc.f, err)
		}

		sql, args, err := sb.ToSql()

		for _, m := range tc.match {
			require.True(t, strings.Contains(sql, m),
				"assertion failed; query %q \n  "+
					"             did not contain  %q", sql, m)
		}

		for _, m := range tc.noMatch {
			require.False(t, strings.Contains(sql, m),
				"assertion failed; query %q \n  "+
					"             must not contain  %q", sql, m)
		}

		tc.args = append(tc.args, m.ID, m.NamespaceID)
		require.True(t, fmt.Sprintf("%+v", args) == fmt.Sprintf("%+v", tc.args),
			"assertion failed; args %+v \n  "+
				"     do not match expected %+v", args, tc.args)
	}
}
