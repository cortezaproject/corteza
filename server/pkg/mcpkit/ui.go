package mcpkit

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// MCP Apps (the io.modelcontextprotocol/ui extension, protocol 2026-01-26): a
// tool names a ui:// resource holding an HTML page, and a host that supports
// the extension renders that page in a sandboxed iframe beside the tool's
// result. Hosts without it ignore the metadata and read the text as before.

// UIMimeType is the MIME type of an MCP Apps HTML resource.
const UIMimeType = "text/html;profile=mcp-app"

// MetaUI is the protocol's own _meta key, not namespaced like the human.dev
// tags: hosts look for it by this exact name.
const MetaUI = "ui"

// UICSP lists the origins a UI page may reach beyond itself. Empty means the
// page is self-contained, which is what the host's default policy assumes.
type UICSP struct {
	ConnectDomains  []string `json:"connectDomains,omitempty"`
	ResourceDomains []string `json:"resourceDomains,omitempty"`
	FrameDomains    []string `json:"frameDomains,omitempty"`
}

// UIResource is an MCP Apps page served at a ui:// URI.
type UIResource struct {
	URI         string
	Name        string
	Description string
	HTML        []byte
	CSP         UICSP
}

// WithUI links a tool to the UI resource that renders its result.
//
// A tool carrying it also gets structuredContent on its results; see
// structuredForUI.
func WithUI(uri string) mcp.ToolOption {
	return func(t *mcp.Tool) {
		ensureMeta(t)
		t.Meta.AdditionalFields[MetaUI] = map[string]any{"resourceUri": uri}
	}
}

// UIOf reports the UI resource URI a tool is linked to, or "" if none.
func UIOf(t mcp.Tool) string {
	if t.Meta == nil || t.Meta.AdditionalFields == nil {
		return ""
	}
	ui, ok := t.Meta.AdditionalFields[MetaUI].(map[string]any)
	if !ok {
		return ""
	}
	uri, _ := ui["resourceUri"].(string)
	return uri
}

// RegisterUIResource serves an MCP Apps page. The CSP travels on the read
// contents rather than the listing, which is where hosts look for it.
func (r *Registry) RegisterUIResource(ui UIResource) {
	res := mcp.NewResource(ui.URI, ui.Name,
		mcp.WithResourceDescription(ui.Description),
		mcp.WithMIMEType(UIMimeType),
	)

	r.RegisterResource(res, func(ctx context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      ui.URI,
				MIMEType: UIMimeType,
				Text:     string(ui.HTML),
				Meta:     map[string]any{MetaUI: map[string]any{"csp": ui.CSP}},
			},
		}, nil
	})
}

// structuredForUI copies a UI-linked tool's JSON text result into
// structuredContent, which the host hands to the page and keeps out of the
// model's context. The text content is left exactly as the handler wrote it.
//
// Results that are not a single JSON object stay text-only: structuredContent
// must be an object, and a page has nothing to render from an error string.
func (m *MCPServer) structuredForUI() server.ToolHandlerMiddleware {
	return func(next server.ToolHandlerFunc) server.ToolHandlerFunc {
		return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			res, err := next(ctx, req)
			if err != nil || res == nil || res.IsError || res.StructuredContent != nil {
				return res, err
			}

			t, ok := m.reg.tools[ResolveToolAlias(req.Params.Name)]
			if !ok || UIOf(t.Tool) == "" || len(res.Content) != 1 {
				return res, err
			}

			text, ok := res.Content[0].(mcp.TextContent)
			if !ok || !strings.HasPrefix(strings.TrimSpace(text.Text), "{") || !json.Valid([]byte(text.Text)) {
				return res, err
			}

			res.StructuredContent = json.RawMessage(text.Text)
			return res, err
		}
	}
}
