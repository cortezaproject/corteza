package automation

import (
	"context"
	"fmt"
	"testing"
	"time"

	aTypes "github.com/cortezaproject/corteza/server/automation/types"
	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/stretchr/testify/require"
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
			&types.ModuleField{Name: "d1", Multi: false, Kind: "DateTime"},
			&types.ModuleField{Name: "d2", Multi: false, Kind: "DateTime"},
			&types.ModuleField{Name: "ref1", Multi: false, Kind: "Record"},
			&types.ModuleField{Name: "ref2", Multi: false, Kind: "Record"},
		}}

		rawValues = types.RecordValueSet{
			&types.RecordValue{Name: "s1", Value: "sVal1"},
			&types.RecordValue{Name: "m1", Value: "mVal1.0", Place: 0},
			&types.RecordValue{Name: "m1", Value: "mVal1.1", Place: 1},
			&types.RecordValue{Name: "m1", Value: "mVal1.2", Place: 2},
			&types.RecordValue{Name: "m2", Value: "mVal2.0", Place: 0},
			&types.RecordValue{Name: "b1", Value: "1", Place: 0},
			&types.RecordValue{Name: "n2", Value: "0", Place: 0},
			&types.RecordValue{Name: "n3", Value: "2", Place: 0},
			&types.RecordValue{Name: "ref2", Value: "", Ref: 2, Place: 0},
			&types.RecordValue{Name: "d1", Value: "", Place: 0},
			&types.RecordValue{Name: "d2", Value: "1999-09-09 09:09:09.000000009 +0000 UTC", Place: 0},
		}
		raw = &types.Record{Values: rawValues}

		scope, _ = expr.NewVars(map[string]interface{}{
			"rec": &ComposeRecord{value: raw},

			// nis is only for testing, gval does not recognise nil (or null) as a keyword!
			"nil": nil,

			// typed value but empty (this happens when you assign next value in expression but do not set a value to it
			"initRec": &ComposeRecord{},

			// same as &ComposeRecord{value: &types.Record{}},
			"validRecZero":    &ComposeRecord{value: &types.Record{ID: 0, Values: types.RecordValueSet{}}},
			"validRecValidID": &ComposeRecord{value: &types.Record{ID: 42, Values: types.RecordValueSet{}}},

			// "record" (not really) set to nil
			"fooRec": nil,

			// record with id and value set (this was) fixme
			"record": &ComposeRecord{value: &types.Record{ID: 99, Values: rawValues}},
		})
	)

	// @todo see note above regarding back-ref to record
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
			{raw, "rec.values"},
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

	t.Run("copy record's values via gval selector", func(t *testing.T) {
		var (
			req      = require.New(t)
			output   *expr.Vars
			inputRec = &types.Record{ID: 99, Values: rawValues}
			cloneRec = &types.Record{}
		)

		// IMPORTANT! set module to record; otherwise values will not be copied properly
		inputRec.SetModule(mod)
		cloneRec.SetModule(mod)

		e, err := aTypes.NewExpr(
			"recordClone.values",
			ComposeRecordValues{}.Type(),
			`record.values`)
		req.NoError(err)

		err = e.SetType(func(s string) (expr.Type, error) {
			return ComposeRecordValues{}, nil
		})
		req.NoError(err)

		evaluable, err := expr.NewParser().Parse(e.GetExpr())
		req.NoError(err)
		e.SetEval(evaluable)

		input, _ := expr.NewVars(map[string]expr.TypedValue{
			"record":      &ComposeRecord{value: inputRec},
			"recordClone": &ComposeRecord{value: cloneRec},
		})

		output, err = (aTypes.ExprSet{e}).Eval(context.Background(), input)
		req.NoError(err)
		rc, err := output.Select("recordClone")
		req.NoError(err)
		req.NotNil(rc.Get())
		req.Equal(rawValues, rc.Get().(*types.Record).Values)
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

			// {true, `rec.values.n1 < 3`}, // invalid op <nil> < 3
			{true, `rec.values.n2 < 3`},
			{true, `rec.values.n3 < 3`},

			// {false, `rec.values.n1 > 1`}, // invalid op <nil> > 3
			{false, `rec.values.n2 > 2`},
			{false, `rec.values.n3 > 2`},

			{true, `rec.values.ref1 != 2`},
			{true, `rec.values.ref2 == 2`},
			{true, `rec.values.ref2 == "2"`},

			{true, `isNil(nil)`},
			{true, `isNil(rec.values.d1)`},
			{true, `rec.values.n3 ? true : false`},
			{false, `rec.values.d1 ? true : false`},
			{true, `rec.values.d2 ? true : false`},
			{true, `rec.values.d1 ? true : false || true `},
			{false, `rec.values.d1 ? true : false && true `},
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

		req.NoError(assignToComposeRecordValues(target, initPather(req, expr.Path("a")), "b"))
		req.Len(target.Values, 1)
		req.True(target.Values.Has("a", 0))
		req.NoError(assignToComposeRecordValues(target, initPather(req, expr.Path("a[1]")), "b"))
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

		req.Error(assignToComposeRecordValues(target, initPather(req, expr.Path("a[2]")), expr.Must(expr.NewAny([]interface{}{"1", "2"}))))
		req.Len(target.Values, 0)

		req.NoError(assignToComposeRecordValues(target, initPather(req, expr.Path("a")), expr.Must(expr.NewAny([]interface{}{"1", "2"}))))
		req.Len(target.Values, 2)

		req.NoError(assignToComposeRecordValues(target, initPather(req, expr.Path("a")), expr.Must(expr.NewAny([]string{"1", "2"}))))
		req.Len(target.Values, 2)
	})

	t.Run("overwrite multiple values", func(t *testing.T) {
		var (
			req    = require.New(t)
			target = &types.Record{Values: types.RecordValueSet{
				&types.RecordValue{Name: "mval", Value: "203221389712893", Place: 0},
				&types.RecordValue{Name: "mval", Value: "203221389712893", Place: 1},
			}}

			scope, _ = expr.NewVars(map[string]interface{}{
				"rec": expr.Must(NewComposeRecord(target)),
			})

			ee = aTypes.ExprSet{
				&aTypes.Expr{Target: "rec.values.mval", Type: "Array", Expr: "[1,2,3,4]"},
			}
		)

		// convert all types
		req.NoError(ee[0].SetType(func(_ string) (expr.Type, error) { return expr.Array{}, nil }))

		// init expression handler
		evaluable, err := expr.NewParser().Parse(ee[0].GetExpr())
		req.NoError(err)
		ee[0].SetEval(evaluable)

		_, err = ee.Eval(context.Background(), scope)
		req.NoError(err)
		req.Len(target.Values, 4)
		req.Equal(target.Values.Get("mval", 0).Value, "1")
		req.Equal(target.Values.Get("mval", 1).Value, "2")
		req.Equal(target.Values.Get("mval", 2).Value, "3")
		req.Equal(target.Values.Get("mval", 3).Value, "4")
	})
}

