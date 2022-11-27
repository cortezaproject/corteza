package automation

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// automation/automation/jsenv_handler.yaml

import (
	"context"
	atypes "github.com/cortezaproject/corteza/server/automation/types"
	"github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/cortezaproject/corteza/server/pkg/wfexec"
	"io"
)

var _ wfexec.ExecResponse

type (
	jsenvHandlerRegistry interface {
		AddFunctions(ff ...*atypes.Function)
		Type(ref string) expr.Type
	}
)

func (h jsenvHandler) register() {
	h.reg.AddFunctions(
		h.Execute(),
	)
}

type (
	jsenvExecuteArgs struct {
		hasScope    bool
		Scope       interface{}
		scopeAny    interface{}
		scopeStream io.Reader

		hasSource bool
		Source    string
	}

	jsenvExecuteResults struct {
		ResultString string
		ResultInt    int64
		ResultBool   bool
		ResultAny    interface{}
	}
)

func (a jsenvExecuteArgs) GetScope() (bool, interface{}, io.Reader) {
	return a.hasScope, a.scopeAny, a.scopeStream
}

// Execute function Process arbitrary data in jsenv
//
// expects implementation of execute function:
//
//	func (h jsenvHandler) execute(ctx context.Context, args *jsenvExecuteArgs) (results *jsenvExecuteResults, err error) {
//	   return
//	}
func (h jsenvHandler) Execute() *atypes.Function {
	return &atypes.Function{
		Ref:    "jsenvExecute",
		Kind:   "function",
		Labels: map[string]string(nil),
		Meta: &atypes.FunctionMeta{
			Short: "Process arbitrary data in jsenv",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "scope",
				Types: []string{"Any", "Reader"}, Required: true,
			},
			{
				Name:  "source",
				Types: []string{"String"}, Required: true,
			},
		},

		Results: []*atypes.Param{

			{
				Name:  "resultString",
				Types: []string{"String"},
			},

			{
				Name:  "resultInt",
				Types: []string{"Integer"},
			},

			{
				Name:  "resultBool",
				Types: []string{"Boolean"},
			},

			{
				Name:  "resultAny",
				Types: []string{"Any"},
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &jsenvExecuteArgs{
					hasScope:  in.Has("scope"),
					hasSource: in.Has("source"),
				}
			)

			if err = in.Decode(args); err != nil {
				return
			}

			// Converting Scope argument
			if args.hasScope {
				aux := expr.Must(expr.Select(in, "scope"))
				switch aux.Type() {
				case h.reg.Type("Any").Type():
					args.scopeAny = aux.Get().(interface{})
				case h.reg.Type("Reader").Type():
					args.scopeStream = aux.Get().(io.Reader)
				}
			}

			var results *jsenvExecuteResults
			if results, err = h.execute(ctx, args); err != nil {
				return
			}

			out = &expr.Vars{}

			{
				// converting results.ResultString (string) to String
				var (
					tval expr.TypedValue
				)

				if tval, err = h.reg.Type("String").Cast(results.ResultString); err != nil {
					return
				} else if err = expr.Assign(out, "resultString", tval); err != nil {
					return
				}
			}

			{
				// converting results.ResultInt (int64) to Integer
				var (
					tval expr.TypedValue
				)

				if tval, err = h.reg.Type("Integer").Cast(results.ResultInt); err != nil {
					return
				} else if err = expr.Assign(out, "resultInt", tval); err != nil {
					return
				}
			}

			{
				// converting results.ResultBool (bool) to Boolean
				var (
					tval expr.TypedValue
				)

				if tval, err = h.reg.Type("Boolean").Cast(results.ResultBool); err != nil {
					return
				} else if err = expr.Assign(out, "resultBool", tval); err != nil {
					return
				}
			}

			{
				// converting results.ResultAny (interface{}) to Any
				var (
					tval expr.TypedValue
				)

				if tval, err = h.reg.Type("Any").Cast(results.ResultAny); err != nil {
					return
				} else if err = expr.Assign(out, "resultAny", tval); err != nil {
					return
				}
			}

			return
		},
	}
}
