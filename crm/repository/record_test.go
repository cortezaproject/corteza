package repository

import (
	"strings"
	"testing"

	"github.com/crusttech/crust/crm/types"
	"github.com/crusttech/crust/internal/test"
)

func TestRecordFinder(t *testing.T) {
	r := record{}
	m := &types.Module{
		ID: 123,
		Fields: types.ModuleFieldSet{
			&types.ModuleField{Name: "foo"},
			&types.ModuleField{Name: "bar"},
		},
	}

	ttc := []struct {
		filter string
		sort   string
		match  []string
		args   []interface{}
	}{
		{
			match: []string{"SELECT * FROM crm_record WHERE module_id = ? AND deleted_at IS NULL"},
			args:  []interface{}{123}},
		{
			filter: "id = 5 AND foo = 7",
			match: []string{
				" AND id = 5",
				" AND (SELECT column_value FROM crm_record_column WHERE column_name = ? AND record_id = crm_record.id) = 7"},
			args: []interface{}{123}},
		{
			sort: "id ASC, foo DESC",
			match: []string{
				" id ASC, (SELECT column_value FROM crm_record_column WHERE column_name = 'foo' AND record_id = crm_record.id) DESC"},
			args: []interface{}{123}},
	}

	for _, tc := range ttc {
		sb, err := r.buildQuery(m, tc.filter, tc.sort)
		test.Assert(t, err == nil, "buildQuery(%q, %q) returned an error: %v", tc.filter, tc.sort, err)
		sb = sb.Column("*")
		sql, args, err := sb.ToSql()

		for _, m := range tc.match {
			test.Assert(t, strings.Contains(sql, m),
				"assertion failed; query %q \n  "+
					"             did not contain  %q", sql, m)
		}

		_ = args
		// test.Assert(t, reflect.DeepEqual(args, tc.args),
		// 	"assertion failed; args %v \n  "+
		// 		"     do not match expected %v", args, tc.args)
	}

}
