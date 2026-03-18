package automation

import (
	"context"
	"fmt"

	atypes "github.com/cortezaproject/corteza/server/automation/types"
	agenticRuntime "github.com/cortezaproject/corteza/server/system/agentic/runtime"
	"github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/cortezaproject/corteza/server/system/types"
)

type (
	typeRegistry interface {
		Type(ref string) expr.Type
	}

	ngAgentConstructSvc interface {
		AddFunctions(ff ...atypes.ConstructFunction)
	}

	ngAgentRuntime interface {
		Run(ctx context.Context, req *agenticRuntime.AgentRequest) (*agenticRuntime.AgentResponse, error)
	}

	ngAgentConversationStore interface {
		FindByID(ctx context.Context, ID uint64) (*types.AiConversation, error)
	}

	ngAgentHandler struct {
		reg           ngAgentConstructSvc
		tReg          typeRegistry
		runtime       ngAgentRuntime
		conversations ngAgentConversationStore
	}

	ngAgentPromptArgs struct {
		hasAgentID bool
		AgentID    uint64

		hasInput bool
		Input    string
	}

	ngAgentPromptResults struct {
		Output         string
		ConversationID uint64
	}

	ngAgentContinueArgs struct {
		hasConversationID bool
		ConversationID    uint64

		hasInput bool
		Input    string
	}

	ngAgentContinueResults struct {
		Output         string
		ConversationID uint64
	}
)

func NgAgentHandler(reg ngAgentConstructSvc, tReg typeRegistry, rt ngAgentRuntime, conversations ngAgentConversationStore) *ngAgentHandler {
	h := &ngAgentHandler{reg: reg, tReg: tReg, runtime: rt, conversations: conversations}
	h.register()
	return h
}

func (h ngAgentHandler) register() {
	h.reg.AddFunctions(
		h.Prompt(),
		h.Continue(),
	)
}

func (h ngAgentHandler) Prompt() atypes.ConstructFunction {
	return atypes.ConstructFunction{
		Ref:    "agentPrompt",
		Kind:   "function",
		Groups: []string{"Agents"},
		Meta: &atypes.ConstructFunctionMeta{
			Short:       "Prompt agent",
			Description: "Send a prompt to an agent",
			Icon:        &atypes.NgAutomationIcon{Type: "name", Value: "microchip-ai"},
			Weight:      1,
		},

		Labels: map[string]string{"agent": "step,workflow"},

		Parameters: []*atypes.Param{
			{
				ArgumentName: "agentID",
				Name:         "",
				Types:        []string{"ID"}, Required: true,
			},
			{
				ArgumentName: "input",
				Name:         "",
				Types:        []string{"String"}, Required: true,
			},
		},

		Results: []*atypes.Param{
			{
				ArgumentName: "output",
				Name:         "",
				Types:        []string{"String"},
			},
			{
				ArgumentName: "conversationID",
				Name:         "",
				Types:        []string{"ID"},
			},
		},

		Segments: []atypes.ConstructSegment{{
			Sections: []atypes.ConstructSection{{
				Elements: []atypes.SectionElement{
					{
						Input: atypes.SectionElementInput{
							Type:     "AgentSelector",
							Label:    "Agent",
							Argument: "agentID",
						},
					},
					{
						Input: atypes.SectionElementInput{
							Type:     "String",
							Label:    "Input",
							Argument: "input",
						},
					},
				},
			}},
		}},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			args := &ngAgentPromptArgs{
				hasAgentID: in.Has("agentID"),
				hasInput:   in.Has("input"),
			}

			if err = in.Decode(args); err != nil {
				return
			}

			var results *ngAgentPromptResults
			if results, err = h.prompt(ctx, args); err != nil {
				return
			}

			out = &expr.Vars{}

			{
				var tval expr.TypedValue
				if tval, err = h.tReg.Type("String").Cast(results.Output); err != nil {
					return
				} else if err = expr.Assign(out, "output", tval); err != nil {
					return
				}
			}

			{
				var tval expr.TypedValue
				if tval, err = h.tReg.Type("ID").Cast(results.ConversationID); err != nil {
					return
				} else if err = expr.Assign(out, "conversationID", tval); err != nil {
					return
				}
			}

			return
		},
	}
}

