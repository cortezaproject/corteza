package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	sysTypes "github.com/cortezaproject/corteza/server/system/types"
)

type (
	anthropicRequest struct {
		Model       string              `json:"model"`
		System      string              `json:"system,omitempty"`
		Messages    []anthropicMessage  `json:"messages"`
		Tools       []anthropicTool     `json:"tools,omitempty"`
		MaxTokens   int                 `json:"max_tokens"`
		Temperature *float64            `json:"temperature,omitempty"`
	}

	anthropicMessage struct {
		Role    string `json:"role"`
		Content any    `json:"content"`
	}

	anthropicContentBlock struct {
		Type      string `json:"type"`
		Text      string `json:"text,omitempty"`
		ID        string `json:"id,omitempty"`
		Name      string `json:"name,omitempty"`
		Input     any    `json:"input,omitempty"`
		ToolUseID string `json:"tool_use_id,omitempty"`
		Content   string `json:"content,omitempty"`
	}

	anthropicTool struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		InputSchema any    `json:"input_schema"`
	}

	anthropicResponse struct {
		Content    []anthropicContentBlock `json:"content"`
		StopReason string                  `json:"stop_reason"`
		Usage      anthropicUsage          `json:"usage"`
	}

	anthropicUsage struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	}
)

const anthropicAPIVersion = "2023-06-01"

func promptAnthropic(ctx context.Context, provider *sysTypes.LlmProvider, cred *sysTypes.Credential, messages []Message, tools []Tool) (*Response, error) {
	req := anthropicRequest{
		Model:     provider.Config.Model,
		MaxTokens: provider.Config.MaxTokens,
	}

	if req.MaxTokens == 0 {
		req.MaxTokens = 4096
	}

	if provider.Config.Temperature > 0 {
		t := provider.Config.Temperature
		req.Temperature = &t
	}

	var (
		filtered    []Message
		systemParts []string
	)
	for _, m := range messages {
		if m.Role == "system" {
			systemParts = append(systemParts, m.Content)
		} else {
			filtered = append(filtered, m)
		}
	}
	if len(systemParts) > 0 {
		req.System = strings.Join(systemParts, "\n\n")
	}

	req.Messages = toAnthropicMessages(filtered)

	if len(tools) > 0 {
		req.Tools = toAnthropicTools(tools)
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := provider.Config.PromptURL + "/messages"
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", cred.Credentials)
	httpReq.Header.Set("anthropic-version", anthropicAPIVersion)

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

	var aResp anthropicResponse
	if err := json.Unmarshal(respBody, &aResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &Response{
		Message:      fromAnthropicResponse(aResp),
		FinishReason: normalizeAnthropicStopReason(aResp.StopReason),
		Usage: Usage{
			PromptTokens:     aResp.Usage.InputTokens,
			CompletionTokens: aResp.Usage.OutputTokens,
			TotalTokens:      aResp.Usage.InputTokens + aResp.Usage.OutputTokens,
		},
	}, nil
}

func toAnthropicMessages(messages []Message) []anthropicMessage {
	out := make([]anthropicMessage, 0, len(messages))
	for _, m := range messages {
		switch m.Role {
		case "assistant":
			if len(m.ToolCalls) > 0 {
				blocks := make([]anthropicContentBlock, 0)
				if m.Content != "" {
					blocks = append(blocks, anthropicContentBlock{Type: "text", Text: m.Content})
				}
				for _, tc := range m.ToolCalls {
					blocks = append(blocks, anthropicContentBlock{
						Type:  "tool_use",
						ID:    tc.ID,
						Name:  tc.Name,
						Input: tc.Arguments,
					})
				}
				out = append(out, anthropicMessage{Role: "assistant", Content: blocks})
			} else {
				out = append(out, anthropicMessage{Role: "assistant", Content: m.Content})
			}
		case "tool":
			out = append(out, anthropicMessage{
				Role: "user",
				Content: []anthropicContentBlock{{
					Type:      "tool_result",
					ToolUseID: m.ToolCallID,
					Content:   m.Content,
				}},
			})
		default:
			out = append(out, anthropicMessage{Role: m.Role, Content: m.Content})
		}
	}
	return out
}

func toAnthropicTools(tools []Tool) []anthropicTool {
	out := make([]anthropicTool, len(tools))
	for i, t := range tools {
		out[i] = anthropicTool{
			Name:        t.Function.Name,
			Description: t.Function.Description,
			InputSchema: t.Function.Parameters,
		}
	}
	return out
}

func fromAnthropicResponse(resp anthropicResponse) Message {
	msg := Message{Role: "assistant"}

	for _, block := range resp.Content {
		switch block.Type {
		case "text":
			msg.Content += block.Text
		case "tool_use":
			args, _ := toMapStringAny(block.Input)
			msg.ToolCalls = append(msg.ToolCalls, ToolCall{
				ID:        block.ID,
				Name:      block.Name,
				Arguments: args,
			})
		}
	}

	return msg
}

func normalizeAnthropicStopReason(reason string) string {
	switch reason {
	case "end_turn":
		return "stop"
	case "tool_use":
		return "tool_calls"
	case "max_tokens":
		return "length"
	default:
		return reason
	}
}

func toMapStringAny(v any) (map[string]any, error) {
	if m, ok := v.(map[string]any); ok {
		return m, nil
	}
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var m map[string]any
	return m, json.Unmarshal(b, &m)
}
