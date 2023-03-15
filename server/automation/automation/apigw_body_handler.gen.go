package automation

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// automation/automation/apigw_body_handler.yaml

import (
	"context"
	atypes "github.com/cortezaproject/corteza/server/automation/types"
	"github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/cortezaproject/corteza/server/pkg/http"
	"github.com/cortezaproject/corteza/server/pkg/wfexec"
)

var _ wfexec.ExecResponse

type (
	apigwBodyHandlerRegistry interface {
		AddFunctions(ff ...*atypes.Function)
		Type(ref string) expr.Type
	}
)

func (h apigwBodyHandler) register() {
	h.reg.AddFunctions(
		h.Read(),
	)
}

type (
	apigwBodyReadArgs struct {
		hasRequest bool
		Request    *http.Request
	}

	apigwBodyReadResults struct {
		Body string
	}
)

// Read function Read request body from integration gateway
//
// expects implementation of read function:
//
//	func (h apigwBodyHandler) read(ctx context.Context, args *apigwBodyReadArgs) (results *apigwBodyReadResults, err error) {
//	   return
//	}
func (h apigwBodyHandler) Read() *atypes.Function {
	return &atypes.Function{
		Ref:    "apigwBodyRead",
		Kind:   "function",
		Labels: map[string]string(nil),
		Meta: &atypes.FunctionMeta{
			Short: "Read request body from integration gateway",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "request",
				Types: []string{"HttpRequest"}, Required: true,
			},
		},

		Results: []*atypes.Param{

			{
				Name:  "body",
				Types: []string{"String"},
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &apigwBodyReadArgs{
					hasRequest: in.Has("request"),
				}
			)

			if err = in.Decode(args); err != nil {
				return
			}

			var results *apigwBodyReadResults
			if results, err = h.read(ctx, args); err != nil {
				return
			}

			out = &expr.Vars{}

			{
				// converting results.Body (string) to String
				var (
					tval expr.TypedValue
				)

				if tval, err = h.reg.Type("String").Cast(results.Body); err != nil {
					return
				} else if err = expr.Assign(out, "body", tval); err != nil {
					return
				}
			}

			return
		},
	}
}
