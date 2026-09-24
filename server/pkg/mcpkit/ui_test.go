package mcpkit

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func uiTestServer(t *testing.T) *MCPServer {
	t.Helper()

	reg := NewRegistry()
	result := func(text string) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return mcp.NewToolResultText(text), nil
		}
	}

	reg.RegisterTool(mcp.NewTool("ui_tool", InGroup(GroupUsage), WithRisk(RiskRead), WithUI("ui://test/view")),
		"UI tool", result(`{"records":[1,2]}`))
	reg.RegisterTool(mcp.NewTool("ui_array", InGroup(GroupUsage), WithRisk(RiskRead), WithUI("ui://test/view")),
		"UI array", result(`[1,2]`))
	reg.RegisterTool(mcp.NewTool("plain_tool", InGroup(GroupUsage), WithRisk(RiskRead)),
		"Plain tool", result(`{"records":[1,2]}`))
	reg.RegisterTool(mcp.NewTool("ui_view", InGroup(GroupUsage), WithRisk(RiskRead), WithUI("ui://test/view")),
		"UI view data", func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return WithViewData(mcp.NewToolResultText(`{"records":[]}`), map[string]any{"fields": []string{"title"}}), nil
		})

	reg.RegisterUIResource(UIResource{
		URI:  "ui://test/view",
		Name: "View",
		HTML: []byte("<p>hi</p>"),
		CSP:  UICSP{ConnectDomains: []string{"https://api.example.tld"}},
	})

	return NewMCPServer(reg, "test", "0")
}

// rpc sends one JSON-RPC request through the server and returns its result
// decoded as generic JSON, the way a host reads it.
func rpc(t *testing.T, m *MCPServer, method string, params any) map[string]any {
	t.Helper()

	req, err := json.Marshal(map[string]any{"jsonrpc": "2.0", "id": 1, "method": method, "params": params})
	require.NoError(t, err)

	out, err := json.Marshal(m.server.HandleMessage(context.Background(), req))
	require.NoError(t, err)

	var resp struct {
		Result map[string]any `json:"result"`
		Error  any            `json:"error"`
	}
	require.NoError(t, json.Unmarshal(out, &resp))
	require.Nilf(t, resp.Error, "%s failed: %s", method, out)
	return resp.Result
}

func TestUIToolListingCarriesResourceURI(t *testing.T) {
	m := uiTestServer(t)

	meta := map[string]any{}
	for _, raw := range rpc(t, m, "tools/list", map[string]any{})["tools"].([]any) {
		tool := raw.(map[string]any)
		tm, _ := tool["_meta"].(map[string]any)
		meta[tool["name"].(string)] = tm["ui"]
	}

	assert.Equal(t, map[string]any{"resourceUri": "ui://test/view"}, meta["ui_tool"])
	assert.Nil(t, meta["plain_tool"])
}

func TestUIResourceReadCarriesMimeAndCSP(t *testing.T) {
	m := uiTestServer(t)

	contents := rpc(t, m, "resources/read", map[string]any{"uri": "ui://test/view"})["contents"].([]any)
	require.Len(t, contents, 1)

	c := contents[0].(map[string]any)
	assert.Equal(t, UIMimeType, c["mimeType"])
	assert.Equal(t, "<p>hi</p>", c["text"])
	assert.Equal(t,
		map[string]any{"csp": map[string]any{"connectDomains": []any{"https://api.example.tld"}}},
		c["_meta"].(map[string]any)["ui"])
}

func TestUIToolResultCarriesStructuredContent(t *testing.T) {
	m := uiTestServer(t)
	call := func(name string) map[string]any {
		return rpc(t, m, "tools/call", map[string]any{"name": name, "arguments": map[string]any{}})
	}

	res := call("ui_tool")
	assert.Equal(t, map[string]any{"records": []any{1.0, 2.0}}, res["structuredContent"])
	assert.Equal(t, `{"records":[1,2]}`, res["content"].([]any)[0].(map[string]any)["text"],
		"the text the model reads is left as the handler wrote it")

	assert.NotContains(t, call("plain_tool"), "structuredContent", "a tool without a UI gets text only")
	assert.NotContains(t, call("ui_array"), "structuredContent", "structuredContent must be an object")
}

func TestViewDataTravelsInResultMeta(t *testing.T) {
	m := uiTestServer(t)

	res := rpc(t, m, "tools/call", map[string]any{"name": "ui_view", "arguments": map[string]any{}})

	assert.Equal(t,
		map[string]any{MetaViewData: map[string]any{"fields": []any{"title"}}},
		res["_meta"])
	assert.Equal(t, `{"records":[]}`, res["content"].([]any)[0].(map[string]any)["text"],
		"view data never reaches the text the model reads")
}
