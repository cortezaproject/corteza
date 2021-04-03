package automation

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSetRecordValuesWithPath(t *testing.T) {
	t.Run("some basic aerobics", func(t *testing.T) {

		var (
			r   = require.New(t)
			rvs = &ComposeRecordValues{&types.Record{}}
		)

		r.NoError(expr.Assign(rvs, "field1", expr.Must(expr.NewString("a"))))
		r.NoError(expr.Assign(rvs, "field1.1", expr.Must(expr.NewString("a"))))
		r.True(rvs.value.Values.Has("field1", 0))
		r.True(rvs.value.Values.Has("field1", 1))
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
			&types.ModuleField{Name: "s1", Multi: false, Kind: "String"},
			&types.ModuleField{Name: "m1", Multi: true, Kind: "String"},
			&types.ModuleField{Name: "m2", Multi: true, Kind: "String"},
			&types.ModuleField{Name: "s2", Multi: false, Kind: "String"},
			&types.ModuleField{Name: "b0", Multi: false, Kind: "Bool"},
			&types.ModuleField{Name: "b1", Multi: false, Kind: "Bool"},
			&types.ModuleField{Name: "n1", Multi: false, Kind: "Number"},
			&types.ModuleField{Name: "n2", Multi: false, Kind: "Number"},
			&types.ModuleField{Name: "n3", Multi: false, Kind: "Number"},
			&types.ModuleField{Name: "ref1", Multi: false, Kind: "Record"},
			&types.ModuleField{Name: "ref2", Multi: false, Kind: "Record"},
		}}

		raw = &types.Record{Values: types.RecordValueSet{
			&types.RecordValue{Name: "s1", Value: "sVal1"},
			&types.RecordValue{Name: "m1", Value: "mVal1.0", Place: 0},
			&types.RecordValue{Name: "m1", Value: "mVal1.1", Place: 1},
			&types.RecordValue{Name: "m1", Value: "mVal1.2", Place: 2},
			&types.RecordValue{Name: "m2", Value: "mVal2.0", Place: 0},
			&types.RecordValue{Name: "b1", Value: "1", Place: 0},
			&types.RecordValue{Name: "n2", Value: "0", Place: 0},
			&types.RecordValue{Name: "n3", Value: "2", Place: 0},
			&types.RecordValue{Name: "ref2", Value: "", Ref: 2, Place: 0},
		}}

		scope, _ = expr.NewVars(map[string]interface{}{
			"rec": &ComposeRecord{value: raw},

			// nis is only for testing, gval does not recognise nil (or null) as a keyword!
			"nil": nil,

			// typed value but empty (this happens when you assign next value in expression but do not set a value to it
			"initRec": &ComposeRecord{},

			// same as &ComposeRecord{value: &types.Record{}},
			"validRecZero":    &ComposeRecord{value: &types.Record{ID: 0, Values: types.RecordValueSet{}}},
			"validRecValidID": &ComposeRecord{&types.Record{ID: 42, Values: types.RecordValueSet{}}},

			// "record" (not really) set to nil
			"fooRec": nil,
		})
	)

	// @todo see not above re. back-ref to record
	raw.SetModule(mod)

	t.Run("via typed value", func(t *testing.T) {
		tcc := []struct {
			expects interface{}
			path    string
		}{
			{"sVal1", "rec.values.s1"},
			{"mVal1.0", "rec.values.m1.0"},
			{"mVal1.1", "rec.values.m1.1"},
			{"mVal2.0", "rec.values.m2.0"},
			// expecting valid value (false)  even when boolean fields are not set
			{false, "rec.values.b0"},
			{true, "rec.values.b1"},
			{uint64(2), "rec.values.ref2"},
		}

		for _, tc := range tcc {
			t.Run(tc.path, func(t *testing.T) {
				var (
					req = require.New(t)
				)

				v, err = expr.Select(scope, tc.path)
				req.NoError(err)
				req.Equal(tc.expects, v.Get())
			})
		}
	})

	t.Run("via gval selector", func(t *testing.T) {
		tcc := []struct {
			test interface{}
			expr string
		}{
			{false, `nil`},
			{true, `rec`},
			{true, `rec.values`},

			// not set, so false
			{fmt.Errorf("failed to select 'notSetRec' on *expr.Vars: no such key 'notSetRec'"), `notSetRec`},

			// set but nil
			{false, `fooRec`},

			// set, initialized
			{true, `initRec`},

			// set, initialized, but recordID is empty
			{false, `initRec.recordID`},

			// set, initialized, but recordID (ID is alias) is empty
			{false, `initRec.ID`},

			// set, initialized, but values are empty
			{false, `initRec.values`},

			{true, `validRecZero`},
			{false, `validRecZero.recordID`},
			{false, `validRecZero.values`},

			{true, `validRecValidID`},
			{true, `validRecValidID.recordID`},
			{false, `validRecValidID.values`},

			// but we can not access the values below...
			{fmt.Errorf("unknown parameter initRec.values.foo"), `initRec.values.foo`},

			// interaction with set values
			{true, `rec.values.s1 == "sVal1"`},
			{false, `rec.values.s1 == "sVal2"`},
			{true, `rec.values.s1`},
			{true, `rec.values.s1 != "foo"`},

			// interaction with unset (= nil) values
			{true, `rec.values.s2 != "foo"`},
			{false, `rec.values.s2 == "foo"`},
			{true, `!rec.values.s2`},
			{false, `rec.values.s2`},

			// multival
			{true, `rec.values.m1[0] == "mVal1.0"`},
			{true, `rec.values.m1[1] == "mVal1.1"`},
			{true, `rec.values.m2[0] == "mVal2.0"`},

			// booleans
			{true, `!rec.values.b0`},
			{false, `rec.values.b0`},
			{true, `rec.values.b1`},
			{false, `!rec.values.b1`},

			// numbers
			{false, `rec.values.n1`},
			{false, `rec.values.n2`},
			{true, `rec.values.n3`},

			{true, `rec.values.n1 == 0`},
			{true, `rec.values.n2 == 0`},
			{false, `rec.values.n3 == 0`},

			{false, `rec.values.n1 == 2`},
			{false, `rec.values.n2 == 2`},
			{true, `rec.values.n3 == 2`},

			//{true, `rec.values.n1 < 3`}, // invalid op <nil> < 3
			{true, `rec.values.n2 < 3`},
			{true, `rec.values.n3 < 3`},

			//{false, `rec.values.n1 > 1`}, // invalid op <nil> > 3
			{false, `rec.values.n2 > 2`},
			{false, `rec.values.n3 > 2`},

			{true, `rec.values.ref1 != 2`},
			{true, `rec.values.ref2 == 2`},
			{true, `rec.values.ref2 == "2"`},
		}

		for _, tc := range tcc {
			t.Run(tc.expr, func(t *testing.T) {
				var (
					req       = require.New(t)
					parser    = expr.NewParser()
					eval, err = parser.Parse(tc.expr)
				)

				req.NoError(err)

				test, err := eval.Test(context.Background(), scope)
				switch tct := tc.test.(type) {
				case error:
					req.EqualError(err, tct.Error())

				case bool:
					req.NoError(err)
					req.Equal(tc.test, test)

				default:
					panic("unexpected test case type")
				}

			})
		}
	})
}

