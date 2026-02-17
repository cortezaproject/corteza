package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/davecgh/go-spew/spew"

	"github.com/cortezaproject/corteza/server/system/types"
)

const (
	defaultMaxIterations = 10
)

func (r *runtime) Run(ctx context.Context, req *AgentRequest) (*AgentResponse, error) {
	// 1. Load and validate agent
	agent, err := r.registry.Get(ctx, req.AgentID)
	if err != nil {
		return nil, fmt.Errorf("failed to load agent: %w", err)
	}

	if err := validateAgent(agent); err != nil {
		return nil, err
	}

	// 2. Get available tools
	tools, err := r.mcp.GetTools(ctx, agent.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tools: %w", err)
	}

	// 3. Load/Create Conversation
	conversation, err := r.resolveConversation(ctx, req.ConversationID, agent.ID)
	if err != nil {
		return nil, err
	}

	// Add user message to history if input exists
	if req.Input != "" {
		conversation.Messages = append(conversation.Messages, types.AiConversationMessage{
			Role:    "user",
			Content: req.Input,
		})
	}

	// 4. Execution Loop
	limits := agent.Execution.Limits
	maxIterations := limits.MaxIterations
	if maxIterations == 0 {
		maxIterations = defaultMaxIterations
	}

	var finalResponse string
	var usage Usage
	var executedTools []ToolCallInfo

	for i := 0; i < maxIterations; i++ {
		config := LLMConfig{
			Model:       agent.Execution.Model.Model,
			Temperature: agent.Execution.Model.Temperature,
			MaxTokens:   agent.Execution.Model.MaxTokens,
		}

		llmResp, err := r.llm.Chat(ctx, agent.Behavior.SystemPrompt, conversation.Messages, tools, config)
		if err != nil {
			return nil, fmt.Errorf("llm chat failed: %w", err)
		}

		usage.accumulate(llmResp.Usage)

		if limits.MaxTokens > 0 && usage.TotalTokens > limits.MaxTokens {
			return nil, fmt.Errorf("limit_exceeded: token limit reached")
		}

		// Process response
		if len(llmResp.ToolCalls) > 0 {
			// Add assistant message with tool calls to history
			conversation.Messages = append(conversation.Messages, types.AiConversationMessage{
				Role:      "assistant",
				Content:   llmResp.Text,
				ToolCalls: toAiToolCalls(llmResp.ToolCalls),
			})

			// Execute tools and append results to conversation
			results, infos := r.executeTools(ctx, llmResp.ToolCalls)
			conversation.Messages = append(conversation.Messages, results...)
			executedTools = append(executedTools, infos...)
		} else {
			// Text response — done
			finalResponse = llmResp.Text
			conversation.Messages = append(conversation.Messages, types.AiConversationMessage{
				Role:    "assistant",
				Content: finalResponse,
			})
			break
		}
	}

	// Save conversation
	conversation.TokenCount = usage.TotalTokens
	if _, err := r.conversationStore.Update(ctx, conversation); err != nil {
		spew.Dump(err)
	}

	return &AgentResponse{
		Output:         finalResponse,
		ConversationID: conversation.ID,
		ToolCalls:      executedTools,
		Usage:          usage,
	}, nil
}

func validateAgent(agent *types.Agent) error {
	if agent.Status != "active" {
		return fmt.Errorf("agent %d is not active (status: %s)", agent.ID, agent.Status)
	}
	return nil
}

func (r *runtime) resolveConversation(ctx context.Context, conversationID, agentID uint64) (*types.AiConversation, error) {
	if conversationID != 0 {
		conv, err := r.conversationStore.FindByID(ctx, conversationID)
		if err != nil {
			return nil, fmt.Errorf("failed to get conversation: %w", err)
		}
		if conv != nil {
			return conv, nil
		}
	}

	conv, err := r.conversationStore.Create(ctx, &types.AiConversation{
		AgentID:  agentID,
		Messages: types.AiConversationMessages{},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create conversation: %w", err)
	}
	return conv, nil
}

func (u *Usage) accumulate(other Usage) {
	u.InputTokens += other.InputTokens
	u.OutputTokens += other.OutputTokens
	u.TotalTokens += other.TotalTokens
}

// executeTools runs each tool call and returns conversation messages + telemetry info.
func (r *runtime) executeTools(ctx context.Context, calls []ToolCall) ([]types.AiConversationMessage, []ToolCallInfo) {
	var (
		messages []types.AiConversationMessage
		infos    []ToolCallInfo
	)

	for _, tc := range calls {
		start := time.Now()

		// TODO: Policy check (allow/deny)

		result, execErr := r.mcp.ExecuteTool(ctx, tc.Name, tc.Args)
		duration := int(time.Since(start).Milliseconds())

		resultData, _ := json.Marshal(result)

		toolResult := types.AiConversationToolResult{
			CallID: tc.ID,
			Data:   string(resultData),
		}

		info := ToolCallInfo{
			Tool:       tc.Name,
			Args:       tc.Args,
			Result:     result,
			DurationMs: duration,
		}

		if execErr != nil {
			toolResult.Error = execErr.Error()
			info.Error = execErr.Error()
		}

		infos = append(infos, info)
		messages = append(messages, types.AiConversationMessage{
			Role: "tool",
			ToolResults: []types.AiConversationToolResult{
				toolResult,
			},
		})
	}

	return messages, infos
}

// toAiToolCalls converts runtime ToolCalls (with parsed Args) to the persisted format.
func toAiToolCalls(tcs []ToolCall) []types.AiConversationToolCall {
	out := make([]types.AiConversationToolCall, len(tcs))
	for i, tc := range tcs {
		data, _ := json.Marshal(tc.Args)
		out[i] = types.AiConversationToolCall{
			CallID: tc.ID,
			Name:   tc.Name,
			Data:   string(data),
		}
	}
	return out
}
