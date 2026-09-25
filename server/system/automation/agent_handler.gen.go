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
	atypes "github.com/crusttech/human/server/automation/types"
	"github.com/crusttech/human/server/pkg/expr"
	"github.com/crusttech/human/server/pkg/wfexec"
	"github.com/crusttech/human/server/system/types"
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
		h.Lookup(),
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
				Meta: &atypes.ParamMeta{
					Label:  "Agent",
					Visual: map[string]interface{}{"input": map[string]interface{}{"type": "agent"}},
				},
			},
			{
				Name:  "input",
				Types: []string{"String"}, Required: true,
				Meta: &atypes.ParamMeta{
					Label: "Input",
				},
			},
			{
				Name:  "conversationID",
				Types: []string{"ID"},
				Meta: &atypes.ParamMeta{
					Label: "Conversation ID",
				},
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

type (
	agentLookupArgs struct {
		hasLookup    bool
		Lookup       interface{}
		lookupID     uint64
		lookupHandle string
		lookupRes    *types.Agent
	}

	agentLookupResults struct {
		Agent *types.Agent
	}
)

func (a agentLookupArgs) GetLookup() (bool, uint64, string, *types.Agent) {
	return a.hasLookup, a.lookupID, a.lookupHandle, a.lookupRes
}

// Lookup function Agent lookup
//
// expects implementation of lookup function:
//
//	func (h agentHandler) lookup(ctx context.Context, args *agentLookupArgs) (results *agentLookupResults, err error) {
//	   return
//	}
func (h agentHandler) Lookup() *atypes.Function {
	return &atypes.Function{
		Ref:    "agentLookup",
		Kind:   "function",
		Labels: map[string]string{"agent": "step,workflow"},
		Meta: &atypes.FunctionMeta{
			Short:       "Agent lookup",
			Description: "Find a specific agent by ID or handle; a handle shared by several agents resolves to the first",
		},

		Parameters: []*atypes.Param{
			{
				Name:  "lookup",
				Types: []string{"ID", "Handle", "Agent"}, Required: true,
				Meta: &atypes.ParamMeta{
					Label: "Lookup",
				},
			},
		},

		Results: []*atypes.Param{

			{
				Name:  "agent",
				Types: []string{"Agent"},
			},
		},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			var (
				args = &agentLookupArgs{
					hasLookup: in.Has("lookup"),
				}
			)

			if err = in.Decode(args); err != nil {
				return
			}

			// Converting Lookup argument
			if args.hasLookup {
				aux := expr.Must(expr.Select(in, "lookup"))
				switch aux.Type() {
				case h.reg.Type("ID").Type():
					args.lookupID = aux.Get().(uint64)
				case h.reg.Type("Handle").Type():
					args.lookupHandle = aux.Get().(string)
				case h.reg.Type("Agent").Type():
					args.lookupRes = aux.Get().(*types.Agent)
				}
			}

			var results *agentLookupResults
			if results, err = h.lookup(ctx, args); err != nil {
				return
			}

			out = &expr.Vars{}

			{
				// converting results.Agent (*types.Agent) to Agent
				var (
					tval expr.TypedValue
				)

				if tval, err = h.reg.Type("Agent").Cast(results.Agent); err != nil {
					return
				} else if err = expr.Assign(out, "agent", tval); err != nil {
					return
				}
			}

			return
		},
	}
}
