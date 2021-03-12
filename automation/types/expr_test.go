package types

import (
	"context"
	. "github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestExprSet_Eval(t *testing.T) {
	var (
		ctx = context.Background()

		cc = []struct {
			name   string
			set    ExprSet
			input  RVars
			output RVars
			err    string
		}{
			{
				name:   "empty",
				set:    ExprSet{},
				output: nil,
			},
			{
				name:   "constant assignment",
				set:    ExprSet{&Expr{Target: "foo", Expr: `"bar"`}},
				output: RVars{"foo": Must(NewAny("bar"))},
			},
			{
				name:   "vars with path",
				set:    ExprSet{&Expr{Target: "l1.l2", Expr: `"bar"`}},
				input:  RVars{"l1": RVars{}.Vars()},
				output: RVars{"l1": RVars{"l2": Must(NewAny("bar"))}.Vars()},
			},
			{
				name: "copy vars with same types",
				set: ExprSet{
					&Expr{Target: "aa", Value: "vv", typ: &String{}},
					&Expr{Target: "bb", Source: "aa", typ: &String{}},
				},
				output: RVars{
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
				output: RVars{
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
				output: RVars{
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
				output: RVars{
					"a": Must(NewKV(map[string]string{
						"b": "c",
					})),
				},
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

			req.Equal(c.output.Vars(), output)
		})
	}
}
