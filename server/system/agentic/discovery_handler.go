package agentic

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	a "github.com/cortezaproject/corteza/server/pkg/auth"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type (
	toolRegistrar interface {
		RegisterTool(tool mcp.Tool, title string, handler server.ToolHandlerFunc)
		RegisterToolWithAvailability(tool mcp.Tool, title string, handler server.ToolHandlerFunc, available func() bool)
	}

	discoveryTokenSigner interface {
		Sign(opt ...a.IssueOptFn) ([]byte, error)
	}

	discoveryHandler struct {
		reg     toolRegistrar
		baseURL string
		signer  discoveryTokenSigner
	}

	discoveryResponse struct {
		Response struct {
			Hits []struct {
				Type  string `json:"type"`
				Value struct {
					RecordID     string            `json:"recordID"`
					CustomValues map[string]string `json:"customValues"`
					Values       []struct {
						Name  string   `json:"name"`
						Label string   `json:"label"`
						Value []string `json:"value"`
					} `json:"values"`
					Module struct {
						Handle string `json:"handle"`
						Name   string `json:"name"`
					} `json:"module"`
					Namespace struct {
						Name string `json:"name"`
					} `json:"namespace"`
					Created struct {
						At string `json:"at"`
						By string `json:"by"`
					} `json:"created"`
				} `json:"value"`
			} `json:"hits"`
			TotalHits    int `json:"total_hits"`
			TotalResults int `json:"total_results"`
		} `json:"response"`
	}
)

func DiscoveryHandler(reg toolRegistrar, baseURL string, signer discoveryTokenSigner) *discoveryHandler {
	h := &discoveryHandler{reg: reg, baseURL: baseURL, signer: signer}
	if baseURL != "" {
		h.register()
	}
	return h
}

func (h *discoveryHandler) isAvailable() bool {
	resp, err := http.Get(h.baseURL + "/healthcheck")
	if err != nil {
		return false
	}
	resp.Body.Close()
	return resp.StatusCode < 500
}

func (h *discoveryHandler) register() {
	h.reg.RegisterToolWithAvailability(
		mcp.NewTool("discovery_search",
			mcp.WithDescription("Full-text search across all indexed records. Use this only when the user explicitly asks to search or find existing records by keyword. Do not use this when creating records or when the user provides a value directly — use it only to look up existing data."),
			mcp.WithString("query", mcp.Required(), mcp.Description("Natural language or keyword search query")),
			mcp.WithString("size", mcp.Description("Number of results to return (default: 10)")),
			mcp.WithString("namespace", mcp.Description("Filter results to a specific namespace slug")),
			mcp.WithString("module", mcp.Description("Filter results to a specific module handle")),
		),
		"Discover Records",
		h.search,
		h.isAvailable,
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

	formatted, err := formatDiscoveryResponse(body)
	if err != nil {
		// Fall back to raw JSON if parsing fails
		return mcp.NewToolResultText(string(body)), nil
	}

	return mcp.NewToolResultText(formatted), nil
}

func formatDiscoveryResponse(body []byte) (string, error) {
	var dr discoveryResponse
	if err := json.Unmarshal(body, &dr); err != nil {
		return "", err
	}

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("found %d result(s) (showing %d):\n\n",
		dr.Response.TotalResults, len(dr.Response.Hits)))

	for i, hit := range dr.Response.Hits {
		v := hit.Value
		sb.WriteString(fmt.Sprintf("--- Result %d ---\n", i+1))
		sb.WriteString(fmt.Sprintf("Record ID : %s\n", v.RecordID))
		sb.WriteString(fmt.Sprintf("Module    : %s\n", v.Module.Name))
		sb.WriteString(fmt.Sprintf("Namespace : %s\n", v.Namespace.Name))
		sb.WriteString(fmt.Sprintf("Created   : %s by %s\n", v.Created.At, v.Created.By))

		if len(v.Values) > 0 {
			sb.WriteString("Fields    :\n")
			for _, field := range v.Values {
				sb.WriteString(fmt.Sprintf("  %s: %s\n", field.Label, strings.Join(field.Value, ", ")))
			}
		}
		sb.WriteString("\n")
	}

	return sb.String(), nil
}
