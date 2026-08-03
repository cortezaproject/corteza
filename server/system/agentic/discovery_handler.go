package agentic

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	cmpService "github.com/crusttech/human/server/compose/service"
	a "github.com/crusttech/human/server/pkg/auth"
	hmcp "github.com/crusttech/human/server/system/agentic/mcp"
	"github.com/crusttech/human/server/system/agentic/toolkit"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Declarations for these handlers are in discovery_tools.go, in the same order.
type (
	// toolRegistrar is the registry seam every handler in this package takes.
	// It lives here rather than in a file of its own because reminder_handler.go
	// and the rest of the package depend on it being declared in this file.
	toolRegistrar interface {
		RegisterTool(tool mcp.Tool, title string, handler server.ToolHandlerFunc, opts ...hmcp.RegisterOption)
	}

	discoveryTokenSigner interface {
		Sign(opt ...a.IssueOptFn) ([]byte, error)
	}

	discoveryHandler struct {
		reg     toolRegistrar
		baseURL string
		signer  discoveryTokenSigner
	}

	// discoveryResponse is the upstream discovery service's payload. Only the
	// fields projected into discoveryItem below reach the caller.
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

	// discoveryItem is the projection returned for each hit.
	discoveryItem struct {
		RecordID  string               `json:"recordID"`
		Module    string               `json:"module"`
		Namespace string               `json:"namespace"`
		CreatedAt string               `json:"createdAt"`
		CreatedBy string               `json:"createdBy"`
		Fields    []discoveryItemField `json:"fields,omitempty"`
	}

	discoveryItemField struct {
		Name  string   `json:"name"`
		Label string   `json:"label"`
		Value []string `json:"value"`
	}
)

// DiscoveryHandler registers the tool only when a discovery service is
// configured; without a base URL there is nothing to search.
func DiscoveryHandler(reg toolRegistrar, baseURL string, signer discoveryTokenSigner) *discoveryHandler {
	h := &discoveryHandler{reg: reg, baseURL: baseURL, signer: signer}
	if baseURL != "" {
		h.register()
	}
	return h
}

// resolveScope turns the declared namespace/module args into the ID filters the
// discovery service expects. Both are optional here even though `namespace` is
// declared Required: the in-process path satisfies the requirement upstream, and
// failing a direct caller that omitted it would be a behaviour change beyond
// making the params work at all.
//
// RBAC does the real work — FindByAny runs as the caller, so a namespace they
// cannot read resolves to an error rather than a silent widening.
func (h *discoveryHandler) resolveScope(ctx context.Context, args map[string]any) (namespaceIDs, moduleIDs []string, err error) {
	nsRef, err := toolkit.Ref(args, "namespace")
	if err != nil {
		return nil, nil, err
	}
	if nsRef == "" {
		return nil, nil, nil
	}

	ns, err := cmpService.DefaultNamespace.FindByAny(ctx, nsRef)
	if err != nil {
		return nil, nil, toolkit.Errf("namespace lookup", err)
	}
	namespaceIDs = []string{strconv.FormatUint(ns.ID, 10)}

	modRef, err := toolkit.Ref(args, "module")
	if err != nil {
		return nil, nil, err
	}
	if modRef == "" {
		return namespaceIDs, nil, nil
	}

	mod, err := cmpService.DefaultModule.FindByAny(ctx, ns.ID, modRef)
	if err != nil {
		return nil, nil, toolkit.Errf("module lookup", err)
	}
	return namespaceIDs, []string{strconv.FormatUint(mod.ID, 10)}, nil
}

func (h *discoveryHandler) search(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	query := toolkit.Str(args, "query")
	size := toolkit.Str(args, "size")

	// namespaceIDs/moduleIDs are injected by the in-process executor after it
	// has checked the agent's allow-list against the declared namespace/module
	// args. They are not declared params.
	namespaceIDs := extractStringSlice(args["namespaceIDs"])
	moduleIDs := extractStringSlice(args["moduleIDs"])

	// Nothing injected means there is no executor — a remote MCP client called
	// this directly. Resolve the declared args ourselves, otherwise they are
	// silently ignored and the search runs across everything the caller's
	// discovery token permits, which is not what the schema promises.
	//
	// The injected IDs deliberately win when present: in-process they are the
	// authorization narrowing, already checked against the agent's allow-list.
	if len(namespaceIDs) == 0 && len(moduleIDs) == 0 {
		nsIDs, modIDs, err := h.resolveScope(ctx, args)
		if err != nil {
			return nil, err
		}
		namespaceIDs, moduleIDs = nsIDs, modIDs
	}

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
		return nil, toolkit.Errf("discovery token signing", err)
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
		return nil, toolkit.Errf("discovery request build", err)
	}
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", string(token)))

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, toolkit.Errf("discovery search", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, toolkit.Errf("discovery response read", err)
	}

	if resp.StatusCode >= 400 {
		return nil, toolkit.Errf("discovery search", fmt.Errorf("status %d: %s", resp.StatusCode, string(body)))
	}

	var dr discoveryResponse
	if err := json.Unmarshal(body, &dr); err != nil {
		// Fall back to the raw payload rather than failing the call — the
		// caller can still read an unexpected shape.
		return toolkit.JSONResult(map[string]any{"raw": string(body)})
	}

	items := make([]discoveryItem, 0, len(dr.Response.Hits))
	for _, hit := range dr.Response.Hits {
		v := hit.Value
		item := discoveryItem{
			RecordID:  v.RecordID,
			Module:    v.Module.Name,
			Namespace: v.Namespace.Name,
			CreatedAt: v.Created.At,
			CreatedBy: v.Created.By,
		}
		for _, field := range v.Values {
			item.Fields = append(item.Fields, discoveryItemField{
				Name:  field.Name,
				Label: field.Label,
				Value: field.Value,
			})
		}
		items = append(items, item)
	}

	return toolkit.JSONResult(map[string]any{
		"records":      items,
		"totalResults": dr.Response.TotalResults,
	})
}

// isAvailable backs hmcp.Available: the discovery service is a separate
// process, so the tool is only advertised while it answers its healthcheck.
func (h *discoveryHandler) isAvailable() bool {
	resp, err := http.Get(h.baseURL + "/healthcheck")
	if err != nil {
		return false
	}
	resp.Body.Close()
	return resp.StatusCode < 500
}

// extractStringSlice reads an injected ID list, which arrives as []string from
// the in-process executor and as []interface{} when it has been through JSON.
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
