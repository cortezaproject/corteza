package agentic

import (
	hmcp "github.com/crusttech/human/server/pkg/mcpkit"
	"github.com/mark3labs/mcp-go/mcp"
)

// This file holds declarations only. Implementations are in agent_handler.go,
// in the same order.

// An agent's configuration is five nested objects rather than a flat field
// list. Each is one parameter that REPLACES that whole section when sent and
// leaves it untouched when omitted (CONVENTIONS.md §8.2), so changing the
// system prompt means sending the whole `behavior` object, not just the prompt.
// Read the agent first with system_agent_lookup and send back an edited copy.
const agentSectionDoc = "An agent's configuration is five objects — meta, behavior, execution, access " +
	"and invocation. Send only the ones you are changing; each REPLACES that whole section, so " +
	"read the agent first and send an edited copy of the section rather than a fragment of it."

const (
	agentMetaDoc = `JSON object: the agent's labels. {"short":"Support triage","description":"Reads new tickets and files them","sidebarRoles":["<roleID>"]}. ` +
		`'short' is the name a person sees — an agent has no separate name field. 'sidebarRoles' limits who sees it in the webapp sidebar and is not an authorization check.`

	agentBehaviorDoc = `JSON object: what the agent is told and what it may draw on. ` +
		`{"systemPrompt":"You triage support tickets.","guardrails":["no-pii"],"injectSystemContext":true,"knowledgeBases":["<knowledgeBaseID>"],"treatyCLEnabled":false}. ` +
		`Setting treatyCLEnabled true makes the server fill tclArticles with its default set when you send none, and merge the hardwired ones into any list you do send — so reading back a longer list than you sent is expected, not a fault.`

	agentExecutionDoc = `JSON object: which model runs it and how far it may go. ` +
		`{"model":{"llmProviderID":"<id>","model":"claude-sonnet-5","temperature":0.2},"limits":{"maxIterations":10,"timeout":"5m","softLimitRatio":0.8,"contextWindow":200000,"outputTokens":8192}}. ` +
		`'timeout' is a Go duration string. 'contextWindow' is NOT the model's window: it is a cumulative token budget for the whole conversation, checked after each iteration and across turns, and 'softLimitRatio' is the fraction of it at which the agent is told to wrap up. ` +
		`Omitting 'llmProviderID' is resolved when the agent RUNS, not when it is written: the instance's sole active provider is used, and where there are several the run fails naming them. So an agent created without one reads back with an empty 'model' object and is not necessarily broken — system_llm_provider_lookup says which case this instance is in, and is where an ID and a model name come from. ` +
		`Sending a 'temperature' makes the server call the provider to check it, which also validates 'llmProviderID' and the model name — an unrecognised model fails with the provider's own error. Omit temperature and none of that is checked: the model name is stored as given.`

	agentAccessDoc = `JSON object: what the agent may reach. ` +
		`{"allow":[{"namespaceID":"<id>","moduleIDs":[]}],"tools":[{"name":"compose_record_lookup","permission":"always","description":"Read leads","allow":[{"namespaceID":"<id>","moduleIDs":["<id>"]}]}],"taqs":[{"id":"<taqID>","description":"Escalate"}],"workflows":[{"id":"<workflowID>","description":"Notify"}]}. ` +
		`Access is deny-by-default: an agent whose 'tools' is EMPTY can call no tool at all — it still answers from its prompt and knowledge, but cannot look anything up or change anything. ` +
		`Granting a tool does not widen anything — an agent runs as the person who invoked it, so RBAC is the ceiling on every call and the agent can never do what that user could not. ` +
		`'allow' on the access object itself scopes the whole agent to namespaces and modules however its tools were granted; that is the setting to use for "this agent is only for namespace X". ` +
		`Granting an agent system_agent_update lets it re-grant itself any tool, so treat that entry the way you would a permission change. ` +
		`Every name is checked against the registry when you write it, and an unknown one is refused with the near matches: the runtime resolves the allow-list as a whole, so a single typo would stop the agent running at all rather than cost it one tool. ` +
		`An entry's "moduleIDs" narrows it to those modules; leave it EMPTY to mean every module in that namespace, now and in future. Prefer empty unless you actually need to withhold a module: an enumerated list has to be edited on every tool entry each time a module is added, and until it is the agent cannot see the new module and nothing says so. ` +
		`Each entry may carry a "permission": "always" runs the tool unannounced, "ask" stops the run and puts it to the user (the exec call comes back with status "awaiting_approval" and the pending call; send it again with that tool in "approvedTools" to carry on), and "deny" refuses it whatever the scope says. Omit it and the mode follows the tool's risk — reading always, anything that writes asks. ` +
		`An entry names ONE tool via "name", or a whole set via "group" plus "maxRisk" — {"group":"usage","maxRisk":"read","allow":[{"namespaceID":"<id>","moduleIDs":[]}]} grants every read-only data tool in that namespace and picks up tools added later. "group" is "usage" (data and execution) or "configuring" (schema and definitions); "maxRisk" is "read", "write" or "destructive" and defaults to "read". Set one or the other, never both. Prefer a group: naming tools one at a time is what makes an agent tedious to set up and stale afterwards.`

	agentInvocationDoc = `JSON object: who may start it. ` +
		`{"user":{"enabled":true},"system":{"enabled":false,"serviceAccount":"<userID>","inputSchema":{},"outputFormat":"json"}}. ` +
		`'user' is a person invoking it from the webapp; 'system' is another automation doing so under the named service account.`
)