func TestCastToComposeRecordValues(t *testing.T) {
	var (
		commonRVS = types.RecordValueSet{
			&types.RecordValue{Name: "bools", Value: "1"},
			&types.RecordValue{Name: "interfaces", Value: "val"},
			&types.RecordValue{Name: "interfaces", Value: "val2", Place: 1},
			&types.RecordValue{Name: "string", Value: "val"},
			&types.RecordValue{Name: "strings", Value: "val"},
			&types.RecordValue{Name: "strings", Value: "val2", Place: 1},
			// false booleans will not be set and that is ok!
			// &types.RecordValue{Name: "bools", Value: "false", Place: 1}
		}

		tt = time.Date(1999, 9, 9, 9, 9, 9, 9, time.UTC)

		nilSlice      []int
		nilUntypedMap map[string]interface{}

		cases = []struct {
			name string
			in   interface{}
			out  types.RecordValueSet
			err  error
		}{
			{
				in: map[string]interface{}{
					"string":        "val",
					"interfaces":    []interface{}{"val", "val2"},
					"strings":       []string{"val", "val2"},
					"bools":         []bool{true, false},
					"nilSlice":      nilSlice,
					"nilUntypedMap": nilUntypedMap,
				},
				// warning! false booleans will not be set and that is ok!
				out: commonRVS,
			},
			{
				in: map[string][]string{
					"strings": {"val", "val2"},
				},
				out: types.RecordValueSet{
					&types.RecordValue{Name: "strings", Value: "val"},
					&types.RecordValue{Name: "strings", Value: "val2", Place: 1},
				},
			},
			{
				in: map[string]string{
					"string": "val",
				},
				out: types.RecordValueSet{
					&types.RecordValue{Name: "string", Value: "val"},
				},
			},
			{
				in: &types.RecordValue{Name: "string", Value: "val"},
				out: types.RecordValueSet{
					&types.RecordValue{Name: "string", Value: "val"},
				},
			},
			{
				in: &types.RecordValue{Name: "datetime", Value: tt.String()},
				out: types.RecordValueSet{
					&types.RecordValue{Name: "datetime", Value: tt.String()},
				},
			},
			{
				in:  commonRVS,
				out: commonRVS,
			},
			{
				in:  &types.Record{Values: commonRVS},
				out: commonRVS,
			},
			{
				in:  42,
				err: fmt.Errorf("unable to cast type int to types.RecordValueSet"),
			},
		}
	)

	for _, c := range cases {
		if c.name == "" && c.in != nil {
			c.name = fmt.Sprintf("%T", c.in)
		}

		t.Run(c.name, func(t *testing.T) {
			var (
				req = require.New(t)
			)

			if c.in == nil {
				t.Skip()
			}

			if out, err := CastToComposeRecordValues(c.in); c.err != nil {
				req.EqualError(err, c.err.Error())
			} else {
				req.EqualValues(c.out, out)
			}
		})
	}
}

