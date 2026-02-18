package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	rt "github.com/cortezaproject/corteza/server/pkg/agentic/runtime"
	"github.com/cortezaproject/corteza/server/pkg/id"
	"github.com/cortezaproject/corteza/server/store"
	sysTypes "github.com/cortezaproject/corteza/server/system/types"
)

type Service struct {
	store store.Storer
}

func New(s store.Storer) (*Service, error) {
	return &Service{store: s}, nil
}

func (svc *Service) Create(ctx context.Context, p *sysTypes.LlmProvider, apiKey string) (*sysTypes.LlmProvider, error) {
	p.ID = id.Next()
	p.CreatedAt = time.Now().Round(time.Second)

	if p.Status == "" {
		p.Status = "active"
	}

	if apiKey != "" {
		cred := &sysTypes.Credential{
			ID:          id.Next(),
			Kind:        "api-key",
			Label:       p.Handle + " API Key",
			Credentials: apiKey,
			CreatedAt:   p.CreatedAt,
		}

		if err := store.CreateCredential(ctx, svc.store, cred); err != nil {
			return nil, fmt.Errorf("could not create credential: %w", err)
		}

		p.CredentialID = cred.ID
	}

	if err := store.CreateLlmProvider(ctx, svc.store, p); err != nil {
		return nil, fmt.Errorf("could not create LLM provider: %w", err)
	}

	return p, nil
}

func (svc *Service) LookupByID(ctx context.Context, providerID uint64) (*sysTypes.LlmProvider, error) {
	if providerID == 0 {
		return nil, fmt.Errorf("invalid LLM provider ID")
	}

	return store.LookupLlmProviderByID(ctx, svc.store, providerID)
}

func (svc *Service) LookupByHandle(ctx context.Context, handle string) (*sysTypes.LlmProvider, error) {
	if handle == "" {
		return nil, fmt.Errorf("invalid LLM provider handle")
	}

	return store.LookupLlmProviderByHandle(ctx, svc.store, handle)
}

func (svc *Service) Update(ctx context.Context, upd *sysTypes.LlmProvider) (*sysTypes.LlmProvider, error) {
	existing, err := store.LookupLlmProviderByID(ctx, svc.store, upd.ID)
	if err != nil {
		return nil, err
	}

	existing.Handle = upd.Handle
	existing.Status = upd.Status
	existing.Provider = upd.Provider
	existing.CredentialID = upd.CredentialID
	existing.Meta = upd.Meta
	existing.Config = upd.Config

	now := time.Now().Round(time.Second)
	existing.UpdatedAt = &now
	existing.UpdatedBy = upd.UpdatedBy

	if err := store.UpdateLlmProvider(ctx, svc.store, existing); err != nil {
		return nil, fmt.Errorf("could not update LLM provider: %w", err)
	}

	return existing, nil
}

func (svc *Service) Delete(ctx context.Context, providerID uint64, deletedBy uint64) error {
	if providerID == 0 {
		return fmt.Errorf("invalid LLM provider ID")
	}

	return store.DeleteLlmProviderByID(ctx, svc.store, providerID)
}

func (svc *Service) Search(ctx context.Context, f sysTypes.LlmProviderFilter) (sysTypes.LlmProviderSet, error) {
	set, _, err := store.SearchLlmProviders(ctx, svc.store, f)
	return set, err
}

// Prompt resolves the provider and its credential, then forwards the conversation to the LLM.
func (svc *Service) Prompt(ctx context.Context, providerID uint64, messages []Message, tools []Tool) (*Response, error) {
	provider, err := store.LookupLlmProviderByID(ctx, svc.store, providerID)
	if err != nil {
		return nil, fmt.Errorf("could not resolve LLM provider: %w", err)
	}

	if provider.Status != "active" {
		return nil, fmt.Errorf("LLM provider %q is not active", provider.Handle)
	}

	cred, err := store.LookupCredentialByID(ctx, svc.store, provider.CredentialID)
	if err != nil {
		return nil, fmt.Errorf("could not resolve credential for LLM provider: %w", err)
	}

	return svc.callProvider(ctx, provider, cred, messages, tools)
}

func (svc *Service) callProvider(ctx context.Context, provider *sysTypes.LlmProvider, cred *sysTypes.Credential, messages []Message, tools []Tool) (*Response, error) {
	switch provider.Provider {
	case "openai", "azure", "mistral":
		return promptOpenAI(ctx, provider, cred, messages, tools)
	case "anthropic":
		return promptAnthropic(ctx, provider, cred, messages, tools)
	default:
		return nil, fmt.Errorf("unsupported LLM provider type: %s", provider.Provider)
	}
}

func (svc *Service) Chat(ctx context.Context, prompt string, history []sysTypes.AiConversationMessage, tools []rt.Tool, config rt.LLMConfig) (*rt.LLMResponse, error) {
	messages := []Message{{Role: "system", Content: prompt}}
	for _, m := range history {
		messages = append(messages, fromConversationMessage(m))
	}

	llmTools := make([]Tool, len(tools))
	for i, t := range tools {
		llmTools[i] = Tool{
			Type: "function",
			Function: ToolFunction{
				Name:        t.Name,
				Description: t.Description,
				Parameters:  t.InputSchema,
			},
		}
	}

	resp, err := svc.Prompt(ctx, config.ProviderID, messages, llmTools)
	if err != nil {
		return nil, err
	}

	var toolCalls []rt.ToolCall
	for _, tc := range resp.Message.ToolCalls {
		toolCalls = append(toolCalls, rt.ToolCall{
			ID:   tc.ID,
			Name: tc.Name,
			Args: tc.Arguments,
		})
	}

	return &rt.LLMResponse{
		Text:      resp.Message.Content,
		ToolCalls: toolCalls,
		Usage: rt.Usage{
			InputTokens:  resp.Usage.PromptTokens,
			OutputTokens: resp.Usage.CompletionTokens,
			TotalTokens:  resp.Usage.TotalTokens,
		},
	}, nil
}

func fromConversationMessage(m sysTypes.AiConversationMessage) Message {
	msg := Message{Role: m.Role, Content: m.Content}

	for _, tc := range m.ToolCalls {
		var args map[string]any
		_ = json.Unmarshal([]byte(tc.Data), &args)
		msg.ToolCalls = append(msg.ToolCalls, ToolCall{
			ID:        tc.CallID,
			Name:      tc.Name,
			Arguments: args,
		})
	}

	if len(m.ToolResults) > 0 {
		msg.Role = "tool"
		msg.ToolCallID = m.ToolResults[0].CallID
		msg.Content = m.ToolResults[0].Data
	}

	return msg
}
