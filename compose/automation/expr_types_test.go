package automation

import (
	"context"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSetRecordValuesWithPath(t *testing.T) {
	t.Run("some basic aerobics", func(t *testing.T) {

		var (
			r   = require.New(t)
			rvs = &ComposeRecordValues{types.RecordValueSet{}}
		)

		r.NoError(expr.Assign(rvs, "field1", "a"))
		r.NoError(expr.Assign(rvs, "field1.1", "a"))
		r.True(rvs.value.Has("field1", 0))
		r.True(rvs.value.Has("field1", 1))
	})

	t.Run("cast string map", func(t *testing.T) {

		var (
			r        = require.New(t)
			rvs, err = CastToComposeRecordValues(map[string]string{"field2": "b"})
		)

		r.NoError(err)
		r.True(rvs.Has("field2", 0))
	})

	t.Run("cast string slice map", func(t *testing.T) {

		var (
			r        = require.New(t)
			rvs, err = CastToComposeRecordValues(map[string][]string{"field2": []string{"a", "b"}})
		)

		r.NoError(err)
		r.True(rvs.Has("field2", 0))
		r.True(rvs.Has("field2", 1))
	})
}

func TestRecordFieldValuesAccess(t *testing.T) {
	var (
		err error
		v   expr.TypedValue

		mod = &types.Module{Fields: types.ModuleFieldSet{
			&types.ModuleField{Name: "s1", Multi: false},
			&types.ModuleField{Name: "m1", Multi: true},
			&types.ModuleField{Name: "m2", Multi: true},
			&types.ModuleField{Name: "s2", Multi: false},
		}}

		raw = &types.Record{Values: types.RecordValueSet{
			&types.RecordValue{Name: "s1", Value: "sVal1"},
			&types.RecordValue{Name: "m1", Value: "mVal1.0"},
			&types.RecordValue{Name: "m1", Value: "mVal1.1", Place: 1},
			&types.RecordValue{Name: "m1", Value: "mVal1.2", Place: 2},
			&types.RecordValue{Name: "m2", Value: "mVal2.0"},
		}}

		tval  = &ComposeRecord{value: raw}
		scope = expr.RVars{"rec": tval}.Vars()
	)

	// @todo see not above re. back-ref to record
	raw.SetModule(mod)

	t.Run("via typed value", func(t *testing.T) {
		var (
			req = require.New(t)
		)

		v, err = expr.Select(scope, "rec.values.s1")
		req.NoError(err)
		req.NotEmpty(v)
		req.Equal("sVal1", v.Get())

		v, err = expr.Select(scope, "rec.values.m1.0")
		req.NoError(err)
		req.NotEmpty(v)
		req.Equal("mVal1.0", v.Get())

		v, err = expr.Select(scope, "rec.values.m1.1")
		req.NoError(err)
		req.NotEmpty(v)
		req.Equal("mVal1.1", v.Get())

		// @todo when RecordValueSet supports back-ref to record,
		//       we can employ better field access:
		//        - no error on missing values when field exists
		//        - proper handling of multi-value field values
		//        - proper value-types that corelate to field types
		//v, err = expr.Select(scope, "rec.values.m2.0")
		//req.NoError(err)
		//req.NotEmpty(v)
		//req.Equal("mVal2.0", v.Get())
	})

	t.Run("via gval selector", func(t *testing.T) {
		var (
			req    = require.New(t)
			parser = expr.NewParser()
		)

		eval, err := parser.Parse(`rec.values.s1 == "sVal1"`)
		req.NoError(err)
		req.True(eval.Test(context.Background(), scope))

		eval, err = parser.Parse(`rec.values.s1 != "foo"`)
		req.NoError(err)
		req.True(eval.Test(context.Background(), scope))

		eval, err = parser.Parse(`rec.values.m1[0] == "mVal1.0"`)
		req.NoError(err)
		req.True(eval.Test(context.Background(), scope))

		eval, err = parser.Parse(`rec.values.m1[1] == "mVal1.1"`)
		req.NoError(err)
		req.True(eval.Test(context.Background(), scope))

		eval, err = parser.Parse(`rec.values.m2[0] == "mVal2.0"`)
		req.NoError(err)
		req.True(eval.Test(context.Background(), scope))

	})

}