func TestIsNil(t *testing.T) {
	var (
		nilSlice        []int
		nilUntypedSlice []interface{}
		nilPtr          *int
		nilStringMap    map[string]string
		nilUntypedMap   map[string]interface{}

		cases = []struct {
			name     string
			input    interface{}
			expected bool
		}{
			{
				"nilSlice",
				nilSlice,
				true,
			},
			{
				"nilUntypedSlice",
				nilUntypedSlice,
				true,
			},
			{
				"nilPtr",
				nilPtr,
				true,
			},
			{
				"nilStringMap",
				nilStringMap,
				true,
			},
			{
				"nilUntypedMap",
				nilUntypedMap,
				true,
			},
		}
	)

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var (
				req = require.New(t)
			)

			req.Equal(c.expected, isNil(c.input))
		})
	}
}

func TestGvalStringFunctionParamsIntCasting(t *testing.T) {
	var (
		ctx = context.Background()
		req = require.New(t)

		scope = &expr.Vars{}
	)

	t.Run("substring", func(t *testing.T) {
		eval, err := expr.NewParser().Parse(`substring("foobar", 1, -1)`)
		req.NoError(err)

		res, err := eval.Eval(ctx, scope)
		req.NoError(err)

		req.Equal("oobar", res.(string))
	})

	t.Run("replace", func(t *testing.T) {
		eval, err := expr.NewParser().Parse(`replace(" foo foo foo", "foo", "bar", 1)`)
		req.NoError(err)

		res, err := eval.Eval(ctx, scope)
		req.NoError(err)

		req.Equal(" bar foo foo", res.(string))
	})

	t.Run("split", func(t *testing.T) {
		eval, err := expr.NewParser().Parse(`split("foo-bar", "-")`)
		req.NoError(err)

		res, err := eval.Eval(ctx, scope)
		req.NoError(err)

		req.Equal([]string{"foo", "bar"}, res.([]string))
	})

	t.Run("has", func(t *testing.T) {
		eval, err := expr.NewParser().Parse(`has([1,2], 2)`)
		req.NoError(err)

		res, err := eval.Eval(ctx, scope)
		req.NoError(err)

		req.Equal(true, res)
	})

	t.Run("repeat", func(t *testing.T) {
		eval, err := expr.NewParser().Parse(`repeat("! ", 3)`)
		req.NoError(err)

		res, err := eval.Eval(ctx, scope)
		req.NoError(err)

		req.Equal("! ! ! ", res.(string))
	})

	t.Run("shorten", func(t *testing.T) {
		eval, err := expr.NewParser().Parse(`shorten("This is a whole sentence", "word", 4)`)
		req.NoError(err)

		res, err := eval.Eval(ctx, scope)
		req.NoError(err)

		req.Equal("This is a whole …", res.(string))
	})

}

