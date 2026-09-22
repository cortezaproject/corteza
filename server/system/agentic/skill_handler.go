package agentic

import (
	"context"
	"fmt"
	"strings"

	"github.com/crusttech/human/server/pkg/mcpkit/toolkit"
	"github.com/crusttech/human/server/system/agentic/skills"
	"github.com/mark3labs/mcp-go/mcp"
)

// Declarations for these handlers are in skill_tools.go, in the same order.
type skillHandler struct {
	reg toolRegistrar
	lib skills.Registry
}

// SkillHandler exposes the skill library the agentic runtime already reads.
//
// The runtime injects a skill AFTER a tool has run, which corrects the next
// call and not the one that was wrong. Over HTTP nothing injected anything at
// all, so a client outside Human — the case this tool exists for — had no way
// to reach the library. A tool is the way in that works on both surfaces,
// because it is only text in a result.
func SkillHandler(reg toolRegistrar, lib skills.Registry) *skillHandler {
	h := &skillHandler{reg: reg, lib: lib}
	h.register()
	h.announce()
	return h
}

// instructionAdder is a registrar that also carries initialize instructions.
type instructionAdder interface {
	AddInstructions(text string)
}

// announce puts each skill that must be read before its tools are called into
// the text a client gets at initialize. A skill found only through its trigger
// tools arrives after the work it governs is already done — a custom app is
// written long before system_application_source_set is called.
func (h *skillHandler) announce() {
	adder, ok := h.reg.(instructionAdder)
	if !ok || h.lib == nil {
		return
	}

	for _, s := range h.lib.All() {
		if s.Announce == "" {
			continue
		}
		adder.AddInstructions(fmt.Sprintf(
			"%s Read it first with system_skill_lookup {\"skill\": %q}.", s.Announce, s.Name,
		))
	}
}

type skillItem struct {
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Tools       []string `json:"tools,omitempty"`
	Body        string   `json:"body,omitempty"`
}

func (h *skillHandler) lookup(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, err := toolkit.Args(req)
	if err != nil {
		return nil, err
	}

	if h.lib == nil {
		return nil, fmt.Errorf("skill library is not loaded on this instance")
	}

	if name := strings.TrimSpace(toolkit.Str(args, "skill")); name != "" {
		for _, s := range h.lib.All() {
			if s.Name == name {
				return toolkit.JSONResult(whole(s))
			}
		}
		return nil, fmt.Errorf("no such skill %q — call this tool with no arguments to list them", name)
	}

	// Skills for one tool come back whole: the caller is about to make that
	// call, so a name and a summary would only cost it a second round trip.
	if tool := strings.TrimSpace(toolkit.Str(args, "tool")); tool != "" {
		hits := h.lib.ForTool(tool)
		out := make([]skillItem, 0, len(hits))
		for _, s := range hits {
			out = append(out, whole(s))
		}
		return toolkit.JSONResult(map[string]any{"tool": tool, "skills": out})
	}

	all := h.lib.All()
	out := make([]skillItem, 0, len(all))
	for _, s := range all {
		out = append(out, skillItem{Name: s.Name, Description: s.Description, Tools: s.Triggers})
	}

	return toolkit.JSONResult(map[string]any{"skills": out})
}

func whole(s *skills.Skill) skillItem {
	return skillItem{Name: s.Name, Description: s.Description, Tools: s.Triggers, Body: s.Body}
}
