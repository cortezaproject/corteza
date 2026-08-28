package agentic

import (
	"context"
	"fmt"
	"strings"

	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/pkg/mcpkit/toolkit"
	sysService "github.com/crusttech/human/server/system/service"
	sysTypes "github.com/crusttech/human/server/system/types"
	"github.com/mark3labs/mcp-go/mcp"
)

// Declarations for this handler are in llm_provider_tools.go.
type llmProviderHandler struct {
	reg toolRegistrar
}

func LlmProviderHandler(reg toolRegistrar) *llmProviderHandler {
	h := &llmProviderHandler{reg: reg}
	h.register()
	return h
}

// llmProviderItem is the compact projection returned by list mode.
type llmProviderItem struct {
	LlmProviderID string `json:"llmProviderID"`
	Handle        string `json:"handle"`
	Name          string `json:"name"`
	Provider      string `json:"provider"`
	Status        string `json:"status"`
	DefaultModel  string `json:"defaultModel"`
}

// lookup lists providers, or fetches one and optionally its model catalogue.
//
// Every path resolves through Search, which is the only method on the LLM
// service that access-checks (CanSearchLlmProviders). FindByID, LookupByHandle
// and ListModels all go straight to the store, so reaching them directly would
// be a tool with no authorization behind it — see CONVENTIONS.md §8.6. Resolving
// through Search first and then calling ListModels keeps the check in front of
// both.
func (h *llmProviderHandler) lookup(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	ref, err := toolkit.Ref(args, "llmProvider")
	if err != nil {
		return nil, err
	}

	wantModels := toolkit.Bool(args, "models")

	if ref == "" {
		if wantModels {
			return nil, fmt.Errorf("llm provider lookup failed: 'models' needs a single provider — name one in 'llmProvider'")
		}
		return h.list(ctx, args)
	}

	provider, err := h.resolve(ctx, ref)
	if err != nil {
		return nil, err
	}

	if !wantModels {
		return toolkit.JSONResult(provider)
	}

	models, err := sysService.DefaultLlmService.ListModels(ctx, provider.ID)
	if err != nil {
		return nil, toolkit.Errf("llm provider model list", err)
	}

	return toolkit.JSONResultWithAny(provider, map[string]any{"models": models})
}

func (h *llmProviderHandler) list(ctx context.Context, args map[string]any) (*mcp.CallToolResult, error) {
	f := sysTypes.LlmProviderFilter{
		Provider: toolkit.Str(args, "provider"),
		Status:   toolkit.Str(args, "status"),
	}

	page := toolkit.Page(args)
	var err error
	if f.Paging, err = filter.NewPaging(page.Limit, page.Cursor); err != nil {
		return nil, fmt.Errorf("invalid pageCursor: %w", err)
	}

	set, out, err := sysService.DefaultLlmService.Search(ctx, f)
	if err != nil {
		return nil, toolkit.Errf("llm provider list", err)
	}

	items := make([]llmProviderItem, 0, len(set))
	for _, p := range set {
		items = append(items, llmProviderItem{
			LlmProviderID: fmt.Sprint(p.ID),
			Handle:        p.Handle,
			Name:          p.Meta.Short,
			Provider:      p.Provider,
			Status:        p.Status,
			DefaultModel:  p.Config.Model,
		})
	}

	return toolkit.JSONResult(map[string]any{
		"llmProviders":   items,
		"nextPageCursor": out.NextPage,
	})
}

// resolve turns a reference — an ID or a handle — into a provider, through the
// access-checked Search.
func (h *llmProviderHandler) resolve(ctx context.Context, ref string) (*sysTypes.LlmProvider, error) {
	f := sysTypes.LlmProviderFilter{}
	if id, err := toolkit.ID(map[string]any{"llmProvider": ref}, "llmProvider"); err == nil && id > 0 {
		f.LlmProviderID = []uint64{id}
	} else {
		f.Handle = ref
	}

	set, _, err := sysService.DefaultLlmService.Search(ctx, f)
	if err != nil {
		return nil, toolkit.Errf("llm provider lookup", err)
	}

	switch len(set) {
	case 0:
		return nil, fmt.Errorf("llm provider %q not found", ref)
	case 1:
		return set[0], nil
	default:
		handles := make([]string, 0, len(set))
		for _, p := range set {
			handles = append(handles, fmt.Sprintf("%d (%s)", p.ID, p.Meta.Short))
		}
		return nil, fmt.Errorf("llm provider %q matches several: %s", ref, strings.Join(handles, ", "))
	}
}
