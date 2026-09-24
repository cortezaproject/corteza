package mcp_test

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/require"

	hmcp "github.com/crusttech/human/server/pkg/mcpkit"
)

// docsPage is the public tool reference, relative to the repository root.
const docsPage = "docs/reference/mcp-tools.gen.md"

// toolFamilies are the page's sections, in page order, keyed by the first
// segment of a tool name.
var toolFamilies = []struct{ prefix, title, intro string }{
	{"compose", "Compose", "Namespaces, modules, records, pages, page layouts and charts."},
	{"system", "System", "Users, roles, groups, permissions, applications, agents, chatbots, reminders and other instance-wide settings."},
	{"automation", "Automation", "TAQs, workflows, triggers and the catalogues they are built from."},
	{"discovery", "Discovery", "Full-text search across the instance."},
}

// riskLabels match the webapp's wording (locale tools.risk.*) and pick the
// VitePress badge colour.
var riskLabels = map[hmcp.Risk]struct{ text, badge string }{
	hmcp.RiskRead:        {"Reads", "tip"},
	hmcp.RiskWrite:       {"Writes", "warning"},
	hmcp.RiskDestructive: {"Deletes", "danger"},
}

// TestToolReferenceIsCurrent holds docs/reference/mcp-tools.gen.md to the
// registry. HUMAN_UPDATE_DOCS=1 rewrites the page instead of comparing.
func TestToolReferenceIsCurrent(t *testing.T) {
	got := renderToolReference(t)
	path := filepath.Join(repoRoot(t), filepath.FromSlash(docsPage))

	if os.Getenv("HUMAN_UPDATE_DOCS") == "1" {
		require.NoError(t, os.WriteFile(path, []byte(got), 0o644))
		return
	}

	want, err := os.ReadFile(path)
	require.NoError(t, err, "%s is missing; run `make docs-reference` in server/", docsPage)
	if string(want) != got {
		t.Fatalf("%s is out of date with the MCP tool registry; run `make docs-reference` in server/ and commit the result", docsPage)
	}
}

func renderToolReference(t *testing.T) string {
	t.Helper()

	reg := buildRegistry(t)
	conditional := conditionalTools(t)

	byFamily := map[string][]mcp.Tool{}
	for _, tool := range reg.Tools() {
		if reg.IsHidden(tool.Name) {
			continue
		}
		family, _, _ := strings.Cut(tool.Name, "_")
		byFamily[family] = append(byFamily[family], tool)
	}

	var b strings.Builder
	b.WriteString(`---
title: MCP tools
description: Every tool Human's MCP server offers, with its risk level and parameters.
outline: [2, 2]
---

<!-- This file is auto-generated from the MCP tool registry by server/tests/mcp/docs_test.go. -->

# MCP tools

Human is an MCP (Model Context Protocol) server: an AI client that speaks MCP
can connect to it, sign in as a Human user and call the tools on this page. The
tools are the same operations Human's own [agents](/platform/agents) use. A
tool only does what the signed-in user may do — every call goes through the
same permission checks as the webapp, so a tool that reads records returns only
the records that user can read.

## Endpoint and scope

The MCP endpoint is ` + "`/api/mcp`" + ` on the server's API. Two variants list a
narrower set of tools, for a client that works on one side of the system:

| Endpoint | Lists |
| --- | --- |
| ` + "`/api/mcp`" + ` | Every tool. |
| ` + "`/api/mcp/configuring`" + ` | Tools that change what the system _is_: modules, pages, charts, TAQ and workflow definitions, roles. |
| ` + "`/api/mcp/usage`" + ` | Tools that work with what the system _holds_, or run it: records, executions, chatbots, reports. |

The variants only shorten the list; they are not a permission boundary.

Add ` + "`?maxRisk=read`" + `, ` + "`?maxRisk=write`" + ` or ` + "`?maxRisk=destructive`" + ` to
any of these to cap a session at a risk level. Tools above the cap are left out
of the list, and calling one by name is refused. This is a safety catch against
accidents, such as a session pointed at a production instance deleting a
namespace; permissions are what protect the data.

## Risk levels

Every tool declares one of three risk levels, shown as a badge on each entry
below and in the webapp's tool picker:

| Badge | Level | Meaning |
| --- | --- | --- |
| <Badge type="tip" text="Reads" /> | ` + "`read`" + ` | Makes no change. It can still return sensitive data, which permissions govern. |
| <Badge type="warning" text="Writes" /> | ` + "`write`" + ` | Creates or changes something. Running a TAQ, workflow or agent is always a write. |
| <Badge type="danger" text="Deletes" /> | ` + "`destructive`" + ` | Removes something. Most deletes can be undone with the matching ` + "`undelete`" + ` tool, but a deleted resource is hidden from everything until then. |

## Finding and loading tools

To keep the tool list small, the server lists each tool with a one-sentence
summary and without per-parameter descriptions. Two more tools, which belong to
the MCP server rather than to any part of Human, return the full documentation
shown on this page:

`)

	for _, meta := range metaTools(t) {
		fmt.Fprintf(&b, "- `%s` — %s\n", meta.Name, mdText(meta.Description))
	}

	for _, fam := range toolFamilies {
		tools := byFamily[fam.prefix]
		delete(byFamily, fam.prefix)
		if len(tools) == 0 {
			continue
		}

		fmt.Fprintf(&b, "\n## %s\n\n%s\n", fam.title, fam.intro)
		for _, tool := range tools {
			writeTool(&b, tool, conditional[tool.Name])
		}
	}

	require.Empty(t, byFamily, "tools with a name prefix the page has no section for; add it to toolFamilies")

	return b.String()
}

