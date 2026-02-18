package runtime

import (
	"context"

	"github.com/cortezaproject/corteza/server/system/types"
)

type (
	// runtime executes the agent conversation loop.
	runtime struct {
		registry          Registry
		llm               LLMClient
		mcp               MCPClient
		conversationStore ConversationStore
	}

	// Registry interface for fetching agent definitions
	Registry interface {
		Get(ctx context.Context, id uint64) (*types.Agent, error)
	}

	// AgentRequest represents the input for an agent execution.
	AgentRequest struct {
		AgentID        uint64         `json:"agentID,string"`
		Input          string         `json:"input"`
		ConversationID uint64         `json:"conversationID,string"`
		ExecContext    map[string]any `json:"execContext"` // User/system context
		Attachments    []Attachment   `json:"attachments"`
	}

	// Attachment represents a file or image attached to the request.
	Attachment struct {
		Name      string `json:"name"`
		MediaType string `json:"mediaType"`
		Content   []byte `json:"content"` // Or a stream/reader if needed
	}

	// AgentResponse represents the output of an agent execution.
	AgentResponse struct {
		Output         string         `json:"output"`
		ConversationID uint64         `json:"conversationID,string"`
		ToolCalls      []ToolCallInfo `json:"toolCalls"`
		Usage          Usage          `json:"usage"`
	}

	// ToolCallInfo describes a tool call that was executed.
	ToolCallInfo struct {
		Tool       string         `json:"tool"`
		Args       map[string]any `json:"args"`
		Result     any            `json:"result"`
		Error      string         `json:"error,omitempty"`
		DurationMs int            `json:"durationMs"`
	}

	// Usage tracks token usage.
	Usage struct {
		InputTokens  int `json:"inputTokens"`
		OutputTokens int `json:"outputTokens"`
		TotalTokens  int `json:"totalTokens"`
	}

	// LLMClient abstracts the LLM provider.
	LLMClient interface {
		// Chat interacts with the LLM.
		// It takes the agent definition (for system prompt/config), current history, and available tools.
		Chat(ctx context.Context, prompt string, history []types.AiConversationMessage, tools []Tool, config LLMConfig) (*LLMResponse, error)
	}

	LLMConfig struct {
		ProviderID  uint64
		Model       string
		Temperature float64
		MaxTokens   int
	}

	LLMResponse struct {
		Text      string
		ToolCalls []ToolCall
		Usage     Usage
	}

	ToolCall struct {
		ID   string
		Name string
		Args map[string]any
	}

	// MCPClient abstracts the Model Context Protocol tools.
	MCPClient interface {
		// GetTools returns available tools for the given agent context.
		GetTools(ctx context.Context, agentID uint64) ([]Tool, error)
		// ExecuteTool executes a specific tool.
		ExecuteTool(ctx context.Context, toolName string, args map[string]any) (any, error)
	}

	Tool struct {
		Name        string
		Description string
		InputSchema map[string]any // JSON schema
	}
)

// Runtime creates a new Agent Runtime.
func Runtime(reg Registry, llm LLMClient, mcp MCPClient, store ConversationStore) *runtime {
	return &runtime{
		registry:          reg,
		llm:               llm,
		mcp:               mcp,
		conversationStore: store,
	}
}
