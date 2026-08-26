package runtime

import (
	"context"

	autoTypes "github.com/crusttech/human/server/automation/types"
	"github.com/crusttech/human/server/system/agentic/guard"
	"github.com/crusttech/human/server/system/agentic/knowledge"
	"github.com/crusttech/human/server/system/agentic/observability"
	"github.com/crusttech/human/server/system/agentic/skills"
	"github.com/crusttech/human/server/system/types"
)

type (
	// runtime executes the agent conversation loop.
	runtime struct {
		registry          Registry
		llm               LLMClient
		mcp               MCPClient
		conversationStore ConversationStore
		knowledgeBase     KnowledgeBaseStore
		namespaceLookup   knowledge.NamespaceLookup
		moduleLookup      knowledge.ModuleLookup
		obs               *observability.Bus
		taqService        TAQService
		workflowService   WorkflowService
		nsModResolver     NsModResolver
		builtinGuard      guard.GuardService
		providerGuard     guard.GuardService
		skills            skills.Registry
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

		// ApprovedTools are the tools the user has agreed this agent may use,
		// by name. A tool whose mode is "ask" runs only if it is named here.
		//
		// They are carried by the request rather than stored on the
		// conversation deliberately: approving is the user's own act, and a
		// caller that sends an approval it was never given has done nothing it
		// could not have done by clicking the prompt. Keeping them out of the
		// database also means an approval never quietly becomes policy.
		ApprovedTools []string `json:"approvedTools,omitempty"`
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
		Decisions      []DecisionInfo `json:"decisions"`
		Usage          Usage          `json:"usage"`
		Context        string         `json:"context,omitempty"`

		// Status is "complete", or "awaiting_approval" when the run stopped to
		// ask. In the second case Output holds what the agent had to say up to
		// that point and PendingApproval says what it wants to do next.
		Status          string           `json:"status,omitempty"`
		PendingApproval *PendingApproval `json:"pendingApproval,omitempty"`
	}

	// PendingApproval is the tool call a run stopped on, described well enough
	// for a person to say yes or no to it.
	PendingApproval struct {
		Tool string         `json:"tool"`
		Args map[string]any `json:"args,omitempty"`
		// Risk is the tool's own classification — "write" or "destructive" —
		// so a prompt can say how much is being asked for.
		Risk string `json:"risk,omitempty"`
	}

	DecisionInfo struct {
		Iteration int      `json:"iteration"`
		Decision  string   `json:"decision"`
		Tools     []string `json:"tools,omitempty"`
		Reasoning string   `json:"reasoning,omitempty"`
		Usage     Usage    `json:"usage"`
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
		InputTokens   int `json:"inputTokens"`
		OutputTokens  int `json:"outputTokens"`
		ContextWindow int `json:"contextWindow"`
	}

	// LLMClient abstracts the LLM provider.
	LLMClient interface {
		// Chat interacts with the LLM.
		// It takes the agent definition (for system prompt/config), current history, and available tools.
		Chat(ctx context.Context, prompt string, history []types.AiConversationMessage, tools []Tool, config LLMConfig) (*LLMResponse, error)
	}

	LLMConfig struct {
		ProviderID   uint64
		Model        string
		Temperature  *float64
		OutputTokens int
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

	KnowledgeBaseStore interface {
		FindByID(ctx context.Context, ID uint64) (*types.KnowledgeBase, error)
	}

	TAQService interface {
		FindByID(ctx context.Context, ID uint64) (*autoTypes.NgAutomation, error)
		LookupByHandle(ctx context.Context, handle string) (*autoTypes.NgAutomation, error)
	}

	WorkflowService interface {
		FindByID(ctx context.Context, ID uint64) (*autoTypes.Workflow, error)
		LookupByHandle(ctx context.Context, handle string) (*autoTypes.Workflow, error)
	}

	NsModResolver interface {
		Resolve(ctx context.Context, namespace, module string) (nsID, modID uint64, err error)
		LookupNamespace(ctx context.Context, id uint64) (NsHandle, error)
		LookupModule(ctx context.Context, nsID, modID uint64) (ModHandle, error)
	}

	NsHandle struct {
		ID             uint64
		Handle         string
		Name           string
		CreatedByAgent uint64
	}

	ModHandle struct {
		ID             uint64
		NamespaceID    uint64
		Handle         string
		Name           string
		CreatedByAgent uint64
	}

	// MCPClient abstracts the Model Context Protocol tools.
	MCPClient interface {
		// GetTools returns tools filtered to the agent's allowed tool list.
		GetTools(ctx context.Context, allowedTools []string) ([]Tool, error)
		// ExecuteTool executes a specific tool.
		ExecuteTool(ctx context.Context, toolName string, args map[string]any) (any, error)
		// ToolNamesIn lists the tools in a group at or below a risk ceiling, so
		// a grant can name a set instead of every member of it.
		ToolNamesIn(group, maxRisk string) []string
	}

	Tool struct {
		Name        string         `json:"name"`
		Title       string         `json:"title,omitempty"`
		Description string         `json:"description"`
		InputSchema map[string]any `json:"inputSchema"`
	}
)

// Runtime creates a new Agent Runtime.
func Runtime(reg Registry, llm LLMClient, mcp MCPClient, convStore ConversationStore, kb KnowledgeBaseStore, ns knowledge.NamespaceLookup, mod knowledge.ModuleLookup, obs *observability.Bus) *runtime {
	return &runtime{
		registry:          reg,
		llm:               llm,
		mcp:               mcp,
		conversationStore: convStore,
		knowledgeBase:     kb,
		namespaceLookup:   ns,
		moduleLookup:      mod,
		obs:               obs,
	}
}

func (r *runtime) SetLookups(ns knowledge.NamespaceLookup, mod knowledge.ModuleLookup) {
	r.namespaceLookup = ns
	r.moduleLookup = mod
}

func (r *runtime) SetTAQService(s TAQService) {
	r.taqService = s
}

func (r *runtime) SetWorkflowService(s WorkflowService) {
	r.workflowService = s
}

func (r *runtime) SetNsModResolver(s NsModResolver) {
	r.nsModResolver = s
}

func (r *runtime) SetBuiltinGuard(g guard.GuardService) {
	r.builtinGuard = g
}

func (r *runtime) SetProviderGuard(g guard.GuardService) {
	r.providerGuard = g
}

func (r *runtime) SetSkillRegistry(s skills.Registry) {
	r.skills = s
}

// What a run ended as. A run that stopped to ask is not an error: the
// conversation is saved as it stands, and coming back with the tool approved
// picks it up.
const (
	StatusComplete         = "complete"
	StatusAwaitingApproval = "awaiting_approval"
)
