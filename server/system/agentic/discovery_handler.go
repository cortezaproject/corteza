package agentic

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	a "github.com/cortezaproject/corteza/server/pkg/auth"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type (
	toolRegistrar interface {
		RegisterTool(tool mcp.Tool, title string, handler server.ToolHandlerFunc)
	}

	discoveryTokenSigner interface {
		Sign(opt ...a.IssueOptFn) ([]byte, error)
	}

	discoveryHandler struct {
		reg     toolRegistrar
		baseURL string
		signer  discoveryTokenSigner
	}
)

func DiscoveryHandler(reg toolRegistrar, baseURL string, signer discoveryTokenSigner) *discoveryHandler {
	h := &discoveryHandler{reg: reg, baseURL: baseURL, signer: signer}
	h.register()
	return h
}

func (h *discoveryHandler) register() {
	h.reg.RegisterTool(
		mcp.NewTool("discovery_search",
			mcp.WithDescription("Use this tool first when the user asks any question about data, records, or information. It searches across all indexed records and returns relevant results. Prefer this over compose_record_lookup unless you already know the exact record ID or module."),
			mcp.WithString("query", mcp.Required(), mcp.Description("Natural language or keyword search query")),
			mcp.WithString("size", mcp.Description("Number of results to return (default: 10)")),
			mcp.WithString("namespace", mcp.Description("Filter results to a specific namespace slug")),
			mcp.WithString("module", mcp.Description("Filter results to a specific module handle")),
		),
		"Discover Records",
		h.search,
	)
}

func (h *discoveryHandler) search(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid request")
	}

	query, _ := args["query"].(string)
	size, _ := args["size"].(string)
	namespace, _ := args["namespace"].(string)
	module, _ := args["module"].(string)

	if size == "" {
		size = "10"
	}

	identity := a.GetIdentityFromContext(ctx)

	token, err := h.signer.Sign(
		a.WithIdentity(identity),
		a.WithScope("api", "discovery"),
		a.WithExpiration(5*time.Minute),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to sign discovery token: %w", err)
	}

	params := url.Values{}
	params.Set("q", query)
	params.Set("resourceTypes", "compose:record")
	params.Set("size", size)
	if namespace != "" {
		params.Set("namespaceAggs", namespace)
	}
	if module != "" {
		params.Set("moduleAggs", module)
	}

	endpoint := fmt.Sprintf("%s/?%s", h.baseURL, params.Encode())

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build discovery request: %w", err)
	}
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", string(token)))

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("discovery search failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read discovery response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("discovery search error %d: %s", resp.StatusCode, string(body))
	}

	return mcp.NewToolResultText(string(body)), nil
}
