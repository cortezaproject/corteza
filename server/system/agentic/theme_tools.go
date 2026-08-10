package agentic

import (
	hmcp "github.com/crusttech/human/server/pkg/mcpkit"
	"github.com/mark3labs/mcp-go/mcp"
)

// The colour names a theme carries. Listed in the descriptions because the
// setting is free-form: a caller that invents a name gets no error from the
// store, just a colour that never reaches the stylesheet.
const themeColorDoc = `Colour names: primary, secondary, success, warning, danger, black, white, light, extra-light, body-bg, sidebar-bg, topbar-bg. Values are hex, e.g. "#09344E".`

// This file holds declarations only. Implementations are in theme_handler.go,
// in the same order.

func (h *themeHandler) register() {
	h.reg.RegisterTool(
		mcp.NewTool("system_theme_lookup",
			mcp.WithDescription(
				"Read the webapp's theme colours. Returns each theme — light, dark, and general — "+
					"with its colours as a plain object, which is what "+
					"system_theme_update expects back. "+
					"Call this before updating: an update merges into what is already there, so seeing "+
					"the current palette is how you avoid changing a colour you did not mean to. "+
					themeColorDoc,
			),
			mcp.WithString("theme", mcp.Description(`Which theme to read: "light", "dark", or "general". Omit for all of them.`)),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskRead),
		),
		"Read Theme Colours",
		h.lookup,
	)

	h.reg.RegisterTool(
		mcp.NewTool("system_theme_update",
			mcp.WithDescription(
				"Change the webapp's theme colours. Colours are merged, not replaced: pass only what "+
					"should change and the rest of the palette is left alone. "+
					"This restyles the webapp for EVERY user of this instance, not just the caller, and "+
					"takes effect on their next page load — it is an instance-wide change, so confirm the "+
					"colours with the user before calling. "+
					"Light and dark are separate palettes and neither inherits from the other: a colour "+
					"changed in one is unchanged in the other, and changing only 'light' leaves anyone "+
					"using dark mode seeing the old palette. "+
					themeColorDoc,
			),
			mcp.WithString("theme", mcp.Required(), mcp.Description(`Which theme to change: "light", "dark", or "general".`)),
			mcp.WithString("colors", mcp.Required(), mcp.Description(`JSON object of colour name to hex value, e.g. {"primary":"#09344E","body-bg":"#F4F4F5"}. Only the names given are changed.`)),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskWrite),
		),
		"Change Theme Colours",
		h.update,
	)
}
