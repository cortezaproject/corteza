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
	"io"
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
		h.ReadFile(),
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

type (
	apigwBodyReadFileArgs struct {
		hasRequest bool
		Request    *http.Request

		hasName bool
		Name    string
	}

	apigwBodyReadFileResults struct {
		File     io.Reader
		FileName string
		Exists   bool
	}
)

// ReadFile function Read file from integration gateway
//
// expects implementation of readFile function:
//
//	func (h apigwBodyHandler) readFile(ctx context.Context, args *apigwBodyReadFileArgs) (results *apigwBodyReadFileResults, err error) {
//	   return
//	}
func (h apigwBodyHandler) ReadFile() *atypes.Function {
	return &atypes.Function{
		Ref:    "apigwBodyReadFile",
		Kind:   "function",
		Labels: map[string]string(nil),
		Meta: &atypes.FunctionMeta{
			Short: "Read file from integration gateway",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "request",
				Types: []string{"HttpRequest"}, Required: true,
			},
			{
				Name:  "name",
				Types: []string{"String"}, Required: true,
			},
		},

		Results: []*atypes.Param{

			{
				Name:  "file",
				Types: []string{"Reader"},
			},

			{
				Name:  "fileName",
				Types: []string{"String"},
			},

			{
				Name:  "exists",
				Types: []string{"Boolean"},
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &apigwBodyReadFileArgs{
					hasRequest: in.Has("request"),
					hasName:    in.Has("name"),
				}
			)

			if err = in.Decode(args); err != nil {
				return
			}

			var results *apigwBodyReadFileResults
			if results, err = h.readFile(ctx, args); err != nil {
				return
			}

			out = &expr.Vars{}

			{
				// converting results.File (io.Reader) to Reader
				var (
					tval expr.TypedValue
				)

				if tval, err = h.reg.Type("Reader").Cast(results.File); err != nil {
					return
				} else if err = expr.Assign(out, "file", tval); err != nil {
					return
				}
			}

			{
				// converting results.FileName (string) to String
				var (
					tval expr.TypedValue
				)

				if tval, err = h.reg.Type("String").Cast(results.FileName); err != nil {
					return
				} else if err = expr.Assign(out, "fileName", tval); err != nil {
					return
				}
			}

			{
				// converting results.Exists (bool) to Boolean
				var (
					tval expr.TypedValue
				)

				if tval, err = h.reg.Type("Boolean").Cast(results.Exists); err != nil {
					return
				} else if err = expr.Assign(out, "exists", tval); err != nil {
					return
				}
			}

			return
		},
	}
}
