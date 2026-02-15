package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	sysTypes "github.com/cortezaproject/corteza/server/system/types"
)

type (
	openaiRequest struct {
		Model       string          `json:"model"`
		Messages    []openaiMessage `json:"messages"`
		Tools       []openaiTool    `json:"tools,omitempty"`
		Temperature *float64        `json:"temperature,omitempty"`
		MaxTokens   *int            `json:"max_tokens,omitempty"`
	}

	openaiMessage struct {
		Role       string           `json:"role"`
		Content    *string          `json:"content"`
		ToolCalls  []openaiToolCall `json:"tool_calls,omitempty"`
		ToolCallID string           `json:"tool_call_id,omitempty"`
	}

	openaiToolCall struct {
		ID       string             `json:"id"`
		Type     string             `json:"type"`
		Function openaiToolCallFunc `json:"function"`
	}

	openaiToolCallFunc struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	}

	openaiTool struct {
		Type     string             `json:"type"`
		Function openaiToolFunction `json:"function"`
	}

	openaiToolFunction struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Parameters  any    `json:"parameters"`
	}

	openaiResponse struct {
		Choices []openaiChoice `json:"choices"`
		Usage   openaiUsage    `json:"usage"`
	}

	openaiChoice struct {
		Message      openaiMessage `json:"message"`
		FinishReason string        `json:"finish_reason"`
	}

	openaiUsage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	}
)

func promptOpenAI(ctx context.Context, provider *sysTypes.LlmProvider, cred *sysTypes.Credential, messages []Message, tools []Tool) (*Response, error) {
	req := openaiRequest{
		Model:    provider.Config.Model,
		Messages: toOpenAIMessages(messages),
	}

	if len(tools) > 0 {
		req.Tools = toOpenAITools(tools)
	}

	if provider.Config.Temperature > 0 {
		t := provider.Config.Temperature
		req.Temperature = &t
	}

	if provider.Config.MaxTokens > 0 {
		m := provider.Config.MaxTokens
		req.MaxTokens = &m
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := provider.Config.PromptURL + "/chat/completions"
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	apiKey := cred.Credentials
	switch provider.Provider {
	case "azure":
		httpReq.Header.Set("api-key", apiKey)
	default:
		httpReq.Header.Set("Authorization", "Bearer "+apiKey)
	}

	timeout := 30 * time.Second
	if provider.Config.Timeout != "" {
		if d, err := time.ParseDuration(provider.Config.Timeout); err == nil {
			timeout = d
		}
	}

	httpResp, err := (&http.Client{Timeout: timeout}).Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("LLM request failed: %w", err)
	}
	defer httpResp.Body.Close()

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("LLM API returned status %d: %s", httpResp.StatusCode, string(respBody))
	}

	var oaiResp openaiResponse
	if err := json.Unmarshal(respBody, &oaiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(oaiResp.Choices) == 0 {
		return nil, fmt.Errorf("LLM returned no choices")
	}

	choice := oaiResp.Choices[0]
	return &Response{
		Message:      fromOpenAIMessage(choice.Message),
		FinishReason: choice.FinishReason,
		Usage: Usage{
			PromptTokens:     oaiResp.Usage.PromptTokens,
			CompletionTokens: oaiResp.Usage.CompletionTokens,
			TotalTokens:      oaiResp.Usage.TotalTokens,
		},
	}, nil
}

func toOpenAIMessages(messages []Message) []openaiMessage {
	out := make([]openaiMessage, len(messages))
	for i, m := range messages {
		msg := openaiMessage{
			Role:       m.Role,
			ToolCallID: m.ToolCallID,
		}

		if m.Content != "" {
			c := m.Content
			msg.Content = &c
		}

		if len(m.ToolCalls) > 0 {
			msg.ToolCalls = make([]openaiToolCall, len(m.ToolCalls))
			for j, tc := range m.ToolCalls {
				argsJSON, _ := json.Marshal(tc.Arguments)
				msg.ToolCalls[j] = openaiToolCall{
					ID:   tc.ID,
					Type: "function",
					Function: openaiToolCallFunc{
						Name:      tc.Name,
						Arguments: string(argsJSON),
					},
				}
			}
		}

		out[i] = msg
	}
	return out
}

func toOpenAITools(tools []Tool) []openaiTool {
	out := make([]openaiTool, len(tools))
	for i, t := range tools {
		out[i] = openaiTool{
			Type: "function",
			Function: openaiToolFunction{
				Name:        t.Function.Name,
				Description: t.Function.Description,
				Parameters:  t.Function.Parameters,
			},
		}
	}
	return out
}

func fromOpenAIMessage(msg openaiMessage) Message {
	m := Message{
		Role: msg.Role,
	}

	if msg.Content != nil {
		m.Content = *msg.Content
	}

	if len(msg.ToolCalls) > 0 {
		m.ToolCalls = make([]ToolCall, len(msg.ToolCalls))
		for i, tc := range msg.ToolCalls {
			var args map[string]any
			_ = json.Unmarshal([]byte(tc.Function.Arguments), &args)

			m.ToolCalls[i] = ToolCall{
				ID:        tc.ID,
				Name:      tc.Function.Name,
				Arguments: args,
			}
		}
	}

	return m
}

