package agentic

import (
	hmcp "github.com/crusttech/human/server/pkg/mcpkit"
	"github.com/mark3labs/mcp-go/mcp"
)

// This file holds declarations only. Implementations are in skill_handler.go,
// in the same order.

func (h *skillHandler) register() {
	h.reg.RegisterTool(
		mcp.NewTool("system_skill_lookup",
			mcp.WithDescription(
				"Read Human's own guidance on how to use its tools — the rules a tool's parameter list "+
					"cannot state, written per topic rather than per tool. "+
					"Call it with 'tool' BEFORE the first call to an unfamiliar tool: it returns the "+
					"skills that apply to that tool, and skipping it is how a call succeeds and is "+
					"quietly wrong. A page block added without reading page_layout, for instance, is "+
					"stored on the page and placed on no layout, so nothing renders it. "+
					"Pass 'skill' to read one whole, or neither argument to list what exists as "+
					"{name, description, tools}. "+
					"This is guidance, not configuration: nothing here changes the instance, and it is "+
					"the same text for every caller. For a single tool's own full documentation, use "+
					"human_tool_load instead — the two are complementary, and a skill spans the several "+
					"tools a job actually takes.",
			),
			mcp.WithString("skill", mcp.Description("Name of one skill to read whole, e.g. \"page_layout\". Omit to list them.")),
			mcp.WithString("tool", mcp.Description("Tool name, e.g. \"compose_page_update\". Returns every skill that applies to that tool, each with its full text. A tool with no skill returns an empty list, which is not an error.")),
			// Both groups on purpose: this documents the tool surface rather
			// than belonging to one half of it, and a usage-scoped session needs
			// record_handling as much as a configuring one needs page_layout.
			hmcp.InGroup(hmcp.GroupConfiguring, hmcp.GroupUsage),
			hmcp.WithRisk(hmcp.RiskRead),
		),
		"Lookup skill",
		h.lookup,
	)
}
