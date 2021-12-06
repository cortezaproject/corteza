package automation

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// automation/automation/jwt_handler.yaml

import (
	"context"
	atypes "github.com/cortezaproject/corteza-server/automation/types"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/pkg/wfexec"
	"io"
)

var _ wfexec.ExecResponse

type (
	jwtHandlerRegistry interface {
		AddFunctions(ff ...*atypes.Function)
		Type(ref string) expr.Type
	}
)

func (h jwtHandler) register() {
	h.reg.AddFunctions(
		h.Generate(),
	)
}

type (
	jwtGenerateArgs struct {
		hasScope bool
		Scope    string

		hasHeader    bool
		Header       interface{}
		headerVars   map[string]expr.TypedValue
		headerString string

		hasPayload    bool
		Payload       interface{}
		payloadVars   map[string]expr.TypedValue
		payloadString string

		hasSecret    bool
		Secret       interface{}
		secretString string
		secretStream io.Reader
	}

	jwtGenerateResults struct {
		Token string
	}
)

func (a jwtGenerateArgs) GetHeader() (bool, map[string]expr.TypedValue, string) {
	return a.hasHeader, a.headerVars, a.headerString
}

func (a jwtGenerateArgs) GetPayload() (bool, map[string]expr.TypedValue, string) {
	return a.hasPayload, a.payloadVars, a.payloadString
}

func (a jwtGenerateArgs) GetSecret() (bool, string, io.Reader) {
	return a.hasSecret, a.secretString, a.secretStream
}

// Generate function Generate JWT
//
// expects implementation of generate function:
// func (h jwtHandler) generate(ctx context.Context, args *jwtGenerateArgs) (results *jwtGenerateResults, err error) {
//    return
// }
func (h jwtHandler) Generate() *atypes.Function {
	return &atypes.Function{
		Ref:    "jwtGenerate",
		Kind:   "function",
		Labels: map[string]string(nil),
		Meta: &atypes.FunctionMeta{
			Short: "Generate JWT",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "scope",
				Types: []string{"String"},
			},
			{
				Name:  "header",
				Types: []string{"Vars", "String"}, Required: true,
			},
			{
				Name:  "payload",
				Types: []string{"Vars", "String"}, Required: true,
			},
			{
				Name:  "secret",
				Types: []string{"String", "Reader"}, Required: true,
			},
		},

		Results: []*atypes.Param{

			{
				Name:  "token",
				Types: []string{"String"},
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &jwtGenerateArgs{
					hasScope:   in.Has("scope"),
					hasHeader:  in.Has("header"),
					hasPayload: in.Has("payload"),
					hasSecret:  in.Has("secret"),
				}
			)

			if err = in.Decode(args); err != nil {
				return
			}

			// Converting Header argument
			if args.hasHeader {
				aux := expr.Must(expr.Select(in, "header"))
				switch aux.Type() {
				case h.reg.Type("Vars").Type():
					args.headerVars = aux.Get().(map[string]expr.TypedValue)
				case h.reg.Type("String").Type():
					args.headerString = aux.Get().(string)
				}
			}

			// Converting Payload argument
			if args.hasPayload {
				aux := expr.Must(expr.Select(in, "payload"))
				switch aux.Type() {
				case h.reg.Type("Vars").Type():
					args.payloadVars = aux.Get().(map[string]expr.TypedValue)
				case h.reg.Type("String").Type():
					args.payloadString = aux.Get().(string)
				}
			}

			// Converting Secret argument
			if args.hasSecret {
				aux := expr.Must(expr.Select(in, "secret"))
				switch aux.Type() {
				case h.reg.Type("String").Type():
					args.secretString = aux.Get().(string)
				case h.reg.Type("Reader").Type():
					args.secretStream = aux.Get().(io.Reader)
				}
			}

			var results *jwtGenerateResults
			if results, err = h.generate(ctx, args); err != nil {
				return
			}

			out = &expr.Vars{}

			{
				// converting results.Token (string) to String
				var (
					tval expr.TypedValue
				)

				if tval, err = h.reg.Type("String").Cast(results.Token); err != nil {
					return
				} else if err = expr.Assign(out, "token", tval); err != nil {
					return
				}
			}

			return
		},
	}
}
