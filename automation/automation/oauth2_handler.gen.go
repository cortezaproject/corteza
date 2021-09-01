package automation

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// automation/automation/oauth2_handler.yaml

import (
	"context"
	atypes "github.com/cortezaproject/corteza-server/automation/types"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/pkg/wfexec"
)

var _ wfexec.ExecResponse

type (
	oauth2HandlerRegistry interface {
		AddFunctions(ff ...*atypes.Function)
		Type(ref string) expr.Type
	}
)

func (h oauth2Handler) register() {
	h.reg.AddFunctions(
		h.Authenticate(),
	)
}

type (
	oauth2AuthenticateArgs struct {
		hasClient bool
		Client    string

		hasSecret bool
		Secret    string

		hasScope bool
		Scope    string

		hasTokenUrl bool
		TokenUrl    string
	}

	oauth2AuthenticateResults struct {
		AccessToken  string
		RefreshToken string
		Token        interface{}
	}
)

// Authenticate function Authentication: OAUTH2
//
// expects implementation of authenticate function:
// func (h oauth2Handler) authenticate(ctx context.Context, args *oauth2AuthenticateArgs) (results *oauth2AuthenticateResults, err error) {
//    return
// }
func (h oauth2Handler) Authenticate() *atypes.Function {
	return &atypes.Function{
		Ref:    "oauth2Authenticate",
		Kind:   "function",
		Labels: map[string]string(nil),
		Meta: &atypes.FunctionMeta{
			Short: "Authentication: OAUTH2",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "client",
				Types: []string{"String"}, Required: true,
			},
			{
				Name:  "secret",
				Types: []string{"String"}, Required: true,
			},
			{
				Name:  "scope",
				Types: []string{"String"}, Required: true,
			},
			{
				Name:  "tokenUrl",
				Types: []string{"String"}, Required: true,
			},
		},

		Results: []*atypes.Param{

			{
				Name:  "accessToken",
				Types: []string{"String"},
			},

			{
				Name:  "refreshToken",
				Types: []string{"String"},
			},

			{
				Name:  "token",
				Types: []string{"Any"},
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &oauth2AuthenticateArgs{
					hasClient:   in.Has("client"),
					hasSecret:   in.Has("secret"),
					hasScope:    in.Has("scope"),
					hasTokenUrl: in.Has("tokenUrl"),
				}
			)

			if err = in.Decode(args); err != nil {
				return
			}

			var results *oauth2AuthenticateResults
			if results, err = h.authenticate(ctx, args); err != nil {
				return
			}

			out = &expr.Vars{}

			{
				// converting results.AccessToken (string) to String
				var (
					tval expr.TypedValue
				)

				if tval, err = h.reg.Type("String").Cast(results.AccessToken); err != nil {
					return
				} else if err = expr.Assign(out, "accessToken", tval); err != nil {
					return
				}
			}

			{
				// converting results.RefreshToken (string) to String
				var (
					tval expr.TypedValue
				)

				if tval, err = h.reg.Type("String").Cast(results.RefreshToken); err != nil {
					return
				} else if err = expr.Assign(out, "refreshToken", tval); err != nil {
					return
				}
			}

			{
				// converting results.Token (interface{}) to Any
				var (
					tval expr.TypedValue
				)

				if tval, err = h.reg.Type("Any").Cast(results.Token); err != nil {
					return
				} else if err = expr.Assign(out, "token", tval); err != nil {
					return
				}
			}

			return
		},
	}
}