func writeTool(b *strings.Builder, tool mcp.Tool, when string) {
	risk, ok := riskLabels[hmcp.RiskOf(tool)]
	if !ok {
		risk.text, risk.badge = string(hmcp.RiskOf(tool)), "info"
	}

	groups := make([]string, 0, 2)
	for _, g := range hmcp.GroupsOf(tool) {
		groups = append(groups, "`/api/mcp/"+string(g)+"`")
	}
	sort.Strings(groups)

	fmt.Fprintf(b, "\n### `%s` {#%s}\n\n", tool.Name, tool.Name)
	fmt.Fprintf(b, "<Badge type=\"%s\" text=\"%s\" />", risk.badge, risk.text)
	if len(groups) > 0 {
		fmt.Fprintf(b, " Listed under %s.", strings.Join(groups, " and "))
	}
	b.WriteString("\n\n")

	if when != "" {
		fmt.Fprintf(b, "::: info\nOnly available when %s.\n:::\n\n", envLinks(when))
	}

	b.WriteString(mdText(tool.Description))
	b.WriteString("\n")

	props := tool.InputSchema.Properties
	if len(props) == 0 {
		b.WriteString("\nNo parameters.\n")
		return
	}

	required := map[string]bool{}
	for _, r := range tool.InputSchema.Required {
		required[r] = true
	}

	names := make([]string, 0, len(props))
	for n := range props {
		names = append(names, n)
	}
	sort.Slice(names, func(i, j int) bool {
		if required[names[i]] != required[names[j]] {
			return required[names[i]]
		}
		return names[i] < names[j]
	})

	b.WriteString("\n| Parameter | Type | Required | Description |\n| --- | --- | --- | --- |\n")
	for _, n := range names {
		prop, _ := props[n].(map[string]any)
		typ, _ := prop["type"].(string)
		desc, _ := prop["description"].(string)

		if enum, ok := prop["enum"].([]any); ok && len(enum) > 0 {
			vals := make([]string, 0, len(enum))
			for _, v := range enum {
				vals = append(vals, fmt.Sprintf("`%v`", v))
			}
			desc = strings.TrimSpace(mdText(desc) + " One of " + strings.Join(vals, ", ") + ".")
		} else {
			desc = mdText(desc)
		}

		req := ""
		if required[n] {
			req = "Yes"
		}

		fmt.Fprintf(b, "| `%s` | %s | %s | %s |\n", n, typ, req, cell(desc))
	}
}

