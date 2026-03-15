package automation

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// automation/automation/json_parser_handler.yaml

import (
	"context"

	atypes "github.com/cortezaproject/corteza/server/automation/types"
	"github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/cortezaproject/corteza/server/pkg/wfexec"
)

var _ wfexec.ExecResponse

type (
	modulesHandlerRegistry interface {
		AddFunctions(ff ...*atypes.Function)
		Type(ref string) expr.Type
	}
)

func (h jsonParserHandler) Parse() *atypes.Function {
	return &atypes.Function{
		Ref:    "jsonParse",
		Kind:   "function",
		Labels: map[string]string{"json": "step", "parser": "step"},
		Meta: &atypes.FunctionMeta{
			Short:       "Parse JSON text and extract fields",
			Description: "Parses JSON text (supports partial JSON) and extracts specified fields. If fields are not specified, returns entire JSON object.",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "jsonText",
				Types: []string{"String"}, Required: true,
			},
			{
				Name:  "fields",
				Types: []string{"Slice"}, Required: false,
			},
		},

		Results: []*atypes.Param{
			{
				Name:  "data",
				Types: []string{"Object"},
			},
			{
				Name:  "error",
				Types: []string{"String"},
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &jsonParserParseArgs{
					hasJsonText: in.Has("jsonText"),
					hasFields:   in.Has("fields"),
				}
			)

			if err = in.Decode(args); err != nil {
				return
			}

			// Converting jsonText argument
			if args.hasJsonText {
				aux := expr.Must(expr.Select(in, "jsonText"))
				if t, ok := aux.Get().(string); ok {
					args.JsonText = t
				}
			}

			// Converting fields argument
			if args.hasFields {
				aux := expr.Must(expr.Select(in, "fields"))
				if sliceVal, ok := aux.Get().([]interface{}); ok {
					args.Fields = make([]string, len(sliceVal))
					for i, v := range sliceVal {
						if s, ok := v.(string); ok {
							args.Fields[i] = s
						}
					}
				}
			}

			var results *jsonParserParseResults
			if results, err = h.parse(ctx, args); err != nil {
				return
			}

			out = &expr.Vars{}

			{
				// converting results.Data (map[string]interface{}) to Object
				var (
					tval expr.TypedValue
				)

				if tval, err = h.reg.Type("Object").Cast(results.Data); err != nil {
					return
				} else if err = expr.Assign(out, "data", tval); err != nil {
					return
				}
			}
			{
				// converting results.Error (string) to String
				var (
					tval expr.TypedValue
				)

				if tval, err = h.reg.Type("String").Cast(results.Error); err != nil {
					return
				} else if err = expr.Assign(out, "error", tval); err != nil {
					return
				}
			}

			return
		},
	}
}