func TestRecordValues_Merge(t *testing.T) {
	var (
		req = require.New(t)

		rv  = &ComposeRecordValues{}
		foo = &ComposeRecordValues{
			value: &types.Record{
				Values: []*types.RecordValue{
					{
						Name:  "k1",
						Value: "testValue1",
					},
					{
						Name:  "k2",
						Value: "testValue2",
					},
				},
			},
		}
		bar = &ComposeRecordValues{
			value: &types.Record{
				Values: []*types.RecordValue{
					{
						Name:  "k3",
						Value: "testValue3",
					},
				},
			},
		}
		expected = &types.Record{
			Values: []*types.RecordValue{
				{
					Name:  "k1",
					Value: "testValue1",
				},
				{
					Name:  "k2",
					Value: "testValue2",
				},
				{
					Name:  "k3",
					Value: "testValue3",
				},
			},
		}
	)

	out, err := rv.Merge(foo, bar)
	req.NoError(err)
	req.Equal(expected, out.Get())
}

func TestRecordValues_Omit(t *testing.T) {
	var (
		req = require.New(t)

		rv = ComposeRecordValues{
			value: &types.Record{
				Values: []*types.RecordValue{
					{
						Name:  "k1",
						Value: "testValue1",
					},
					{
						Name:  "k2",
						Value: "testValue2",
					},
					{
						Name:  "k3",
						Value: "testValue3",
					},
				},
			},
		}
		expected = &types.Record{
			Values: []*types.RecordValue{
				{
					Name:  "k2",
					Value: "testValue2",
				},
			},
		}
	)

	out, err := rv.Delete("k1", "k3")
	req.NoError(err)
	req.Equal(expected, out.Get())
}

// goos: darwin
// goarch: arm64
// pkg: github.com/cortezaproject/corteza/server/compose/automation
// BenchmarkAssignToComposeRecordValues-12    	22131764	        53.22 ns/op	      96 B/op	       1 allocs/op
// PASS
func BenchmarkAssignToComposeRecordValues(b *testing.B) {
	var (
		req    = require.New(b)
		target = &types.Record{Values: types.RecordValueSet{}}
		p      = initPather(req, expr.Path("a"))
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		assignToComposeRecordValues(target, p, "b")
	}
}

// goos: darwin
// goarch: arm64
// pkg: github.com/cortezaproject/corteza/server/compose/automation
// BenchmarkRecordFieldValuesAccess-12    	 1947148	       593.6 ns/op	     544 B/op	      19 allocs/op
// PASS
func BenchmarkRecordFieldValuesAccess(b *testing.B) {
	var (
		mod = &types.Module{Fields: types.ModuleFieldSet{
			&types.ModuleField{Name: "s1", Multi: true, Kind: "String"},
		}}

		rawValues = types.RecordValueSet{
			&types.RecordValue{Name: "s1", Value: "sVal1.0", Place: 0},
			&types.RecordValue{Name: "s1", Value: "sVal1.1", Place: 1},
			&types.RecordValue{Name: "s1", Value: "sVal1.2", Place: 2},
		}
		raw = &types.Record{Values: rawValues}

		scope, _ = expr.NewVars(map[string]interface{}{
			"rec": &ComposeRecord{value: raw},
		})
	)

	// @todo see note above regarding back-ref to record
	raw.SetModule(mod)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		expr.Select(scope, "rec.values.s1[1]")
	}
}

func initPather(req *require.Assertions, p expr.Pather) (out expr.Pather) {
	err := p.Next()
	req.NoError(err)

	return p
}
