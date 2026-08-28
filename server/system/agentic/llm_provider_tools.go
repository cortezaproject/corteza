package agentic

import (
	hmcp "github.com/crusttech/human/server/pkg/mcpkit"
	"github.com/mark3labs/mcp-go/mcp"
)

// An LLM provider is the connection to a model vendor. It is read-only here:
// creating one takes an API key, and a tool that accepts a credential is the
// mirror image of one that returns it (CONVENTIONS.md §8.6b). Configuring a
// provider is a person's job in the admin UI; naming one on an agent is not,
// and that is what this exists for.
//
// The API key itself is never in reach: it lives in a separate credential row
// and the provider carries only its ID.
//
// This file holds declarations only. The implementation is in
// llm_provider_handler.go.

func (h *llmProviderHandler) register() {
	h.reg.RegisterTool(
		mcp.NewTool("system_llm_provider_lookup",
			mcp.WithDescription(
				"List the LLM providers configured on this instance, or fetch one. A provider is the "+
					"connection to a model vendor — its type (anthropic, openai, mistral, …), its endpoint "+
					"and its default model — and it is what 'execution.model.llmProviderID' on an agent "+
					"names. "+
					"Read this before system_agent_create or system_agent_update. The ID is not guessable "+
					"and there is no other way to discover one, so an agent gets its provider from here. "+
					"An agent that names none runs on the instance's sole active provider; where there is "+
					"more than one, an unnamed provider is an agent that fails when it is first run rather "+
					"than when it is written — this is where you find out which case you are in. "+
					"Set 'models' to fetch the model names the provider itself advertises, which is where "+
					"'execution.model.model' comes from. That calls the vendor over the network with the "+
					"stored credential, so it is slower than the rest of this tool and fails when the "+
					"credential is wrong — which is itself worth knowing before an agent is built on it. "+
					"'models' needs a single provider named. "+
					"No response here contains an API key. A provider carries only the ID of the "+
					"credential it uses, and nothing on this surface reads a credential. "+
					"Providers are created and their keys set by a person in the admin UI; there is no "+
					"tool for either.",
			),
			mcp.WithString("llmProvider", mcp.Description("Provider ID as a string (to prevent precision loss), or its handle. Omit to list instead.")),
			mcp.WithString("provider", mcp.Description("List only providers of this vendor type, e.g. \"anthropic\", \"openai\", \"mistral\". Ignored when 'llmProvider' is given.")),
			mcp.WithString("status", mcp.Description("List only providers in this state, typically \"active\". A provider that is not active is not eligible to run an agent. Ignored when 'llmProvider' is given.")),
			mcp.WithBoolean("models", mcp.Description("Also return the model names this provider advertises, fetched from the vendor. Requires 'llmProvider'. Slower, and it surfaces a bad credential as an error rather than a silent failure later.")),
			mcp.WithString("limit", mcp.Description("Maximum results, default 50, capped at 200.")),
			mcp.WithString("pageCursor", mcp.Description("Cursor from a previous response, to fetch the next page.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskRead),
		),
		"Lookup LLM providers",
		h.lookup,
	)
}