var (
	mdSpecial = strings.NewReplacer(
		`\`, `\\`,
		"*", `\*`,
		"<", "&lt;",
		"{{", "&#123;&#123;",
		"`", "\\`",
	)
	trailingSpace = regexp.MustCompile(`[ \t]+\n`)
	envName       = regexp.MustCompile("`([A-Z][A-Z0-9_]+)`")
)

// mdText makes tool prose safe for VitePress: no emphasis or HTML by accident,
// no Vue interpolation, no trailing whitespace.
func mdText(s string) string {
	s = trailingSpace.ReplaceAllString(strings.TrimSpace(s), "\n")
	return mdSpecial.Replace(s)
}

// cell fits text into one table cell.
func cell(s string) string {
	s = strings.ReplaceAll(s, "|", `\|`)
	s = strings.ReplaceAll(s, "\n\n", "<br><br>")
	return strings.ReplaceAll(s, "\n", "<br>")
}

// envLinks links each `ENV_NAME` in s to its entry in the environment reference.
func envLinks(s string) string {
	return envName.ReplaceAllString(s, "[`$1`](/reference/environment#$1)")
}

// metaTools returns human_tool_search and human_tool_load as a client sees
// them. They are added by the MCP server itself rather than the registry, so
// the only way to read them is to ask a server.
func metaTools(t *testing.T) []mcp.Tool {
	t.Helper()

	r := chi.NewRouter()
	r.Route("/api/mcp", hmcp.NewMCPServer(hmcp.NewRegistry(), "human", "docs").MountRoutes)
	srv := httptest.NewServer(r)
	defer srv.Close()

	session := ""
	call := func(method string, params any, id int) json.RawMessage {
		msg := map[string]any{"jsonrpc": "2.0", "method": method}
		if params != nil {
			msg["params"] = params
		}
		if id > 0 {
			msg["id"] = id
		}
		body, err := json.Marshal(msg)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, srv.URL+"/api/mcp?docs=full", bytes.NewReader(body))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json, text/event-stream")
		if session != "" {
			req.Header.Set("Mcp-Session-Id", session)
		}

		res, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer res.Body.Close()
		if s := res.Header.Get("Mcp-Session-Id"); s != "" {
			session = s
		}

		raw, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		require.Lessf(t, res.StatusCode, 300, "%s: %s", method, raw)
		if id == 0 {
			return nil
		}
		return rpcResult(t, res.Header.Get("Content-Type"), raw)
	}

	call("initialize", map[string]any{
		"protocolVersion": mcp.LATEST_PROTOCOL_VERSION,
		"capabilities":    map[string]any{},
		"clientInfo":      map[string]any{"name": "docs", "version": "0"},
	}, 1)
	call("notifications/initialized", nil, 0)

	var listed struct {
		Tools []mcp.Tool `json:"tools"`
	}
	require.NoError(t, json.Unmarshal(call("tools/list", map[string]any{}, 2), &listed))
	require.NotEmpty(t, listed.Tools, "an MCP server over an empty registry listed no meta-tools")

	sort.Slice(listed.Tools, func(i, j int) bool { return listed.Tools[i].Name < listed.Tools[j].Name })
	return listed.Tools
}

// rpcResult extracts the JSON-RPC result from a plain JSON or an SSE response.
func rpcResult(t *testing.T, contentType string, raw []byte) json.RawMessage {
	t.Helper()

	payload := raw
	if strings.HasPrefix(contentType, "text/event-stream") {
		payload = nil
		sc := bufio.NewScanner(bytes.NewReader(raw))
		sc.Buffer(make([]byte, 0, 64*1024), 16*1024*1024)
		for sc.Scan() {
			if data, ok := strings.CutPrefix(sc.Text(), "data:"); ok {
				payload = []byte(strings.TrimSpace(data))
			}
		}
	}

	var env struct {
		Result json.RawMessage `json:"result"`
		Error  *struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	require.NoError(t, json.Unmarshal(payload, &env), "%s", raw)
	require.Nil(t, env.Error, "%s", raw)
	return env.Result
}
