package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	rt "github.com/cortezaproject/corteza/server/system/agentic/runtime"
	"github.com/cortezaproject/corteza/server/pkg/id"
	"github.com/cortezaproject/corteza/server/store"
	sysTypes "github.com/cortezaproject/corteza/server/system/types"
)

type (
	Service struct {
		store               store.Storer
		ac                  llmProviderAccessController
		anthropicAPIVersion string
	}

	llmProviderAccessController interface {
		CanCreateLlmProvider(ctx context.Context) bool
		CanSearchLlmProviders(ctx context.Context) bool
		CanReadLlmProvider(ctx context.Context, p *sysTypes.LlmProvider) bool
		CanUpdateLlmProvider(ctx context.Context, p *sysTypes.LlmProvider) bool
		CanDeleteLlmProvider(ctx context.Context, p *sysTypes.LlmProvider) bool
	}
)

func New(s store.Storer, ac llmProviderAccessController, anthropicAPIVersion string) (*Service, error) {
	return &Service{store: s, ac: ac, anthropicAPIVersion: anthropicAPIVersion}, nil
}

func (svc *Service) Create(ctx context.Context, p *sysTypes.LlmProvider, apiKey string) (*sysTypes.LlmProvider, error) {
	if !svc.ac.CanCreateLlmProvider(ctx) {
		return nil, fmt.Errorf("not allowed to create LLM providers")
	}

	p.ID = id.Next()
	p.CreatedAt = time.Now().Round(time.Second)

	var cred *sysTypes.Credential
	if apiKey != "" {
		cred = &sysTypes.Credential{
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

	if cred != nil {
		if _, err := fetchModels(ctx, p, cred, svc.anthropicAPIVersion); err != nil {
			p.Status = "unauthorized"
		} else {
			p.Status = "active"
		}
	} else {
		p.Status = "unauthorized"
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

func (svc *Service) Update(ctx context.Context, upd *sysTypes.LlmProvider, apiKey string) (*sysTypes.LlmProvider, error) {
	existing, err := store.LookupLlmProviderByID(ctx, svc.store, upd.ID)
	if err != nil {
		return nil, err
	}

	if !svc.ac.CanUpdateLlmProvider(ctx, existing) {
		return nil, fmt.Errorf("not allowed to update LLM provider")
	}

	existing.Handle = upd.Handle
	existing.Provider = upd.Provider
	existing.Meta = upd.Meta
	existing.Config = upd.Config

	var cred *sysTypes.Credential
	if apiKey != "" {
		if existing.CredentialID == 0 {
			cred = &sysTypes.Credential{
				ID:          id.Next(),
				Kind:        "api-key",
				Label:       existing.Handle + " API Key",
				Credentials: apiKey,
				CreatedAt:   existing.CreatedAt,
			}
			if err := store.CreateCredential(ctx, svc.store, cred); err != nil {
				return nil, fmt.Errorf("could not create credential: %w", err)
			}
			existing.CredentialID = cred.ID
		} else {
			var err error
			cred, err = store.LookupCredentialByID(ctx, svc.store, existing.CredentialID)
			if err != nil {
				return nil, fmt.Errorf("could not load credential: %w", err)
			}
			cred.Credentials = apiKey
			if err := store.UpdateCredential(ctx, svc.store, cred); err != nil {
				return nil, fmt.Errorf("could not update credential: %w", err)
			}
		}
	} else if existing.CredentialID != 0 {
		var err error
		cred, err = store.LookupCredentialByID(ctx, svc.store, existing.CredentialID)
		if err != nil {
			return nil, fmt.Errorf("could not load credential: %w", err)
		}
	}

	if cred != nil {
		if _, err := fetchModels(ctx, existing, cred, svc.anthropicAPIVersion); err != nil {
			existing.Status = "unauthorized"
		} else {
			existing.Status = "active"
		}
	} else {
		existing.Status = "unauthorized"
	}

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

	existing, err := store.LookupLlmProviderByID(ctx, svc.store, providerID)
	if err != nil {
		return err
	}

	if !svc.ac.CanDeleteLlmProvider(ctx, existing) {
		return fmt.Errorf("not allowed to delete LLM provider")
	}

	return store.DeleteLlmProviderByID(ctx, svc.store, providerID)
}

func (svc *Service) Search(ctx context.Context, f sysTypes.LlmProviderFilter) (sysTypes.LlmProviderSet, error) {
	if !svc.ac.CanSearchLlmProviders(ctx) {
		return nil, fmt.Errorf("not allowed to search LLM providers")
	}

	set, _, err := store.SearchLlmProviders(ctx, svc.store, f)
	return set, err
}

func (svc *Service) Validate(ctx context.Context, providerID uint64) error {
	provider, err := store.LookupLlmProviderByID(ctx, svc.store, providerID)
	if err != nil {
		return fmt.Errorf("could not resolve LLM provider: %w", err)
	}

	if provider.CredentialID == 0 {
		provider.Status = "unauthorized"
		now := time.Now().Round(time.Second)
		provider.UpdatedAt = &now
		_ = store.UpdateLlmProvider(ctx, svc.store, provider)
		return fmt.Errorf("no credential configured for this provider")
	}

	cred, err := store.LookupCredentialByID(ctx, svc.store, provider.CredentialID)
	if err != nil {
		return fmt.Errorf("could not resolve credential: %w", err)
	}

	_, checkErr := fetchModels(ctx, provider, cred, svc.anthropicAPIVersion)

	if checkErr != nil {
		provider.Status = "unauthorized"
	} else {
		provider.Status = "active"
	}

	now := time.Now().Round(time.Second)
	provider.UpdatedAt = &now
	if storeErr := store.UpdateLlmProvider(ctx, svc.store, provider); storeErr != nil {
		return fmt.Errorf("could not update provider status: %w", storeErr)
	}

	return checkErr
}

func (svc *Service) ListModels(ctx context.Context, providerID uint64) ([]string, error) {
	provider, err := store.LookupLlmProviderByID(ctx, svc.store, providerID)
	if err != nil {
		return nil, fmt.Errorf("could not resolve LLM provider: %w", err)
	}

	if provider.CredentialID == 0 {
		return nil, fmt.Errorf("no credential configured for this provider")
	}

	cred, err := store.LookupCredentialByID(ctx, svc.store, provider.CredentialID)
	if err != nil {
		return nil, fmt.Errorf("could not resolve credential for LLM provider: %w", err)
	}

	return fetchModels(ctx, provider, cred, svc.anthropicAPIVersion)
}

func fetchModels(ctx context.Context, provider *sysTypes.LlmProvider, cred *sysTypes.Credential, anthropicAPIVersion string) ([]string, error) {
	switch provider.Provider {
	case "anthropic":
		return fetchAnthropicModels(ctx, provider, cred, anthropicAPIVersion)
	default:
		return fetchOpenAIModels(ctx, provider, cred)
	}
}

// Prompt resolves the provider and its credential, then forwards the conversation to the LLM.
// If model is non-empty it overrides the provider's configured model without changing the DB record.
func (svc *Service) Prompt(ctx context.Context, providerID uint64, model string, maxTokens int, messages []Message, tools []Tool) (*Response, error) {
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

	return svc.callProvider(ctx, provider, cred, model, maxTokens, messages, tools)
}

func (svc *Service) callProvider(ctx context.Context, provider *sysTypes.LlmProvider, cred *sysTypes.Credential, model string, maxTokens int, messages []Message, tools []Tool) (*Response, error) {
	switch provider.Provider {
	case "anthropic":
		return promptAnthropic(ctx, provider, cred, model, maxTokens, messages, tools, svc.anthropicAPIVersion)
	default:
		return promptOpenAI(ctx, provider, cred, model, maxTokens, messages, tools)
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

	resp, err := svc.Prompt(ctx, config.ProviderID, config.Model, config.MaxTokens, messages, llmTools)
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