func (h ngAgentHandler) Continue() atypes.ConstructFunction {
	return atypes.ConstructFunction{
		Ref:    "agentContinue",
		Kind:   "function",
		Groups: []string{"Agents"},
		Meta: &atypes.ConstructFunctionMeta{
			Short:       "Continue agent conversation",
			Description: "Continue an existing agent conversation",
			Icon:        &atypes.NgAutomationIcon{Type: "name", Value: "microchip-ai"},
			Weight:      2,
		},

		Labels: map[string]string{"agent": "step,workflow"},

		Parameters: []*atypes.Param{
			{
				ArgumentName: "conversationID",
				Name:         "",
				Types:        []string{"ID"}, Required: true,
			},
			{
				ArgumentName: "input",
				Name:         "",
				Types:        []string{"String"}, Required: true,
			},
		},

		Results: []*atypes.Param{
			{
				ArgumentName: "output",
				Name:         "",
				Types:        []string{"String"},
			},
			{
				ArgumentName: "conversationID",
				Name:         "",
				Types:        []string{"ID"},
			},
		},

		Segments: []atypes.ConstructSegment{{
			Sections: []atypes.ConstructSection{{
				Elements: []atypes.SectionElement{
					{
						Input: atypes.SectionElementInput{
							Type:     "ID",
							Label:    "Conversation ID",
							Argument: "conversationID",
						},
					},
					{
						Input: atypes.SectionElementInput{
							Type:     "String",
							Label:    "Input",
							Argument: "input",
						},
					},
				},
			}},
		}},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			args := &ngAgentContinueArgs{
				hasConversationID: in.Has("conversationID"),
				hasInput:          in.Has("input"),
			}

			if err = in.Decode(args); err != nil {
				return
			}

			var results *ngAgentContinueResults
			if results, err = h.continueConversation(ctx, args); err != nil {
				return
			}

			out = &expr.Vars{}

			{
				var tval expr.TypedValue
				if tval, err = h.tReg.Type("String").Cast(results.Output); err != nil {
					return
				} else if err = expr.Assign(out, "output", tval); err != nil {
					return
				}
			}

			{
				var tval expr.TypedValue
				if tval, err = h.tReg.Type("ID").Cast(results.ConversationID); err != nil {
					return
				} else if err = expr.Assign(out, "conversationID", tval); err != nil {
					return
				}
			}

			return
		},
	}
}

func (h ngAgentHandler) prompt(ctx context.Context, args *ngAgentPromptArgs) (*ngAgentPromptResults, error) {
	resp, err := h.runtime.Run(ctx, &agenticRuntime.AgentRequest{
		AgentID: args.AgentID,
		Input:   args.Input,
	})
	if err != nil {
		return nil, err
	}

	return &ngAgentPromptResults{
		Output:         resp.Output,
		ConversationID: resp.ConversationID,
	}, nil
}

func (h ngAgentHandler) continueConversation(ctx context.Context, args *ngAgentContinueArgs) (*ngAgentContinueResults, error) {
	conv, err := h.conversations.FindByID(ctx, args.ConversationID)
	if err != nil {
		return nil, fmt.Errorf("could not load conversation: %w", err)
	}

	resp, err := h.runtime.Run(ctx, &agenticRuntime.AgentRequest{
		AgentID:        conv.AgentID,
		ConversationID: args.ConversationID,
		Input:          args.Input,
	})
	if err != nil {
		return nil, err
	}

	return &ngAgentContinueResults{
		Output:         resp.Output,
		ConversationID: resp.ConversationID,
	}, nil
}
