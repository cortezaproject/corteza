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

	a "github.com/crusttech/human/server/pkg/auth"
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
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithDescription("Search or list records you have access to. Use this when the user asks to find, list, or show existing records. Only namespaces listed in your DISCOVERY ACCESS section are permitted — the executor will deny any other namespace. Leave query empty to list all accessible records, or provide a specific field value (name, email, phone) to search within them."),
			mcp.WithString("namespace", mcp.Required(), mcp.Description("The exact namespace name from your DISCOVERY ACCESS section. Do not use module names here.")),
			mcp.WithString("module", mcp.Description("The exact module name from your DISCOVERY ACCESS section. Leave empty to search across all accessible modules in the namespace.")),
			mcp.WithString("query", mcp.Description("A specific value to search for inside record fields (e.g. a name, email, phone). Leave empty to list all records.")),
			mcp.WithString("size", mcp.Description("Number of results to return (default: 10)")),
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
	namespaceIDs := extractStringSlice(args["namespaceIDs"])
	moduleIDs := extractStringSlice(args["moduleIDs"])

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
	for _, id := range namespaceIDs {
		params.Add("namespaceID", id)
	}
	for _, id := range moduleIDs {
		params.Add("moduleID", id)
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

func extractStringSlice(v any) []string {
	switch ids := v.(type) {
	case []string:
		return ids
	case []interface{}:
		out := make([]string, 0, len(ids))
		for _, id := range ids {
			if s, ok := id.(string); ok {
				out = append(out, s)
			}
		}
		return out
	}
	return nil
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
		sb.WriteString(fmt.Sprintf("Record ID : \"%s\"\n", v.RecordID))
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
