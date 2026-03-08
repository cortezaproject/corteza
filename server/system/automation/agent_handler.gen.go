package automation

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// system/automation/agent_handler.yaml

import (
	"context"
	atypes "github.com/cortezaproject/corteza/server/automation/types"
	"github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/cortezaproject/corteza/server/pkg/wfexec"
)

var _ wfexec.ExecResponse

type (
	agentHandlerRegistry interface {
		AddFunctions(ff ...*atypes.Function)
		Type(ref string) expr.Type
	}
)

func (h agentHandler) register() {
	h.reg.AddFunctions(
		h.Run(),
	)
}

type (
	agentRunArgs struct {
		hasAgentID bool
		AgentID    uint64

		hasInput bool
		Input    string

		hasConversationID bool
		ConversationID    uint64
	}

	agentRunResults struct {
		Output         string
		ConversationID uint64
	}
)

// Run function Invoke agent
//
// expects implementation of run function:
//
//	func (h agentHandler) run(ctx context.Context, args *agentRunArgs) (results *agentRunResults, err error) {
//	   return
//	}
func (h agentHandler) Run() *atypes.Function {
	return &atypes.Function{
		Ref:    "agentRun",
		Kind:   "function",
		Labels: map[string]string{"agent": "step,workflow"},
		Meta: &atypes.FunctionMeta{
			Short:       "Invoke agent",
			Description: "Executes an agent with the given input and returns its response",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "agentID",
				Types: []string{"ID"}, Required: true,
			},
			{
				Name:  "input",
				Types: []string{"String"}, Required: true,
			},
			{
				Name:  "conversationID",
				Types: []string{"ID"},
			},
		},

		Results: []*atypes.Param{

			{
				Name:  "output",
				Types: []string{"String"},
			},

			{
				Name:  "conversationID",
				Types: []string{"ID"},
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &agentRunArgs{
					hasAgentID:        in.Has("agentID"),
					hasInput:          in.Has("input"),
					hasConversationID: in.Has("conversationID"),
				}
			)

			if err = in.Decode(args); err != nil {
				return
			}

			var results *agentRunResults
			if results, err = h.run(ctx, args); err != nil {
				return
			}

			out = &expr.Vars{}

			{
				// converting results.Output (string) to String
				var (
					tval expr.TypedValue
				)

				if tval, err = h.reg.Type("String").Cast(results.Output); err != nil {
					return
				} else if err = expr.Assign(out, "output", tval); err != nil {
					return
				}
			}

			{
				// converting results.ConversationID (uint64) to ID
				var (
					tval expr.TypedValue
				)

				if tval, err = h.reg.Type("ID").Cast(results.ConversationID); err != nil {
					return
				} else if err = expr.Assign(out, "conversationID", tval); err != nil {
					return
				}
			}

			return
		},
	}
}
