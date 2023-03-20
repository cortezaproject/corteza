package automation

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// automation/automation/queue_handler.yaml

import (
	"context"
	atypes "github.com/cortezaproject/corteza/server/automation/types"
	"github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/cortezaproject/corteza/server/pkg/wfexec"
	"io"
)

var _ wfexec.ExecResponse

type (
	queueHandlerRegistry interface {
		AddFunctions(ff ...*atypes.Function)
		Type(ref string) expr.Type
	}
)

func (h queueHandler) register() {
	h.reg.AddFunctions(
		h.Write(),
	)
}

type (
	queueWriteArgs struct {
		hasPayload    bool
		Payload       interface{}
		payloadString string
		payloadStream io.Reader

		hasQueue bool
		Queue    string
	}
)

func (a queueWriteArgs) GetPayload() (bool, string, io.Reader) {
	return a.hasPayload, a.payloadString, a.payloadStream
}

// Write function Queue message send
//
// expects implementation of write function:
// func (h queueHandler) write(ctx context.Context, args *queueWriteArgs) (err error) {
//    return
// }
func (h queueHandler) Write() *atypes.Function {
	return &atypes.Function{
		Ref:    "queueWrite",
		Kind:   "function",
		Labels: map[string]string(nil),
		Meta: &atypes.FunctionMeta{
			Short: "Queue message send",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "payload",
				Types: []string{"String", "Reader"}, Required: true,
			},
			{
				Name:  "queue",
				Types: []string{"String"}, Required: true,
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &queueWriteArgs{
					hasPayload: in.Has("payload"),
					hasQueue:   in.Has("queue"),
				}
			)

			if err = in.Decode(args); err != nil {
				return
			}

			// Converting Payload argument
			if args.hasPayload {
				aux := expr.Must(expr.Select(in, "payload"))
				switch aux.Type() {
				case h.reg.Type("String").Type():
					args.payloadString = aux.Get().(string)
				case h.reg.Type("Reader").Type():
					args.payloadStream = aux.Get().(io.Reader)
				}
			}

			return out, h.write(ctx, args)
		},
	}
}