func TestAssignToComposeRecordValues(t *testing.T) {
	t.Run("assign simple", func(t *testing.T) {
		var (
			req    = require.New(t)
			target = &types.Record{Values: types.RecordValueSet{}}
		)

		req.NoError(assignToComposeRecordValues(target, []string{"a"}, "b"))
		req.Len(target.Values, 1)
		req.True(target.Values.Has("a", 0))
		req.NoError(assignToComposeRecordValues(target, []string{"a", "1"}, "b"))
		req.Len(target.Values, 2)
		req.True(target.Values.Has("a", 0))
		req.True(target.Values.Has("a", 1))
	})

	t.Run("assign rvs", func(t *testing.T) {
		var (
			req    = require.New(t)
			target = &types.Record{Values: types.RecordValueSet{}}
		)

		req.NoError(assignToComposeRecordValues(target, nil, types.RecordValueSet{{}}))
		req.Len(target.Values, 1)
	})

	t.Run("assign record", func(t *testing.T) {
		var (
			req    = require.New(t)
			target = &types.Record{Values: types.RecordValueSet{}}
		)

		req.NoError(assignToComposeRecordValues(target, nil, &types.Record{Values: types.RecordValueSet{{}}}))
		req.Len(target.Values, 1)
	})

	t.Run("overwrite rvs", func(t *testing.T) {
		var (
			req    = require.New(t)
			target = &types.Record{Values: types.RecordValueSet{{Name: "a"}}}
		)

		req.NoError(assignToComposeRecordValues(target, nil, types.RecordValueSet{{Name: "b"}}))
		req.Len(target.Values, 1)
		req.False(target.Values.Has("a", 0))
		req.True(target.Values.Has("b", 0))
	})

	t.Run("assign multiple values", func(t *testing.T) {
		var (
			req    = require.New(t)
			target = &types.Record{Values: types.RecordValueSet{}}
		)

		req.Error(assignToComposeRecordValues(target, []string{"a", "2"}, expr.Must(expr.NewAny([]interface{}{"1", "2"}))))
		req.Len(target.Values, 0)

		req.NoError(assignToComposeRecordValues(target, []string{"a"}, expr.Must(expr.NewAny([]interface{}{"1", "2"}))))
		req.Len(target.Values, 2)

		req.NoError(assignToComposeRecordValues(target, []string{"a"}, expr.Must(expr.NewAny([]string{"1", "2"}))))
		req.Len(target.Values, 2)
	})
}
