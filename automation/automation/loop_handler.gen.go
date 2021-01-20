package automation

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// automation/automation/loop_handler.yaml

import (
	"context"
	atypes "github.com/cortezaproject/corteza-server/automation/types"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/pkg/wfexec"
	"io"
)

var _ wfexec.ExecResponse

type (
	loopHandlerRegistry interface {
		AddFunctions(ff ...*atypes.Function)
		Type(ref string) expr.Type
	}
)

func (h loopHandler) register() {
	h.reg.AddFunctions(
		h.Sequence(),
		h.Do(),
		h.Each(),
		h.Lines(),
	)
}

type (
	loopSequenceArgs struct {
		hasFirst bool
		First    int64

		hasLast bool
		Last    int64

		hasStep bool
		Step    int64
	}

	loopSequenceResults struct {
		Counter int64
		IsFirst bool
		IsLast  bool
	}
)

// Sequence function Iterates over sequence of numbers
//
// expects implementation of sequence function:
// func (h loopHandler) sequence(ctx context.Context, args *loopSequenceArgs) (results *loopSequenceResults, err error) {
//    return
// }
func (h loopHandler) Sequence() *atypes.Function {
	return &atypes.Function{
		Ref:  "loopSequence",
		Kind: "iterator",
		Meta: &atypes.FunctionMeta{
			Short: "Iterates over sequence of numbers",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "first",
				Types: []string{"Integer"},
			},
			{
				Name:  "last",
				Types: []string{"Integer"},
			},
			{
				Name:  "step",
				Types: []string{"Integer"},
			},
		},

		Results: []*atypes.Param{

			{
				Name:  "counter",
				Types: []string{"Integer"},
			},

			{
				Name:  "isFirst",
				Types: []string{"Boolean"},
			},

			{
				Name:  "isLast",
				Types: []string{"Boolean"},
			},
		},

		Iterator: func(ctx context.Context, in *expr.Vars) (out wfexec.IteratorHandler, err error) {
			var (
				args = &loopSequenceArgs{
					hasFirst: in.Has("first"),
					hasLast:  in.Has("last"),
					hasStep:  in.Has("step"),
				}
			)

			if err = in.Decode(args); err != nil {
				return
			}

			return h.sequence(ctx, args)
		},
	}
}

type (
	loopDoArgs struct {
		hasWhile bool
		While    string
	}
)

// Do function Iterates while condition is true
//
// expects implementation of do function:
// func (h loopHandler) do(ctx context.Context, args *loopDoArgs) (err error) {
//    return
// }
func (h loopHandler) Do() *atypes.Function {
	return &atypes.Function{
		Ref:  "loopDo",
		Kind: "iterator",
		Meta: &atypes.FunctionMeta{
			Short: "Iterates while condition is true",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "while",
				Types: []string{"String"}, Required: true,
				Meta: &atypes.ParamMeta{
					Label:       "Expression tested before each iteration",
					Description: "Expression to be evaluated each iteration; loop will continue until expression is true",
				},
			},
		},

		Iterator: func(ctx context.Context, in *expr.Vars) (out wfexec.IteratorHandler, err error) {
			var (
				args = &loopDoArgs{
					hasWhile: in.Has("while"),
				}
			)

			if err = in.Decode(args); err != nil {
				return
			}

			return h.do(ctx, args)
		},
	}
}

type (
	loopEachArgs struct {
		hasItems bool
		Items    []expr.TypedValue
	}

	loopEachResults struct {
		Item interface{}
	}
)

// Each function Iterates over set of items
//
// expects implementation of each function:
// func (h loopHandler) each(ctx context.Context, args *loopEachArgs) (results *loopEachResults, err error) {
//    return
// }
func (h loopHandler) Each() *atypes.Function {
	return &atypes.Function{
		Ref:  "loopEach",
		Kind: "iterator",
		Meta: &atypes.FunctionMeta{
			Short: "Iterates over set of items",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "items",
				Types: []string{"Any"}, Required: true, IsArray: true,
			},
		},

		Results: []*atypes.Param{

			{
				Name:  "item",
				Types: []string{"Any"},
			},
		},

		Iterator: func(ctx context.Context, in *expr.Vars) (out wfexec.IteratorHandler, err error) {
			var (
				args = &loopEachArgs{
					hasItems: in.Has("items"),
				}
			)

			if err = in.Decode(args); err != nil {
				return
			}

			return h.each(ctx, args)
		},
	}
}

type (
	loopLinesArgs struct {
		hasStream bool
		Stream    io.Reader
	}

	loopLinesResults struct {
		Line string
	}
)

// Lines function Iterates over lines from stream
//
// expects implementation of lines function:
// func (h loopHandler) lines(ctx context.Context, args *loopLinesArgs) (results *loopLinesResults, err error) {
//    return
// }
func (h loopHandler) Lines() *atypes.Function {
	return &atypes.Function{
		Ref:  "loopLines",
		Kind: "iterator",
		Meta: &atypes.FunctionMeta{
			Short: "Iterates over lines from stream",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "stream",
				Types: []string{"Reader"}, Required: true,
			},
		},

		Results: []*atypes.Param{

			{
				Name:  "line",
				Types: []string{"String"},
			},
		},

		Iterator: func(ctx context.Context, in *expr.Vars) (out wfexec.IteratorHandler, err error) {
			var (
				args = &loopLinesArgs{
					hasStream: in.Has("stream"),
				}
			)

			if err = in.Decode(args); err != nil {
				return
			}

			return h.lines(ctx, args)
		},
	}
}
