package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/davecgh/go-spew/spew"

	"github.com/cortezaproject/corteza/server/pkg/agentic/observability"
	"github.com/cortezaproject/corteza/server/pkg/agentic/policy"
	"github.com/cortezaproject/corteza/server/pkg/auth"
	"github.com/cortezaproject/corteza/server/pkg/id"
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

	// Observability setup
	traceID := sid()
	agentIDStr := strconv.FormatUint(agent.ID, 10)
	convIDStr := strconv.FormatUint(conversation.ID, 10)
	userIDStr := strconv.FormatUint(auth.GetIdentityFromContext(ctx).Identity(), 10)
	rootSpanID := sid()
	rootStartedAt := time.Now()

	r.emitEvent(observability.AgentEvent{
		ID:             sid(),
		TraceID:        traceID,
		SpanID:         rootSpanID,
		Timestamp:      time.Now(),
		Event:          "agent.invoked",
		AgentID:        agentIDStr,
		UserID:         userIDStr,
		ConversationID: convIDStr,
		Details:        map[string]any{"input": req.Input},
	})

	// 4. prompt.build span — system prompt preparation
	promptBuildStart := time.Now()
	systemPrompt := agent.Behavior.SystemPrompt
	r.emitSpan(observability.AgentSpan{
		ID:             sid(),
		ParentID:       rootSpanID,
		TraceID:        traceID,
		Name:           "prompt.build",
		AgentID:        agentIDStr,
		UserID:         userIDStr,
		ConversationID: convIDStr,
		StartedAt:      promptBuildStart,
		EndedAt:        time.Now(),
		Status:         observability.StatusOK,
		Attributes:     map[string]any{"promptLength": len(systemPrompt)},
	})

	// 5. Execution Loop
	limits := agent.Execution.Limits
	maxIterations := limits.MaxIterations
	if maxIterations == 0 {
		maxIterations = defaultMaxIterations
	}

	var finalResponse string
	var usage Usage
	var executedTools []ToolCallInfo
	var decisions []DecisionInfo
	var runErr error

	for i := 0; i < maxIterations; i++ {
		config := LLMConfig{
			ProviderID:  agent.Execution.Model.LLMProviderID,
			Model:       agent.Execution.Model.Model,
			Temperature: agent.Execution.Model.Temperature,
			MaxTokens:   agent.Execution.Model.MaxTokens,
		}

		llmSpanID := sid()
		llmStart := time.Now()

		llmResp, llmErr := r.llm.Chat(ctx, systemPrompt, conversation.Messages, tools, config)

		llmSpan := observability.AgentSpan{
			ID:             llmSpanID,
			ParentID:       rootSpanID,
			TraceID:        traceID,
			Name:           "llm.chat",
			AgentID:        agentIDStr,
			UserID:         userIDStr,
			ConversationID: convIDStr,
			StartedAt:      llmStart,
			EndedAt:        time.Now(),
		}
		if llmErr != nil {
			llmSpan.Status = observability.StatusError
			llmSpan.Error = llmErr
			r.emitSpan(llmSpan)
			runErr = fmt.Errorf("llm chat failed: %w", llmErr)
			break
		}
		llmSpan.Status = observability.StatusOK
		llmSpan.Attributes = map[string]any{
			"inputTokens":  llmResp.Usage.InputTokens,
			"outputTokens": llmResp.Usage.OutputTokens,
		}
		r.emitSpan(llmSpan)

		usage.accumulate(llmResp.Usage)

		if limits.MaxTokens > 0 && usage.TotalTokens > limits.MaxTokens {
			runErr = fmt.Errorf("limit_exceeded: token limit reached")
			break
		}

		// Process response
		if len(llmResp.ToolCalls) > 0 {
			toolNames := make([]string, len(llmResp.ToolCalls))
			for j, tc := range llmResp.ToolCalls {
				toolNames[j] = tc.Name
			}
			d := DecisionInfo{
				Iteration: i + 1,
				Decision:  "tool_call",
				Tools:     toolNames,
				Reasoning: llmResp.Text,
			}
			decisions = append(decisions, d)
			r.emitEvent(observability.AgentEvent{
				ID:             sid(),
				TraceID:        traceID,
				SpanID:         rootSpanID,
				Timestamp:      time.Now(),
				Event:          "agent.decision",
				AgentID:        agentIDStr,
				UserID:         userIDStr,
				ConversationID: convIDStr,
				Details:        map[string]any{"iteration": d.Iteration, "decision": d.Decision, "tools": d.Tools},
			})

			// Add assistant message with tool calls to history
			conversation.Messages = append(conversation.Messages, types.AiConversationMessage{
				Role:      "assistant",
				Content:   llmResp.Text,
				ToolCalls: toAiToolCalls(llmResp.ToolCalls),
			})

			// Execute tools and append results to conversation
			results, infos := r.executeTools(ctx, agent, llmResp.ToolCalls, traceID, rootSpanID, agentIDStr, userIDStr, convIDStr)
			conversation.Messages = append(conversation.Messages, results...)
			executedTools = append(executedTools, infos...)
		} else {
			d := DecisionInfo{
				Iteration: i + 1,
				Decision:  "respond",
			}
			decisions = append(decisions, d)
			r.emitEvent(observability.AgentEvent{
				ID:             sid(),
				TraceID:        traceID,
				SpanID:         rootSpanID,
				Timestamp:      time.Now(),
				Event:          "agent.decision",
				AgentID:        agentIDStr,
				UserID:         userIDStr,
				ConversationID: convIDStr,
				Details:        map[string]any{"iteration": d.Iteration, "decision": d.Decision},
			})

			// agent.respond span — final response assembly
			respondStart := time.Now()
			finalResponse = llmResp.Text
			conversation.Messages = append(conversation.Messages, types.AiConversationMessage{
				Role:    "assistant",
				Content: finalResponse,
			})
			r.emitSpan(observability.AgentSpan{
				ID:             sid(),
				ParentID:       rootSpanID,
				TraceID:        traceID,
				Name:           "agent.respond",
				AgentID:        agentIDStr,
				UserID:         userIDStr,
				ConversationID: convIDStr,
				StartedAt:      respondStart,
				EndedAt:        time.Now(),
				Status:         observability.StatusOK,
				Attributes:     map[string]any{"responseLength": len(finalResponse)},
			})
			break
		}
	}

	// End root span
	rootStatus := observability.StatusOK
	if runErr != nil {
		rootStatus = observability.StatusError
	}
	r.emitSpan(observability.AgentSpan{
		ID:             rootSpanID,
		TraceID:        traceID,
		Name:           "agent.run",
		AgentID:        agentIDStr,
		UserID:         userIDStr,
		ConversationID: convIDStr,
		StartedAt:      rootStartedAt,
		EndedAt:        time.Now(),
		Status:         rootStatus,
		Error:          runErr,
	})

	r.emitEvent(observability.AgentEvent{
		ID:             sid(),
		TraceID:        traceID,
		SpanID:         rootSpanID,
		Timestamp:      time.Now(),
		Event:          "agent.completed",
		AgentID:        agentIDStr,
		UserID:         userIDStr,
		ConversationID: convIDStr,
		Details: map[string]any{
			"totalTokens": usage.TotalTokens,
			"toolCalls":   len(executedTools),
		},
	})

	if runErr != nil {
		return nil, runErr
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
		Decisions:      decisions,
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
func (r *runtime) executeTools(ctx context.Context, agent *types.Agent, calls []ToolCall, traceID, parentSpanID, agentID, userID, convID string) ([]types.AiConversationMessage, []ToolCallInfo) {
	var (
		messages []types.AiConversationMessage
		infos    []ToolCallInfo
	)

	for _, tc := range calls {
		start := time.Now()

		policyStart := time.Now()
		decision := policy.Evaluate(agent, tc.Name, tc.Args)
		policySpan := observability.AgentSpan{
			ID:             sid(),
			ParentID:       parentSpanID,
			TraceID:        traceID,
			Name:           "policy.evaluate",
			AgentID:        agentID,
			UserID:         userID,
			ConversationID: convID,
			StartedAt:      policyStart,
			EndedAt:        time.Now(),
			Attributes:     map[string]any{"tool": tc.Name, "allowed": decision.Allowed, "reason": decision.Reason},
		}
		if decision.Allowed {
			policySpan.Status = observability.StatusOK
		} else {
			policySpan.Status = observability.StatusError
		}
		r.emitSpan(policySpan)

		if !decision.Allowed {
			r.emitEvent(observability.AgentEvent{
				ID:             sid(),
				TraceID:        traceID,
				SpanID:         parentSpanID,
				Timestamp:      time.Now(),
				Event:          "tool.denied",
				AgentID:        agentID,
				UserID:         userID,
				ConversationID: convID,
				Details:        map[string]any{"tool": tc.Name, "reason": decision.Reason},
			})
			infos = append(infos, ToolCallInfo{
				Tool:  tc.Name,
				Args:  tc.Args,
				Error: decision.Reason,
			})
			messages = append(messages, types.AiConversationMessage{
				Role: "tool",
				ToolResults: []types.AiConversationToolResult{
					{CallID: tc.ID, Data: decision.Reason, Error: decision.Reason},
				},
			})
			continue
		}

		r.emitEvent(observability.AgentEvent{
			ID:             sid(),
			TraceID:        traceID,
			SpanID:         parentSpanID,
			Timestamp:      time.Now(),
			Event:          "tool.called",
			AgentID:        agentID,
			UserID:         userID,
			ConversationID: convID,
			Details:        map[string]any{"tool": tc.Name, "args": decision.SanitizedArgs},
		})

		result, execErr := r.mcp.ExecuteTool(ctx, tc.Name, decision.SanitizedArgs)
		duration := int(time.Since(start).Milliseconds())

		toolSpan := observability.AgentSpan{
			ID:             sid(),
			ParentID:       parentSpanID,
			TraceID:        traceID,
			Name:           "tool.execute",
			AgentID:        agentID,
			UserID:         userID,
			ConversationID: convID,
			StartedAt:      start,
			EndedAt:        time.Now(),
			Attributes:     map[string]any{"tool": tc.Name},
		}
		if execErr != nil {
			toolSpan.Status = observability.StatusError
			toolSpan.Error = execErr
		} else {
			toolSpan.Status = observability.StatusOK
		}
		r.emitSpan(toolSpan)

		r.emitEvent(observability.AgentEvent{
			ID:             sid(),
			TraceID:        traceID,
			SpanID:         parentSpanID,
			Timestamp:      time.Now(),
			Event:          "tool.completed",
			AgentID:        agentID,
			UserID:         userID,
			ConversationID: convID,
			Details:        map[string]any{"tool": tc.Name, "durationMs": duration},
		})

		if resultMap, ok := result.(map[string]any); ok {
			result = policy.FilterResponse(ctx, agent, tc.Name, resultMap)
		}

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

func (r *runtime) emitSpan(span observability.AgentSpan) {
	if r.obs != nil {
		r.obs.EmitSpan(span)
	}
}

func (r *runtime) emitEvent(event observability.AgentEvent) {
	if r.obs != nil {
		r.obs.EmitEvent(event)
	}
}

func sid() string {
	return strconv.FormatUint(id.Next(), 10)
}
