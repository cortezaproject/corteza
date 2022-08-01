package automation

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// system/automation/valuestore_handler.yaml

import (
	"context"
	atypes "github.com/cortezaproject/corteza-server/automation/types"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/pkg/wfexec"
)

var _ wfexec.ExecResponse

type (
	valuestoreHandlerRegistry interface {
		AddFunctions(ff ...*atypes.Function)
		Type(ref string) expr.Type
	}
)

func (h valuestoreHandler) register() {
	h.reg.AddFunctions(
		h.Env(),
	)
}

type (
	valuestoreEnvArgs struct {
		hasKey bool
		Key    string
	}

	valuestoreEnvResults struct {
		Value interface{}
	}
)

// Env function Get ENV variable
//
// expects implementation of env function:
// func (h valuestoreHandler) env(ctx context.Context, args *valuestoreEnvArgs) (results *valuestoreEnvResults, err error) {
//    return
// }
func (h valuestoreHandler) Env() *atypes.Function {
	return &atypes.Function{
		Ref:    "valuestoreEnv",
		Kind:   "function",
		Labels: map[string]string{"templates": "step,workflow"},
		Meta: &atypes.FunctionMeta{
			Short:       "Get ENV variable",
			Description: "Get ENV variable for the specified key. If the key doesn't correspond to any value, nil is returned",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "key",
				Types: []string{"String"}, Required: true,
			},
		},

		Results: []*atypes.Param{

			{
				Name:  "value",
				Types: []string{"Any"},
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &valuestoreEnvArgs{
					hasKey: in.Has("key"),
				}
			)

			if err = in.Decode(args); err != nil {
				return
			}

			var results *valuestoreEnvResults
			if results, err = h.env(ctx, args); err != nil {
				return
			}

			out = &expr.Vars{}

			{
				// converting results.Value (interface{}) to Any
				var (
					tval expr.TypedValue
				)

				if tval, err = h.reg.Type("Any").Cast(results.Value); err != nil {
					return
				} else if err = expr.Assign(out, "value", tval); err != nil {
					return
				}
			}

			return
		},
	}
}
