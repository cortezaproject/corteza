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
)

func (h jsonParserHandler) Parse() *atypes.Function {
	return &atypes.Function{
		Ref:    "jsonParse",
		Kind:   "function",
		Labels: map[string]string{"json": "step", "parser": "step"},
		Meta: &atypes.FunctionMeta{
			Short:       "JSON Parse string into a JSON object",
			Description: "Takes a JSON string and returns the parsed JSON value. The result can be an object or array depending on the input.",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "jsonText",
				Types: []string{"String"}, Required: true,
			},
		},

		Results: []*atypes.Param{
			{
				Name:  "result",
				Types: []string{"Any"},
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

			var results *jsonParserParseResults
			if results, err = h.parse(ctx, args); err != nil {
				return
			}

			out = &expr.Vars{}

			{
				// converting results.Result (interface{}) to Any
				var (
					tval expr.TypedValue
				)

				if tval, err = h.reg.Type("Any").Cast(results.Result); err != nil {
					return
				} else if err = expr.Assign(out, "result", tval); err != nil {
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

func (h jsonParserHandler) Stringify() *atypes.Function {
	return &atypes.Function{
		Ref:    "jsonStringify",
		Kind:   "function",
		Labels: map[string]string{"json": "step", "stringify": "step"},
		Meta: &atypes.FunctionMeta{
			Short:       "JSON Convert a object to a string",
			Description: "Takes a JSON object (Any type) and returns a JSON string representation.",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "jsonObject",
				Types: []string{"Any"}, Required: true,
			},
		},

		Results: []*atypes.Param{
			{
				Name:  "result",
				Types: []string{"String"},
			},
			{
				Name:  "error",
				Types: []string{"String"},
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &jsonParserStringifyArgs{
					hasJsonObject: in.Has("jsonObject"),
				}
			)

			if err = in.Decode(args); err != nil {
				return
			}

			// Converting jsonObject argument
			if args.hasJsonObject {
				aux := expr.Must(expr.Select(in, "jsonObject"))
				args.JsonObject = aux.Get()
			}

			var results *jsonParserStringifyResults
			if results, err = h.stringify(ctx, args); err != nil {
				return
			}

			out = &expr.Vars{}

			{
				// converting results.Result (string) to String
				var (
					tval expr.TypedValue
				)

				if tval, err = h.reg.Type("String").Cast(results.Result); err != nil {
					return
				} else if err = expr.Assign(out, "result", tval); err != nil {
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

func (h jsonParserHandler) Template() *atypes.Function {
	return &atypes.Function{
		Ref:    "jsonTemplate",
		Kind:   "function",
		Labels: map[string]string{"json": "step", "template": "step"},
		Meta: &atypes.FunctionMeta{
			Short:       "JSON Template with variable substitution",
			Description: "Takes a JSON template string with ${var.key} placeholders and a vars map. Substitutes the placeholders with values from vars and returns the final JSON string. Supports dot-notation for nested values (e.g., ${var.filename}, ${var.data.content}).",
		},

		Parameters: []*atypes.Param{
			{
				Name:     "template",
				Types:    []string{"String"},
				Required: true,
			},
			{
				Name:     "vars",
				Types:    []string{"Vars"},
				Required: false,
			},
		},

		Results: []*atypes.Param{
			{
				Name:  "result",
				Types: []string{"String"},
			},
			{
				Name:  "error",
				Types: []string{"String"},
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &jsonParserTemplateArgs{
					hasVars: in.Has("vars"),
				}
			)

			// Extract template manually
			if in.Has("template") {
				aux := expr.Must(expr.Select(in, "template"))
				if t, ok := aux.Get().(string); ok {
					args.Template = t
				}
			}

			// Extract vars manually, unwrapping TypedValue
			if args.hasVars {
				aux := expr.Must(expr.Select(in, "vars"))
				switch m := aux.Get().(type) {
				case *expr.Vars:
					args.Vars = make(map[string]interface{})
					m.Each(func(k string, v expr.TypedValue) error {
						args.Vars[k] = v.Get()
						return nil
					})
				case map[string]expr.TypedValue:
					args.Vars = make(map[string]interface{})
					for k, v := range m {
						args.Vars[k] = v.Get()
					}
				}
			}

			var results *jsonParserTemplateResults
			if results, err = h.template(ctx, args); err != nil {
				return
			}

			out = &expr.Vars{}

			{
				// converting results.Result (string) to String
				var (
					tval expr.TypedValue
				)

				if tval, err = h.reg.Type("String").Cast(results.Result); err != nil {
					return
				} else if err = expr.Assign(out, "result", tval); err != nil {
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
