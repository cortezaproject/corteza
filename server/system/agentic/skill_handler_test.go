package agentic

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/crusttech/human/server/system/agentic/skills"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/require"
)

func callSkillLookup(t *testing.T, args map[string]any) map[string]any {
	t.Helper()

	lib, err := skills.LoadLibrary()
	require.NoError(t, err)

	h := &skillHandler{reg: &stubRegistrar{}, lib: lib}

	req := mcp.CallToolRequest{}
	req.Params.Arguments = args

	res, err := h.lookup(context.Background(), req)
	require.NoError(t, err)

	text, ok := mcp.AsTextContent(res.Content[0])
	require.True(t, ok)

	var out map[string]any
	require.NoError(t, json.Unmarshal([]byte(text.Text), &out))
	return out
}

func TestSkillLookupListsTheLibrary(t *testing.T) {
	out := callSkillLookup(t, map[string]any{})

	list, ok := out["skills"].([]any)
	require.True(t, ok)
	require.NotEmpty(t, list)

	// A listing is name and description only — the bodies are what makes the
	// library expensive, and a caller listing has not chosen one yet.
	first := list[0].(map[string]any)
	require.NotEmpty(t, first["name"])
	require.NotContains(t, first, "body")
}

func TestSkillLookupForAToolReturnsWholeSkills(t *testing.T) {
	out := callSkillLookup(t, map[string]any{"tool": "compose_page_update"})

	list, ok := out["skills"].([]any)
	require.True(t, ok)
	require.NotEmpty(t, list, "compose_page_update has skills; none came back")

	var names []string
	for _, raw := range list {
		s := raw.(map[string]any)
		names = append(names, s["name"].(string))
		require.NotEmpty(t, s["body"], "a skill fetched for a tool must carry its text")
	}
	require.Contains(t, names, "page_layout")
}

func TestSkillLookupReadsOneWhole(t *testing.T) {
	out := callSkillLookup(t, map[string]any{"skill": "page_layout"})

	require.Equal(t, "page_layout", out["name"])
	// Matched against the body with its line breaks collapsed: these are
	// sentences, and where they happen to wrap is not part of the claim.
	body := strings.Join(strings.Fields(out["body"].(string)), " ")

	// The two facts an agent gets wrong without this file: which object places
	// a block, and that an xywh it sends is honoured rather than ignored.
	require.Contains(t, body, "A block the layout does not name is never rendered")
	require.Contains(t, body, "`xywh` is honoured")
	require.Contains(t, body, "48 columns")
}

func TestSkillLookupRejectsAnUnknownSkill(t *testing.T) {
	lib, err := skills.LoadLibrary()
	require.NoError(t, err)

	h := &skillHandler{reg: &stubRegistrar{}, lib: lib}
	req := mcp.CallToolRequest{}
	req.Params.Arguments = map[string]any{"skill": "nope"}

	_, err = h.lookup(context.Background(), req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "no such skill")
}

type announcingRegistrar struct {
	stubRegistrar
	said []string
}

func (a *announcingRegistrar) AddInstructions(text string) { a.said = append(a.said, text) }

// Every announced skill reaches the initialize text by name, and nothing else
// does; a registrar with no instructions to carry is left alone.
func TestSkillHandlerAnnouncesSkills(t *testing.T) {
	lib, err := skills.LoadLibrary()
	require.NoError(t, err)

	reg := &announcingRegistrar{}
	SkillHandler(reg, lib)

	announced := 0
	for _, s := range lib.All() {
		if s.Announce != "" {
			announced++
		}
	}
	require.Len(t, reg.said, announced)
	require.Contains(t, strings.Join(reg.said, "\n"), `system_skill_lookup {"skill": "custom_app"}`)

	require.NotPanics(t, func() { SkillHandler(&stubRegistrar{}, lib) })
}