func (h *agentHandler) register() {
	h.reg.RegisterTool(
		mcp.NewTool("system_agent_lookup",
			mcp.WithDescription(
				"Read the AI agents configured on this instance — an agent is a model, a system prompt and "+
					"an allow-list of tools it may call. "+
					"Provide 'agent' to fetch one whole, with all five configuration sections; omit it to "+
					"list them as {agentID, handle, short, status}. "+
					"Call this before system_agent_update: an update replaces whole sections, so you need "+
					"the current one to send back an edited copy. "+
					"An agent is not a chatbot — a chatbot is the widget and the conversation, and it "+
					"delegates to an agent; see system_chatbot_lookup.",
			),
			mcp.WithString("agent", mcp.Description("Agent handle, or ID as a string to prevent precision loss. Omit to list. An agent has no name to resolve by — 'short' is a label, not an identifier.")),
			mcp.WithString("query", mcp.Description("Free-text search across the agents.")),
			mcp.WithString("handle", mcp.Description("Filter the listing to an exact handle.")),
			mcp.WithString("status", mcp.Description(`Filter the listing by status, e.g. "active". New agents are created active.`)),
			mcp.WithBoolean("includeDeleted", mcp.Description("Include soft-deleted agents in the listing. They are excluded by default, which is why undelete needs an ID.")),
			mcp.WithString("limit", mcp.Description("Maximum agents to return when listing, default 50, capped at 200.")),
			mcp.WithString("pageCursor", mcp.Description("Cursor from a previous response, to fetch the next page.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskRead),
		),
		"Lookup agent",
		h.lookup,
	)

	h.reg.RegisterTool(
		mcp.NewTool("system_agent_create",
			mcp.WithDescription(
				"Create an AI agent. "+agentSectionDoc+" "+
					"A new agent is 'active' unless you say otherwise, and with no tools in 'access' it "+
					"answers from its prompt alone, reaching nothing — grant it something in the same "+
					"call. Name a provider and a model in 'execution.model' — system_llm_provider_lookup lists "+
					"both — unless this instance has exactly one LLM provider, which is resolved at run "+
					"time and leaves 'execution.model' empty until then. "+
					"To let people talk to it in a widget, create a chatbot with a conversation scenario "+
					"pointing at this agent — system_chatbot_create.",
			),
			mcp.WithString("handle", mcp.Description("URL-friendly identifier, unique across agents. Use snake_case (lowercase letters, digits, underscores).")),
			mcp.WithString("status", mcp.Description(`Lifecycle status; defaults to "active".`)),
			mcp.WithString("meta", mcp.Description(agentMetaDoc)),
			mcp.WithString("behavior", mcp.Description(agentBehaviorDoc)),
			mcp.WithString("execution", mcp.Description(agentExecutionDoc)),
			mcp.WithString("access", mcp.Description(agentAccessDoc)),
			mcp.WithString("namespace", mcp.Description(
				"Compose namespace this agent is for, by handle, slug or ID. Sets 'access.allow' to that "+
					"namespace and turns on the platform context, so whatever tools you grant are confined to "+
					"it and the agent knows what it is looking at. It does NOT grant any tool: an agent scoped "+
					"here and granted nothing still reaches nothing. Ignored when 'access' already says "+
					"something. It is not stored — what is stored is the scope it produced.")),
			mcp.WithString("invocation", mcp.Description(agentInvocationDoc)),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskWrite),
			hmcp.NeedsFullDocs(),
		),
		"Create agent",
		h.create,
	)

	h.reg.RegisterTool(
		mcp.NewTool("system_agent_update",
			mcp.WithDescription(
				"Change an existing agent. "+agentSectionDoc+" "+
					"An argument you omit is left unchanged, so this is safe to call with one section. "+
					"What is NOT safe is sending a section you built from nothing: 'behavior' without the "+
					"knowledge bases the agent had drops them, because the section replaces rather than "+
					"merges. Call system_agent_lookup with 'agent' first. "+
					"The revision number is the server's and increments on every update; it is not a "+
					"parameter and sending one changes nothing.",
			),
			mcp.WithString("agent", mcp.Required(), mcp.Description("Agent handle, or ID as a string to prevent precision loss.")),
			mcp.WithString("handle", mcp.Description("New handle. Pass an empty string to clear it.")),
			mcp.WithString("status", mcp.Description("New lifecycle status.")),
			mcp.WithString("meta", mcp.Description(agentMetaDoc)),
			mcp.WithString("behavior", mcp.Description(agentBehaviorDoc)),
			mcp.WithString("execution", mcp.Description(agentExecutionDoc)),
			mcp.WithString("access", mcp.Description(agentAccessDoc)),
			mcp.WithString("invocation", mcp.Description(agentInvocationDoc)),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskWrite),
			hmcp.NeedsFullDocs(),
		),
		"Update agent",
		h.update,
	)

	h.reg.RegisterTool(
		mcp.NewTool("system_agent_delete",
			mcp.WithDescription(
				"Delete an agent. The delete is soft: the configuration is retained and "+
					"system_agent_undelete brings it back intact, but the agent stops appearing in "+
					"lookups and stops being invocable. "+
					"Note the agentID before calling — a deleted agent is excluded from every listing, so "+
					"a handle no longer resolves and undelete needs the ID. "+
					"A chatbot scenario pointing at this agent is NOT updated and will fail to find it, so "+
					"check system_chatbot_lookup for scenarios naming it first.",
			),
			mcp.WithString("agent", mcp.Required(), mcp.Description("Agent handle, or ID as a string to prevent precision loss.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskDestructive),
		),
		"Delete agent",
		h.del,
	)

	h.reg.RegisterTool(
		mcp.NewTool("system_agent_undelete",
			mcp.WithDescription(
				"Restore a soft-deleted agent, reversing system_agent_delete. Nothing was erased, so the "+
					"agent returns with its prompt, model, limits and tool allow-list exactly as they were. "+
					"Requires the numeric agentID: a deleted agent is excluded from the default listing, "+
					"so its handle resolves to nothing. Use the ID system_agent_delete reported, or call "+
					"system_agent_lookup with includeDeleted first. "+
					"Calling this on an agent that is not deleted is accepted and changes nothing.",
			),
			mcp.WithString("agentID", mcp.Required(), mcp.Description("ID of the deleted agent, as a string to prevent precision loss. A handle will not work.")),
			hmcp.InGroup(hmcp.GroupConfiguring),
			hmcp.WithRisk(hmcp.RiskWrite),
		),
		"Undelete agent",
		h.undelete,
	)
}
