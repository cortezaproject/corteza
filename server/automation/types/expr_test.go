package types

import (
	"context"
	. "github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestExprSet_Eval(t *testing.T) {
	var (
		ctx = context.Background()

		cc = []struct {
			name   string
			set    ExprSet
			input  map[string]interface{}
			output map[string]interface{}
			err    string
		}{
			{
				name:   "empty",
				set:    ExprSet{},
				output: make(map[string]interface{}),
			},
			{
				name:   "constant assignment",
				set:    ExprSet{&Expr{Target: "foo", Expr: `"bar"`}},
				output: map[string]interface{}{"foo": Must(NewString("bar"))},
			},
			{
				name:   "vars with path",
				set:    ExprSet{&Expr{Target: "l1.l2", Expr: `"bar"`}},
				input:  map[string]interface{}{"l1": map[string]interface{}{}},
				output: map[string]interface{}{"l1": map[string]interface{}{"l2": Must(Typify("bar"))}},
			},
			{
				name: "copy vars with same types",
				set: ExprSet{
					&Expr{Target: "aa", Value: "vv", typ: &String{}},
					&Expr{Target: "bb", Source: "aa", typ: &String{}},
				},
				output: map[string]interface{}{
					"aa": Must(NewString("vv")),
					"bb": Must(NewString("vv")),
				},
			},
			{
				name: "copy var with type",
				set: ExprSet{
					&Expr{Target: "aa", Value: "should be always String", typ: &String{}},
					&Expr{Target: "bb", Source: "aa"},
				},
				output: map[string]interface{}{
					"aa": Must(NewString("should be always String")),
					"bb": Must(NewString("should be always String")),
				},
			},
			{
				name: "copy var to target with type",
				set: ExprSet{
					&Expr{Target: "aa", Value: "42", typ: &String{}},
					&Expr{Target: "bb", Source: "aa", typ: &Integer{}},
				},
				output: map[string]interface{}{
					"aa": Must(NewString("42")),
					"bb": Must(NewInteger(42)),
				},
			},
			{
				name: "assign into incompatible",
				set: ExprSet{
					&Expr{Target: "aa", Value: "foo", typ: &String{}},
					&Expr{Target: "bb", Source: "aa", typ: &Integer{}},
				},
				err: "unable to cast \"foo\" of type string to int64",
			},
			{
				name: "deep set into generated type",
				set: ExprSet{
					&Expr{Target: "a", typ: &KV{}},
					&Expr{Target: "a.b", Value: "c", typ: &String{}},
				},
				output: map[string]interface{}{
					"a": Must(NewKV(map[string]string{
						"b": "c",
					})),
				},
			},
			{
				name: "slice push",
				set: ExprSet{
					&Expr{Target: "arr", typ: &Array{}, Expr: `[]`},
					&Expr{Target: "arr", typ: &Array{}, Expr: `push(arr, "foo")`},
				},
				output: map[string]interface{}{
					"arr": Must(NewArray([]string{
						"foo",
					})),
				},
			},
			{
				name: "slice push w/o type",
				set: ExprSet{
					&Expr{Target: "arr", typ: &Array{}, Expr: `[]`},
					&Expr{Target: "arr", typ: &Any{}, Expr: `push(arr, "foo")`},
				},
				output: map[string]interface{}{
					"arr": Must(NewArray([]string{
						"foo",
					})),
				},
			},
			{
				name: "slice create & push w/o type",
				set: ExprSet{
					&Expr{Target: "arr", typ: &Array{}, Expr: `push([], "foo")`},
				},
				output: map[string]interface{}{
					"arr": Must(NewArray([]TypedValue{
						Must(Typify("foo")),
					})),
				},
			},
			{
				name:   "vars with nested path",
				set:    ExprSet{&Expr{Target: "t1.t11", typ: &String{}, Expr: `e1.e11`}},
				input:  map[string]interface{}{"t1": map[string]interface{}{}, "e1": map[string]interface{}{"e11": Must(Typify("bar"))}},
				output: map[string]interface{}{"t1": map[string]interface{}{"t11": Must(Typify("bar"))}},
			},
		}
	)

	for _, c := range cc {
		t.Run(c.name, func(t *testing.T) {
			var (
				req = require.New(t)
			)

			for _, e := range c.set {
				if e.Expr != "" {
					req.NoError(NewGvalParser().ParseEvaluators(e))
				}

				if e.typ == nil {
					e.typ = Any{}
				}
			}

			var (
				aux, _      = NewVars(c.input)
				output, err = c.set.Eval(ctx, aux)
			)

			if c.err == "" {
				req.NoError(err)
			} else {
				req.Error(err, c.err)
				return
			}

			req.Equal(Must(Typify(c.output)), output)
		})
	}
}
